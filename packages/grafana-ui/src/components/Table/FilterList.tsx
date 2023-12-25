import { css } from '@emotion/css';
import React, { useCallback, useMemo, useState } from 'react';
import { FixedSizeList as List } from 'react-window';

import { GrafanaTheme2, formattedValueToString, getValueFormat, SelectableValue } from '@grafana/data';

import { ButtonSelect, Checkbox, FilterInput, HorizontalGroup, Label, VerticalGroup } from '..';
import { useStyles2, useTheme2 } from '../../themes';

interface Props {
  values: SelectableValue[];
  options: SelectableValue[];
  onChange: (options: SelectableValue[]) => void;
  caseSensitive?: boolean;
  showOperators?: boolean;
}

const ITEM_HEIGHT = 28;
const MIN_HEIGHT = ITEM_HEIGHT * 5;

const operatorSelectableValue = (op: string): SelectableValue<string> => {
  const result: SelectableValue<string> = { label: op, value: op };
  switch (op) {
    case 'contains':
      result.description = 'Contains';
      break;
    case '=':
      result.description = 'Equals';
      break;
    case '!=':
      result.description = 'Not equals';
      break;
    case '>':
      result.description = 'Greater';
      break;
    case '>=':
      result.description = 'Greater or Equal';
      break;
    case '<':
      result.description = 'Less';
      break;
    case '<=':
      result.description = 'Less or Equal';
      break;
    case 'expression':
      result.description =
        'Bool Expression (Char v represents the column value in the expression, e.g. "v >= 10 && v <= 12")';
      break;
  }
  return result;
};
const OPERATORS = ['contains', '=', '!=', '<', '<=', '>', '>=', 'expression'].map(operatorSelectableValue);
const REGEX_OPERATOR = OPERATORS.filter((op) => op.value === 'contains')[0];
const XPR_OPERATOR = OPERATORS.filter((op) => op.value === 'expression')[0];

const comparableValue = (value: string): string | number | Date | boolean => {
  value = value.trim().replace(/\\/g, '');

  // Does it look like a Date (Starting with pattern YYYY-MM-DD* or YYYY/MM/DD*)?
  if (/^(\d{4}-\d{2}-\d{2}|\d{4}\/\d{2}\/\d{2})/.test(value)) {
    const date = new Date(value);
    if (!isNaN(date.getTime())) {
      const fmt = getValueFormat('dateTimeAsIso');
      return formattedValueToString(fmt(date.getTime()));
    }
  }
  // Does it look like a Number?
  const num = parseFloat(value);
  if (!isNaN(num)) {
    return num;
  }
  // Does it look like a Bool?
  const lvalue = value.toLowerCase();
  if (lvalue === 'true' || lvalue === 'false') {
    return lvalue === 'true';
  }
  // Anything else
  return value;
};

export const FilterList = ({ options, values, caseSensitive, showOperators, onChange }: Props) => {
  const [operator, setOperator] = useState<SelectableValue<string>>(REGEX_OPERATOR);
  const [searchFilter, setSearchFilter] = useState('');
  const regex = useMemo(() => new RegExp(searchFilter, caseSensitive ? undefined : 'i'), [searchFilter, caseSensitive]);
  const items = useMemo(
    () =>
      options.filter((option) => {
        if (!showOperators || !searchFilter || operator.value === REGEX_OPERATOR.value) {
          if (option.label === undefined) {
            return false;
          }
          return regex.test(option.label);
        } else if (operator.value === XPR_OPERATOR.value) {
          if (option.value === undefined) {
            return false;
          }
          try {
            const xpr = searchFilter.replace(/\\/g, '');
            const fnc = new Function('v', `'use strict'; return ${xpr};`);
            const val = comparableValue(option.value);
            return fnc(val);
          } catch (_) {}
          return false;
        } else {
          if (option.value === undefined) {
            return false;
          }

          const value1 = comparableValue(option.value);
          const value2 = comparableValue(searchFilter);

          switch (operator.value) {
            case '=':
              return value1 === value2;
            case '!=':
              return value1 !== value2;
            case '>':
              return value1 > value2;
            case '>=':
              return value1 >= value2;
            case '<':
              return value1 < value2;
            case '<=':
              return value1 <= value2;
          }
          return false;
        }
      }),
    [options, regex, showOperators, operator, searchFilter]
  );

  const styles = useStyles2(getStyles);
  const theme = useTheme2();
  const gutter = theme.spacing.gridSize;
  const height = useMemo(() => Math.min(items.length * ITEM_HEIGHT, MIN_HEIGHT) + gutter, [gutter, items.length]);

  const onCheckedChanged = useCallback(
    (option: SelectableValue) => (event: React.FormEvent<HTMLInputElement>) => {
      const newValues = event.currentTarget.checked
        ? values.concat(option)
        : values.filter((c) => c.value !== option.value);

      onChange(newValues);
    },
    [onChange, values]
  );

  return (
    <VerticalGroup spacing="md">
      {!showOperators && <FilterInput placeholder="Filter values" onChange={setSearchFilter} value={searchFilter} />}
      {showOperators && (
        <HorizontalGroup>
          <ButtonSelect<string>
            variant="canvas"
            options={OPERATORS}
            onChange={setOperator}
            value={operator}
            tooltip={operator.description}
          />
          <FilterInput placeholder="Filter values" onChange={setSearchFilter} value={searchFilter} />
        </HorizontalGroup>
      )}
      {!items.length && <Label>No values</Label>}
      {items.length && (
        <List
          height={height}
          itemCount={items.length}
          itemSize={ITEM_HEIGHT}
          width="100%"
          className={styles.filterList}
        >
          {({ index, style }) => {
            const option = items[index];
            const { value, label } = option;
            const isChecked = values.find((s) => s.value === value) !== undefined;

            return (
              <div className={styles.filterListRow} style={style} title={label}>
                <Checkbox value={isChecked} label={label} onChange={onCheckedChanged(option)} />
              </div>
            );
          }}
        </List>
      )}
    </VerticalGroup>
  );
};

const getStyles = (theme: GrafanaTheme2) => ({
  filterList: css({
    label: 'filterList',
  }),
  filterListRow: css({
    label: 'filterListRow',
    cursor: 'pointer',
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    padding: theme.spacing(0.5),

    ':hover': {
      backgroundColor: theme.colors.action.hover,
    },
  }),
});
