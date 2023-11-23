import { DataFrame, Labels, roundDecimals } from '@grafana/data';
import { CombinedRuleNamespace } from 'app/types/unified-alerting';

import { isCloudRulesSource } from '../../utils/datasource';

/**
 * ⚠️ `frame.fields` could be an empty array ⚠️
 *
 * TypeScript will NOT complain about it when accessing items via index signatures.
 * Make sure to check for empty array or use optional chaining!
 *
 * see https://github.com/Microsoft/TypeScript/issues/13778
 */

const getSeriesName = (frame: DataFrame): string | undefined => {
  const firstField = frame.fields[0];

  const displayNameFromDS = firstField?.config?.displayNameFromDS;
  return displayNameFromDS ?? frame.name ?? firstField?.labels?.__name__;
};

const getSeriesValue = (frame: DataFrame) => {
  const value = frame.fields[0]?.values[0];

  if (Number.isFinite(value)) {
    return roundDecimals(value, 5);
  }

  return value;
};

const getSeriesLabels = (frame: DataFrame): Record<string, string> => {
  const firstField = frame.fields[0];
  return firstField?.labels ?? {};
};

const formatLabels = (labels: Labels): string => {
  return Object.entries(labels)
    .map(([key, value]) => key + '=' + value)
    .join(', ');
};

/**
 * After https://github.com/grafana/grafana/pull/74600,
 * Grafana rule names will be returned from the API as a combination of the folder name and parent UID separated by /.
 * For this reason, we need to format this namespace , removing the folder uid and the / separator.
 */
const GRAFANA_PARENT_FOLDER_SEPARATOR = '/';
const decodeGrafanaNamespace = (namespace: CombinedRuleNamespace): string => {
  const nameSpaceName = namespace.name;

  if (isCloudRulesSource(namespace.rulesSource)) {
    return nameSpaceName;
  }
  const nameSpaceWithoutSeparator = nameSpaceName.substring(nameSpaceName.indexOf(GRAFANA_PARENT_FOLDER_SEPARATOR) + 1);
  // replace \/ with / and \\ with \
  const nameSpaceWitoutSeparatorAndSlashesUnescaped = nameSpaceWithoutSeparator
    .replace(/\\\//g, '/')
    .replace(/\\\\/g, '\\');
  return nameSpaceWitoutSeparatorAndSlashesUnescaped;
};

const isEmptySeries = (series: DataFrame[]): boolean => {
  const isEmpty = series.every((serie) => serie.fields.every((field) => field.values.every((value) => value == null)));

  return isEmpty;
};

export {
  GRAFANA_PARENT_FOLDER_SEPARATOR,
  decodeGrafanaNamespace,
  formatLabels,
  getSeriesLabels,
  getSeriesName,
  getSeriesValue,
  isEmptySeries,
};
