import { css } from '@emotion/css';
import React from 'react';

import {
  AbsoluteTimeRange,
  DataQueryResponse,
  LoadingState,
  SplitOpen,
  TimeZone,
  EventBus,
  isLogsVolumeLimited,
  getLogsVolumeAbsoluteRange,
  GrafanaTheme2,
} from '@grafana/data';
import { TooltipDisplayMode, useStyles2, useTheme2 } from '@grafana/ui';

import { ExploreGraph } from './Graph/ExploreGraph';

type Props = {
  extraInfo: string;
  logsVolumeData: DataQueryResponse | undefined;
  absoluteRange: AbsoluteTimeRange;
  timeZone: TimeZone;
  splitOpen: SplitOpen;
  width: number;
  onUpdateTimeRange: (timeRange: AbsoluteTimeRange) => void;
  onLoadLogsVolume: () => void;
  onHiddenSeriesChanged: (hiddenSeries: string[]) => void;
  eventBus: EventBus;
};

export function LogsVolumePanel(props: Props) {
  const { width, timeZone, splitOpen, onUpdateTimeRange, onHiddenSeriesChanged, extraInfo } = props;
  const theme = useTheme2();
  const styles = useStyles2(getStyles);
  const spacing = parseInt(theme.spacing(2).slice(0, -2), 10);
  const height = 150;

  if (props.logsVolumeData === undefined) {
    return null;
  }

  const logsVolumeData = props.logsVolumeData;

  const range = isLogsVolumeLimited(logsVolumeData.data)
    ? getLogsVolumeAbsoluteRange(logsVolumeData.data, props.absoluteRange)
    : props.absoluteRange;

  let LogsVolumePanelContent;

  if (logsVolumeData?.data) {
    if (logsVolumeData.data.length > 0) {
      LogsVolumePanelContent = (
        <ExploreGraph
          graphStyle="lines"
          loadingState={LoadingState.Done}
          data={logsVolumeData.data}
          height={height}
          width={width - spacing * 2}
          absoluteRange={range}
          onChangeTime={onUpdateTimeRange}
          timeZone={timeZone}
          splitOpenFn={splitOpen}
          tooltipDisplayMode={TooltipDisplayMode.Multi}
          onHiddenSeriesChanged={onHiddenSeriesChanged}
          anchorToZero
          eventBus={props.eventBus}
        />
      );
    } else {
      LogsVolumePanelContent = <span>No volume data.</span>;
    }
  }

  return (
    <div style={{ height }} className={styles.contentContainer}>
      {LogsVolumePanelContent}
      {extraInfo && <div className={styles.extraInfoContainer}>{extraInfo}</div>}
    </div>
  );
}

const getStyles = (theme: GrafanaTheme2) => {
  return {
    extraInfoContainer: css`
      display: flex;
      justify-content: end;
      position: absolute;
      right: 5px;
      top: -10px;
      font-size: ${theme.typography.size.sm};
      color: ${theme.colors.text.secondary};
    `,
    contentContainer: css`
      display: flex;
      align-items: center;
      justify-content: center;
      position: relative;
    `,
  };
};
