import { __awaiter, __generator } from "tslib";
import { createGraphFrames } from './graphTransform';
import { testResponse, testResponseEdgesFields, testResponseNodesFields, toEdgesFrame, toNodesFrame, } from './testResponse';
describe('createGraphFrames', function () {
    it('transforms basic response into nodes and edges frame', function () { return __awaiter(void 0, void 0, void 0, function () {
        var frames;
        return __generator(this, function (_a) {
            frames = createGraphFrames(testResponse);
            expect(frames.length).toBe(2);
            expect(frames[0].fields).toMatchObject(testResponseNodesFields);
            expect(frames[1].fields).toMatchObject(testResponseEdgesFields);
            return [2 /*return*/];
        });
    }); });
    it('handles single span response', function () { return __awaiter(void 0, void 0, void 0, function () {
        var frames;
        return __generator(this, function (_a) {
            frames = createGraphFrames(singleSpanResponse);
            expect(frames.length).toBe(2);
            expect(frames[0].fields).toMatchObject(toNodesFrame([
                ['3fa414edcef6ad90'],
                ['tempo-querier'],
                ['HTTP GET - api_traces_traceid'],
                ['1049.14ms (100%)'],
                ['1049.14ms (100%)'],
                [1],
            ]));
            expect(frames[1].fields).toMatchObject(toEdgesFrame([[], [], []]));
            return [2 /*return*/];
        });
    }); });
    it('handles missing spans', function () { return __awaiter(void 0, void 0, void 0, function () {
        var frames;
        return __generator(this, function (_a) {
            frames = createGraphFrames(missingSpanResponse);
            expect(frames.length).toBe(2);
            expect(frames[0].length).toBe(2);
            expect(frames[1].length).toBe(0);
            return [2 /*return*/];
        });
    }); });
});
export var singleSpanResponse = {
    traceID: '3fa414edcef6ad90',
    spans: [
        {
            traceID: '3fa414edcef6ad90',
            spanID: '3fa414edcef6ad90',
            operationName: 'HTTP GET - api_traces_traceid',
            references: [],
            startTime: 1605873894680409,
            duration: 1049141,
            tags: [
                { key: 'sampler.type', type: 'string', value: 'probabilistic' },
                { key: 'sampler.param', type: 'float64', value: 1 },
            ],
            logs: [],
            processID: 'p1',
            warnings: null,
            flags: 0,
        },
    ],
    processes: {
        p1: {
            serviceName: 'tempo-querier',
            tags: [
                { key: 'cluster', type: 'string', value: 'ops-tools1' },
                { key: 'container', type: 'string', value: 'tempo-query' },
            ],
        },
    },
    warnings: null,
};
export var missingSpanResponse = {
    traceID: '3fa414edcef6ad90',
    spans: [
        {
            traceID: '3fa414edcef6ad90',
            spanID: '1',
            operationName: 'HTTP GET - api_traces_traceid',
            references: [],
            startTime: 1605873894680409,
            duration: 1049141,
            tags: [
                { key: 'sampler.type', type: 'string', value: 'probabilistic' },
                { key: 'sampler.param', type: 'float64', value: 1 },
            ],
            logs: [],
            processID: 'p1',
            warnings: null,
            flags: 0,
        },
        {
            traceID: '3fa414edcef6ad90',
            spanID: '2',
            operationName: 'HTTP GET - api_traces_traceid',
            references: [{ refType: 'CHILD_OF', traceID: '3fa414edcef6ad90', spanID: '3' }],
            startTime: 1605873894680409,
            duration: 1049141,
            tags: [
                { key: 'sampler.type', type: 'string', value: 'probabilistic' },
                { key: 'sampler.param', type: 'float64', value: 1 },
            ],
            logs: [],
            processID: 'p1',
            warnings: null,
            flags: 0,
        },
    ],
    processes: {
        p1: {
            serviceName: 'tempo-querier',
            tags: [
                { key: 'cluster', type: 'string', value: 'ops-tools1' },
                { key: 'container', type: 'string', value: 'tempo-query' },
            ],
        },
    },
    warnings: null,
};
//# sourceMappingURL=graphTransform.test.js.map