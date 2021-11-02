import { useTheme2 } from '@grafana/ui';
import { STAT, TIMESERIES } from '../utils/constants';
export function useVizHeight(data, pluginId, frameIndex) {
    var theme = useTheme2();
    if (pluginId === TIMESERIES || pluginId === STAT || dataIsEmpty(data)) {
        return 200;
    }
    var values = data.series[frameIndex].fields[0].values.length;
    var rowHeight = theme.spacing.gridSize * 5;
    /*
     Calculate how if we can make  the table smaller than 200px
     for when we only have 1-2 values
     The extra rowHeight is to accommodate the header.
    */
    var tableHeight = values * rowHeight + rowHeight;
    return tableHeight >= 200 ? 200 : tableHeight;
}
function dataIsEmpty(data) {
    return !data || !data.series[0] || !data.series[0].fields[0] || !data.series[0].fields[0].values;
}
//# sourceMappingURL=useVizHeight.js.map