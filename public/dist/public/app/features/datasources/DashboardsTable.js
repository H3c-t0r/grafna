import React from 'react';
import { Button, Icon } from '@grafana/ui';
var DashboardsTable = function (_a) {
    var dashboards = _a.dashboards, onImport = _a.onImport, onRemove = _a.onRemove;
    function buttonText(dashboard) {
        return dashboard.revision !== dashboard.importedRevision ? 'Update' : 'Re-import';
    }
    return (React.createElement("table", { className: "filter-table" },
        React.createElement("tbody", null, dashboards.map(function (dashboard, index) {
            return (React.createElement("tr", { key: dashboard.dashboardId + "-" + index },
                React.createElement("td", { className: "width-1" },
                    React.createElement(Icon, { name: "apps" })),
                React.createElement("td", null, dashboard.imported ? (React.createElement("a", { href: dashboard.importedUrl }, dashboard.title)) : (React.createElement("span", null, dashboard.title))),
                React.createElement("td", { style: { textAlign: 'right' } },
                    !dashboard.imported ? (React.createElement(Button, { variant: "secondary", size: "sm", onClick: function () { return onImport(dashboard, false); } }, "Import")) : (React.createElement(Button, { variant: "secondary", size: "sm", onClick: function () { return onImport(dashboard, true); } }, buttonText(dashboard))),
                    dashboard.imported && (React.createElement(Button, { icon: "trash-alt", variant: "destructive", size: "sm", onClick: function () { return onRemove(dashboard); } })))));
        }))));
};
export default DashboardsTable;
//# sourceMappingURL=DashboardsTable.js.map