import { __awaiter } from "tslib";
import { css } from '@emotion/css';
import React, { useEffect, useState } from 'react';
import useAsync from 'react-use/lib/useAsync';
import { Alert, InlineField, InlineFieldRow, useStyles2 } from '@grafana/ui';
import { AdHocFilter } from '../../../features/variables/adhoc/picker/AdHocFilter';
import { getDS } from './utils';
export function ServiceGraphSection({ graphDatasourceUid, query, onChange, }) {
    const styles = useStyles2(getStyles);
    const dsState = useAsync(() => getDS(graphDatasourceUid), [graphDatasourceUid]);
    // Check if service graph metrics are being collected. If not, displays a warning
    const [hasKeys, setHasKeys] = useState(undefined);
    useEffect(() => {
        function fn(ds) {
            return __awaiter(this, void 0, void 0, function* () {
                const keys = yield ds.getTagKeys({
                    filters: [
                        {
                            key: '__name__',
                            operator: '=~',
                            value: 'traces_service_graph_request_server_seconds_sum|traces_service_graph_request_total|traces_service_graph_request_failed_total',
                            condition: '',
                        },
                    ],
                });
                setHasKeys(Boolean(keys.length));
            });
        }
        if (!dsState.loading && dsState.value) {
            fn(dsState.value);
        }
    }, [dsState]);
    if (dsState.loading) {
        return null;
    }
    const ds = dsState.value;
    if (!graphDatasourceUid) {
        return getWarning('No service graph datasource selected', 'Please set up a service graph datasource in the datasource settings', styles);
    }
    if (graphDatasourceUid && !ds) {
        return getWarning('No service graph data found', 'Service graph datasource is configured but the data source no longer exists. Please configure existing data source to use the service graph functionality', styles);
    }
    const filters = queryToFilter(query.serviceMapQuery || '');
    return (React.createElement("div", null,
        React.createElement(InlineFieldRow, null,
            React.createElement(InlineField, { label: "Filter", labelWidth: 14, grow: true },
                React.createElement(AdHocFilter, { datasource: { uid: graphDatasourceUid }, filters: filters, baseFilters: [
                        {
                            key: '__name__',
                            operator: '=~',
                            value: 'traces_service_graph_request_total|traces_spanmetrics_calls_total',
                            condition: '',
                        },
                    ], addFilter: (filter) => {
                        onChange(Object.assign(Object.assign({}, query), { serviceMapQuery: filtersToQuery([...filters, filter]) }));
                    }, removeFilter: (index) => {
                        const newFilters = [...filters];
                        newFilters.splice(index, 1);
                        onChange(Object.assign(Object.assign({}, query), { serviceMapQuery: filtersToQuery(newFilters) }));
                    }, changeFilter: (index, filter) => {
                        const newFilters = [...filters];
                        newFilters.splice(index, 1, filter);
                        onChange(Object.assign(Object.assign({}, query), { serviceMapQuery: filtersToQuery(newFilters) }));
                    } }))),
        hasKeys === false
            ? getWarning('No service graph data found', 'Please ensure that service graph metrics are set up correctly', styles)
            : null));
}
function getWarning(title, description, styles) {
    return (React.createElement(Alert, { title: title, severity: "info", className: styles.alert },
        description,
        " according to the",
        ' ',
        React.createElement("a", { target: "_blank", rel: "noreferrer noopener", href: "https://grafana.com/docs/grafana/latest/datasources/tempo/service-graph/", className: styles.link }, "Tempo documentation"),
        "."));
}
function queryToFilter(query) {
    let match;
    let filters = [];
    const re = /([\w_]+)(=|!=|<|>|=~|!~)"(.*?)"/g;
    while ((match = re.exec(query)) !== null) {
        filters.push({
            key: match[1],
            operator: match[2],
            value: match[3],
            condition: '',
        });
    }
    return filters;
}
function filtersToQuery(filters) {
    return `{${filters.map((f) => `${f.key}${f.operator}"${f.value}"`).join(',')}}`;
}
const getStyles = (theme) => ({
    alert: css `
    max-width: 75ch;
    margin-top: ${theme.spacing(2)};
  `,
    link: css `
    color: ${theme.colors.text.link};
    text-decoration: underline;
  `,
});
//# sourceMappingURL=ServiceGraphSection.js.map