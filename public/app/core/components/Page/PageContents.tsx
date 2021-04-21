// Libraries
import React, { FC } from 'react';
import { cx } from '@emotion/css';

// Components
import PageLoader from '../PageLoader/PageLoader';

interface Props {
  isLoading?: boolean;
  children: React.ReactNode;
  className?: string;
}

export const PageContents: FC<Props> = ({ isLoading, children, className }) => {
  return <div className={cx('page-container', 'page-body', className)}>{isLoading ? <PageLoader /> : children}</div>;
};
