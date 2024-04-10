import React from 'react';

import { InlineField, RadioButtonGroup } from '@grafana/ui';

import { AzureQueryEditorFieldProps } from '../../types';

import { setBasicLogsQuery, setKustoQuery } from './setQueryValue';

export function LogsManagement({ query, onQueryChange: onChange }: AzureQueryEditorFieldProps) {
  return (
    <>
      <InlineField label="Logs" tooltip={<span>Specifies whether to run a Basic or Analytics Logs query.</span>}>
        <RadioButtonGroup
          options={[
            { label: 'Analytics', value: false },
            { label: 'Basic', value: true },
          ]}
          value={query.azureLogAnalytics?.basicLogsQuery ?? false}
          size={'md'}
          onChange={(val) => {
            const updatedBasicLogsQuery = setBasicLogsQuery(query, val);
            onChange(setKustoQuery(updatedBasicLogsQuery, ''));
          }}
        />
      </InlineField>
    </>
  );
}
