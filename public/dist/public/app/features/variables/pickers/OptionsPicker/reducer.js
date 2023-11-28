import { createSlice } from '@reduxjs/toolkit';
import { cloneDeep, isString, trimStart } from 'lodash';
import { containsSearchFilter } from '@grafana/data';
import { applyStateChanges } from '../../../../core/utils/applyStateChanges';
import { ALL_VARIABLE_VALUE } from '../../constants';
import { isMulti, isQuery } from '../../guard';
export const initialOptionPickerState = {
    id: '',
    highlightIndex: -1,
    queryValue: '',
    selectedValues: [],
    options: [],
    multi: false,
};
export const OPTIONS_LIMIT = 1000;
const optionsToRecord = (options) => {
    if (!Array.isArray(options)) {
        return {};
    }
    return options.reduce((all, option) => {
        if (isString(option.value)) {
            all[option.value] = option;
        }
        return all;
    }, {});
};
const updateOptions = (state) => {
    if (!Array.isArray(state.options)) {
        state.options = [];
        return state;
    }
    const selectedOptions = optionsToRecord(state.selectedValues);
    state.selectedValues = Object.values(selectedOptions);
    state.options = state.options.map((option) => {
        if (!isString(option.value)) {
            return option;
        }
        const selected = !!selectedOptions[option.value];
        if (option.selected === selected) {
            return option;
        }
        return Object.assign(Object.assign({}, option), { selected });
    });
    state.options = applyLimit(state.options);
    return state;
};
const applyLimit = (options) => {
    if (!Array.isArray(options)) {
        return [];
    }
    if (options.length <= OPTIONS_LIMIT) {
        return options;
    }
    return options.slice(0, OPTIONS_LIMIT);
};
const updateDefaultSelection = (state) => {
    const { options, selectedValues } = state;
    if (options.length === 0 || selectedValues.length > 0) {
        return state;
    }
    if (!options[0] || options[0].value !== ALL_VARIABLE_VALUE) {
        return state;
    }
    state.selectedValues = [Object.assign(Object.assign({}, options[0]), { selected: true })];
    return state;
};
const updateAllSelection = (state) => {
    const { selectedValues } = state;
    if (selectedValues.length > 1) {
        state.selectedValues = selectedValues.filter((option) => option.value !== ALL_VARIABLE_VALUE);
    }
    return state;
};
// Utility function to select all options except 'ALL_VARIABLE_VALUE'
const selectAllOptions = (options) => options
    .filter((option) => option.value !== ALL_VARIABLE_VALUE)
    .map((option) => (Object.assign(Object.assign({}, option), { selected: true })));
const optionsPickerSlice = createSlice({
    name: 'templating/optionsPicker',
    initialState: initialOptionPickerState,
    reducers: {
        showOptions: (state, action) => {
            var _a;
            const { query, options } = action.payload;
            state.highlightIndex = -1;
            state.options = cloneDeep(options);
            state.id = action.payload.id;
            state.queryValue = '';
            state.multi = false;
            if (isMulti(action.payload)) {
                state.multi = (_a = action.payload.multi) !== null && _a !== void 0 ? _a : false;
            }
            if (isQuery(action.payload)) {
                const { queryValue } = action.payload;
                const queryHasSearchFilter = containsSearchFilter(query);
                state.queryValue = queryHasSearchFilter && queryValue ? queryValue : '';
            }
            state.selectedValues = state.options.filter((option) => option.selected);
            return applyStateChanges(state, updateDefaultSelection, updateOptions);
        },
        hideOptions: (state, action) => {
            return Object.assign({}, initialOptionPickerState);
        },
        toggleOption: (state, action) => {
            const { option, clearOthers, forceSelect } = action.payload;
            const { multi, selectedValues } = state;
            if (option) {
                const selected = !selectedValues.find((o) => o.value === option.value && o.text === option.text);
                if (option.value === ALL_VARIABLE_VALUE || !multi || clearOthers) {
                    if (selected || forceSelect) {
                        state.selectedValues = [Object.assign(Object.assign({}, option), { selected: true })];
                    }
                    else {
                        state.selectedValues = [];
                    }
                    return applyStateChanges(state, updateDefaultSelection, updateAllSelection, updateOptions);
                }
                if (forceSelect || selected) {
                    state.selectedValues.push(Object.assign(Object.assign({}, option), { selected: true }));
                    return applyStateChanges(state, updateDefaultSelection, updateAllSelection, updateOptions);
                }
                state.selectedValues = selectedValues.filter((o) => o.value !== option.value && o.text !== option.text);
            }
            else {
                state.selectedValues = [];
            }
            return applyStateChanges(state, updateDefaultSelection, updateAllSelection, updateOptions);
        },
        moveOptionsHighlight: (state, action) => {
            let nextIndex = state.highlightIndex + action.payload;
            if (nextIndex < 0) {
                nextIndex = -1;
            }
            else if (nextIndex >= state.options.length) {
                nextIndex = state.options.length - 1;
            }
            return Object.assign(Object.assign({}, state), { highlightIndex: nextIndex });
        },
        /**
         * Toggle the 'All' option or clear selections in the Options Picker dropdown.
         * 1. If 'All' is configured but not selected, and some other options are selected, it deselects all other options and selects only 'All'.
         * 2. If only 'All' is selected, it deselects 'All' and selects all other available options.
         * 3. If some options are selected but 'All' is not configured in the variable,
         *    it clears all selections and defaults to the current behavior for scenarios where 'All' is not configured.
         * 4. If no options are selected, it selects all available options.
         */
        toggleAllOptions: (state, action) => {
            var _a;
            // Check if 'All' option is configured by the user and if it's selected in the dropdown
            const isAllSelected = state.selectedValues.find((option) => option.value === ALL_VARIABLE_VALUE);
            const allOptionConfigured = state.options.find((option) => option.value === ALL_VARIABLE_VALUE);
            // If 'All' option is not selected from the dropdown, but some options are, clear all options and select 'All'
            if (state.selectedValues.length > 0 && !!allOptionConfigured && !isAllSelected) {
                state.selectedValues = [];
                state.selectedValues.push({
                    text: (_a = allOptionConfigured.text) !== null && _a !== void 0 ? _a : 'All',
                    value: allOptionConfigured.value,
                    selected: true,
                });
                return applyStateChanges(state, updateOptions);
            }
            // If 'All' option is the only one selected in the dropdown, unselect "All" and select each one of the other options.
            if (isAllSelected && state.selectedValues.length === 1) {
                state.selectedValues = selectAllOptions(state.options);
                return applyStateChanges(state, updateOptions);
            }
            // If some options are selected, but 'All' is not configured by the user, clear the selection and let the
            // current behavior when "All" does not exist and user clear the selected items.
            if (state.selectedValues.length > 0 && !allOptionConfigured) {
                state.selectedValues = [];
                return applyStateChanges(state, updateOptions);
            }
            // If no options are selected and 'All' is not selected, select all options
            state.selectedValues = selectAllOptions(state.options);
            return applyStateChanges(state, updateOptions);
        },
        updateSearchQuery: (state, action) => {
            state.queryValue = action.payload;
            return state;
        },
        updateOptionsAndFilter: (state, action) => {
            var _a;
            const searchQuery = trimStart(((_a = state.queryValue) !== null && _a !== void 0 ? _a : '').toLowerCase());
            state.options = action.payload.filter((option) => {
                var _a;
                const optionsText = (_a = option.text) !== null && _a !== void 0 ? _a : '';
                const text = Array.isArray(optionsText) ? optionsText.toString() : optionsText;
                return text.toLowerCase().indexOf(searchQuery) !== -1;
            });
            state.highlightIndex = 0;
            return applyStateChanges(state, updateDefaultSelection, updateOptions);
        },
        updateOptionsFromSearch: (state, action) => {
            state.options = action.payload;
            state.highlightIndex = 0;
            return applyStateChanges(state, updateDefaultSelection, updateOptions);
        },
        cleanPickerState: () => initialOptionPickerState,
    },
});
export const { toggleOption, showOptions, hideOptions, moveOptionsHighlight, toggleAllOptions, updateSearchQuery, updateOptionsAndFilter, updateOptionsFromSearch, cleanPickerState, } = optionsPickerSlice.actions;
export const optionsPickerReducer = optionsPickerSlice.reducer;
//# sourceMappingURL=reducer.js.map