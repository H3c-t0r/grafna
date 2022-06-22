import React from 'react';

import { getMount } from 'app/percona/shared/helpers/testUtils';

import { dbClustersStub, mongoDBClusterConnectionStub } from '../__mocks__/dbClustersStubs';

import { DBClusterConnection } from './DBClusterConnection';

jest.mock('app/core/app_events');
jest.mock('../XtraDB.service');
jest.mock('../PSMDB.service');

jest.mock('@percona/platform-core', () => {
  const originalModule = jest.requireActual('@percona/platform-core');
  return {
    ...originalModule,
    logger: {
      error: jest.fn(),
    },
  };
});

describe('DBClusterConnection::', () => {
  it('renders correctly connection items', async () => {
    const root = await getMount(<DBClusterConnection dbCluster={dbClustersStub[0]} />);

    expect(root.find('[data-testid="cluster-connection-host"]')).toBeTruthy();
    expect(root.find('[data-testid="cluster-connection-port"]')).toBeTruthy();
    expect(root.find('[data-testid="cluster-connection-username"]')).toBeTruthy();
    expect(root.find('[data-testid="cluster-connection-password"]')).toBeTruthy();
  });
  it('renders correctly connection items with MongoDB cluster', async () => {
    const root = await getMount(<DBClusterConnection dbCluster={dbClustersStub[2]} />);

    root.update();

    const host = root.find('[data-testid="cluster-connection-host"]');

    expect(host).toBeTruthy();
    expect(host.text()).toContain(mongoDBClusterConnectionStub.host);
    expect(root.find('[data-testid="cluster-connection-port"]')).toBeTruthy();
    expect(root.find('[data-testid="cluster-connection-username"]')).toBeTruthy();
    expect(root.find('[data-testid="cluster-connection-password"]')).toBeTruthy();
  });
});
