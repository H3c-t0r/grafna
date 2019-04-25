import { FieldType } from '../../types/data';
import { SeriesDataMatcherID } from '../matchers/index';
import { SeriesTransformerID, transformSeriesData } from './transformers';
export const simpleSeriesWithTypes = {
  fields: [
    { name: 'A', type: FieldType.time },
    { name: 'B', type: FieldType.boolean },
    { name: 'C', type: FieldType.string },
    { name: 'D', type: FieldType.number },
  ],
  rows: [[1, 2, 3, 4], [4, 5, 6, 7]],
};

describe('Filter Transformer', () => {
  it('filters by include', () => {
    const cfg = {
      id: SeriesTransformerID.filter,
      options: {
        include: { id: SeriesDataMatcherID.numericFields },
      },
    };

    const filtered = transformSeriesData([cfg], [simpleSeriesWithTypes])[0];
    expect(filtered.fields.length).toBe(1);
    expect(filtered.fields[0].name).toBe('D');
  });
});
