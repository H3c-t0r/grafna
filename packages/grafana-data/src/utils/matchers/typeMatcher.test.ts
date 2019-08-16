import { FieldType } from '../../types/dataFrame';
import { dataMatchers } from './matchers';
import { DataMatcherID } from './ids';
import { toDataFrame } from '../processDataFrame';

export const simpleSeriesWithTypes = toDataFrame({
  fields: [
    { name: 'A', type: FieldType.time },
    { name: 'B', type: FieldType.boolean },
    { name: 'C', type: FieldType.string },
  ],
});

describe('Field Type Matcher', () => {
  const matcher = dataMatchers.get(DataMatcherID.fieldType);
  it('finds numbers', () => {
    for (const field of simpleSeriesWithTypes.fields) {
      expect(matcher.matcher(FieldType.number)(simpleSeriesWithTypes, field)).toBe(field.type === FieldType.number);
    }
  });
});
