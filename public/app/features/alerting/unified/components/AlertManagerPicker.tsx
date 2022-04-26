import { css } from '@emotion/css';
import React, { FC, useMemo } from 'react';

import { GrafanaTheme2, SelectableValue } from '@grafana/data';
import { Field, Select, useStyles2 } from '@grafana/ui';

import { AlertManagerDataSource } from '../utils/datasource';

interface Props {
  onChange: (alertManagerSourceName: string) => void;
  current?: string;
  disabled?: boolean;
  dataSources: AlertManagerDataSource[];
}

export const AlertManagerPicker: FC<Props> = ({ onChange, current, dataSources, disabled = false }) => {
  const styles = useStyles2(getStyles);

  const options: Array<SelectableValue<string>> = useMemo(() => {
    return dataSources.map((ds) => ({
      label: ds.displayName,
      value: ds.name,
      imgUrl: ds.imgUrl,
      meta: ds.meta,
    }));
  }, [dataSources]);

  return (
    <Field
      className={styles.field}
      label={disabled ? 'Alertmanager' : 'Choose Alertmanager'}
      disabled={disabled || options.length === 1}
      data-testid="alertmanager-picker"
    >
      <Select
        aria-label={disabled ? 'Alertmanager' : 'Choose Alertmanager'}
        menuShouldPortal
        width={29}
        className="ds-picker select-container"
        backspaceRemovesValue={false}
        onChange={(value) => value.value && onChange(value.value)}
        options={options}
        maxMenuHeight={500}
        noOptionsMessage="No datasources found"
        value={current}
        getOptionLabel={(o) => o.label}
      />
    </Field>
  );
};

const getStyles = (theme: GrafanaTheme2) => ({
  field: css`
    margin-bottom: ${theme.spacing(4)};
  `,
});
