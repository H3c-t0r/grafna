// THIS FILE IS GENERATED. EDITING IS FUTILE.
//
// Generated by:
//     kinds/gen.go
// Using jennies:
//     TSTypesJenny
//     LatestMajorsOrXJenny
//
// Run 'make gen-cue' from repository root to regenerate.

export interface NavLink {
  URL?: string;
  id: string;
  target?: string;
  text?: string;
}

export interface NavbarPreference {
  savedItems?: Array<NavLink>;
}

export const defaultNavbarPreference: Partial<NavbarPreference> = {
  savedItems: [],
};

export interface QueryHistoryPreference {
  /**
   * one of: '' | 'query' | 'starred';
   */
  homeTab?: string;
}

export interface Preferences {
  /**
   * UID for the home dashboard
   */
  homeDashboardUID?: string;
  /**
   * Selected language (beta)
   */
  language?: string;
  navbar?: NavbarPreference;
  queryHistory?: QueryHistoryPreference;
  /**
   * light, dark, empty is default
   */
  theme?: string;
  /**
   * The timezone selection
   * Would be nice it this used:
   * import { TimeZone } from '@grafana/data';
   */
  timezone?: string;
  /**
   * day of the week (sunday, monday, etc)
   */
  weekStart?: string;
}
