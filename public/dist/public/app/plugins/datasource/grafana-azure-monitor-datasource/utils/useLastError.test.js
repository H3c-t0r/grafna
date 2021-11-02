import { renderHook, act } from '@testing-library/react-hooks';
import useLastError from './useLastError';
describe('AzureMonitor: useLastError', function () {
    it('returns the set error', function () {
        var result = renderHook(function () { return useLastError(); }).result;
        act(function () {
            result.current[1]('component-a', new Error('an error'));
        });
        expect(result.current[0]).toBe('an error');
    });
    it('returns the most recent error', function () {
        var result = renderHook(function () { return useLastError(); }).result;
        act(function () {
            result.current[1]('component-a', new Error('component a error'));
            result.current[1]('component-b', new Error('component b error'));
            result.current[1]('component-a', new Error('second component a error'));
        });
        expect(result.current[0]).toBe('second component a error');
    });
});
//# sourceMappingURL=useLastError.test.js.map