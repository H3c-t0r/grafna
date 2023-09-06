import { LanguageMap, languages as prismLanguages } from 'prismjs';
import React, { ReactNode } from 'react';
import { Node, Plugin } from 'slate';
import { Editor } from 'slate-react';

import { AbsoluteTimeRange, QueryEditorProps } from '@grafana/data';
import {
  BracesPlugin,
  QueryField,
  SlatePrism,
  Themeable2,
  TypeaheadInput,
  TypeaheadOutput,
  withTheme2,
} from '@grafana/ui';

// Utils & Services
// dom also includes Element polyfills
import { CloudWatchDatasource } from '../../../datasource';
import syntax from '../../../language/cloudwatch-logs/syntax';
import { CloudWatchJsonData, CloudWatchLogsQuery, CloudWatchQuery, LogGroup } from '../../../types';
import { getStatsGroups } from '../../../utils/query/getStatsGroups';
import { LogGroupsFieldWrapper } from '../../shared/LogGroups/LogGroupsField';

export interface CloudWatchLogsQueryFieldProps
  extends QueryEditorProps<CloudWatchDatasource, CloudWatchQuery, CloudWatchJsonData>,
    Themeable2 {
  absoluteRange: AbsoluteTimeRange;
  onLabelsRefresh?: () => void;
  ExtraFieldElement?: ReactNode;
  exploreId: string;
  query: CloudWatchLogsQuery;
}
const plugins: Array<Plugin<Editor>> = [
  BracesPlugin(),
  SlatePrism(
    {
      onlyIn: (node: Node) => node.object === 'block' && node.type === 'code_block',
      getSyntax: (node: Node) => 'cloudwatch',
    },
    { ...(prismLanguages as LanguageMap), cloudwatch: syntax }
  ),
];
export const CloudWatchLogsQueryField = (props: CloudWatchLogsQueryFieldProps) => {
  const { query, datasource, onChange, ExtraFieldElement, data } = props;

  const showError = data?.error?.refId === query.refId;
  const cleanText = datasource.languageProvider.cleanText;

  const onChangeQuery = (value: string) => {
    // Send text change to parent
    const nextQuery = {
      ...query,
      expression: value,
      statsGroups: getStatsGroups(value),
    };
    onChange(nextQuery);
  };

  const onTypeahead = async (typeahead: TypeaheadInput): Promise<TypeaheadOutput> => {
    const { datasource, query } = props;
    const { logGroups } = query;

    if (!datasource.languageProvider) {
      return { suggestions: [] };
    }

    const { history, absoluteRange } = props;
    const { prefix, text, value, wrapperClasses, labelKey, editor } = typeahead;

    return await datasource.languageProvider.provideCompletionItems(
      { text, value, prefix, wrapperClasses, labelKey, editor },
      {
        history,
        absoluteRange,
        logGroups: logGroups,
        region: query.region,
      }
    );
  };

  return (
    <>
      <LogGroupsFieldWrapper
        region={query.region}
        datasource={datasource}
        legacyLogGroupNames={query.logGroupNames}
        logGroups={query.logGroups}
        onChange={(logGroups: LogGroup[]) => {
          onChange({ ...query, logGroups, logGroupNames: undefined });
        }}
        //legacy props can be removed once we remove support for Legacy Log Group Selector
        legacyOnChange={(logGroups: string[]) => {
          onChange({ ...query, logGroupNames: logGroups });
        }}
      />
      <div className="gf-form-inline gf-form-inline--nowrap flex-grow-1">
        <div className="gf-form gf-form--grow flex-shrink-1">
          <QueryField
            additionalPlugins={plugins}
            query={query.expression ?? ''}
            onChange={onChangeQuery}
            onTypeahead={onTypeahead}
            cleanText={cleanText}
            placeholder="Enter a CloudWatch Logs Insights query (run with Shift+Enter)"
            portalOrigin="cloudwatch"
          />
        </div>
        {ExtraFieldElement}
      </div>
      {showError ? (
        <div className="query-row-break">
          <div className="prom-query-field-info text-error">{data?.error?.message}</div>
        </div>
      ) : null}
    </>
  );
};

export default withTheme2(CloudWatchLogsQueryField);
