import { css } from '@emotion/css';
import React, { useCallback } from 'react';

import {
  DataTransformerID,
  ReducerID,
  SelectableValue,
  standardTransformers,
  TransformerRegistryItem,
  TransformerUIProps,
  TransformerCategory,
  GrafanaTheme2,
} from '@grafana/data';
import {
  GroupByFieldOptions,
  GroupByOperationID,
  GroupByTransformerOptions,
} from '@grafana/data/src/transformations/transformers/groupBy';
import {
  GroupToSubframeTransformerOptions,
  SHOW_SUBFRAME_HEADERS_DEFAULT,
} from '@grafana/data/src/transformations/transformers/groupToSubframe';
import { Stack } from '@grafana/experimental';
import { useTheme2, Select, StatsPicker, InlineField, Field, Switch } from '@grafana/ui';

import { useAllFieldNamesFromDataFrames } from '../utils';

interface FieldProps {
  fieldName: string;
  config?: GroupByFieldOptions;
  onConfigChange: (config: GroupByFieldOptions) => void;
}

export const GroupToSubframeTransformerEditor = ({
  input,
  options,
  onChange,
}: TransformerUIProps<GroupToSubframeTransformerOptions>) => {
  const fieldNames = useAllFieldNamesFromDataFrames(input);
  const showHeaders =
    options.showSubframeHeaders === undefined ? SHOW_SUBFRAME_HEADERS_DEFAULT : options.showSubframeHeaders;

  console.log(options);

  const onConfigChange = useCallback(
    (fieldName: string) => (config: GroupByFieldOptions) => {
      onChange({
        ...options,
        fields: {
          ...options.fields,
          [fieldName]: config,
        },
      });
    },
    // Adding options to the dependency array causes infinite loop here.
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [onChange]
  );

  const onShowFieldNamesChange = useCallback(
    () => {
      const showSubframeHeaders =
        options.showSubframeHeaders === undefined ? !SHOW_SUBFRAME_HEADERS_DEFAULT : !options.showSubframeHeaders;

      onChange({
        showSubframeHeaders,
        fields: {
          ...options.fields,
        },
      });
    },
    // Adding options to the dependency array causes infinite loop here.
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [onChange]
  );

  return (
    <Stack direction="row">
      <div>
        {fieldNames.map((key) => (
          <GroupByFieldConfiguration
            onConfigChange={onConfigChange(key)}
            fieldName={key}
            config={options.fields[key]}
            key={key}
          />
        ))}
      </div>
      <Field label="Show field names in subframe" description="If enabled subframes will show field names">
        <Switch value={showHeaders} onChange={onShowFieldNamesChange} />
      </Field>
    </Stack>
  );
};

const options = [
  { label: 'Group by', value: GroupByOperationID.groupBy },
  { label: 'Calculate', value: GroupByOperationID.aggregate },
];

export const GroupByFieldConfiguration = ({ fieldName, config, onConfigChange }: FieldProps) => {
  const theme = useTheme2();
  const styles = getStyles(theme);

  const onChange = useCallback(
    (value: SelectableValue<GroupByOperationID | null>) => {
      onConfigChange({
        aggregations: config?.aggregations ?? [],
        operation: value?.value ?? null,
      });
    },
    [config, onConfigChange]
  );

  return (
    <InlineField className={styles.label} label={fieldName} grow shrink>
      <Stack gap={0.5} direction="row" wrap={false}>
        <div className={styles.operation}>
          <Select options={options} value={config?.operation} placeholder="Ignored" onChange={onChange} isClearable />
        </div>

        {config?.operation === GroupByOperationID.aggregate && (
          <StatsPicker
            className={styles.aggregations}
            placeholder="Select Stats"
            allowMultiple
            stats={config.aggregations}
            onChange={(stats) => {
              onConfigChange({ ...config, aggregations: stats as ReducerID[] });
            }}
          />
        )}
      </Stack>
    </InlineField>
  );
};

const getStyles = (theme: GrafanaTheme2) => {
  return {
    label: css`
      label {
        min-width: ${theme.spacing(32)};
      }
    `,
    operation: css`
      flex-shrink: 0;
      height: 100%;
      width: ${theme.spacing(24)};
    `,
    aggregations: css`
      flex-grow: 1;
    `,
  };
};

export const groupToSubframeTransformRegistryItem: TransformerRegistryItem<GroupByTransformerOptions> = {
  id: DataTransformerID.groupToSubframe,
  editor: GroupToSubframeTransformerEditor,
  transformation: standardTransformers.groupToSubframeTransformer,
  name: standardTransformers.groupToSubframeTransformer.name,
  description: standardTransformers.groupToSubframeTransformer.description,
  categories: new Set([
    TransformerCategory.Combine,
    TransformerCategory.CalculateNewFields,
    TransformerCategory.Reformat,
  ]),
};
