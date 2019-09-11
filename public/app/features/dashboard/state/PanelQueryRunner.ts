// Libraries
import { cloneDeep } from 'lodash';
import { ReplaySubject, Unsubscribable, Observable } from 'rxjs';
import { map } from 'rxjs/operators';

// Services & Utils
import { config } from 'app/core/config';
import { getDatasourceSrv } from 'app/features/plugins/datasource_srv';
import kbn from 'app/core/utils/kbn';
import templateSrv from 'app/features/templating/template_srv';
import { runRequest, postProcessPanelData } from './runRequest';
import { runSharedRequest, isSharedDashboardQuery } from '../../../plugins/datasource/dashboard';

// Types
import {
  PanelData,
  PanelDataFormat,
  DataQuery,
  DataQueryRequest,
  DataSourceApi,
  DataSourceJsonData,
} from '@grafana/ui';
import { TimeRange, DataTransformerConfig, transformDataFrame, ScopedVars } from '@grafana/data';

export interface QueryRunnerOptions<
  TQuery extends DataQuery = DataQuery,
  TOptions extends DataSourceJsonData = DataSourceJsonData
> {
  datasource: string | DataSourceApi<TQuery, TOptions>;
  queries: TQuery[];
  panelId: number;
  dashboardId?: number;
  timezone?: string;
  timeRange: TimeRange;
  timeInfo?: string; // String description of time range for display
  widthPixels: number;
  maxDataPoints: number | undefined | null;
  minInterval: string | undefined | null;
  scopedVars?: ScopedVars;
  cacheTimeout?: string;
  delayStateNotification?: number; // default 100ms.
  transformations?: DataTransformerConfig[];
}

let counter = 100;
function getNextRequestId() {
  return 'Q' + counter++;
}

export class PanelQueryRunner {
  private subject?: ReplaySubject<PanelData>;
  private subscription?: Unsubscribable;
  private transformations?: DataTransformerConfig[];

  constructor() {
    this.subject = new ReplaySubject(1);
  }

  /**
   * Returns an observable that subscribes to the shared multi-cast subject (that reply last result).
   * Here the caller can also control if transformations should be applied or not
   */
  getData(format: PanelDataFormat = PanelDataFormat.Frames, applyTransforms = true): Observable<PanelData> {
    const ensureFormats = postProcessPanelData(format);

    return this.subject.pipe(
      map((data: PanelData) => {
        const transformedData = ensureFormats(data);

        if (applyTransforms && this.hasTransformations()) {
          const newSeries = transformDataFrame(this.transformations, transformedData.series);
          return { ...transformedData, series: newSeries };
        }

        return transformedData;
      })
    );
  }

  hasTransformations() {
    return config.featureToggles.transformations && this.transformations && this.transformations.length > 0;
  }

  /**
   * Useful when chaining PanelQueryRunners (for shared query results feature)
   * To avoid double post-processing & transformations
   */
  getDataRaw(): Observable<PanelData> {
    return this.subject.pipe();
  }

  async run(options: QueryRunnerOptions) {
    const {
      queries,
      timezone,
      datasource,
      panelId,
      dashboardId,
      timeRange,
      timeInfo,
      cacheTimeout,
      widthPixels,
      maxDataPoints,
      scopedVars,
      minInterval,
      // delayStateNotification,
    } = options;

    if (isSharedDashboardQuery(datasource)) {
      this.pipeToSubject(runSharedRequest(options));
      return;
    }

    const request: DataQueryRequest = {
      requestId: getNextRequestId(),
      timezone,
      panelId,
      dashboardId,
      range: timeRange,
      timeInfo,
      interval: '',
      intervalMs: 0,
      targets: cloneDeep(queries),
      maxDataPoints: maxDataPoints || widthPixels,
      scopedVars: scopedVars || {},
      cacheTimeout,
      startTime: Date.now(),
    };

    // Add deprecated property
    (request as any).rangeRaw = timeRange.raw;

    try {
      const ds = await getDataSource(datasource, request.scopedVars);

      if (ds.meta && !ds.meta.hiddenQueries) {
        request.targets = request.targets.filter(q => !q.hide);
      }

      // Attach the datasource name to each query
      request.targets = request.targets.map(query => {
        if (!query.datasource) {
          query.datasource = ds.name;
        }
        return query;
      });

      const lowerIntervalLimit = minInterval ? templateSrv.replace(minInterval, request.scopedVars) : ds.interval;
      const norm = kbn.calculateInterval(timeRange, widthPixels, lowerIntervalLimit);

      // make shallow copy of scoped vars,
      // and add built in variables interval and interval_ms
      request.scopedVars = Object.assign({}, request.scopedVars, {
        __interval: { text: norm.interval, value: norm.interval },
        __interval_ms: { text: norm.intervalMs.toString(), value: norm.intervalMs },
      });

      request.interval = norm.interval;
      request.intervalMs = norm.intervalMs;

      this.pipeToSubject(runRequest(ds, request));
    } catch (err) {
      console.log('PanelQueryRunner Error', err);
    }
  }

  private pipeToSubject(observable: Observable<PanelData>) {
    if (this.subscription) {
      this.subscription.unsubscribe();
    }

    this.subscription = observable.subscribe({
      next: (data: PanelData) => {
        this.subject.next(data);
      },
    });
  }

  setTransformations(transformations?: DataTransformerConfig[]) {
    this.transformations = transformations;
  }

  /**
   * Called when the panel is closed
   */
  destroy() {
    // Tell anyone listening that we are done
    if (this.subject) {
      this.subject.complete();
    }

    if (this.subscription) {
      this.subscription.unsubscribe();
    }
  }
}

async function getDataSource(
  datasource: string | DataSourceApi | null,
  scopedVars: ScopedVars
): Promise<DataSourceApi> {
  if (datasource && (datasource as any).query) {
    return datasource as DataSourceApi;
  }
  return await getDatasourceSrv().get(datasource as string, scopedVars);
}
