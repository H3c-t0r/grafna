import { css } from '@emotion/css';
import React, { useState } from 'react';

import { dateTime, getDefaultRelativeTimeRange, GrafanaTheme2, RelativeTimeRange } from '@grafana/data';
import { relativeToTimeRange } from '@grafana/data/src/datetime/rangeutil';
import { clearButtonStyles, Icon, RelativeTimeRangePicker, Toggletip, useStyles2 } from '@grafana/ui';
import { AlertQuery } from 'app/types/unified-alerting-dto';

import { AlertQueryOptions, DEFAULT_MAX_DATA_POINTS, MaxDataPointsOption } from './QueryWrapper';

export interface QueryOptionsProps {
  query: AlertQuery;
  queryOptions: AlertQueryOptions;
  onChangeTimeRange?: (timeRange: RelativeTimeRange, index: number) => void;
  onChangeQueryOptions: (options: AlertQueryOptions, index: number) => void;
  index: number;
}

export const QueryOptions = ({
  query,
  queryOptions,
  onChangeTimeRange,
  onChangeQueryOptions,
  index,
}: QueryOptionsProps) => {
  const styles = useStyles2(getStyles);

  const [showOptions, setShowOptions] = useState(false);

  const timeRange = query.relativeTimeRange ? relativeToTimeRange(query.relativeTimeRange) : undefined;

  return (
    <>
      <Toggletip
        content={
          <>
            {onChangeTimeRange && (
              <RelativeTimeRangePicker
                timeRange={query.relativeTimeRange ?? getDefaultRelativeTimeRange()}
                onChange={(range) => onChangeTimeRange(range, index)}
              />
            )}
            <div className={styles.queryOptions}>
              <MaxDataPointsOption
                options={queryOptions}
                onChange={(options) => onChangeQueryOptions(options, index)}
              />
            </div>
          </>
        }
        closeButton={true}
        placement="bottom-start"
      >
        <button type="button" className={styles.actionLink} onClick={() => setShowOptions(!showOptions)}>
          Options {showOptions ? <Icon name="angle-down" /> : <Icon name="angle-right" />}
        </button>
      </Toggletip>

      <div className={styles.staticValues}>
        <span>{dateTime(timeRange?.from).locale('en').fromNow(true)}, </span>
        <span>MD {queryOptions?.maxDataPoints || DEFAULT_MAX_DATA_POINTS.toLocaleString()}</span>
      </div>
    </>
  );
};

const getStyles = (theme: GrafanaTheme2) => {
  const clearButton = clearButtonStyles(theme);

  return {
    queryOptions: css`
      margin-bottom: -${theme.spacing(2)};
    `,

    staticValues: css`
      color: ${theme.colors.text.secondary};
      margin-right: ${theme.spacing(1)};
    `,

    actionLink: css`
      ${clearButton};
      color: ${theme.colors.text.link};
      cursor: pointer;

      &:hover {
        text-decoration: underline;
      }
    `,
  };
};
