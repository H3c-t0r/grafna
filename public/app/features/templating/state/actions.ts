import castArray from 'lodash/castArray';
import { v4 } from 'uuid';
import { createAction, PrepareAction } from '@reduxjs/toolkit';
import { UrlQueryMap, UrlQueryValue } from '@grafana/runtime';

import {
  QueryVariableModel,
  VariableModel,
  VariableOption,
  VariableRefresh,
  VariableType,
  VariableWithOptions,
} from '../variable';
import { StoreState, ThunkResult } from '../../../types';
import { getVariable, getVariables } from './selectors';
import { variableAdapters } from '../adapters';
import { Graph } from '../../../core/utils/dag';
import { updateLocation } from 'app/core/actions';

// process flow queryVariable
// thunk => processVariables
//    adapter => setValueFromUrl
//      thunk => setOptionFromUrl
//        adapter => updateOptions
//          thunk => updateQueryVariableOptions
//            action => updateVariableOptions
//            action => updateVariableTags
//            thunk => validateVariableSelectionState
//              adapter => setValue
//                thunk => setOptionAsCurrent
//                  action => setCurrentVariableValue
//                  thunk => variableUpdated
//                    adapter => updateOptions for dependent nodes
//        adapter => setValue
//          thunk => setOptionAsCurrent
//            action => setCurrentVariableValue
//            thunk => variableUpdated
//              adapter => updateOptions for dependent nodes
//    adapter => updateOptions
//      thunk => updateQueryVariableOptions
//        action => updateVariableOptions
//        action => updateVariableTags
//        thunk => validateVariableSelectionState
//          adapter => setValue
//            thunk => setOptionAsCurrent
//              action => setCurrentVariableValue
//              thunk => variableUpdated
//                adapter => updateOptions for dependent nodes

export interface VariableIdentifier {
  type: VariableType;
  uuid: string;
}

export interface VariablePayload<T extends any = undefined> extends VariableIdentifier {
  data: T;
}

export interface AddVariable<T extends VariableModel = VariableModel> {
  global: boolean; // part of dashboard or global
  index: number; // the order in variables list
  model: T;
}

export const addVariable = createAction<PrepareAction<VariablePayload<AddVariable>>>(
  'templating/addVariable',
  (payload: VariablePayload<AddVariable>) => {
    return {
      payload: {
        ...payload,
        uuid: payload.uuid ?? v4(), // for testing purposes we allow to pass existing uuid
      },
    };
  }
);

export const removeVariable = createAction<VariablePayload>('templating/removeVariable');
export const storeNewVariable = createAction<VariablePayload>('templating/storeNewVariable');
export interface DuplicateVariable {
  newUuid: string;
}
export const duplicateVariable = createAction<PrepareAction<VariablePayload<DuplicateVariable>>>(
  'templating/duplicateVariable',
  (payload: VariablePayload<DuplicateVariable>) => {
    return {
      payload: {
        ...payload,
        data: {
          ...payload.data,
          newUuid: payload.data.newUuid ?? v4(), // for testing purposes we allow to pass existing newUuid
        },
      },
    };
  }
);
export const changeVariableOrder = createAction<VariablePayload<{ fromIndex: number; toIndex: number }>>(
  'templating/changeVariableOrder'
);
export const changeToEditorListMode = createAction<VariablePayload>('templating/changeToEditorListMode');
export const changeToEditorEditMode = createAction<VariablePayload>('templating/changeToEditorEditMode');
export const resolveInitLock = createAction<VariablePayload>('templating/resolveInitLock');
export const removeInitLock = createAction<VariablePayload>('templating/removeInitLock');
export const setCurrentVariableValue = createAction<VariablePayload<VariableOption>>(
  'templating/setCurrentVariableValue'
);

export const changeVariableType = createAction<VariablePayload<VariableType>>('templating/changeVariableType');
export const updateVariableOptions = createAction<VariablePayload<any[]>>('templating/updateVariableOptions');
export const updateVariableTags = createAction<VariablePayload<any[]>>('templating/updateVariableTags');
export const changeVariableProp = createAction<VariablePayload<{ propName: string; propValue: any }>>(
  'templating/changeVariableProp'
);

export const toVariableIdentifier = (variable: VariableModel): VariableIdentifier => {
  return { type: variable.type, uuid: variable.uuid! };
};

export function toVariablePayload<T extends any = undefined>(
  identifier: VariableIdentifier,
  data?: T
): VariablePayload<T>;

export function toVariablePayload<T extends any = undefined>(model: VariableModel, data?: T): VariablePayload<T>;

export function toVariablePayload<T extends any = undefined>(
  obj: VariableIdentifier | VariableModel,
  data?: T
): VariablePayload<T> {
  return { type: obj.type, uuid: obj.uuid!, data: data as T };
}

export const initDashboardTemplating = (list: VariableModel[]): ThunkResult<void> => {
  return (dispatch, getState) => {
    for (let index = 0; index < list.length; index++) {
      const model = list[index];
      if (!variableAdapters.contains(model.type)) {
        continue;
      }

      dispatch(addVariable(toVariablePayload(model, { global: false, index, model })));
    }
  };
};

export const processVariableDependencies = async (variable: VariableModel, state: StoreState) => {
  let dependencies: Array<Promise<any>> = [];

  for (const otherVariable of getVariables(state)) {
    if (variable === otherVariable) {
      continue;
    }

    if (variableAdapters.contains(variable.type)) {
      if (variableAdapters.get(variable.type).dependsOn(variable, otherVariable)) {
        dependencies.push(otherVariable.initLock!.promise);
      }
    }
  }

  await Promise.all(dependencies);
};

export const processVariable = (identifier: VariableIdentifier, queryParams: UrlQueryMap): ThunkResult<void> => {
  return async (dispatch, getState) => {
    const variable = getVariable(identifier.uuid!, getState());
    await processVariableDependencies(variable, getState());

    const urlValue = queryParams['var-' + variable.name];
    if (urlValue !== void 0) {
      await variableAdapters.get(variable.type).setValueFromUrl(variable, urlValue ?? '');
    }

    if (variable.hasOwnProperty('refresh')) {
      const refreshableVariable = variable as QueryVariableModel;
      if (
        refreshableVariable.refresh === VariableRefresh.onDashboardLoad ||
        refreshableVariable.refresh === VariableRefresh.onTimeRangeChanged
      ) {
        await variableAdapters.get(variable.type).updateOptions(refreshableVariable);
      }
    }

    await dispatch(resolveInitLock(toVariablePayload(variable)));
  };
};

export const processVariables = (): ThunkResult<void> => {
  return async (dispatch, getState) => {
    const queryParams = getState().location.query;
    const promises = getVariables(getState()).map(
      async (variable: VariableModel) => await dispatch(processVariable(toVariableIdentifier(variable), queryParams))
    );

    await Promise.all(promises);

    for (let index = 0; index < getVariables(getState()).length; index++) {
      await dispatch(removeInitLock(toVariablePayload(getVariables(getState())[index])));
    }
  };
};

export const setOptionFromUrl = (identifier: VariableIdentifier, urlValue: UrlQueryValue): ThunkResult<void> => {
  return async (dispatch, getState) => {
    const variable = getVariable(identifier.uuid!, getState());
    if (variable.hasOwnProperty('refresh') && (variable as QueryVariableModel).refresh !== VariableRefresh.never) {
      // updates options
      await variableAdapters.get(variable.type).updateOptions(variable);
    }

    // get variable from state
    const variableFromState = getVariable<VariableWithOptions>(variable.uuid!, getState());
    if (!variableFromState) {
      throw new Error(`Couldn't find variable with name: ${variable.name}`);
    }
    // Simple case. Value in url matches existing options text or value.
    let option = variableFromState.options.find(op => {
      return op.text === urlValue || op.value === urlValue;
    });

    if (!option) {
      let defaultText = urlValue as string | string[];
      const defaultValue = urlValue as string | string[];

      if (Array.isArray(urlValue)) {
        // Multiple values in the url. We construct text as a list of texts from all matched options.
        const urlValueArray = urlValue as string[];
        defaultText = urlValueArray.reduce((acc: string[], item: string) => {
          const foundOption = variableFromState.options.find(o => o.value === item);
          if (!foundOption) {
            // @ts-ignore according to strict null errors this can never happen
            // TODO: investigate this further or refactor code
            return [].concat(acc, [item]);
          }

          // @ts-ignore according to strict null errors this can never happen
          // TODO: investigate this further or refactor code
          return [].concat(acc, [foundOption.text]);
        }, []);
      }

      // It is possible that we did not match the value to any existing option. In that case the url value will be
      // used anyway for both text and value.
      option = { text: defaultText, value: defaultValue, selected: false };
    }

    if (variableFromState.hasOwnProperty('multi')) {
      // In case variable is multiple choice, we cast to array to preserve the same behaviour as when selecting
      // the option directly, which will return even single value in an array.
      option = { text: castArray(option.text), value: castArray(option.value), selected: false };
    }

    await variableAdapters.get(variable.type).setValue(variableFromState, option);
  };
};

export const selectOptionsForCurrentValue = (variable: VariableWithOptions): VariableOption[] => {
  let i, y, value, option;
  const selected: VariableOption[] = [];

  for (i = 0; i < variable.options.length; i++) {
    option = { ...variable.options[i] };
    option.selected = false;
    if (Array.isArray(variable.current.value)) {
      for (y = 0; y < variable.current.value.length; y++) {
        value = variable.current.value[y];
        if (option.value === value) {
          option.selected = true;
          selected.push(option);
        }
      }
    } else if (option.value === variable.current.value) {
      option.selected = true;
      selected.push(option);
    }
  }

  return selected;
};

export const validateVariableSelectionState = (
  identifier: VariableIdentifier,
  defaultValue?: string
): ThunkResult<void> => {
  return async (dispatch, getState) => {
    const variableInState = getVariable<VariableWithOptions>(identifier.uuid!, getState());
    const setValue = variableAdapters.get(variableInState.type).setValue;
    if (!variableInState.current) {
      return setValue(variableInState, {} as VariableOption);
    }

    if (Array.isArray(variableInState.current.value)) {
      const selected = selectOptionsForCurrentValue(variableInState);

      // if none pick first
      if (selected.length === 0) {
        const option = variableInState.options[0];
        return setValue(variableInState, option);
      }

      const option: VariableOption = {
        value: selected.map(v => v.value) as string[],
        text: selected.map(v => v.text) as string[],
        selected: true,
      };
      return setValue(variableInState, option);
    }

    let option: VariableOption | undefined | null = null;

    // 1. find the current value
    option = variableInState.options.find(v => v.text === variableInState.current.text);
    if (option) {
      return setValue(variableInState, option);
    }

    // 2. find the default value
    if (defaultValue) {
      option = variableInState.options.find(v => v.text === defaultValue);
      if (option) {
        return setValue(variableInState, option);
      }
    }

    // 3. use the first value
    if (variableInState.options) {
      const option = variableInState.options[0];
      return setValue(variableInState, option);
    }

    // 4... give up
    return Promise.resolve();
  };
};

export const setOptionAsCurrent = (
  identifier: VariableIdentifier,
  current: VariableOption,
  emitChanges: boolean
): ThunkResult<void> => {
  return async (dispatch, getState) => {
    dispatch(setCurrentVariableValue(toVariablePayload(identifier, current)));
    return dispatch(variableUpdated(identifier, emitChanges));
  };
};

const createGraph = (variables: VariableModel[]) => {
  const g = new Graph();

  variables.forEach(v => {
    g.createNode(v.name);
  });

  variables.forEach(v1 => {
    variables.forEach(v2 => {
      if (v1 === v2) {
        return;
      }

      if (variableAdapters.get(v1.type).dependsOn(v1, v2)) {
        g.link(v1.name, v2.name);
      }
    });
  });

  return g;
};

export const variableUpdated = (identifier: VariableIdentifier, emitChangeEvents: boolean): ThunkResult<void> => {
  return (dispatch, getState) => {
    // if there is a variable lock ignore cascading update because we are in a boot up scenario
    const variable = getVariable(identifier.uuid!, getState());
    if (variable.initLock) {
      return Promise.resolve();
    }

    const variables = getVariables(getState());
    const g = createGraph(variables);

    const node = g.getNode(variable.name);
    let promises: Array<Promise<any>> = [];
    if (node) {
      promises = node.getOptimizedInputEdges().map(e => {
        const variable = variables.find(v => v.name === e.inputNode.name);
        if (!variable) {
          return Promise.resolve();
        }
        return variableAdapters.get(variable.type).updateOptions(variable);
      });
    }

    return Promise.all(promises).then(() => {
      if (emitChangeEvents) {
        const dashboard = getState().dashboard.getModel();
        dashboard?.processRepeats();
        dispatch(updateLocation({ query: getQueryWithVariables(getState) }));
        dashboard?.startRefresh();
      }
    });
  };
};

const getQueryWithVariables = (getState: () => StoreState): UrlQueryMap => {
  const queryParams = getState().location.query;

  const queryParamsNew = Object.keys(queryParams)
    .filter(key => key.indexOf('var-') === -1)
    .reduce((obj, key) => {
      obj[key] = queryParams[key];
      return obj;
    }, {} as UrlQueryMap);

  for (const variable of getVariables(getState())) {
    if (variable.skipUrlSync) {
      continue;
    }

    const adapter = variableAdapters.get(variable.type);
    queryParamsNew['var-' + variable.name] = adapter.getValueForUrl(variable);
  }

  return queryParamsNew;
};
