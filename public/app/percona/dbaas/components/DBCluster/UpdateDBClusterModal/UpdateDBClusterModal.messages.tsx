/* eslint-disable react/display-name */
import React, { ReactNode } from 'react';

import { Databases } from 'app/percona/shared/core';
import { DATABASE_LABELS } from 'app/percona/shared/core/constants';

export const Messages = {
  cancel: 'Cancel',
  confirm: 'Update',
  title: 'Confirm database update',
  buildUpdateDatabaseMessage: (
    databaseType: Databases,
    installedVersion: ReactNode,
    availableVersion: ReactNode,
    clusterName: ReactNode
  ) => (
    <>
      Are you sure you want to update {DATABASE_LABELS[databaseType]} {installedVersion} to version {availableVersion}{' '}
      in {clusterName} cluster?
    </>
  ),
};
