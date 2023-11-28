import { css } from '@emotion/css';
import React from 'react';
import { Components } from '@grafana/e2e-selectors';
import { Icon, IconButton, ToolbarButton, useStyles2 } from '@grafana/ui';
import { t } from 'app/core/internationalization';
import { HOME_NAV_ID } from 'app/core/reducers/navModel';
import { useSelector } from 'app/types';
import { Breadcrumbs } from '../../Breadcrumbs/Breadcrumbs';
import { buildBreadcrumbs } from '../../Breadcrumbs/utils';
import { TOP_BAR_LEVEL_HEIGHT } from '../types';
import { NavToolbarSeparator } from './NavToolbarSeparator';
export const TOGGLE_BUTTON_ID = 'mega-menu-toggle';
export function NavToolbar({ actions, searchBarHidden, sectionNav, pageNav, onToggleMegaMenu, onToggleSearchBar, onToggleKioskMode, }) {
    const homeNav = useSelector((state) => state.navIndex)[HOME_NAV_ID];
    const styles = useStyles2(getStyles);
    const breadcrumbs = buildBreadcrumbs(sectionNav, pageNav, homeNav);
    return (React.createElement("div", { "data-testid": Components.NavToolbar.container, className: styles.pageToolbar },
        React.createElement("div", { className: styles.menuButton },
            React.createElement(IconButton, { id: TOGGLE_BUTTON_ID, name: "bars", tooltip: t('navigation.toolbar.toggle-menu', 'Toggle menu'), tooltipPlacement: "bottom", size: "xl", onClick: onToggleMegaMenu })),
        React.createElement(Breadcrumbs, { breadcrumbs: breadcrumbs, className: styles.breadcrumbsWrapper }),
        React.createElement("div", { className: styles.actions },
            actions,
            actions && React.createElement(NavToolbarSeparator, null),
            searchBarHidden && (React.createElement(ToolbarButton, { onClick: onToggleKioskMode, narrow: true, title: t('navigation.toolbar.enable-kiosk', 'Enable kiosk mode'), icon: "monitor" })),
            React.createElement(ToolbarButton, { onClick: onToggleSearchBar, narrow: true, title: t('navigation.toolbar.toggle-search-bar', 'Toggle top search bar') },
                React.createElement(Icon, { name: searchBarHidden ? 'angle-down' : 'angle-up', size: "xl" })))));
}
const getStyles = (theme) => {
    return {
        breadcrumbsWrapper: css({
            display: 'flex',
            overflow: 'hidden',
            [theme.breakpoints.down('sm')]: {
                minWidth: '50%',
            },
        }),
        pageToolbar: css({
            height: TOP_BAR_LEVEL_HEIGHT,
            display: 'flex',
            padding: theme.spacing(0, 1, 0, 2),
            alignItems: 'center',
        }),
        menuButton: css({
            display: 'flex',
            alignItems: 'center',
            marginRight: theme.spacing(1),
        }),
        actions: css({
            label: 'NavToolbar-actions',
            display: 'flex',
            alignItems: 'center',
            flexWrap: 'nowrap',
            justifyContent: 'flex-end',
            paddingLeft: theme.spacing(1),
            flexGrow: 1,
            gap: theme.spacing(0.5),
            minWidth: 0,
            '.body-drawer-open &': {
                display: 'none',
            },
        }),
    };
};
//# sourceMappingURL=NavToolbar.js.map