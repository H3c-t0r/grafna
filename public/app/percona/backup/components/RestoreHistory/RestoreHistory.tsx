import { logger } from '@percona/platform-core';
import React, { FC, useState, useMemo, useEffect } from 'react';
import { Column, Row } from 'react-table';

import { Table } from 'app/percona/integrated-alerting/components/Table';
import { ExpandableCell } from 'app/percona/shared/components/Elements/ExpandableCell';
import { useCancelToken } from 'app/percona/shared/components/hooks/cancelToken.hook';
import { DATABASE_LABELS } from 'app/percona/shared/core';
import { isApiCancelError } from 'app/percona/shared/helpers/api';

import { Messages } from '../../Backup.messages';
import { useRecurringCall } from '../../hooks/recurringCall.hook';
import { DetailedDate } from '../DetailedDate';
import { Status } from '../Status';

import { DATA_INTERVAL, LIST_RESTORES_CANCEL_TOKEN } from './RestoreHistory.constants';
import { RestoreHistoryService } from './RestoreHistory.service';
import { Restore } from './RestoreHistory.types';
import { RestoreHistoryDetails } from './RestoreHistoryDetails';

export const RestoreHistory: FC = () => {
  const [pending, setPending] = useState(true);
  const [data, setData] = useState<Restore[]>([]);
  const [generateToken] = useCancelToken();
  const [triggerTimeout] = useRecurringCall();
  const columns = useMemo(
    (): Column[] => [
      {
        Header: Messages.backupInventory.table.columns.name,
        accessor: 'name',
        id: 'name',
        width: '250px',
        Cell: ({ row, value }) => <ExpandableCell row={row} value={value} />,
      },
      {
        Header: Messages.backupInventory.table.columns.vendor,
        accessor: ({ vendor }: Restore) => DATABASE_LABELS[vendor],
        width: '150px',
      },
      {
        Header: Messages.restoreHistory.table.columns.started,
        accessor: 'started',
        Cell: ({ value }) => <DetailedDate date={value} />,
      },
      {
        Header: Messages.backupInventory.table.columns.location,
        accessor: 'locationName',
      },
      {
        Header: Messages.backupInventory.table.columns.status,
        accessor: 'status',
        Cell: ({ value }) => <Status status={value} />,
      },
    ],
    []
  );

  const renderSelectedSubRow = React.useCallback(
    (row: Row<Restore>) => (
      <RestoreHistoryDetails
        name={row.original.name}
        finished={row.original.finished}
        dataModel={row.original.dataModel}
      />
    ),
    []
  );

  useEffect(() => {
    const getData = async (showLoading = false) => {
      showLoading && setPending(true);

      try {
        const restores = await RestoreHistoryService.list(generateToken(LIST_RESTORES_CANCEL_TOKEN));
        setData(restores);
      } catch (e) {
        if (isApiCancelError(e)) {
          return;
        }
        logger.error(e);
      }
      setPending(false);
    };

    getData(true).then(() => triggerTimeout(getData, DATA_INTERVAL));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [triggerTimeout]);

  return (
    <Table
      columns={columns}
      data={data}
      totalItems={data.length}
      emptyMessage={Messages.restoreHistory.table.noData}
      pendingRequest={pending}
      renderExpandedRow={renderSelectedSubRow}
    />
  );
};
