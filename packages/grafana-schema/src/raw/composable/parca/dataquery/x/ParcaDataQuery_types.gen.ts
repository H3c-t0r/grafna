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

export const pluginVersion = "10.0.10";

export type ParcaQueryType = ('metrics' | 'profile' | 'both');

export const defaultParcaQueryType: ParcaQueryType = 'both';

export interface ParcaDataQuery extends common.DataQuery {
  /**
   * Specifies the query label selectors.
   */
  labelSelector: string;
  /**
   * Specifies the type of profile to query.
   */
  profileTypeId: string;
}

export const defaultParcaDataQuery: Partial<ParcaDataQuery> = {
  labelSelector: '{}',
};
