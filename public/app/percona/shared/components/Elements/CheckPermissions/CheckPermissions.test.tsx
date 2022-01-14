import { dataQa } from '@percona/platform-core';
import React from 'react';

import { SettingsService } from 'app/percona/settings/Settings.service';
import { getMount } from 'app/percona/shared/helpers/testUtils';

import { CheckPermissions } from './CheckPermissions';

jest.mock('app/percona/settings/Settings.service');
jest.mock('@percona/platform-core', () => {
  const originalModule = jest.requireActual('@percona/platform-core');
  return {
    ...originalModule,
    logger: {
      error: jest.fn(),
    },
  };
});

describe('CheckPermissions::', () => {
  it('should render children', async () => {
    const wrapper = await getMount(
      <CheckPermissions>
        <div>Test</div>
      </CheckPermissions>
    );

    wrapper.update();

    expect(wrapper.find('div').text()).toEqual('Test');
  });

  it('should render unauthorized message', async () => {
    const errorObj = { response: { status: 401 } };
    jest.spyOn(SettingsService, 'getSettings').mockImplementationOnce(() => {
      throw errorObj;
    });
    const wrapper = await getMount(
      <CheckPermissions>
        <div>Test</div>
      </CheckPermissions>
    );

    wrapper.update();

    expect(wrapper.find(dataQa('unauthorized'))).not.toBeNull();
  });
});
