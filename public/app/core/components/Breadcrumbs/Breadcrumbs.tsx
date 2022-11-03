import { css } from '@emotion/css';
import React from 'react';

import { GrafanaTheme2 } from '@grafana/data';
import { useStyles2 } from '@grafana/ui';

import { BreadcrumbItem } from './BreadcrumbItem';
import { Breadcrumb } from './types';

export interface Props {
  breadcrumbs: Breadcrumb[];
  className?: string;
}

export function Breadcrumbs({ breadcrumbs, className }: Props) {
  const styles = useStyles2(getStyles);

  return (
    <nav aria-label="Breadcrumbs" className={className}>
      <ol className={styles.breadcrumbs}>
        {breadcrumbs.map((breadcrumb, index) => (
          <BreadcrumbItem {...breadcrumb} isCurrent={index === breadcrumbs.length - 1} key={index} />
        ))}
      </ol>
    </nav>
  );
}

const getStyles = (theme: GrafanaTheme2) => {
  return {
    breadcrumbs: css({
      display: 'flex',
      alignItems: 'center',
      flexWrap: 'nowrap',
    }),
  };
};
