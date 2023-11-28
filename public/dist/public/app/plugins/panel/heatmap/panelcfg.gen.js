// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     public/app/plugins/gen.go
// Using jennies:
//     TSTypesJenny
//     PluginTSTypesJenny
//
// Run 'make gen-cue' from repository root to regenerate.
import * as ui from '@grafana/schema';
/**
 * Controls the color mode of the heatmap
 */
export var HeatmapColorMode;
(function (HeatmapColorMode) {
    HeatmapColorMode["Opacity"] = "opacity";
    HeatmapColorMode["Scheme"] = "scheme";
})(HeatmapColorMode || (HeatmapColorMode = {}));
/**
 * Controls the color scale of the heatmap
 */
export var HeatmapColorScale;
(function (HeatmapColorScale) {
    HeatmapColorScale["Exponential"] = "exponential";
    HeatmapColorScale["Linear"] = "linear";
})(HeatmapColorScale || (HeatmapColorScale = {}));
export const defaultOptions = {
    calculate: false,
    cellGap: 1,
    cellValues: {},
    color: {
        /**
         * mode:     HeatmapColorMode // TODO: fix after remove when https://github.com/grafana/cuetsy/issues/74 is fixed
         */
        scheme: 'Oranges',
        fill: 'dark-orange',
        /**
         * scale:    HeatmapColorScale // TODO: fix after remove when https://github.com/grafana/cuetsy/issues/74 is fixed
         */
        reverse: false,
        exponent: 0.5,
        steps: 64,
    },
    exemplars: {
        color: 'rgba(255,0,255,0.7)',
    },
    filterValues: {
        le: 1e-09,
    },
    legend: {
        show: true,
    },
    showValue: ui.VisibilityMode.Auto,
    tooltip: {
        show: true,
        yHistogram: false,
    },
};
//# sourceMappingURL=panelcfg.gen.js.map