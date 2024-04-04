import React from 'react';
import { render, fireEvent} from '@testing-library/react';
import { FilterByValueTransformerEditor } from './FilterByValueTransformerEditor';
import {DataFrame, FieldType, ValueMatcherID, valueMatchers} from "@grafana/data";
import {FilterByValueMatch, FilterByValueType} from "@grafana/data/src/transformations/transformers/filterByValue";

describe('FilterByValueTransformerEditor', () => {
  it('adds a filter with default isNull option when onAddFilter is called', () => {
    // Mock onChange function
    const onChangeMock = jest.fn();

    // Mock options
    const options = {
      type: FilterByValueType.include,
      match: FilterByValueMatch.all,
      filters: [],
    };

    // Mock input
    const input: DataFrame[] = [
      {
        fields: [
          { name: 'person', type: FieldType.string, config: { displayName: 'Person' }, values: ['john', 'jill', 'jeremy', ''] },
          { name: 'city', type: FieldType.string, config: { displayName: 'City' }, values: ['london', 'budapest', '', 'lisbon'] }
        ],
        length: 4
      }
    ];

    // Render the component
    const { getByText } = render(
      <FilterByValueTransformerEditor input={input} options={options} onChange={onChangeMock} />
    );

    // Find and click the "Add condition" button
    fireEvent.click(getByText('Add condition'));

    // Check if onChange was called with the correct filter
    expect(onChangeMock).toHaveBeenCalledWith({
      filters: [
        {
          fieldName: 'Person',
          config: {
            id: ValueMatcherID.isNull,
            options: valueMatchers.get(ValueMatcherID.isNull).getDefaultOptions(
              { name: 'person', type: FieldType.string, config: { displayName: 'Person' }, values: ['john', 'jill', 'jeremy', ''] }
            ),
          },
        },
      ],
      match: FilterByValueMatch.all,
      type: FilterByValueType.include,
    });
  });
});
