import { Subscription } from 'rxjs';
import { getDataSourceSrv, toDataQueryError } from '@grafana/runtime';
import { DataSourceRef } from '@grafana/data';

import { updateOptions } from '../state/actions';
import { QueryVariableModel } from '../types';
import { ThunkResult } from '../../../types';
import { getVariable, getVariablesState } from '../state/selectors';
import {
  addVariableEditorError,
  changeVariableEditorExtended,
  removeVariableEditorError,
  VariableEditorState,
} from '../editor/reducer';
import { changeVariableProp } from '../state/sharedReducer';
import { KeyedVariableIdentifier } from '../state/types';
import { getVariableQueryEditor } from '../editor/getVariableQueryEditor';
import { getVariableQueryRunner } from './VariableQueryRunner';
import { variableQueryObserver } from './variableQueryObserver';
import { QueryVariableEditorState } from './reducer';
import { hasOngoingTransaction, toKeyedVariableIdentifier, toVariablePayload } from '../utils';
import { toKeyedAction } from '../state/keyedVariablesReducer';

export const updateQueryVariableOptions = (
  identifier: KeyedVariableIdentifier,
  searchFilter?: string
): ThunkResult<void> => {
  return async (dispatch, getState) => {
    try {
      const { rootStateKey } = identifier;
      if (!hasOngoingTransaction(rootStateKey, getState())) {
        // we might have cancelled a batch so then variable state is removed
        return;
      }

      const variableInState = getVariable<QueryVariableModel>(identifier, getState());
      if (getVariablesState(rootStateKey, getState()).editor.id === variableInState.id) {
        dispatch(toKeyedAction(rootStateKey, removeVariableEditorError({ errorProp: 'update' })));
      }
      const datasource = await getDataSourceSrv().get(variableInState.datasource ?? '');

      // We need to await the result from variableQueryRunner before moving on otherwise variables dependent on this
      // variable will have the wrong current value as input
      await new Promise((resolve, reject) => {
        const subscription: Subscription = new Subscription();
        const observer = variableQueryObserver(resolve, reject, subscription);
        const responseSubscription = getVariableQueryRunner().getResponse(identifier).subscribe(observer);
        subscription.add(responseSubscription);

        getVariableQueryRunner().queueRequest({ identifier, datasource, searchFilter });
      });
    } catch (err) {
      const error = toDataQueryError(err);
      const { rootStateKey } = identifier;
      if (getVariablesState(rootStateKey, getState()).editor.id === identifier.id) {
        dispatch(
          toKeyedAction(rootStateKey, addVariableEditorError({ errorProp: 'update', errorText: error.message }))
        );
      }

      throw error;
    }
  };
};

export const initQueryVariableEditor = (identifier: KeyedVariableIdentifier): ThunkResult<void> => async (
  dispatch,
  getState
) => {
  const variable = getVariable<QueryVariableModel>(identifier, getState());
  await dispatch(changeQueryVariableDataSource(toKeyedVariableIdentifier(variable), variable.datasource));
};

export const changeQueryVariableDataSource = (
  identifier: KeyedVariableIdentifier,
  name: DataSourceRef | null
): ThunkResult<void> => {
  return async (dispatch, getState) => {
    try {
      const { rootStateKey } = identifier;
      const state = getVariablesState(rootStateKey, getState()).editor as VariableEditorState<QueryVariableEditorState>;
      const previousDatasource = state.extended?.dataSource;
      const dataSource = await getDataSourceSrv().get(name ?? '');
      if (previousDatasource && previousDatasource.type !== dataSource?.type) {
        dispatch(
          toKeyedAction(
            rootStateKey,
            changeVariableProp(toVariablePayload(identifier, { propName: 'query', propValue: '' }))
          )
        );
      }
      dispatch(
        toKeyedAction(rootStateKey, changeVariableEditorExtended({ propName: 'dataSource', propValue: dataSource }))
      );

      const VariableQueryEditor = await getVariableQueryEditor(dataSource);
      dispatch(
        toKeyedAction(
          rootStateKey,
          changeVariableEditorExtended({ propName: 'VariableQueryEditor', propValue: VariableQueryEditor })
        )
      );
    } catch (err) {
      console.error(err);
    }
  };
};

export const changeQueryVariableQuery = (
  identifier: KeyedVariableIdentifier,
  query: any,
  definition?: string
): ThunkResult<void> => async (dispatch, getState) => {
  const { rootStateKey } = identifier;
  const variableInState = getVariable<QueryVariableModel>(identifier, getState());
  if (hasSelfReferencingQuery(variableInState.name, query)) {
    const errorText = 'Query cannot contain a reference to itself. Variable: $' + variableInState.name;
    dispatch(toKeyedAction(rootStateKey, addVariableEditorError({ errorProp: 'query', errorText })));
    return;
  }

  dispatch(toKeyedAction(rootStateKey, removeVariableEditorError({ errorProp: 'query' })));
  dispatch(
    toKeyedAction(
      rootStateKey,
      changeVariableProp(toVariablePayload(identifier, { propName: 'query', propValue: query }))
    )
  );

  if (definition) {
    dispatch(
      toKeyedAction(
        rootStateKey,
        changeVariableProp(toVariablePayload(identifier, { propName: 'definition', propValue: definition }))
      )
    );
  } else if (typeof query === 'string') {
    dispatch(
      toKeyedAction(
        rootStateKey,
        changeVariableProp(toVariablePayload(identifier, { propName: 'definition', propValue: query }))
      )
    );
  }

  await dispatch(updateOptions(identifier));
};

export function hasSelfReferencingQuery(name: string, query: any): boolean {
  if (typeof query === 'string' && query.match(new RegExp('\\$' + name + '(/| |$)'))) {
    return true;
  }

  const flattened = flattenQuery(query);

  for (let prop in flattened) {
    if (flattened.hasOwnProperty(prop)) {
      const value = flattened[prop];
      if (typeof value === 'string' && value.match(new RegExp('\\$' + name + '(/| |$)'))) {
        return true;
      }
    }
  }

  return false;
}

/*
 * Function that takes any object and flattens all props into one level deep object
 * */
export function flattenQuery(query: any): any {
  if (typeof query !== 'object') {
    return { query };
  }

  const keys = Object.keys(query);
  const flattened = keys.reduce((all, key) => {
    const value = query[key];
    if (typeof value !== 'object') {
      all[key] = value;
      return all;
    }

    const result = flattenQuery(value);
    for (let childProp in result) {
      if (result.hasOwnProperty(childProp)) {
        all[`${key}_${childProp}`] = result[childProp];
      }
    }

    return all;
  }, {} as Record<string, any>);

  return flattened;
}
