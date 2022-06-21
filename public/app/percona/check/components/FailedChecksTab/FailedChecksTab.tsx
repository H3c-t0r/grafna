import { LoaderButton, logger } from '@percona/platform-core';
import React, { FC, useEffect, useState, useCallback } from 'react';

import { AppEvents } from '@grafana/data';
import { Spinner, Switch, useStyles } from '@grafana/ui';
import { AlertsReloadContext } from 'app/percona/check/Check.context';
import { CheckService } from 'app/percona/check/Check.service';
import { COLUMNS } from 'app/percona/check/CheckPanel.constants';
import { Table } from 'app/percona/check/components';
import { ActiveCheck } from 'app/percona/check/types';
import { useCancelToken } from 'app/percona/shared/components/hooks/cancelToken.hook';
import { isApiCancelError } from 'app/percona/shared/helpers/api';

import { appEvents } from '../../../../core/app_events';

import { GET_ACTIVE_ALERTS_CANCEL_TOKEN } from './FailedChecksTab.constants';
import { Messages } from './FailedChecksTab.messages';
import { getStyles } from './FailedChecksTab.styles';
import { loadShowSilencedValue, saveShowSilencedValue } from './FailedChecksTab.utils';
import { FailedChecksTabProps } from './types';

export const FailedChecksTab: FC<FailedChecksTabProps> = ({ hasNoAccess }) => {
  const [fetchAlertsPending, setFetchAlertsPending] = useState(true);
  const [runChecksPending, setRunChecksPending] = useState(false);
  const [showSilenced, setShowSilenced] = useState(loadShowSilencedValue());
  const [dataSource, setDataSource] = useState<ActiveCheck[] | undefined>();
  const styles = useStyles(getStyles);
  const [generateToken] = useCancelToken();

  const fetchAlerts = useCallback(async (): Promise<void> => {
    setFetchAlertsPending(true);

    try {
      const dataSource = await CheckService.getActiveAlerts(
        showSilenced,
        generateToken(GET_ACTIVE_ALERTS_CANCEL_TOKEN)
      );
      setDataSource(dataSource);
    } catch (e) {
      if (isApiCancelError(e)) {
        return;
      }
      logger.error(e);
    }
    setFetchAlertsPending(false);
  }, [generateToken, showSilenced]);

  const handleRunChecksClick = async () => {
    setRunChecksPending(true);
    try {
      await CheckService.runDbChecks();
    } catch (e) {
      logger.error(e);
    }
    // TODO (nicolalamacchia): remove this timeout when the API will become synchronous
    setTimeout(async () => {
      setRunChecksPending(false);
      await fetchAlerts();
      appEvents.emit(AppEvents.alertSuccess, ['Done running DB checks. The latest results are displayed.']);
    }, 10000);
  };

  const toggleShowSilenced = () => {
    setShowSilenced((currentValue) => !currentValue);
  };

  useEffect(() => {
    fetchAlerts();
    saveShowSilencedValue(showSilenced);
  }, [showSilenced, fetchAlerts]);

  return (
    <>
      <div className={styles.header}>
        <div className={styles.actionButtons} data-qa="db-check-panel-actions">
          <span className={styles.showAll}>
            <span data-qa="db-checks-failed-checks-toggle-silenced">
              <Switch value={showSilenced} onChange={toggleShowSilenced} />
            </span>
            <span>{Messages.showAll}</span>
          </span>
          <LoaderButton
            type="button"
            size="md"
            loading={runChecksPending}
            disabled={hasNoAccess}
            onClick={handleRunChecksClick}
            className={styles.runChecksButton}
          >
            {Messages.runDbChecks}
          </LoaderButton>
        </div>
      </div>
      <AlertsReloadContext.Provider value={{ fetchAlerts }}>
        {fetchAlertsPending ? (
          <div className={styles.spinner} data-qa="db-checks-failed-checks-spinner">
            <Spinner />
          </div>
        ) : (
          <Table data={dataSource} columns={COLUMNS} hasNoAccess={hasNoAccess} />
        )}
      </AlertsReloadContext.Provider>
    </>
  );
};
