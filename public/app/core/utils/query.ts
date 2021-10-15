import { DataQuery, DataSourceInstanceSettings, DataSourceRef } from '@grafana/data';

export const getNextRefIdChar = (queries: DataQuery[]): string => {
  for (let num = 0; ; num++) {
    const refId = getRefId(num);
    if (!queries.some((query) => query.refId === refId)) {
      return refId;
    }
  }
};

export function addQuery(queries: DataQuery[], query?: Partial<DataQuery>, datasource?: DataSourceRef): DataQuery[] {
  const q = query || {};
  q.refId = getNextRefIdChar(queries);
  q.hide = false;

  if (!q.datasource && datasource) {
    q.datasource = datasource;
  }

  return [...queries, q as DataQuery];
}

export function updateQueries(
  newSettings: DataSourceInstanceSettings,
  queries: DataQuery[],
  extensionID: string, // pass this in because importing it creates a circular dependency
  dsSettings?: DataSourceInstanceSettings
): DataQuery[] {
  const datasource: DataSourceRef = { uid: newSettings.uid, type: newSettings.type };

  if (!newSettings.meta.mixed && dsSettings?.meta.mixed) {
    return queries.map((q) => {
      if (q.datasource !== extensionID) {
        q.datasource = datasource;
      }
      return q;
    });
  } else if (!newSettings.meta.mixed && dsSettings?.meta.id !== newSettings.meta.id) {
    // we are changing data source type, clear queries
    return [{ refId: 'A', datasource }];
  }

  return queries;
}

export function isDataQuery(url: string): boolean {
  if (
    url.indexOf('api/datasources/proxy') !== -1 ||
    url.indexOf('api/tsdb/query') !== -1 ||
    url.indexOf('api/ds/query') !== -1
  ) {
    return true;
  }

  return false;
}

export function isLocalUrl(url: string) {
  return !url.match(/^http/);
}

function getRefId(num: number): string {
  const letters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';

  if (num < letters.length) {
    return letters[num];
  } else {
    return getRefId(Math.floor(num / letters.length) - 1) + letters[num % letters.length];
  }
}
