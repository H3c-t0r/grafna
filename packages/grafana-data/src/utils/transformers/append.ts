import { DataTransformerInfo } from './transformers';
import { DataFrame } from '../../types/data';
import { DataTransformerID } from './ids';

interface AppendOptions {}

interface FieldIndex {
  [field: string]: number;
}

export const appendTransformer: DataTransformerInfo<AppendOptions> = {
  id: DataTransformerID.append,
  name: 'Append',
  description: 'Append all series',
  defaultOptions: {},

  /**
   * Return a modified copy of the series.  If the transform is not or should not
   * be applied, just return the input series
   */
  transformer: (options: AppendOptions) => {
    return (data: DataFrame[]) => {
      if (data.length < 2) {
        return data;
      }

      const processed: DataFrame = {
        fields: [],
        rows: [],
      };
      const index: FieldIndex = {};
      let width = data[0].fields.length;
      for (let i = 0; i < width; i++) {
        const f = data[0].fields[i];
        processed.fields.push(f);
        index[f.name] = i;
      }

      for (const series of data) {
        const arr: number[] = [];
        for (let i = 0; i < series.fields.length; i++) {
          const f = series.fields[i];
          if (!index.hasOwnProperty(f.name)) {
            index[f.name] = width++;
            processed.fields.push(f);
          }
          arr[i] = index[f.name];
        }
        for (const row of series.rows) {
          const norm: any[] = [];
          for (let i = 0; i < row.length; i++) {
            norm[arr[i]] = row[i];
          }
          processed.rows.push(norm);
        }
      }

      return [processed];
    };
  },
};
