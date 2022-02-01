import { ArrayVector, FieldType, MutableDataFrame } from '@grafana/data';
import { nullInsertThreshold } from './nullInsertThreshold';

function randInt(min: number, max: number) {
  return Math.floor(Math.random() * (max - min + 1) + min);
}

function genFrame() {
  let fieldCount = 10;
  let valueCount = 3000;
  let step = 1000;
  let skipProb = 0.5;
  let skipSteps = [1, 5]; // min, max

  let allValues = Array(fieldCount);

  allValues[0] = Array(valueCount);

  for (let i = 0, curStep = Date.now(); i < valueCount; i++) {
    curStep = allValues[0][i] = curStep + step * (Math.random() < skipProb ? randInt(skipSteps[0], skipSteps[1]) : 1);
  }

  for (let fi = 1; fi < fieldCount; fi++) {
    let values = Array(valueCount);

    for (let i = 0; i < valueCount; i++) {
      values[i] = Math.random() * 100;
    }

    allValues[fi] = values;
  }

  return {
    length: valueCount,
    fields: allValues.map((values, i) => {
      return {
        name: 'A-' + i,
        type: i === 0 ? FieldType.time : FieldType.number,
        config: {
          interval: i === 0 ? step : null,
        },
        values: new ArrayVector(values),
      };
    }),
  };
}

describe('nullInsertThreshold Transformer', () => {
  test('should insert nulls at midpoints between adjacent > threshold: 1', () => {
    const df = new MutableDataFrame({
      refId: 'A',
      fields: [
        { name: 'Time', type: FieldType.time, values: [1, 3, 10] },
        { name: 'One', type: FieldType.number, values: [4, 6, 8] },
        { name: 'Two', type: FieldType.string, values: ['a', 'b', 'c'] },
      ],
    });

    const result = nullInsertThreshold(df, 1);

    expect(result.fields[0].values.toArray()).toStrictEqual([1, 2, 3, 4, 10]);
    expect(result.fields[1].values.toArray()).toStrictEqual([4, null, 6, null, 8]);
    expect(result.fields[2].values.toArray()).toStrictEqual(['a', null, 'b', null, 'c']);
  });

  test('should insert nulls at midpoints between adjacent > threshold: 2', () => {
    const df = new MutableDataFrame({
      refId: 'A',
      fields: [
        { name: 'Time', type: FieldType.time, values: [5, 7, 11] },
        { name: 'One', type: FieldType.number, values: [4, 6, 8] },
        { name: 'Two', type: FieldType.string, values: ['a', 'b', 'c'] },
      ],
    });

    const result = nullInsertThreshold(df, 2);

    expect(result.fields[0].values.toArray()).toStrictEqual([5, 7, 9, 11]);
    expect(result.fields[1].values.toArray()).toStrictEqual([4, 6, null, 8]);
    expect(result.fields[2].values.toArray()).toStrictEqual(['a', 'b', null, 'c']);
  });

  test('should insert nulls at midpoints between adjacent > interval', () => {
    const df = new MutableDataFrame({
      refId: 'A',
      fields: [
        { name: 'Time', type: FieldType.time, config: { interval: 1 }, values: [1, 3, 10] },
        { name: 'One', type: FieldType.number, values: [4, 6, 8] },
        { name: 'Two', type: FieldType.string, values: ['a', 'b', 'c'] },
      ],
    });

    const result = nullInsertThreshold(df);

    expect(result.fields[0].values.toArray()).toStrictEqual([1, 2, 3, 4, 10]);
    expect(result.fields[1].values.toArray()).toStrictEqual([4, null, 6, null, 8]);
    expect(result.fields[2].values.toArray()).toStrictEqual(['a', null, 'b', null, 'c']);
  });

  test('should insert nulls at midpoints between adjacent > interval: 2', () => {
    const df = new MutableDataFrame({
      refId: 'A',
      fields: [
        { name: 'Time', type: FieldType.time, config: { interval: 2 }, values: [5, 7, 11] },
        { name: 'One', type: FieldType.number, values: [4, 6, 8] },
        { name: 'Two', type: FieldType.string, values: ['a', 'b', 'c'] },
      ],
    });

    const result = nullInsertThreshold(df, 2);

    expect(result.fields[0].values.toArray()).toStrictEqual([5, 7, 9, 11]);
    expect(result.fields[1].values.toArray()).toStrictEqual([4, 6, null, 8]);
    expect(result.fields[2].values.toArray()).toStrictEqual(['a', 'b', null, 'c']);
  });

  test('should noop on fewer than two values', () => {
    const df = new MutableDataFrame({
      refId: 'A',
      fields: [
        { name: 'Time', type: FieldType.time, config: { interval: 1 }, values: [1] },
        { name: 'Value', type: FieldType.number, values: [1] },
      ],
    });

    const result = nullInsertThreshold(df);

    expect(result.fields[0].values.toArray()).toStrictEqual([1]);
    expect(result.fields[1].values.toArray()).toStrictEqual([1]);
  });

  test('should noop on invalid threshold', () => {
    const df = new MutableDataFrame({
      refId: 'A',
      fields: [
        { name: 'Time', type: FieldType.time, values: [1, 2, 4] },
        { name: 'Value', type: FieldType.number, values: [1, 1, 1] },
      ],
    });

    const result = nullInsertThreshold(df, -1);

    expect(result.fields[0].values.toArray()).toStrictEqual([1, 2, 4]);
    expect(result.fields[1].values.toArray()).toStrictEqual([1, 1, 1]);
  });

  test('should noop on invalid interval', () => {
    const df = new MutableDataFrame({
      refId: 'A',
      fields: [
        { name: 'Time', type: FieldType.time, config: { interval: -1 }, values: [1, 2, 4] },
        { name: 'Value', type: FieldType.number, values: [1, 1, 1] },
      ],
    });

    const result = nullInsertThreshold(df);

    expect(result.fields[0].values.toArray()).toStrictEqual([1, 2, 4]);
    expect(result.fields[1].values.toArray()).toStrictEqual([1, 1, 1]);
  });

  test('perf stress test should be <= 10ms', () => {
    // 10 fields x 3,000 values with 50% skip (output = 10 fields x 6,000 values)
    let bigFrameA = genFrame();

    // eslint-disable-next-line no-console
    console.time('insertValues-10x3k');
    nullInsertThreshold(bigFrameA);
    // eslint-disable-next-line no-console
    console.timeEnd('insertValues-10x3k');
  });
});
