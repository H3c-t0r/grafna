import { DataQuery, DataSourceJsonData, QueryResultMeta, ScopedVars } from '@grafana/data';
import { FetchError } from '@grafana/runtime';

export interface PromQuery extends DataQuery {
  expr: string;
  format?: string;
  instant?: boolean;
  range?: boolean;
  exemplar?: boolean;
  hinting?: boolean;
  interval?: string;
  intervalFactor?: number;
  legendFormat?: string;
  valueWithRefId?: boolean;
  requestId?: string;
  showingGraph?: boolean;
  showingTable?: boolean;
}

export interface PromOptions extends DataSourceJsonData {
  timeInterval: string;
  queryTimeout: string;
  httpMethod: string;
  directUrl: string;
  customQueryParameters?: string;
  disableMetricsLookup?: boolean;
  exemplarTraceIdDestination?: ExemplarTraceIdDestination;
}

export type ExemplarTraceIdDestination = {
  name: string;
  url?: string;
  datasourceUid?: string;
};

export interface PromQueryRequest extends PromQuery {
  step?: number;
  requestId?: string;
  start: number;
  end: number;
  headers?: any;
}

export interface PromMetricsMetadataItem {
  type: string;
  help: string;
  unit?: string;
}

export interface PromMetricsMetadata {
  [metric: string]: PromMetricsMetadataItem[];
}

export interface PromDataSuccessResponse<T = PromData> {
  status: 'success';
  data: T;
}

export interface PromDataErrorResponse<T = PromData> {
  status: 'error';
  errorType: string;
  error: string;
  data: T;
}

export type PromData = PromMatrixData | PromVectorData | PromScalarData | PromExemplarData[] | null;

export interface Labels {
  [index: string]: any;
}

export interface ScrapeExemplar {
  exemplar: Exemplar;
  scrapeTimestamp: number;
}

export interface Exemplar {
  labels: Labels;
  value: number;
  timestamp: number;
  hasTimestamp: boolean;
}

export interface PromExemplarData {
  seriesLabels: PromMetric;
  exemplars: ScrapeExemplar[];
}

export interface PromVectorData {
  resultType: 'vector';
  result: Array<{
    metric: PromMetric;
    value: PromValue;
  }>;
}

export interface PromMatrixData {
  resultType: 'matrix';
  result: Array<{
    metric: PromMetric;
    values: PromValue[];
  }>;
}

export interface PromScalarData {
  resultType: 'scalar';
  result: PromValue;
}

export type PromValue = [number, any];

export interface PromMetric {
  __name__?: string;
  [index: string]: any;
}

export function isFetchErrorResponse(response: any): response is FetchError {
  return 'cancelled' in response;
}

export function isMatrixData(result: MatrixOrVectorResult): result is PromMatrixData['result'][0] {
  return 'values' in result;
}

export function isExemplarData(result: PromData): result is PromExemplarData[] {
  if (result == null || !Array.isArray(result)) {
    return false;
  }
  return 'exemplars' in (result as any)[0];
}

export type MatrixOrVectorResult = PromMatrixData['result'][0] | PromVectorData['result'][0];

export interface TransformOptions {
  format?: string;
  step?: number;
  legendFormat?: string;
  start: number;
  end: number;
  query: string;
  responseListLength: number;
  scopedVars?: ScopedVars;
  refId: string;
  valueWithRefId?: boolean;
  meta: QueryResultMeta;
}

export interface PromLabelQueryResponse {
  data: {
    status: string;
    data: string[];
  };
  cancelled?: boolean;
}
