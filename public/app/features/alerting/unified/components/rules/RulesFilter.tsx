import React, { FormEvent, useState } from 'react';
import { Button, Icon, Input, Label, RadioButtonGroup, useStyles } from '@grafana/ui';
import { GrafanaTheme, SelectableValue } from '@grafana/data';
import { css, cx } from '@emotion/css';
import { debounce } from 'lodash';

import { PromAlertingRuleState } from 'app/types/unified-alerting-dto';
import { useQueryParams } from 'app/core/hooks/useQueryParams';
import DataSourceSelect from '../DataSourceSelect';
import { getFiltersFromUrlParams } from '../../utils/misc';
import { DataSourcePicker } from '@grafana/runtime';

const RulesFilter = () => {
  const [queryParams, setQueryParams] = useQueryParams();
  // This key is used to force a rerender on the inputs when the filters are cleared
  const [filterKey, setFilterKey] = useState<number>(Math.floor(Math.random() * 100));
  const dataSourceKey = `dataSource-${filterKey}`;
  const queryStringKey = `queryString-${filterKey}`;

  const { dataSource, alertState, queryString } = getFiltersFromUrlParams(queryParams);

  const styles = useStyles(getStyles);
  const stateOptions = Object.entries(PromAlertingRuleState).map(([key, value]) => ({ label: key, value }));

  const handleDataSourceChange = (dataSourceValue: SelectableValue<string>) => {
    setQueryParams({ dataSource: dataSourceValue.value });
  };

  const handleQueryStringChange = debounce((e: FormEvent<HTMLInputElement>) => {
    const target = e.target as HTMLInputElement;
    setQueryParams({ queryString: target.value || null });
  }, 600);

  const handleAlertStateChange = (value: string) => {
    setQueryParams({ alertState: value });
  };

  const handleClearFiltersClick = () => {
    setQueryParams({
      alertState: null,
      queryString: null,
      dataSource: null,
    });
    setFilterKey(filterKey + 1);
  };

  const searchIcon = <Icon name={'search'} />;
  return (
    <div className={styles.container}>
      <div className={styles.inputWidth}>
        <Label>Select data source</Label>
        <DataSourcePicker alerting={true} />
        <DataSourceSelect key={dataSourceKey} value={dataSource} onChange={handleDataSourceChange} />
      </div>
      <div className={cx(styles.flexRow, styles.spaceBetween)}>
        <div className={styles.flexRow}>
          <div className={styles.rowChild}>
            <Label>Search by name or label</Label>
            <Input
              key={queryStringKey}
              className={styles.inputWidth}
              prefix={searchIcon}
              onChange={handleQueryStringChange}
              defaultValue={queryString}
            />
          </div>
          <div className={styles.rowChild}>
            <RadioButtonGroup options={stateOptions} value={alertState} onChange={handleAlertStateChange} />
          </div>
        </div>
        {(dataSource || alertState || queryString) && (
          <div className={styles.flexRow}>
            <Button fullWidth={false} icon="times" variant="secondary" onClick={handleClearFiltersClick}>
              Clear filters
            </Button>
          </div>
        )}
      </div>
    </div>
  );
};

const getStyles = (theme: GrafanaTheme) => {
  return {
    container: css`
      display: flex;
      flex-direction: column;
      border-bottom: 1px solid ${theme.colors.border1};
      padding-bottom: ${theme.spacing.sm};
      margin-bottom: ${theme.spacing.sm};

      & > div {
        margin-bottom: ${theme.spacing.sm};
      }
    `,
    inputWidth: css`
      width: 340px;
      flex-grow: 0;
    `,
    flexRow: css`
      display: flex;
      flex-direction: row;
      align-items: flex-end;
    `,
    spaceBetween: css`
      justify-content: space-between;
    `,
    rowChild: css`
      & + & {
        margin-left: ${theme.spacing.sm};
      }
    `,
    clearButton: css`
      align-self: flex-end;
    `,
  };
};

export default RulesFilter;
