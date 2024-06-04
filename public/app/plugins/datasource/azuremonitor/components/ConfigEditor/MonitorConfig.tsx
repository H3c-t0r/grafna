import React, { useMemo, useState } from 'react';
import { useEffectOnce } from 'react-use';

import { SelectableValue } from '@grafana/data';
import { config } from '@grafana/runtime';

import { getCredentials, updateCredentials } from '../../credentials';
import { AzureDataSourceSettings, AzureCredentials } from '../../types';

import { AzureCredentialsForm, getAzureCloudOptions } from './AzureCredentialsForm';
import { BasicLogsToggle } from './BasicLogsToggle';
import { DefaultSubscription } from './DefaultSubscription';

export interface Props {
  options: AzureDataSourceSettings;
  updateOptions: (optionsFunc: (options: AzureDataSourceSettings) => AzureDataSourceSettings) => void;
  getSubscriptions: () => Promise<Array<SelectableValue<string>>>;
}

export const MonitorConfig = (props: Props) => {
  const { updateOptions, getSubscriptions, options } = props;
  const [subscriptions, setSubscriptions] = useState<Array<SelectableValue<string>>>([]);
  const credentials = useMemo(() => getCredentials(props.options), [props.options]);

  const onCredentialsChange = (credentials: AzureCredentials, subscriptionId?: string): void => {
    if (!subscriptionId) {
      setSubscriptions([]);
    }
    updateOptions((options) =>
      updateCredentials({ ...options, jsonData: { ...options.jsonData, subscriptionId } }, credentials)
    );
  };

  const onSubscriptionsChange = (receivedSubscriptions: Array<SelectableValue<string>>) =>
    setSubscriptions(receivedSubscriptions);

  const onSubscriptionChange = (subscriptionId?: string) =>
    updateOptions((options) => ({ ...options, jsonData: { ...options.jsonData, subscriptionId } }));

  const onBasicLogsEnabledChange = (enableBasicLogs: boolean) =>
    updateOptions((options) => ({ ...options, jsonData: { ...options.jsonData, basicLogsEnabled: enableBasicLogs } }));

  // The auth type needs to be set on the first load of the data source
  useEffectOnce(() => {
    if (!options.jsonData.authType) {
      onCredentialsChange(credentials);
    }
  });

  return (
    <>
      <AzureCredentialsForm
        managedIdentityEnabled={config.azure.managedIdentityEnabled}
        workloadIdentityEnabled={config.azure.workloadIdentityEnabled}
        userIdentityEnabled={config.azure.userIdentityEnabled}
        credentials={credentials}
        azureCloudOptions={getAzureCloudOptions()}
        onCredentialsChange={onCredentialsChange}
        disabled={props.options.readOnly}
      >
        <>
          <DefaultSubscription
            subscriptions={subscriptions}
            credentials={credentials}
            getSubscriptions={getSubscriptions}
            disabled={props.options.readOnly}
            onSubscriptionsChange={onSubscriptionsChange}
            onSubscriptionChange={onSubscriptionChange}
            options={options.jsonData}
          />
          <BasicLogsToggle options={options.jsonData} onBasicLogsEnabledChange={onBasicLogsEnabledChange} />
        </>
      </AzureCredentialsForm>
    </>
  );
};

export default MonitorConfig;
