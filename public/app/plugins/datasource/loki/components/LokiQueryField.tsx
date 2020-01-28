import React, { FunctionComponent } from 'react';
import { LokiQueryFieldForm, LokiQueryFieldFormProps } from './LokiQueryFieldForm';
import { useLokiSyntaxAndLabels } from './useLokiSyntaxAndLabels';
import LokiLanguageProvider from '../language_provider';

type LokiQueryFieldProps = Omit<
  LokiQueryFieldFormProps,
  'syntax' | 'syntaxLoaded' | 'onLoadOptions' | 'onLabelsRefresh' | 'logLabelOptions'
>;

export const LokiQueryField: FunctionComponent<LokiQueryFieldProps> = props => {
  const { datasource, absoluteRange, ...otherProps } = props;
  const { isSyntaxReady, setActiveOption, refreshLabels, syntax, logLabelOptions } = useLokiSyntaxAndLabels(
    datasource.languageProvider as LokiLanguageProvider,
    absoluteRange
  );

  return (
    <LokiQueryFieldForm
      datasource={datasource}
      /**
       * setActiveOption name is intentional. Because of the way rc-cascader requests additional data
       * https://github.com/react-component/cascader/blob/master/src/Cascader.jsx#L165
       * we are notyfing useLokiSyntax hook, what the active option is, and then it's up to the hook logic
       * to fetch data of options that aren't fetched yet
       */
      onLoadOptions={setActiveOption}
      onLabelsRefresh={refreshLabels}
      absoluteRange={absoluteRange}
      syntax={syntax}
      syntaxLoaded={isSyntaxReady}
      logLabelOptions={logLabelOptions}
      {...otherProps}
    />
  );
};

export default LokiQueryField;
