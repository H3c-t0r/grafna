import { transformDataFrame } from './transformers';
import { ReducerID } from '../fieldReducer';
import { DataTransformerID } from './ids';

const seriesWithValues = {
  fields: [{ name: 'A' }, { name: 'B' }],
  rows: [
    [1, 2], // 1
    [2, 3], //
    [3, 4], //
    [4, 5], //
    [5, 6], //
    [6, 7], //
    [7, 8], //
  ],
};

describe('Calc Transformer', () => {
  it('filters by include', () => {
    const cfg = {
      id: DataTransformerID.calc,
      options: {
        calcs: [ReducerID.min, ReducerID.max, ReducerID.mean, ReducerID.delta],
      },
    };
    const filtered = transformDataFrame([cfg], [seriesWithValues])[0];
    expect(filtered.fields.length).toBe(5);
    expect(filtered).toMatchSnapshot();
  });
});
