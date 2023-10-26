import { css } from '@emotion/css';
import { isEqual } from 'lodash';
import React, { memo, useCallback, useEffect, useMemo, useRef, useState } from 'react';

import { AbsoluteTimeRange, GrafanaTheme2, LogsSortOrder } from '@grafana/data';
import { reportInteraction } from '@grafana/runtime';
import { DataQuery, TimeZone } from '@grafana/schema';
import { Button, Icon, Spinner, useTheme2 } from '@grafana/ui';
import { TOP_BAR_LEVEL_HEIGHT } from 'app/core/components/AppChrome/types';

import { LogsNavigationPages } from './LogsNavigationPages';

type Props = {
  absoluteRange: AbsoluteTimeRange;
  timeZone: TimeZone;
  queries: DataQuery[];
  loading: boolean;
  visibleRange: AbsoluteTimeRange;
  logsSortOrder?: LogsSortOrder | null;
  onChangeTime: (range: AbsoluteTimeRange) => void;
  scrollToTopLogs: () => void;
  addResultsToCache: () => void;
  clearCache: () => void;
  loadMore?: (range: AbsoluteTimeRange) => void;
};

export type LogsPage = {
  logsRange: AbsoluteTimeRange;
  queryRange: AbsoluteTimeRange;
};

function LogsNavigation({
  absoluteRange,
  logsSortOrder,
  timeZone,
  loading,
  onChangeTime,
  scrollToTopLogs,
  visibleRange,
  queries,
  clearCache,
  addResultsToCache,
  loadMore,
}: Props) {
  const [pages, setPages] = useState<LogsPage[]>([]);

  // These refs are to determine, if we want to clear up logs navigation when totally new query is run
  const expectedQueriesRef = useRef<DataQuery[]>();
  const expectedRangeRef = useRef<AbsoluteTimeRange>();
  // This ref is to store range span for future queres based on firstly selected time range
  // e.g. if last 5 min selected, always run 5 min range
  const rangeSpanRef = useRef(0);

  const currentPageIndex = useMemo(
    () =>
      pages.findIndex((page) => {
        return page.queryRange.to === absoluteRange.to;
      }),
    [absoluteRange.to, pages]
  );

  const oldestLogsFirst = logsSortOrder === LogsSortOrder.Ascending;
  const onFirstPage = oldestLogsFirst ? currentPageIndex === pages.length - 1 : currentPageIndex === 0;
  const onLastPage = oldestLogsFirst ? currentPageIndex === 0 : currentPageIndex === pages.length - 1;
  const theme = useTheme2();
  const styles = getStyles(theme, oldestLogsFirst);

  // Main effect to set pages and index
  useEffect(() => {
    const newPage = { logsRange: visibleRange, queryRange: absoluteRange };
    let newPages: LogsPage[] = [];
    // We want to start new pagination if queries change or if absolute range is different than expected
    if (!isEqual(expectedRangeRef.current, absoluteRange) || !isEqual(expectedQueriesRef.current, queries)) {
      clearCache();
      setPages([newPage]);
      expectedQueriesRef.current = queries;
      rangeSpanRef.current = absoluteRange.to - absoluteRange.from;
    } else {
      setPages((pages) => {
        // Remove duplicates with new query
        newPages = pages.filter((page) => !isEqual(newPage.queryRange, page.queryRange));
        // Sort pages based on logsOrder so they visually align with displayed logs
        newPages = [...newPages, newPage].sort((a, b) => sortPages(a, b, logsSortOrder));
        return newPages;
      });
    }
  }, [visibleRange, absoluteRange, logsSortOrder, queries, clearCache, addResultsToCache]);

  useEffect(() => {
    const scrollView = document.querySelector('.scrollbar-view');
    if (!scrollView) {
      return;
    }
    const delta = 1;
    function handleScroll(e: Event) {
      if (!e.target || !loadMore) {
        return;
      }
      const target: HTMLDivElement = e.target as HTMLDivElement;
      const diff = (target.scrollHeight - target.scrollTop) - target.clientHeight;
      if (diff > delta) {
        return;
      }
      if (!onLastPage) {
        const indexChange = oldestLogsFirst ? -1 : 1;
          loadMore({
            from: pages[currentPageIndex + indexChange].queryRange.from,
            to: pages[currentPageIndex + indexChange].queryRange.to,
          });
      } else {
        loadMore({ from: visibleRange.from - rangeSpanRef.current, to: visibleRange.from });
      }
      scrollView?.removeEventListener('scroll', handleScroll);
    }
    scrollView.addEventListener('scroll', handleScroll);

    return () => {
      scrollView.removeEventListener('scroll', handleScroll);
    }
  }, [currentPageIndex, loadMore, oldestLogsFirst, onLastPage, pages, visibleRange.from]);

  const changeTime = useCallback(
    ({ from, to }: AbsoluteTimeRange) => {
      addResultsToCache();
      expectedRangeRef.current = { from, to };
      onChangeTime({ from, to });
    },
    [onChangeTime, addResultsToCache]
  );

  const sortPages = (a: LogsPage, b: LogsPage, logsSortOrder?: LogsSortOrder | null) => {
    if (logsSortOrder === LogsSortOrder.Ascending) {
      return a.queryRange.to > b.queryRange.to ? 1 : -1;
    }
    return a.queryRange.to > b.queryRange.to ? -1 : 1;
  };

  const handleLoadMore = () => {
    //If we are not on the last page, use next page's range
    reportInteraction('grafana_explore_logs_pagination_clicked', {
      pageType: 'olderLogsButton',
    });
    if (!onLastPage) {
      const indexChange = oldestLogsFirst ? -1 : 1;
      if (loadMore) {
        loadMore({
          from: pages[currentPageIndex + indexChange].queryRange.from,
          to: pages[currentPageIndex + indexChange].queryRange.to,
        });
      } else {
        changeTime({
          from: pages[currentPageIndex + indexChange].queryRange.from,
          to: pages[currentPageIndex + indexChange].queryRange.to,
        });
      }
    } else {
      //If we are on the last page, create new range
      if (loadMore) {
        loadMore({ from: visibleRange.from - rangeSpanRef.current, to: visibleRange.from });
      } else {
        changeTime({ from: visibleRange.from - rangeSpanRef.current, to: visibleRange.from });
      }
    }
  };

  const olderLogsButton = (
    <Button
      data-testid="olderLogsButton"
      className={styles.navButton}
      variant="secondary"
      onClick={handleLoadMore}
      disabled={loading}
    >
      <div className={styles.navButtonContent}>
        {loading ? <Spinner /> : <Icon name={oldestLogsFirst ? 'angle-up' : 'angle-down'} size="lg" />}
        Older logs
      </div>
    </Button>
  );

  const newerLogsButton = (
    <Button
      data-testid="newerLogsButton"
      className={styles.navButton}
      variant="secondary"
      onClick={() => {
        reportInteraction('grafana_explore_logs_pagination_clicked', {
          pageType: 'newerLogsButton',
        });
        //If we are not on the first page, use previous page's range
        if (!onFirstPage) {
          const indexChange = oldestLogsFirst ? 1 : -1;
          changeTime({
            from: pages[currentPageIndex + indexChange].queryRange.from,
            to: pages[currentPageIndex + indexChange].queryRange.to,
          });
        }
        //If we are on the first page, button is disabled and we do nothing
      }}
      disabled={loading || onFirstPage}
    >
      <div className={styles.navButtonContent}>
        {loading && <Spinner />}
        {onFirstPage || loading ? null : <Icon name={oldestLogsFirst ? 'angle-down' : 'angle-up'} size="lg" />}
        {onFirstPage ? 'Start of range' : 'Newer logs'}
      </div>
    </Button>
  );

  const onPageClick = useCallback(
    (page: LogsPage, pageNumber: number) => {
      reportInteraction('grafana_explore_logs_pagination_clicked', {
        pageType: 'page',
        pageNumber,
      });
      changeTime({ from: page.queryRange.from, to: page.queryRange.to });
      scrollToTopLogs();
    },
    [changeTime, scrollToTopLogs]
  );

  return (
    <div className={styles.navContainer}>
      {oldestLogsFirst ? olderLogsButton : newerLogsButton}
      <LogsNavigationPages
        pages={pages}
        currentPageIndex={currentPageIndex}
        oldestLogsFirst={oldestLogsFirst}
        timeZone={timeZone}
        loading={loading}
        onClick={onPageClick}
      />
      {oldestLogsFirst ? newerLogsButton : olderLogsButton}
      <Button
        data-testid="scrollToTop"
        className={styles.scrollToTopButton}
        variant="secondary"
        onClick={scrollToTopLogs}
        title="Scroll to top"
      >
        <Icon name="arrow-up" size="lg" />
      </Button>
    </div>
  );
}

export default memo(LogsNavigation);

const getStyles = (theme: GrafanaTheme2, oldestLogsFirst: boolean) => {
  const navContainerHeight = `calc(100vh - 2*${theme.spacing(2)} - 2*${TOP_BAR_LEVEL_HEIGHT}px)`;
  return {
    navContainer: css`
      max-height: ${navContainerHeight};
      display: flex;
      flex-direction: column;
      justify-content: ${oldestLogsFirst ? 'flex-start' : 'space-between'};
      position: sticky;
      top: ${theme.spacing(2)};
      right: 0;
    `,
    navButton: css`
      width: 58px;
      height: 68px;
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      line-height: 1;
    `,
    navButtonContent: css`
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      width: 100%;
      height: 100%;
      white-space: normal;
    `,
    scrollToTopButton: css`
      width: 40px;
      height: 40px;
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      margin-top: ${theme.spacing(1)};
    `,
  };
};
