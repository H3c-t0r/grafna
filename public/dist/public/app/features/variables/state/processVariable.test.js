import { __awaiter } from "tslib";
import { setDataSourceSrv } from '@grafana/runtime';
import { reduxTester } from '../../../../test/core/redux/reduxTester';
import { variableAdapters } from '../adapters';
import { createCustomVariableAdapter } from '../custom/adapter';
import { createCustomOptionsFromQuery } from '../custom/reducer';
import { setVariableQueryRunner, VariableQueryRunner } from '../query/VariableQueryRunner';
import { createQueryVariableAdapter } from '../query/adapter';
import { updateVariableOptions } from '../query/reducer';
import { customBuilder, queryBuilder } from '../shared/testing/builders';
import { VariableRefresh } from '../types';
import { toKeyedVariableIdentifier, toVariablePayload } from '../utils';
import { initDashboardTemplating, processVariable } from './actions';
import { getTemplatingRootReducer } from './helpers';
import { toKeyedAction } from './keyedVariablesReducer';
import { setCurrentVariableValue, variableStateCompleted, variableStateFetching } from './sharedReducer';
import { variablesInitTransaction } from './transactionReducer';
jest.mock('app/features/dashboard/services/TimeSrv', () => ({
    getTimeSrv: jest.fn().mockReturnValue({
        timeRange: jest.fn().mockReturnValue({
            from: '2001-01-01T01:00:00.000Z',
            to: '2001-01-01T02:00:00.000Z',
            raw: {
                from: 'now-1h',
                to: 'now',
            },
        }),
    }),
}));
setDataSourceSrv({
    get: jest.fn().mockResolvedValue({
        metricFindQuery: jest.fn().mockImplementation((query, options) => {
            if (query === '$custom.*') {
                return Promise.resolve([
                    { value: 'AA', text: 'AA' },
                    { value: 'AB', text: 'AB' },
                    { value: 'AC', text: 'AC' },
                ]);
            }
            if (query === '$custom.$queryDependsOnCustom.*') {
                return Promise.resolve([
                    { value: 'AAA', text: 'AAA' },
                    { value: 'AAB', text: 'AAB' },
                    { value: 'AAC', text: 'AAC' },
                ]);
            }
            if (query === '*') {
                return Promise.resolve([
                    { value: 'A', text: 'A' },
                    { value: 'B', text: 'B' },
                    { value: 'C', text: 'C' },
                ]);
            }
            return Promise.resolve([]);
        }),
    }),
});
variableAdapters.setInit(() => [createCustomVariableAdapter(), createQueryVariableAdapter()]);
describe('processVariable', () => {
    // these following processVariable tests will test the following base setup
    // custom doesn't depend on any other variable
    // queryDependsOnCustom depends on custom
    // queryNoDepends doesn't depend on any other variable
    const key = 'key';
    const getTestContext = () => {
        const custom = customBuilder()
            .withId('custom')
            .withRootStateKey(key)
            .withName('custom')
            .withQuery('A,B,C')
            .withOptions('A', 'B', 'C')
            .withCurrent('A')
            .build();
        const queryDependsOnCustom = queryBuilder()
            .withId('queryDependsOnCustom')
            .withRootStateKey(key)
            .withName('queryDependsOnCustom')
            .withQuery('$custom.*')
            .withOptions('AA', 'AB', 'AC')
            .withCurrent('AA')
            .build();
        const queryNoDepends = queryBuilder()
            .withId('queryNoDepends')
            .withRootStateKey(key)
            .withName('queryNoDepends')
            .withQuery('*')
            .withOptions('A', 'B', 'C')
            .withCurrent('A')
            .build();
        const list = [custom, queryDependsOnCustom, queryNoDepends];
        const dashboard = { templating: { list } };
        setVariableQueryRunner(new VariableQueryRunner());
        return {
            key,
            custom,
            queryDependsOnCustom,
            queryNoDepends,
            dashboard,
        };
    };
    // testing processVariable for the custom variable from case described above
    describe('when processVariable is dispatched for a custom variable without dependencies', () => {
        describe('and queryParams does not match variable', () => {
            it('then correct actions are dispatched', () => __awaiter(void 0, void 0, void 0, function* () {
                const { key, dashboard, custom } = getTestContext();
                const queryParams = {};
                const tester = yield reduxTester()
                    .givenRootReducer(getTemplatingRootReducer())
                    .whenActionIsDispatched(toKeyedAction(key, variablesInitTransaction({ uid: key })))
                    .whenActionIsDispatched(initDashboardTemplating(key, dashboard))
                    .whenAsyncActionIsDispatched(processVariable(toKeyedVariableIdentifier(custom), queryParams), true);
                yield tester.thenDispatchedActionsPredicateShouldEqual((dispatchedActions) => {
                    expect(dispatchedActions.length).toEqual(4);
                    expect(dispatchedActions[0]).toEqual(toKeyedAction(key, variableStateFetching(toVariablePayload(custom))));
                    expect(dispatchedActions[1]).toEqual(toKeyedAction(key, createCustomOptionsFromQuery(toVariablePayload(custom))));
                    expect(dispatchedActions[2].type).toEqual('templating/keyed/shared/setCurrentVariableValue');
                    expect(dispatchedActions[3]).toEqual(toKeyedAction(key, variableStateCompleted(toVariablePayload(custom))));
                    return true;
                });
            }));
        });
        describe('and queryParams does match variable', () => {
            it('then correct actions are dispatched', () => __awaiter(void 0, void 0, void 0, function* () {
                const { key, dashboard, custom } = getTestContext();
                const queryParams = { 'var-custom': 'B' };
                const tester = yield reduxTester()
                    .givenRootReducer(getTemplatingRootReducer())
                    .whenActionIsDispatched(toKeyedAction(key, variablesInitTransaction({ uid: key })))
                    .whenActionIsDispatched(initDashboardTemplating(key, dashboard))
                    .whenAsyncActionIsDispatched(processVariable(toKeyedVariableIdentifier(custom), queryParams), true);
                yield tester.thenDispatchedActionsShouldEqual(toKeyedAction(key, setCurrentVariableValue(toVariablePayload({ type: 'custom', id: 'custom' }, { option: { text: 'B', value: 'B', selected: false } }))), toKeyedAction(key, variableStateCompleted(toVariablePayload(custom))));
            }));
        });
    });
    // testing processVariable for the queryNoDepends variable from case described above
    describe('when processVariable is dispatched for a query variable without dependencies', () => {
        describe('and queryParams does not match variable', () => {
            const queryParams = {};
            describe('and refresh is VariableRefresh.never', () => {
                const refresh = VariableRefresh.never;
                it('then correct actions are dispatched', () => __awaiter(void 0, void 0, void 0, function* () {
                    const { dashboard, key, queryNoDepends } = getTestContext();
                    queryNoDepends.refresh = refresh;
                    const tester = yield reduxTester()
                        .givenRootReducer(getTemplatingRootReducer())
                        .whenActionIsDispatched(toKeyedAction(key, variablesInitTransaction({ uid: key })))
                        .whenActionIsDispatched(initDashboardTemplating(key, dashboard))
                        .whenAsyncActionIsDispatched(processVariable(toKeyedVariableIdentifier(queryNoDepends), queryParams), true);
                    yield tester.thenDispatchedActionsShouldEqual(toKeyedAction(key, variableStateCompleted(toVariablePayload(queryNoDepends))));
                }));
            });
            it.each `
        refresh
        ${VariableRefresh.onDashboardLoad}
        ${VariableRefresh.onTimeRangeChanged}
      `('and refresh is $refresh then correct actions are dispatched', ({ refresh }) => __awaiter(void 0, void 0, void 0, function* () {
                const { dashboard, key, queryNoDepends } = getTestContext();
                queryNoDepends.refresh = refresh;
                const tester = yield reduxTester()
                    .givenRootReducer(getTemplatingRootReducer())
                    .whenActionIsDispatched(toKeyedAction(key, variablesInitTransaction({ uid: key })))
                    .whenActionIsDispatched(initDashboardTemplating(key, dashboard))
                    .whenAsyncActionIsDispatched(processVariable(toKeyedVariableIdentifier(queryNoDepends), queryParams), true);
                yield tester.thenDispatchedActionsShouldEqual(toKeyedAction(key, variableStateFetching(toVariablePayload({ type: 'query', id: 'queryNoDepends' }))), toKeyedAction(key, updateVariableOptions(toVariablePayload({ type: 'query', id: 'queryNoDepends' }, {
                    results: [
                        { value: 'A', text: 'A' },
                        { value: 'B', text: 'B' },
                        { value: 'C', text: 'C' },
                    ],
                    templatedRegex: '',
                }))), toKeyedAction(key, setCurrentVariableValue(toVariablePayload({ type: 'query', id: 'queryNoDepends' }, { option: { text: 'A', value: 'A', selected: false } }))), toKeyedAction(key, variableStateCompleted(toVariablePayload({ type: 'query', id: 'queryNoDepends' }))));
            }));
        });
        describe('and queryParams does match variable', () => {
            const queryParams = { 'var-queryNoDepends': 'B' };
            describe('and refresh is VariableRefresh.never', () => {
                const refresh = VariableRefresh.never;
                it('then correct actions are dispatched', () => __awaiter(void 0, void 0, void 0, function* () {
                    const { dashboard, key, queryNoDepends } = getTestContext();
                    queryNoDepends.refresh = refresh;
                    const tester = yield reduxTester()
                        .givenRootReducer(getTemplatingRootReducer())
                        .whenActionIsDispatched(toKeyedAction(key, variablesInitTransaction({ uid: key })))
                        .whenActionIsDispatched(initDashboardTemplating(key, dashboard))
                        .whenAsyncActionIsDispatched(processVariable(toKeyedVariableIdentifier(queryNoDepends), queryParams), true);
                    yield tester.thenDispatchedActionsShouldEqual(toKeyedAction(key, setCurrentVariableValue(toVariablePayload({ type: 'query', id: 'queryNoDepends' }, { option: { text: 'B', value: 'B', selected: false } }))), toKeyedAction(key, variableStateCompleted(toVariablePayload({ type: 'query', id: 'queryNoDepends' }))));
                }));
            });
            it.each `
        refresh
        ${VariableRefresh.onDashboardLoad}
        ${VariableRefresh.onTimeRangeChanged}
      `('and refresh is $refresh then correct actions are dispatched', ({ refresh }) => __awaiter(void 0, void 0, void 0, function* () {
                const { dashboard, key, queryNoDepends } = getTestContext();
                queryNoDepends.refresh = refresh;
                const tester = yield reduxTester()
                    .givenRootReducer(getTemplatingRootReducer())
                    .whenActionIsDispatched(toKeyedAction(key, variablesInitTransaction({ uid: key })))
                    .whenActionIsDispatched(initDashboardTemplating(key, dashboard))
                    .whenAsyncActionIsDispatched(processVariable(toKeyedVariableIdentifier(queryNoDepends), queryParams), true);
                yield tester.thenDispatchedActionsShouldEqual(toKeyedAction(key, variableStateFetching(toVariablePayload({ type: 'query', id: 'queryNoDepends' }))), toKeyedAction(key, updateVariableOptions(toVariablePayload({ type: 'query', id: 'queryNoDepends' }, {
                    results: [
                        { value: 'A', text: 'A' },
                        { value: 'B', text: 'B' },
                        { value: 'C', text: 'C' },
                    ],
                    templatedRegex: '',
                }))), toKeyedAction(key, setCurrentVariableValue(toVariablePayload({ type: 'query', id: 'queryNoDepends' }, { option: { text: 'A', value: 'A', selected: false } }))), toKeyedAction(key, variableStateCompleted(toVariablePayload({ type: 'query', id: 'queryNoDepends' }))), toKeyedAction(key, setCurrentVariableValue(toVariablePayload({ type: 'query', id: 'queryNoDepends' }, { option: { text: 'B', value: 'B', selected: false } }))));
            }));
        });
    });
    // testing processVariable for the queryDependsOnCustom variable from case described above
    describe('when processVariable is dispatched for a query variable with one dependency', () => {
        describe('and queryParams does not match variable', () => {
            const queryParams = {};
            describe('and refresh is VariableRefresh.never', () => {
                const refresh = VariableRefresh.never;
                it('then correct actions are dispatched', () => __awaiter(void 0, void 0, void 0, function* () {
                    const { key, dashboard, custom, queryDependsOnCustom } = getTestContext();
                    queryDependsOnCustom.refresh = refresh;
                    const customProcessed = yield reduxTester()
                        .givenRootReducer(getTemplatingRootReducer())
                        .whenActionIsDispatched(toKeyedAction(key, variablesInitTransaction({ uid: key })))
                        .whenActionIsDispatched(initDashboardTemplating(key, dashboard))
                        .whenAsyncActionIsDispatched(processVariable(toKeyedVariableIdentifier(custom), queryParams)); // Need to process this dependency otherwise we never complete the promise chain
                    const tester = yield customProcessed.whenAsyncActionIsDispatched(processVariable(toKeyedVariableIdentifier(queryDependsOnCustom), queryParams), true);
                    yield tester.thenDispatchedActionsShouldEqual(toKeyedAction(key, variableStateCompleted(toVariablePayload({ type: 'query', id: 'queryDependsOnCustom' }))));
                }));
            });
            it.each `
        refresh
        ${VariableRefresh.onDashboardLoad}
        ${VariableRefresh.onTimeRangeChanged}
      `('and refresh is $refresh then correct actions are dispatched', ({ refresh }) => __awaiter(void 0, void 0, void 0, function* () {
                const { key, dashboard, custom, queryDependsOnCustom } = getTestContext();
                queryDependsOnCustom.refresh = refresh;
                const customProcessed = yield reduxTester()
                    .givenRootReducer(getTemplatingRootReducer())
                    .whenActionIsDispatched(toKeyedAction(key, variablesInitTransaction({ uid: key })))
                    .whenActionIsDispatched(initDashboardTemplating(key, dashboard))
                    .whenAsyncActionIsDispatched(processVariable(toKeyedVariableIdentifier(custom), queryParams)); // Need to process this dependency otherwise we never complete the promise chain
                const tester = yield customProcessed.whenAsyncActionIsDispatched(processVariable(toKeyedVariableIdentifier(queryDependsOnCustom), queryParams), true);
                yield tester.thenDispatchedActionsShouldEqual(toKeyedAction(key, variableStateFetching(toVariablePayload({ type: 'query', id: 'queryDependsOnCustom' }))), toKeyedAction(key, updateVariableOptions(toVariablePayload({ type: 'query', id: 'queryDependsOnCustom' }, {
                    results: [
                        { value: 'AA', text: 'AA' },
                        { value: 'AB', text: 'AB' },
                        { value: 'AC', text: 'AC' },
                    ],
                    templatedRegex: '',
                }))), toKeyedAction(key, setCurrentVariableValue(toVariablePayload({ type: 'query', id: 'queryDependsOnCustom' }, { option: { text: 'AA', value: 'AA', selected: false } }))), toKeyedAction(key, variableStateCompleted(toVariablePayload({ type: 'query', id: 'queryDependsOnCustom' }))));
            }));
        });
        describe('and queryParams does match variable', () => {
            const queryParams = { 'var-queryDependsOnCustom': 'AB' };
            describe('and refresh is VariableRefresh.never', () => {
                const refresh = VariableRefresh.never;
                it('then correct actions are dispatched', () => __awaiter(void 0, void 0, void 0, function* () {
                    const { key, dashboard, custom, queryDependsOnCustom } = getTestContext();
                    queryDependsOnCustom.refresh = refresh;
                    const customProcessed = yield reduxTester()
                        .givenRootReducer(getTemplatingRootReducer())
                        .whenActionIsDispatched(toKeyedAction(key, variablesInitTransaction({ uid: key })))
                        .whenActionIsDispatched(initDashboardTemplating(key, dashboard))
                        .whenAsyncActionIsDispatched(processVariable(toKeyedVariableIdentifier(custom), queryParams)); // Need to process this dependency otherwise we never complete the promise chain
                    const tester = yield customProcessed.whenAsyncActionIsDispatched(processVariable(toKeyedVariableIdentifier(queryDependsOnCustom), queryParams), true);
                    yield tester.thenDispatchedActionsShouldEqual(toKeyedAction(key, setCurrentVariableValue(toVariablePayload({ type: 'query', id: 'queryDependsOnCustom' }, { option: { text: 'AB', value: 'AB', selected: false } }))), toKeyedAction(key, variableStateCompleted(toVariablePayload({ type: 'query', id: 'queryDependsOnCustom' }))));
                }));
            });
            it.each `
        refresh
        ${VariableRefresh.onDashboardLoad}
        ${VariableRefresh.onTimeRangeChanged}
      `('and refresh is $refresh then correct actions are dispatched', ({ refresh }) => __awaiter(void 0, void 0, void 0, function* () {
                const { key, dashboard, custom, queryDependsOnCustom } = getTestContext();
                queryDependsOnCustom.refresh = refresh;
                const customProcessed = yield reduxTester()
                    .givenRootReducer(getTemplatingRootReducer())
                    .whenActionIsDispatched(toKeyedAction(key, variablesInitTransaction({ uid: key })))
                    .whenActionIsDispatched(initDashboardTemplating(key, dashboard))
                    .whenAsyncActionIsDispatched(processVariable(toKeyedVariableIdentifier(custom), queryParams)); // Need to process this dependency otherwise we never complete the promise chain
                const tester = yield customProcessed.whenAsyncActionIsDispatched(processVariable(toKeyedVariableIdentifier(queryDependsOnCustom), queryParams), true);
                yield tester.thenDispatchedActionsShouldEqual(toKeyedAction(key, variableStateFetching(toVariablePayload({ type: 'query', id: 'queryDependsOnCustom' }))), toKeyedAction(key, updateVariableOptions(toVariablePayload({ type: 'query', id: 'queryDependsOnCustom' }, {
                    results: [
                        { value: 'AA', text: 'AA' },
                        { value: 'AB', text: 'AB' },
                        { value: 'AC', text: 'AC' },
                    ],
                    templatedRegex: '',
                }))), toKeyedAction(key, setCurrentVariableValue(toVariablePayload({ type: 'query', id: 'queryDependsOnCustom' }, { option: { text: 'AA', value: 'AA', selected: false } }))), toKeyedAction(key, variableStateCompleted(toVariablePayload({ type: 'query', id: 'queryDependsOnCustom' }))), toKeyedAction(key, setCurrentVariableValue(toVariablePayload({ type: 'query', id: 'queryDependsOnCustom' }, { option: { text: 'AB', value: 'AB', selected: false } }))));
            }));
        });
    });
});
//# sourceMappingURL=processVariable.test.js.map