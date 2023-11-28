import { __awaiter } from "tslib";
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';
import { getDefaultTimeRange, LoadingState } from '@grafana/data';
import { setDataSourceSrv } from '@grafana/runtime';
import { mockDataSource, MockDataSourceSrv } from 'app/features/alerting/unified/mocks';
import { getDashboardSrv } from 'app/features/dashboard/services/DashboardSrv';
import { createDashboardModelFixture, createPanelSaveModel, } from '../../../features/dashboard/state/__fixtures__/dashboardFixtures';
import { DashboardQueryEditor } from './DashboardQueryEditor';
import { SHARED_DASHBOARD_QUERY } from './types';
jest.mock('app/core/config', () => (Object.assign(Object.assign({}, jest.requireActual('app/core/config')), { panels: {
        timeseries: {
            info: {
                logos: {
                    small: '',
                },
            },
        },
    } })));
setDataSourceSrv(new MockDataSourceSrv({
    test: mockDataSource({ isDefault: true }),
}));
describe('DashboardQueryEditor', () => {
    const mockOnChange = jest.fn();
    const mockOnRunQueries = jest.fn();
    const mockPanelData = {
        state: LoadingState.Done,
        series: [],
        timeRange: getDefaultTimeRange(),
    };
    const mockQueries = [{ refId: 'A' }];
    let mockDashboard;
    beforeEach(() => {
        mockDashboard = createDashboardModelFixture({
            panels: [
                createPanelSaveModel({
                    targets: [],
                    type: 'timeseries',
                    id: 1,
                    title: 'My first panel',
                }),
                createPanelSaveModel({
                    targets: [],
                    id: 2,
                    type: 'timeseries',
                    title: 'Another panel',
                }),
                createPanelSaveModel({
                    datasource: {
                        uid: SHARED_DASHBOARD_QUERY,
                    },
                    targets: [],
                    id: 3,
                    type: 'timeseries',
                    title: 'A dashboard query panel',
                }),
            ],
        });
        jest.spyOn(getDashboardSrv(), 'getCurrent').mockImplementation(() => mockDashboard);
    });
    it('does not show a panel with the SHARED_DASHBOARD_QUERY datasource as an option in the dropdown', () => __awaiter(void 0, void 0, void 0, function* () {
        render(React.createElement(DashboardQueryEditor, { queries: mockQueries, panelData: mockPanelData, onChange: mockOnChange, onRunQueries: mockOnRunQueries }));
        const select = screen.getByText('Choose panel');
        yield userEvent.click(select);
        const myFirstPanel = yield screen.findByText('My first panel');
        expect(myFirstPanel).toBeInTheDocument();
        const anotherPanel = yield screen.findByText('Another panel');
        expect(anotherPanel).toBeInTheDocument();
        expect(screen.queryByText('A dashboard query panel')).not.toBeInTheDocument();
    }));
    it('does not show the current panelInEdit as an option in the dropdown', () => __awaiter(void 0, void 0, void 0, function* () {
        mockDashboard.initEditPanel(mockDashboard.panels[0]);
        render(React.createElement(DashboardQueryEditor, { queries: mockQueries, panelData: mockPanelData, onChange: mockOnChange, onRunQueries: mockOnRunQueries }));
        const select = screen.getByText('Choose panel');
        yield userEvent.click(select);
        expect(screen.queryByText('My first panel')).not.toBeInTheDocument();
        const anotherPanel = yield screen.findByText('Another panel');
        expect(anotherPanel).toBeInTheDocument();
        expect(screen.queryByText('A dashboard query panel')).not.toBeInTheDocument();
    }));
});
//# sourceMappingURL=DashboardQueryEditor.test.js.map