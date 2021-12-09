import React from 'react';
import { css } from '@emotion/css';
import { Collapse } from '@grafana/ui';
import { AbsoluteTimeRange, LoadingState, TimeZone, DataFrame, DataSourceApi, GrafanaTheme2 } from '@grafana/data';
import { ExploreGraph } from './ExploreGraph';
import { EXPLORE_GRAPH_STYLES } from 'app/core/utils/explore';

export interface Props {
  timeZone: TimeZone;
  graphStyle: any;
  datasourceInstance?: DataSourceApi | null;
  autoBreakdownRange?: AbsoluteTimeRange;
}

const spacing = css({
  marginRight: '10px',
  marginBottom: '10px',
});

function calcWidth(df: DataFrame): number {
  const baseWidth = 300;
  const margin = 10;
  const panelPadding = 8;
  const panelBorder = 1;
  return df.length <= 8 ? baseWidth : 2 * baseWidth + margin + 2 * panelPadding + 2 * panelBorder;
}

export const AutoBreakdowns = (props: Props) => {
  const { datasourceInstance, graphStyle, autoBreakdownRange, timeZone } = props;
  if (!datasourceInstance || !autoBreakdownRange) {
    return null;
  }
  //@ts-ignore
  const dataFrames = datasourceInstance?.createAutoBreakdowns(autoBreakdownRange);

  return (
    <div style={{ display: 'flex', flexWrap: 'wrap' }}>
      {dataFrames.map((df: DataFrame) => {
        return (
          <div className={spacing} key={`${df.name}`}>
            <Collapse label={`Breakdown: ${df.name}`} isOpen>
              <ExploreGraph
                graphStyle="bars"
                data={[df]}
                height={150}
                width={calcWidth(df)}
                timeZone={timeZone}
                absoluteRange={autoBreakdownRange}
                loadingState={LoadingState.Done}
                onChangeTime={() => {}}
                pluginId="barchart"
              />
            </Collapse>
          </div>
        );
      })}
    </div>
  );
};
