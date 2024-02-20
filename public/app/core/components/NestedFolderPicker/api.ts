import { createSelector } from '@reduxjs/toolkit';
import { QueryDefinition, BaseQueryFn } from '@reduxjs/toolkit/dist/query';
import { QueryActionCreatorResult } from '@reduxjs/toolkit/dist/query/core/buildInitiate';
import { RequestOptions } from 'http';
import { useCallback, useEffect, useMemo, useRef } from 'react';

import { ListFolderArgs, browseDashboardsAPI } from 'app/features/browse-dashboards/api/browseDashboardsAPI';
import { PAGE_SIZE } from 'app/features/browse-dashboards/api/services';
import { getPaginationPlaceholders } from 'app/features/browse-dashboards/state/utils';
import { DashboardViewItemWithUIItems, DashboardsTreeItem } from 'app/features/browse-dashboards/types';
import { RootState } from 'app/store/configureStore';
import { FolderDTO } from 'app/types';
import { useDispatch, useSelector } from 'app/types/store';

type TODOFolderListPage = ReturnType<typeof listFoldersSelector>;
type TODORequestPromise = QueryActionCreatorResult<
  QueryDefinition<ListFolderArgs, BaseQueryFn<RequestOptions>, 'getFolder', FolderDTO[], 'browseDashboardsAPI'>
>;

const createListFoldersSelector = createSelector(
  [
    (parentUid: ListFolderArgs['parentUid']) => parentUid,
    (parentUid: ListFolderArgs['parentUid'], page: ListFolderArgs['page']) => page,
    (parentUid: ListFolderArgs['parentUid'], page: ListFolderArgs['page'], limit: ListFolderArgs['limit']) => limit,
  ],
  (parentUid, page, limit) => {
    return browseDashboardsAPI.endpoints.listFolders.select({ parentUid, page, limit });
  }
);

const listFoldersSelector = createSelector(
  (state: RootState) => state,
  (
    state: RootState,
    parentUid: ListFolderArgs['parentUid'],
    page: ListFolderArgs['page'],
    limit: ListFolderArgs['limit']
  ) => createListFoldersSelector(parentUid, page, limit),
  (state, selectFolderList) => selectFolderList(state)
);

const listAllFoldersSelector = createSelector(
  [(state: RootState) => state, (state: RootState, requests: TODORequestPromise[]) => requests],
  (state: RootState, requests: TODORequestPromise[]) => {
    const seenRequests = new Set<string>();

    const rootPages: TODOFolderListPage[] = [];
    const pagesByParent: Record<string, TODOFolderListPage[]> = {};
    let isLoading = false;

    for (const req of requests) {
      if (seenRequests.has(req.requestId)) {
        continue;
      }

      const page = listFoldersSelector(state, req.arg.parentUid, req.arg.page, req.arg.limit);
      if (page.status === 'pending') {
        isLoading = true;
      }

      const parentUid = page.originalArgs?.parentUid;
      if (parentUid) {
        if (!pagesByParent[parentUid]) {
          pagesByParent[parentUid] = [];
        }

        pagesByParent[parentUid].push(page);
      } else {
        rootPages.push(page);
      }
    }

    return {
      isLoading,
      rootPages,
      pagesByParent,
    };
  }
);

/**
 * Returns the whether the set of pages are 'fully loaded', and the last page number
 */
function getPagesLoadStatus(pages: TODOFolderListPage[]): [boolean, number | undefined] {
  const lastPage = pages.at(-1);
  const lastPageNumber = lastPage?.originalArgs?.page;

  if (!lastPage?.data) {
    // If there's no pages yet, or the last page is still loading
    return [false, lastPageNumber];
  } else {
    return [lastPage.data.length < lastPage.originalArgs.limit, lastPageNumber];
  }
}

/**
 * Returns a loaded folder hierarchy as a flat list and a function to load more pages.
 */
export function useFolderList(isBrowsing: boolean, openFolders: Record<string, boolean>) {
  const dispatch = useDispatch();

  // Keep a list of all requests so we can
  //   a) unsubscribe from them when the component is unmounted
  //   b) use them to select the responses out of the state
  const requestsRef = useRef<TODORequestPromise[]>([]);

  const state = useSelector((rootState: RootState) => {
    return listAllFoldersSelector(rootState, requestsRef.current);
  });

  // Loads the next page of folders for the given parent UID by inspecting the
  // state to determine what the next page is
  const requestNextPage = useCallback(
    (parentUid: string | undefined) => {
      const pages = parentUid ? state.pagesByParent[parentUid] : state.rootPages;
      const [fullyLoaded, pageNumber] = getPagesLoadStatus(pages ?? []);
      if (fullyLoaded) {
        return;
      }

      const args = { parentUid, page: (pageNumber ?? 0) + 1, limit: PAGE_SIZE };
      const promise = dispatch(browseDashboardsAPI.endpoints.listFolders.initiate(args));

      // It's important that we create a new array so we can correctly memoize with it
      requestsRef.current = requestsRef.current.concat([promise]);
    },
    [state, dispatch]
  );

  // Unsubscribe from all requests when the component is unmounted
  useEffect(() => {
    return () => {
      for (const req of requestsRef.current) {
        req.unsubscribe();
      }
    };
  }, []);

  // Convert the individual responses into a flat list of folders, with level indicating
  // the depth in the hierarchy.
  // TODO: this will probably go up in the parent component so it can also do search
  const treeList = useMemo(() => {
    if (!isBrowsing) {
      return [];
    }

    function createFlatList(
      parentUid: string | undefined,
      pages: TODOFolderListPage[],
      level: number
    ): Array<DashboardsTreeItem<DashboardViewItemWithUIItems>> {
      const flatList = pages.flatMap((page) => {
        const pageItems = page.data ?? [];

        return pageItems.flatMap((item) => {
          const folderIsOpen = openFolders[item.uid];

          const flatItem: DashboardsTreeItem<DashboardViewItemWithUIItems> = {
            isOpen: Boolean(folderIsOpen),
            level: level,
            item: {
              kind: 'folder' as const,
              title: item.title,
              uid: item.uid,
            },
          };

          const childPages = folderIsOpen && state.pagesByParent[item.uid];
          if (childPages) {
            const childFlatItems = createFlatList(item.uid, childPages, level + 1);
            return [flatItem, ...childFlatItems];
          }

          return flatItem;
        });
      });

      const [fullyLoaded] = getPagesLoadStatus(pages);
      if (!fullyLoaded) {
        flatList.push(...getPaginationPlaceholders(PAGE_SIZE, parentUid, level));
      }

      return flatList;
    }

    const rootFlatTree = createFlatList(undefined, state.rootPages, 1);
    rootFlatTree.unshift(ROOT_FOLDER_ITEM);

    return rootFlatTree;
  }, [state, isBrowsing, openFolders]);

  return [treeList, state.isLoading, requestNextPage] as const;
}

const ROOT_FOLDER_ITEM = {
  isOpen: true,
  level: 0,
  item: {
    kind: 'folder' as const,
    title: 'Dashboards',
    uid: '',
  },
};
