import { ArrayVector, DataFrame, FieldType } from '@grafana/data';

const INSERT_MODES = {
  threshold: (prev: number, next: number, threshold: number) => prev + threshold,
  midpoint: (prev: number, next: number, threshold: number) => (prev + next) / 2,
  // previous time + 1ms to prevent StateTimeline from forward-interpolating prior state
  plusone: (prev: number, next: number, threshold: number) => prev + 1,
};

type InsertMode = keyof typeof INSERT_MODES;

export function nullInsertThreshold(
  frame: DataFrame,
  threshold?: number,
  refFieldName?: string,
  insertMode: InsertMode = 'threshold'
): DataFrame {
  if (frame.length < 2) {
    return frame;
  }

  let refField = frame.fields.find((field) => {
    // note: getFieldDisplayName() would require full DF[]
    return refFieldName != null ? field.name === refFieldName : field.type === FieldType.time;
  });

  if (refField == null) {
    return frame;
  }

  const getInsertValue = INSERT_MODES[insertMode];

  if (threshold == null) {
    threshold = refField.config.interval ?? 0;
  }

  if (threshold <= 0) {
    return frame;
  }

  let refValues = refField.values.toArray();
  let len = refValues.length;

  let prevValue: number = refValues[0];
  let refValuesNew: number[] = [prevValue];

  for (let i = 1; i < len; i++) {
    let curValue = refValues[i];

    if (curValue - prevValue > threshold) {
      refValuesNew.push(getInsertValue(prevValue, curValue, threshold));
    }

    refValuesNew.push(curValue);

    prevValue = curValue;
  }

  let filledLen = refValuesNew.length;

  let filledFieldValues: any[][] = [];

  for (let field of frame.fields) {
    let filledValues;

    if (field !== refField) {
      filledValues = Array(filledLen);

      let fieldValues = field.values.toArray();

      for (let i = 0, j = 0; i < filledLen; i++) {
        filledValues[i] = refValues[j] === refValuesNew[i] ? fieldValues[j++] : null;
      }
    } else {
      filledValues = refValuesNew;
    }

    filledFieldValues.push(filledValues);
  }

  let outFrame: DataFrame = {
    ...frame,
    length: filledLen,
    fields: frame.fields.map((field, i) => ({
      ...field,
      values: new ArrayVector(filledFieldValues[i]),
    })),
  };

  return outFrame;
}
