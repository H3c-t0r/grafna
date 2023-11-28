import { __awaiter } from "tslib";
import { trimEnd } from 'lodash';
import { escapeLabelValueInExactSelector } from '../../../languageUtils';
import { isQueryWithParser } from '../../../queryUtils';
import { explainOperator } from '../../../querybuilder/operations';
import { LokiOperationId } from '../../../querybuilder/types';
import { AGGREGATION_OPERATORS, RANGE_VEC_FUNCTIONS, BUILT_IN_FUNCTIONS } from '../../../syntax';
import { NeverCaseError } from './NeverCaseError';
const LOG_COMPLETIONS = [
    {
        type: 'PATTERN',
        label: '{}',
        insertText: '{$0}',
        isSnippet: true,
        triggerOnInsert: true,
    },
];
const AGGREGATION_COMPLETIONS = AGGREGATION_OPERATORS.map((f) => {
    var _a;
    return ({
        type: 'FUNCTION',
        label: f.label,
        insertText: `${(_a = f.insertText) !== null && _a !== void 0 ? _a : ''}($0)`,
        isSnippet: true,
        triggerOnInsert: true,
        detail: f.detail,
        documentation: f.documentation,
    });
});
const FUNCTION_COMPLETIONS = RANGE_VEC_FUNCTIONS.map((f) => {
    var _a;
    return ({
        type: 'FUNCTION',
        label: f.label,
        insertText: `${(_a = f.insertText) !== null && _a !== void 0 ? _a : ''}({$0}[\\$__auto])`,
        isSnippet: true,
        triggerOnInsert: true,
        detail: f.detail,
        documentation: f.documentation,
    });
});
const BUILT_IN_FUNCTIONS_COMPLETIONS = BUILT_IN_FUNCTIONS.map((f) => {
    var _a;
    return ({
        type: 'FUNCTION',
        label: f.label,
        insertText: `${(_a = f.insertText) !== null && _a !== void 0 ? _a : ''}($0)`,
        isSnippet: true,
        triggerOnInsert: true,
        detail: f.detail,
        documentation: f.documentation,
    });
});
const DURATION_COMPLETIONS = ['$__auto', '1m', '5m', '10m', '30m', '1h', '1d'].map((text) => ({
    type: 'DURATION',
    label: text,
    insertText: text,
}));
const UNWRAP_FUNCTION_COMPLETIONS = [
    {
        type: 'FUNCTION',
        label: 'duration_seconds',
        documentation: 'Will convert the label value in seconds from the go duration format (e.g 5m, 24s30ms).',
        insertText: 'duration_seconds()',
    },
    {
        type: 'FUNCTION',
        label: 'duration',
        documentation: 'Short version of duration_seconds().',
        insertText: 'duration()',
    },
    {
        type: 'FUNCTION',
        label: 'bytes',
        documentation: 'Will convert the label value to raw bytes applying the bytes unit (e.g. 5 MiB, 3k, 1G).',
        insertText: 'bytes()',
    },
];
const LOGFMT_ARGUMENT_COMPLETIONS = [
    {
        type: 'FUNCTION',
        label: '--strict',
        documentation: 'Strict parsing. The logfmt parser stops scanning the log line and returns early with an error when it encounters any poorly formatted key/value pair.',
        insertText: '--strict',
    },
    {
        type: 'FUNCTION',
        label: '--keep-empty',
        documentation: 'Retain standalone keys with empty value. The logfmt parser retains standalone keys (keys without a value) as labels with value set to empty string.',
        insertText: '--keep-empty',
    },
];
const LINE_FILTER_COMPLETIONS = [
    {
        operator: '|=',
        documentation: explainOperator(LokiOperationId.LineContains),
        afterPipe: true,
    },
    {
        operator: '!=',
        documentation: explainOperator(LokiOperationId.LineContainsNot),
    },
    {
        operator: '|~',
        documentation: explainOperator(LokiOperationId.LineMatchesRegex),
        afterPipe: true,
    },
    {
        operator: '!~',
        documentation: explainOperator(LokiOperationId.LineMatchesRegexNot),
    },
];
function getLineFilterCompletions(afterPipe) {
    return LINE_FILTER_COMPLETIONS.filter((completion) => !afterPipe || completion.afterPipe).map(({ operator, documentation }) => ({
        type: 'LINE_FILTER',
        label: `${operator} ""`,
        insertText: `${afterPipe ? operator.replace('|', '') : operator} "$0"`,
        isSnippet: true,
        documentation,
    }));
}
function getPipeOperationsCompletions(prefix = '') {
    const completions = [];
    completions.push({
        type: 'PIPE_OPERATION',
        label: 'line_format',
        insertText: `${prefix}line_format "{{.$0}}"`,
        isSnippet: true,
        documentation: explainOperator(LokiOperationId.LineFormat),
    });
    completions.push({
        type: 'PIPE_OPERATION',
        label: 'label_format',
        insertText: `${prefix}label_format`,
        isSnippet: true,
        documentation: explainOperator(LokiOperationId.LabelFormat),
    });
    completions.push({
        type: 'PIPE_OPERATION',
        label: 'unwrap',
        insertText: `${prefix}unwrap`,
        documentation: explainOperator(LokiOperationId.Unwrap),
    });
    completions.push({
        type: 'PIPE_OPERATION',
        label: 'decolorize',
        insertText: `${prefix}decolorize`,
        documentation: explainOperator(LokiOperationId.Decolorize),
    });
    completions.push({
        type: 'PIPE_OPERATION',
        label: 'drop',
        insertText: `${prefix}drop`,
        documentation: explainOperator(LokiOperationId.Drop),
    });
    completions.push({
        type: 'PIPE_OPERATION',
        label: 'keep',
        insertText: `${prefix}keep`,
        documentation: explainOperator(LokiOperationId.Keep),
    });
    return completions;
}
function getAllHistoryCompletions(dataProvider) {
    return __awaiter(this, void 0, void 0, function* () {
        const history = yield dataProvider.getHistory();
        return history.map((expr) => ({
            type: 'HISTORY',
            label: expr,
            insertText: expr,
        }));
    });
}
function getLabelNamesForSelectorCompletions(otherLabels, dataProvider) {
    return __awaiter(this, void 0, void 0, function* () {
        const labelNames = yield dataProvider.getLabelNames(otherLabels);
        return labelNames.map((label) => ({
            type: 'LABEL_NAME',
            label,
            insertText: `${label}=`,
            triggerOnInsert: true,
        }));
    });
}
function getInGroupingCompletions(logQuery, dataProvider) {
    return __awaiter(this, void 0, void 0, function* () {
        const { extractedLabelKeys } = yield dataProvider.getParserAndLabelKeys(logQuery);
        return extractedLabelKeys.map((label) => ({
            type: 'LABEL_NAME',
            label,
            insertText: label,
            triggerOnInsert: false,
        }));
    });
}
const PARSERS = ['json', 'logfmt', 'pattern', 'regexp', 'unpack'];
function getParserCompletions(prefix, hasJSON, hasLogfmt, hasPack, extractedLabelKeys, hasParserInQuery) {
    return __awaiter(this, void 0, void 0, function* () {
        const allParsers = new Set(PARSERS);
        const completions = [];
        // We use this to improve documentation specifically for level label as it is tied to showing color-coded logs volume
        const hasLevelInExtractedLabels = extractedLabelKeys.some((key) => key === 'level');
        if (hasJSON) {
            // We show "detected" label only if there is no previous parser in the query
            const extra = hasParserInQuery ? '' : ' (detected)';
            if (hasPack) {
                allParsers.delete('unpack');
                completions.push({
                    type: 'PARSER',
                    label: `unpack${extra}`,
                    insertText: `${prefix}unpack`,
                    documentation: explainOperator(LokiOperationId.Unpack),
                });
            }
            else {
                allParsers.delete('json');
                completions.push({
                    type: 'PARSER',
                    label: `json${extra}`,
                    insertText: `${prefix}json`,
                    documentation: hasLevelInExtractedLabels
                        ? 'Use it to get log-levels in the histogram'
                        : explainOperator(LokiOperationId.Json),
                });
            }
        }
        if (hasLogfmt) {
            allParsers.delete('logfmt');
            // We show "detected" label only if there is no previous parser in the query
            const extra = hasParserInQuery ? '' : ' (detected)';
            completions.push({
                type: 'PARSER',
                label: `logfmt${extra}`,
                insertText: `${prefix}logfmt`,
                documentation: hasLevelInExtractedLabels
                    ? 'Get detected levels in the histogram'
                    : explainOperator(LokiOperationId.Logfmt),
            });
        }
        const remainingParsers = Array.from(allParsers).sort();
        remainingParsers.forEach((parser) => {
            completions.push({
                type: 'PARSER',
                label: parser,
                insertText: `${prefix}${parser}`,
                documentation: explainOperator(parser),
            });
        });
        return completions;
    });
}
export function getAfterSelectorCompletions(logQuery, afterPipe, hasSpace, dataProvider) {
    return __awaiter(this, void 0, void 0, function* () {
        let query = logQuery;
        if (afterPipe) {
            query = trimEnd(logQuery, '| ');
        }
        const { extractedLabelKeys, hasJSON, hasLogfmt, hasPack } = yield dataProvider.getParserAndLabelKeys(query);
        const hasQueryParser = isQueryWithParser(query).queryWithParser;
        const prefix = `${hasSpace ? '' : ' '}${afterPipe ? '' : '| '}`;
        const parserCompletions = yield getParserCompletions(prefix, hasJSON, hasLogfmt, hasPack, extractedLabelKeys, hasQueryParser);
        const pipeOperations = getPipeOperationsCompletions(prefix);
        const completions = [...parserCompletions, ...pipeOperations];
        // Let's show label options only if query has parser
        if (hasQueryParser) {
            extractedLabelKeys.forEach((key) => {
                completions.push({
                    type: 'LABEL_NAME',
                    label: `${key} (detected)`,
                    insertText: `${prefix}${key}`,
                    documentation: `"${key}" was suggested based on the content of your log lines for the label filter expression.`,
                });
            });
        }
        // If we have parser, we don't need to consider line filters
        if (hasQueryParser) {
            return [...completions];
        }
        // With a space between the pipe and the cursor, we omit line filters
        // E.g. `{label="value"} | `
        const lineFilters = afterPipe && hasSpace ? [] : getLineFilterCompletions(afterPipe);
        return [...lineFilters, ...completions];
    });
}
export function getLogfmtCompletions(logQuery, flags, trailingComma, trailingSpace, otherLabels, dataProvider) {
    return __awaiter(this, void 0, void 0, function* () {
        let completions = [];
        if (trailingComma) {
            // Remove the trailing comma, otherwise the sample query will fail.
            logQuery = trimEnd(logQuery, ', ');
        }
        const { extractedLabelKeys, hasJSON, hasLogfmt, hasPack } = yield dataProvider.getParserAndLabelKeys(logQuery);
        const pipeOperations = getPipeOperationsCompletions('| ');
        /**
         * The user is not in the process of writing another label, and has not specified 2 flags.
         * The current grammar doesn't allow us to know which flags were used (by node name), so we consider flags = true
         * when 2 have been used.
         * For example:
         * - {label="value"} | logfmt ^
         * - {label="value"} | logfmt --strict ^
         * - {label="value"} | logfmt --strict --keep-empty ^
         */
        if (!trailingComma && !flags) {
            completions = [...LOGFMT_ARGUMENT_COMPLETIONS];
        }
        /**
         * If the user has no trailing comma and has a trailing space it can mean that they finished writing the logfmt
         * part and want to move on, for example, with other parsers or pipe operations.
         * For example:
         * - {label="value"} | logfmt --flag ^
         * - {label="value"} | logfmt label, label2 ^
         */
        if (!trailingComma && trailingSpace) {
            /**
             * Don't offer parsers if there is no label argument: {label="value"} | logfmt ^
             * The reason is that it would be unusual that they would want to use another parser just after logfmt, and
             * more likely that they would want a flag, labels, or continue with pipe operations.
             *
             * Offer parsers with at least one label argument: {label="value"} | logfmt label ^
             * The rationale here is to offer the same completions as getAfterSelectorCompletions().
             */
            const parserCompletions = otherLabels.length > 0
                ? yield getParserCompletions('| ', hasJSON, hasLogfmt, hasPack, extractedLabelKeys, true)
                : [];
            completions = [...completions, ...parserCompletions, ...pipeOperations];
        }
        const labels = extractedLabelKeys.filter((label) => !otherLabels.includes(label));
        /**
         * We want to decide whether to use a trailing comma or not based on the data we have of the current
         * situation. In particular, the following scenarios will not lead to a trailing comma:
         * {label="value"} | logfmt ^
         * - trailingSpace: true, trailingComma: false, otherLabels: []
         * {label="value"} | logfmt lab^
         * trailingSpace: false, trailignComma: false, otherLabels: [lab]
         * {label="value"} | logfmt label,^
         * trailingSpace: false, trailingComma: true, otherLabels: [label]
         * {label="value"} | logfmt label, ^
         * trailingSpace: true, trailingComma: true, otherLabels: [label]
         */
        let labelPrefix = '';
        if (otherLabels.length > 0 && trailingSpace) {
            labelPrefix = trailingComma ? '' : ', ';
        }
        const labelCompletions = labels.map((label) => ({
            type: 'LABEL_NAME',
            label,
            insertText: labelPrefix + label,
            triggerOnInsert: false,
        }));
        completions = [...completions, ...labelCompletions];
        return completions;
    });
}
function getLabelValuesForMetricCompletions(labelName, betweenQuotes, otherLabels, dataProvider) {
    return __awaiter(this, void 0, void 0, function* () {
        const values = yield dataProvider.getLabelValues(labelName, otherLabels);
        return values.map((text) => ({
            type: 'LABEL_VALUE',
            label: text,
            insertText: betweenQuotes ? escapeLabelValueInExactSelector(text) : `"${escapeLabelValueInExactSelector(text)}"`,
        }));
    });
}
function getAfterUnwrapCompletions(logQuery, dataProvider) {
    return __awaiter(this, void 0, void 0, function* () {
        const { unwrapLabelKeys } = yield dataProvider.getParserAndLabelKeys(logQuery);
        const labelCompletions = unwrapLabelKeys.map((label) => ({
            type: 'LABEL_NAME',
            label,
            insertText: label,
            triggerOnInsert: false,
        }));
        return [...labelCompletions, ...UNWRAP_FUNCTION_COMPLETIONS];
    });
}
function getAfterKeepAndDropCompletions(logQuery, dataProvider) {
    return __awaiter(this, void 0, void 0, function* () {
        const { extractedLabelKeys } = yield dataProvider.getParserAndLabelKeys(logQuery);
        const labelCompletions = extractedLabelKeys.map((label) => ({
            type: 'LABEL_NAME',
            label,
            insertText: label,
            triggerOnInsert: false,
        }));
        return [...labelCompletions];
    });
}
export function getCompletions(situation, dataProvider) {
    return __awaiter(this, void 0, void 0, function* () {
        switch (situation.type) {
            case 'EMPTY':
            case 'AT_ROOT':
                const historyCompletions = yield getAllHistoryCompletions(dataProvider);
                return [
                    ...historyCompletions,
                    ...LOG_COMPLETIONS,
                    ...AGGREGATION_COMPLETIONS,
                    ...BUILT_IN_FUNCTIONS_COMPLETIONS,
                    ...FUNCTION_COMPLETIONS,
                ];
            case 'IN_RANGE':
                return DURATION_COMPLETIONS;
            case 'IN_GROUPING':
                return getInGroupingCompletions(situation.logQuery, dataProvider);
            case 'IN_LABEL_SELECTOR_NO_LABEL_NAME':
                return getLabelNamesForSelectorCompletions(situation.otherLabels, dataProvider);
            case 'IN_LABEL_SELECTOR_WITH_LABEL_NAME':
                return getLabelValuesForMetricCompletions(situation.labelName, situation.betweenQuotes, situation.otherLabels, dataProvider);
            case 'AFTER_SELECTOR':
                return getAfterSelectorCompletions(situation.logQuery, situation.afterPipe, situation.hasSpace, dataProvider);
            case 'AFTER_UNWRAP':
                return getAfterUnwrapCompletions(situation.logQuery, dataProvider);
            case 'IN_AGGREGATION':
                return [...FUNCTION_COMPLETIONS, ...AGGREGATION_COMPLETIONS];
            case 'AFTER_KEEP_AND_DROP':
                return getAfterKeepAndDropCompletions(situation.logQuery, dataProvider);
            case 'IN_LOGFMT':
                return getLogfmtCompletions(situation.logQuery, situation.flags, situation.trailingComma, situation.trailingSpace, situation.otherLabels, dataProvider);
            default:
                throw new NeverCaseError(situation);
        }
    });
}
//# sourceMappingURL=completions.js.map