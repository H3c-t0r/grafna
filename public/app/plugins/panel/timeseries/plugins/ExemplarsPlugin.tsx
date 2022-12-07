import React, { useCallback, useLayoutEffect, useRef } from 'react';
import uPlot from 'uplot';

import {
  DataFrame,
  DataFrameFieldIndex,
  Field,
  Labels,
  LinkModel,
  TIME_SERIES_TIME_FIELD_NAME,
  TIME_SERIES_VALUE_FIELD_NAME,
  TimeZone,
} from '@grafana/data';
import { EventsCanvas, FIXED_UNIT, UPlotConfigBuilder } from '@grafana/ui';

import { ExemplarMarker } from './ExemplarMarker';

interface ExemplarsPluginProps {
  config: UPlotConfigBuilder;
  exemplars: DataFrame[];
  timeZone: TimeZone;
  getFieldLinks: (field: Field, rowIndex: number) => Array<LinkModel<Field>>;
  visibleSeries: VisibleExemplarLabels;
}

export const ExemplarsPlugin: React.FC<ExemplarsPluginProps> = ({
  exemplars,
  timeZone,
  getFieldLinks,
  config,
  visibleSeries,
}) => {
  const plotInstance = useRef<uPlot>();

  useLayoutEffect(() => {
    config.addHook('init', (u) => {
      plotInstance.current = u;
    });
  }, [config]);

  const mapExemplarToXYCoords = useCallback((dataFrame: DataFrame, dataFrameFieldIndex: DataFrameFieldIndex) => {
    const time = dataFrame.fields.find((f) => f.name === TIME_SERIES_TIME_FIELD_NAME);
    const value = dataFrame.fields.find((f) => f.name === TIME_SERIES_VALUE_FIELD_NAME);

    if (!time || !value || !plotInstance.current) {
      return undefined;
    }

    // Filter x, y scales out
    const yScale =
      Object.keys(plotInstance.current.scales).find((scale) => !['x', 'y'].some((key) => key === scale)) ?? FIXED_UNIT;

    const yMin = plotInstance.current.scales[yScale].min;
    const yMax = plotInstance.current.scales[yScale].max;

    let y = value.values.get(dataFrameFieldIndex.fieldIndex);
    // To not to show exemplars outside of the graph we set the y value to min if it is smaller and max if it is bigger than the size of the graph
    if (yMin != null && y < yMin) {
      y = yMin;
    }

    if (yMax != null && y > yMax) {
      y = yMax;
    }

    return {
      x: plotInstance.current.valToPos(time.values.get(dataFrameFieldIndex.fieldIndex), 'x'),
      y: plotInstance.current.valToPos(y, yScale),
    };
  }, []);

  const renderMarker = useCallback(
    (dataFrame: DataFrame, dataFrameFieldIndex: DataFrameFieldIndex) => {
      let showMarker = false;

      const markerColor = getExemplarColor(dataFrame, dataFrameFieldIndex, visibleSeries);

      showMarker = showExemplar(visibleSeries, showMarker, dataFrame, dataFrameFieldIndex);

      if (!showMarker) {
        return <></>;
      }

      return (
        <ExemplarMarker
          timeZone={timeZone}
          getFieldLinks={getFieldLinks}
          dataFrame={dataFrame}
          dataFrameFieldIndex={dataFrameFieldIndex}
          config={config}
          exemplarColor={markerColor}
        />
      );
    },
    [config, timeZone, getFieldLinks, visibleSeries]
  );

  return (
    <EventsCanvas
      config={config}
      id="exemplars"
      events={exemplars}
      renderEventMarker={renderMarker}
      mapEventToXYCoords={mapExemplarToXYCoords}
    />
  );
};

export type VisibleExemplarLabels = { labels: LabelWithExemplarUIData[]; totalSeriesCount: number };
/**
 * Get labels that are currently visible/active in the legend
 */
export const getVisibleLabels = (config: UPlotConfigBuilder, frames: DataFrame[] | null): VisibleExemplarLabels => {
  const visibleSeries = config.series.filter((series) => series.props.show);
  const visibleLabels: LabelWithExemplarUIData[] = [];
  if (frames?.length) {
    visibleSeries.forEach((plotInstance) => {
      const frameIndex = plotInstance.props?.dataFrameFieldIndex?.frameIndex;
      const fieldIndex = plotInstance.props?.dataFrameFieldIndex?.fieldIndex;

      if (frameIndex !== undefined && fieldIndex !== undefined) {
        const field = frames[frameIndex].fields[fieldIndex];
        if (field.labels) {
          // Note that this may be an empty object in the case of a metric being rendered with no labels
          visibleLabels.push({
            labels: field.labels,
            color: plotInstance.props?.lineColor ?? '',
          });
        }
      }
    });
  }

  return { labels: visibleLabels, totalSeriesCount: config.series.length };
};

interface LabelWithExemplarUIData {
  labels: Labels;
  color?: string;
}
/**
 * Get color of active series in legend
 */
const getExemplarColor = (
  dataFrame: DataFrame,
  dataFrameFieldIndex: DataFrameFieldIndex,
  visibleLabels: VisibleExemplarLabels
) => {
  let exemplarColor;
  visibleLabels.labels.forEach((visibleLabel) => {
    const labelKeys = Object.keys(visibleLabel.labels);
    const fields = dataFrame.fields.filter((field) => {
      return labelKeys.find((labelKey) => labelKey === field.name);
    });
    if (fields.length) {
      const hasMatch = fields.every((field, index, fields) => {
        const value = field.values.get(dataFrameFieldIndex.fieldIndex);
        return visibleLabel.labels[field.name] === value;
      });

      if (hasMatch) {
        exemplarColor = visibleLabel.color;
        return;
      }
    }
  });
  return exemplarColor;
};

/**
 * Determine if the current exemplar marker is filtered by what series are selected in the legend UI
 */
function showExemplar(
  visibleSeries: VisibleExemplarLabels,
  showMarker: boolean,
  dataFrame: DataFrame,
  dataFrameFieldIndex: DataFrameFieldIndex
) {
  // If all series are visible, don't filter any exemplars
  if (visibleSeries.labels.length === visibleSeries.totalSeriesCount) {
    showMarker = true;
  } else {
    visibleSeries.labels.forEach((visibleLabel) => {
      // Get the label names
      const labelKeys = Object.keys(visibleLabel.labels);

      // If there aren't any labels, the graph is only displaying a single series with exemplars, let's show all exemplars in this case as well
      if (Object.keys(visibleLabel.labels).length === 0) {
        showMarker = true;
      } else {
        // If there are labels, lets only show the exemplars with labels associated with series that are currently visible
        const fields = dataFrame.fields.filter((field) => {
          return labelKeys.find((labelKey) => labelKey === field.name);
        });

        if (fields.length) {
          // Check to see if at least one value matches each field
          showMarker = visibleSeries.labels.some((series) => {
            return Object.keys(series.labels).every((label) => {
              const value = series.labels[label];
              return fields.find((field) => field.values.get(dataFrameFieldIndex.fieldIndex) === value);
            });
          });
        }
      }
    });
  }
  return showMarker;
}
