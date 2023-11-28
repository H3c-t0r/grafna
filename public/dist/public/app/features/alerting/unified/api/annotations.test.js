import '@grafana/runtime';
import { fetchAnnotations, sortStateHistory } from './annotations';
const get = jest.fn(() => {
    return new Promise((resolve) => {
        resolve(undefined);
    });
});
jest.mock('@grafana/runtime', () => ({
    getBackendSrv: () => ({
        get,
    }),
}));
describe('annotations', () => {
    beforeEach(() => get.mockClear());
    it('should fetch annotation for an alertId', () => {
        const ALERT_ID = 'abc123';
        fetchAnnotations(ALERT_ID);
        expect(get).toBeCalledWith('/api/annotations', { alertId: ALERT_ID });
    });
});
describe(sortStateHistory, () => {
    describe('should stably sort', () => {
        describe('when timeEnd is different', () => {
            it('should not sort by rule id', () => {
                let data = [
                    { timeEnd: 23, time: 22, id: 1 },
                    { timeEnd: 22, time: 21, id: 3 },
                    { timeEnd: 22, time: 22, id: 2 },
                    { timeEnd: 24, id: 3 },
                ];
                data.sort(sortStateHistory);
                expect(data[0].timeEnd).toBe(24);
                expect(data[1].timeEnd).toBe(23);
                expect(data[2].time).toBe(22);
                expect(data[3].id).toBe(3);
            });
        });
        describe('when only the rule id is different', () => {
            it('should sort by rule id', () => {
                let data = [
                    { timeEnd: 23, time: 22, id: 1 },
                    { timeEnd: 23, time: 22, id: 3 },
                    { timeEnd: 23, time: 22, id: 2 },
                    { timeEnd: 23, time: 22, id: 6 },
                ];
                data.sort(sortStateHistory);
                expect(data[0].id).toBe(6);
                expect(data[1].id).toBe(3);
                expect(data[2].id).toBe(2);
                expect(data[3].id).toBe(1);
            });
        });
    });
});
//# sourceMappingURL=annotations.test.js.map