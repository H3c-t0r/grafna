import cloneDeep from 'lodash/cloneDeep';
import groupBy from 'lodash/groupBy';
import map from 'lodash/map';
import flatten from 'lodash/flatten';

import {
  DataSourceApi,
  DataQuery,
  DataQueryRequest,
  DataQueryResponse,
  DataStreamState,
  DataStreamObserver,
  DataSourceInstanceSettings,
} from '@grafana/ui';
import { LoadingState } from '@grafana/data';

import { getDatasourceSrv } from 'app/features/plugins/datasource_srv';
import { getProcessedDataFrames } from 'app/features/dashboard/state/PanelQueryState';
import { toDataQueryError } from 'app/features/dashboard/state/PanelQueryState';
import { getBackendSrv } from 'app/core/services/backend_srv';

export const MIXED_DATASOURCE_NAME = '-- Mixed --';

export class MixedDatasource extends DataSourceApi<DataQuery> {
  constructor(instanceSettings: DataSourceInstanceSettings) {
    super(instanceSettings);
  }

  async query(request: DataQueryRequest<DataQuery>, observer: DataStreamObserver): Promise<DataQueryResponse> {
    // Remove any hidden or invalid queries
    const queries = request.targets.filter(t => {
      return !t.hide || t.datasource === MIXED_DATASOURCE_NAME;
    });

    if (!queries.length) {
      return Promise.resolve({ data: [] }); // nothing
    }

    // When multiple datasources are used, merge their results
    // after all queries have finished.  (<=6.3 behavior)
    const sets = groupBy(queries, 'datasource');
    const names = Object.keys(sets);
    if (names.length > 1) {
      const promises = map(sets, (targets: DataQuery[]) => {
        const dsName = targets[0].datasource;
        return getDatasourceSrv()
          .get(dsName)
          .then(ds => {
            const opt = cloneDeep(request);
            opt.targets = targets;
            return ds.query(opt);
          });
      });

      return Promise.all(promises).then(results => {
        return { data: flatten(map(results, 'data')) };
      });
    }

    // If there is only one datasource used, make multiple requests
    // and stream the results.  (new in 6.4)
    const ds = await getDatasourceSrv().get(queries[0].datasource);
    request.subRequests = [];
    for (const query of queries) {
      const sub = cloneDeep(request);
      sub.requestId = request.requestId + '_' + request.subRequests.length;
      sub.targets = [query];
      sub.startTime = Date.now();
      request.subRequests.push(sub);
      this.startStreamingQuery(
        ds,
        query.refId, // Replace existing data by refId
        sub, // the modified sub-request
        observer
      );
    }
    return { data: [] }; // maybe wait for first result?
  }

  startStreamingQuery(
    datasource: DataSourceApi,
    key: string,
    request: DataQueryRequest<DataQuery>,
    observer: DataStreamObserver
  ) {
    const event: DataStreamState = {
      key,
      state: LoadingState.Loading,
      request,
      unsubscribe: () => {
        console.log('Cancel async query', request);
        getBackendSrv().resolveCancelerIfExists(request.requestId);
      },
    };
    // Send an initial 'loading' state -- necessary?
    observer(event);

    // Starts background process
    datasource
      .query(request, observer)
      .then(res => {
        request.endTime = Date.now();
        event.state = LoadingState.Done;
        event.data = getProcessedDataFrames(res.data).map(series => {
          if (!series.meta) {
            series.meta = {};
          }
          series.meta.requestId = request.requestId;
          return series;
        });
        observer(event);
      })
      .catch(err => {
        request.endTime = Date.now();
        event.state = LoadingState.Error;
        event.error = toDataQueryError(err);
        observer(event);
      });
  }

  testDatasource() {
    return Promise.resolve({});
  }
}
