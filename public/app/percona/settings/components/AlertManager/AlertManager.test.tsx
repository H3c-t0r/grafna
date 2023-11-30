import { render, screen, fireEvent, waitForElementToBeRemoved } from '@testing-library/react';
import React from 'react';
import { Provider } from 'react-redux';

import * as reducers from 'app/percona/shared/core/reducers';
import { wrapWithGrafanaContextMock } from 'app/percona/shared/helpers/testUtils';
import { configureStore } from 'app/store/configureStore';
import { StoreState } from 'app/types';

import { AlertManager } from './AlertManager';

jest.mock('app/percona/settings/Settings.service');

describe('AlertManager::', () => {
  it('Renders correctly with props', async () => {
    render(
      <Provider
        store={configureStore({
          percona: {
            user: { isAuthorized: true },
            settings: { loading: false, result: { alertManagerUrl: 'fake.url', alertManagerRules: 'rule' } },
          },
        } as StoreState)}
      >
        {wrapWithGrafanaContextMock(<AlertManager />)}
      </Provider>
    );

    expect(screen.getByDisplayValue('fake.url')).toBeInTheDocument();
    expect(screen.getByTestId('alertmanager-rules').textContent).toBe('rule');
  });

  it('Disables apply changes on initial values', () => {
    render(
      <Provider
        store={configureStore({
          percona: {
            user: { isAuthorized: true },
            settings: { loading: false, result: { alertManagerUrl: 'fake.url', alertManagerRules: 'rule' } },
          },
        } as StoreState)}
      >
        {wrapWithGrafanaContextMock(<AlertManager />)}
      </Provider>
    );

    expect(screen.getByTestId('alertmanager-button')).toBeDisabled();
  });

  it('Calls apply changes', async () => {
    const spy = jest.spyOn(reducers, 'updateSettingsAction');
    const { container } = render(
      <Provider
        store={configureStore({
          percona: {
            user: { isAuthorized: true },
            settings: { loading: false, result: { alertManagerUrl: 'fake.url', alertManagerRules: 'rule' } },
          },
        } as StoreState)}
      >
        {wrapWithGrafanaContextMock(<AlertManager />)}
      </Provider>
    );

    fireEvent.change(screen.getByTestId('alertmanager-rules'), { target: { value: 'new key' } });
    fireEvent.submit(screen.getByTestId('alertmanager-button'));
    await waitForElementToBeRemoved(() => container.querySelector('.fa-spin'));

    expect(spy).toHaveBeenCalled();
  });
});
