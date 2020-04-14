import { ArrayDataFrame } from './ArrayDataFrame';
import { toDataFrameDTO } from './processDataFrame';
import { FieldType } from '../types';

describe('Array DataFrame', () => {
  const input = [
    { name: 'first', value: 1, time: 123 },
    { name: 'second', value: 2, time: 456, extra: 'here' },
    { name: 'third', value: 3, time: 789 },
  ];

  const frame = new ArrayDataFrame(input);
  frame.name = 'Hello';
  frame.refId = 'Z';
  frame.setFieldType('phantom', FieldType.string, v => '🦥');
  const field = frame.fields.find(f => f.name == 'value');
  field!.config.unit = 'kwh';

  test('Should support functional methods', () => {
    const expectedNames = input.map(row => row.name);

    // Check map
    expect(frame.map(row => row.name)).toEqual(expectedNames);

    let names: string[] = [];
    for (const row of frame) {
      names.push(row.name);
    }
    expect(names).toEqual(expectedNames);

    names = [];
    frame.forEach(row => {
      names.push(row.name);
    });
    expect(names).toEqual(expectedNames);
  });

  test('Should convert an array of objects to a dataframe', () => {
    expect(toDataFrameDTO(frame)).toMatchInlineSnapshot(`
      Object {
        "fields": Array [
          Object {
            "config": Object {},
            "labels": undefined,
            "name": "name",
            "type": "string",
            "values": Array [
              "first",
              "second",
              "third",
            ],
          },
          Object {
            "config": Object {
              "unit": "kwh",
            },
            "labels": undefined,
            "name": "value",
            "type": "number",
            "values": Array [
              1,
              2,
              3,
            ],
          },
          Object {
            "config": Object {},
            "labels": undefined,
            "name": "time",
            "type": "time",
            "values": Array [
              123,
              456,
              789,
            ],
          },
          Object {
            "config": Object {},
            "labels": undefined,
            "name": "phantom",
            "type": "string",
            "values": Array [
              "🦥",
              "🦥",
              "🦥",
            ],
          },
        ],
        "meta": undefined,
        "name": "Hello",
        "refId": "Z",
      }
    `);
  });
});
