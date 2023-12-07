import { css, cx } from '@emotion/css';
import React from 'react';

import { GrafanaTheme2 } from '@grafana/data';
import { IconButton, useStyles2 } from '@grafana/ui';

const getStyles = (theme: GrafanaTheme2) => ({
  logsStatsRow: css`
    label: logs-stats-row;
    margin: ${parseInt(theme.spacing(2), 10) / 1.75}px 0;
  `,
  logsStatsRowActive: css`
    label: logs-stats-row--active;
    color: ${theme.colors.primary.text};
    position: relative;
  `,
  logsStatsRowLabel: css`
    label: logs-stats-row__label;
    display: flex;
    margin-bottom: 1px;
  `,
  logsStatsRowValue: css`
    label: logs-stats-row__value;
    flex: 1;
    text-overflow: ellipsis;
    overflow: hidden;
  `,
  logsStatsRowCount: css`
    label: logs-stats-row__count;
    text-align: right;
    margin-left: ${theme.spacing(0.75)};
  `,
  logsStatsRowPercent: css`
    label: logs-stats-row__percent;
    text-align: right;
    margin-left: ${theme.spacing(0.75)};
    width: ${theme.spacing(4.5)};
  `,
  logsStatsRowBar: css`
    label: logs-stats-row__bar;
    height: ${theme.spacing(0.5)};
    overflow: hidden;
    background: ${theme.colors.text.disabled};
  `,
  logsStatsRowInnerBar: css`
    label: logs-stats-row__innerbar;
    height: ${theme.spacing(0.5)};
    overflow: hidden;
    background: ${theme.colors.primary.main};
  `,
});

export interface Props {
  active?: boolean;
  count: number;
  proportion: number;
  value?: string;
  total: number;
  shouldFilter: boolean;
}

export const LogLabelStatsRow = ({ active, count, proportion, value, total, shouldFilter }: Props) => {
  const style = useStyles2(getStyles);
  const percent = `${Math.round(proportion * 100)}%`;
  // const barStyle = { width: percent };
  // const className = active ? cx([style.logsStatsRow, style.logsStatsRowActive]) : cx([style.logsStatsRow]);

  return (
    <div style={{ display: 'flex' }}>
      {shouldFilter && (
        <div style={{ width: '40px', marginTop: '3px' }}>
          <IconButton size="xs" name="search-plus" aria-label="search-plus" />
          <IconButton size="xs" name="search-minus" aria-label="search-minus" />
        </div>
      )}
      <div
        style={{
          width: '180px',
          textOverflow: 'ellipsis',
          overflow: 'hidden',
          whiteSpace: 'nowrap',
          marginRight: '20px',
        }}
      >
        {value}
      </div>
      <div style={{ width: '80px' }}>
        {count}/{total}
      </div>
      <div style={{ width: '44px' }}>{percent}</div>
    </div>
  );

  // return (
  //   <div className={className}>
  //     <div className={cx([style.logsStatsRowLabel])}>
  //       <div className={cx([style.logsStatsRowValue])} title={value}>
  //         {value}
  //       </div>
  //       <div className={cx([style.logsStatsRowCount])}>{count}</div>
  //       <div className={cx([style.logsStatsRowPercent])}>{percent}</div>
  //     </div>
  //     <div className={cx([style.logsStatsRowBar])}>
  //       <div className={cx([style.logsStatsRowInnerBar])} style={barStyle} />
  //     </div>
  //   </div>
  // );
};

LogLabelStatsRow.displayName = 'LogLabelStatsRow';
