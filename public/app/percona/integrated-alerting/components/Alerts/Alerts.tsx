/* eslint-disable react/display-name */
import { cx } from '@emotion/css';
import { logger, Chip } from '@percona/platform-core';
import React, { FC, useCallback, useEffect, useState } from 'react';
import { Cell, Column } from 'react-table';

import { Button, useStyles } from '@grafana/ui';
import { ExpandableCell } from 'app/percona/shared/components/Elements/ExpandableCell';
import { useCancelToken } from 'app/percona/shared/components/hooks/cancelToken.hook';
import { isApiCancelError } from 'app/percona/shared/helpers/api';

import { Messages } from '../../IntegratedAlerting.messages';
import { Severity } from '../Severity';
import { useStoredTablePageSize } from '../Table/Pagination';
import { Table } from '../Table/Table';

import { AlertDetails } from './AlertDetails/AlertDetails';
import { ALERT_RULE_TEMPLATES_TABLE_ID, GET_ALERTS_CANCEL_TOKEN } from './Alerts.constants';
import { AlertsService } from './Alerts.service';
import { getStyles } from './Alerts.styles';
import { Alert, AlertStatus, AlertTogglePayload } from './Alerts.types';
import { formatAlerts } from './Alerts.utils';
import { AlertsActions } from './AlertsActions';

const { noData, columns } = Messages.alerts.table;
const {
  activeSince: activeSinceColumn,
  labels: labelsColumn,
  lastNotified: lastNotifiedColumn,
  severity: severityColumn,
  state: stateColumn,
  summary: summaryColumn,
  actions: actionsColumn,
} = columns;

export const Alerts: FC = () => {
  const style = useStyles(getStyles);
  const [pendingRequest, setPendingRequest] = useState(true);
  const [data, setData] = useState<Alert[]>([]);
  const [pageSize, setPageSize] = useStoredTablePageSize(ALERT_RULE_TEMPLATES_TABLE_ID);
  const [pageIndex, setPageindex] = useState(0);
  const [totalItems, setTotalItems] = useState(0);
  const [totalPages, setTotalPages] = useState(0);
  const [generateToken] = useCancelToken();

  const getAlerts = useCallback(async () => {
    setPendingRequest(true);
    try {
      const { alerts, totals } = await AlertsService.list(pageSize, pageIndex, generateToken(GET_ALERTS_CANCEL_TOKEN));
      setData(formatAlerts(alerts));
      setTotalItems(totals.total_items || 0);
      setTotalPages(totals.total_pages || 0);
    } catch (e) {
      if (isApiCancelError(e)) {
        return;
      }
      logger.error(e);
    }
    setPendingRequest(false);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pageSize, pageIndex]);

  const columns = React.useMemo(
    (): Array<Column<Alert>> => [
      {
        Header: summaryColumn,
        accessor: 'summary',
        width: '30%',
        Cell: ({ row, value }) => <ExpandableCell row={row} value={value} />,
      },
      {
        Header: severityColumn,
        accessor: 'severity',
        Cell: ({ row, value }) => (
          <Severity
            severity={value}
            className={cx({ [style.silencedSeverity]: (row.original as Alert).status === AlertStatus.SILENCED })}
          />
        ),
        width: '5%',
      },
      {
        Header: stateColumn,
        accessor: 'status',
        width: '5%',
      },
      {
        Header: labelsColumn,
        accessor: ({ labels }: Alert) => (
          <div className={style.labelsWrapper}>
            {labels.primary.map((label) => (
              <Chip text={label} key={label} />
            ))}
          </div>
        ),
        width: '40%',
      },
      {
        Header: activeSinceColumn,
        accessor: 'activeSince',
        width: '10%',
      },
      {
        Header: lastNotifiedColumn,
        accessor: 'lastNotified',
        width: '10%',
      },
      {
        Header: actionsColumn,
        accessor: (alert: Alert) => <AlertsActions alert={alert} getAlerts={getAlerts} />,
      },
    ],
    [style, getAlerts]
  );

  const getCellProps = useCallback(
    (cell: Cell<Alert>) => ({
      className: cell.row.original.status === AlertStatus.SILENCED ? style.disabledRow : '',
      key: cell.row.original.alertId,
    }),
    [style.disabledRow]
  );

  const handlePaginationChanged = useCallback(
    (pageSize: number, pageIndex: number) => {
      setPageSize(pageSize);
      setPageindex(pageIndex);
    },
    [setPageSize]
  );

  const renderSelectedSubRow = React.useCallback(
    (row: Row<Alert>) => (
      <AlertDetails ruleExpression={row.original.rule?.expr} labels={row.original.labels.secondary} />
    ),
    []
  );

  const handleSilenceAll = useCallback(async (silenced: AlertTogglePayload['silenced']) => {
    await AlertsService.toggle({ silenced, alert_ids: [] });
    getAlerts();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    getAlerts();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <>
      <div className={style.actionsWrapper}>
        <Button
          size="md"
          icon="bell-slash"
          variant="link"
          onClick={() => handleSilenceAll('TRUE')}
          data-testid="alert-rule-template-add-modal-button"
        >
          {Messages.alerts.silenceAllAction}
        </Button>
        <Button
          size="md"
          icon="bell"
          variant="link"
          onClick={() => handleSilenceAll('FALSE')}
          data-testid="alert-rule-template-add-modal-button"
        >
          {Messages.alerts.unsilenceAllAction}
        </Button>
      </div>
      <Table
        showPagination
        totalItems={totalItems}
        totalPages={totalPages}
        pageSize={pageSize}
        pageIndex={pageIndex}
        onPaginationChanged={handlePaginationChanged}
        data={data}
        columns={columns}
        pendingRequest={pendingRequest}
        emptyMessage={noData}
        getCellProps={getCellProps}
        renderExpandedRow={renderSelectedSubRow}
      />
    </>
  );
};
