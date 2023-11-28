import { __awaiter } from "tslib";
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import React from 'react';
import { Provider } from 'react-redux';
import { byRole, byTestId, byText } from 'testing-library-selector';
import { logInfo } from '@grafana/runtime';
import { contextSrv } from 'app/core/services/context_srv';
import { configureStore } from 'app/store/configureStore';
import { AccessControlAction } from 'app/types';
import { LogMessages } from '../../Analytics';
import { useHasRuler } from '../../hooks/useHasRuler';
import { mockExportApi, mockFolderApi, setupMswServer } from '../../mockApi';
import { grantUserPermissions, mockCombinedRule, mockDataSource, mockFolder, mockGrafanaRulerRule } from '../../mocks';
import { RulesGroup } from './RulesGroup';
jest.mock('../../hooks/useHasRuler');
jest.mock('@grafana/runtime', () => {
    const original = jest.requireActual('@grafana/runtime');
    return Object.assign(Object.assign({}, original), { logInfo: jest.fn() });
});
jest.mock('react-virtualized-auto-sizer', () => {
    return ({ children }) => children({ height: 600, width: 1 });
});
jest.mock('@grafana/ui', () => (Object.assign(Object.assign({}, jest.requireActual('@grafana/ui')), { CodeEditor: ({ value }) => React.createElement("textarea", { "data-testid": "code-editor", value: value, readOnly: true }) })));
const mocks = {
    useHasRuler: jest.mocked(useHasRuler),
};
function mockUseHasRuler(hasRuler, rulerRulesLoaded) {
    mocks.useHasRuler.mockReturnValue({
        hasRuler: () => hasRuler,
        rulerRulesLoaded: () => rulerRulesLoaded,
    });
}
beforeEach(() => {
    mocks.useHasRuler.mockReset();
    // FIXME: scope down
    grantUserPermissions(Object.values(AccessControlAction));
});
const ui = {
    editGroupButton: byTestId('edit-group'),
    deleteGroupButton: byTestId('delete-group'),
    exportGroupButton: byRole('button', { name: 'Export rule group' }),
    confirmDeleteModal: {
        header: byText('Delete group'),
        confirmButton: byText('Delete'),
    },
    moreActionsButton: byRole('button', { name: 'More' }),
    export: {
        dialog: byRole('dialog', { name: /Drawer title Export .* rules/ }),
        jsonTab: byRole('tab', { name: /JSON/ }),
        yamlTab: byRole('tab', { name: /YAML/ }),
        editor: byTestId('code-editor'),
        copyCodeButton: byRole('button', { name: 'Copy code' }),
        downloadButton: byRole('button', { name: 'Download' }),
    },
    loadingSpinner: byTestId('spinner'),
};
const server = setupMswServer();
afterEach(() => {
    server.resetHandlers();
});
describe('Rules group tests', () => {
    const store = configureStore();
    function renderRulesGroup(namespace, group) {
        return render(React.createElement(Provider, { store: store },
            React.createElement(RulesGroup, { group: group, namespace: namespace, expandAll: false, viewMode: 'grouped' })));
    }
    describe('Grafana rules', () => {
        const group = {
            name: 'TestGroup',
            rules: [
                mockCombinedRule({
                    rulerRule: mockGrafanaRulerRule({
                        namespace_uid: 'cpu-usage',
                    }),
                }),
            ],
            totals: {},
        };
        const namespace = {
            name: 'TestNamespace',
            rulesSource: 'grafana',
            groups: [group],
        };
        it('Should hide delete and edit group buttons', () => {
            // Act
            mockUseHasRuler(true, true);
            renderRulesGroup(namespace, group);
            // Assert
            expect(ui.deleteGroupButton.query()).not.toBeInTheDocument();
            expect(ui.editGroupButton.query()).not.toBeInTheDocument();
        });
        it('Should allow exporting rules group', () => __awaiter(void 0, void 0, void 0, function* () {
            // Arrange
            mockUseHasRuler(true, true);
            mockFolderApi(server).folder('cpu-usage', mockFolder({ uid: 'cpu-usage' }));
            mockExportApi(server).exportRulesGroup('cpu-usage', 'TestGroup', {
                yaml: 'Yaml Export Content',
                json: 'Json Export Content',
            });
            const user = userEvent.setup();
            // Act
            renderRulesGroup(namespace, group);
            yield user.click(yield ui.exportGroupButton.find());
            // Assert
            const drawer = yield ui.export.dialog.find();
            expect(ui.export.yamlTab.get(drawer)).toHaveAttribute('aria-selected', 'true');
            yield waitFor(() => {
                expect(ui.export.editor.get(drawer)).toHaveTextContent('Yaml Export Content');
            });
            yield user.click(ui.export.jsonTab.get(drawer));
            yield waitFor(() => {
                expect(ui.export.editor.get(drawer)).toHaveTextContent('Json Export Content');
            });
            expect(ui.export.copyCodeButton.get(drawer)).toBeInTheDocument();
            expect(ui.export.downloadButton.get(drawer)).toBeInTheDocument();
        }));
    });
    describe('Cloud rules', () => {
        beforeEach(() => {
            contextSrv.isEditor = true;
        });
        const group = {
            name: 'TestGroup',
            rules: [mockCombinedRule()],
            totals: {},
        };
        const namespace = {
            name: 'TestNamespace',
            rulesSource: mockDataSource(),
            groups: [group],
        };
        it('When ruler enabled should display delete and edit group buttons', () => {
            // Arrange
            mockUseHasRuler(true, true);
            // Act
            renderRulesGroup(namespace, group);
            // Assert
            expect(mocks.useHasRuler).toHaveBeenCalled();
            expect(ui.deleteGroupButton.get()).toBeInTheDocument();
            expect(ui.editGroupButton.get()).toBeInTheDocument();
        });
        it('When ruler disabled should hide delete and edit group buttons', () => {
            // Arrange
            mockUseHasRuler(false, false);
            // Act
            renderRulesGroup(namespace, group);
            // Assert
            expect(mocks.useHasRuler).toHaveBeenCalled();
            expect(ui.deleteGroupButton.query()).not.toBeInTheDocument();
            expect(ui.editGroupButton.query()).not.toBeInTheDocument();
        });
        it('Delete button click should display confirmation modal', () => __awaiter(void 0, void 0, void 0, function* () {
            // Arrange
            mockUseHasRuler(true, true);
            // Act
            renderRulesGroup(namespace, group);
            yield userEvent.click(ui.deleteGroupButton.get());
            // Assert
            expect(ui.confirmDeleteModal.header.get()).toBeInTheDocument();
            expect(ui.confirmDeleteModal.confirmButton.get()).toBeInTheDocument();
        }));
    });
    describe('Analytics', () => {
        beforeEach(() => {
            contextSrv.isEditor = true;
        });
        const group = {
            name: 'TestGroup',
            rules: [mockCombinedRule()],
            totals: {},
        };
        const namespace = {
            name: 'TestNamespace',
            rulesSource: mockDataSource(),
            groups: [group],
        };
        it('Should log info when closing the edit group rule modal without saving', () => __awaiter(void 0, void 0, void 0, function* () {
            mockUseHasRuler(true, true);
            renderRulesGroup(namespace, group);
            yield userEvent.click(ui.editGroupButton.get());
            expect(screen.getByText('Cancel')).toBeInTheDocument();
            yield userEvent.click(screen.getByText('Cancel'));
            expect(screen.queryByText('Cancel')).not.toBeInTheDocument();
            expect(logInfo).toHaveBeenCalledWith(LogMessages.leavingRuleGroupEdit);
        }));
    });
});
//# sourceMappingURL=RulesGroup.test.js.map