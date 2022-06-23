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

import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';

import { createTheme } from '@grafana/data';

import SpanDetail from './SpanDetail';
import DetailState from './SpanDetail/DetailState';
import { UnthemedSpanDetailRow } from './SpanDetailRow';
import SpanTreeOffset from './SpanTreeOffset';

const testSpan = {
  spanID: 'testSpanID',
  traceID: 'testTraceID',
  depth: 3,
  process: {
    serviceName: 'some-service',
    tags: [{ key: 'tag-key', value: 'tag-value' }],
  },
};
const setup = (propOverrides) => {
  const props = {
    color: 'some-color',
    columnDivision: 0.5,
    detailState: new DetailState(),
    onDetailToggled: jest.fn(),
    linksGetter: jest.fn(),
    isFilteredOut: false,
    logItemToggle: jest.fn(),
    logsToggle: jest.fn(),
    processToggle: jest.fn(),
    createFocusSpanLink: jest.fn(),
    hoverIndentGuideIds: new Map(),
    span: testSpan,
    tagsToggle: jest.fn(),
    traceStartTime: 1000,
    theme: createTheme(),
    ...propOverrides,
  };
  return render(<UnthemedSpanDetailRow {...props} />);
};

describe('SpanDetailRow tests', () => {
  it('renders without exploding', () => {
    expect(() => setup()).not.toThrow();
  });

  it('calls toggle on click', async () => {
    const mockToggle = jest.fn();
    setup({ onDetailToggled: mockToggle });
    expect(mockToggle).not.toHaveBeenCalled();

    const detailRow = screen.getByTestId('detail-row-expanded-accent');
    await userEvent.click(detailRow);
    expect(mockToggle).toHaveBeenCalled();
  });

  it('renders the span tree offset', () => {
    setup();

    expect(screen.getByTestId('SpanTreeOffset--indentGuide')).toBeInTheDocument();
  });

  it('renders the SpanDetail', () => {
    setup();

    expect(screen.getByTestId('span-detail-component')).toBeInTheDocument();
  });

  // it('adds span when calling linksGetter', () => {
  //   const spanDetail = wrapper.find(SpanDetail);
  //   const linksGetter = spanDetail.prop('linksGetter');
  //   const tags = [{ key: 'myKey', value: 'myValue' }];
  //   const linksGetterResponse = {};
  //   props.linksGetter.mockReturnValueOnce(linksGetterResponse);
  //   const result = linksGetter(tags, 0);
  //   expect(result).toBe(linksGetterResponse);
  //   expect(props.linksGetter).toHaveBeenCalledTimes(1);
  //   expect(props.linksGetter).toHaveBeenCalledWith(props.span, tags, 0);
  // });
});
