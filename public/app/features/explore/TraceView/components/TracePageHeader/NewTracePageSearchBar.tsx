// Copyright (c) 2018 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import { css } from '@emotion/css';
import React, { memo, Dispatch, SetStateAction, useEffect, useMemo, useState } from 'react';

import { config, reportInteraction } from '@grafana/runtime';
import { Button, Icon, Switch, Tooltip, useStyles2 } from '@grafana/ui';

import { SearchProps } from '../../useSearch';
import { convertTimeFilter } from '../utils/filter-spans';

export type TracePageSearchBarProps = {
  search: SearchProps;
  spanFilterMatches: Set<string> | undefined;
  showSpanFilterMatchesOnly: boolean;
  setShowSpanFilterMatchesOnly: (showMatchesOnly: boolean) => void;
  focusedSpanIdForSearch: string;
  setFocusedSpanIdForSearch: Dispatch<SetStateAction<string>>;
  datasourceType: string;
  reset: () => void;
  totalSpans: number;
};

export default memo(function NewTracePageSearchBar(props: TracePageSearchBarProps) {
  const {
    search,
    spanFilterMatches,
    showSpanFilterMatchesOnly,
    setShowSpanFilterMatchesOnly,
    focusedSpanIdForSearch,
    setFocusedSpanIdForSearch,
    datasourceType,
    reset,
    totalSpans,
  } = props;
  const [currentSpanIndex, setCurrentSpanIndex] = useState<number>();
  const styles = useStyles2(getStyles);

  useEffect(() => {
    setFocusedSpanIdForSearch('');
  }, [search, setFocusedSpanIdForSearch]);

  useEffect(() => {
    setCurrentSpanIndex(undefined);
  }, [spanFilterMatches?.size]);

  const nextResult = () => {
    reportInteraction('grafana_traces_trace_view_find_next_prev_clicked', {
      datasourceType: datasourceType,
      grafana_version: config.buildInfo.version,
      direction: 'next',
    });

    const spanMatches = Array.from(spanFilterMatches!);
    const prevMatchedIndex = spanMatches.indexOf(focusedSpanIdForSearch);

    // new query || at end, go to start
    if (prevMatchedIndex === -1 || prevMatchedIndex === spanMatches.length - 1) {
      setFocusedSpanIdForSearch(spanMatches[0]);
      setCurrentSpanIndex(1);
      return;
    }

    // get next
    setFocusedSpanIdForSearch(spanMatches[prevMatchedIndex + 1]);
    setCurrentSpanIndex(prevMatchedIndex + 2);
  };

  const prevResult = () => {
    reportInteraction('grafana_traces_trace_view_find_next_prev_clicked', {
      datasourceType: datasourceType,
      grafana_version: config.buildInfo.version,
      direction: 'prev',
    });

    const spanMatches = Array.from(spanFilterMatches!);
    const prevMatchedIndex = spanMatches.indexOf(focusedSpanIdForSearch);

    // new query || at start, go to end
    if (prevMatchedIndex === -1 || prevMatchedIndex === 0) {
      setFocusedSpanIdForSearch(spanMatches[spanMatches.length - 1]);
      setCurrentSpanIndex(spanMatches.length);
      return;
    }

    // get prev
    setFocusedSpanIdForSearch(spanMatches[prevMatchedIndex - 1]);
    setCurrentSpanIndex(prevMatchedIndex);
  };

  const buttonEnabled = spanFilterMatches && spanFilterMatches?.size > 0;
  const resetEnabled = useMemo(() => {
    return (
      (search.serviceName && search.serviceName !== '') ||
      (search.spanName && search.spanName !== '') ||
      convertTimeFilter(search.from || '') ||
      convertTimeFilter(search.to || '') ||
      search.tags.length > 1 ||
      search.tags.some((tag) => {
        return tag.key;
      })
    );
  }, [search.serviceName, search.spanName, search.from, search.to, search.tags]);

  const matchesText = () => {
    if (spanFilterMatches?.size === 0) {
      return (
        <>
          <span>0 matches</span>
          <Tooltip
            content="There are 0 span matches for the filters selected. Please try removing some of the selected filters."
            placement="left"
          >
            <span className={styles.matchesTooltip}>
              <Icon name="info-circle" size="lg" />
            </span>
          </Tooltip>
        </>
      );
    }
    const amountText = spanFilterMatches?.size === 1 ? 'match' : 'matches';
    return currentSpanIndex
      ? `${currentSpanIndex}/${spanFilterMatches?.size} ${amountText}`
      : `${spanFilterMatches?.size} ${amountText}`;
  };

  return (
    <div className={styles.searchBar}>
      <div className={styles.buttons}>
        <>
          <div className={styles.resetButton}>
            <Button
              variant="destructive"
              disabled={!resetEnabled}
              type="button"
              fill="outline"
              aria-label="Reset filters button"
              onClick={reset}
            >
              Reset
            </Button>
            <div className={styles.matchesOnly}>
              <Switch
                value={showSpanFilterMatchesOnly}
                onChange={(value) => setShowSpanFilterMatchesOnly(value.currentTarget.checked ?? false)}
                label="Show matches only switch"
              />
              <span onClick={() => setShowSpanFilterMatchesOnly(!showSpanFilterMatchesOnly)}>Show matches only</span>
            </div>
          </div>
          <div className={styles.nextPrevButtons}>
            <span className={styles.matches}>{spanFilterMatches ? matchesText() : `${totalSpans} spans`}</span>
            <Button
              variant="secondary"
              disabled={!buttonEnabled}
              type="button"
              fill="outline"
              aria-label="Prev result button"
              onClick={prevResult}
            >
              Prev
            </Button>
            <Button
              variant="secondary"
              disabled={!buttonEnabled}
              type="button"
              fill="outline"
              aria-label="Next result button"
              onClick={nextResult}
            >
              Next
            </Button>
          </div>
        </>
      </div>
    </div>
  );
});

export const getStyles = () => {
  return {
    searchBar: css`
      display: inline;
    `,
    matchesOnly: css`
      display: inline-flex;
      margin: 0 0 0 10px;
      vertical-align: middle;

      span {
        cursor: pointer;
        margin: -3px 0 0 5px;
      }
    `,
    buttons: css`
      display: flex;
      justify-content: flex-end;
      margin: 5px 0 0 0;
    `,
    resetButton: css`
      order: 1;
    `,
    nextPrevButtons: css`
      margin-left: auto;
      order: 2;

      button {
        margin-left: 8px;
      }
    `,
    matches: css`
      margin-right: 5px;
    `,
    matchesTooltip: css`
      color: #aaa;
      margin: -2px 0 0 10px;
    `,
  };
};
