import { css, cx } from '@emotion/css';
import React, { useEffect } from 'react';
import { useStyles2, FilterInput } from '@grafana/ui';
import { contextSrv } from 'app/core/services/context_srv';
import { AccessControlAction } from 'app/types';
import { useKeyNavigationListener } from '../hooks/useSearchKeyboardSelection';
import { SearchView } from '../page/components/SearchView';
import { getSearchStateManager } from '../state/SearchStateManager';
import { getSearchPlaceholder } from '../tempI18nPhrases';
import { DashboardActions } from './DashboardActions';
export const ManageDashboardsNew = React.memo(({ folder }) => {
    var _a;
    const styles = useStyles2(getStyles);
    // since we don't use "query" from use search... it is not actually loaded from the URL!
    const stateManager = getSearchStateManager();
    const state = stateManager.useState();
    const { onKeyDown, keyboardEvents } = useKeyNavigationListener();
    // TODO: we need to refactor DashboardActions to use folder.uid instead
    const folderUid = folder === null || folder === void 0 ? void 0 : folder.uid;
    const canSave = folder === null || folder === void 0 ? void 0 : folder.canSave;
    const { isEditor } = contextSrv;
    const hasEditPermissionInFolders = folder ? canSave : contextSrv.hasEditPermissionInFolders;
    const canCreateFolders = contextSrv.hasPermission(AccessControlAction.FoldersCreate);
    const canCreateDashboards = folderUid
        ? contextSrv.hasPermissionInMetadata(AccessControlAction.DashboardsCreate, folder)
        : contextSrv.hasPermission(AccessControlAction.DashboardsCreate);
    const viewActions = (folder === undefined && canCreateFolders) || canCreateDashboards;
    useEffect(() => stateManager.initStateFromUrl(folder === null || folder === void 0 ? void 0 : folder.uid), [folder === null || folder === void 0 ? void 0 : folder.uid, stateManager]);
    return (React.createElement(React.Fragment, null,
        React.createElement("div", { className: cx(styles.actionBar, 'page-action-bar') },
            React.createElement("div", { className: cx(styles.inputWrapper, 'gf-form gf-form--grow m-r-2') },
                React.createElement(FilterInput, { value: (_a = state.query) !== null && _a !== void 0 ? _a : '', onChange: (e) => stateManager.onQueryChange(e), onKeyDown: onKeyDown, 
                    // eslint-disable-next-line jsx-a11y/no-autofocus
                    autoFocus: true, spellCheck: false, placeholder: getSearchPlaceholder(state.includePanels), escapeRegex: false, className: styles.searchInput })),
            viewActions && (React.createElement(DashboardActions, { folder: folder, canCreateFolders: canCreateFolders, canCreateDashboards: canCreateDashboards }))),
        React.createElement(SearchView, { showManage: Boolean(isEditor || hasEditPermissionInFolders || canSave), folderDTO: folder, hidePseudoFolders: true, keyboardEvents: keyboardEvents })));
});
ManageDashboardsNew.displayName = 'ManageDashboardsNew';
export default ManageDashboardsNew;
const getStyles = (theme) => ({
    actionBar: css `
    ${theme.breakpoints.down('sm')} {
      flex-wrap: wrap;
    }
  `,
    inputWrapper: css `
    ${theme.breakpoints.down('sm')} {
      margin-right: 0 !important;
    }
  `,
    searchInput: css `
    margin-bottom: 6px;
    min-height: ${theme.spacing(4)};
  `,
    unsupported: css `
    padding: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    font-size: 18px;
  `,
    noResults: css `
    padding: ${theme.v1.spacing.md};
    background: ${theme.v1.colors.bg2};
    font-style: italic;
    margin-top: ${theme.v1.spacing.md};
  `,
});
//# sourceMappingURL=ManageDashboardsNew.js.map