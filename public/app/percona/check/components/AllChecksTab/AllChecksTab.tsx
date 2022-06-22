import { cx } from '@emotion/css';
import { logger } from '@percona/platform-core';
import React, { FC, useEffect, useState } from 'react';

import { Spinner, useTheme, useStyles } from '@grafana/ui';
import { CheckService } from 'app/percona/check/Check.service';
// TODO: make a shared table style
import { getStyles as getCheckPanelStyles } from 'app/percona/check/CheckPanel.styles';
import { getStyles as getTableStyles } from 'app/percona/check/components/Table/Table.styles';
import { CheckDetails } from 'app/percona/check/types';
import { useCancelToken } from 'app/percona/shared/components/hooks/cancelToken.hook';
import { isApiCancelError } from 'app/percona/shared/helpers/api';

import { ChecksReloadContext } from './AllChecks.context';
import { GET_ALL_CHECKS_CANCEL_TOKEN } from './AllChecksTab.constants';
import { Messages } from './AllChecksTab.messages';
import * as styles from './AllChecksTab.styles';
import { CheckTableRow } from './CheckTableRow';
import { FetchChecks } from './types';

export const AllChecksTab: FC = () => {
  const [fetchChecksPending, setFetchChecksPending] = useState(false);
  const [checks, setChecks] = useState<CheckDetails[] | undefined>();
  const theme = useTheme();
  const tableStyles = getTableStyles(theme);
  const checkPanelStyles = useStyles(getCheckPanelStyles);
  const [generateToken] = useCancelToken();

  const updateUI = (check: CheckDetails) => {
    const { name, disabled } = check;

    setChecks((oldChecks) =>
      oldChecks?.map((oldCheck) => {
        if (oldCheck.name !== name) {
          return oldCheck;
        }

        return { ...oldCheck, disabled };
      })
    );
  };

  useEffect(() => {
    const fetchChecks: FetchChecks = async () => {
      setFetchChecksPending(true);

      try {
        const checks = await CheckService.getAllChecks(generateToken(GET_ALL_CHECKS_CANCEL_TOKEN));

        setChecks(checks);
      } catch (e) {
        if (isApiCancelError(e)) {
          return;
        }
        logger.error(e);
      }
      setFetchChecksPending(false);
    };

    fetchChecks();
  }, [generateToken]);

  return (
    <div className={cx(tableStyles.wrapper, styles.wrapper)} data-testid="db-checks-all-checks-wrapper">
      {fetchChecksPending ? (
        <div className={checkPanelStyles.spinner} data-testid="db-checks-all-checks-spinner">
          <Spinner />
        </div>
      ) : (
        <table className={tableStyles.table} data-testid="db-checks-all-checks-table">
          <colgroup>
            <col className={styles.nameColumn} />
            <col />
            <col className={styles.statusColumn} />
            <col className={styles.intervalColumn} />
            <col className={styles.actionsColumn} />
          </colgroup>
          <thead data-testid="db-checks-all-checks-thead">
            <tr>
              <th>{Messages.name}</th>
              <th>{Messages.description}</th>
              <th>{Messages.status}</th>
              <th>{Messages.interval}</th>
              <th>{Messages.actions}</th>
            </tr>
          </thead>
          <tbody data-testid="db-checks-all-checks-tbody">
            <ChecksReloadContext.Provider value={{ fetchChecks }}>
              {checks?.map((check) => (
                <CheckTableRow key={check.name} check={check} onSuccess={updateUI} />
              ))}
            </ChecksReloadContext.Provider>
          </tbody>
        </table>
      )}
    </div>
  );
};
