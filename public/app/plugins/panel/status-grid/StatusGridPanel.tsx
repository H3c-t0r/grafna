import React, { useMemo } from 'react';
import { PanelProps } from '@grafana/data';
import { useTheme2, ZoomPlugin } from '@grafana/ui';
import { StatusPanelOptions } from './types';
import { TimelineChart } from '../state-timeline/TimelineChart';
import { TimelineMode } from '../state-timeline/types';
import { prepareTimelineFields, prepareTimelineLegendItems } from '../state-timeline/utils';

interface TimelinePanelProps extends PanelProps<StatusPanelOptions> {}

/**
 * @alpha
 */
export const StatusGridPanel: React.FC<TimelinePanelProps> = ({
  data,
  timeRange,
  timeZone,
  options,
  width,
  height,
  onChangeTimeRange,
}) => {
  const theme = useTheme2();

  const { frames, warn } = useMemo(() => prepareTimelineFields(data?.series, false), [data]);

  const legendItems = useMemo(() => prepareTimelineLegendItems(frames, options.legend), [frames, options.legend]);

  if (!frames || warn) {
    return (
      <div className="panel-empty">
        <p>{warn ?? 'No data found in response'}</p>
      </div>
    );
  }

  return (
    <TimelineChart
      theme={theme}
      frames={frames}
      structureRev={data.structureRev}
      timeRange={timeRange}
      timeZone={timeZone}
      width={width}
      height={height}
      legendItems={legendItems}
      {...options}
      // hardcoded
      mode={TimelineMode.Samples}
    >
      {(config) => <ZoomPlugin config={config} onZoom={onChangeTimeRange} />}
    </TimelineChart>
  );
};
