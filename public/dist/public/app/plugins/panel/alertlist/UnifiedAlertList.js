import { css } from '@emotion/css';
import { sortBy } from 'lodash';
import React, { useEffect, useMemo } from 'react';
import { useEffectOnce, useToggle } from 'react-use';
import { TimeRangeUpdatedEvent } from '@grafana/runtime';
import { Alert, BigValue, BigValueGraphMode, BigValueJustifyMode, BigValueTextMode, CustomScrollbar, LoadingPlaceholder, useStyles2, } from '@grafana/ui';
import { config } from 'app/core/config';
import { contextSrv } from 'app/core/services/context_srv';
import alertDef from 'app/features/alerting/state/alertDef';
import { alertRuleApi } from 'app/features/alerting/unified/api/alertRuleApi';
import { INSTANCES_DISPLAY_LIMIT } from 'app/features/alerting/unified/components/rules/RuleDetails';
import { useCombinedRuleNamespaces } from 'app/features/alerting/unified/hooks/useCombinedRuleNamespaces';
import { useUnifiedAlertingSelector } from 'app/features/alerting/unified/hooks/useUnifiedAlertingSelector';
import { fetchAllPromAndRulerRulesAction, fetchPromAndRulerRulesAction, } from 'app/features/alerting/unified/state/actions';
import { parseMatchers } from 'app/features/alerting/unified/utils/alertmanager';
import { Annotation } from 'app/features/alerting/unified/utils/constants';
import { GRAFANA_DATASOURCE_NAME, GRAFANA_RULES_SOURCE_NAME } from 'app/features/alerting/unified/utils/datasource';
import { isAsyncRequestMapSlicePartiallyDispatched, isAsyncRequestMapSlicePartiallyFulfilled, isAsyncRequestMapSlicePending, } from 'app/features/alerting/unified/utils/redux';
import { flattenCombinedRules, getFirstActiveAt } from 'app/features/alerting/unified/utils/rules';
import { getDashboardSrv } from 'app/features/dashboard/services/DashboardSrv';
import { AccessControlAction, useDispatch } from 'app/types';
import { PromAlertingRuleState } from 'app/types/unified-alerting-dto';
import { getAlertingRule } from '../../../features/alerting/unified/utils/rules';
import { GroupMode, SortOrder, ViewMode } from './types';
import GroupedModeView from './unified-alerting/GroupedView';
import UngroupedModeView from './unified-alerting/UngroupedView';
import { filterAlerts } from './util';
function getStateList(state) {
    const reducer = (list, [stateKey, value]) => {
        if (Boolean(value)) {
            return [...list, stateKey];
        }
        else {
            return list;
        }
    };
    return Object.entries(state).reduce(reducer, []);
}
const fetchPromAndRuler = ({ dispatch, limitInstances, matcherList, dataSourceName, stateList, }) => {
    if (dataSourceName) {
        dispatch(fetchPromAndRulerRulesAction({
            rulesSourceName: dataSourceName,
            limitAlerts: limitInstances ? INSTANCES_DISPLAY_LIMIT : undefined,
            matcher: matcherList,
            state: stateList,
        }));
    }
    else {
        dispatch(fetchAllPromAndRulerRulesAction(false, {
            limitAlerts: limitInstances ? INSTANCES_DISPLAY_LIMIT : undefined,
            matcher: matcherList,
            state: stateList,
        }));
    }
};
export function UnifiedAlertList(props) {
    const dispatch = useDispatch();
    const [limitInstances, toggleLimit] = useToggle(true);
    const { usePrometheusRulesByNamespaceQuery } = alertRuleApi;
    const promRulesRequests = useUnifiedAlertingSelector((state) => state.promRules);
    const rulerRulesRequests = useUnifiedAlertingSelector((state) => state.rulerRules);
    const somePromRulesDispatched = isAsyncRequestMapSlicePartiallyDispatched(promRulesRequests);
    const hideViewRuleLinkText = props.width < 320;
    // backwards compat for "Inactive" state filter
    useEffect(() => {
        if (props.options.stateFilter.inactive === true) {
            props.options.stateFilter.normal = true; // enable the normal filter
        }
        props.options.stateFilter.inactive = undefined; // now disable inactive
    }, [props.options.stateFilter]);
    let dashboard = undefined;
    useEffectOnce(() => {
        dashboard = getDashboardSrv().getCurrent();
    });
    const stateList = useMemo(() => getStateList(props.options.stateFilter), [props.options.stateFilter]);
    const { options, replaceVariables } = props;
    const dataSourceName = options.datasource === GRAFANA_DATASOURCE_NAME ? GRAFANA_RULES_SOURCE_NAME : options.datasource;
    const parsedOptions = Object.assign(Object.assign({}, props.options), { alertName: replaceVariables(options.alertName), alertInstanceLabelFilter: replaceVariables(options.alertInstanceLabelFilter) });
    const matcherList = useMemo(() => parseMatchers(parsedOptions.alertInstanceLabelFilter), [parsedOptions.alertInstanceLabelFilter]);
    // If the datasource is not defined we should NOT skip the query
    // Undefined dataSourceName means that there is no datasource filter applied and we should fetch all the rules
    const shouldFetchGrafanaRules = !dataSourceName || dataSourceName === GRAFANA_RULES_SOURCE_NAME;
    //For grafana managed rules, get the result using RTK Query to avoid the need of using the redux store
    //See https://github.com/grafana/grafana/pull/70482
    const { currentData: grafanaPromRules = [], isLoading: grafanaRulesLoading, refetch: refetchGrafanaPromRules, } = usePrometheusRulesByNamespaceQuery({
        limitAlerts: limitInstances ? INSTANCES_DISPLAY_LIMIT : undefined,
        matcher: matcherList,
        state: stateList,
    }, { skip: !shouldFetchGrafanaRules });
    useEffect(() => {
        //we need promRules and rulerRules for getting the uid when creating the alert link in panel in case of being a rulerRule.
        if (!promRulesRequests.loading) {
            fetchPromAndRuler({ dispatch, limitInstances, matcherList, dataSourceName, stateList });
        }
        const sub = dashboard === null || dashboard === void 0 ? void 0 : dashboard.events.subscribe(TimeRangeUpdatedEvent, () => {
            if (shouldFetchGrafanaRules) {
                refetchGrafanaPromRules();
            }
            if (!dataSourceName || dataSourceName !== GRAFANA_RULES_SOURCE_NAME) {
                fetchPromAndRuler({ dispatch, limitInstances, matcherList, dataSourceName, stateList });
            }
        });
        return () => {
            sub === null || sub === void 0 ? void 0 : sub.unsubscribe();
        };
    }, [
        dispatch,
        dashboard,
        matcherList,
        stateList,
        limitInstances,
        dataSourceName,
        refetchGrafanaPromRules,
        shouldFetchGrafanaRules,
        promRulesRequests.loading,
    ]);
    const handleInstancesLimit = (limit) => {
        if (limit) {
            fetchPromAndRuler({ dispatch, limitInstances, matcherList, dataSourceName, stateList });
            toggleLimit(true);
        }
        else {
            fetchPromAndRuler({ dispatch, limitInstances: false, matcherList, dataSourceName, stateList });
            toggleLimit(false);
        }
    };
    const combinedRules = useCombinedRuleNamespaces(undefined, grafanaPromRules);
    const someRulerRulesDispatched = isAsyncRequestMapSlicePartiallyDispatched(rulerRulesRequests);
    const haveResults = isAsyncRequestMapSlicePartiallyFulfilled(promRulesRequests);
    const dispatched = somePromRulesDispatched || someRulerRulesDispatched;
    const loading = isAsyncRequestMapSlicePending(promRulesRequests);
    const styles = useStyles2(getStyles);
    const flattenedCombinedRules = flattenCombinedRules(combinedRules);
    const order = props.options.sortOrder;
    const rules = useMemo(() => filterRules(props, sortRules(order, flattenedCombinedRules)), [flattenedCombinedRules, order, props]);
    const noAlertsMessage = rules.length === 0 ? 'No alerts matching filters' : undefined;
    if (!contextSrv.hasPermission(AccessControlAction.AlertingRuleRead) &&
        !contextSrv.hasPermission(AccessControlAction.AlertingRuleExternalRead)) {
        return (React.createElement(Alert, { title: "Permission required" }, "Sorry, you do not have the required permissions to read alert rules"));
    }
    return (React.createElement(CustomScrollbar, { autoHeightMin: "100%", autoHeightMax: "100%" },
        React.createElement("div", { className: styles.container },
            (grafanaRulesLoading || (dispatched && loading && !haveResults)) && React.createElement(LoadingPlaceholder, { text: "Loading..." }),
            noAlertsMessage && React.createElement("div", { className: styles.noAlertsMessage }, noAlertsMessage),
            React.createElement("section", null,
                props.options.viewMode === ViewMode.Stat && haveResults && (React.createElement(BigValue, { width: props.width, height: props.height, graphMode: BigValueGraphMode.None, textMode: BigValueTextMode.Auto, justifyMode: BigValueJustifyMode.Auto, theme: config.theme2, value: { text: `${rules.length}`, numeric: rules.length } })),
                props.options.viewMode === ViewMode.List && props.options.groupMode === GroupMode.Custom && haveResults && (React.createElement(GroupedModeView, { rules: rules, options: parsedOptions })),
                props.options.viewMode === ViewMode.List && props.options.groupMode === GroupMode.Default && haveResults && (React.createElement(UngroupedModeView, { rules: rules, options: parsedOptions, handleInstancesLimit: handleInstancesLimit, limitInstances: limitInstances, hideViewRuleLinkText: hideViewRuleLinkText }))))));
}
function sortRules(sortOrder, rules) {
    if (sortOrder === SortOrder.Importance) {
        // @ts-ignore
        return sortBy(rules, (rule) => alertDef.alertStateSortScore[rule.state]);
    }
    else if (sortOrder === SortOrder.TimeAsc) {
        return sortBy(rules, (rule) => {
            var _a;
            //at this point rules are all AlertingRule, this check is only needed for Typescript checks
            const alertingRule = (_a = getAlertingRule(rule)) !== null && _a !== void 0 ? _a : undefined;
            return getFirstActiveAt(alertingRule) || new Date();
        });
    }
    else if (sortOrder === SortOrder.TimeDesc) {
        return sortBy(rules, (rule) => {
            var _a;
            //at this point rules are all AlertingRule, this check is only needed for Typescript checks
            const alertingRule = (_a = getAlertingRule(rule)) !== null && _a !== void 0 ? _a : undefined;
            return getFirstActiveAt(alertingRule) || new Date();
        }).reverse();
    }
    const result = sortBy(rules, (rule) => rule.name.toLowerCase());
    if (sortOrder === SortOrder.AlphaDesc) {
        result.reverse();
    }
    return result;
}
function filterRules(props, rules) {
    var _a;
    const { options, replaceVariables } = props;
    let filteredRules = [...rules];
    if (options.dashboardAlerts) {
        const dashboardUid = (_a = getDashboardSrv().getCurrent()) === null || _a === void 0 ? void 0 : _a.uid;
        filteredRules = filteredRules.filter(({ annotations = {} }) => Object.entries(annotations).some(([key, value]) => key === Annotation.dashboardUID && value === dashboardUid));
    }
    if (options.alertName) {
        const replacedName = replaceVariables(options.alertName);
        filteredRules = filteredRules.filter(({ name }) => name.toLocaleLowerCase().includes(replacedName.toLocaleLowerCase()));
    }
    filteredRules = filteredRules.filter((rule) => {
        const alertingRule = getAlertingRule(rule);
        if (!alertingRule) {
            return false;
        }
        return ((options.stateFilter.firing && alertingRule.state === PromAlertingRuleState.Firing) ||
            (options.stateFilter.pending && alertingRule.state === PromAlertingRuleState.Pending) ||
            (options.stateFilter.normal && alertingRule.state === PromAlertingRuleState.Inactive));
    });
    if (options.folder) {
        filteredRules = filteredRules.filter((rule) => {
            return rule.namespaceName === options.folder.title;
        });
    }
    if (options.datasource) {
        const isGrafanaDS = options.datasource === GRAFANA_DATASOURCE_NAME;
        filteredRules = filteredRules.filter(isGrafanaDS
            ? ({ dataSourceName }) => dataSourceName === GRAFANA_RULES_SOURCE_NAME
            : ({ dataSourceName }) => dataSourceName === options.datasource);
    }
    // Remove rules having 0 instances
    // AlertInstances filters instances and we need to prevent situation
    // when we display a rule with 0 instances
    filteredRules = filteredRules.reduce((rules, rule) => {
        var _a;
        const alertingRule = getAlertingRule(rule);
        const filteredAlerts = alertingRule
            ? filterAlerts({
                stateFilter: options.stateFilter,
                alertInstanceLabelFilter: replaceVariables(options.alertInstanceLabelFilter),
            }, (_a = alertingRule.alerts) !== null && _a !== void 0 ? _a : [])
            : [];
        if (filteredAlerts.length) {
            // We intentionally don't set alerts to filteredAlerts
            // because later we couldn't display that some alerts are hidden (ref AlertInstances filtering)
            rules.push(rule);
        }
        return rules;
    }, []);
    return filteredRules;
}
export const getStyles = (theme) => ({
    cardContainer: css `
    padding: ${theme.spacing(0.5)} 0 ${theme.spacing(0.25)} 0;
    line-height: ${theme.typography.body.lineHeight};
    margin-bottom: 0px;
  `,
    container: css `
    overflow-y: auto;
    height: 100%;
  `,
    alertRuleList: css `
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    list-style-type: none;
  `,
    alertRuleItem: css `
    display: flex;
    align-items: center;
    width: 100%;
    height: 100%;
    background: ${theme.colors.background.secondary};
    padding: ${theme.spacing(0.5)} ${theme.spacing(1)};
    border-radius: ${theme.shape.radius.default};
    margin-bottom: ${theme.spacing(0.5)};

    gap: ${theme.spacing(2)};
  `,
    alertName: css `
    font-size: ${theme.typography.h6.fontSize};
    font-weight: ${theme.typography.fontWeightBold};

    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  `,
    alertNameWrapper: css `
    display: flex;
    flex: 1;
    flex-wrap: nowrap;
    flex-direction: column;

    min-width: 100px;
  `,
    alertLabels: css `
    > * {
      margin-right: ${theme.spacing(0.5)};
    }
  `,
    alertDuration: css `
    font-size: ${theme.typography.bodySmall.fontSize};
  `,
    alertRuleItemText: css `
    font-weight: ${theme.typography.fontWeightBold};
    font-size: ${theme.typography.bodySmall.fontSize};
    margin: 0;
  `,
    alertRuleItemTime: css `
    color: ${theme.colors.text.secondary};
    font-weight: normal;
    white-space: nowrap;
  `,
    alertRuleItemInfo: css `
    font-weight: normal;
    flex-grow: 2;
    display: flex;
    align-items: flex-end;
  `,
    noAlertsMessage: css `
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
  `,
    alertIcon: css `
    margin-right: ${theme.spacing(0.5)};
  `,
    instanceDetails: css `
    min-width: 1px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  `,
    customGroupDetails: css `
    margin-bottom: ${theme.spacing(0.5)};
  `,
    link: css `
    word-break: break-all;
    color: ${theme.colors.primary.text};
    display: flex;
    align-items: center;
    gap: ${theme.spacing(1)};
  `,
    hidden: css `
    display: none;
  `,
});
//# sourceMappingURL=UnifiedAlertList.js.map