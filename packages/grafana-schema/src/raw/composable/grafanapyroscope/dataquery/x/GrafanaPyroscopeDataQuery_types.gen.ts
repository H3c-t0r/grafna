// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     public/app/plugins/gen.go
// Using jennies:
//     TSTypesJenny
//     LatestMajorsOrXJenny
//     PluginEachMajorJenny
//
// Run 'make gen-cue' from repository root to regenerate.

import * as common from '@grafana/schema';

export const pluginVersion = "10.2.0-pre";

export type PyroscopeQueryType = ('metrics' | 'profile' | 'both');

export const defaultPyroscopeQueryType: PyroscopeQueryType = 'both';

export interface GrafanaPyroscopeDataQuery extends common.DataQuery {
  /**
   * Allows to group the results.
   */
  groupBy: Array<string>;
  /**
   * Specifies the query label selectors.
   */
  labelSelector: string;
  /**
   * Sets the maximum number of nodes in the flamegraph.
   */
  maxNodes?: number;
  /**
   * Specifies the type of profile to query.
   */
  profileTypeId: string;
}

export const defaultGrafanaPyroscopeDataQuery: Partial<GrafanaPyroscopeDataQuery> = {
  groupBy: [],
  labelSelector: '{}',
};
