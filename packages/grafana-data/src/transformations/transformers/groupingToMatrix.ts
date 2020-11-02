import { DataFrame, DataTransformerInfo, Field, FieldType } from '../../types';
import { DataTransformerID } from './ids';
import { MutableDataFrame } from '../../dataframe';
import { getFieldDisplayName } from '../../field/fieldState';

export interface GroupingToMatrixTransformerOptions {
  columnField?: string;
  rowField?: string;
  valueField?: string;
}

const DEFAULT_COLUMN_FIELD = 'Time';
const DEFAULT_ROW_FIELD = 'Time';
const DEFAULT_VALUE_FIELD = 'Value';

export const groupingToMatrixTransformer: DataTransformerInfo<GroupingToMatrixTransformerOptions> = {
  id: DataTransformerID.groupingToMatrix,
  name: 'Grouping to Matrix',
  description: 'Groups series by field and return a matrix visualisation',
  defaultOptions: {
    columnField: DEFAULT_COLUMN_FIELD,
    rowField: DEFAULT_ROW_FIELD,
    valueField: DEFAULT_VALUE_FIELD,
  },
  transformer: options => (data: DataFrame[]) => {
    const columnFieldMatch = options.columnField || DEFAULT_COLUMN_FIELD;
    const rowFieldMatch = options.rowField || DEFAULT_ROW_FIELD;
    const valueFieldMatch = options.valueField || DEFAULT_VALUE_FIELD;

    // Accept only single queries
    if (data.length !== 1) {
      return data;
    }

    const frame = data[0];
    const keyColumnField = findKeyField(frame, columnFieldMatch);
    const keyRowField = findKeyField(frame, rowFieldMatch);
    const valueField = findKeyField(frame, valueFieldMatch);
    const rowColumnField = `${rowFieldMatch}\\${columnFieldMatch}`;

    if (!keyColumnField || !keyRowField || !valueField) {
      return data;
    }

    const columnValues = [...new Set(keyColumnField.values.toArray())];
    const rowValues = [...new Set(keyRowField.values.toArray())];

    const matrixValues: { [key: string]: { [key: string]: any } } = {};

    for (let index = 0; index < valueField.values.length; index++) {
      const columnName = keyColumnField.values.get(index);
      const rowName = keyRowField.values.get(index);
      const value = valueField.values.get(index);

      if (!matrixValues[columnName]) {
        matrixValues[columnName] = {};
      }

      matrixValues[columnName][rowName] = value;
    }

    // FIELDS TO ADD
    // FIELD[0]:
    //  NAME: ROW/COLUMN
    //  VALUES: FIELD[ROW] -> VALUES
    // FIELD[1..N]:
    //  NAME: FIELD[COLUMN] -> VALUES[N]
    //  VALUES: FIELD[ROW][COLUMN] -> VALUE

    const resultFrame = new MutableDataFrame();

    resultFrame.addField({
      name: rowColumnField,
      values: rowValues,
      type: FieldType.string,
    });

    for (const columnName of columnValues) {
      let values = [];
      for (const rowName of rowValues) {
        const value = matrixValues[columnName][rowName] ?? '';
        values.push(value);
      }

      resultFrame.addField({
        name: columnName,
        values: values,
        config: valueField.config,
        type: valueField.type,
      });
    }

    return [resultFrame];
  },
};

function findKeyField(frame: DataFrame, matchTitle: string): Field | null {
  for (let fieldIndex = 0; fieldIndex < frame.fields.length; fieldIndex++) {
    const field = frame.fields[fieldIndex];

    if (matchTitle === getFieldDisplayName(field)) {
      return field;
    }
  }

  return null;
}
