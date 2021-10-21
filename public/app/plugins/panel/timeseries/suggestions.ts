import { VisualizationSuggestionsBuilder } from '@grafana/data';
import {
  GraphDrawStyle,
  GraphFieldConfig,
  GraphGradientMode,
  LegendDisplayMode,
  LineInterpolation,
  StackingMode,
} from '@grafana/schema';
import { SuggestionName } from 'app/types/suggestions';
import { TimeSeriesOptions } from './types';

export class TimeSeriesSuggestionsSupplier {
  getSuggestions(builder: VisualizationSuggestionsBuilder) {
    const list = builder.getListAppender<TimeSeriesOptions, GraphFieldConfig>({
      name: SuggestionName.LineChart,
      pluginId: 'timeseries',
      options: {
        legend: {} as any,
      },
      fieldConfig: {
        defaults: {
          custom: {},
        },
        overrides: [],
      },
      previewModifier: (s) => {
        s.options!.legend.displayMode = LegendDisplayMode.Hidden;

        if (s.fieldConfig?.defaults.custom?.drawStyle !== GraphDrawStyle.Bars) {
          s.fieldConfig!.defaults.custom!.lineWidth = 3;
        }
      },
    });

    const { dataSummary } = builder;

    if (!dataSummary.hasTimeField || !dataSummary.hasNumberField) {
      return;
    }

    list.append({
      name: SuggestionName.LineChart,
    });

    if (dataSummary.rowCountMax < 200) {
      list.append({
        name: SuggestionName.LineChartSmooth,
        fieldConfig: {
          defaults: {
            custom: {
              lineInterpolation: LineInterpolation.Smooth,
            },
          },
          overrides: [],
        },
      });
    }

    // Single series suggestions
    if (dataSummary.numberFieldCount === 1) {
      list.append({
        name: SuggestionName.AreaChart,
        fieldConfig: {
          defaults: {
            custom: {
              fillOpacity: 25,
            },
          },
          overrides: [],
        },
      });

      list.append({
        name: SuggestionName.BarChart,
        fieldConfig: {
          defaults: {
            custom: {
              drawStyle: GraphDrawStyle.Bars,
              fillOpacity: 100,
              lineWidth: 1,
              gradientMode: GraphGradientMode.Hue,
            },
          },
          overrides: [],
        },
      });

      return;
    }

    // Multiple series suggestions

    list.append({
      name: SuggestionName.AreaChartStacked,
      fieldConfig: {
        defaults: {
          custom: {
            fillOpacity: 25,
            stacking: {
              mode: StackingMode.Normal,
              group: 'A',
            },
          },
        },
        overrides: [],
      },
    });

    list.append({
      name: SuggestionName.AreaChartStackedPercent,
      fieldConfig: {
        defaults: {
          custom: {
            fillOpacity: 25,
            stacking: {
              mode: StackingMode.Percent,
              group: 'A',
            },
          },
        },
        overrides: [],
      },
    });

    if (dataSummary.rowCountTotal / dataSummary.numberFieldCount < 100) {
      list.append({
        name: SuggestionName.BarChartStacked,
        fieldConfig: {
          defaults: {
            custom: {
              drawStyle: GraphDrawStyle.Bars,
              fillOpacity: 100,
              lineWidth: 1,
              gradientMode: GraphGradientMode.Hue,
              stacking: {
                mode: StackingMode.Normal,
                group: 'A',
              },
            },
          },
          overrides: [],
        },
      });

      list.append({
        name: SuggestionName.BarChartStackedPercent,
        fieldConfig: {
          defaults: {
            custom: {
              drawStyle: GraphDrawStyle.Bars,
              fillOpacity: 100,
              lineWidth: 1,
              gradientMode: GraphGradientMode.Hue,
              stacking: {
                mode: StackingMode.Percent,
                group: 'A',
              },
            },
          },
          overrides: [],
        },
      });
    }
  }
}
