/* eslint-disable react/jsx-no-undef */
import { css } from '@emotion/css';
import React, { useEffect, useRef, useCallback } from 'react';
import { FixedSizeList } from 'react-window';
import InfiniteLoader from 'react-window-infinite-loader';

import { GrafanaTheme2 } from '@grafana/data';
import { useStyles2 } from '@grafana/ui';

import { SearchItem } from '../../components/SearchItem';
import { useSearchKeyboardNavigation } from '../../hooks/useSearchKeyboardSelection';
import { SearchResultMeta } from '../../service';
import { DashboardSearchItemType, DashboardSectionItem } from '../../types';

import { SearchResultsProps } from './SearchResultsTable';

export const SearchResultsCards = React.memo(
  ({
    response,
    width,
    height,
    selection,
    selectionToggle,
    clearSelection,
    onTagSelected,
    onDatasourceChange,
    keyboardEvents,
  }: SearchResultsProps) => {
    const styles = useStyles2(getStyles);
    const infiniteLoaderRef = useRef<InfiniteLoader>(null);
    const listRef = useRef<FixedSizeList>(null);
    const highlightIndex = useSearchKeyboardNavigation(keyboardEvents, 0, response);

    // Scroll to the top and clear loader cache when the query results change
    useEffect(() => {
      if (infiniteLoaderRef.current) {
        infiniteLoaderRef.current.resetloadMoreItemsCache();
      }
      if (listRef.current) {
        listRef.current.scrollTo(0);
      }
    }, [response]);

    const onToggleChecked = useCallback(
      (item: DashboardSectionItem) => {
        if (selectionToggle) {
          selectionToggle('dashboard', item.uid!);
        }
      },
      [selectionToggle]
    );

    const RenderRow = useCallback(
      ({ index: rowIndex, style }) => {
        const meta = response.view.dataFrame.meta?.custom as SearchResultMeta;

        let className = '';
        if (rowIndex === highlightIndex.y) {
          className += ' ' + styles.selectedRow;
        }

        const item = response.view.get(rowIndex);
        let v: DashboardSectionItem = {
          uid: item.uid,
          title: item.name,
          url: item.url,
          uri: item.url,
          type: item.kind === 'folder' ? DashboardSearchItemType.DashFolder : DashboardSearchItemType.DashDB,
          id: 666, // do not use me!
          isStarred: false,
          tags: item.tags ?? [],
        };

        if (item.location) {
          const first = item.location.split('/')[0];
          const finfo = meta.locationInfo[first];
          if (finfo) {
            v.folderUid = item.location;
            v.folderTitle = finfo.name;
          }
        }

        if (selection && selectionToggle) {
          const type = v.type === DashboardSearchItemType.DashFolder ? 'folder' : 'dashboard';
          v = {
            ...v,
            checked: selection(type, v.uid!),
          };
        }
        return (
          <div style={style} key={item.uid} className={className}>
            <SearchItem
              item={v}
              onTagSelected={onTagSelected}
              onToggleChecked={onToggleChecked as any}
              editable={Boolean(selection != null)}
            />
          </div>
        );
      },
      [response.view, highlightIndex, styles, onTagSelected, selection, selectionToggle, onToggleChecked]
    );

    if (!response.totalRows) {
      return (
        <div className={styles.noData} style={{ width }}>
          No data
        </div>
      );
    }

    return (
      <div aria-label="Search result cards" style={{ width }}>
        <InfiniteLoader
          ref={infiniteLoaderRef}
          isItemLoaded={response.isItemLoaded}
          itemCount={response.totalRows}
          loadMoreItems={response.loadMoreItems}
        >
          {({ onItemsRendered }) => (
            <FixedSizeList
              ref={listRef}
              onItemsRendered={onItemsRendered}
              height={height}
              itemCount={response.totalRows}
              itemSize={72}
              width="100%"
              style={{ overflow: 'hidden auto' }}
            >
              {RenderRow}
            </FixedSizeList>
          )}
        </InfiniteLoader>
      </div>
    );
  }
);
SearchResultsCards.displayName = 'SearchResultsCards';

const getStyles = (theme: GrafanaTheme2) => {
  return {
    noData: css`
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 100%;
    `,
    selectedRow: css`
      border-left: 3px solid ${theme.colors.primary.border};
    `,
  };
};
