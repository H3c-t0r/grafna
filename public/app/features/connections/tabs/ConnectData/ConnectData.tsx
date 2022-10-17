import { css } from '@emotion/css';
import React, { useMemo, useState } from 'react';

import { LoadingPlaceholder, useStyles2 } from '@grafana/ui';
import { useGetAllWithFilters } from 'app/features/plugins/admin/state/hooks';

import { CardGrid } from './CardGrid';
import { CategoryHeader } from './CategoryHeader';
import { NoResults } from './NoResults';
import { Search } from './Search';

const getStyles = () => ({
  spacer: css`
    height: 16px;
  `,
});

export function ConnectData() {
  const [searchTerm, setSearchTerm] = useState('');
  const styles = useStyles2(getStyles);

  const handleSearchChange = (e: React.FormEvent<HTMLInputElement>) => {
    setSearchTerm(e.currentTarget.value.toLowerCase());
  };

  const { isLoading, error, plugins } = useGetAllWithFilters({ query: searchTerm, filterBy: '' });

  const cardGridItems = useMemo(
    () => plugins.map((plugin) => ({ id: plugin.id, name: plugin.name, logo: plugin.info.logos.small })),
    [plugins]
  );

  return (
    <>
      <Search onChange={handleSearchChange} />
      {/* We need this extra spacing when there are no filters */}
      <div className={styles.spacer} />
      {isLoading ? (
        <LoadingPlaceholder text="Loading..." />
      ) : !!error ? (
        <div>Error: {error.message}</div>
      ) : (
        <>
          <CategoryHeader iconName="database" label="Data sources" />
          <CardGrid items={cardGridItems} />
          {plugins.length < 1 && <NoResults />}
        </>
      )}
    </>
  );
}
