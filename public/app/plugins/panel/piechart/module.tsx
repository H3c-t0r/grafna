import { FieldColorModeId, FieldConfigProperty, PanelPlugin } from '@grafana/data';
import { PieChartPanel } from './PieChartPanel';
import { PieChartOptions } from './types';
import { addStandardDataReduceOptions } from '../stat/types';
import { LegendDisplayMode, PieChartType } from '@grafana/ui';
import { PieChartLabels } from '@grafana/ui/src/components/PieChart/PieChart';

export const plugin = new PanelPlugin<PieChartOptions>(PieChartPanel)
  .useFieldConfig({
    standardOptions: {
      [FieldConfigProperty.Color]: {
        settings: {
          byValueSupport: false,
          bySeriesSupport: true,
          preferThresholdsMode: false,
        },
        defaultValue: {
          mode: FieldColorModeId.PaletteClassic,
        },
      },
    },
  })
  .setPanelOptions((builder) => {
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
      .addSelect({
        name: 'Labels',
        path: 'displayLabels',
        defaultValue: [],
        description: 'Select the labels to be displayed in the pie chart',
        settings: {
          options: [
            { value: PieChartLabels.Percent, label: 'Percent' },
            { value: PieChartLabels.Name, label: 'Name' },
            { value: PieChartLabels.Value, label: 'Value' },
          ],
        },
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
      })
      .addBooleanSwitch({
        name: 'Show percent in legend',
        path: 'legend.showPercent',
        defaultValue: false,
      })
      .addBooleanSwitch({
        name: 'Show value in legend',
        path: 'legend.showValue',
        defaultValue: false,
      })
      .addRadio({
        path: 'legend.displayMode',
        name: 'Legend mode',
        description: '',
        defaultValue: LegendDisplayMode.List,
        settings: {
          options: [
            { value: LegendDisplayMode.List, label: 'List' },
            { value: LegendDisplayMode.Table, label: 'Table' },
            { value: LegendDisplayMode.Hidden, label: 'Hidden' },
          ],
        },
      })
      .addRadio({
        path: 'legend.placement',
        name: 'Legend placement',
        description: '',
        defaultValue: 'right',
        settings: {
          options: [
            { value: 'bottom', label: 'Bottom' },
            { value: 'right', label: 'Right' },
          ],
        },
        showIf: (c) => c.legend.displayMode !== LegendDisplayMode.Hidden,
      });
  });
