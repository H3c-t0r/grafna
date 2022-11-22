import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';
import { Provider } from 'react-redux';
import { mockToolkitActionCreator } from 'test/core/redux/mocks';

import { NavModel } from '@grafana/data';
import { ModalManager } from 'app/core/services/ModalManager';
import { configureStore } from 'app/store/configureStore';

import { backendSrv } from '../../core/services/backend_srv';
import { Organization } from '../../types';

import { OrgDetailsPage, Props } from './OrgDetailsPage';
import { setOrganizationName } from './state/reducers';

jest.mock('app/core/core', () => {
  return {
    ...jest.requireActual('app/core/core'),
    contextSrv: {
      hasPermission: () => true,
    },
  };
});

jest.mock('@grafana/runtime', () => {
  const originalModule = jest.requireActual('@grafana/runtime');
  return {
    ...originalModule,
    config: {
      ...originalModule.config,
      featureToggles: {
        internationalization: true,
      },
    },
  };
});

const setup = (propOverrides?: object) => {
  jest.clearAllMocks();
  // needed because SharedPreferences is rendered in the test
  jest.spyOn(backendSrv, 'put');
  jest
    .spyOn(backendSrv, 'get')
    .mockResolvedValue({ timezone: 'UTC', homeDashboardUID: 'home-dashboard', theme: 'dark' });
  jest.spyOn(backendSrv, 'search').mockResolvedValue([]);

  const props: Props = {
    organization: {} as Organization,
    navModel: {
      main: {
        text: 'Configuration',
      },
      node: {
        text: 'Org details',
      },
    } as NavModel,
    loadOrganization: jest.fn(),
    setOrganizationName: mockToolkitActionCreator(setOrganizationName),
    updateOrganization: jest.fn(),
  };
  Object.assign(props, propOverrides);

  render(
    <Provider store={configureStore()}>
      <OrgDetailsPage {...props} />
    </Provider>
  );
};

describe('Render', () => {
  beforeEach(() => {
    jest.spyOn(console, 'error').mockImplementation(() => {});
  });

  it('should render component', () => {
    expect(() => setup()).not.toThrow();
  });

  it('should render organization and preferences', () => {
    expect(() =>
      setup({
        organization: {
          name: 'Cool org',
          id: 1,
        },
        preferences: {
          homeDashboardUID: 'home-dashboard',
          theme: 'Default',
          timezone: 'Default',
          locale: '',
        },
      })
    ).not.toThrow();
  });

  it('should show a modal when submitting', async () => {
    new ModalManager().init();
    setup({
      organization: {
        name: 'Cool org',
        id: 1,
      },
      preferences: {
        homeDashboardUID: 'home-dashboard',
        theme: 'Default',
        timezone: 'Default',
        locale: '',
      },
    });

    await userEvent.click(screen.getByRole('button', { name: 'Save' }));

    expect(screen.getByText('Confirm preferences update')).toBeInTheDocument();
  });
});
