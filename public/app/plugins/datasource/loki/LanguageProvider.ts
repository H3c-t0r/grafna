import { flatten } from 'lodash';
import { LRUCache } from 'lru-cache';

import { LanguageProvider, AbstractQuery, KeyValue, getDefaultTimeRange, TimeRange } from '@grafana/data';
import { config } from '@grafana/runtime';

import { DEFAULT_MAX_LINES_SAMPLE, LokiDatasource } from './datasource';
import { abstractQueryToExpr, mapAbstractOperatorsToOp, processLabels } from './languageUtils';
import { getStreamSelectorsFromQuery } from './queryUtils';
import { buildVisualQueryFromString } from './querybuilder/parsing';
import {
  extractLabelKeysFromDataFrame,
  extractLogParserFromDataFrame,
  extractUnwrapLabelKeysFromDataFrame,
} from './responseUtils';
import { ParserAndLabelKeysResult, LokiQuery, LokiQueryType, LabelType } from './types';

const NS_IN_MS = 1000000;

export default class LokiLanguageProvider extends LanguageProvider {
  labelKeys: string[];
  started = false;
  startedTimeRange?: TimeRange;
  datasource: LokiDatasource;

  /**
   *  Cache for labels of series. This is bit simplistic in the sense that it just counts responses each as a 1 and does
   *  not account for different size of a response. If that is needed a `length` function can be added in the options.
   *  10 as a max size is totally arbitrary right now.
   */
  private seriesCache = new LRUCache<string, Record<string, string[]>>({ max: 10 });
  private labelsCache = new LRUCache<string, string[]>({ max: 10 });

  constructor(datasource: LokiDatasource, initialValues?: any) {
    super();

    this.datasource = datasource;
    this.labelKeys = [];

    Object.assign(this, initialValues);
  }

  request = async (url: string, params?: any) => {
    try {
      return await this.datasource.metadataRequest(url, params);
    } catch (error) {
      console.error(error);
    }

    return undefined;
  };

  /**
   * Initialize the language provider by fetching set of labels.
   */
  start = (timeRange?: TimeRange) => {
    const range = timeRange ?? this.getDefaultTimeRange();
    // refetch labels if either there's not already a start task or the time range has changed
    if (
      !this.startTask ||
      this.startedTimeRange?.from.isSame(range.from) === false ||
      this.startedTimeRange?.to.isSame(range.to) === false
    ) {
      this.startedTimeRange = range;
      this.startTask = this.fetchLabels({ timeRange: range }).then(() => {
        this.started = true;
        return [];
      });
    }

    return this.startTask;
  };

  /**
   * Returns the label keys that have been fetched.
   * If labels have not been fetched yet, it will return an empty array.
   * For updated labels (which should not happen often), use fetchLabels.
   * @todo It is quite complicated to know when to use fetchLabels and when to use getLabelKeys.
   * We should consider simplifying this and use caching in the same way as with seriesCache and labelsCache
   * and just always use fetchLabels.
   * Caching should be thought out properly, so we are not fetching this often, as labelKeys should not be changing often.
   *
   * @returns {string[]} An array of label keys or an empty array if labels have not been fetched.
   */
  getLabelKeys(): string[] {
    return this.labelKeys;
  }

  importFromAbstractQuery(labelBasedQuery: AbstractQuery): LokiQuery {
    return {
      refId: labelBasedQuery.refId,
      expr: abstractQueryToExpr(labelBasedQuery),
      queryType: LokiQueryType.Range,
    };
  }

  exportToAbstractQuery(query: LokiQuery): AbstractQuery {
    if (!query.expr || query.expr.length === 0) {
      return { refId: query.refId, labelMatchers: [] };
    }
    const streamSelectors = getStreamSelectorsFromQuery(query.expr);

    const labelMatchers = streamSelectors.map((streamSelector) => {
      const visualQuery = buildVisualQueryFromString(streamSelector).query;
      const matchers = visualQuery.labels.map((label) => {
        return {
          name: label.label,
          value: label.value,
          operator: mapAbstractOperatorsToOp[label.op],
        };
      });
      return matchers;
    });

    return {
      refId: query.refId,
      labelMatchers: flatten(labelMatchers),
    };
  }

  /**
   * Fetch all label keys
   * This asynchronous function returns all available label keys from the data source.
   * It returns a promise that resolves to an array of strings containing the label keys.
   *
   * @param options - (Optional) An object containing additional options - currently only time range.
   * @param options.timeRange - (Optional) The time range for which you want to retrieve label keys. If not provided, the default time range is used.
   * @returns A promise containing an array of label keys.
   * @throws An error if the fetch operation fails.
   */
  async fetchLabels(options?: { timeRange?: TimeRange }): Promise<string[]> {
    const url = 'labels';
    const range = options?.timeRange ?? this.getDefaultTimeRange();
    const timeRange = this.datasource.getTimeRangeParams(range);

    const res = await this.request(url, timeRange);
    if (Array.isArray(res)) {
      const labels = res
        .slice()
        .sort()
        .filter((label) => label !== '__name__');
      this.labelKeys = labels;
      return this.labelKeys;
    }

    return [];
  }

  /**
   * Fetch series labels for a selector
   *
   * This method fetches labels for a given stream selector, such as `{job="grafana"}`.
   * It returns a promise that resolves to a record mapping label names to their corresponding values.
   *
   * @param streamSelector - The stream selector for which you want to retrieve labels.
   * @param options - (Optional) An object containing additional options - currently only time range.
   * @param options.timeRange - (Optional) The time range for which you want to retrieve label keys. If not provided, the default time range is used.
   * @returns A promise containing a record of label names and their values.
   * @throws An error if the fetch operation fails.
   */
  fetchSeriesLabels = async (
    streamSelector: string,
    options?: { timeRange?: TimeRange }
  ): Promise<Record<string, string[]>> => {
    const interpolatedMatch = this.datasource.interpolateString(streamSelector);
    const url = 'series';
    const range = options?.timeRange ?? this.getDefaultTimeRange();
    const { start, end } = this.datasource.getTimeRangeParams(range);

    const cacheKey = this.generateCacheKey(url, start, end, interpolatedMatch);
    let value = this.seriesCache.get(cacheKey);
    if (!value) {
      const params = { 'match[]': interpolatedMatch, start, end };
      const data = await this.request(url, params);
      const { values } = processLabels(data);
      value = values;
      this.seriesCache.set(cacheKey, value);
    }
    return value;
  };

  /**
   * Fetch series for a selector. Use this for raw results. Use fetchSeriesLabels() to get labels.
   * @param match
   * @param streamSelector - The stream selector for which you want to retrieve labels.
   * @param options - (Optional) An object containing additional options.
   * @param options.timeRange - (Optional) The time range for which you want to retrieve label keys. If not provided, the default time range is used.
   * @returns A promise containing array with records of label names and their value.
   */
  fetchSeries = async (match: string, options?: { timeRange?: TimeRange }): Promise<Array<Record<string, string>>> => {
    const url = 'series';
    const range = options?.timeRange ?? this.getDefaultTimeRange();
    const { start, end } = this.datasource.getTimeRangeParams(range);
    const params = { 'match[]': match, start, end };
    return await this.request(url, params);
  };

  // Cache key is a bit different here. We round up to a minute the intervals.
  // The rounding may seem strange but makes relative intervals like now-1h less prone to need separate request every
  // millisecond while still actually getting all the keys for the correct interval. This still can create problems
  // when user does not the newest values for a minute if already cached.
  private generateCacheKey(url: string, start: number, end: number, param: string): string {
    return [url, this.roundTime(start), this.roundTime(end), param].join();
  }

  // Round nanoseconds epoch to nearest 5 minute interval
  private roundTime(nanoseconds: number): number {
    return nanoseconds ? Math.floor(nanoseconds / NS_IN_MS / 1000 / 60 / 5) : 0;
  }

  /**
   * Fetch label values
   *
   * This asynchronous function fetches values associated with a specified label name.
   * It returns a promise that resolves to an array of strings containing the label values.
   *
   * @param labelName - The name of the label for which you want to retrieve values.
   * @param options - (Optional) An object containing additional options.
   * @param options.streamSelector - (Optional) The stream selector to filter label values. If not provided, all label values are fetched.
   * @param options.timeRange - (Optional) The time range for which you want to retrieve label values. If not provided, the default time range is used.
   * @returns A promise containing an array of label values.
   * @throws An error if the fetch operation fails.
   */
  async fetchLabelValues(
    labelName: string,
    options?: { streamSelector?: string; timeRange?: TimeRange }
  ): Promise<string[]> {
    const label = encodeURIComponent(this.datasource.interpolateString(labelName));
    const streamParam = options?.streamSelector
      ? encodeURIComponent(this.datasource.interpolateString(options.streamSelector))
      : undefined;

    const url = `label/${label}/values`;
    const range = options?.timeRange ?? this.getDefaultTimeRange();
    const rangeParams = this.datasource.getTimeRangeParams(range);
    const { start, end } = rangeParams;
    const params: KeyValue<string | number> = { start, end };
    let paramCacheKey = label;

    if (streamParam) {
      params.query = streamParam;
      paramCacheKey += streamParam;
    }

    const cacheKey = this.generateCacheKey(url, start, end, paramCacheKey);

    let labelValues = this.labelsCache.get(cacheKey);
    if (!labelValues) {
      // Clear value when requesting new one. Empty object being truthy also makes sure we don't request twice.
      this.labelsCache.set(cacheKey, []);
      const res = await this.request(url, params);
      if (Array.isArray(res)) {
        labelValues = res.slice().sort();
        this.labelsCache.set(cacheKey, labelValues);
      }
    }

    return labelValues ?? [];
  }

  /**
   * Get parser and label keys for a selector
   *
   * This asynchronous function is used to fetch parsers and label keys for a selected log stream based on sampled lines.
   * It returns a promise that resolves to an object with the following properties:
   *
   * - `extractedLabelKeys`: An array of available label keys associated with the log stream.
   * - `hasJSON`: A boolean indicating whether JSON parsing is available for the stream.
   * - `hasLogfmt`: A boolean indicating whether Logfmt parsing is available for the stream.
   * - `hasPack`: A boolean indicating whether Pack parsing is available for the stream.
   * - `unwrapLabelKeys`: An array of label keys that can be used for unwrapping log data.
   *
   * @param streamSelector - The selector for the log stream you want to analyze.
   * @param options - (Optional) An object containing additional options.
   * @param options.maxLines - (Optional) The number of log lines requested when determining parsers and label keys.
   * @param options.timeRange - (Optional) The time range for which you want to retrieve label keys. If not provided, the default time range is used.
   * Smaller maxLines is recommended for improved query performance. The default count is 10.
   * @returns A promise containing an object with parser and label key information.
   * @throws An error if the fetch operation fails.
   */
  async getParserAndLabelKeys(
    streamSelector: string,
    options?: { maxLines?: number; timeRange?: TimeRange }
  ): Promise<ParserAndLabelKeysResult> {
    const empty = {
      extractedLabelKeys: [],
      structuredMetadataKeys: [],
      unwrapLabelKeys: [],
      hasJSON: false,
      hasLogfmt: false,
      hasPack: false,
    };
    if (!config.featureToggles.lokiQueryHints) {
      return empty;
    }

    const series = await this.datasource.getDataSamples(
      {
        expr: streamSelector,
        refId: 'data-samples',
        maxLines: options?.maxLines || DEFAULT_MAX_LINES_SAMPLE,
      },
      options?.timeRange ?? this.getDefaultTimeRange()
    );

    if (!series.length) {
      return empty;
    }

    const { hasLogfmt, hasJSON, hasPack } = extractLogParserFromDataFrame(series[0]);

    return {
      extractedLabelKeys: [
        ...extractLabelKeysFromDataFrame(series[0], LabelType.Indexed),
        ...extractLabelKeysFromDataFrame(series[0], LabelType.Parsed),
      ],
      structuredMetadataKeys: extractLabelKeysFromDataFrame(series[0], LabelType.StructuredMetadata),
      unwrapLabelKeys: extractUnwrapLabelKeysFromDataFrame(series[0]),
      hasJSON,
      hasPack,
      hasLogfmt,
    };
  }

  /**
   * Get the default time range
   *
   * @returns {TimeRange} The default time range
   */
  private getDefaultTimeRange(): TimeRange {
    return getDefaultTimeRange();
  }
}
