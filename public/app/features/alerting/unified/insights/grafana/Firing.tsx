import { ThresholdsMode } from '@grafana/data';
import { PanelBuilders, SceneFlexItem, SceneQueryRunner, SceneTimeRange } from '@grafana/scenes';
import { DataSourceRef } from '@grafana/schema';

const QUERY = 'sum by (state) (grafanacloud_grafana_instance_alerting_rule_group_rules{state="active"})';

export function getFiringGrafanaAlertsScene(timeRange: SceneTimeRange, datasource: DataSourceRef, panelTitle: string) {
  const query = new SceneQueryRunner({
    datasource,
    queries: [
      {
        refId: 'A',
        instant: true,
        expr: QUERY,
      },
    ],
    $timeRange: timeRange,
  });

  return new SceneFlexItem({
    width: 'calc(25% - 4px)',
    height: 300,
    body: PanelBuilders.stat()
      .setTitle(panelTitle)
      .setData(query)
      .setThresholds({
        mode: ThresholdsMode.Absolute,
        steps: [
          {
            color: 'red',
            value: 0,
          },
          {
            color: 'red',
            value: 80,
          },
        ],
      })
      .build(),
  });
}
