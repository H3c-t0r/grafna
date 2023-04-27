// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     public/app/plugins/gen.go
// Using jennies:
//     TSTypesJenny
//     PluginTSTypesJenny
//
// Run 'make gen-cue' from repository root to regenerate.

import * as common from '@grafana/schema';

export const DataQueryModelVersion = Object.freeze([0, 0]);

export type PhlareQueryType = ('metrics' | 'profile' | 'both');

export const defaultPhlareQueryType: PhlareQueryType = 'both';

export interface GrafanaPyroscope extends common.DataQuery {
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

export const defaultGrafanaPyroscope: Partial<GrafanaPyroscope> = {
  groupBy: [],
  labelSelector: '{}',
};
