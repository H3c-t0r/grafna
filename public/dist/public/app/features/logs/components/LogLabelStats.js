import { css } from '@emotion/css';
import React, { PureComponent } from 'react';
import { stylesFactory, withTheme2 } from '@grafana/ui';
//Components
import { LogLabelStatsRow } from './LogLabelStatsRow';
const STATS_ROW_LIMIT = 5;
const getStyles = stylesFactory((theme) => {
    return {
        logsStats: css `
      label: logs-stats;
      background: inherit;
      color: ${theme.colors.text.primary};
      word-break: break-all;
      width: fit-content;
      max-width: 100%;
    `,
        logsStatsHeader: css `
      label: logs-stats__header;
      border-bottom: 1px solid ${theme.colors.border.medium};
      display: flex;
    `,
        logsStatsTitle: css `
      label: logs-stats__title;
      font-weight: ${theme.typography.fontWeightMedium};
      padding-right: ${theme.spacing(2)};
      display: inline-block;
      white-space: nowrap;
      text-overflow: ellipsis;
      flex-grow: 1;
    `,
        logsStatsClose: css `
      label: logs-stats__close;
      cursor: pointer;
    `,
        logsStatsBody: css `
      label: logs-stats__body;
      padding: 5px 0px;
    `,
    };
});
class UnThemedLogLabelStats extends PureComponent {
    render() {
        const { label, rowCount, stats, value, theme, isLabel } = this.props;
        const style = getStyles(theme);
        const topRows = stats.slice(0, STATS_ROW_LIMIT);
        let activeRow = topRows.find((row) => row.value === value);
        let otherRows = stats.slice(STATS_ROW_LIMIT);
        const insertActiveRow = !activeRow;
        // Remove active row from other to show extra
        if (insertActiveRow) {
            activeRow = otherRows.find((row) => row.value === value);
            otherRows = otherRows.filter((row) => row.value !== value);
        }
        const otherCount = otherRows.reduce((sum, row) => sum + row.count, 0);
        const topCount = topRows.reduce((sum, row) => sum + row.count, 0);
        const total = topCount + otherCount;
        const otherProportion = otherCount / total;
        return (React.createElement("div", { className: style.logsStats, "data-testid": "logLabelStats" },
            React.createElement("div", { className: style.logsStatsHeader },
                React.createElement("div", { className: style.logsStatsTitle },
                    label,
                    ": ",
                    total,
                    " of ",
                    rowCount,
                    " rows have that ",
                    isLabel ? 'label' : 'field')),
            React.createElement("div", { className: style.logsStatsBody },
                topRows.map((stat) => (React.createElement(LogLabelStatsRow, Object.assign({ key: stat.value }, stat, { active: stat.value === value })))),
                insertActiveRow && activeRow && React.createElement(LogLabelStatsRow, Object.assign({ key: activeRow.value }, activeRow, { active: true })),
                otherCount > 0 && (React.createElement(LogLabelStatsRow, { key: "__OTHERS__", count: otherCount, value: "Other", proportion: otherProportion })))));
    }
}
export const LogLabelStats = withTheme2(UnThemedLogLabelStats);
LogLabelStats.displayName = 'LogLabelStats';
//# sourceMappingURL=LogLabelStats.js.map