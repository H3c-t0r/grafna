import { toUtc, dateTime } from '@grafana/data';
export const intervalMap = {
    Hourly: { startOf: 'hour', amount: 'hours' },
    Daily: { startOf: 'day', amount: 'days' },
    Weekly: { startOf: 'isoWeek', amount: 'weeks' },
    Monthly: { startOf: 'month', amount: 'months' },
    Yearly: { startOf: 'year', amount: 'years' },
};
export class IndexPattern {
    constructor(pattern, interval) {
        this.pattern = pattern;
        this.interval = interval;
        this.dateLocale = 'en';
    }
    getIndexForToday() {
        if (this.interval) {
            return toUtc().locale(this.dateLocale).format(this.pattern);
        }
        else {
            return this.pattern;
        }
    }
    getIndexList(from, to) {
        // When no `from` or `to` is provided, we request data from 7 subsequent/previous indices
        // for the provided index pattern.
        // This is useful when requesting log context where the only time data we have is the log
        // timestamp.
        // TODO: Remove when enableBackendMigration toggle is removed
        const indexOffset = 7;
        if (!this.interval) {
            return this.pattern;
        }
        const intervalInfo = intervalMap[this.interval];
        const start = dateTime(from || dateTime(to).add(-indexOffset, intervalInfo.amount))
            .utc()
            .startOf(intervalInfo.startOf);
        const endEpoch = dateTime(to || dateTime(from).add(indexOffset, intervalInfo.amount))
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
//# sourceMappingURL=IndexPattern.js.map