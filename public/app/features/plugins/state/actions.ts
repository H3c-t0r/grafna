import { getBackendSrv } from '@grafana/runtime';
import { PanelPlugin } from '@grafana/data';
import { ThunkResult } from 'app/types';
import { pluginDashboardsLoad, pluginDashboardsLoaded, pluginsLoaded, panelPluginLoaded } from './reducers';
import { importPanelPlugin } from 'app/features/plugins/plugin_loader';

export function loadPlugins(): ThunkResult<void> {
  return async dispatch => {
    const result = await getBackendSrv().get('api/plugins', { embedded: 0 });
    dispatch(pluginsLoaded(result));
  };
}

export function loadPluginDashboards(): ThunkResult<void> {
  return async (dispatch, getStore) => {
    dispatch(pluginDashboardsLoad());
    const dataSourceType = getStore().dataSources.dataSource.type;
    const response = await getBackendSrv().get(`api/plugins/${dataSourceType}/dashboards`);
    dispatch(pluginDashboardsLoaded(response));
  };
}

export function loadPanelPlugin(pluginId: string): ThunkResult<Promise<PanelPlugin>> {
  return async (dispatch, getStore) => {
    let plugin = getStore().plugins.panels[pluginId];

    if (!plugin) {
      plugin = await importPanelPlugin(pluginId);

      dispatch(panelPluginLoaded(plugin));
    }

    return plugin;
  };
}
