import { css } from '@emotion/css';
import React, { useEffect } from 'react';

import { GrafanaTheme2 } from '@grafana/data';
import { Alert, LoadingPlaceholder, useStyles2 } from '@grafana/ui';
import { useQueryParams } from 'app/core/hooks/useQueryParams';
import { useDispatch } from 'app/types';

import { AlertmanagerChoice } from '../../../plugins/datasource/alertmanager/types';

import { alertmanagerApi } from './api/alertmanagerApi';
import { AlertmanagerPageWrapper } from './components/AlertingPageWrapper';
import { AlertGroup } from './components/alert-groups/AlertGroup';
import { AlertGroupFilter } from './components/alert-groups/AlertGroupFilter';
import { useFilteredAmGroups } from './hooks/useFilteredAmGroups';
import { useGroupedAlerts } from './hooks/useGroupedAlerts';
import { useUnifiedAlertingSelector } from './hooks/useUnifiedAlertingSelector';
import { useSelectedAlertmanager } from './state/AlertmanagerContext';
import { fetchAlertGroupsAction } from './state/actions';
import { NOTIFICATIONS_POLL_INTERVAL_MS } from './utils/constants';
import { GRAFANA_RULES_SOURCE_NAME } from './utils/datasource';
import { getFiltersFromUrlParams } from './utils/misc';
import { initialAsyncRequestState } from './utils/redux';

const AlertGroups = () => {
  const { useGetAlertmanagerChoiceStatusQuery } = alertmanagerApi;

  const { selectedAlertmanager } = useSelectedAlertmanager();
  const dispatch = useDispatch();
  const [queryParams] = useQueryParams();
  const { groupBy = [] } = getFiltersFromUrlParams(queryParams);
  const styles = useStyles2(getStyles);

  const { currentData: amConfigStatus } = useGetAlertmanagerChoiceStatusQuery();

  const alertGroups = useUnifiedAlertingSelector((state) => state.amAlertGroups);
  const { loading, error, result: results = [] } = alertGroups[selectedAlertmanager || ''] ?? initialAsyncRequestState;

  const groupedAlerts = useGroupedAlerts(results, groupBy);
  const filteredAlertGroups = useFilteredAmGroups(groupedAlerts);

  const grafanaAmDeliveryDisabled =
    selectedAlertmanager === GRAFANA_RULES_SOURCE_NAME &&
    amConfigStatus?.alertmanagersChoice === AlertmanagerChoice.External;

  useEffect(() => {
    function fetchNotifications() {
      if (selectedAlertmanager) {
        dispatch(fetchAlertGroupsAction(selectedAlertmanager));
      }
    }
    fetchNotifications();
    const interval = setInterval(fetchNotifications, NOTIFICATIONS_POLL_INTERVAL_MS);
    return () => {
      clearInterval(interval);
    };
  }, [dispatch, selectedAlertmanager]);

  return (
    <>
      <AlertGroupFilter groups={results} />
      {loading && <LoadingPlaceholder text="Loading notifications" />}
      {error && !loading && (
        <Alert title={'Error loading notifications'} severity={'error'}>
          {error.message || 'Unknown error'}
        </Alert>
      )}

      {grafanaAmDeliveryDisabled && (
        <Alert title="Grafana alerts are not delivered to Grafana Alertmanager">
          Grafana is configured to send alerts to external alertmanagers only. No alerts are expected to be available
          here for the selected Alertmanager.
        </Alert>
      )}

      {results &&
        filteredAlertGroups.map((group, index) => {
          return (
            <React.Fragment key={`${JSON.stringify(group.labels)}-group-${index}`}>
              {((index === 1 && Object.keys(filteredAlertGroups[0].labels).length === 0) ||
                (index === 0 && Object.keys(group.labels).length > 0)) && (
                <p className={styles.groupingBanner}>Grouped by: {Object.keys(group.labels).join(', ')}</p>
              )}
              <AlertGroup alertManagerSourceName={selectedAlertmanager || ''} group={group} />
            </React.Fragment>
          );
        })}
      {results && !filteredAlertGroups.length && <p>No results.</p>}
    </>
  );
};

const AlertGroupsPage = () => (
  <AlertmanagerPageWrapper pageId="groups" accessType="instance">
    <AlertGroups />
  </AlertmanagerPageWrapper>
);

const getStyles = (theme: GrafanaTheme2) => ({
  groupingBanner: css`
    margin: ${theme.spacing(2, 0)};
  `,
});

export default AlertGroupsPage;
