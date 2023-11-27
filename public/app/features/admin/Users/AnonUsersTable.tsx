import React, { useMemo } from 'react';

import {
  Avatar,
  CellProps,
  Column,
  FetchDataFunc,
  InteractiveTable,
  Pagination,
  Stack,
  Badge,
} from '@grafana/ui';
import { UserAnonymousDTO } from 'app/types';

type Cell<T extends keyof UserAnonymousDTO = keyof UserAnonymousDTO> = CellProps<UserAnonymousDTO, UserAnonymousDTO[T]>;

interface AnonUsersTableProps {
  users: UserAnonymousDTO[];
  showPaging?: boolean;
  totalPages: number;
  currentPage: number;
  fetchData?: FetchDataFunc<UserAnonymousDTO>;
}

export const AnonUsersTable = ({ users, showPaging, totalPages, currentPage, fetchData }: AnonUsersTableProps) => {
  const columns: Array<Column<UserAnonymousDTO>> = useMemo(
    () => [
      {
        id: 'avatarUrl',
        header: '',
        cell: ({ cell: { value } }: Cell<'avatarUrl'>) => value && <Avatar src={value} alt={'User avatar'} />,
      },
      {
        id: 'login',
        header: 'Login',
        cell: ({ cell: { value } }: Cell<'login'>) => "Anonymous",
      },
      {
        id: 'user_agent',
        header: 'User Agent',
        cell: ({ cell: { value } }: Cell<'user_agent'>) => value,
        sortType: 'string',
      },
      {
        id: 'lastSeenAt',
        header: 'Last active',
        cell: ({ cell: { value } }: Cell<'lastSeenAt'>) => value,
        sortType: (a, b) => new Date(a.original.updated_at).getTime() - new Date(b.original.updated_at).getTime(),
      },
      {
        id: 'client_ip',
        header: 'Origin IP (address)',
        cell: ({ cell: { value } }: Cell<'client_ip'>) => value && <Badge text={value} color='orange'  />,
      },
    ],
    []
  );
  // FIXME: pagining is implemented
  showPaging = false;
  return (
    <Stack direction={'column'} gap={2}>
      <InteractiveTable
        columns={columns}
        data={users}
        getRowId={(user) => user.device_id}
      />
      {showPaging && (
        <Stack justifyContent={'flex-end'}>
          <Pagination numberOfPages={totalPages} currentPage={currentPage} onNavigate={() => {}} />
        </Stack>
      )}
    </Stack>
  );
};
