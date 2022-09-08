import { Dispatch } from '@reduxjs/toolkit';
import { Location } from 'history';

import { locationUtil, NavModelItem, NavSection } from '@grafana/data';
import { reportInteraction } from '@grafana/runtime';
import { config } from 'app/core/config';
import { updateMenuTree } from 'app/core/reducers/navBarTree';
import { contextSrv } from 'app/core/services/context_srv';

import { ShowModalReactEvent } from '../../../types/events';
import appEvents from '../../app_events';
import { getFooterLinks } from '../Footer/Footer';
import { HelpModal } from '../help/HelpModal';

import { PMM_ADD_INSTANCE_PAGE } from './constants';

export const SEARCH_ITEM_ID = 'search';
export const NAV_MENU_PORTAL_CONTAINER_ID = 'navbar-menu-portal-container';

export const getNavMenuPortalContainer = () => document.getElementById(NAV_MENU_PORTAL_CONTAINER_ID) ?? document.body;

const DIVIDER = {
  id: 'divider',
  text: 'Divider',
  divider: true,
  hideFromTabs: true,
};

export const getForcedLoginUrl = (url: string) => {
  const queryParams = new URLSearchParams(url.split('?')[1]);
  queryParams.append('forceLogin', 'true');

  return `${config.appSubUrl}${url.split('?')[0]}?${queryParams.toString()}`;
};

export const enrichConfigItems = (
  items: NavModelItem[],
  location: Location<unknown>,
  toggleOrgSwitcher: () => void
) => {
  const { isSignedIn, user } = contextSrv;
  const onOpenShortcuts = () => {
    appEvents.publish(new ShowModalReactEvent({ component: HelpModal }));
  };

  if (user && user.orgCount > 1) {
    const profileNode = items.find((bottomNavItem) => bottomNavItem.id === 'profile');
    if (profileNode) {
      profileNode.showOrgSwitcher = true;
      profileNode.subTitle = `Current Org.: ${user?.orgName}`;
    }
  }

  if (!isSignedIn) {
    const forcedLoginUrl = getForcedLoginUrl(location.pathname + location.search);

    items.unshift({
      icon: 'signout',
      id: 'signin',
      section: NavSection.Config,
      target: '_self',
      text: 'Sign in',
      url: forcedLoginUrl,
    });
  }

  items.forEach((link, index) => {
    let menuItems = link.children || [];

    if (link.id === 'help') {
      link.children = [
        ...getFooterLinks(),
        {
          id: 'keyboard-shortcuts',
          text: 'Keyboard shortcuts',
          icon: 'keyboard',
          onClick: onOpenShortcuts,
        },
      ];
    }

    if (link.showOrgSwitcher) {
      link.children = [
        ...menuItems,
        {
          id: 'switch-organization',
          text: 'Switch organization',
          icon: 'arrow-random',
          onClick: toggleOrgSwitcher,
        },
      ];
    }
  });
  return items;
};

export const enrichClickWith = (item: NavModelItem, withFc: (item: NavModelItem) => void) => {
  const onClick = item.onClick;
  item.onClick = () => {
    withFc(item);

    onClick?.();
  };
  if (item.children) {
    item.children = item.children.map((item) => enrichClickWith(item, withFc));
  }
  return item;
};

export const enrichWithInteractionTracking = (item: NavModelItem, expandedState: boolean) =>
  enrichClickWith(item, (item) => {
    reportInteraction('grafana_navigation_item_clicked', {
      path: item.url ?? item.id,
      state: expandedState ? 'expanded' : 'collapsed',
    });
  });

// @Percona
export const enrichWithClickDispatch = (item: NavModelItem, dispatch: Dispatch, dispatchOffset: number) =>
  enrichClickWith(item, (link) => {
    // let the animation play out, dispatch action after that
    setTimeout(() => {
      if (link.id) {
        dispatch(updateMenuTree({ id: link.id, active: true }));
      }
    }, dispatchOffset);
  });

export const isMatchOrChildMatch = (itemToCheck: NavModelItem, searchItem?: NavModelItem): boolean => {
  return Boolean(itemToCheck === searchItem || itemToCheck.children?.some((child) => child === searchItem));
};

// @Percona
export const isMatchOrInnerMatch = (itemToCheck: NavModelItem, searchItem?: NavModelItem): boolean => {
  return Boolean(
    itemToCheck === searchItem || itemToCheck.children?.some((child) => isMatchOrInnerMatch(child, searchItem))
  );
};

const stripQueryParams = (url?: string) => {
  return url?.split('?')[0] ?? '';
};

const isBetterMatch = (newMatch: NavModelItem, currentMatch?: NavModelItem) => {
  const currentMatchUrl = stripQueryParams(currentMatch?.url);
  const newMatchUrl = stripQueryParams(newMatch.url);
  return newMatchUrl && newMatchUrl.length > currentMatchUrl?.length;
};

export const getActiveItem = (
  navTree: NavModelItem[],
  pathname: string,
  currentBestMatch?: NavModelItem
): NavModelItem | undefined => {
  const dashboardLinkMatch = '/dashboards';

  for (const link of navTree) {
    const linkWithoutParams = stripQueryParams(link.url);
    const linkPathname = locationUtil.stripBaseFromUrl(linkWithoutParams);
    if (linkPathname) {
      if (linkPathname === pathname) {
        // exact match
        currentBestMatch = link;
        break;
      } else if (linkPathname !== '/' && pathname.startsWith(linkPathname)) {
        // partial match
        if (isBetterMatch(link, currentBestMatch)) {
          currentBestMatch = link;
        }
      } else if (linkPathname === '/alerting/list' && pathname.startsWith('/alerting/notification/')) {
        // alert channel match
        // TODO refactor routes such that we don't need this custom logic
        currentBestMatch = link;
        break;
      } else if (linkPathname === dashboardLinkMatch && pathname.startsWith('/d/')) {
        // dashboard match
        // TODO refactor routes such that we don't need this custom logic
        if (isBetterMatch(link, currentBestMatch)) {
          currentBestMatch = link;
        }
      }
    }
    if (link.children) {
      currentBestMatch = getActiveItem(link.children, pathname, currentBestMatch);
    }
    if (stripQueryParams(currentBestMatch?.url) === pathname) {
      return currentBestMatch;
    }
  }
  return currentBestMatch;
};

export const isSearchActive = (location: Location<unknown>) => {
  const query = new URLSearchParams(location.search);
  return query.get('search') === 'open';
};

export function getNavModelItemKey(item: NavModelItem) {
  return item.id ?? item.text;
}

export const buildIntegratedAlertingMenuItem = (mainLinks: NavModelItem[]): NavModelItem[] => {
  const integratedAlertingLink = {
    id: 'integrated-alerting',
    text: 'Integrated Alerting',
    icon: 'bell',
    url: `${config.appSubUrl}/integrated-alerting`,
  };

  const alertingIndex = mainLinks.findIndex(({ id }) => id === 'alerting');

  if (alertingIndex === -1) {
    mainLinks.push({
      id: 'alerting',
      text: 'Alerting',
      icon: 'bell',
      url: `${config.appSubUrl}/integrated-alerting/alerts`,
      subTitle: 'Alert rules & notifications',
      children: [integratedAlertingLink],
    });
  } else {
    mainLinks[alertingIndex].children?.unshift(integratedAlertingLink, DIVIDER);
  }

  return mainLinks;
};

export const buildInventoryAndSettings = (mainLinks: NavModelItem[]): NavModelItem[] => {
  const inventoryLink: NavModelItem = {
    id: 'inventory',
    icon: 'percona-inventory',
    text: 'PMM Inventory',
    url: `${config.appSubUrl}/inventory`,
    children: [
      {
        id: 'inventory-list',
        url: `${config.appSubUrl}/inventory`,
        icon: 'percona-inventory',
        text: 'Inventory List',
        hideFromTabs: true,
        children: [
          {
            id: 'inventory-services',
            text: 'Services',
            url: `${config.appSubUrl}/inventory/services`,
            hideFromMenu: true,
          },
          {
            id: 'inventory-agents',
            text: 'Agents',
            url: `${config.appSubUrl}/inventory/agents`,
            hideFromMenu: true,
          },
          {
            id: 'inventory-nodes',
            text: 'Nodes',
            url: `${config.appSubUrl}/inventory/nodes`,
            hideFromMenu: true,
          },
        ],
      },
    ],
  };
  const settingsLink = {
    id: 'settings',
    icon: 'percona-setting',
    text: 'PMM Settings',
    url: `${config.appSubUrl}/settings`,
  };
  const configNode = mainLinks.find((link) => link.id === 'cfg');

  if (!configNode) {
    mainLinks.push({
      id: 'cfg',
      text: 'Configuration',
      icon: 'cog',
      url: `${config.appSubUrl}/inventory`,
      subTitle: 'Configuration',
      children: [inventoryLink, settingsLink, DIVIDER, PMM_ADD_INSTANCE_PAGE],
    });
  } else {
    if (!configNode.children) {
      configNode.children = [];
    }
    configNode.url = `${config.appSubUrl}/inventory`;
    configNode.children.unshift(PMM_ADD_INSTANCE_PAGE, DIVIDER, inventoryLink, settingsLink);
  }

  return mainLinks;
};
