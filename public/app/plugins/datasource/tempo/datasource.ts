import { groupBy } from 'lodash';
import {
  DataQuery,
  DataQueryRequest,
  DataQueryResponse,
  DataSourceApi,
  DataSourceInstanceSettings,
  DataSourceJsonData,
  LoadingState,
} from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { TraceToLogsOptions } from 'app/core/components/TraceToLogsSettings';
import { getDatasourceSrv } from 'app/features/plugins/datasource_srv';
import { from, merge, Observable, throwError } from 'rxjs';
import { map, mergeMap, toArray } from 'rxjs/operators';
import { LokiOptions, LokiQuery } from '../loki/types';
import { transformTrace, transformTraceList } from './resultTransformer';
import { PrometheusDatasource } from '../prometheus/datasource';
import { PromQuery } from '../prometheus/types';
import { mapPromMetricsToServiceMap, serviceMapMetrics } from './graphTransform';

export type TempoQueryType = 'search' | 'traceId' | 'serviceMap';

export interface TempoJsonData extends DataSourceJsonData {
  tracesToLogs?: TraceToLogsOptions;
  serviceMap?: {
    datasourceUid?: string;
  };
}

export type TempoQuery = {
  query: string;
  // Query to find list of traces, e.g., via Loki
  linkedQuery?: LokiQuery;
  queryType: TempoQueryType;
} & DataQuery;

export class TempoDatasource extends DataSourceWithBackend<TempoQuery, TempoJsonData> {
  tracesToLogs?: TraceToLogsOptions;
  serviceMap?: {
    datasourceUid?: string;
  };

  constructor(instanceSettings: DataSourceInstanceSettings<TempoJsonData>) {
    super(instanceSettings);
    this.tracesToLogs = instanceSettings.jsonData.tracesToLogs;
    this.serviceMap = instanceSettings.jsonData.serviceMap;
  }

  query(options: DataQueryRequest<TempoQuery>): Observable<DataQueryResponse> {
    const subQueries: Array<Observable<DataQueryResponse>> = [];
    const filteredTargets = options.targets.filter((target) => !target.hide);
    const targets: { [type: string]: TempoQuery[] } = groupBy(filteredTargets, (t) => t.queryType || 'traceId');

    // Run search queries on linked datasource
    if (this.tracesToLogs?.datasourceUid && targets.search?.length > 0) {
      const dsSrv = getDatasourceSrv();
      subQueries.push(
        from(dsSrv.get(this.tracesToLogs.datasourceUid)).pipe(
          mergeMap((linkedDatasource: DataSourceApi) => {
            // Wrap linked query into a data request based on original request
            const linkedRequest: DataQueryRequest = { ...options, targets: targets.search.map((t) => t.linkedQuery!) };
            // Find trace matchers in derived fields of the linked datasource that's identical to this datasource
            const settings: DataSourceInstanceSettings<LokiOptions> = (linkedDatasource as any).instanceSettings;
            const traceLinkMatcher: string[] =
              settings.jsonData.derivedFields
                ?.filter((field) => field.datasourceUid === this.uid && field.matcherRegex)
                .map((field) => field.matcherRegex) || [];
            if (!traceLinkMatcher || traceLinkMatcher.length === 0) {
              return throwError(
                'No Loki datasource configured for search. Set up Derived Fields for traces in a Loki datasource settings and link it to this Tempo datasource.'
              );
            } else {
              return (linkedDatasource.query(linkedRequest) as Observable<DataQueryResponse>).pipe(
                map((response) =>
                  response.error ? response : transformTraceList(response, this.uid, this.name, traceLinkMatcher)
                )
              );
            }
          })
        )
      );
    }

    if (this.serviceMap?.datasourceUid && targets.serviceMap?.length > 0) {
      subQueries.push(this.serviceMapQuery(options));
    }

    if (targets.traceId?.length > 0) {
      const traceRequest: DataQueryRequest<TempoQuery> = { ...options, targets: targets.traceId };
      subQueries.push(
        super.query(traceRequest).pipe(
          map((response) => {
            if (response.error) {
              return response;
            }
            return transformTrace(response);
          })
        )
      );
    }

    return merge(...subQueries);
  }

  async testDatasource(): Promise<any> {
    // to test Tempo we send a dummy traceID and verify Tempo answers with 'trace not found'
    const response = await super.query({ targets: [{ query: '0' }] } as any).toPromise();

    const errorMessage = response.error?.message;
    if (
      errorMessage &&
      errorMessage.startsWith('failed to get trace') &&
      errorMessage.endsWith('trace not found in Tempo')
    ) {
      return { status: 'success', message: 'Data source is working' };
    }

    return { status: 'error', message: 'Data source is not working' + (errorMessage ? `: ${errorMessage}` : '') };
  }

  getQueryDisplayText(query: TempoQuery) {
    return query.query;
  }

  private queryServiceMapPrometheus(request: DataQueryRequest<PromQuery>) {
    return from(getDatasourceSrv().get(this.serviceMap!.datasourceUid)).pipe(
      mergeMap((ds) => {
        return (ds as PrometheusDatasource).query(request);
      })
    );
  }

  private serviceMapQuery(request: DataQueryRequest<TempoQuery>) {
    return this.queryServiceMapPrometheus(makePromServiceMapRequest(request)).pipe(
      // Just collect all the responses first before processing into node graph data
      toArray(),
      map((responses: DataQueryResponse[]) => {
        return {
          data: mapPromMetricsToServiceMap(request, responses),
          state: LoadingState.Done,
        };
      })
    );
  }
}

function makePromServiceMapRequest(options: DataQueryRequest<TempoQuery>): DataQueryRequest<PromQuery> {
  return {
    ...options,
    targets: serviceMapMetrics.map((metric) => {
      return {
        refId: metric,
        expr: `delta(${metric}[$__range])`,
        instant: true,
      };
    }),
  };
}
