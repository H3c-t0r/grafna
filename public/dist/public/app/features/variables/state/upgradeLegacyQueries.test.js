import { __awaiter } from "tslib";
import { VariableSupportType } from '@grafana/data';
import { thunkTester } from '../../../../test/core/thunk/thunkTester';
import { customBuilder, queryBuilder } from '../shared/testing/builders';
import { TransactionStatus } from '../types';
import { toKeyedVariableIdentifier } from '../utils';
import { upgradeLegacyQueries } from './actions';
import { getPreloadedState } from './helpers';
import { toKeyedAction } from './keyedVariablesReducer';
import { changeVariableProp } from './sharedReducer';
function getTestContext({ query = '', variable, datasource, transactionStatus = TransactionStatus.Fetching, } = {}) {
    const key = 'key';
    variable =
        variable !== null && variable !== void 0 ? variable : queryBuilder()
            .withId('query')
            .withRootStateKey(key)
            .withName('query')
            .withQuery(query)
            .withDatasource({ uid: 'test-data', type: 'test-data' })
            .build();
    const templatingState = {
        transaction: { status: transactionStatus, uid: key, isDirty: false },
        variables: {
            [variable.id]: variable,
        },
    };
    const state = getPreloadedState(key, templatingState);
    datasource = datasource !== null && datasource !== void 0 ? datasource : {
        name: 'TestData',
        metricFindQuery: () => undefined,
        variables: { getType: () => VariableSupportType.Standard, toDataQuery: () => undefined },
    };
    const get = jest.fn().mockResolvedValue(datasource);
    const getDatasourceSrv = jest.fn().mockReturnValue({ get });
    const identifier = toKeyedVariableIdentifier(variable);
    return { key, state, get, getDatasourceSrv, identifier };
}
describe('upgradeLegacyQueries', () => {
    describe('when called with a query variable for a standard variable supported data source that has not been upgraded', () => {
        it('then it should dispatch changeVariableProp', () => __awaiter(void 0, void 0, void 0, function* () {
            const { key, state, identifier, get, getDatasourceSrv } = getTestContext({ query: '*' });
            const dispatchedActions = yield thunkTester(state)
                .givenThunk(upgradeLegacyQueries)
                .whenThunkIsDispatched(identifier, getDatasourceSrv);
            expect(dispatchedActions).toEqual([
                toKeyedAction(key, changeVariableProp({
                    type: 'query',
                    id: 'query',
                    data: {
                        propName: 'query',
                        propValue: {
                            refId: 'TestData-query-Variable-Query',
                            query: '*',
                        },
                    },
                })),
            ]);
            expect(get).toHaveBeenCalledTimes(1);
            expect(get).toHaveBeenCalledWith({ uid: 'test-data', type: 'test-data' });
        }));
        describe('but there is no ongoing transaction', () => {
            it('then it should not dispatch changeVariableProp', () => __awaiter(void 0, void 0, void 0, function* () {
                const { state, identifier, get, getDatasourceSrv } = getTestContext({
                    query: '*',
                    transactionStatus: TransactionStatus.NotStarted,
                });
                const dispatchedActions = yield thunkTester(state)
                    .givenThunk(upgradeLegacyQueries)
                    .whenThunkIsDispatched(identifier, getDatasourceSrv);
                expect(dispatchedActions).toEqual([]);
                expect(get).toHaveBeenCalledTimes(0);
            }));
        });
    });
    describe('when called with a query variable for a standard variable supported data source that has been upgraded', () => {
        it('then it should not dispatch any actions', () => __awaiter(void 0, void 0, void 0, function* () {
            const { state, identifier, get, getDatasourceSrv } = getTestContext({ query: { refId: 'A' } });
            const dispatchedActions = yield thunkTester(state)
                .givenThunk(upgradeLegacyQueries)
                .whenThunkIsDispatched(identifier, getDatasourceSrv);
            expect(dispatchedActions).toEqual([]);
            expect(get).toHaveBeenCalledTimes(1);
            expect(get).toHaveBeenCalledWith({ uid: 'test-data', type: 'test-data' });
        }));
    });
    describe('when called with a query variable for a legacy variable supported data source', () => {
        it('then it should not dispatch any actions', () => __awaiter(void 0, void 0, void 0, function* () {
            const datasource = {
                name: 'TestData',
                metricFindQuery: () => undefined,
            };
            const { state, identifier, get, getDatasourceSrv } = getTestContext({ datasource });
            const dispatchedActions = yield thunkTester(state)
                .givenThunk(upgradeLegacyQueries)
                .whenThunkIsDispatched(identifier, getDatasourceSrv);
            expect(dispatchedActions).toEqual([]);
            expect(get).toHaveBeenCalledTimes(1);
            expect(get).toHaveBeenCalledWith({ uid: 'test-data', type: 'test-data' });
        }));
    });
    describe('when called with a query variable for a custom variable supported data source', () => {
        it('then it should not dispatch any actions', () => __awaiter(void 0, void 0, void 0, function* () {
            const datasource = {
                name: 'TestData',
                metricFindQuery: () => undefined,
                variables: { getType: () => VariableSupportType.Custom, query: () => undefined, editor: {} },
            };
            const { state, identifier, get, getDatasourceSrv } = getTestContext({ datasource });
            const dispatchedActions = yield thunkTester(state)
                .givenThunk(upgradeLegacyQueries)
                .whenThunkIsDispatched(identifier, getDatasourceSrv);
            expect(dispatchedActions).toEqual([]);
            expect(get).toHaveBeenCalledTimes(1);
            expect(get).toHaveBeenCalledWith({ uid: 'test-data', type: 'test-data' });
        }));
    });
    describe('when called with a query variable for a datasource variable supported data source', () => {
        it('then it should not dispatch any actions', () => __awaiter(void 0, void 0, void 0, function* () {
            const datasource = {
                name: 'TestData',
                metricFindQuery: () => undefined,
                variables: { getType: () => VariableSupportType.Datasource },
            };
            const { state, identifier, get, getDatasourceSrv } = getTestContext({ datasource });
            const dispatchedActions = yield thunkTester(state)
                .givenThunk(upgradeLegacyQueries)
                .whenThunkIsDispatched(identifier, getDatasourceSrv);
            expect(dispatchedActions).toEqual([]);
            expect(get).toHaveBeenCalledTimes(1);
            expect(get).toHaveBeenCalledWith({ uid: 'test-data', type: 'test-data' });
        }));
    });
    describe('when called with a custom variable', () => {
        it('then it should not dispatch any actions', () => __awaiter(void 0, void 0, void 0, function* () {
            const variable = customBuilder().withId('custom').withRootStateKey('key').withName('custom').build();
            const { state, identifier, get, getDatasourceSrv } = getTestContext({ variable });
            const dispatchedActions = yield thunkTester(state)
                .givenThunk(upgradeLegacyQueries)
                .whenThunkIsDispatched(identifier, getDatasourceSrv);
            expect(dispatchedActions).toEqual([]);
            expect(get).toHaveBeenCalledTimes(0);
        }));
    });
});
//# sourceMappingURL=upgradeLegacyQueries.test.js.map