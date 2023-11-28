import { __awaiter } from "tslib";
import { wellFormedPanelModel } from 'test/fixtures/panelModel.fixture';
import { mockFolderDTO } from 'app/features/browse-dashboards/fixtures/folder.fixture';
import { dashlistMigrationHandler } from './migrations';
const getMock = jest.fn();
jest.mock('@grafana/runtime', () => (Object.assign(Object.assign({}, jest.requireActual('@grafana/runtime')), { getBackendSrv: () => ({
        get: getMock,
    }) })));
describe('dashlist migrations', () => {
    it('migrates angular panel model to react model', () => __awaiter(void 0, void 0, void 0, function* () {
        const basePanelModel = wellFormedPanelModel({});
        basePanelModel.pluginVersion = '5.1';
        const angularPanel = Object.assign(Object.assign({}, basePanelModel), { 
            // pluginVersion: '5.1',
            starred: true, recent: true, search: true, headings: true, limit: 7, query: 'hello, query' });
        const newOptions = yield dashlistMigrationHandler(angularPanel);
        expect(newOptions).toEqual({
            showStarred: true,
            showRecentlyViewed: true,
            showSearch: true,
            showHeadings: true,
            maxItems: 7,
            query: 'hello, query',
            includeVars: undefined,
            keepTime: undefined,
        });
        expect(angularPanel).toStrictEqual(basePanelModel);
    }));
    it('migrates folder id to folder UID', () => __awaiter(void 0, void 0, void 0, function* () {
        const folderDTO = mockFolderDTO(1, {
            id: 77,
            uid: 'abc-124',
        });
        getMock.mockResolvedValue(folderDTO);
        const basePanelOptions = {
            showStarred: true,
            showRecentlyViewed: true,
            showSearch: true,
            showHeadings: true,
            maxItems: 7,
            query: 'hello, query',
            includeVars: false,
            keepTime: false,
            tags: [],
        };
        const panelModel = wellFormedPanelModel(Object.assign(Object.assign({}, basePanelOptions), { folderId: 77 }));
        const newOptions = yield dashlistMigrationHandler(panelModel);
        expect(newOptions).toStrictEqual(Object.assign(Object.assign({}, basePanelOptions), { folderUID: 'abc-124' }));
    }));
    it("doesn't fail if the api request fails", () => __awaiter(void 0, void 0, void 0, function* () {
        const spyConsoleWarn = jest.spyOn(console, 'warn').mockImplementation();
        getMock.mockRejectedValue({
            status: 403,
            statusText: 'Forbidden',
            data: {
                accessErrorId: 'ACE0577385389',
                message: "You'll need additional permissions to perform this action. Permissions needed: folders:read",
                title: 'Access denied',
            },
            config: {
                showErrorAlert: false,
                method: 'GET',
                url: 'api/folders/id/0',
                retry: 0,
                headers: {
                    'X-Grafana-Org-Id': 1,
                },
                hideFromInspector: true,
            },
        });
        const basePanelOptions = {
            showStarred: true,
            showRecentlyViewed: true,
            showSearch: true,
            showHeadings: true,
            maxItems: 7,
            query: 'hello, query',
            includeVars: false,
            keepTime: false,
            tags: [],
            folderId: 77,
        };
        const panelModel = wellFormedPanelModel(basePanelOptions);
        // We expect it to not reject
        const newOptions = yield dashlistMigrationHandler(panelModel);
        expect(newOptions).toStrictEqual(basePanelOptions);
        expect(spyConsoleWarn).toHaveBeenCalledTimes(1);
    }));
});
//# sourceMappingURL=migrations.test.js.map