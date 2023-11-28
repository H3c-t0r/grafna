import { __awaiter } from "tslib";
import { act, waitFor } from '@testing-library/react';
import { renderHook } from '@testing-library/react-hooks';
import { createMemoryHistory } from 'history';
import { stringify } from 'querystring';
import React from 'react';
import { TestProvider } from 'test/helpers/TestProvider';
import { getGrafanaContextMock } from 'test/mocks/getGrafanaContextMock';
import { HistoryWrapper, setDataSourceSrv } from '@grafana/runtime';
import { setLastUsedDatasourceUID } from 'app/core/utils/explore';
import { MIXED_DATASOURCE_NAME } from 'app/plugins/datasource/mixed/MixedDataSource';
import { configureStore } from 'app/store/configureStore';
import { makeDatasourceSetup } from '../../spec/helper/setup';
import { splitClose, splitOpen } from '../../state/main';
import { useStateSync } from './';
const fetch = jest.fn().mockResolvedValue({ correlations: [] });
jest.mock('@grafana/runtime', () => (Object.assign(Object.assign({}, jest.requireActual('@grafana/runtime')), { getBackendSrv: () => ({ fetch }) })));
jest.mock('rxjs', () => (Object.assign(Object.assign({}, jest.requireActual('rxjs')), { lastValueFrom: () => new Promise((resolve, reject) => {
        resolve({ data: { correlations: [] } });
    }) })));
function defaultDsGetter(datasources) {
    return (datasource) => {
        var _a, _b;
        let ds;
        if (!datasource) {
            ds = (_a = datasources[0]) === null || _a === void 0 ? void 0 : _a.api;
        }
        else {
            ds = (_b = datasources.find((ds) => typeof datasource === 'string'
                ? ds.api.name === datasource || ds.api.uid === datasource
                : ds.api.uid === (datasource === null || datasource === void 0 ? void 0 : datasource.uid))) === null || _b === void 0 ? void 0 : _b.api;
        }
        if (ds) {
            return Promise.resolve(ds);
        }
        return Promise.reject();
    };
}
function setup({ queryParams = {}, datasourceGetter = defaultDsGetter }) {
    const history = createMemoryHistory({
        initialEntries: [{ pathname: '/explore', search: stringify(queryParams) }],
    });
    const location = new HistoryWrapper(history);
    const datasources = [
        makeDatasourceSetup({ name: 'loki', uid: 'loki-uid' }),
        makeDatasourceSetup({ name: 'elastic', uid: 'elastic-uid' }),
        makeDatasourceSetup({ name: MIXED_DATASOURCE_NAME, uid: MIXED_DATASOURCE_NAME, id: 999 }),
    ];
    setDataSourceSrv({
        get: datasourceGetter(datasources),
        getInstanceSettings: jest.fn(),
        getList: jest.fn(),
        reload: jest.fn(),
    });
    const store = configureStore({
        user: {
            orgId: 1,
            fiscalYearStartMonth: 0,
            isUpdating: false,
            orgs: [],
            orgsAreLoading: false,
            sessions: [],
            sessionsAreLoading: false,
            teams: [],
            teamsAreLoading: false,
            timeZone: 'utc',
            user: null,
            weekStart: 'monday',
        },
    });
    const context = getGrafanaContextMock();
    const wrapper = ({ children }) => (React.createElement(TestProvider, { grafanaContext: Object.assign(Object.assign({}, context), { location }), store: store }, children));
    return Object.assign(Object.assign({}, renderHook(({ params }) => useStateSync(params), {
        wrapper,
        initialProps: {
            children: null,
            params: queryParams,
        },
    })), { location,
        store });
}
describe('useStateSync', () => {
    it('does not push a new entry to history on first render', () => __awaiter(void 0, void 0, void 0, function* () {
        const { location, waitForNextUpdate } = setup({});
        const initialHistoryLength = location.getHistory().length;
        yield waitForNextUpdate();
        expect(location.getHistory().length).toBe(initialHistoryLength);
        const search = location.getSearchObject();
        expect(search.panes).toBeDefined();
    }));
    it('correctly inits an explore pane for each key in the panes search object', () => __awaiter(void 0, void 0, void 0, function* () {
        const { location, waitForNextUpdate, store } = setup({
            queryParams: {
                panes: JSON.stringify({
                    one: { datasource: 'loki-uid', queries: [{ datasource: { name: 'loki', uid: 'loki-uid' }, refId: '1+2' }] },
                    two: { datasource: 'elastic-uid', queries: [{ datasource: { name: 'elastic', uid: 'elastic-uid' } }] },
                }),
                schemaVersion: 1,
            },
            datasourceGetter: (datasources) => {
                return (datasource) => {
                    var _a, _b;
                    let ds;
                    if (!datasource) {
                        ds = (_a = datasources[0]) === null || _a === void 0 ? void 0 : _a.api;
                    }
                    else {
                        ds = (_b = datasources.find((ds) => typeof datasource === 'string'
                            ? ds.api.name === datasource || ds.api.uid === datasource
                            : ds.api.uid === (datasource === null || datasource === void 0 ? void 0 : datasource.uid))) === null || _b === void 0 ? void 0 : _b.api;
                    }
                    return new Promise((resolve, reject) => {
                        if (ds) {
                            if (typeof datasource === 'string' && datasource === 'loki-uid') {
                                setTimeout(() => Promise.resolve(ds));
                            }
                            else {
                                resolve(ds);
                            }
                        }
                        reject();
                    });
                };
            },
        });
        const initialHistoryLength = location.getHistory().length;
        yield waitForNextUpdate();
        expect(location.getHistory().length).toBe(initialHistoryLength);
        const panes = location.getSearch().get('panes');
        expect(panes).not.toBeNull();
        if (panes) {
            // check if the URL is properly encoded when finishing rendering the hook. (this would be '1 2' otherwise)
            expect(JSON.parse(panes)['one'].queries[0].refId).toBe('1+2');
            // we expect panes in the state to be in the same order as the ones in the URL
            expect(Object.keys(store.getState().explore.panes)).toStrictEqual(Object.keys(JSON.parse(panes)));
        }
    }));
    it('inits with a default query from the root level datasource when there are no valid queries in the URL', () => __awaiter(void 0, void 0, void 0, function* () {
        var _a, _b;
        const { location, waitForNextUpdate, store } = setup({
            queryParams: {
                panes: JSON.stringify({
                    one: { datasource: 'loki-uid', queries: [{ datasource: { name: 'UNKNOWN', uid: 'UNKNOWN-DS' } }] },
                }),
                schemaVersion: 1,
            },
        });
        const initialHistoryLength = location.getHistory().length;
        yield waitForNextUpdate();
        expect(location.getHistory().length).toBe(initialHistoryLength);
        const search = location.getSearchObject();
        expect(search.panes).toBeDefined();
        const queries = (_a = store.getState().explore.panes['one']) === null || _a === void 0 ? void 0 : _a.queries;
        expect(queries).toHaveLength(1);
        expect((_b = queries === null || queries === void 0 ? void 0 : queries[0].datasource) === null || _b === void 0 ? void 0 : _b.uid).toBe('loki-uid');
    }));
    it('inits with the last used datasource from localStorage', () => __awaiter(void 0, void 0, void 0, function* () {
        var _c, _d;
        setLastUsedDatasourceUID(1, 'elastic-uid');
        const { waitForNextUpdate, store } = setup({
            queryParams: {},
        });
        yield waitForNextUpdate();
        expect((_d = (_c = Object.values(store.getState().explore.panes)[0]) === null || _c === void 0 ? void 0 : _c.datasourceInstance) === null || _d === void 0 ? void 0 : _d.uid).toBe('elastic-uid');
    }));
    it('inits with the default datasource if the last used in localStorage does not exits', () => __awaiter(void 0, void 0, void 0, function* () {
        var _e, _f;
        setLastUsedDatasourceUID(1, 'unknown-ds-uid');
        const { waitForNextUpdate, store } = setup({
            queryParams: {},
        });
        yield waitForNextUpdate();
        expect((_f = (_e = Object.values(store.getState().explore.panes)[0]) === null || _e === void 0 ? void 0 : _e.datasourceInstance) === null || _f === void 0 ? void 0 : _f.uid).toBe('loki-uid');
    }));
    it('updates the state with correct queries from URL', () => __awaiter(void 0, void 0, void 0, function* () {
        var _g, _h, _j;
        const { waitForNextUpdate, rerender, store } = setup({
            queryParams: {
                panes: JSON.stringify({
                    one: { datasource: 'loki-uid', queries: [{ expr: 'a' }] },
                }),
                schemaVersion: 1,
            },
        });
        yield waitForNextUpdate();
        let queries = (_g = store.getState().explore.panes['one']) === null || _g === void 0 ? void 0 : _g.queries;
        expect(queries).toHaveLength(1);
        expect(queries === null || queries === void 0 ? void 0 : queries[0]).toMatchObject({ expr: 'a' });
        rerender({
            children: null,
            params: {
                panes: JSON.stringify({
                    one: { datasource: 'loki-uid', queries: [{ expr: 'a' }, { expr: 'b' }] },
                }),
                schemaVersion: 1,
            },
        });
        yield waitForNextUpdate();
        queries = (_h = store.getState().explore.panes['one']) === null || _h === void 0 ? void 0 : _h.queries;
        expect(queries).toHaveLength(2);
        expect(queries === null || queries === void 0 ? void 0 : queries[0]).toMatchObject({ expr: 'a' });
        expect(queries === null || queries === void 0 ? void 0 : queries[1]).toMatchObject({ expr: 'b' });
        rerender({
            children: null,
            params: {
                panes: JSON.stringify({
                    one: { datasource: 'loki-uid', queries: [{ expr: 'a' }] },
                }),
                schemaVersion: 1,
            },
        });
        yield waitForNextUpdate();
        queries = (_j = store.getState().explore.panes['one']) === null || _j === void 0 ? void 0 : _j.queries;
        expect(queries).toHaveLength(1);
        expect(queries === null || queries === void 0 ? void 0 : queries[0]).toMatchObject({ expr: 'a' });
    }));
    it('Opens and closes the split pane if an a new pane is added or removed in the URL', () => __awaiter(void 0, void 0, void 0, function* () {
        const { waitForNextUpdate, rerender, store } = setup({
            queryParams: {
                panes: JSON.stringify({
                    one: { datasource: 'loki-uid', queries: [{ expr: 'a' }] },
                }),
                schemaVersion: 1,
            },
        });
        yield waitForNextUpdate();
        let panes = Object.keys(store.getState().explore.panes);
        expect(panes).toHaveLength(1);
        rerender({
            children: null,
            params: {
                panes: JSON.stringify({
                    one: { datasource: 'loki-uid', queries: [{ expr: 'a' }, { expr: 'b' }] },
                    two: { datasource: 'loki-uid', queries: [{ expr: 'a' }, { expr: 'b' }] },
                }),
                schemaVersion: 1,
            },
        });
        yield waitForNextUpdate();
        expect(Object.keys(store.getState().explore.panes)).toHaveLength(2);
        rerender({
            children: null,
            params: {
                panes: JSON.stringify({
                    one: { datasource: 'loki-uid', queries: [{ expr: 'a' }] },
                }),
                schemaVersion: 1,
            },
        });
        yield waitForNextUpdate();
        yield waitFor(() => {
            expect(Object.keys(store.getState().explore.panes)).toHaveLength(1);
        });
    }));
    it('Changes datasource when the datasource in the URL is updated', () => __awaiter(void 0, void 0, void 0, function* () {
        var _k, _l, _m, _o;
        const { waitForNextUpdate, rerender, store } = setup({
            queryParams: {
                panes: JSON.stringify({
                    one: { datasource: 'loki-uid', queries: [{ expr: 'a' }] },
                }),
                schemaVersion: 1,
            },
        });
        yield waitForNextUpdate();
        expect((_l = (_k = store.getState().explore.panes['one']) === null || _k === void 0 ? void 0 : _k.datasourceInstance) === null || _l === void 0 ? void 0 : _l.getRef()).toMatchObject({
            type: 'logs',
            uid: 'loki-uid',
        });
        rerender({
            children: null,
            params: {
                panes: JSON.stringify({
                    one: { datasource: 'elastic-uid', queries: [{ expr: 'a' }, { expr: 'b' }] },
                }),
                schemaVersion: 1,
            },
        });
        yield waitForNextUpdate();
        expect((_o = (_m = store.getState().explore.panes['one']) === null || _m === void 0 ? void 0 : _m.datasourceInstance) === null || _o === void 0 ? void 0 : _o.getRef()).toMatchObject({
            type: 'logs',
            uid: 'elastic-uid',
        });
    }));
    it('Changes time rage when the range in the URL is updated', () => __awaiter(void 0, void 0, void 0, function* () {
        var _p, _q;
        const { waitForNextUpdate, rerender, store } = setup({
            queryParams: {
                panes: JSON.stringify({
                    one: { datasource: 'loki-uid', queries: [{ expr: 'a' }], range: { from: 'now-1h', to: 'now' } },
                }),
                schemaVersion: 1,
            },
        });
        yield waitForNextUpdate();
        expect((_p = store.getState().explore.panes['one']) === null || _p === void 0 ? void 0 : _p.range.raw).toMatchObject({ from: 'now-1h', to: 'now' });
        rerender({
            children: null,
            params: {
                panes: JSON.stringify({
                    one: { datasource: 'loki-uid', queries: [{ expr: 'a' }], range: { from: 'now-6h', to: 'now' } },
                }),
                schemaVersion: 1,
            },
        });
        yield waitForNextUpdate();
        expect((_q = store.getState().explore.panes['one']) === null || _q === void 0 ? void 0 : _q.range.raw).toMatchObject({ from: 'now-6h', to: 'now' });
    }));
    it('Changes time range when the range in the URL is updated to absolute range', () => __awaiter(void 0, void 0, void 0, function* () {
        var _r, _s, _t;
        const { waitForNextUpdate, rerender, store } = setup({
            queryParams: {
                panes: JSON.stringify({
                    one: { datasource: 'loki-uid', queries: [{ expr: 'a' }], range: { from: 'now-1h', to: 'now' } },
                }),
                schemaVersion: 1,
            },
        });
        yield waitForNextUpdate();
        expect((_r = store.getState().explore.panes['one']) === null || _r === void 0 ? void 0 : _r.range.raw).toMatchObject({ from: 'now-1h', to: 'now' });
        rerender({
            children: null,
            params: {
                panes: JSON.stringify({
                    one: {
                        datasource: 'loki-uid',
                        queries: [{ expr: 'a' }],
                        range: { from: '1500000000000', to: '1500000001000' },
                    },
                }),
                schemaVersion: 1,
            },
        });
        yield waitForNextUpdate();
        expect((_s = store.getState().explore.panes['one']) === null || _s === void 0 ? void 0 : _s.range.raw.from.valueOf().toString()).toEqual('1500000000000');
        expect((_t = store.getState().explore.panes['one']) === null || _t === void 0 ? void 0 : _t.range.raw.to.valueOf().toString()).toEqual('1500000001000');
    }));
    it('uses the first query datasource if no root datasource is specified in the URL', () => __awaiter(void 0, void 0, void 0, function* () {
        var _u, _v;
        const { waitForNextUpdate, store } = setup({
            queryParams: {
                panes: JSON.stringify({
                    one: {
                        queries: [{ expr: 'b', datasource: { uid: 'loki-uid', type: 'logs' }, refId: 'B' }],
                    },
                }),
                schemaVersion: 1,
            },
        });
        yield waitForNextUpdate();
        expect((_v = (_u = store.getState().explore.panes['one']) === null || _u === void 0 ? void 0 : _u.datasourceInstance) === null || _v === void 0 ? void 0 : _v.getRef()).toMatchObject({
            uid: 'loki-uid',
            type: 'logs',
        });
    }));
    it('updates the URL opening and closing a pane datasource changes', () => __awaiter(void 0, void 0, void 0, function* () {
        var _w, _x;
        const { waitForNextUpdate, store, location } = setup({
            queryParams: {
                panes: JSON.stringify({
                    one: {
                        datasource: 'loki-uid',
                        queries: [{ expr: 'a', datasource: { uid: 'loki-uid', type: 'logs' }, refId: 'A' }],
                    },
                }),
                schemaVersion: 1,
            },
        });
        yield waitForNextUpdate();
        expect(location.getHistory().length).toBe(1);
        expect((_x = (_w = store.getState().explore.panes['one']) === null || _w === void 0 ? void 0 : _w.datasourceInstance) === null || _x === void 0 ? void 0 : _x.uid).toBe('loki-uid');
        act(() => {
            store.dispatch(splitOpen());
        });
        yield waitForNextUpdate();
        yield waitFor(() => __awaiter(void 0, void 0, void 0, function* () {
            expect(location.getHistory().length).toBe(2);
        }));
        expect(Object.keys(store.getState().explore.panes)).toHaveLength(2);
        act(() => {
            store.dispatch(splitClose('one'));
        });
        yield waitFor(() => __awaiter(void 0, void 0, void 0, function* () {
            expect(location.getHistory()).toHaveLength(3);
        }));
    }));
    it('filters out queries from the URL that do not have a datasource', () => __awaiter(void 0, void 0, void 0, function* () {
        var _y, _z;
        const { waitForNextUpdate, store } = setup({
            queryParams: {
                panes: JSON.stringify({
                    one: {
                        datasource: MIXED_DATASOURCE_NAME,
                        queries: [
                            { expr: 'a', refId: 'A' },
                            { expr: 'b', datasource: { uid: 'loki-uid', type: 'logs' }, refId: 'B' },
                        ],
                    },
                }),
                schemaVersion: 1,
            },
        });
        yield waitForNextUpdate();
        expect((_y = store.getState().explore.panes['one']) === null || _y === void 0 ? void 0 : _y.queries.length).toBe(1);
        expect((_z = store.getState().explore.panes['one']) === null || _z === void 0 ? void 0 : _z.queries[0]).toMatchObject({ expr: 'b', refId: 'B' });
    }));
});
//# sourceMappingURL=index.test.js.map