import { __assign } from "tslib";
import React from 'react';
import { FieldType } from '@grafana/data';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { InspectDataTab } from './InspectDataTab';
var createProps = function (propsOverride) {
    var defaultProps = {
        isLoading: false,
        options: {
            withTransforms: false,
            withFieldConfig: false,
        },
        data: [
            {
                name: 'First data frame',
                fields: [
                    { name: 'time', type: FieldType.time, values: [100, 200, 300] },
                    { name: 'name', type: FieldType.string, values: ['uniqueA', 'b', 'c'] },
                    { name: 'value', type: FieldType.number, values: [1, 2, 3] },
                ],
                length: 3,
            },
            {
                name: 'Second data frame',
                fields: [
                    { name: 'time', type: FieldType.time, values: [400, 500, 600] },
                    { name: 'name', type: FieldType.string, values: ['d', 'e', 'g'] },
                    { name: 'value', type: FieldType.number, values: [4, 5, 6] },
                ],
                length: 3,
            },
        ],
    };
    return Object.assign(defaultProps, propsOverride);
};
describe('InspectDataTab', function () {
    describe('when panel is not passed as prop (Explore)', function () {
        it('should render InspectDataTab', function () {
            render(React.createElement(InspectDataTab, __assign({}, createProps())));
            expect(screen.getByLabelText(/Panel inspector Data content/i)).toBeInTheDocument();
        });
        it('should render Data Option row', function () {
            render(React.createElement(InspectDataTab, __assign({}, createProps())));
            expect(screen.getByText(/Data options/i)).toBeInTheDocument();
        });
        it('should show available options', function () {
            render(React.createElement(InspectDataTab, __assign({}, createProps())));
            var dataOptions = screen.getByText(/Data options/i);
            userEvent.click(dataOptions);
            expect(screen.getByText(/Show data frame/i)).toBeInTheDocument();
            expect(screen.getByText(/Download for Excel/i)).toBeInTheDocument();
        });
        it('should show available dataFrame options', function () {
            render(React.createElement(InspectDataTab, __assign({}, createProps())));
            var dataOptions = screen.getByText(/Data options/i);
            userEvent.click(dataOptions);
            var dataFrameInput = screen.getByRole('textbox', { name: /Select dataframe/i });
            userEvent.click(dataFrameInput);
            expect(screen.getByText(/Second data frame/i)).toBeInTheDocument();
        });
        it('should show download logs button if logs data', function () {
            var dataWithLogs = [
                {
                    name: 'Data frame with logs',
                    fields: [
                        { name: 'time', type: FieldType.time, values: [100, 200, 300] },
                        { name: 'name', type: FieldType.string, values: ['uniqueA', 'b', 'c'] },
                        { name: 'value', type: FieldType.number, values: [1, 2, 3] },
                    ],
                    length: 3,
                    meta: {
                        preferredVisualisationType: 'logs',
                    },
                },
            ];
            render(React.createElement(InspectDataTab, __assign({}, createProps({ data: dataWithLogs }))));
            expect(screen.getByText(/Download logs/i)).toBeInTheDocument();
        });
        it('should not show download logs button if no logs data', function () {
            render(React.createElement(InspectDataTab, __assign({}, createProps())));
            expect(screen.queryByText(/Download logs/i)).not.toBeInTheDocument();
        });
        it('should show download traces button if traces data', function () {
            var dataWithtraces = [
                {
                    name: 'Data frame with traces',
                    fields: [
                        { name: 'traceID', values: ['3fa414edcef6ad90', '3fa414edcef6ad90'] },
                        { name: 'spanID', values: ['3fa414edcef6ad90', '0f5c1808567e4403'] },
                        { name: 'parentSpanID', values: [undefined, '3fa414edcef6ad90'] },
                        { name: 'operationName', values: ['HTTP GET - api_traces_traceid', '/tempopb.Querier/FindTraceByID'] },
                        { name: 'serviceName', values: ['tempo-querier', 'tempo-querier'] },
                        {
                            name: 'serviceTags',
                            values: [
                                [
                                    { key: 'cluster', type: 'string', value: 'ops-tools1' },
                                    { key: 'container', type: 'string', value: 'tempo-query' },
                                ],
                                [
                                    { key: 'cluster', type: 'string', value: 'ops-tools1' },
                                    { key: 'container', type: 'string', value: 'tempo-query' },
                                ],
                            ],
                        },
                        { name: 'startTime', values: [1605873894680.409, 1605873894680.587] },
                        { name: 'duration', values: [1049.141, 1.847] },
                        { name: 'logs', values: [[], []] },
                        {
                            name: 'tags',
                            values: [
                                [
                                    { key: 'sampler.type', type: 'string', value: 'probabilistic' },
                                    { key: 'sampler.param', type: 'float64', value: 1 },
                                ],
                                [
                                    { key: 'component', type: 'string', value: 'gRPC' },
                                    { key: 'span.kind', type: 'string', value: 'client' },
                                ],
                            ],
                        },
                        { name: 'warnings', values: [undefined, undefined] },
                        { name: 'stackTraces', values: [undefined, undefined] },
                    ],
                    length: 2,
                    meta: {
                        preferredVisualisationType: 'trace',
                        custom: {
                            traceFormat: 'jaeger',
                        },
                    },
                },
            ];
            render(React.createElement(InspectDataTab, __assign({}, createProps({ data: dataWithtraces }))));
            expect(screen.getByText(/Download traces/i)).toBeInTheDocument();
        });
        it('should not show download traces button if no traces data', function () {
            render(React.createElement(InspectDataTab, __assign({}, createProps())));
            expect(screen.queryByText(/Download traces/i)).not.toBeInTheDocument();
        });
    });
});
//# sourceMappingURL=InspectDataTab.test.js.map