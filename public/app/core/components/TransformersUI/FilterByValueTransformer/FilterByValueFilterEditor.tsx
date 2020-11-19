import React, { useCallback } from 'react';
import { Button, Select } from '@grafana/ui';
import { Field, SelectableValue, valueMatchers } from '@grafana/data';
import { FilterByValueFilter } from '@grafana/data/src/transformations/transformers/filterByValue';
import { valueMatchersUI } from './ValueMatchers/valueMatchersUI';

interface Props {
  onDelete: () => void;
  onChange: (filter: FilterByValueFilter) => void;
  filter: FilterByValueFilter;
  fieldsInfo: DataFrameFieldsInfo;
}

export interface DataFrameFieldsInfo {
  fieldsAsOptions: Array<SelectableValue<string>>;
  fieldByDisplayName: Record<string, Field>;
}

export const FilterByValueFilterEditor: React.FC<Props> = props => {
  const { onDelete, onChange, filter, fieldsInfo } = props;
  const { fieldsAsOptions, fieldByDisplayName } = fieldsInfo;
  const fieldName = getFieldName(filter, fieldsAsOptions);
  const field = fieldByDisplayName[fieldName];

  if (!field) {
    return null;
  }

  const matcherOptions = getMatcherOptions(field);
  const matcherId = getSelectedMatcherId(filter, matcherOptions);
  const editor = valueMatchersUI.getIfExists(matcherId);

  console.log('mactcherId', matcherId);

  if (!editor || !editor.component) {
    return null;
  }

  const onChangeField = useCallback(
    (selectable?: SelectableValue<string>) => {
      if (!selectable?.value) {
        return;
      }
      onChange({
        ...filter,
        fieldName: selectable.value,
      });
    },
    [onChange, filter]
  );

  const onChangeMatcher = useCallback(
    (selectable?: SelectableValue<string>) => {
      if (!selectable?.value) {
        return;
      }
      onChange({
        ...filter,
        config: {
          id: selectable?.value,
        },
      });
    },
    [onChange, filter]
  );

  const onChangeMatcherOptions = useCallback(() => {}, [onChange, filter]);

  return (
    <div className="gf-form-inline">
      <div className="gf-form gf-form-spacing">
        <div className="gf-form-label width-4">Field</div>
        <Select
          className="width-24"
          placeholder="Field Name"
          options={fieldsAsOptions}
          value={filter.fieldName}
          onChange={onChangeField}
          menuPlacement="bottom"
        />
      </div>
      <div className="gf-form gf-form-spacing">
        <div className="gf-form-label width-8">Match</div>
        <Select
          className="width-8"
          placeholder="Select test"
          options={matcherOptions}
          value={matcherId}
          onChange={onChangeMatcher}
          menuPlacement="bottom"
        />
      </div>
      <div className="gf-form gf-form--grow gf-form-spacing ">
        <editor.component options={filter.config.options} onChange={onChangeMatcherOptions} />
      </div>
      <div className="gf-form">
        <Button icon="times" onClick={onDelete} style={{ height: '100%' }} size="sm" variant="secondary" />
      </div>
    </div>
  );
};

const getMatcherOptions = (field: Field): Array<SelectableValue<string>> => {
  const options = [];

  for (const matcher of valueMatchers.list()) {
    if (!matcher.isApplicable(field)) {
      continue;
    }

    const editor = valueMatchersUI.getIfExists(matcher.id);

    if (!editor) {
      continue;
    }

    options.push({
      value: matcher.id,
      label: matcher.name,
      description: matcher.description,
    });
  }

  return options;
};

const getSelectedMatcherId = (filter: FilterByValueFilter, matcherOptions: Array<SelectableValue<string>>): string => {
  const matcher = matcherOptions.find(m => m.value === filter.config.id);

  if (matcher && matcher.value) {
    return matcher.value;
  }

  if (matcherOptions[0]?.value) {
    return matcherOptions[0]?.value;
  }

  return '';
};

const getFieldName = (filter: FilterByValueFilter, fieldOptions: Array<SelectableValue<string>>): string => {
  const fieldName = fieldOptions.find(m => m.value === filter.fieldName);

  if (fieldName && fieldName.value) {
    return fieldName.value;
  }

  if (fieldOptions[0]?.value) {
    return fieldOptions[0]?.value;
  }

  return '';
};
