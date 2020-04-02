// Copyright (c) 2017 Uber Technologies, Inc.
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

import React from 'react';
import { css } from 'emotion';

import TimelineHeaderRow from './TimelineHeaderRow';
import VirtualizedTraceView from './VirtualizedTraceView';
import { merge as mergeShortcuts } from '../keyboard-shortcuts';
import { Accessors } from '../ScrollManager';
import { TUpdateViewRangeTimeFunction, ViewRange, ViewRangeTimeUpdate } from './types';
import { TNil } from '../types';
import { Span, Trace, Log, KeyValuePair, Link } from '../types/trace';
import TTraceTimeline from '../types/TTraceTimeline';
import { createStyle } from '../Theme';
import ExternalLinkContext from '../url/externalLinkContext';

type TExtractUiFindFromStateReturn = {
  uiFind: string | undefined;
};

const getStyles = createStyle(() => {
  return {
    TraceTimelineViewer: css`
      border-bottom: 1px solid #bbb;

      & .json-markup {
        line-height: 17px;
        font-size: 13px;
        font-family: monospace;
        white-space: pre-wrap;
      }

      & .json-markup-key {
        font-weight: bold;
      }

      & .json-markup-bool {
        color: firebrick;
      }

      & .json-markup-string {
        color: teal;
      }

      & .json-markup-null {
        color: teal;
      }

      & .json-markup-number {
        color: blue;
      }
    `,
  };
});

type TProps = TExtractUiFindFromStateReturn & {
  registerAccessors: (accessors: Accessors) => void;
  findMatchesIDs: Set<string> | TNil;
  scrollToFirstVisibleSpan: () => void;
  traceTimeline: TTraceTimeline;
  trace: Trace;
  updateNextViewRangeTime: (update: ViewRangeTimeUpdate) => void;
  updateViewRangeTime: TUpdateViewRangeTimeFunction;
  viewRange: ViewRange;
  focusSpan: (uiFind: string) => void;
  createLinkToExternalSpan: (traceID: string, spanID: string) => string;

  setSpanNameColumnWidth: (width: number) => void;
  collapseAll: (spans: Span[]) => void;
  collapseOne: (spans: Span[]) => void;
  expandAll: () => void;
  expandOne: (spans: Span[]) => void;

  childrenToggle: (spanID: string) => void;
  clearShouldScrollToFirstUiFindMatch: () => void;
  detailLogItemToggle: (spanID: string, log: Log) => void;
  detailLogsToggle: (spanID: string) => void;
  detailWarningsToggle: (spanID: string) => void;
  detailReferencesToggle: (spanID: string) => void;
  detailProcessToggle: (spanID: string) => void;
  detailTagsToggle: (spanID: string) => void;
  detailToggle: (spanID: string) => void;
  setTrace: (trace: Trace | TNil, uiFind: string | TNil) => void;
  addHoverIndentGuideId: (spanID: string) => void;
  removeHoverIndentGuideId: (spanID: string) => void;
  linksGetter: (span: Span, items: KeyValuePair[], itemIndex: number) => Link[];
};

const NUM_TICKS = 5;

/**
 * `TraceTimelineViewer` now renders the header row because it is sensitive to
 * `props.viewRange.time.cursor`. If `VirtualizedTraceView` renders it, it will
 * re-render the ListView every time the cursor is moved on the trace minimap
 * or `TimelineHeaderRow`.
 */
export default class TraceTimelineViewer extends React.PureComponent<TProps> {
  componentDidMount() {
    mergeShortcuts({
      collapseAll: this.collapseAll,
      expandAll: this.expandAll,
      collapseOne: this.collapseOne,
      expandOne: this.expandOne,
    });
  }

  collapseAll = () => {
    this.props.collapseAll(this.props.trace.spans);
  };

  collapseOne = () => {
    this.props.collapseOne(this.props.trace.spans);
  };

  expandAll = () => {
    this.props.expandAll();
  };

  expandOne = () => {
    this.props.expandOne(this.props.trace.spans);
  };

  render() {
    const {
      setSpanNameColumnWidth,
      updateNextViewRangeTime,
      updateViewRangeTime,
      viewRange,
      createLinkToExternalSpan,
      traceTimeline,
      ...rest
    } = this.props;
    const { trace } = rest;
    const styles = getStyles();

    return (
      <ExternalLinkContext.Provider value={createLinkToExternalSpan}>
        <div className={styles.TraceTimelineViewer}>
          <TimelineHeaderRow
            duration={trace.duration}
            nameColumnWidth={traceTimeline.spanNameColumnWidth}
            numTicks={NUM_TICKS}
            onCollapseAll={this.collapseAll}
            onCollapseOne={this.collapseOne}
            onColummWidthChange={setSpanNameColumnWidth}
            onExpandAll={this.expandAll}
            onExpandOne={this.expandOne}
            viewRangeTime={viewRange.time}
            updateNextViewRangeTime={updateNextViewRangeTime}
            updateViewRangeTime={updateViewRangeTime}
          />
          <VirtualizedTraceView
            {...rest}
            {...traceTimeline}
            setSpanNameColumnWidth={setSpanNameColumnWidth}
            currentViewRangeTime={viewRange.time.current}
          />
        </div>
      </ExternalLinkContext.Provider>
    );
  }
}
