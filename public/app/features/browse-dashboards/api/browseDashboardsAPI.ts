import { BaseQueryFn, createApi } from '@reduxjs/toolkit/query/react';
import { lastValueFrom } from 'rxjs';

import { isTruthy, locationUtil } from '@grafana/data';
import { BackendSrvRequest, getBackendSrv, locationService } from '@grafana/runtime';
import { notifyApp } from 'app/core/actions';
import { createSuccessNotification } from 'app/core/copy/appNotification';
import { contextSrv } from 'app/core/core';
import { SaveDashboardCommand } from 'app/features/dashboard/components/SaveDashboard/types';
import { dashboardWatcher } from 'app/features/live/dashboard/dashboardWatcher';
import { DashboardViewItem } from 'app/features/search/types';
import { DashboardDTO, DescendantCount, DescendantCountDTO, FolderDTO, SaveDashboardResponseDTO } from 'app/types';

import { refetchChildren } from '../state';
import { findItem } from '../state/utils';
import { DashboardTreeSelection, DashboardViewItemCollection } from '../types';

import { PAGE_SIZE, ROOT_PAGE_SIZE } from './services';

interface RequestOptions extends BackendSrvRequest {
  manageError?: (err: unknown) => { error: unknown };
  showErrorAlert?: boolean;
}

interface DeleteItemsArgs {
  selectedItems: Omit<DashboardTreeSelection, 'panel' | '$all'>;
  rootItems: DashboardViewItem[] | undefined;
  childrenByParentUID: Record<string, DashboardViewItemCollection | undefined>;
}

interface MoveItemsArgs extends DeleteItemsArgs {
  destinationUID: string;
}

function createBackendSrvBaseQuery({ baseURL }: { baseURL: string }): BaseQueryFn<RequestOptions> {
  async function backendSrvBaseQuery(requestOptions: RequestOptions) {
    try {
      const { data: responseData, ...meta } = await lastValueFrom(
        getBackendSrv().fetch({
          ...requestOptions,
          url: baseURL + requestOptions.url,
          showErrorAlert: requestOptions.showErrorAlert,
        })
      );
      return { data: responseData, meta };
    } catch (error) {
      return requestOptions.manageError ? requestOptions.manageError(error) : { error };
    }
  }

  return backendSrvBaseQuery;
}

export const browseDashboardsAPI = createApi({
  tagTypes: ['getFolder'],
  reducerPath: 'browseDashboardsAPI',
  baseQuery: createBackendSrvBaseQuery({ baseURL: '/api' }),
  endpoints: (builder) => ({
    deleteFolder: builder.mutation<void, FolderDTO>({
      invalidatesTags: (_result, _error, folder) => [{ type: 'getFolder', id: folder.uid }],
      query: ({ uid }) => ({
        url: `/folders/${uid}`,
        method: 'DELETE',
        params: {
          // TODO: Once backend returns alert rule counts, set this back to true
          // when this is merged https://github.com/grafana/grafana/pull/67259
          forceDeleteRules: false,
        },
      }),
      onQueryStarted: (arg, { queryFulfilled, dispatch }) => {
        const { parentUid } = arg;
        queryFulfilled.then(() => {
          dispatch(
            refetchChildren({
              parentUID: parentUid,
              pageSize: parentUid ? PAGE_SIZE : ROOT_PAGE_SIZE,
            })
          );
        });
      },
    }),
    getFolder: builder.query<FolderDTO, string>({
      providesTags: (_result, _error, arg) => [{ type: 'getFolder', id: arg }],
      query: (folderUID) => ({ url: `/folders/${folderUID}`, params: { accesscontrol: true } }),
    }),
    moveFolder: builder.mutation<void, { folder: FolderDTO; destinationUID: string }>({
      invalidatesTags: (_result, _error, arg) => [{ type: 'getFolder', id: arg.folder.uid }],
      query: ({ folder, destinationUID }) => ({
        url: `/folders/${folder.uid}/move`,
        method: 'POST',
        data: { parentUID: destinationUID },
      }),
      onQueryStarted: (arg, { queryFulfilled, dispatch }) => {
        const { folder, destinationUID } = arg;
        const { parentUid } = folder;
        queryFulfilled.then(() => {
          dispatch(
            refetchChildren({
              parentUID: parentUid,
              pageSize: parentUid ? PAGE_SIZE : ROOT_PAGE_SIZE,
            })
          );
          dispatch(
            refetchChildren({
              parentUID: destinationUID,
              pageSize: destinationUID ? PAGE_SIZE : ROOT_PAGE_SIZE,
            })
          );
        });
      },
    }),
    newFolder: builder.mutation<FolderDTO, { title: string; parentUid?: string }>({
      invalidatesTags: (_result, _error, args) => [{ type: 'getFolder', id: args.parentUid }],
      query: ({ title, parentUid }) => ({
        method: 'POST',
        url: '/folders',
        data: {
          title,
          parentUid,
        },
      }),
      onQueryStarted: (arg, { queryFulfilled, dispatch }) => {
        const { parentUid } = arg;
        queryFulfilled.then(async ({ data: folder }) => {
          await contextSrv.fetchUserPermissions();
          dispatch(notifyApp(createSuccessNotification('Folder Created', 'OK')));
          dispatch(
            refetchChildren({
              parentUID: parentUid,
              pageSize: parentUid ? PAGE_SIZE : ROOT_PAGE_SIZE,
            })
          );
          locationService.push(locationUtil.stripBaseFromUrl(folder.url));
        });
      },
    }),
    saveFolder: builder.mutation<FolderDTO, FolderDTO>({
      invalidatesTags: (_result, _error, args) => [{ type: 'getFolder', id: args.uid }],
      query: (folder) => ({
        method: 'PUT',
        url: `/folders/${folder.uid}`,
        data: {
          title: folder.title,
          version: folder.version,
        },
      }),
      onQueryStarted: (arg, { queryFulfilled, dispatch }) => {
        const { parentUid } = arg;
        queryFulfilled.then(() => {
          dispatch(
            refetchChildren({
              parentUID: parentUid,
              pageSize: parentUid ? PAGE_SIZE : ROOT_PAGE_SIZE,
            })
          );
        });
      },
    }),
    deleteItems: builder.mutation<void, DeleteItemsArgs>({
      invalidatesTags: (_result, _error, args) => {
        const { selectedItems } = args;
        const selectedFolders = Object.keys(selectedItems.folder).filter((uid) => selectedItems.folder[uid]);
        return selectedFolders.map((folderUID) => ({ type: 'getFolder', id: folderUID }));
      },
      queryFn: async (args, _api, _extraOptions, baseQuery) => {
        const { selectedItems } = args;
        const selectedDashboards = Object.keys(selectedItems.dashboard).filter((uid) => selectedItems.dashboard[uid]);
        const selectedFolders = Object.keys(selectedItems.folder).filter((uid) => selectedItems.folder[uid]);
        // Delete all the folders sequentially
        // TODO error handling here
        for (const folderUID of selectedFolders) {
          await baseQuery({
            url: `/folders/${folderUID}`,
            method: 'DELETE',
            params: {
              // TODO: Once backend returns alert rule counts, set this back to true
              // when this is merged https://github.com/grafana/grafana/pull/67259
              forceDeleteRules: false,
            },
          });
        }

        // Delete all the dashboards sequentially
        // TODO error handling here
        for (const dashboardUID of selectedDashboards) {
          await baseQuery({
            url: `/dashboards/uid/${dashboardUID}`,
            method: 'DELETE',
          });
        }
        return { data: undefined };
      },
      onQueryStarted: (arg, { queryFulfilled, dispatch }) => {
        const { selectedItems, rootItems, childrenByParentUID } = arg;
        const selectedDashboards = Object.keys(selectedItems.dashboard).filter((uid) => selectedItems.dashboard[uid]);
        const selectedFolders = Object.keys(selectedItems.folder).filter((uid) => selectedItems.folder[uid]);
        queryFulfilled.then(() => {
          const parentsToRefresh = new Set<string | undefined>();
          for (const folderUID of selectedFolders) {
            // find the parent folder uid to invalidate the children
            const folder = findItem(rootItems ?? [], childrenByParentUID, folderUID);
            parentsToRefresh.add(folder?.parentUID);
          }

          for (const dashboardUID of selectedDashboards) {
            // find the parent folder uid to invalidate the children
            const dashboard = findItem(rootItems ?? [], childrenByParentUID, dashboardUID);
            parentsToRefresh.add(dashboard?.parentUID);
          }
          for (const parentUID of parentsToRefresh) {
            dispatch(refetchChildren({ parentUID, pageSize: parentUID ? PAGE_SIZE : ROOT_PAGE_SIZE }));
          }
        });
      },
    }),
    moveItems: builder.mutation<void, MoveItemsArgs>({
      invalidatesTags: (_result, _error, args) => {
        const { selectedItems } = args;
        const selectedFolders = Object.keys(selectedItems.folder).filter((uid) => selectedItems.folder[uid]);
        return selectedFolders.map((folderUID) => ({ type: 'getFolder', id: folderUID }));
      },
      queryFn: async (args, _api, _extraOptions, baseQuery) => {
        const { selectedItems, destinationUID } = args;
        const selectedDashboards = Object.keys(selectedItems.dashboard).filter((uid) => selectedItems.dashboard[uid]);
        const selectedFolders = Object.keys(selectedItems.folder).filter((uid) => selectedItems.folder[uid]);

        // Move all the folders sequentially
        // TODO error handling here
        for (const folderUID of selectedFolders) {
          await baseQuery({
            url: `/folders/${folderUID}/move`,
            method: 'POST',
            data: { parentUID: destinationUID },
          });
        }

        // Move all the dashboards sequentially
        // TODO error handling here
        for (const dashboardUID of selectedDashboards) {
          const fullDash: DashboardDTO = await getBackendSrv().get(`/api/dashboards/uid/${dashboardUID}`);

          const options = {
            dashboard: fullDash.dashboard,
            folderUid: destinationUID,
            overwrite: false,
            message: '',
          };

          await baseQuery({
            url: `/dashboards/db`,
            method: 'POST',
            data: options,
          });
        }
        return { data: undefined };
      },
      onQueryStarted: (arg, { queryFulfilled, dispatch }) => {
        const { destinationUID, selectedItems, rootItems, childrenByParentUID } = arg;
        queryFulfilled.then(() => {
          const parentsToRefresh = new Set<string | undefined>();

          const selectedDashboards = Object.keys(selectedItems.dashboard).filter((uid) => selectedItems.dashboard[uid]);
          const selectedFolders = Object.keys(selectedItems.folder).filter((uid) => selectedItems.folder[uid]);

          parentsToRefresh.add(destinationUID);

          for (const folderUID of selectedFolders) {
            // find the parent folder uid to invalidate the children
            const folder = findItem(rootItems ?? [], childrenByParentUID, folderUID);
            parentsToRefresh.add(folder?.parentUID);
          }

          for (const dashboardUID of selectedDashboards) {
            // find the parent folder uid to invalidate the children
            const dashboard = findItem(rootItems ?? [], childrenByParentUID, dashboardUID);
            parentsToRefresh.add(dashboard?.parentUID);
          }

          for (const parentUID of parentsToRefresh) {
            dispatch(refetchChildren({ parentUID, pageSize: parentUID ? PAGE_SIZE : ROOT_PAGE_SIZE }));
          }
        });
      },
    }),
    getAffectedItems: builder.query<DescendantCount, DashboardTreeSelection>({
      queryFn: async (selectedItems) => {
        const folderUIDs = Object.keys(selectedItems.folder).filter((uid) => selectedItems.folder[uid]);

        const promises = folderUIDs.map((folderUID) => {
          return getBackendSrv().get<DescendantCountDTO>(`/api/folders/${folderUID}/counts`);
        });

        const results = await Promise.all(promises);

        const totalCounts = {
          folder: Object.values(selectedItems.folder).filter(isTruthy).length,
          dashboard: Object.values(selectedItems.dashboard).filter(isTruthy).length,
          libraryPanel: 0,
          alertRule: 0,
        };

        for (const folderCounts of results) {
          totalCounts.folder += folderCounts.folder;
          totalCounts.dashboard += folderCounts.dashboard;
          totalCounts.alertRule += folderCounts.alertrule ?? 0;

          // TODO enable these once the backend correctly returns them
          // totalCounts.libraryPanel += folderCounts.libraryPanel;
        }

        return { data: totalCounts };
      },
    }),
    saveDashboard: builder.mutation<SaveDashboardResponseDTO, SaveDashboardCommand>({
      query: ({ dashboard, folderUid, message, overwrite }) => ({
        url: `/dashboards/db`,
        method: 'POST',
        data: {
          dashboard,
          folderUid,
          message: message ?? '',
          overwrite: Boolean(overwrite),
        },
      }),
      onQueryStarted: (arg, { queryFulfilled, dispatch }) => {
        const { folderUid } = arg;
        dashboardWatcher.ignoreNextSave();
        queryFulfilled.then(async () => {
          await contextSrv.fetchUserPermissions();
          dispatch(
            refetchChildren({
              parentUID: folderUid,
              pageSize: folderUid ? PAGE_SIZE : ROOT_PAGE_SIZE,
            })
          );
        });
      },
    }),
  }),
});

export const {
  endpoints,
  useDeleteFolderMutation,
  useDeleteItemsMutation,
  useGetAffectedItemsQuery,
  useGetFolderQuery,
  useMoveFolderMutation,
  useMoveItemsMutation,
  useNewFolderMutation,
  useSaveDashboardMutation,
  useSaveFolderMutation,
} = browseDashboardsAPI;
export { skipToken } from '@reduxjs/toolkit/query/react';
