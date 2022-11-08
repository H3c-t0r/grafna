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

import fuzzy from 'fuzzy';
import { createSelector } from 'reselect';

import { TraceSpan, TraceSpanData, TraceSpanReference } from '../types/trace';

import { getProcessServiceName } from './process';

export const getSpanId = (span: TraceSpanData) => span.spanID;
export const getSpanName = (span: TraceSpanData) => span.operationName;
export const getSpanDuration = (span: TraceSpanData) => span.duration;
export const getSpanTimestamp = (span: TraceSpanData) => span.startTime;
export const getSpanProcessId = (span: TraceSpanData) => span.processID;
export const getSpanReferences = (span: TraceSpanData) => span.references || [];
export const getSpanReferenceByType = createSelector(
  createSelector(({ span }: { span: TraceSpanData }) => span, getSpanReferences),
  ({ type }: { type: string }) => type,
  (references, type) => references.find((ref: TraceSpanReference) => ref.refType === type)
);
export const getSpanParentId = createSelector(
  (span: TraceSpanData) => getSpanReferenceByType({ span, type: 'CHILD_OF' }),
  (childOfRef) => (childOfRef ? childOfRef.spanID : null)
);

export const getSpanProcess = (span: TraceSpan | TraceSpanData) => {
  if ('process' in span) {
    return span.process;
  } else {
    throw new Error(
      `
      you must hydrate the spans with the processes, perhaps
      using hydrateSpansWithProcesses(), before accessing a span's process
    `
    );
  }
};

export const getSpanServiceName = createSelector(getSpanProcess, getProcessServiceName);

export const filterSpansForTimestamps = createSelector(
  ({ spans }: { spans: TraceSpanData[] }) => spans,
  ({ leftBound }: { leftBound: number }) => leftBound,
  ({ rightBound }: { rightBound: number }) => rightBound,
  (spans, leftBound, rightBound) =>
    spans.filter((span) => getSpanTimestamp(span) >= leftBound && getSpanTimestamp(span) <= rightBound)
);

export const filterSpansForText = createSelector(
  ({ spans }: { spans: TraceSpanData[] }) => spans,
  ({ text }: { text: string }) => text,
  (spans, text) =>
    fuzzy
      .filter(text, spans, {
        extract: (span) => `${getSpanServiceName(span)} ${getSpanName(span)}`,
      })
      .map(({ original }) => original)
);

const getTextFilteredSpansAsMap = createSelector(filterSpansForText, (matchingSpans) =>
  matchingSpans.reduce(
    (obj, span) => ({
      ...obj,
      [getSpanId(span)]: span,
    }),
    {}
  )
);

// TODO: delete this function as it is not used?
export const highlightSpansForTextFilter = createSelector(
  ({ spans }: { spans: TraceSpanData[] }) => spans,
  getTextFilteredSpansAsMap,
  (spans, textFilteredSpansMap: { [key: string]: TraceSpanData }) =>
    spans.map((span: TraceSpanData) => ({
      ...span,
      muted: !textFilteredSpansMap[getSpanId(span)],
    }))
);
