// Libraries
import React, { memo } from 'react';
import { css, cx } from '@emotion/css';
import { LokiQuery, StepType } from '../types';
import { map } from 'lodash';

// Types
import { InlineFormLabel, RadioButtonGroup, InlineField, Input, Select } from '@grafana/ui';
import { SelectableValue } from '@grafana/data';

export interface LokiOptionFieldsProps {
  lineLimitValue: string;
  stepInterval: string;
  stepMode: StepType;
  resolution: number;
  queryType: LokiQueryType;
  query: LokiQuery;
  onChange: (value: LokiQuery) => void;
  onRunQuery: () => void;
  runOnBlur?: boolean;
}

type LokiQueryType = 'instant' | 'range';

const queryTypeOptions = [
  { value: 'range', label: 'Range', description: 'Run query over a range of time.' },
  {
    value: 'instant',
    label: 'Instant',
    description: 'Run query against a single point in time. For this query, the "To" time is used.',
  },
];

const INTERVAL_FACTOR_OPTIONS: Array<SelectableValue<number>> = map([1, 2, 3, 4, 5, 10], (value: number) => ({
  value,
  label: '1/' + value,
}));

export const DEFAULT_STEP_OPTION: SelectableValue<StepType> = {
  value: 'min',
  label: 'Minimum',
};

const STEP_OPTIONS: Array<SelectableValue<StepType>> = [
  DEFAULT_STEP_OPTION,
  {
    value: 'max',
    label: 'Maximum',
  },
  {
    value: 'exact',
    label: 'Exact',
  },
];

export function LokiOptionFields(props: LokiOptionFieldsProps) {
  const {
    lineLimitValue,
    stepInterval,
    resolution,
    stepMode,
    queryType,
    query,
    onRunQuery,
    runOnBlur,
    onChange,
  } = props;

  function onChangeQueryLimit(value: string) {
    const nextQuery = { ...query, maxLines: preprocessMaxLines(value) };
    onChange(nextQuery);
  }

  function onQueryTypeChange(value: LokiQueryType) {
    let nextQuery;
    if (value === 'instant') {
      nextQuery = { ...query, instant: true, range: false };
    } else {
      nextQuery = { ...query, instant: false, range: true };
    }
    onChange(nextQuery);
  }

  function preprocessMaxLines(value: string): number {
    if (value.length === 0) {
      // empty input - falls back to dataSource.maxLines limit
      return NaN;
    } else if (value.length > 0 && (isNaN(+value) || +value < 0)) {
      // input with at least 1 character and that is either incorrect (value in the input field is not a number) or negative
      // falls back to the limit of 0 lines
      return 0;
    } else {
      // default case - correct input
      return +value;
    }
  }

  function onMaxLinesChange(e: React.SyntheticEvent<HTMLInputElement>) {
    if (query.maxLines !== preprocessMaxLines(e.currentTarget.value)) {
      onChangeQueryLimit(e.currentTarget.value);
    }
  }

  function onReturnKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
    if (e.key === 'Enter') {
      onRunQuery();
    }
  }

  function onStepIntervalChange(e: React.KeyboardEvent<HTMLInputElement>) {
    const nextQuery = { ...query, stepInterval: e.currentTarget.value };
    onChange(nextQuery);
  }

  function onStepModeChange(option: SelectableValue<StepType>) {
    const nextQuery = { ...query, stepMode: option.value };
    onChange(nextQuery);
  }

  function onResolutionChange(option: SelectableValue<number>) {
    const nextQuery = { ...query, resolution: option.value };
    onChange(nextQuery);
  }

  return (
    <div aria-label="Loki extra field" className="gf-form-inline">
      {/*Query type field*/}
      <div
        data-testid="queryTypeField"
        className={cx(
          'gf-form explore-input-margin',
          css`
            flex-wrap: nowrap;
          `
        )}
        aria-label="Query type field"
      >
        <InlineFormLabel width="auto">Query type</InlineFormLabel>

        <RadioButtonGroup
          options={queryTypeOptions}
          value={queryType}
          onChange={(type: LokiQueryType) => {
            onQueryTypeChange(type);
            if (runOnBlur) {
              onRunQuery();
            }
          }}
        />
      </div>
      {/*Line limit field*/}
      <div
        data-testid="lineLimitField"
        className={cx(
          'gf-form',
          css`
            flex-wrap: nowrap;
          `
        )}
        aria-label="Line limit field"
      >
        <InlineField label="Line limit">
          <Input
            className="width-4"
            placeholder="auto"
            type="number"
            min={0}
            onChange={onMaxLinesChange}
            onKeyDown={onReturnKeyDown}
            value={lineLimitValue}
            onBlur={() => {
              if (runOnBlur) {
                onRunQuery();
              }
            }}
          />
        </InlineField>
        <InlineField
          label="Step"
          tooltip={
            'Optionally, set the lower or upper bounds on the interval between data points, for example, set "minimum 1h" to hint that measurements were not taken more frequently. `$__interval` and `$__rate_interval` are supported.'
          }
        >
          <Select isSearchable={false} width={16} onChange={onStepModeChange} options={STEP_OPTIONS} value={stepMode} />
        </InlineField>
        <InlineField>
          <Input
            className="width-4"
            width={12}
            placeholder="15s"
            min={0}
            onChange={onStepIntervalChange}
            onKeyDown={onReturnKeyDown}
            value={stepInterval}
            onBlur={() => {
              if (runOnBlur) {
                onRunQuery();
              }
            }}
          />
        </InlineField>
        <InlineField label="Resolution">
          <Select
            isSearchable={false}
            options={INTERVAL_FACTOR_OPTIONS}
            onChange={onResolutionChange}
            value={resolution}
          />
        </InlineField>
      </div>
    </div>
  );
}

export default memo(LokiOptionFields);
