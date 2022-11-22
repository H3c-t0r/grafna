import { dateMath, getTimeZone, isDateTime, TimeRange, TimeZone, toUtc } from '@grafana/data';

import { SceneObjectUrlSyncConfig } from '../services/SceneObjectUrlSyncConfig';

import { SceneObjectBase } from './SceneObjectBase';
import { SceneTimeRangeLike, SceneTimeRangeState } from './types';

export class SceneTimeRange extends SceneObjectBase<SceneTimeRangeState> implements SceneTimeRangeLike {
  protected _urlSync = new SceneObjectUrlSyncConfig({
    keys: ['from', 'to'],
    getUrlState: () => this.getUrlState(),
    updateFromUrl: (values) => this.updateFromUrl(values),
  });

  public constructor(state: Partial<SceneTimeRangeState> = {}) {
    const from = state.from ?? 'now-6h';
    const to = state.to ?? 'now';
    const timeZone = state.timeZone ?? getTimeZone();
    const value = evaluateTimeRange(state.from ?? 'now-6h', state.to ?? 'now', state.timeZone ?? getTimeZone());
    super({ from, to, timeZone, value, ...state });
  }

  public onTimeRangeChange = (timeRange: TimeRange) => {
    const update: Partial<SceneTimeRangeState> = {};

    if (typeof timeRange.raw.from === 'string') {
      update.from = timeRange.raw.from;
    } else {
      update.from = timeRange.raw.from.toISOString();
    }

    if (typeof timeRange.raw.to === 'string') {
      update.to = timeRange.raw.to;
    } else {
      update.to = timeRange.raw.to.toISOString();
    }

    update.value = evaluateTimeRange(update.from, update.to, this.state.timeZone);

    this.setState(update);
  };

  public onRefresh = () => {
    this.setState({ value: evaluateTimeRange(this.state.from, this.state.to, this.state.timeZone) });
  };

  public onIntervalChanged = (_: string) => {};

  private getUrlState() {
    const value = this.state.value;
    let from = this.state.from;
    let to = this.state.from;

    if (isDateTime(value.raw.from)) {
      from = value.from.toISOString();
    }

    if (isDateTime(value.raw.to)) {
      to = value.raw.to.toISOString();
    }

    return new Map<string, string>([
      ['from', from],
      ['to', to],
    ]);
  }

  private updateFromUrl(values: Map<string, string>) {
    const update: Partial<SceneTimeRangeState> = {};

    const from = parseUrlParam(values.get('from'));
    if (from) {
      update.from = from;
    }

    const to = parseUrlParam(values.get('to'));
    if (to) {
      update.to = to;
    }

    this.setState(update);
  }
}

function parseUrlParam(value: string | undefined): string | null {
  if (typeof value !== 'string') {
    return null;
  }

  if (value.indexOf('now') !== -1) {
    return value;
  }

  if (value.length === 8) {
    const utcValue = toUtc(value, 'YYYYMMDD');
    if (utcValue.isValid()) {
      return utcValue.toISOString();
    }
  } else if (value.length === 15) {
    const utcValue = toUtc(value, 'YYYYMMDDTHHmmss');
    if (utcValue.isValid()) {
      return utcValue.toISOString();
    }
  }

  const epoch = parseInt(value, 10);
  if (!isNaN(epoch)) {
    return toUtc(epoch).toISOString();
  }

  return null;
}

function evaluateTimeRange(from: string, to: string, timeZone: TimeZone, fiscalYearStartMonth?: number): TimeRange {
  return {
    from: dateMath.parse(from, false, timeZone, fiscalYearStartMonth)!,
    to: dateMath.parse(to, true, timeZone, fiscalYearStartMonth)!,
    raw: {
      from: from,
      to: to,
    },
  };
}
