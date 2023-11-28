import { __awaiter } from "tslib";
import { render, screen, waitFor } from '@testing-library/react';
import React from 'react';
import { Provider } from 'react-redux';
import { MemoryRouter } from 'react-router-dom';
import { configureStore } from 'app/store/configureStore';
import { stubWithLabels } from '../inventory/__mocks__/Inventory.service';
import EditInstancePage from './EditInstance';
import { fromPayload } from './EditInstance.utils';
jest.mock('app/percona/inventory/Inventory.service');
const renderWithDefaults = () => render(React.createElement(MemoryRouter, { initialEntries: ['/edit-instance/service_id'] },
    React.createElement(Provider, { store: configureStore() },
        React.createElement(EditInstancePage, null))));
describe('EditInstance::', () => {
    it('prefills current values', () => __awaiter(void 0, void 0, void 0, function* () {
        renderWithDefaults();
        yield waitFor(() => expect(screen.queryByLabelText('Environment')).toHaveValue(stubWithLabels.environment));
        yield waitFor(() => expect(screen.queryByLabelText('Cluster')).toHaveValue(stubWithLabels.cluster));
        yield waitFor(() => expect(screen.queryByLabelText('Replication set')).toHaveValue(stubWithLabels.replication_set));
        yield waitFor(() => expect(screen.queryByLabelText('Custom labels')).toHaveValue(fromPayload(stubWithLabels.custom_labels)));
    }));
});
//# sourceMappingURL=EditInstance.test.js.map