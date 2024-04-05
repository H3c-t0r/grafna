import { config } from '@grafana/runtime';
import { SceneDataProvider, SceneDataTransformer, SceneQueryRunner } from '@grafana/scenes';
import { PanelModel } from 'app/features/dashboard/state';
import { SHARED_DASHBOARD_QUERY } from 'app/plugins/datasource/dashboard';

import { DashboardDatasourceBehaviour } from '../scene/DashboardDatasourceBehaviour';

export function createPanelDataProvider(panel: PanelModel): SceneDataProvider | undefined {
  // Skip setting query runner for panels without queries
  if (!panel.targets?.length) {
    return undefined;
  }

  // Skip setting query runner for panel plugins with skipDataQuery
  if (config.panels[panel.type]?.skipDataQuery) {
    return undefined;
  }

  let dataProvider: SceneDataProvider | undefined = undefined;

  dataProvider = new SceneQueryRunner({
    datasource: panel.datasource ?? undefined,
    queries: panel.targets,
    maxDataPoints: panel.maxDataPoints ?? undefined,
    maxDataPointsFromWidth: true,
    cacheTimeout: panel.cacheTimeout,
    queryCachingTTL: panel.queryCachingTTL,
    dataLayerFilter: {
      panelId: panel.id,
    },
    $behaviors: panel.datasource?.uid === SHARED_DASHBOARD_QUERY ? [new DashboardDatasourceBehaviour({})] : [],
  });

  // Wrap inner data provider in a data transformer
  return new SceneDataTransformer({
    $data: dataProvider,
    transformations: panel.transformations || [],
  });
}
