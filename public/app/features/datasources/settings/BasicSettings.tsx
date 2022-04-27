import React, { FC } from 'react';

import { InlineField, InlineSwitch, Input } from '@grafana/ui';

export interface Props {
  dataSourceName: string;
  isDefault: boolean;
  onNameChange: (name: string) => void;
  onDefaultChange: (value: boolean) => void;
}

const BasicSettings: FC<Props> = ({ dataSourceName, isDefault, onDefaultChange, onNameChange }) => {
  return (
    <div className="gf-form-group" aria-label="Datasource settings page basic settings">
      <div className="gf-form-inline">
        <div className="gf-form max-width-30">
          <InlineField
            label="Name"
            tooltip="The name is used when you select the data source in panels. The default data source is
              'preselected in new panels."
            grow
          >
            <Input
              id="basic-settings-name"
              type="text"
              value={dataSourceName}
              placeholder="Name"
              onChange={(event) => onNameChange(event.currentTarget.value)}
              required
            />
          </InlineField>
        </div>

        <InlineField label="Default" labelWidth={8}>
          <InlineSwitch
            id="basic-settings-default"
            value={isDefault}
            onChange={(event: React.FormEvent<HTMLInputElement>) => {
              onDefaultChange(event.currentTarget.checked);
            }}
          />
        </InlineField>
      </div>
    </div>
  );
};

export default BasicSettings;
