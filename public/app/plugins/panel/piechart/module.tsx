import { PanelPlugin } from '@grafana/data';
import { PieChartPanel } from './PieChartPanel';
import { PieChartOptions } from './types';
import { addStandardDataReduceOptions } from '../stat/types';
import { PieChartType } from '@grafana/ui';

export const plugin = new PanelPlugin<PieChartOptions>(PieChartPanel).useFieldConfig().setPanelOptions((builder) => {
  addStandardDataReduceOptions(builder, false);

  builder
    .addRadio({
      name: 'Piechart type',
      description: 'How the piechart should be rendered',
      path: 'pieType',
      settings: {
        options: [
          { value: PieChartType.Pie, label: 'Pie' },
          { value: PieChartType.Donut, label: 'Donut' },
        ],
      },
      defaultValue: PieChartType.Pie,
    })
    .addBooleanSwitch({
      name: 'Show name',
      path: 'labelOptions.showName',
      defaultValue: true,
    })
    .addBooleanSwitch({
      name: 'Show value',
      path: 'labelOptions.showValue',
      defaultValue: false,
    })
    .addBooleanSwitch({
      name: 'Show percent',
      path: 'labelOptions.showPercent',
      defaultValue: false,
    });
});
