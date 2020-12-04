// Libraries
import React, { ReactNode } from 'react';

import {
  CascaderOption,
  SlatePrism,
  TypeaheadOutput,
  SuggestionsState,
  QueryField,
  TypeaheadInput,
  BracesPlugin,
  DOMUtil,
} from '@grafana/ui';

// Utils & Services
// dom also includes Element polyfills
import { Plugin, Node } from 'slate';
import { LokiLabelBrowser } from './LokiLabelBrowser';

// Types
import { ExploreQueryFieldProps, AbsoluteTimeRange } from '@grafana/data';
import { LokiQuery, LokiOptions } from '../types';
import { LanguageMap, languages as prismLanguages } from 'prismjs';
import LokiLanguageProvider, { LokiHistoryItem } from '../language_provider';
import LokiDatasource from '../datasource';
import LokiOptionFields from './LokiOptionFields';

function getChooserText(hasSyntax: boolean, hasLogLabels: boolean) {
  if (!hasSyntax) {
    return 'Loading labels...';
  }
  if (!hasLogLabels) {
    return '(No labels found)';
  }
  return 'Log labels';
}

function willApplySuggestion(suggestion: string, { typeaheadContext, typeaheadText }: SuggestionsState): string {
  // Modify suggestion based on context
  switch (typeaheadContext) {
    case 'context-labels': {
      const nextChar = DOMUtil.getNextCharacter();
      if (!nextChar || nextChar === '}' || nextChar === ',') {
        suggestion += '=';
      }
      break;
    }

    case 'context-label-values': {
      // Always add quotes and remove existing ones instead
      if (!typeaheadText.match(/^(!?=~?"|")/)) {
        suggestion = `"${suggestion}`;
      }
      if (DOMUtil.getNextCharacter() !== '"') {
        suggestion = `${suggestion}"`;
      }
      break;
    }

    default:
  }
  return suggestion;
}

export interface LokiQueryFieldFormProps extends ExploreQueryFieldProps<LokiDatasource, LokiQuery, LokiOptions> {
  history: LokiHistoryItem[];
  logLabelOptions: CascaderOption[];
  labelsLoaded: boolean;
  absoluteRange: AbsoluteTimeRange;
  ExtraFieldElement?: ReactNode;
  runOnBlur?: boolean;
}

export class LokiQueryFieldForm extends React.PureComponent<LokiQueryFieldFormProps> {
  plugins: Plugin[];

  constructor(props: LokiQueryFieldFormProps, context: React.Context<any>) {
    super(props, context);

    this.plugins = [
      BracesPlugin(),
      SlatePrism(
        {
          onlyIn: (node: Node) => node.object === 'block' && node.type === 'code_block',
          getSyntax: (node: Node) => 'logql',
        },
        { ...(prismLanguages as LanguageMap), logql: this.props.datasource.languageProvider.getSyntax() }
      ),
    ];
  }

  onChangeLogLabels = (selector: string) => {
    this.onChangeQuery(selector, true);
  };

  onChangeQuery = (value: string, override?: boolean) => {
    // Send text change to parent
    const { query, onChange, onRunQuery } = this.props;
    if (onChange) {
      const nextQuery = { ...query, expr: value };
      onChange(nextQuery);

      if (override && onRunQuery) {
        onRunQuery();
      }
    }
  };

  onTypeahead = async (typeahead: TypeaheadInput): Promise<TypeaheadOutput> => {
    const { datasource } = this.props;

    if (!datasource.languageProvider) {
      return { suggestions: [] };
    }

    const lokiLanguageProvider = datasource.languageProvider as LokiLanguageProvider;
    const { history, absoluteRange } = this.props;
    const { prefix, text, value, wrapperClasses, labelKey } = typeahead;

    const result = await lokiLanguageProvider.provideCompletionItems(
      { text, value, prefix, wrapperClasses, labelKey },
      { history, absoluteRange }
    );
    return result;
  };

  render() {
    const { ExtraFieldElement, query, labelsLoaded, logLabelOptions, datasource, runOnBlur } = this.props;

    const lokiLanguageProvider = datasource.languageProvider as LokiLanguageProvider;
    const cleanText = datasource.languageProvider ? lokiLanguageProvider.cleanText : undefined;
    const hasLogLabels = logLabelOptions && logLabelOptions.length > 0;
    const chooserText = getChooserText(labelsLoaded, hasLogLabels);
    const buttonDisabled = !(labelsLoaded && hasLogLabels);

    return (
      <>
        <div className="gf-form-inline gf-form-inline--xs-view-flex-column flex-grow-1">
          <div className="gf-form flex-shrink-0 min-width-5">
            <LokiLabelBrowser
              buttonClass="gf-form-label"
              buttonText={chooserText}
              disabled={buttonDisabled}
              languageProvider={lokiLanguageProvider}
              onChange={this.onChangeLogLabels}
            />
          </div>
          <div className="gf-form gf-form--grow flex-shrink-1 min-width-15">
            <QueryField
              additionalPlugins={this.plugins}
              cleanText={cleanText}
              query={query.expr}
              onTypeahead={this.onTypeahead}
              onWillApplySuggestion={willApplySuggestion}
              onChange={this.onChangeQuery}
              onBlur={this.props.onBlur}
              onRunQuery={this.props.onRunQuery}
              placeholder="Enter a Loki query (run with Shift+Enter)"
              portalOrigin="loki"
            />
          </div>
        </div>
        <LokiOptionFields
          queryType={query.instant ? 'instant' : 'range'}
          lineLimitValue={query?.maxLines?.toString() || ''}
          query={query}
          onRunQuery={this.props.onRunQuery}
          onChange={this.props.onChange}
          runOnBlur={runOnBlur}
        />
        {ExtraFieldElement}
      </>
    );
  }
}
