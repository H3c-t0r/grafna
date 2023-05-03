import { NavModel, NavModelItem } from '@grafana/data';
import { config } from '@grafana/runtime';
import { contextSrv } from 'app/core/services/context_srv';
import { AccessControlAction, FolderDTO } from 'app/types';

export function buildNavModel(folder: FolderDTO, parents = folder.parents): NavModelItem {
  const url = `/dashboards/f/${folder.uid}`;
  const model: NavModelItem = {
    icon: 'folder',
    id: 'manage-folder',
    subTitle: 'Manage folder dashboards and permissions',
    url,
    text: folder.title,
    breadcrumbs: [{ title: 'Dashboards', url: 'dashboards' }],
    children: [
      {
        active: false,
        icon: 'apps',
        id: `folder-dashboards-${folder.uid}`,
        text: 'Dashboards',
        url,
        hideFromBreadcrumbs: true,
      },
    ],
  };

  if (parents && parents.length > 0) {
    const parent = parents[parents.length - 1];
    const remainingParents = parents.slice(0, parents.length - 1);
    model.parentItem = buildNavModel(parent, remainingParents);
  }

  model.children!.push({
    active: false,
    icon: 'library-panel',
    id: `folder-library-panels-${folder.uid}`,
    text: 'Panels',
    url: `${url}/library-panels`,
    hideFromBreadcrumbs: true,
  });

  if (contextSrv.hasPermission(AccessControlAction.AlertingRuleRead) && config.unifiedAlertingEnabled) {
    model.children!.push({
      active: false,
      icon: 'bell',
      id: `folder-alerting-${folder.uid}`,
      text: 'Alert rules',
      url: `${url}/alerting`,
      hideFromBreadcrumbs: true,
    });
  }

  if (folder.canAdmin) {
    model.children!.push({
      active: false,
      icon: 'lock',
      id: `folder-permissions-${folder.uid}`,
      text: 'Permissions',
      url: `${url}/permissions`,
      hideFromBreadcrumbs: true,
    });
  }

  if (folder.canSave) {
    model.children!.push({
      active: false,
      icon: 'cog',
      id: `folder-settings-${folder.uid}`,
      text: 'Settings',
      url: `${url}/settings`,
      hideFromBreadcrumbs: true,
    });
  }

  return model;
}

export function getLoadingNav(tabIndex: number): NavModel {
  const main = buildNavModel({
    created: '',
    createdBy: '',
    hasAcl: false,
    updated: '',
    updatedBy: '',
    id: 1,
    uid: 'loading',
    title: 'Loading',
    url: 'url',
    canSave: true,
    canEdit: true,
    canAdmin: true,
    canDelete: true,
    version: 0,
  });

  main.children![tabIndex].active = true;

  return {
    main: main,
    node: main.children![tabIndex],
  };
}
