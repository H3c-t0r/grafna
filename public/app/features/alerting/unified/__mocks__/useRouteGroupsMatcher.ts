import { useCallback } from 'react';

import { Labels } from '@grafana/data';

import { AlertmanagerGroup, RouteWithID } from '../../../../plugins/datasource/alertmanager/types';
import { routeGroupsMatcher } from '../routeGroupsMatcher';

export function useRouteGroupsMatcher() {
  const getRouteGroupsMap = useCallback(async (route: RouteWithID, groups: AlertmanagerGroup[]) => {
    return Promise.resolve(routeGroupsMatcher.getRouteGroupsMap(route, groups));
  }, []);

  const matchInstancesToRoute = useCallback(async (rootRoute: RouteWithID, instancesToMatch: Labels[]) => {
    return Promise.resolve(routeGroupsMatcher.matchInstancesToRoute(rootRoute, instancesToMatch));
  }, []);

  return { getRouteGroupsMap, matchInstancesToRoute };
}
