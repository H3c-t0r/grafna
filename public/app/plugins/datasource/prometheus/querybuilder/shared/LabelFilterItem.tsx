import { css, cx } from '@emotion/css';
import { uniqBy } from 'lodash';
import React, { useState } from 'react';

import { GrafanaTheme2, SelectableValue, toOption } from '@grafana/data';
import { selectors } from '@grafana/e2e-selectors';
import { AccessoryButton, InputGroup } from '@grafana/experimental';
import { Select, useStyles2 } from '@grafana/ui';

import { isConflictingSelector } from './operationUtils';
import { QueryBuilderLabelFilter } from './types';

export interface Props {
  defaultOp: string;
  item: Partial<QueryBuilderLabelFilter>;
  items: Array<Partial<QueryBuilderLabelFilter>>;
  onChange: (value: QueryBuilderLabelFilter) => void;
  onGetLabelNames: (forLabel: Partial<QueryBuilderLabelFilter>) => Promise<SelectableValue[]>;
  onGetLabelValues: (forLabel: Partial<QueryBuilderLabelFilter>) => Promise<SelectableValue[]>;
  onDelete: () => void;
  invalidLabel?: boolean;
  invalidValue?: boolean;
}

export function LabelFilterItem({
  item,
  items,
  defaultOp,
  onChange,
  onDelete,
  onGetLabelNames,
  onGetLabelValues,
  invalidLabel,
  invalidValue,
}: Props) {
  const [state, setState] = useState<{
    labelNames?: SelectableValue[];
    labelValues?: SelectableValue[];
    isLoadingLabelNames?: boolean;
    isLoadingLabelValues?: boolean;
  }>({});
  const styles = useStyles2(getStyles);

  const isMultiSelect = (operator = item.op) => {
    return operators.find((op) => op.label === operator)?.isMultiValue;
  };

  const getSelectOptionsFromString = (item?: string): string[] => {
    if (item) {
      if (item.indexOf('|') > 0) {
        return item.split('|');
      }
      return [item];
    }
    return [];
  };

  const getOptions = (): SelectableValue[] => {
    const labelValues = state.labelValues ? [...state.labelValues] : [];
    const selectedOptions = getSelectOptionsFromString(item?.value).map(toOption);

    // Remove possible duplicated values
    return uniqBy([...selectedOptions, ...labelValues], 'value');
  };

  const isConflicting = isConflictingSelector(item, items);

  return (
    <div data-testid="prometheus-dimensions-filter-item" className={cx(isConflicting && styles.isConflicting)}>
      <InputGroup>
        <Select
          placeholder="Select label"
          aria-label={selectors.components.QueryBuilder.labelSelect}
          inputId="prometheus-dimensions-filter-item-key"
          width="auto"
          value={item.label ? toOption(item.label) : null}
          allowCustomValue
          onOpenMenu={async () => {
            setState({ isLoadingLabelNames: true });
            const labelNames = await onGetLabelNames(item);
            setState({ labelNames, isLoadingLabelNames: undefined });
          }}
          isLoading={state.isLoadingLabelNames}
          options={state.labelNames}
          onChange={(change) => {
            if (change.label) {
              onChange({
                ...item,
                op: item.op ?? defaultOp,
                label: change.label,
              } as unknown as QueryBuilderLabelFilter);
            }
          }}
          invalid={invalidLabel}
        />

        <Select
          aria-label={selectors.components.QueryBuilder.matchOperatorSelect}
          value={toOption(item.op ?? defaultOp)}
          options={operators}
          width="auto"
          onChange={(change) => {
            if (change.value != null) {
              onChange({
                ...item,
                op: change.value,
                value: isMultiSelect(change.value) ? item.value : getSelectOptionsFromString(item?.value)[0],
              } as unknown as QueryBuilderLabelFilter);
            }
          }}
        />

        <Select
          placeholder="Select value"
          aria-label={selectors.components.QueryBuilder.valueSelect}
          inputId="prometheus-dimensions-filter-item-value"
          width="auto"
          value={
            isMultiSelect()
              ? getSelectOptionsFromString(item?.value).map(toOption)
              : getSelectOptionsFromString(item?.value).map(toOption)[0]
          }
          allowCustomValue
          onOpenMenu={async () => {
            setState({ isLoadingLabelValues: true });
            const labelValues = await onGetLabelValues(item);
            setState({
              ...state,
              labelValues,
              isLoadingLabelValues: undefined,
            });
          }}
          isMulti={isMultiSelect()}
          isLoading={state.isLoadingLabelValues}
          options={getOptions()}
          onChange={(change) => {
            if (change.value) {
              onChange({
                ...item,
                value: change.value,
                op: item.op ?? defaultOp,
              } as unknown as QueryBuilderLabelFilter);
            } else {
              const changes = change
                .map((change: any) => {
                  return change.label;
                })
                .join('|');
              onChange({ ...item, value: changes, op: item.op ?? defaultOp } as unknown as QueryBuilderLabelFilter);
            }
          }}
          invalid={invalidValue}
        />
        <AccessoryButton aria-label="remove" icon="times" variant="secondary" onClick={onDelete} />
      </InputGroup>
    </div>
  );
}

const operators = [
  { label: '=~', value: '=~', isMultiValue: true },
  { label: '=', value: '=', isMultiValue: false },
  { label: '!=', value: '!=', isMultiValue: false },
  { label: '!~', value: '!~', isMultiValue: true },
];

const getStyles = (theme: GrafanaTheme2) => {
  return {
    isConflicting: css({
      boxShadow: `0px 0px 4px 0px ${theme.colors.error.main}`,
      border: `1px solid ${theme.colors.error.main}`,
    }),
  };
};
