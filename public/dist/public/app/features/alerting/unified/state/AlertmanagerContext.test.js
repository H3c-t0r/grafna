import { renderHook } from '@testing-library/react';
import { createMemoryHistory } from 'history';
import React from 'react';
import { MemoryRouter, Router } from 'react-router-dom';
import store from 'app/core/store';
import { AlertManagerImplementation } from 'app/plugins/datasource/alertmanager/types';
import * as useAlertManagerSources from '../hooks/useAlertManagerSources';
import { ALERTMANAGER_NAME_LOCAL_STORAGE_KEY } from '../utils/constants';
import { GRAFANA_RULES_SOURCE_NAME } from '../utils/datasource';
import { AlertmanagerProvider, isAlertManagerWithConfigAPI, useAlertmanager } from './AlertmanagerContext';
const grafanaAm = {
    name: GRAFANA_RULES_SOURCE_NAME,
    imgUrl: '',
};
const externalAmProm = {
    name: 'PrometheusAm',
    imgUrl: '',
};
const externalAmMimir = {
    name: 'MimirAm',
    imgUrl: '',
};
describe('useAlertmanager', () => {
    it('Should return undefined alert manager name when there are no available alert managers', () => {
        jest.spyOn(useAlertManagerSources, 'useAlertManagersByPermission').mockReturnValueOnce([]);
        const wrapper = ({ children }) => (React.createElement(MemoryRouter, null,
            React.createElement(AlertmanagerProvider, { accessType: "instance" }, children)));
        const { result } = renderHook(() => useAlertmanager(), { wrapper });
        expect(result.current.selectedAlertmanager).toBe(undefined);
    });
    it('Should return Grafana AM when it is available and no alert manager query param exists', () => {
        jest.spyOn(useAlertManagerSources, 'useAlertManagersByPermission').mockReturnValueOnce([grafanaAm]);
        const wrapper = ({ children }) => (React.createElement(MemoryRouter, null,
            React.createElement(AlertmanagerProvider, { accessType: "instance" }, children)));
        const { result } = renderHook(() => useAlertmanager(), { wrapper });
        expect(result.current.selectedAlertmanager).toBe(grafanaAm.name);
    });
    it('Should return alert manager included in the query param when available', () => {
        jest.spyOn(useAlertManagerSources, 'useAlertManagersByPermission').mockReturnValueOnce([externalAmProm]);
        const history = createMemoryHistory();
        history.push({ search: `alertmanager=${externalAmProm.name}` });
        const wrapper = ({ children }) => (React.createElement(Router, { history: history },
            React.createElement(AlertmanagerProvider, { accessType: "instance" }, children)));
        const { result } = renderHook(() => useAlertmanager(), { wrapper });
        expect(result.current.selectedAlertmanager).toBe(externalAmProm.name);
    });
    it('Should return undefined if alert manager included in the query is not available', () => {
        jest.spyOn(useAlertManagerSources, 'useAlertManagersByPermission').mockReturnValueOnce([]);
        const history = createMemoryHistory();
        history.push({ search: `alertmanager=Not available external AM` });
        const wrapper = ({ children }) => (React.createElement(Router, { history: history },
            React.createElement(AlertmanagerProvider, { accessType: "instance" }, children)));
        const { result } = renderHook(() => useAlertmanager(), { wrapper });
        expect(result.current.selectedAlertmanager).toBe(undefined);
    });
    it('Should return alert manager from store if available and query is empty', () => {
        jest.spyOn(useAlertManagerSources, 'useAlertManagersByPermission').mockReturnValueOnce([externalAmProm]);
        const wrapper = ({ children }) => (React.createElement(MemoryRouter, null,
            React.createElement(AlertmanagerProvider, { accessType: "instance" }, children)));
        store.set(ALERTMANAGER_NAME_LOCAL_STORAGE_KEY, externalAmProm.name);
        const { result } = renderHook(() => useAlertmanager(), { wrapper });
        expect(result.current.selectedAlertmanager).toBe(externalAmProm.name);
    });
    it('Should prioritize the alert manager from query over store', () => {
        jest
            .spyOn(useAlertManagerSources, 'useAlertManagersByPermission')
            .mockReturnValueOnce([externalAmProm, externalAmMimir]);
        const history = createMemoryHistory();
        history.push({ search: `alertmanager=${externalAmProm.name}` });
        const wrapper = ({ children }) => (React.createElement(Router, { history: history },
            React.createElement(AlertmanagerProvider, { accessType: "instance" }, children)));
        store.set(ALERTMANAGER_NAME_LOCAL_STORAGE_KEY, externalAmMimir.name);
        const { result } = renderHook(() => useAlertmanager(), { wrapper });
        expect(result.current.selectedAlertmanager).toBe(externalAmProm.name);
    });
});
test('isAlertManagerWithConfigAPI', () => {
    expect(isAlertManagerWithConfigAPI({
        implementation: undefined,
    })).toBe(true);
    expect(isAlertManagerWithConfigAPI({
        implementation: AlertManagerImplementation.mimir,
    })).toBe(true);
    expect(isAlertManagerWithConfigAPI({
        implementation: AlertManagerImplementation.prometheus,
    })).toBe(false);
});
//# sourceMappingURL=AlertmanagerContext.test.js.map