import { useEffect } from 'react';
import { useAsyncFn } from 'react-use';

import { locationUtil } from '@grafana/data';
import { locationService } from '@grafana/runtime';
import { useAppNotification } from 'app/core/copy/appNotification';
import { getHistorySrv } from 'app/features/dashboard-scene/settings/version-history/HistorySrv';
import { useSelector } from 'app/types';

import { dashboardWatcher } from '../../../live/dashboard/dashboardWatcher';
import { DashboardModel } from '../../state';

const restoreDashboard = async (version: number, dashboard: DashboardModel) => {
  // Skip the watcher logic for this save since it's handled by the hook
  dashboardWatcher.ignoreNextSave();
  return await getHistorySrv().restoreDashboard(dashboard.uid, version);
};

export const useDashboardRestore = (version: number) => {
  const dashboard = useSelector((state) => state.dashboard.getModel());
  const [state, onRestoreDashboard] = useAsyncFn(async () => await restoreDashboard(version, dashboard!), []);
  const notifyApp = useAppNotification();

  useEffect(() => {
    if (state.value) {
      const location = locationService.getLocation();
      const newUrl = locationUtil.stripBaseFromUrl(state.value.url);
      const prevState = (location.state as any)?.routeReloadCounter;
      locationService.replace({
        ...location,
        pathname: newUrl,
        state: { routeReloadCounter: prevState ? prevState + 1 : 1 },
      });
      notifyApp.success('Dashboard restored', `Restored from version ${version}`);
    }
  }, [state, version, notifyApp]);
  return { state, onRestoreDashboard };
};
