import React, { Context, createContext, FunctionComponent, useCallback, useContext } from 'react';
import { ElasticDatasource } from '../../datasource';
import { combineReducers, useStatelessReducer, DispatchContext } from '../../hooks/useStatelessReducer';
import { ElasticsearchQuery } from '../../types';

import { reducer as metricsReducer } from './MetricAggregationsEditor/state/reducer';
import { reducer as bucketAggsReducer } from './BucketAggregationsEditor/state/reducer';
import { aliasPatternReducer, queryReducer, initQuery } from './state';
import { TimeRange } from '@grafana/data';

const DatasourceContext = createContext<ElasticDatasource | undefined>(undefined);
const QueryContext = createContext<ElasticsearchQuery | undefined>(undefined);
const RangeContext = createContext<TimeRange | undefined>(undefined);

interface Props {
  query: ElasticsearchQuery;
  onChange: (query: ElasticsearchQuery) => void;
  onRunQuery: () => void;
  datasource: ElasticDatasource;
  range: TimeRange;
}

export const ElasticsearchProvider: FunctionComponent<Props> = ({
  children,
  onChange,
  onRunQuery,
  query,
  datasource,
  range,
}) => {
  const onStateChange = useCallback(
    (query: ElasticsearchQuery) => {
      onChange(query);
      onRunQuery();
    },
    [onChange, onRunQuery]
  );

  const reducer = combineReducers({
    query: queryReducer,
    alias: aliasPatternReducer,
    metrics: metricsReducer,
    bucketAggs: bucketAggsReducer,
  });

  const dispatch = useStatelessReducer(
    // timeField is part of the query model, but its value is always set to be the one from datasource settings.
    (newState) => onStateChange({ ...query, ...newState, timeField: datasource.timeField }),
    query,
    reducer
  );

  // This initializes the query by dispatching an init action to each reducer.
  // useStatelessReducer will then call `onChange` with the newly generated query
  if (!query.metrics || !query.bucketAggs || query.query === undefined) {
    dispatch(initQuery());

    return null;
  }

  return (
    <DatasourceContext.Provider value={datasource}>
      <QueryContext.Provider value={query}>
        <RangeContext.Provider value={range}>
          <DispatchContext.Provider value={dispatch}>{children}</DispatchContext.Provider>
        </RangeContext.Provider>
      </QueryContext.Provider>
    </DatasourceContext.Provider>
  );
};

interface GetHook {
  <T>(context: Context<T>): () => NonNullable<T>;
}

const getHook: GetHook = (c) => () => {
  const contextValue = useContext(c);

  if (!contextValue) {
    throw new Error('use ElasticsearchProvider first.');
  }

  return contextValue as NonNullable<typeof contextValue>;
};

export const useQuery = getHook(QueryContext);
export const useDatasource = getHook(DatasourceContext);
export const useRange = getHook(RangeContext);
