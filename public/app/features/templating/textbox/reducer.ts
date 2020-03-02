import { createReducer } from '@reduxjs/toolkit';

import { TextBoxVariableModel, VariableHide, VariableOption } from '../variable';
import {
  addVariable,
  changeVariableProp,
  removeInitLock,
  resolveInitLock,
  setCurrentVariableValue,
} from '../state/actions';
import { Deferred } from '../deferred';
import { ALL_VARIABLE_VALUE, emptyUuid, getInstanceState, VariableState } from '../state/types';
import { initialTemplatingState } from '../state/reducers';
import cloneDeep from 'lodash/cloneDeep';
import { createTextBoxOptions } from './actions';

export interface TextBoxVariableState extends VariableState<TextBoxVariableModel> {}

export const initialTextBoxVariableModelState: TextBoxVariableModel = {
  uuid: emptyUuid,
  global: false,
  index: -1,
  type: 'textbox',
  name: '',
  label: '',
  hide: VariableHide.dontHide,
  query: '',
  current: {} as VariableOption,
  options: [],
  skipUrlSync: false,
  initLock: null,
};

export const initialTextBoxVariableState: TextBoxVariableState = {
  variable: initialTextBoxVariableModelState,
};

export const textBoxVariableReducer = createReducer(initialTemplatingState, builder =>
  builder
    .addCase(addVariable, (state, action) => {
      state.variables[action.payload.uuid!] = cloneDeep(initialTextBoxVariableState);
      state.variables[action.payload.uuid!].variable = {
        ...state.variables[action.payload.uuid!].variable,
        ...action.payload.data.model,
      };
      state.variables[action.payload.uuid!].variable.uuid = action.payload.uuid;
      state.variables[action.payload.uuid!].variable.index = action.payload.data.index;
      state.variables[action.payload.uuid!].variable.global = action.payload.data.global;
      state.variables[action.payload.uuid!].variable.initLock = new Deferred();
    })
    .addCase(setCurrentVariableValue, (state, action) => {
      const instanceState = getInstanceState<TextBoxVariableState>(state, action.payload.uuid);
      const current = { ...action.payload.data };

      if (Array.isArray(current.text) && current.text.length > 0) {
        current.text = current.text.join(' + ');
      } else if (Array.isArray(current.value) && current.value[0] !== ALL_VARIABLE_VALUE) {
        current.text = current.value.join(' + ');
      }

      instanceState.variable.current = current;
      instanceState.variable.options = instanceState.variable.options.map(option => {
        let selected = false;
        if (Array.isArray(current.value)) {
          for (let index = 0; index < current.value.length; index++) {
            const value = current.value[index];
            if (option.value === value) {
              selected = true;
              break;
            }
          }
        } else if (option.value === current.value) {
          selected = true;
        }
        option.selected = selected;
        return option;
      });
    })
    .addCase(resolveInitLock, (state, action) => {
      const instanceState = getInstanceState<TextBoxVariableState>(state, action.payload.uuid!);
      instanceState.variable.initLock?.resolve();
    })
    .addCase(removeInitLock, (state, action) => {
      const instanceState = getInstanceState<TextBoxVariableState>(state, action.payload.uuid!);
      instanceState.variable.initLock = null;
    })
    .addCase(changeVariableProp, (state, action) => {
      const instanceState = getInstanceState<TextBoxVariableState>(state, action.payload.uuid!);
      (instanceState.variable as Record<string, any>)[action.payload.data.propName] = action.payload.data.propValue;
    })
    .addCase(createTextBoxOptions, (state, action) => {
      const instanceState = getInstanceState<TextBoxVariableState>(state, action.payload.uuid!);
      instanceState.variable.options = [
        { text: instanceState.variable.query.trim(), value: instanceState.variable.query.trim(), selected: false },
      ];
      instanceState.variable.current = instanceState.variable.options[0];
    })
);
