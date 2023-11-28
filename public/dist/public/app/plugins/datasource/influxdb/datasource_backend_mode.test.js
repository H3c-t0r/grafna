import { __awaiter } from "tslib";
import { of } from 'rxjs';
import { dateTime } from '@grafana/data/src';
import config from 'app/core/config';
import { TemplateSrv } from '../../../features/templating/template_srv';
import InfluxDatasource from './datasource';
import { getMockDSInstanceSettings, getMockInfluxDS, mockBackendService, mockInfluxFetchResponse, mockInfluxQueryWithTemplateVars, mockTemplateSrv, } from './mocks';
import { InfluxVersion } from './types';
config.featureToggles.influxdbBackendMigration = true;
const fetchMock = mockBackendService(mockInfluxFetchResponse());
describe('InfluxDataSource Backend Mode', () => {
    const text = 'interpolationText';
    const text2 = 'interpolationText2';
    const textWithoutFormatRegex = 'interpolationText,interpolationText2';
    const textWithFormatRegex = 'interpolationText|interpolationText2';
    const variableMap = {
        $interpolationVar: text,
        $interpolationVar2: text2,
    };
    const adhocFilters = [
        {
            key: 'adhoc',
            operator: '=',
            value: 'val',
            condition: '',
        },
    ];
    const templateSrv = mockTemplateSrv(jest.fn(() => {
        return adhocFilters;
    }), jest.fn((target, scopedVars, format) => {
        if (!format) {
            return variableMap[target] || '';
        }
        if (format === 'regex') {
            return textWithFormatRegex;
        }
        return textWithoutFormatRegex;
    }));
    let queryOptions;
    let influxQuery;
    const now = dateTime('2023-09-16T21:26:00Z');
    beforeEach(() => {
        queryOptions = {
            app: 'dashboard',
            interval: '10',
            intervalMs: 10,
            requestId: 'A-testing',
            startTime: 0,
            range: {
                from: dateTime(now).subtract(15, 'minutes'),
                to: now,
                raw: {
                    from: 'now-15m',
                    to: 'now',
                },
            },
            rangeRaw: {
                from: 'now-15m',
                to: 'now',
            },
            targets: [],
            timezone: 'UTC',
            scopedVars: {
                interval: { text: '1m', value: '1m' },
                __interval: { text: '1m', value: '1m' },
                __interval_ms: { text: 60000, value: 60000 },
            },
        };
        influxQuery = {
            refId: 'x',
            alias: '$interpolationVar',
            measurement: '$interpolationVar',
            policy: '$interpolationVar',
            limit: '$interpolationVar',
            slimit: '$interpolationVar',
            tz: '$interpolationVar',
            tags: [
                {
                    key: 'cpu',
                    operator: '=~',
                    value: '/^$interpolationVar,$interpolationVar2$/',
                },
            ],
            groupBy: [
                {
                    params: ['$interpolationVar'],
                    type: 'tag',
                },
            ],
            select: [
                [
                    {
                        params: ['$interpolationVar'],
                        type: 'field',
                    },
                ],
            ],
        };
    });
    describe('adhoc filters', () => {
        let fetchReq;
        const ctx = {
            ds: getMockInfluxDS(getMockDSInstanceSettings(), templateSrv),
        };
        beforeEach(() => __awaiter(void 0, void 0, void 0, function* () {
            fetchMock.mockImplementation((req) => {
                fetchReq = req.data;
                return of(mockInfluxFetchResponse());
            });
            const req = Object.assign(Object.assign({}, queryOptions), { targets: [...queryOptions.targets, Object.assign(Object.assign({}, influxQuery), { adhocFilters })] });
            yield ctx.ds.query(req);
        }));
        it('should add adhocFilters to the tags in the query', () => {
            var _a, _b, _c;
            expect(fetchMock).toHaveBeenCalled();
            expect(fetchReq).not.toBeNull();
            expect(fetchReq.queries.length).toBe(1);
            expect(fetchReq.queries[0].tags).toBeDefined();
            expect((_a = fetchReq.queries[0].tags) === null || _a === void 0 ? void 0 : _a.length).toBe(2);
            expect((_b = fetchReq.queries[0].tags) === null || _b === void 0 ? void 0 : _b[1].key).toBe(adhocFilters[0].key);
            expect((_c = fetchReq.queries[0].tags) === null || _c === void 0 ? void 0 : _c[1].value).toBe(adhocFilters[0].value);
        });
    });
    describe('when interpolating template variables', () => {
        const text = 'interpolationText';
        const text2 = 'interpolationText2';
        const textWithoutFormatRegex = 'interpolationText,interpolationText2';
        const textWithFormatRegex = 'interpolationText,interpolationText2';
        const variableMap = {
            $interpolationVar: text,
            $interpolationVar2: text2,
        };
        const adhocFilters = [
            {
                key: 'adhoc',
                operator: '=',
                value: 'val',
                condition: '',
            },
        ];
        const templateSrv = mockTemplateSrv(jest.fn((_) => adhocFilters), jest.fn((target, scopedVars, format) => {
            if (!format) {
                return variableMap[target] || '';
            }
            if (format === 'regex') {
                return textWithFormatRegex;
            }
            return textWithoutFormatRegex;
        }));
        const ds = new InfluxDatasource(getMockDSInstanceSettings(), templateSrv);
        function influxChecks(query) {
            var _a;
            expect(templateSrv.replace).toBeCalledTimes(10);
            expect(query.alias).toBe(text);
            expect(query.measurement).toBe(textWithFormatRegex);
            expect(query.policy).toBe(textWithFormatRegex);
            expect(query.limit).toBe(textWithFormatRegex);
            expect(query.slimit).toBe(textWithFormatRegex);
            expect(query.tz).toBe(text);
            expect(query.tags[0].value).toBe(textWithFormatRegex);
            expect(query.groupBy[0].params[0]).toBe(textWithFormatRegex);
            expect(query.select[0][0].params[0]).toBe(textWithFormatRegex);
            expect((_a = query.adhocFilters) === null || _a === void 0 ? void 0 : _a[0].key).toBe(adhocFilters[0].key);
        }
        it('should apply all template variables with InfluxQL mode', () => {
            ds.version = ds.version = InfluxVersion.InfluxQL;
            ds.access = 'proxy';
            const query = ds.applyTemplateVariables(mockInfluxQueryWithTemplateVars(adhocFilters), {
                interpolationVar: { text: text, value: text },
                interpolationVar2: { text: 'interpolationText2', value: 'interpolationText2' },
            });
            influxChecks(query);
        });
        it('should apply all scopedVars to tags', () => {
            var _a;
            ds.version = InfluxVersion.InfluxQL;
            ds.access = 'proxy';
            const query = ds.applyTemplateVariables(mockInfluxQueryWithTemplateVars(adhocFilters), {
                interpolationVar: { text: text, value: text },
                interpolationVar2: { text: 'interpolationText2', value: 'interpolationText2' },
            });
            if (!((_a = query.tags) === null || _a === void 0 ? void 0 : _a.length)) {
                throw new Error('Tags are not defined');
            }
            const value = query.tags[0].value;
            const scopedVars = 'interpolationText,interpolationText2';
            expect(value).toBe(scopedVars);
        });
    });
    describe('variable interpolation with chained variables with backend mode', () => {
        const mockTemplateService = new TemplateSrv();
        mockTemplateService.getAdhocFilters = jest.fn((_) => []);
        let ds = getMockInfluxDS(getMockDSInstanceSettings(), mockTemplateService);
        const fetchMockImpl = () => of({
            data: {
                status: 'success',
                results: [
                    {
                        series: [
                            {
                                name: 'measurement',
                                columns: ['name'],
                                values: [['cpu']],
                            },
                        ],
                    },
                ],
            },
        });
        beforeEach(() => {
            jest.clearAllMocks();
            fetchMock.mockImplementation(fetchMockImpl);
        });
        it('should render chained regex variables with floating point number', () => {
            ds.metricFindQuery(`SELECT sum("piece_count") FROM "rp"."pdata" WHERE diameter <= $maxSED`, {
                scopedVars: { maxSED: { text: '8.1', value: '8.1' } },
            });
            const qe = `SELECT sum("piece_count") FROM "rp"."pdata" WHERE diameter <= 8.1`;
            const qData = fetchMock.mock.calls[0][0].data.queries[0].query;
            expect(qData).toBe(qe);
        });
        it('should render chained regex variables with URL', () => {
            ds.metricFindQuery('SHOW TAG VALUES WITH KEY = "agent_url" WHERE agent_url =~ /^$var1$/', {
                scopedVars: {
                    var1: {
                        text: 'https://aaaa-aa-aaa.bbb.ccc.ddd:8443/ggggg',
                        value: 'https://aaaa-aa-aaa.bbb.ccc.ddd:8443/ggggg',
                    },
                },
            });
            const qe = `SHOW TAG VALUES WITH KEY = "agent_url" WHERE agent_url =~ /^https:\\/\\/aaaa-aa-aaa\\.bbb\\.ccc\\.ddd:8443\\/ggggg$/`;
            const qData = fetchMock.mock.calls[0][0].data.queries[0].query;
            expect(qData).toBe(qe);
        });
        it('should render chained regex variables with floating point number and url', () => {
            ds.metricFindQuery('SELECT sum("piece_count") FROM "rp"."pdata" WHERE diameter <= $maxSED AND agent_url =~ /^$var1$/', {
                scopedVars: {
                    var1: {
                        text: 'https://aaaa-aa-aaa.bbb.ccc.ddd:8443/ggggg',
                        value: 'https://aaaa-aa-aaa.bbb.ccc.ddd:8443/ggggg',
                    },
                    maxSED: { text: '8.1', value: '8.1' },
                },
            });
            const qe = `SELECT sum("piece_count") FROM "rp"."pdata" WHERE diameter <= 8.1 AND agent_url =~ /^https:\\/\\/aaaa-aa-aaa\\.bbb\\.ccc\\.ddd:8443\\/ggggg$/`;
            const qData = fetchMock.mock.calls[0][0].data.queries[0].query;
            expect(qData).toBe(qe);
        });
    });
});
//# sourceMappingURL=datasource_backend_mode.test.js.map