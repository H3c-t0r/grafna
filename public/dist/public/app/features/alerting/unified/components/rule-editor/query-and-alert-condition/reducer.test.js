import { getDefaultRelativeTimeRange } from '@grafana/data';
import { getDataSourceSrv } from '@grafana/runtime/src/services/__mocks__/dataSourceSrv';
import { dataSource as expressionDatasource } from 'app/features/expressions/ExpressionDatasource';
import { ExpressionDatasourceUID, ExpressionQueryType } from 'app/features/expressions/types';
import { defaultCondition } from 'app/features/expressions/utils/expressionTypes';
import { addNewDataQuery, addNewExpression, duplicateQuery, queriesAndExpressionsReducer, removeExpression, rewireExpressions, setDataQueries, updateExpression, updateExpressionRefId, updateExpressionTimeRange, updateExpressionType, } from './reducer';
jest.mock('@grafana/runtime', () => (Object.assign(Object.assign({}, jest.requireActual('@grafana/runtime')), { getDataSourceSrv: getDataSourceSrv })));
const alertQuery = {
    refId: 'A',
    queryType: 'query',
    datasourceUid: 'abc123',
    model: {
        refId: 'A',
    },
};
const expressionQuery = {
    datasourceUid: ExpressionDatasourceUID,
    model: expressionDatasource.newQuery({
        type: ExpressionQueryType.classic,
        conditions: [Object.assign(Object.assign({}, defaultCondition), { query: { params: ['A'] } })],
        expression: '',
        refId: 'B',
    }),
    refId: 'B',
    queryType: '',
};
describe('Query and expressions reducer', () => {
    it('should return initial state', () => {
        expect(queriesAndExpressionsReducer(undefined, { type: undefined })).toEqual({
            queries: [],
        });
    });
    it('should duplicate query', () => {
        const initialState = {
            queries: [alertQuery],
        };
        const newState = queriesAndExpressionsReducer(initialState, duplicateQuery(alertQuery));
        const newQuery = newState.queries.at(-1);
        expect(newState).toMatchSnapshot();
        expect(newQuery).toHaveProperty('relativeTimeRange', getDefaultRelativeTimeRange());
    });
    it('should duplicate query and copy time range', () => {
        const initialState = {
            queries: [alertQuery],
        };
        const customTimeRange = {
            from: -200,
            to: 800,
        };
        const query = Object.assign(Object.assign({}, initialState.queries[0]), { relativeTimeRange: customTimeRange });
        const previousState = {
            queries: [query],
        };
        const newState = queriesAndExpressionsReducer(previousState, duplicateQuery(query));
        const newQuery = newState.queries.at(-1);
        expect(newQuery).toHaveProperty('relativeTimeRange', customTimeRange);
    });
    it('should add query', () => {
        const initialState = {
            queries: [alertQuery],
        };
        const newState = queriesAndExpressionsReducer(initialState, addNewDataQuery());
        expect(newState.queries).toHaveLength(2);
        expect(newState).toMatchSnapshot();
    });
    it('should set data queries', () => {
        const initialState = {
            queries: [alertQuery, expressionQuery],
        };
        const newState = queriesAndExpressionsReducer(initialState, setDataQueries([]));
        expect(newState.queries).toHaveLength(1);
        expect(newState).toMatchSnapshot();
    });
    it('should add a new expression', () => {
        const initialState = {
            queries: [alertQuery],
        };
        const newState = queriesAndExpressionsReducer(initialState, addNewExpression(ExpressionQueryType.math));
        expect(newState.queries).toHaveLength(2);
        expect(newState).toMatchSnapshot();
    });
    it('should remove an expression or alert query', () => {
        const initialState = {
            queries: [alertQuery, expressionQuery],
        };
        let stateWithoutB = queriesAndExpressionsReducer(initialState, removeExpression('B'));
        expect(stateWithoutB.queries).toHaveLength(1);
        expect(stateWithoutB).toMatchSnapshot();
        let stateWithoutAOrB = queriesAndExpressionsReducer(stateWithoutB, removeExpression('A'));
        expect(stateWithoutAOrB.queries).toHaveLength(0);
    });
    it('should update an expression', () => {
        const newExpression = Object.assign(Object.assign({}, expressionQuery.model), { type: ExpressionQueryType.math });
        const initialState = {
            queries: [expressionQuery],
        };
        const newState = queriesAndExpressionsReducer(initialState, updateExpression(newExpression));
        expect(newState).toMatchSnapshot();
    });
    it('should use time range from data source when updating an expression', () => {
        const expressionQuery = {
            refId: 'B',
            queryType: 'expression',
            datasourceUid: '__expr__',
            relativeTimeRange: { from: 900, to: 1000 },
            model: {
                queryType: 'query',
                datasource: '__expr__',
                refId: 'B',
                expression: 'C',
                type: ExpressionQueryType.classic,
                window: '10s',
            },
        };
        const expressionQuery2 = {
            refId: 'C',
            queryType: 'expression',
            datasourceUid: '__expr__',
            relativeTimeRange: { from: 1, to: 3 },
            model: {
                queryType: 'query',
                datasource: '__expr__',
                refId: 'C',
                expression: 'A',
                type: ExpressionQueryType.classic,
                window: '10s',
            },
        };
        const queryA = {
            refId: 'A',
            relativeTimeRange: { from: 900, to: 1000 },
            datasourceUid: 'dsuid',
            model: { refId: 'A' },
            queryType: 'query',
        };
        const newExpression = Object.assign(Object.assign({}, expressionQuery2.model), { type: ExpressionQueryType.resample });
        const initialState = {
            queries: [queryA, expressionQuery, expressionQuery2],
        };
        const newState = queriesAndExpressionsReducer(initialState, updateExpression(newExpression));
        expect(newState).toStrictEqual({
            queries: [
                {
                    refId: 'A',
                    relativeTimeRange: { from: 900, to: 1000 },
                    datasourceUid: 'dsuid',
                    model: { refId: 'A' },
                    queryType: 'query',
                },
                {
                    datasourceUid: '__expr__',
                    relativeTimeRange: { from: 900, to: 1000 },
                    model: {
                        datasource: '__expr__',
                        expression: 'C',
                        queryType: 'query',
                        refId: 'B',
                        type: 'classic_conditions',
                        window: '10s',
                    },
                    queryType: 'expression',
                    refId: 'B',
                },
                {
                    datasourceUid: '__expr__',
                    relativeTimeRange: { from: 900, to: 1000 },
                    model: {
                        datasource: '__expr__',
                        expression: 'A',
                        queryType: 'query',
                        refId: 'C',
                        type: 'resample',
                        window: '10s',
                    },
                    queryType: 'expression',
                    refId: 'C',
                },
            ],
        });
    });
    it('Should update time range for all expressions that have this data source when dispatching updateExpressionTimeRange', () => {
        const expressionQuery = {
            refId: 'B',
            queryType: 'expression',
            datasourceUid: '__expr__',
            model: {
                queryType: 'query',
                datasource: '__expr__',
                refId: 'B',
                expression: 'A',
                type: ExpressionQueryType.classic,
                window: '10s',
            },
        };
        const customTimeRange = { from: 900, to: 1000 };
        const queryA = {
            refId: 'A',
            relativeTimeRange: customTimeRange,
            datasourceUid: 'dsuid',
            model: { refId: 'A' },
            queryType: 'query',
        };
        const initialState = {
            queries: [queryA, expressionQuery],
        };
        const newState = queriesAndExpressionsReducer(initialState, updateExpressionTimeRange());
        expect(newState).toStrictEqual({
            queries: [
                {
                    refId: 'A',
                    relativeTimeRange: customTimeRange,
                    datasourceUid: 'dsuid',
                    model: { refId: 'A' },
                    queryType: 'query',
                },
                {
                    datasourceUid: '__expr__',
                    model: {
                        datasource: '__expr__',
                        expression: 'A',
                        queryType: 'query',
                        refId: 'B',
                        type: ExpressionQueryType.classic,
                        window: '10s',
                    },
                    queryType: 'expression',
                    refId: 'B',
                    relativeTimeRange: {
                        from: 900,
                        to: 1000,
                    },
                },
            ],
        });
    });
    it('should update an expression refId and rewire expressions', () => {
        const initialState = {
            queries: [alertQuery, expressionQuery],
        };
        const newState = queriesAndExpressionsReducer(initialState, updateExpressionRefId({
            oldRefId: 'A',
            newRefId: 'C',
        }));
        expect(newState).toMatchSnapshot();
    });
    it('should not update an expression when the refId exists', () => {
        const initialState = {
            queries: [alertQuery, expressionQuery],
        };
        const newState = queriesAndExpressionsReducer(initialState, updateExpressionRefId({
            oldRefId: 'A',
            newRefId: 'B',
        }));
        expect(newState).toEqual(initialState);
    });
    it('should rewire expressions', () => {
        const initialState = {
            queries: [alertQuery, expressionQuery],
        };
        const newState = queriesAndExpressionsReducer(initialState, rewireExpressions({
            oldRefId: 'A',
            newRefId: 'C',
        }));
        expect(newState).toMatchSnapshot();
    });
    it('should update expression type', () => {
        const initialState = {
            queries: [alertQuery, expressionQuery],
        };
        const newState = queriesAndExpressionsReducer(initialState, updateExpressionType({
            refId: 'B',
            type: ExpressionQueryType.reduce,
        }));
        expect(newState).toMatchSnapshot();
    });
});
//# sourceMappingURL=reducer.test.js.map