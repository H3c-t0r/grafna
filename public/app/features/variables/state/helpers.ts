import { combineReducers } from '@reduxjs/toolkit';
import { LoadingState } from '@grafana/data';

import { VariablesState } from './types';
import {
  DashboardVariableModel,
  initialVariableModelState,
  OrgVariableModel,
  UserVariableModel,
  VariableHide,
  VariableModel,
} from '../types';

import { VariableAdapter } from '../adapters';
import { dashboardReducer } from 'app/features/dashboard/state/reducers';
import { DashboardState, StoreState } from '../../../types';
import { NEW_VARIABLE_ID } from '../constants';
import { keyedVariablesReducer, DashboardVariablesState } from './keyedVariablesReducer';
import { getInitialTemplatingState, TemplatingState } from './reducers';

export const getVariableState = (
  noOfVariables: number,
  inEditorIndex = -1,
  includeEmpty = false,
  includeSystem = false
): Record<string, VariableModel> => {
  const variables: Record<string, VariableModel> = {};

  if (includeSystem) {
    const dashboardModel: DashboardVariableModel = {
      ...initialVariableModelState,
      id: '__dashboard',
      name: '__dashboard',
      type: 'system',
      index: -3,
      skipUrlSync: true,
      hide: VariableHide.hideVariable,
      current: {
        value: {
          name: 'A dashboard title',
          uid: 'An dashboard UID',
          toString: () => 'A dashboard title',
        },
      },
    };

    const orgModel: OrgVariableModel = {
      ...initialVariableModelState,
      id: '__org',
      name: '__org',
      type: 'system',
      index: -2,
      skipUrlSync: true,
      hide: VariableHide.hideVariable,
      current: {
        value: {
          name: 'An org name',
          id: 1,
          toString: () => '1',
        },
      },
    };

    const userModel: UserVariableModel = {
      ...initialVariableModelState,
      id: '__user',
      name: '__user',
      type: 'system',
      index: -1,
      skipUrlSync: true,
      hide: VariableHide.hideVariable,
      current: {
        value: {
          login: 'admin',
          id: 1,
          email: 'admin@test',
          toString: () => '1',
        },
      },
    };

    variables[dashboardModel.id] = dashboardModel;
    variables[orgModel.id] = orgModel;
    variables[userModel.id] = userModel;
  }

  for (let index = 0; index < noOfVariables; index++) {
    variables[index] = {
      id: index.toString(),
      dashboardUid: 'uid',
      type: 'query',
      name: `Name-${index}`,
      hide: VariableHide.dontHide,
      index,
      label: `Label-${index}`,
      skipUrlSync: false,
      global: false,
      state: LoadingState.NotStarted,
      error: null,
      description: null,
    };
  }

  if (includeEmpty) {
    variables[NEW_VARIABLE_ID] = {
      id: NEW_VARIABLE_ID,
      dashboardUid: 'uid',
      type: 'query',
      name: `Name-${NEW_VARIABLE_ID}`,
      hide: VariableHide.dontHide,
      index: noOfVariables,
      label: `Label-${NEW_VARIABLE_ID}`,
      skipUrlSync: false,
      global: false,
      state: LoadingState.NotStarted,
      error: null,
      description: null,
    };
  }

  return variables;
};

export const getVariableTestContext = <Model extends VariableModel>(
  adapter: VariableAdapter<Model>,
  variableOverrides: Partial<Model> = {}
) => {
  const defaultVariable = {
    ...adapter.initialState,
    id: '0',
    dashboardUid: 'uid',
    index: 0,
    name: '0',
  };

  const initialState: VariablesState = {
    '0': { ...defaultVariable, ...variableOverrides },
  };

  return { initialState };
};

export const getRootReducer = () =>
  combineReducers({
    dashboard: dashboardReducer,
    dashboardVariables: keyedVariablesReducer,
  });

export type RootReducerType = { dashboard: DashboardState; dashboardVariables: DashboardVariablesState };

export const getTemplatingRootReducer = () =>
  combineReducers({
    dashboardVariables: keyedVariablesReducer,
  });

export type TemplatingReducerType = { dashboardVariables: DashboardVariablesState };

export function getPreloadedState(
  key: string,
  templatingState: Partial<TemplatingState>
): Pick<StoreState, 'dashboardVariables'> {
  return {
    dashboardVariables: {
      lastKey: key,
      keys: {
        [key]: {
          ...getInitialTemplatingState(),
          ...templatingState,
        },
      },
    },
  };
}
