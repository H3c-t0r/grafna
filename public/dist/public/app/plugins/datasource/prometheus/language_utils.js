import { __rest } from "tslib";
import { invert } from 'lodash';
import { Token } from 'prismjs';
import { AbstractLabelOperator, dateMath, incrRoundDn, } from '@grafana/data';
import { addLabelToQuery } from './add_label_to_query';
import { SUGGESTIONS_LIMIT } from './language_provider';
import { PROMETHEUS_QUERY_BUILDER_MAX_RESULTS } from './querybuilder/components/MetricSelect';
import { PrometheusCacheLevel } from './types';
export const processHistogramMetrics = (metrics) => {
    const resultSet = new Set();
    const regexp = new RegExp('_bucket($|:)');
    for (let index = 0; index < metrics.length; index++) {
        const metric = metrics[index];
        const isHistogramValue = regexp.test(metric);
        if (isHistogramValue) {
            resultSet.add(metric);
        }
    }
    return [...resultSet];
};
export function processLabels(labels, withName = false) {
    // For processing we are going to use sets as they have significantly better performance than arrays
    // After we process labels, we will convert sets to arrays and return object with label values in arrays
    const valueSet = {};
    labels.forEach((label) => {
        const { __name__ } = label, rest = __rest(label, ["__name__"]);
        if (withName) {
            valueSet['__name__'] = valueSet['__name__'] || new Set();
            if (!valueSet['__name__'].has(__name__)) {
                valueSet['__name__'].add(__name__);
            }
        }
        Object.keys(rest).forEach((key) => {
            if (!valueSet[key]) {
                valueSet[key] = new Set();
            }
            if (!valueSet[key].has(rest[key])) {
                valueSet[key].add(rest[key]);
            }
        });
    });
    // valueArray that we are going to return in the object
    const valueArray = {};
    limitSuggestions(Object.keys(valueSet)).forEach((key) => {
        valueArray[key] = limitSuggestions(Array.from(valueSet[key]));
    });
    return { values: valueArray, keys: Object.keys(valueArray) };
}
// const cleanSelectorRegexp = /\{(\w+="[^"\n]*?")(,\w+="[^"\n]*?")*\}/;
export const selectorRegexp = /\{[^}]*?(\}|$)/;
export const labelRegexp = /\b(\w+)(!?=~?)("[^"\n]*?")/g;
export function parseSelector(query, cursorOffset = 1) {
    if (!query.match(selectorRegexp)) {
        // Special matcher for metrics
        if (query.match(/^[A-Za-z:][\w:]*$/)) {
            return {
                selector: `{__name__="${query}"}`,
                labelKeys: ['__name__'],
            };
        }
        throw new Error('Query must contain a selector: ' + query);
    }
    // Check if inside a selector
    const prefix = query.slice(0, cursorOffset);
    const prefixOpen = prefix.lastIndexOf('{');
    const prefixClose = prefix.lastIndexOf('}');
    if (prefixOpen === -1) {
        throw new Error('Not inside selector, missing open brace: ' + prefix);
    }
    if (prefixClose > -1 && prefixClose > prefixOpen) {
        throw new Error('Not inside selector, previous selector already closed: ' + prefix);
    }
    const suffix = query.slice(cursorOffset);
    const suffixCloseIndex = suffix.indexOf('}');
    const suffixClose = suffixCloseIndex + cursorOffset;
    const suffixOpenIndex = suffix.indexOf('{');
    const suffixOpen = suffixOpenIndex + cursorOffset;
    if (suffixClose === -1) {
        throw new Error('Not inside selector, missing closing brace in suffix: ' + suffix);
    }
    if (suffixOpenIndex > -1 && suffixOpen < suffixClose) {
        throw new Error('Not inside selector, next selector opens before this one closed: ' + suffix);
    }
    // Extract clean labels to form clean selector, incomplete labels are dropped
    const selector = query.slice(prefixOpen, suffixClose);
    const labels = {};
    selector.replace(labelRegexp, (label, key, operator, value) => {
        const labelOffset = query.indexOf(label);
        const valueStart = labelOffset + key.length + operator.length + 1;
        const valueEnd = labelOffset + key.length + operator.length + value.length - 1;
        // Skip label if cursor is in value
        if (cursorOffset < valueStart || cursorOffset > valueEnd) {
            labels[key] = { value, operator };
        }
        return '';
    });
    // Add metric if there is one before the selector
    const metricPrefix = query.slice(0, prefixOpen);
    const metricMatch = metricPrefix.match(/[A-Za-z:][\w:]*$/);
    if (metricMatch) {
        labels['__name__'] = { value: `"${metricMatch[0]}"`, operator: '=' };
    }
    // Build sorted selector
    const labelKeys = Object.keys(labels).sort();
    const cleanSelector = labelKeys.map((key) => `${key}${labels[key].operator}${labels[key].value}`).join(',');
    const selectorString = ['{', cleanSelector, '}'].join('');
    return { labelKeys, selector: selectorString };
}
export function expandRecordingRules(query, mapping) {
    const ruleNames = Object.keys(mapping);
    const rulesRegex = new RegExp(`(\\s|^)(${ruleNames.join('|')})(\\s|$|\\(|\\[|\\{)`, 'ig');
    const expandedQuery = query.replace(rulesRegex, (match, pre, name, post) => `${pre}${mapping[name]}${post}`);
    // Split query into array, so if query uses operators, we can correctly add labels to each individual part.
    const queryArray = expandedQuery.split(/(\+|\-|\*|\/|\%|\^)/);
    // Regex that matches occurrences of ){ or }{ or ]{ which is a sign of incorrecly added labels.
    const invalidLabelsRegex = /(\)\{|\}\{|\]\{)/;
    const correctlyExpandedQueryArray = queryArray.map((query) => {
        return addLabelsToExpression(query, invalidLabelsRegex);
    });
    return correctlyExpandedQueryArray.join('');
}
function addLabelsToExpression(expr, invalidLabelsRegexp) {
    var _a;
    const match = expr.match(invalidLabelsRegexp);
    if (!match) {
        return expr;
    }
    // Split query into 2 parts - before the invalidLabelsRegex match and after.
    const indexOfRegexMatch = (_a = match.index) !== null && _a !== void 0 ? _a : 0;
    const exprBeforeRegexMatch = expr.slice(0, indexOfRegexMatch + 1);
    const exprAfterRegexMatch = expr.slice(indexOfRegexMatch + 1);
    // Create arrayOfLabelObjects with label objects that have key, operator and value.
    const arrayOfLabelObjects = [];
    exprAfterRegexMatch.replace(labelRegexp, (label, key, operator, value) => {
        arrayOfLabelObjects.push({ key, operator, value });
        return '';
    });
    // Loop through all label objects and add them to query.
    // As a starting point we have valid query without the labels.
    let result = exprBeforeRegexMatch;
    arrayOfLabelObjects.filter(Boolean).forEach((obj) => {
        // Remove extra set of quotes from obj.value
        const value = obj.value.slice(1, -1);
        result = addLabelToQuery(result, obj.key, value, obj.operator);
    });
    return result;
}
/**
 * Adds metadata for synthetic metrics for which the API does not provide metadata.
 * See https://github.com/grafana/grafana/issues/22337 for details.
 *
 * @param metadata HELP and TYPE metadata from /api/v1/metadata
 */
export function fixSummariesMetadata(metadata) {
    if (!metadata) {
        return metadata;
    }
    const baseMetadata = {};
    const summaryMetadata = {};
    for (const metric in metadata) {
        // NOTE: based on prometheus-documentation, we can receive
        // multiple metadata-entries for the given metric, it seems
        // it happens when the same metric is on multiple targets
        // and their help-text differs
        // (https://prometheus.io/docs/prometheus/latest/querying/api/#querying-metric-metadata)
        // for now we just use the first entry.
        const item = metadata[metric][0];
        baseMetadata[metric] = item;
        if (item.type === 'histogram') {
            summaryMetadata[`${metric}_bucket`] = {
                type: 'counter',
                help: `Cumulative counters for the observation buckets (${item.help})`,
            };
            summaryMetadata[`${metric}_count`] = {
                type: 'counter',
                help: `Count of events that have been observed for the histogram metric (${item.help})`,
            };
            summaryMetadata[`${metric}_sum`] = {
                type: 'counter',
                help: `Total sum of all observed values for the histogram metric (${item.help})`,
            };
        }
        if (item.type === 'summary') {
            summaryMetadata[`${metric}_count`] = {
                type: 'counter',
                help: `Count of events that have been observed for the base metric (${item.help})`,
            };
            summaryMetadata[`${metric}_sum`] = {
                type: 'counter',
                help: `Total sum of all observed values for the base metric (${item.help})`,
            };
        }
    }
    // Synthetic series
    const syntheticMetadata = {};
    syntheticMetadata['ALERTS'] = {
        type: 'counter',
        help: 'Time series showing pending and firing alerts. The sample value is set to 1 as long as the alert is in the indicated active (pending or firing) state.',
    };
    return Object.assign(Object.assign(Object.assign({}, baseMetadata), summaryMetadata), syntheticMetadata);
}
export function roundMsToMin(milliseconds) {
    return roundSecToMin(milliseconds / 1000);
}
export function roundSecToMin(seconds) {
    return Math.floor(seconds / 60);
}
// Returns number of minutes rounded up to the nearest nth minute
export function roundSecToNextMin(seconds, secondsToRound = 1) {
    return Math.ceil(seconds / 60) - (Math.ceil(seconds / 60) % secondsToRound);
}
export function limitSuggestions(items) {
    return items.slice(0, SUGGESTIONS_LIMIT);
}
export function addLimitInfo(items) {
    return items && items.length >= SUGGESTIONS_LIMIT ? `, limited to the first ${SUGGESTIONS_LIMIT} received items` : '';
}
// NOTE: the following 2 exported functions are very similar to the prometheus*Escape
// functions in datasource.ts, but they are not exactly the same algorithm, and we found
// no way to reuse one in the another or vice versa.
// Prometheus regular-expressions use the RE2 syntax (https://github.com/google/re2/wiki/Syntax),
// so every character that matches something in that list has to be escaped.
// the list of metacharacters is: *+?()|\.[]{}^$
// we make a javascript regular expression that matches those characters:
const RE2_METACHARACTERS = /[*+?()|\\.\[\]{}^$]/g;
function escapePrometheusRegexp(value) {
    return value.replace(RE2_METACHARACTERS, '\\$&');
}
// based on the openmetrics-documentation, the 3 symbols we have to handle are:
// - \n ... the newline character
// - \  ... the backslash character
// - "  ... the double-quote character
export function escapeLabelValueInExactSelector(labelValue) {
    return labelValue.replace(/\\/g, '\\\\').replace(/\n/g, '\\n').replace(/"/g, '\\"');
}
export function escapeLabelValueInRegexSelector(labelValue) {
    return escapeLabelValueInExactSelector(escapePrometheusRegexp(labelValue));
}
const FromPromLikeMap = {
    '=': AbstractLabelOperator.Equal,
    '!=': AbstractLabelOperator.NotEqual,
    '=~': AbstractLabelOperator.EqualRegEx,
    '!~': AbstractLabelOperator.NotEqualRegEx,
};
const ToPromLikeMap = invert(FromPromLikeMap);
export function toPromLikeExpr(labelBasedQuery) {
    const expr = labelBasedQuery.labelMatchers
        .map((selector) => {
        const operator = ToPromLikeMap[selector.operator];
        if (operator) {
            return `${selector.name}${operator}"${selector.value}"`;
        }
        else {
            return '';
        }
    })
        .filter((e) => e !== '')
        .join(', ');
    return expr ? `{${expr}}` : '';
}
export function toPromLikeQuery(labelBasedQuery) {
    return {
        refId: labelBasedQuery.refId,
        expr: toPromLikeExpr(labelBasedQuery),
        range: true,
    };
}
export function extractLabelMatchers(tokens) {
    const labelMatchers = [];
    for (let prop in tokens) {
        if (tokens[prop] instanceof Token) {
            let token = tokens[prop];
            if (token.type === 'context-labels') {
                let labelKey = '';
                let labelValue = '';
                let labelOperator = '';
                let contentTokens = token.content;
                for (let currentToken in contentTokens) {
                    if (typeof contentTokens[currentToken] === 'string') {
                        let currentStr;
                        currentStr = contentTokens[currentToken];
                        if (currentStr === '=' || currentStr === '!=' || currentStr === '=~' || currentStr === '!~') {
                            labelOperator = currentStr;
                        }
                    }
                    else if (contentTokens[currentToken] instanceof Token) {
                        switch (contentTokens[currentToken].type) {
                            case 'label-key':
                                labelKey = contentTokens[currentToken].content;
                                break;
                            case 'label-value':
                                labelValue = contentTokens[currentToken].content;
                                labelValue = labelValue.substring(1, labelValue.length - 1);
                                const labelComparator = FromPromLikeMap[labelOperator];
                                if (labelComparator) {
                                    labelMatchers.push({ name: labelKey, operator: labelComparator, value: labelValue });
                                }
                                break;
                        }
                    }
                }
            }
        }
    }
    return labelMatchers;
}
/**
 * Calculates new interval "snapped" to the closest Nth minute, depending on cache level datasource setting
 * @param cacheLevel
 * @param range
 */
export function getRangeSnapInterval(cacheLevel, range) {
    // Don't round the range if we're not caching
    if (cacheLevel === PrometheusCacheLevel.None) {
        return {
            start: getPrometheusTime(range.from, false).toString(),
            end: getPrometheusTime(range.to, true).toString(),
        };
    }
    // Otherwise round down to the nearest nth minute for the start time
    const startTime = getPrometheusTime(range.from, false);
    // const startTimeQuantizedSeconds = roundSecToLastMin(startTime, getClientCacheDurationInMinutes(cacheLevel)) * 60;
    const startTimeQuantizedSeconds = incrRoundDn(startTime, getClientCacheDurationInMinutes(cacheLevel) * 60);
    // And round up to the nearest nth minute for the end time
    const endTime = getPrometheusTime(range.to, true);
    const endTimeQuantizedSeconds = roundSecToNextMin(endTime, getClientCacheDurationInMinutes(cacheLevel)) * 60;
    // If the interval was too short, we could have rounded both start and end to the same time, if so let's add one step to the end
    if (startTimeQuantizedSeconds === endTimeQuantizedSeconds) {
        const endTimePlusOneStep = endTimeQuantizedSeconds + getClientCacheDurationInMinutes(cacheLevel) * 60;
        return { start: startTimeQuantizedSeconds.toString(), end: endTimePlusOneStep.toString() };
    }
    const start = startTimeQuantizedSeconds.toString();
    const end = endTimeQuantizedSeconds.toString();
    return { start, end };
}
export function getClientCacheDurationInMinutes(cacheLevel) {
    switch (cacheLevel) {
        case PrometheusCacheLevel.Medium:
            return 10;
        case PrometheusCacheLevel.High:
            return 60;
        default:
            return 1;
    }
}
export function getPrometheusTime(date, roundUp) {
    if (typeof date === 'string') {
        date = dateMath.parse(date, roundUp);
    }
    return Math.ceil(date.valueOf() / 1000);
}
export function truncateResult(array, limit) {
    if (limit === undefined) {
        limit = PROMETHEUS_QUERY_BUILDER_MAX_RESULTS;
    }
    array.length = Math.min(array.length, limit);
    return array;
}
//# sourceMappingURL=language_utils.js.map