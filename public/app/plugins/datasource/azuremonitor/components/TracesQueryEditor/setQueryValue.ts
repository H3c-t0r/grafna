import { AzureMonitorQuery, AzureTracesFilter, ResultFormat } from '../../types';

export function setQueryOperationId(query: AzureMonitorQuery, operationId?: string): AzureMonitorQuery {
  return {
    ...query,
    azureTraces: {
      ...query.azureTraces,
      operationId,
    },
  };
}

export function setFormatAs(query: AzureMonitorQuery, formatAs: ResultFormat): AzureMonitorQuery {
  return {
    ...query,
    azureTraces: {
      ...query.azureTraces,
      resultFormat: formatAs,
    },
  };
}

export function setTraceTypes(query: AzureMonitorQuery, traceTypes: string[]): AzureMonitorQuery {
  return {
    ...query,
    azureTraces: {
      ...query.azureTraces,
      traceTypes,
    },
  };
}

export function setFilters(query: AzureMonitorQuery, filters: AzureTracesFilter[]): AzureMonitorQuery {
  return {
    ...query,
    azureTraces: {
      ...query.azureTraces,
      filters,
    },
  };
}
