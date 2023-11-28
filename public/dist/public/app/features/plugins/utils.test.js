import { HOME_NAV_ID } from 'app/core/reducers/navModel';
import { buildPluginSectionNav } from './utils';
describe('buildPluginSectionNav', () => {
    const pluginNav = { main: { text: 'Plugin nav' }, node: { text: 'Plugin nav' } };
    const app1 = {
        text: 'App1',
        id: 'plugin-page-app1',
        url: '/a/plugin1',
        children: [
            {
                text: 'page1',
                url: '/a/plugin1/page1',
            },
            {
                text: 'page2',
                url: '/a/plugin1/page2',
            },
            {
                text: 'page3',
                url: '/a/plugin1/page3',
                children: [
                    {
                        text: 'page4',
                        url: '/a/plugin1/page3/page4',
                    },
                ],
            },
        ],
    };
    const appsSection = {
        text: 'apps',
        id: 'apps',
        children: [app1],
    };
    const home = {
        id: HOME_NAV_ID,
        text: 'Home',
    };
    const adminSection = {
        text: 'Admin',
        id: 'admin',
        children: [],
        parentItem: home,
    };
    const standalonePluginPage = {
        id: 'standalone-plugin-page-/a/app2/config',
        text: 'Standalone page',
        parentItem: adminSection,
    };
    adminSection.children = [standalonePluginPage];
    app1.parentItem = appsSection;
    it('Should return return section nav', () => {
        const result = buildPluginSectionNav(appsSection, pluginNav, '/a/plugin1/page1');
        expect(result === null || result === void 0 ? void 0 : result.main.text).toBe('apps');
    });
    it('Should set active page', () => {
        const result = buildPluginSectionNav(appsSection, null, '/a/plugin1/page2');
        expect(result === null || result === void 0 ? void 0 : result.main.children[0].children[1].active).toBe(true);
        expect(result === null || result === void 0 ? void 0 : result.node.text).toBe('page2');
    });
    it('Should only set the most specific match as active (not the parents)', () => {
        const result = buildPluginSectionNav(appsSection, null, '/a/plugin1/page2');
        expect(result === null || result === void 0 ? void 0 : result.main.children[0].children[1].active).toBe(true);
        expect(result === null || result === void 0 ? void 0 : result.main.children[0].active).not.toBe(true); // Parent should not be active
    });
    it('Should set app section to active', () => {
        const result = buildPluginSectionNav(appsSection, null, '/a/plugin1');
        expect(result === null || result === void 0 ? void 0 : result.main.children[0].active).toBe(true);
        expect(result === null || result === void 0 ? void 0 : result.node.text).toBe('App1');
    });
    it('Should handle standalone page', () => {
        const result = buildPluginSectionNav(adminSection, pluginNav, '/a/app2/config');
        expect(result === null || result === void 0 ? void 0 : result.main.text).toBe('Admin');
        expect(result === null || result === void 0 ? void 0 : result.node.text).toBe('Standalone page');
    });
    it('Should set nested active page', () => {
        const result = buildPluginSectionNav(appsSection, null, '/a/plugin1/page3/page4');
        expect(result === null || result === void 0 ? void 0 : result.main.children[0].children[2].children[0].active).toBe(true);
        expect(result === null || result === void 0 ? void 0 : result.node.text).toBe('page4');
    });
});
//# sourceMappingURL=utils.test.js.map