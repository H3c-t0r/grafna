import React, { FC, MutableRefObject } from 'react';
import { css } from 'emotion';
import { FixedSizeList } from 'react-window';
import { GrafanaTheme } from '@grafana/data';
import { stylesFactory, useTheme, Spinner } from '@grafana/ui';
import { DashboardSection, OnToggleChecked, SearchLayout } from '../types';
import { getItemsHeight } from '../utils';
import { ITEM_HEIGHT } from '../constants';
import { useListHeight } from '../hooks/useListHeight';
import { SearchItem } from './SearchItem';
import { SectionHeader } from './SectionHeader';

export interface Props {
  editable?: boolean;
  loading?: boolean;
  onTagSelected: (name: string) => any;
  onToggleChecked?: OnToggleChecked;
  onToggleSection: (section: DashboardSection) => void;
  results: DashboardSection[] | undefined;
  layout?: string;
  wrapperRef?: MutableRefObject<HTMLDivElement | null>;
}

export const SearchResults: FC<Props> = ({
  editable,
  loading,
  onTagSelected,
  onToggleChecked,
  onToggleSection,
  results,
  wrapperRef,
  layout,
}) => {
  const theme = useTheme();
  const styles = getSectionStyles(theme);
  const listHeight = useListHeight(wrapperRef);

  const renderItems = (section: DashboardSection) => {
    if (!section.expanded && layout !== SearchLayout.List) {
      return null;
    }

    return section.items.map(item => (
      <SearchItem key={item.id} {...{ item, editable, onToggleChecked, onTagSelected }} />
    ));
  };

  if (loading) {
    return <Spinner className={styles.spinner} />;
  } else if (!results || !results.length) {
    return <h6>No dashboards matching your query were found.</h6>;
  }

  return (
    <div className="search-results-container">
      <ul className={styles.wrapper}>
        {results.map(section => {
          let height = getItemsHeight(section, listHeight);
          return (
            <li aria-label="Search section" className={styles.section} key={section.title}>
              <SectionHeader onSectionClick={onToggleSection} {...{ onToggleChecked, editable, section }} />
              {section.expanded && section.items.length && (
                <FixedSizeList
                  aria-label="Search items"
                  className={styles.wrapper}
                  innerElementType="ul"
                  itemSize={ITEM_HEIGHT}
                  height={height}
                  itemCount={section.items.length}
                  width="100%"
                >
                  {({ index, style }) => {
                    const item = section.items[index];
                    return (
                      <SearchItem style={style} key={item.id} {...{ item, editable, onToggleChecked, onTagSelected }} />
                    );
                  }}
                </FixedSizeList>
              )}
            </li>
          );
        })}
      </ul>
    </div>
  );
};

const getSectionStyles = stylesFactory((theme: GrafanaTheme) => {
  return {
    wrapper: css`
      list-style: none;
    `,
    section: css`
      background: ${theme.colors.panelBg};
      border-bottom: solid 1px ${theme.isLight ? theme.palette.gray95 : theme.palette.gray25};
      padding: 0px 4px 4px 4px;
      margin-bottom: 3px;
    `,
    spinner: css`
      display: flex;
      justify-content: center;
      align-items: center;
      min-height: 100px;
    `,
  };
});
