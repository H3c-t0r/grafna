import React, { FormEvent } from 'react';
import { Label, Tooltip, Input, Icon, useStyles2 } from '@grafana/ui';
import { GrafanaTheme2 } from '@grafana/data';
import { css } from '@emotion/css';

interface Props {
  className?: string;
  queryString?: string;
  defaultQueryString?: string;
  onFilterChange: (filterString: string) => void;
}

export const MatcherFilter = ({ className, onFilterChange, defaultQueryString, queryString }: Props) => {
  const styles = useStyles2(getStyles);
  const handleSearchChange = (e: FormEvent<HTMLInputElement>) => {
    const target = e.target as HTMLInputElement;
    onFilterChange(target.value);
  };
  const searchIcon = <Icon name={'search'} />;
  return (
    <div className={className}>
      <Label>
        <Tooltip
          content={
            <div>
              Filter alerts using label querying, ex:
              <pre>{`{severity="critical", instance=~"cluster-us-.+"}`}</pre>
            </div>
          }
        >
          <Icon className={styles.icon} name="info-circle" size="xs" />
        </Tooltip>
        Search by label
      </Label>
      <Input
        placeholder="Search"
        defaultValue={defaultQueryString}
        value={queryString}
        onChange={handleSearchChange}
        data-testid="search-query-input"
        prefix={searchIcon}
        className={styles.inputWidth}
      />
    </div>
  );
};

const getStyles = (theme: GrafanaTheme2) => ({
  icon: css`
    margin-right: ${theme.spacing(0.5)};
  `,
  inputWidth: css`
    width: 340px;
    flex-grow: 0;
  `,
});
