import React, { useMemo, useCallback } from 'react';
import { css } from 'emotion';
import {
  DataTransformerID,
  standardTransformers,
  TransformerRegistyItem,
  TransformerUIProps,
  getFieldDisplayName,
  DataFrame,
  SelectableValue,
} from '@grafana/data';
import { Button, RadioButtonGroup, stylesFactory } from '@grafana/ui';
import cloneDeep from 'lodash/cloneDeep';
import {
  FilterByValueFilter,
  FilterByValueMatch,
  FilterByValueTransformerOptions,
  FilterByValueType,
} from '@grafana/data/src/transformations/transformers/filterByValue';

import { ValueFilterID } from '@grafana/data/src/transformations/valueFilters';
import { DataFrameFieldsInfo, FilterByValueFilterEditor } from './FilterByValueFilterEditor';

const filterTypes: Array<SelectableValue<FilterByValueType>> = [
  { label: 'Include values', value: FilterByValueType.include },
  { label: 'Exclude values', value: FilterByValueType.exclude },
];

const filterMatch: Array<SelectableValue<FilterByValueMatch>> = [
  { label: 'Match all', value: FilterByValueMatch.all },
  { label: 'Match any', value: FilterByValueMatch.any },
];

export const FilterByValueTransformerEditor: React.FC<TransformerUIProps<FilterByValueTransformerOptions>> = props => {
  const { input, options, onChange } = props;
  const styles = getEditorStyles();
  const fieldsInfo = useFieldsInfo(input);

  const onAddFilter = useCallback(() => {
    const filters = cloneDeep(options.filters);
    filters.push({
      fieldName: '',
      config: {
        id: ValueFilterID.greater,
        options: {},
      },
    });
    onChange({ ...options, filters });
  }, [onChange, options]);

  const onDeleteFilter = useCallback(
    (index: number) => () => {
      let filters = cloneDeep(options.filters);
      filters.splice(index, 1);
      onChange({ ...options, filters });
    },
    [options, onChange]
  );

  const onChangeFilter = useCallback(
    (filter: FilterByValueFilter, index: number) => {
      let filters = cloneDeep(options.filters);
      filters[index] = filter;
      onChange({ ...options, filters });
    },
    [options, onChange]
  );

  const onChangeType = useCallback(
    (type?: FilterByValueType) => {
      onChange({
        ...options,
        type: type ?? FilterByValueType.include,
      });
    },
    [options, onChange]
  );

  const onChangeMatch = useCallback(
    (match?: FilterByValueMatch) => {
      onChange({
        ...options,
        match: match ?? FilterByValueMatch.all,
      });
    },
    [options, onChange]
  );

  return (
    <div>
      <div className="gf-form gf-form-inline">
        <div className="gf-form-label">Filter type</div>
        <RadioButtonGroup options={filterTypes} value={options.type} onChange={onChangeType} />
      </div>
      <div className="gf-form gf-form-inline">
        <div className="gf-form-label gf-form--grow">Conditions</div>
        <RadioButtonGroup options={filterMatch} value={options.match} onChange={onChangeMatch} />
      </div>
      <div className={styles.conditions}>
        {options.filters.map((filter, idx) => (
          <FilterByValueFilterEditor
            key={idx}
            filter={filter}
            fieldsInfo={fieldsInfo}
            onChange={filter => onChangeFilter(filter, idx)}
            onDelete={() => onDeleteFilter(idx)}
          />
        ))}
        <div className="gf-form">
          <Button icon="plus" size="sm" onClick={onAddFilter} variant="secondary">
            Add condition
          </Button>
        </div>
      </div>
    </div>
  );
};

export const filterByValueTransformRegistryItem: TransformerRegistyItem<FilterByValueTransformerOptions> = {
  id: DataTransformerID.filterByValue,
  editor: FilterByValueTransformerEditor,
  transformation: standardTransformers.filterByValueTransformer,
  name: standardTransformers.filterByValueTransformer.name,
  description: standardTransformers.filterByValueTransformer.description,
};

const getEditorStyles = stylesFactory(() => ({
  conditions: css`
    padding-left: 16px;
  `,
}));

const useFieldsInfo = (data: DataFrame[]): DataFrameFieldsInfo => {
  return useMemo(() => {
    const meta = {
      fieldsAsOptions: [],
      fieldByDisplayName: {},
    };

    if (!Array.isArray(data)) {
      return meta;
    }

    return data.reduce((meta: DataFrameFieldsInfo, frame) => {
      return frame.fields.reduce((meta, field) => {
        const fieldName = getFieldDisplayName(field, frame, data);

        if (meta.fieldByDisplayName[fieldName]) {
          return meta;
        }

        meta.fieldsAsOptions.push({
          label: fieldName,
          value: fieldName,
          type: field.type,
        });

        meta.fieldByDisplayName[fieldName] = field;

        return meta;
      }, meta);
    }, meta);
  }, [data]);
};
