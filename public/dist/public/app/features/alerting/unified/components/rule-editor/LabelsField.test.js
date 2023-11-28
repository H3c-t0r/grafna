import { __awaiter } from "tslib";
import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { Provider } from 'react-redux';
import { configureStore } from 'app/store/configureStore';
import LabelsField from './LabelsField';
const labels = [
    { key: 'key1', value: 'value1' },
    { key: 'key2', value: 'value2' },
];
const FormProviderWrapper = ({ children }) => {
    const methods = useForm({ defaultValues: { labels } });
    return React.createElement(FormProvider, Object.assign({}, methods), children);
};
function renderAlertLabels(dataSourceName) {
    const store = configureStore({});
    render(React.createElement(Provider, { store: store }, dataSourceName ? React.createElement(LabelsField, { dataSourceName: dataSourceName }) : React.createElement(LabelsField, null)), { wrapper: FormProviderWrapper });
}
describe('LabelsField with suggestions', () => {
    it('Should display two dropdowns with the existing labels', () => __awaiter(void 0, void 0, void 0, function* () {
        renderAlertLabels('grafana');
        yield waitFor(() => expect(screen.getAllByTestId('alertlabel-key-picker')).toHaveLength(2));
        expect(screen.getByTestId('label-key-0').textContent).toBe('key1');
        expect(screen.getByTestId('label-key-1').textContent).toBe('key2');
        expect(screen.getAllByTestId('alertlabel-value-picker')).toHaveLength(2);
        expect(screen.getByTestId('label-value-0').textContent).toBe('value1');
        expect(screen.getByTestId('label-value-1').textContent).toBe('value2');
    }));
    it('Should delete a key-value combination', () => __awaiter(void 0, void 0, void 0, function* () {
        renderAlertLabels('grafana');
        yield waitFor(() => expect(screen.getAllByTestId('alertlabel-key-picker')).toHaveLength(2));
        expect(screen.getAllByTestId('alertlabel-key-picker')).toHaveLength(2);
        expect(screen.getAllByTestId('alertlabel-value-picker')).toHaveLength(2);
        yield userEvent.click(screen.getByTestId('delete-label-1'));
        expect(screen.getAllByTestId('alertlabel-key-picker')).toHaveLength(1);
        expect(screen.getAllByTestId('alertlabel-value-picker')).toHaveLength(1);
    }));
    it('Should add new key-value dropdowns', () => __awaiter(void 0, void 0, void 0, function* () {
        renderAlertLabels('grafana');
        yield waitFor(() => expect(screen.getByText('Add label')).toBeVisible());
        yield userEvent.click(screen.getByText('Add label'));
        expect(screen.getAllByTestId('alertlabel-key-picker')).toHaveLength(3);
        expect(screen.getByTestId('label-key-0').textContent).toBe('key1');
        expect(screen.getByTestId('label-key-1').textContent).toBe('key2');
        expect(screen.getByTestId('label-key-2').textContent).toBe('Choose key');
        expect(screen.getAllByTestId('alertlabel-value-picker')).toHaveLength(3);
        expect(screen.getByTestId('label-value-0').textContent).toBe('value1');
        expect(screen.getByTestId('label-value-1').textContent).toBe('value2');
        expect(screen.getByTestId('label-value-2').textContent).toBe('Choose value');
    }));
    it('Should be able to write new keys and values using the dropdowns', () => __awaiter(void 0, void 0, void 0, function* () {
        renderAlertLabels('grafana');
        yield waitFor(() => expect(screen.getByText('Add label')).toBeVisible());
        yield userEvent.click(screen.getByText('Add label'));
        const LastKeyDropdown = within(screen.getByTestId('label-key-2'));
        const LastValueDropdown = within(screen.getByTestId('label-value-2'));
        yield userEvent.type(LastKeyDropdown.getByRole('combobox'), 'key3{enter}');
        yield userEvent.type(LastValueDropdown.getByRole('combobox'), 'value3{enter}');
        expect(screen.getByTestId('label-key-2').textContent).toBe('key3');
        expect(screen.getByTestId('label-value-2').textContent).toBe('value3');
    }));
    it('Should be able to write new keys and values using the dropdowns, case sensitive', () => __awaiter(void 0, void 0, void 0, function* () {
        renderAlertLabels('grafana');
        yield waitFor(() => expect(screen.getAllByTestId('alertlabel-key-picker')).toHaveLength(2));
        expect(screen.getByTestId('label-key-0').textContent).toBe('key1');
        expect(screen.getByTestId('label-key-1').textContent).toBe('key2');
        expect(screen.getByTestId('label-value-0').textContent).toBe('value1');
        expect(screen.getByTestId('label-value-1').textContent).toBe('value2');
        const LastKeyDropdown = within(screen.getByTestId('label-key-1'));
        const LastValueDropdown = within(screen.getByTestId('label-value-1'));
        yield userEvent.type(LastKeyDropdown.getByRole('combobox'), 'KEY2{enter}');
        expect(screen.getByTestId('label-key-0').textContent).toBe('key1');
        expect(screen.getByTestId('label-key-1').textContent).toBe('KEY2');
        yield userEvent.type(LastValueDropdown.getByRole('combobox'), 'VALUE2{enter}');
        expect(screen.getByTestId('label-value-0').textContent).toBe('value1');
        expect(screen.getByTestId('label-value-1').textContent).toBe('VALUE2');
    }));
});
describe('LabelsField without suggestions', () => {
    it('Should display two inputs without label suggestions', () => __awaiter(void 0, void 0, void 0, function* () {
        renderAlertLabels();
        yield waitFor(() => expect(screen.getAllByTestId('alertlabel-input-wrapper')).toHaveLength(2));
        expect(screen.queryAllByTestId('alertlabel-key-picker')).toHaveLength(0);
        expect(screen.getByTestId('label-key-0')).toHaveValue('key1');
        expect(screen.getByTestId('label-key-1')).toHaveValue('key2');
        expect(screen.getByTestId('label-value-0')).toHaveValue('value1');
        expect(screen.getByTestId('label-value-1')).toHaveValue('value2');
    }));
});
//# sourceMappingURL=LabelsField.test.js.map