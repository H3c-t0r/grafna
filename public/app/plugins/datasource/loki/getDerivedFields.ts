import { groupBy } from 'lodash';

import { FieldType, DataFrame, DataLink, Field } from '@grafana/data';
import { getDataSourceSrv, config } from '@grafana/runtime';

import { DerivedFieldConfig } from './types';

export function getDerivedFields(dataFrame: DataFrame, derivedFieldConfigs: DerivedFieldConfig[]): Field[] {
  if (!derivedFieldConfigs.length) {
    return [];
  }
  const derivedFieldsGrouped = groupBy(derivedFieldConfigs, 'name');

  const newFields = Object.values(derivedFieldsGrouped).map(fieldFromDerivedFieldConfig);

  const fieldsByName = new Map(newFields.map((field) => [field.name, field]));

  // line-field is the first string-field
  // NOTE: we should create some common log-frame-extra-string-field code somewhere
  const lineField = dataFrame.fields.find((f) => f.type === FieldType.string);

  if (lineField === undefined) {
    // if this is happening, something went wrong, let's raise an error
    throw new Error('invalid logs-dataframe, string-field missing');
  }

  const labelFields = dataFrame.fields.find((f) => f.type === FieldType.other && f.name === 'labels');

  for (let i = 0; i < lineField.values.length; i++) {
    for (const field of newFields) {
      if (
        config.featureToggles.lokiEnableNameMatcherOption &&
        derivedFieldsGrouped[field.name][0].enableNameMatcher &&
        labelFields
      ) {
        const label = labelFields.values[i];
        if (label) {
          const intersectingKey = Object.keys(label).find((key) => fieldsByName.has(key));

          if (intersectingKey) {
            field.values.push(label[intersectingKey]);
            continue;
          }
        }
      }

      const line = lineField.values[i];
      const logMatch = line.match(derivedFieldsGrouped[field.name][0].matcherRegex);

      if (logMatch && logMatch[1]) {
        field.values.push(logMatch[1]);
        continue;
      }

      field.values.push(null);
    }
  }

  return newFields;
}

/**
 * Transform derivedField config into dataframe field with config that contains link.
 */
function fieldFromDerivedFieldConfig(derivedFieldConfigs: DerivedFieldConfig[]): Field {
  const dataSourceSrv = getDataSourceSrv();

  const dataLinks = derivedFieldConfigs.reduce<DataLink[]>((acc, derivedFieldConfig) => {
    // Having field.datasourceUid means it is an internal link.
    if (derivedFieldConfig.datasourceUid) {
      const dsSettings = dataSourceSrv.getInstanceSettings(derivedFieldConfig.datasourceUid);
      const queryType = (type: string | undefined): string | undefined => {
        switch (type) {
          case 'tempo':
            return 'traceql';
          case 'grafana-x-ray-datasource':
            return 'getTrace';
          default:
            return undefined;
        }
      };

      acc.push({
        // Will be filled out later
        title: derivedFieldConfig.urlDisplayLabel || '',
        url: '',
        // This is hardcoded for Jaeger or Zipkin not way right now to specify datasource specific query object
        internal: {
          query: { query: derivedFieldConfig.url, queryType: queryType(dsSettings?.type) },
          datasourceUid: derivedFieldConfig.datasourceUid,
          datasourceName: dsSettings?.name ?? 'Data source not found',
        },
      });
    } else if (derivedFieldConfig.url) {
      acc.push({
        // We do not know what title to give here so we count on presentation layer to create a title from metadata.
        title: derivedFieldConfig.urlDisplayLabel || '',
        // This is hardcoded for Jaeger or Zipkin not way right now to specify datasource specific query object
        url: derivedFieldConfig.url,
      });
    }
    return acc;
  }, []);

  return {
    name: derivedFieldConfigs[0].name,
    type: FieldType.string,
    config: {
      links: dataLinks,
    },
    // We are adding values later on
    values: [],
  };
}
