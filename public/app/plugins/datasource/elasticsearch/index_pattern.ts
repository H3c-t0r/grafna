import { toUtc, dateTime, DateTime, DurationUnit } from '@grafana/data';
import { Interval } from './types';

type IntervalMap = Record<
  Interval,
  {
    startOf: DurationUnit;
    amount: DurationUnit;
  }
>;

const intervalMap: IntervalMap = {
  Hourly: { startOf: 'hour', amount: 'hours' },
  Daily: { startOf: 'day', amount: 'days' },
  Weekly: { startOf: 'week', amount: 'weeks' },
  Monthly: { startOf: 'month', amount: 'months' },
  Yearly: { startOf: 'year', amount: 'years' },
};

export class IndexPattern {
  private dateLocale = 'en';

  constructor(private pattern: string, private interval?: keyof typeof intervalMap) {}

  getIndexForToday() {
    if (this.interval) {
      return toUtc()
        .locale(this.dateLocale)
        .format(this.pattern);
    } else {
      return this.pattern;
    }
  }

  getIndexList(from?: DateTime, to?: DateTime) {
    if (!this.interval) {
      return this.pattern;
    }

    const intervalInfo = intervalMap[this.interval];

    if (!from) {
      from = dateTime(to).add(-7, intervalInfo.amount);
    }
    if (!to) {
      to = dateTime(from).add(7, intervalInfo.amount);
    }

    const start = dateTime(from)
      .utc()
      .startOf(intervalInfo.startOf);
    const endEpoch = dateTime(to)
      .utc()
      .startOf(intervalInfo.startOf)
      .valueOf();
    const indexList = [];

    while (start.valueOf() <= endEpoch) {
      indexList.push(start.locale(this.dateLocale).format(this.pattern));
      start.add(1, intervalInfo.amount);
    }

    return indexList;
  }
}
