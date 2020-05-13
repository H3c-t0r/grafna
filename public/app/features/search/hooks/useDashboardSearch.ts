import { KeyboardEvent, useReducer } from 'react';
import { getLocationSrv } from '@grafana/runtime';
import { DashboardQuery, DashboardSearchItemType, DashboardSection } from '../types';
import { MOVE_SELECTION_DOWN, MOVE_SELECTION_UP } from '../reducers/actionTypes';
import { dashboardsSearchState, DashboardsSearchState, searchReducer } from '../reducers/dashboardSearch';
import { findSelected } from '../utils';
import { useSearch } from './useSearch';
import config from 'app/core/config';

export const useDashboardSearch = (query: DashboardQuery, onCloseSearch: () => void) => {
  const reducer = useReducer(searchReducer, dashboardsSearchState);
  const {
    state: { results, loading },
    onToggleSection,
    dispatch,
  } = useSearch<DashboardsSearchState>(query, reducer, { queryParsing: true });

  const onKeyDown = (event: KeyboardEvent<HTMLInputElement>) => {
    switch (event.key) {
      case 'Escape':
        onCloseSearch();
        break;
      case 'ArrowUp':
        dispatch({ type: MOVE_SELECTION_UP });
        break;
      case 'ArrowDown':
        dispatch({ type: MOVE_SELECTION_DOWN });
        break;
      case 'Enter':
        const selectedItem = findSelected(results);
        if (selectedItem) {
          if (selectedItem.type === DashboardSearchItemType.DashFolder) {
            onToggleSection(selectedItem as DashboardSection);
          } else {
            getLocationSrv().update({
              path: selectedItem.url.startsWith(config.appSubUrl)
                ? selectedItem.url.split(config.appSubUrl + '/')[1]
                : selectedItem.url,
            });
            // Delay closing to prevent current page flicker
            setTimeout(onCloseSearch, 0);
          }
        }
    }
  };

  return {
    results,
    loading,
    onToggleSection,
    onKeyDown,
  };
};
