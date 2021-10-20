import { ThresholdsMode, VisualizationSuggestionsBuilder } from '@grafana/data';
import { GaugeOptions } from './types';

export class GaugeSuggestionsSupplier {
  getSuggestions(builder: VisualizationSuggestionsBuilder) {
    const list = builder.getListAppender<GaugeOptions, {}>({
      name: 'Gauge',
      pluginId: 'gauge',
      options: {},
      fieldConfig: {
        defaults: {
          thresholds: {
            steps: [
              { value: -Infinity, color: 'green' },
              { value: 70, color: 'orange' },
              { value: 85, color: 'red' },
            ],
            mode: ThresholdsMode.Percentage,
          },
          custom: {},
        },
        overrides: [],
      },
      previewModifier: (s) => {
        if (s.options!.reduceOptions.values) {
          s.options!.reduceOptions.limit = 2;
        }
      },
    });

    const { dataSummary } = builder;

    if (!dataSummary.hasNumberField) {
      return;
    }

    if (dataSummary.frameCount === 1 && dataSummary.rowCountTotal < 10) {
      list.append({
        options: {
          reduceOptions: {
            values: true,
            calcs: [],
          },
        },
      });
      list.append({
        name: 'Gauge without thresholds markers',
        options: {
          reduceOptions: {
            values: true,
            calcs: [],
          },
          showThresholdMarkers: false,
        },
      });
    } else {
      list.append({
        options: {
          reduceOptions: {
            values: false,
            calcs: ['lastNotNull'],
          },
        },
      });
      list.append({
        name: 'Gauge without thresholds markers',
        options: {
          reduceOptions: {
            values: false,
            calcs: ['lastNotNull'],
          },
          showThresholdMarkers: false,
        },
      });
    }
  }
}
