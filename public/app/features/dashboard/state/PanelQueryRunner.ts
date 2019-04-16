import { getDatasourceSrv } from 'app/features/plugins/datasource_srv';
import { Subject, Unsubscribable, PartialObserver } from 'rxjs';
import {
  guessFieldTypes,
  toSeriesData,
  PanelData,
  LoadingState,
  DataQuery,
  TimeRange,
  ScopedVars,
  DataRequestInfo,
  SeriesData,
  DataQueryError,
  toLegacyResponseData,
  isSeriesData,
  DataSourceApi,
} from '@grafana/ui';

import kbn from 'app/core/utils/kbn';

/**
 * Query response may come as events. any null|missing value just means
 * that is unknown info in the request
 */
export interface PanelDataEvent extends Partial<PanelData> {
  eventId: number;
}

export interface QueryRunnerOptions {
  ds?: DataSourceApi; // if they already have the datasource, don't look it up
  datasource: string | null;
  queries: DataQuery[];
  panelId: number;
  dashboardId?: number;
  timezone?: string;
  timeRange?: TimeRange;
  widthPixels: number;
  minInterval?: string;
  maxDataPoints?: number;
  scopedVars?: ScopedVars;
  cacheTimeout?: string;
  delayStateNotification?: number; // default 100ms.
}

export enum PanelQueryRunnerFormat {
  series = 'series',
  legacy = 'legacy',
}

export class PanelQueryRunner {
  private counter = 0;
  private subject?: Subject<PanelDataEvent>;

  private sendSeries = false;
  private sendLegacy = false;

  subscribe(observer: PartialObserver<PanelDataEvent>, format = PanelQueryRunnerFormat.series): Unsubscribable {
    if (!this.subject) {
      this.subject = new Subject();
    }
    if (format === PanelQueryRunnerFormat.legacy) {
      this.sendLegacy = true;
    } else {
      this.sendSeries = true;
    }
    return this.subject.subscribe(observer);
  }

  async run(options: QueryRunnerOptions) {
    if (!this.subject || !this.subject.observers.length) {
      return; // Don't run if nobody cares!
    }

    const {
      queries,
      timezone,
      datasource,
      panelId,
      dashboardId,
      timeRange,
      cacheTimeout,
      widthPixels,
      maxDataPoints,
      scopedVars,
      delayStateNotification,
    } = options;

    const request: DataRequestInfo = {
      timezone,
      panelId,
      dashboardId,
      range: timeRange,
      rangeRaw: timeRange.raw,
      interval: '',
      intervalMs: 0,
      targets: queries,
      maxDataPoints: maxDataPoints || widthPixels,
      scopedVars: scopedVars || {},
      cacheTimeout,
      startTime: Date.now(),
    };

    if (!queries) {
      this.subject.next({
        eventId: this.counter++,
        state: LoadingState.Done,
        series: [], // Clear the data
        legacy: [],
        request,
      });
      return;
    }

    try {
      const ds = options.ds ? options.ds : await getDatasourceSrv().get(datasource, request.scopedVars);

      const minInterval = options.minInterval || ds.interval;
      const norm = kbn.calculateInterval(timeRange, widthPixels, minInterval);

      // make shallow copy of scoped vars,
      // and add built in variables interval and interval_ms
      request.scopedVars = Object.assign({}, request.scopedVars, {
        __interval: { text: norm.interval, value: norm.interval },
        __interval_ms: { text: norm.intervalMs, value: norm.intervalMs },
      });
      request.interval = norm.interval;
      request.intervalMs = norm.intervalMs;

      // Send a loading status event on slower queries
      setTimeout(() => {
        if (!request.endTime) {
          this.subject.next({
            eventId: this.counter++,
            state: LoadingState.Loading,
            request,
          });
        }
      }, delayStateNotification || 100);

      const resp = await ds.query(request);
      request.endTime = Date.now();

      // Make sure the response is in a supported format
      const series = this.sendSeries ? getProcessedSeriesData(resp.data) : [];
      const legacy = this.sendLegacy
        ? resp.data.map(v => {
            if (isSeriesData(v)) {
              return toLegacyResponseData(v);
            }
            return v;
          })
        : undefined;

      // Notify results
      this.subject.next({
        eventId: this.counter++,
        state: LoadingState.Done,
        request,
        series,
        legacy,
      });
    } catch (err) {
      request.endTime = Date.now();

      const error = err as DataQueryError;
      if (!error.message) {
        err.message = 'Query Error';
      }

      this.subject.next({
        eventId: this.counter++,
        state: LoadingState.Error,
        error: error,
        request,
        // ?? Should an error clear the last data ???
      });
    }
  }
}

/**
 * All panels will be passed tables that have our best guess at colum type set
 *
 * This is also used by PanelChrome for snapshot support
 */
export function getProcessedSeriesData(results?: any[]): SeriesData[] {
  if (!results) {
    return [];
  }

  const series: SeriesData[] = [];
  for (const r of results) {
    if (r) {
      series.push(guessFieldTypes(toSeriesData(r)));
    }
  }
  return series;
}
