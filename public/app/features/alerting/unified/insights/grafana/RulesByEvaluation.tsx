import { PanelBuilders, SceneFlexItem, SceneQueryRunner, SceneTimeRange } from '@grafana/scenes';
import { DataSourceRef, GraphDrawStyle, TooltipDisplayMode } from '@grafana/schema';

import { getPanelMenu, overrideToFixedColor, PANEL_STYLES } from '../../home/Insights';

export function getGrafanaRulesByEvaluationScene(
  timeRange: SceneTimeRange,
  datasource: DataSourceRef,
  panelTitle: string
) {
  const query = new SceneQueryRunner({
    datasource,
    queries: [
      {
        refId: 'A',
        expr: 'sum by (state) (grafanacloud_grafana_instance_alerting_rule_group_rules)',
        range: true,
        legendFormat: '{{state}} evaluation',
      },
    ],
    $timeRange: timeRange,
  });

  return new SceneFlexItem({
    ...PANEL_STYLES,
    body: PanelBuilders.timeseries()
      .setTitle(panelTitle)
      .setDescription(panelTitle)
      .setData(query)
      .setCustomFieldConfig('drawStyle', GraphDrawStyle.Line)
      .setOption('tooltip', { mode: TooltipDisplayMode.Multi })
      .setOverrides((b) =>
        b.matchFieldsWithName('active evaluation').overrideColor(overrideToFixedColor('active evaluation'))
      )
      .setMenu(getPanelMenu(panelTitle))
      .build(),
  });
}
