import { css } from '@emotion/css';
import React, { ReactNode, useEffect, useState } from 'react';

import { AbsoluteTimeRange, LogRowModel, TimeRange } from '@grafana/data';
import { convertRawToRange, isRelativeTimeRange } from '@grafana/data/src/datetime/rangeutil';
import { LogsSortOrder, TimeZone } from '@grafana/schema';

type Props = {
  children: ReactNode;
  loading: boolean;
  loadMoreLogs?: (range: AbsoluteTimeRange) => void;
  range: TimeRange;
  rows: LogRowModel[];
  scrollElement?: HTMLDivElement;
  sortOrder: LogsSortOrder;
  timeZone: TimeZone;
};

export const InfiniteScroll = ({ children, loading, loadMoreLogs, range, rows, scrollElement, sortOrder, timeZone }: Props) => {
  const [lowerOutOfRange, setLowerOutOfRange] = useState(false);
  const [upperOutOfRange, setUpperOutOfRange] = useState(false);
  const [lastScroll, setLastScroll] = useState(scrollElement?.scrollTop || 0);

  useEffect(() => {
    if (!scrollElement || !loadMoreLogs) {
      return;
    }

    function handleScroll(e: Event) {
      if (!scrollElement || !loadMoreLogs || !rows.length || loading) {
        return;
      }
      e.stopImmediatePropagation();
      setLastScroll(scrollElement.scrollTop);
      const scrollDirection = shouldLoadMore(scrollElement, lastScroll);
      if (scrollDirection === ScrollDirection.NoScroll) {
        return;
      } else if (scrollDirection === ScrollDirection.Top) {
        scrollTop();
      } else {
        scrollBottom();
      }
    }

    function scrollTop() {
      if (!canScrollTop(getVisibleRange(rows), range, timeZone, sortOrder)) {
        setUpperOutOfRange(true);
        return;
      }
      setUpperOutOfRange(false);
      const newRange = sortOrder === LogsSortOrder.Descending ? getNextRange(getVisibleRange(rows), range, timeZone) : getPrevRange(getVisibleRange(rows), range);
      loadMoreLogs?.(newRange);
    };
  
    function scrollBottom() {
      if (!canScrollBottom(getVisibleRange(rows), range, timeZone, sortOrder)) {
        setLowerOutOfRange(true);
        return;
      }
      setLowerOutOfRange(false);
      const newRange = sortOrder === LogsSortOrder.Descending ? getPrevRange(getVisibleRange(rows), range) : getNextRange(getVisibleRange(rows), range, timeZone);
      loadMoreLogs?.(newRange);
    }

    scrollElement.addEventListener('scroll', handleScroll);
    scrollElement.addEventListener('wheel', handleScroll);

    return () => {
      scrollElement.removeEventListener('scroll', handleScroll);
      scrollElement.removeEventListener('wheel', handleScroll);
    };
  }, [lastScroll, loadMoreLogs, loading, range, rows, scrollElement, sortOrder, timeZone]);

  return (
    <>
      {upperOutOfRange && outOfRangeMessage}
      {children}
      {lowerOutOfRange && outOfRangeMessage}
    </>
  );
};

const styles = {
  limitReached: css({
    textAlign: 'center',
    padding: 0.25,
  })
}

const outOfRangeMessage = <div className={styles.limitReached}>Limit reached for the current time range.</div>;

enum ScrollDirection {
  Top = -1,
  Bottom = 1,
  NoScroll = 0
}
function shouldLoadMore(element: HTMLDivElement, lastScroll: number): ScrollDirection {
  const delta = element.scrollTop - lastScroll;
  const scrollDirection = delta <= 0 ? ScrollDirection.Top : ScrollDirection.Bottom;
  const diff = scrollDirection === ScrollDirection.Top ? 
    element.scrollTop :
    element.scrollHeight - element.scrollTop - element.clientHeight;
  const coef = 1;
  
  return diff <= coef ? scrollDirection : ScrollDirection.NoScroll;
}

function getVisibleRange(rows: LogRowModel[]) {
  const firstTimeStamp = rows[0].timeEpochMs;
  const lastTimeStamp = rows[rows.length - 1].timeEpochMs;

  const visibleRange =
    lastTimeStamp < firstTimeStamp
      ? { from: lastTimeStamp, to: firstTimeStamp }
      : { from: firstTimeStamp, to: lastTimeStamp };

  return visibleRange;
}

function getPrevRange(visibleRange: AbsoluteTimeRange, currentRange: TimeRange) {
  return { from: currentRange.from.valueOf(), to: visibleRange.from };
}

function getNextRange(visibleRange: AbsoluteTimeRange, currentRange: TimeRange, timeZone: TimeZone) {
  // When requesting new logs, update the current range if using relative time ranges.
  currentRange = updateCurrentRange(currentRange, timeZone);
  return { from: visibleRange.to, to: currentRange.to.valueOf() };
}

// To get more logs, the difference between the visible range and the current range should be 1 second or more.
function canScrollTop(visibleRange: AbsoluteTimeRange, currentRange: TimeRange, timeZone: TimeZone, sortOrder: LogsSortOrder) {
  if (sortOrder === LogsSortOrder.Descending) {
    // When requesting new logs, update the current range if using relative time ranges.
    currentRange = updateCurrentRange(currentRange, timeZone);
    return (currentRange.to.valueOf() - visibleRange.to) > 1e3;
  }
  return Math.abs(currentRange.from.valueOf() - visibleRange.from) > 1e3;
}

function canScrollBottom(visibleRange: AbsoluteTimeRange, currentRange: TimeRange, timeZone: TimeZone, sortOrder: LogsSortOrder) {
  if (sortOrder === LogsSortOrder.Descending) {
    return Math.abs(currentRange.from.valueOf() - visibleRange.from) > 1e3;  
  }
  // When requesting new logs, update the current range if using relative time ranges.
  currentRange = updateCurrentRange(currentRange, timeZone);
  return (currentRange.to.valueOf() - visibleRange.to) > 1e3;
}

// Given a TimeRange, returns a new instance if using relative time, or else the same.
function updateCurrentRange(timeRange: TimeRange, timeZone: TimeZone) {
  return isRelativeTimeRange(timeRange.raw) ? convertRawToRange(timeRange.raw, timeZone) : timeRange;
}
