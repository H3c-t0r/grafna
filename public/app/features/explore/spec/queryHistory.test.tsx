import { screen } from '@testing-library/dom';
import React from 'react';
import { of } from 'rxjs';

import { serializeStateToUrlParam } from '@grafana/data';
import { config } from '@grafana/runtime';

import { silenceConsoleOutput } from '../../../../test/core/utils/silenceConsoleOutput';

import {
  assertDataSourceFilterVisibility,
  assertLoadMoreQueryHistoryNotVisible,
  assertModalIsOpen,
  assertNoQueryHistory,
  assertQueryHistory,
  assertQueryHistoryComment,
  assertQueryHistoryElementsShown,
  assertQueryHistoryExists,
  assertQueryHistoryIsStarred,
  assertQueryHistoryTabIsSelected,
  assertStarredQueryHistoryExists,
} from './helper/assert';
import {
  commentQueryHistory,
  closeQueryHistory,
  deleteQueryHistory,
  inputQuery,
  loadMoreQueryHistory,
  openQueryHistory,
  runQuery,
  selectOnlyActiveDataSource,
  selectStarredTabFirst,
  starQueryHistory,
  switchToQueryHistoryTab,
  deleteStarredQueryHistoryItem,
  confirmDeletingStarredQueryHistoryItem,
  openDeleteModal,
  confirmDeletingInModal,
} from './helper/interactions';
import { makeLogsQueryResponse } from './helper/query';
import { setupExplore, tearDown, waitForExplore } from './helper/setup';

const fetchMock = jest.fn();
const postMock = jest.fn();
const getMock = jest.fn();
const reportInteractionMock = jest.fn();
jest.mock('@grafana/runtime', () => ({
  ...jest.requireActual('@grafana/runtime'),
  getBackendSrv: () => ({ fetch: fetchMock, post: postMock, get: getMock }),
  reportInteraction: (...args: object[]) => {
    reportInteractionMock(...args);
  },
}));

jest.mock('app/core/core', () => ({
  contextSrv: {
    hasPermission: () => true,
    hasAccess: () => true,
    isSignedIn: true,
  },
}));

jest.mock('app/core/services/PreferencesService', () => ({
  PreferencesService: function () {
    return {
      patch: jest.fn(),
      load: jest.fn().mockResolvedValue({
        queryHistory: {
          homeTab: 'query',
        },
      }),
    };
  },
}));

jest.mock('react-virtualized-auto-sizer', () => {
  return {
    __esModule: true,
    default(props: any) {
      return <div>{props.children({ width: 1000 })}</div>;
    },
  };
});

describe('Explore: Query History', () => {
  const USER_INPUT = 'my query';
  const RAW_QUERY = `{"expr":"${USER_INPUT}"}`;

  // silenceConsoleOutput();

  afterEach(() => {
    config.queryHistoryEnabled = false;
    fetchMock.mockClear();
    postMock.mockClear();
    getMock.mockClear();
    reportInteractionMock.mockClear();
    tearDown();
  });

  it.skip('adds new query history items after the query is run.', async () => {
    // when Explore is opened
    const { datasources, unmount } = setupExplore();
    jest.mocked(datasources.loki.query).mockReturnValueOnce(makeLogsQueryResponse());
    await waitForExplore();

    // and a user runs a query and opens query history
    await inputQuery(USER_INPUT);
    await runQuery();
    await openQueryHistory();

    // the query that was run is in query history
    await assertQueryHistoryExists(RAW_QUERY);

    // when Explore is opened again
    unmount();
    setupExplore({ clearLocalStorage: false });
    await waitForExplore();

    // previously added query is in query history
    await openQueryHistory();
    await assertQueryHistoryExists(RAW_QUERY);

    expect(reportInteractionMock).toBeCalledWith('grafana_explore_query_history_opened', {
      queryHistoryEnabled: false,
    });
  });

  it.skip('adds recently added query if the query history panel is already open', async () => {
    const urlParams = {
      left: serializeStateToUrlParam({
        datasource: 'loki',
        queries: [{ refId: 'A', expr: 'query #1' }],
        range: { from: 'now-1h', to: 'now' },
      }),
    };

    const { datasources } = setupExplore({ urlParams });
    jest.mocked(datasources.loki.query).mockReturnValueOnce(makeLogsQueryResponse());
    await waitForExplore();
    await openQueryHistory();

    await inputQuery('query #2');
    await runQuery();
    await assertQueryHistory(['{"expr":"query #2"}', '{"expr":"query #1"}']);
  });

  // FIXME: flaky test
  it.skip('updates the state in both Explore panes', async () => {
    const urlParams = {
      left: serializeStateToUrlParam({
        datasource: 'loki',
        queries: [{ refId: 'A', expr: 'query #1' }],
        range: { from: 'now-1h', to: 'now' },
      }),
      right: serializeStateToUrlParam({
        datasource: 'loki',
        queries: [{ refId: 'A', expr: 'query #2' }],
        range: { from: 'now-1h', to: 'now' },
      }),
    };

    const { datasources } = setupExplore({ urlParams });
    jest.mocked(datasources.loki.query).mockReturnValue(makeLogsQueryResponse());
    await waitForExplore();
    await waitForExplore('right');

    // queries in history
    await openQueryHistory('left');
    await assertQueryHistory(['{"expr":"query #2"}', '{"expr":"query #1"}'], 'left');
    await openQueryHistory('right');
    await assertQueryHistory(['{"expr":"query #2"}', '{"expr":"query #1"}'], 'right');

    // star one one query
    await starQueryHistory(1, 'left');
    await assertQueryHistoryIsStarred([false, true], 'left');
    await assertQueryHistoryIsStarred([false, true], 'right');
    expect(reportInteractionMock).toBeCalledWith('grafana_explore_query_history_starred', {
      queryHistoryEnabled: false,
      newValue: true,
    });

    await deleteQueryHistory(0, 'left');
    await assertQueryHistory(['{"expr":"query #1"}'], 'left');
    await assertQueryHistory(['{"expr":"query #1"}'], 'right');
    expect(reportInteractionMock).toBeCalledWith('grafana_explore_query_history_deleted', {
      queryHistoryEnabled: false,
    });
  });

  it.skip('add comments to query history', async () => {
    const urlParams = {
      left: serializeStateToUrlParam({
        datasource: 'loki',
        queries: [{ refId: 'A', expr: 'query #1' }],
        range: { from: 'now-1h', to: 'now' },
      }),
    };

    const { datasources } = setupExplore({ urlParams });
    jest.mocked(datasources.loki.query).mockReturnValueOnce(makeLogsQueryResponse());
    await waitForExplore();
    await openQueryHistory();
    await assertQueryHistory(['{"expr":"query #1"}'], 'left');

    await commentQueryHistory(0, 'test comment');
    await assertQueryHistoryComment(['test comment'], 'left');
  });

  it.skip('updates query history settings', async () => {
    // open settings page
    setupExplore();
    await waitForExplore();
    await openQueryHistory();

    // assert default values
    assertQueryHistoryTabIsSelected('Query history');
    assertDataSourceFilterVisibility(true);
    await switchToQueryHistoryTab('Settings');

    // change settings
    await selectStarredTabFirst();
    await selectOnlyActiveDataSource();
    await closeQueryHistory();
    await openQueryHistory();

    // assert new settings
    assertQueryHistoryTabIsSelected('Starred');
    assertDataSourceFilterVisibility(false);
  });

  it.skip('removes the query item from the history panel when the user deletes a regular query', async () => {
    const { datasources } = setupExplore();
    jest.mocked(datasources.loki.query).mockReturnValueOnce(makeLogsQueryResponse());

    // run query and open query history
    await waitForExplore();
    await inputQuery(USER_INPUT);
    await runQuery();
    await openQueryHistory();

    // the query that was run is in query history
    await assertQueryHistoryExists(RAW_QUERY);

    // delete the query
    await deleteQueryHistory(0);

    // the query is removed from the query history
    await assertNoQueryHistory();
  });

  it('removes the query item from the history panel when the user deletes a starred query', async () => {
    setupExplore();

    // run query and open query history
    await waitForExplore();
    await inputQuery(USER_INPUT);
    await runQuery();
    await openQueryHistory();

    // the query that was run is in query history
    await assertQueryHistoryExists(RAW_QUERY);

    // star the query
    await starQueryHistory(0);

    // the starred query exist in query history
    await assertStarredQueryHistoryExists();

    // delete the starred query
    await deleteQueryHistory(0);

    // modal is open with a question about deleting the query
    await assertModalIsOpen();

    // confirm deleting in modal
    await confirmDeletingInModal();

    // the query is removed from the query history
    await assertNoQueryHistory();
  });

  it.skip('pagination', async () => {
    config.queryHistoryEnabled = true;
    const { datasources } = setupExplore();
    jest.mocked(datasources.loki.query).mockReturnValueOnce(makeLogsQueryResponse());
    fetchMock.mockReturnValue(
      of({
        data: { result: { queryHistory: [{ datasourceUid: 'loki', queries: [{ expr: 'query' }] }], totalCount: 2 } },
      })
    );
    await waitForExplore();

    await openQueryHistory();
    await assertQueryHistory(['{"expr":"query"}']);
    assertQueryHistoryElementsShown(1, 2);

    await loadMoreQueryHistory();
    await assertQueryHistory(['{"expr":"query"}', '{"expr":"query"}']);

    assertLoadMoreQueryHistoryNotVisible();
  });
});
