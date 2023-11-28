import { debounce, memoize } from 'lodash';
export default (func, wait = 7000) => {
    const mem = memoize((...args) => debounce(func, wait, {
        leading: true,
    }), (...args) => JSON.stringify(args));
    return (...args) => mem(...args)(...args);
};
//# sourceMappingURL=memoizedDebounce.js.map