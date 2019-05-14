import { VizOrientation } from '@grafana/ui';
export var defaults = {
    valueOptions: {
        prefix: '',
        suffix: '',
        decimals: null,
        stat: 'avg',
        unit: 'none',
    },
    valueMappings: [],
    thresholds: [{ index: 0, value: -Infinity, color: 'green' }, { index: 1, value: 80, color: 'red' }],
    orientation: VizOrientation.Auto,
};
//# sourceMappingURL=types.js.map