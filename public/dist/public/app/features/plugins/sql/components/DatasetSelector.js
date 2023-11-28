import { __awaiter } from "tslib";
import React, { useEffect } from 'react';
import { useAsync } from 'react-use';
import { Select } from '@grafana/ui';
import { toOption } from '../types';
import { isSqlDatasourceDatabaseSelectionFeatureFlagEnabled } from './QueryEditorFeatureFlag.utils';
export const DatasetSelector = ({ dataset, db, isPostgresInstance, onChange, preconfiguredDataset, }) => {
    /*
      The behavior of this component - for MSSQL and MySQL datasources - is based on whether the user chose to create a datasource
      with or without a default database (preconfiguredDataset). If the user configured a default database, this selector
      should only allow that single preconfigured database option to be selected. If the user chose to NOT assign/configure a default database,
      then the user should be able to use this component to choose between multiple databases available to the datasource.
    */
    // `hasPreconfigCondition` is true if either 1) the sql datasource has a preconfigured default database,
    // OR if 2) the datasource is Postgres. In either case the only option available to the user is the preconfigured database.
    const hasPreconfigCondition = !!preconfiguredDataset || isPostgresInstance;
    const state = useAsync(() => __awaiter(void 0, void 0, void 0, function* () {
        if (isSqlDatasourceDatabaseSelectionFeatureFlagEnabled()) {
            // If a default database is already configured for a MSSQL or MySQL data source, OR the data source is Postgres, no need to fetch other databases.
            if (hasPreconfigCondition) {
                // Set the current database to the preconfigured database.
                onChange(toOption(preconfiguredDataset));
                return [toOption(preconfiguredDataset)];
            }
        }
        // If there is no preconfigured database, but there is a selected dataset, set the current database to the selected dataset.
        if (dataset) {
            onChange(toOption(dataset));
        }
        // Otherwise, fetch all databases available to the datasource.
        const datasets = yield db.datasets();
        return datasets.map(toOption);
    }), []);
    useEffect(() => {
        if (!isSqlDatasourceDatabaseSelectionFeatureFlagEnabled()) {
            // Set default dataset when values are fetched
            if (!dataset) {
                if (state.value && state.value[0]) {
                    onChange(state.value[0]);
                }
            }
            else {
                if (state.value && state.value.find((v) => v.value === dataset) === undefined) {
                    // if value is set and newly fetched values does not contain selected value
                    if (state.value.length > 0) {
                        onChange(state.value[0]);
                    }
                }
            }
        }
    }, [state.value, onChange, dataset]);
    return (React.createElement(Select, { "aria-label": "Dataset selector", value: dataset, options: state.value, onChange: onChange, disabled: state.loading, isLoading: state.loading, menuShouldPortal: true }));
};
//# sourceMappingURL=DatasetSelector.js.map