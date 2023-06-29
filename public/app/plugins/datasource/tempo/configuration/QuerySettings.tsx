import { css } from '@emotion/css';
import React, { useState } from 'react';

import { DataSourcePluginOptionsEditorProps, GrafanaTheme2, updateDatasourcePluginJsonDataOption } from '@grafana/data';
import { InlineField, InlineSwitch, useStyles2 } from '@grafana/ui';
import { TimeRangeShift } from 'app/core/components/TimeRangeShift/TimeRangeShift';
import { validateTimeShift } from 'app/core/components/TimeRangeShift/validation';

import { TempoJsonData } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<TempoJsonData> {}

export function QuerySettings({ options, onOptionsChange }: Props) {
  const styles = useStyles2(getStyles);
  const [spanStartTimeShiftIsInvalid, setSpanStartTimeShiftIsInvalid] = useState(() => {
    return options.jsonData.traceQuery?.spanStartTimeShift
      ? validateTimeShift(options.jsonData.traceQuery?.spanStartTimeShift)
      : false;
  });
  const [spanEndTimeShiftIsInvalid, setSpanEndTimeShiftIsInvalid] = useState(() => {
    return options.jsonData.traceQuery?.spanEndTimeShift
      ? validateTimeShift(options.jsonData.traceQuery?.spanEndTimeShift)
      : false;
  });

  const getLabel = (type: 'start' | 'end') => {
    return `Time shift for ${type} of search`;
  };

  const getTooltip = (type: 'start' | 'end') => {
    return `Shifts the ${type} of the time range when searching by TraceID. Searching can return traces that do not fully fall into the search time range, so we recommend using higher time shifts for longer traces. Default: 30m (Time units can be used here, for example: 5s, 1m, 3h`;
  };

  return (
    <div className={styles.container}>
      <InlineField
        label="Use time range in query"
        tooltip="The time range can be used when there are performance issues or timeouts since it will narrow down the search to the defined range. Default: disabled"
        labelWidth={26}
      >
        <InlineSwitch
          id="enable-time-shift"
          value={options.jsonData.traceQuery?.timeShiftEnabled || false}
          onChange={(event) => {
            updateDatasourcePluginJsonDataOption({ onOptionsChange, options }, 'traceQuery', {
              ...options.jsonData.traceQuery,
              timeShiftEnabled: event.currentTarget.checked,
            });
          }}
        />
      </InlineField>

      <TimeRangeShift
        label={getLabel('start')}
        tooltip={getTooltip('start')}
        value={options.jsonData.traceQuery?.spanStartTimeShift || ''}
        disabled={!options.jsonData.traceQuery?.timeShiftEnabled}
        onChange={(val) => {
          setSpanStartTimeShiftIsInvalid(validateTimeShift(val));
          updateDatasourcePluginJsonDataOption({ onOptionsChange, options }, 'traceQuery', {
            ...options.jsonData.traceQuery,
            spanStartTimeShift: val,
          });
        }}
        isInvalid={spanStartTimeShiftIsInvalid}
      />

      <TimeRangeShift
        label={getLabel('end')}
        tooltip={getTooltip('end')}
        value={options.jsonData.traceQuery?.spanEndTimeShift || ''}
        disabled={!options.jsonData.traceQuery?.timeShiftEnabled}
        onChange={(val) => {
          setSpanEndTimeShiftIsInvalid(validateTimeShift(val));
          updateDatasourcePluginJsonDataOption({ onOptionsChange, options }, 'traceQuery', {
            ...options.jsonData.traceQuery,
            spanEndTimeShift: val,
          });
        }}
        isInvalid={spanEndTimeShiftIsInvalid}
      />
    </div>
  );
}

export const getStyles = (theme: GrafanaTheme2) => ({
  infoText: css`
    padding-bottom: ${theme.spacing(2)};
    color: ${theme.colors.text.secondary};
  `,
  container: css`
    width: 100%;
  `,
  row: css`
    align-items: baseline;
  `,
});
