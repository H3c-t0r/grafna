import { VisualizationSuggestionBuilderUtil, VisualizationSuggestionsInput } from '@grafana/data';
import { GraphFieldConfig, LegendDisplayMode } from '@grafana/schema';
import { TimeSeriesOptions } from './types';

export function getSuggestions({ data }: VisualizationSuggestionsInput) {
  if (!data || !data.series || data.series.length === 0) {
    return [];
  }

  const frames = data.series;
  const builder = new VisualizationSuggestionBuilderUtil<TimeSeriesOptions, GraphFieldConfig>({
    name: 'Line graph',
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
    },
  });

  if (frames.length === 1) {
    builder.add({});
  }

  return builder.getList();
}
