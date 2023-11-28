import { __awaiter } from "tslib";
import { fireEvent, render, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';
import { MultiFilter } from './MultiFilter';
describe('MultiFilters', () => {
    describe('when rendered with two existing multifilters', () => {
        it('should render two filter items', () => __awaiter(void 0, void 0, void 0, function* () {
            const filters = {
                InstanceId: ['a', 'b'],
                InstanceGroup: ['Group1'],
            };
            const onChange = jest.fn();
            render(React.createElement(MultiFilter, { filters: filters, onChange: onChange }));
            const filterItems = screen.getAllByTestId('cloudwatch-multifilter-item');
            expect(filterItems.length).toBe(2);
            expect(within(filterItems[0]).getByDisplayValue('InstanceId')).toBeInTheDocument();
            expect(within(filterItems[0]).getByDisplayValue('a, b')).toBeInTheDocument();
            expect(within(filterItems[1]).getByDisplayValue('InstanceGroup')).toBeInTheDocument();
            expect(within(filterItems[1]).getByDisplayValue('Group1')).toBeInTheDocument();
        }));
    });
    describe('when adding a new filter item', () => {
        it('it should add the new item but not call onChange', () => __awaiter(void 0, void 0, void 0, function* () {
            const filters = {};
            const onChange = jest.fn();
            render(React.createElement(MultiFilter, { filters: filters, onChange: onChange }));
            yield userEvent.click(screen.getByLabelText('Add'));
            expect(screen.getByTestId('cloudwatch-multifilter-item')).toBeInTheDocument();
            expect(onChange).not.toHaveBeenCalled();
        }));
    });
    describe('when adding a new filter item with key', () => {
        it('it should add the new item but not call onChange', () => __awaiter(void 0, void 0, void 0, function* () {
            const filters = {};
            const onChange = jest.fn();
            render(React.createElement(MultiFilter, { filters: filters, onChange: onChange }));
            yield userEvent.click(screen.getByLabelText('Add'));
            const filterItemElement = screen.getByTestId('cloudwatch-multifilter-item');
            expect(filterItemElement).toBeInTheDocument();
            const keyElement = screen.getByTestId('cloudwatch-multifilter-item-key');
            expect(keyElement).toBeInTheDocument();
            yield userEvent.type(keyElement, 'my-key');
            fireEvent.blur(keyElement);
            expect(within(filterItemElement).getByDisplayValue('my-key')).toBeInTheDocument();
            expect(onChange).not.toHaveBeenCalled();
        }));
    });
    describe('when adding a new filter item with key and value', () => {
        it('it should add the new item and trigger onChange', () => __awaiter(void 0, void 0, void 0, function* () {
            const filters = {};
            const onChange = jest.fn();
            render(React.createElement(MultiFilter, { filters: filters, onChange: onChange }));
            const label = yield screen.findByLabelText('Add');
            yield userEvent.click(label);
            const filterItemElement = screen.getByTestId('cloudwatch-multifilter-item');
            expect(filterItemElement).toBeInTheDocument();
            const keyElement = screen.getByTestId('cloudwatch-multifilter-item-key');
            expect(keyElement).toBeInTheDocument();
            yield userEvent.type(keyElement, 'my-key');
            fireEvent.blur(keyElement);
            expect(within(filterItemElement).getByDisplayValue('my-key')).toBeInTheDocument();
            expect(onChange).not.toHaveBeenCalled();
            const valueElement = screen.getByTestId('cloudwatch-multifilter-item-value');
            expect(valueElement).toBeInTheDocument();
            yield userEvent.type(valueElement, 'my-value1,my-value2');
            fireEvent.blur(valueElement);
            expect(within(filterItemElement).getByDisplayValue('my-value1, my-value2')).toBeInTheDocument();
            expect(onChange).toHaveBeenCalledWith({
                'my-key': ['my-value1', 'my-value2'],
            });
        }));
    });
    describe('when editing an existing filter item key', () => {
        it('it should change the key and call onChange', () => __awaiter(void 0, void 0, void 0, function* () {
            const filters = { 'my-key': ['my-value'] };
            const onChange = jest.fn();
            render(React.createElement(MultiFilter, { filters: filters, onChange: onChange }));
            const filterItemElement = screen.getByTestId('cloudwatch-multifilter-item');
            expect(filterItemElement).toBeInTheDocument();
            expect(within(filterItemElement).getByDisplayValue('my-key')).toBeInTheDocument();
            expect(within(filterItemElement).getByDisplayValue('my-value')).toBeInTheDocument();
            const keyElement = screen.getByTestId('cloudwatch-multifilter-item-key');
            expect(keyElement).toBeInTheDocument();
            yield userEvent.type(keyElement, '2');
            fireEvent.blur(keyElement);
            expect(within(filterItemElement).getByDisplayValue('my-key2')).toBeInTheDocument();
            expect(onChange).toHaveBeenCalledWith({
                'my-key2': ['my-value'],
            });
        }));
    });
});
//# sourceMappingURL=MultiFilter.test.js.map