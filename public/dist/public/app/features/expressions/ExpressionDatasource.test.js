import { backendSrv } from 'app/core/services/backend_srv';
import { ExpressionDatasourceApi } from './ExpressionDatasource';
import { ExpressionQueryType } from './types';
jest.mock('@grafana/runtime', () => (Object.assign(Object.assign({}, jest.requireActual('@grafana/runtime')), { getBackendSrv: () => backendSrv, getTemplateSrv: () => ({
        replace: (val) => (val ? val.replace('$input', '10').replace('$window', '10s') : val),
    }) })));
describe('ExpressionDatasourceApi', () => {
    beforeEach(() => {
        jest.clearAllMocks();
    });
    describe('expression queries with template variables', () => {
        it('should interpolate template variables in expression query', () => {
            const ds = new ExpressionDatasourceApi({});
            const query = ds.applyTemplateVariables({ type: ExpressionQueryType.math, refId: 'B', expression: '$input + 5 + $A' }, {});
            expect(query.expression).toBe('10 + 5 + $A');
        });
        it('should interpolate template variables in expression query', () => {
            const ds = new ExpressionDatasourceApi({});
            const query = ds.applyTemplateVariables({ type: ExpressionQueryType.resample, refId: 'B', window: '$window' }, {});
            expect(query.window).toBe('10s');
        });
    });
});
//# sourceMappingURL=ExpressionDatasource.test.js.map