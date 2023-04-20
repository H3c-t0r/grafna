import { css } from '@emotion/css';
import React from 'react';

import { dateTime, GrafanaTheme2 } from '@grafana/data';
import { Alert, Badge, LoadingPlaceholder, useStyles2 } from '@grafana/ui';
import { AlertmanagerAlert, Matcher } from 'app/plugins/datasource/alertmanager/types';

import { alertmanagerApi } from '../../api/alertmanagerApi';
import { AlertLabels } from '../AlertLabels';
import { DynamicTable, DynamicTableColumnProps, DynamicTableItemProps } from '../DynamicTable';

import { AmAlertStateTag } from './AmAlertStateTag';

interface Props {
  amSourceName: string;
  matchers: Matcher[];
}

export const MatchedSilencedRules = ({ amSourceName, matchers }: Props) => {
  const { useGetAlertmanagerAlertsQuery } = alertmanagerApi;
  const styles = useStyles2(getStyles);
  const columns = useColumns();

  const {
    currentData: alerts = [],
    isFetching,
    isError,
  } = useGetAlertmanagerAlertsQuery(
    { amSourceName, filter: { matchers } },
    { skip: matchers.length === 0, refetchOnMountOrArgChange: true }
  );

  const tableItemAlerts = alerts.map<DynamicTableItemProps<AlertmanagerAlert>>((alert) => ({
    id: alert.fingerprint,
    data: alert,
  }));

  const hasValidMatchers = matchers.some((matcher) => matcher.value && matcher.name);

  return (
    <div>
      <h4 className={styles.title}>
        Affected alert instances
        {tableItemAlerts.length > 0 ? (
          <Badge className={styles.badge} color="blue" text={tableItemAlerts.length} />
        ) : null}
      </h4>
      {!hasValidMatchers && <span>Add a valid matcher to see affected alerts</span>}
      {isError && (
        <Alert title="Preview not available" severity="error">
          Error occured when generating affected alerts preview. Are you matchers valid?
        </Alert>
      )}
      {isFetching && <LoadingPlaceholder text="Loading..." />}
      {!isFetching && !isError && (
        <div className={styles.table}>
          {tableItemAlerts.length > 0 ? (
            <DynamicTable
              items={tableItemAlerts}
              isExpandable={false}
              cols={columns}
              pagination={{ itemsPerPage: 10 }}
            />
          ) : (
            <span>No matching alert instances found</span>
          )}
        </div>
      )}
    </div>
  );
};

function useColumns(): Array<DynamicTableColumnProps<AlertmanagerAlert>> {
  const styles = useStyles2(getStyles);

  return [
    {
      id: 'state',
      label: 'State',
      renderCell: function renderStateTag({ data }) {
        return <AmAlertStateTag state={data.status.state} />;
      },
      size: '120px',
      className: styles.stateColumn,
    },
    {
      id: 'labels',
      label: 'Labels',
      renderCell: function renderName({ data }) {
        return <AlertLabels labels={data.labels} className={styles.alertLabels} />;
      },
      size: 'auto',
    },
    {
      id: 'created',
      label: 'Created',
      renderCell: function renderSummary({ data }) {
        return <>{data.startsAt.startsWith('0001') ? '-' : dateTime(data.startsAt).format('YYYY-MM-DD HH:mm:ss')}</>;
      },
      size: '180px',
    },
  ];
}

const getStyles = (theme: GrafanaTheme2) => ({
  table: css`
    max-width: ${theme.breakpoints.values.lg}px;
  `,
  moreMatches: css`
    margin-top: ${theme.spacing(1)};
  `,
  title: css`
    display: flex;
    align-items: center;
  `,
  badge: css`
    margin-left: ${theme.spacing(1)};
  `,
  stateColumn: css`
    display: flex;
    align-items: center;
  `,
  alertLabels: css`
    justify-content: flex-start;
  `,
});
