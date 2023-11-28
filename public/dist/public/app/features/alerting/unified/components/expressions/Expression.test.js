import { __awaiter } from "tslib";
import { screen, render } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { times } from 'lodash';
import React from 'react';
import { toDataFrame } from '@grafana/data';
import { ExpressionResult } from './Expression';
describe('TestResult', () => {
    it('should be able to render', () => {
        expect(() => {
            render(React.createElement(ExpressionResult, { series: [] }));
        }).not.toThrow();
    });
    it('should show labels and values', () => {
        const series = [
            toDataFrame({ fields: [{ name: 'temp', values: [0.1234], labels: { label1: 'value1', label2: 'value2' } }] }),
            toDataFrame({ fields: [{ name: 'temp', values: [0.5678], labels: { label1: 'value3', label2: 'value4' } }] }),
        ];
        render(React.createElement(ExpressionResult, { series: series }));
        expect(screen.getByTitle('{label1=value1, label2=value2}')).toBeInTheDocument();
        expect(screen.getByText('0.1234')).toBeInTheDocument();
        expect(screen.getByTitle('{label1=value3, label2=value4}')).toBeInTheDocument();
        expect(screen.getByText('0.5678')).toBeInTheDocument();
    });
    it('should not paginate with less than PAGE_SIZE', () => {
        const series = [
            toDataFrame({
                fields: [
                    {
                        name: 'temp',
                        values: [23, 11, 10],
                    },
                ],
            }),
        ];
        render(React.createElement(ExpressionResult, { series: series }));
        expect(screen.queryByTestId('paginate-expression')).not.toBeInTheDocument();
    });
    it('should paginate with greater than PAGE_SIZE', () => __awaiter(void 0, void 0, void 0, function* () {
        const series = makeSeries(50);
        render(React.createElement(ExpressionResult, { series: series }));
        expect(screen.getByTestId('paginate-expression')).toBeInTheDocument();
        expect(screen.getByText(`1 - 20 of ${50}`)).toBeInTheDocument();
        // click previous page
        yield userEvent.click(screen.getByLabelText('previous-page'));
        expect(screen.getByText(`1 - 20 of ${50}`)).toBeInTheDocument();
        // keep clicking next page, should clamp
        yield userEvent.click(screen.getByLabelText('next-page'));
        expect(screen.getByText(`21 - 40 of ${50}`)).toBeInTheDocument();
        yield userEvent.click(screen.getByLabelText('next-page'));
        expect(screen.getByText(`41 - 50 of ${50}`)).toBeInTheDocument();
        // click one more time, should still be on the last page
        yield userEvent.click(screen.getByLabelText('next-page'));
        expect(screen.getByText(`41 - 50 of ${50}`)).toBeInTheDocument();
    }));
});
function makeSeries(n) {
    return times(n, () => toDataFrame({
        fields: [
            {
                name: 'temp',
                values: [0.1234],
                labels: {
                    label1: 'value1',
                    label2: 'value2',
                },
            },
        ],
    }));
}
//# sourceMappingURL=Expression.test.js.map