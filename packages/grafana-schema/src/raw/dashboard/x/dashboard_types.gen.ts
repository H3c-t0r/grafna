// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     kinds/gen.go
// Using jennies:
//     TSTypesJenny
//     LatestMajorsOrXJenny
//
// Run 'make gen-cue' from repository root to regenerate.

/**
 * TODO docs
 */
export interface AnnotationTarget {
  limit: number;
  matchAny: boolean;
  tags: Array<string>;
  type: string;
}

export const defaultAnnotationTarget: Partial<AnnotationTarget> = {
  tags: [],
};

/**
 * TODO docs
 * FROM: AnnotationQuery in grafana-data/src/types/annotations.ts
 */
export interface AnnotationQuery {
  builtIn: number;
  /**
   * Datasource to use for annotation.
   */
  datasource: {
    type?: string;
    uid?: string;
  };
  /**
   * Whether annotation is enabled.
   */
  enable: boolean;
  /**
   * Whether to hide annotation.
   */
  hide?: boolean;
  /**
   * Annotation icon color.
   */
  iconColor?: string;
  /**
   * Name of annotation.
   */
  name?: string;
  /**
   * Query for annotation data.
   */
  rawQuery?: string;
  showIn: number;
  target?: AnnotationTarget;
  type: string;
}

export const defaultAnnotationQuery: Partial<AnnotationQuery> = {
  builtIn: 0,
  enable: true,
  hide: false,
  showIn: 0,
  type: 'dashboard',
};

/**
 * FROM: packages/grafana-data/src/types/templateVars.ts
 * TODO docs
 * TODO what about what's in public/app/features/types.ts?
 * TODO there appear to be a lot of different kinds of [template] vars here? if so need a disjunction
 */
export interface VariableModel {
  datasource?: DataSourceRef;
  description?: string;
  error?: Record<string, unknown>;
  global: boolean;
  hide: VariableHide;
  id: string;
  index: number;
  label?: string;
  name: string;
  /**
   * TODO: Move this into a separated QueryVariableModel type
   */
  query?: (string | Record<string, unknown>);
  rootStateKey?: string;
  skipUrlSync: boolean;
  state: LoadingState;
  type: VariableType;
}

export const defaultVariableModel: Partial<VariableModel> = {
  global: false,
  id: '00000000-0000-0000-0000-000000000000',
  index: -1,
  skipUrlSync: false,
};

export enum VariableHide {
  dontHide = 0,
  hideLabel = 1,
  hideVariable = 2,
}

export enum LoadingState {
  Done = 'Done',
  Error = 'Error',
  Loading = 'Loading',
  NotStarted = 'NotStarted',
  Streaming = 'Streaming',
}

/**
 * Ref to a DataSource instance
 */
export interface DataSourceRef {
  /**
   * The plugin type-id
   */
  type?: string;
  /**
   * Specific datasource instance
   */
  uid?: string;
}

/**
 * FROM public/app/features/dashboard/state/DashboardModels.ts - ish
 * TODO docs
 */
export interface DashboardLink {
  asDropdown: boolean;
  icon: string;
  includeVars: boolean;
  keepTime: boolean;
  tags: Array<string>;
  targetBlank: boolean;
  title: string;
  tooltip: string;
  type: DashboardLinkType;
  url: string;
}

export const defaultDashboardLink: Partial<DashboardLink> = {
  asDropdown: false,
  includeVars: false,
  keepTime: false,
  tags: [],
  targetBlank: false,
};

/**
 * TODO docs
 */
export type DashboardLinkType = ('link' | 'dashboards');

/**
 * FROM: packages/grafana-data/src/types/templateVars.ts
 * TODO docs
 * TODO this implies some wider pattern/discriminated union, probably?
 */
export type VariableType = ('query' | 'adhoc' | 'constant' | 'datasource' | 'interval' | 'textbox' | 'custom' | 'system');

/**
 * TODO docs
 */
export enum FieldColorModeId {
  ContinuousGrYlRd = 'continuous-GrYlRd',
  Fixed = 'fixed',
  PaletteClassic = 'palette-classic',
  PaletteSaturated = 'palette-saturated',
  Thresholds = 'thresholds',
}

/**
 * TODO docs
 */
export type FieldColorSeriesByMode = ('min' | 'max' | 'last');

/**
 * TODO docs
 */
export interface FieldColor {
  /**
   * Stores the fixed color value if mode is fixed
   */
  fixedColor?: string;
  /**
   * The main color scheme mode
   */
  mode: FieldColorModeId;
  /**
   * Some visualizations need to know how to assign a series color from by value color schemes
   */
  seriesBy?: FieldColorSeriesByMode;
}

export interface GridPos {
  /**
   * Panel
   */
  h: number;
  /**
   * true if fixed
   */
  static?: boolean;
  /**
   * Panel
   */
  w: number;
  /**
   * Panel x
   */
  x: number;
  /**
   * Panel y
   */
  y: number;
}

export const defaultGridPos: Partial<GridPos> = {
  h: 9,
  w: 12,
  x: 0,
  y: 0,
};

/**
 * TODO docs
 */
export interface Threshold {
  /**
   * TODO docs
   */
  color: string;
  /**
   * TODO docs
   * TODO are the values here enumerable into a disjunction?
   * Some seem to be listed in typescript comment
   */
  state?: string;
  /**
   * TODO docs
   * FIXME the corresponding typescript field is required/non-optional, but nulls currently appear here when serializing -Infinity to JSON
   */
  value?: number;
}

export enum ThresholdsMode {
  Absolute = 'absolute',
  Percentage = 'percentage',
}

export interface ThresholdsConfig {
  mode: ThresholdsMode;
  /**
   * Must be sorted by 'value', first value is always -Infinity
   */
  steps: Array<Threshold>;
}

export const defaultThresholdsConfig: Partial<ThresholdsConfig> = {
  steps: [],
};

/**
 * TODO docs
 */
export type ValueMapping = (ValueMap | RangeMap | RegexMap | SpecialValueMap);

/**
 * TODO docs
 */
export enum MappingType {
  RangeToText = 'range',
  RegexToText = 'regex',
  SpecialValue = 'special',
  ValueToText = 'value',
}

/**
 * TODO docs
 */
export interface ValueMap {
  options: Record<string, ValueMappingResult>;
  type: MappingType.ValueToText;
}

/**
 * TODO docs
 */
export interface RangeMap {
  options: {
    /**
     * to and from are `number | null` in current ts, really not sure what to do
     */
    from: number;
    to: number;
    result: ValueMappingResult;
  };
  type: MappingType.RangeToText;
}

/**
 * TODO docs
 */
export interface RegexMap {
  options: {
    pattern: string;
    result: ValueMappingResult;
  };
  type: MappingType.RegexToText;
}

/**
 * TODO docs
 */
export interface SpecialValueMap {
  options: {
    match: ('true' | 'false');
    pattern: string;
    result: ValueMappingResult;
  };
  type: MappingType.SpecialValue;
}

/**
 * TODO docs
 */
export enum SpecialValueMatch {
  Empty = 'empty',
  False = 'false',
  NaN = 'nan',
  Null = 'null',
  NullAndNan = 'null+nan',
  True = 'true',
}

/**
 * TODO docs
 */
export interface ValueMappingResult {
  color?: string;
  icon?: string;
  index?: number;
  text?: string;
}

/**
 * TODO docs
 */
export interface DataTransformerConfig {
  /**
   * Disabled transformations are skipped
   */
  disabled?: boolean;
  /**
   * Optional frame matcher.  When missing it will be applied to all results
   */
  filter?: MatcherConfig;
  /**
   * Unique identifier of transformer
   */
  id: string;
  /**
   * Options to be passed to the transformer
   * Valid options depend on the transformer id
   */
  options: unknown;
}

/**
 * 0 for no shared crosshair or tooltip (default).
 * 1 for shared crosshair.
 * 2 for shared crosshair AND shared tooltip.
 */
export enum DashboardCursorSync {
  Crosshair = 1,
  Off = 0,
  Tooltip = 2,
}

export const defaultDashboardCursorSync: DashboardCursorSync = DashboardCursorSync.Off;

/**
 * Dashboard panels. Panels are canonically defined inline
 * because they share a version timeline with the dashboard
 * schema; they do not evolve independently.
 */
export interface Panel {
  /**
   * The datasource used in all targets.
   */
  datasource?: {
    type?: string;
    uid?: string;
  };
  /**
   * Description.
   */
  description?: string;
  fieldConfig: FieldConfigSource;
  /**
   * Grid position.
   */
  gridPos?: GridPos;
  /**
   * TODO docs
   */
  id?: number;
  /**
   * TODO docs
   * TODO tighter constraint
   */
  interval?: string;
  /**
   * Dynamically load the panel
   */
  libraryPanel?: LibraryPanelRef;
  /**
   * Panel links.
   * TODO fill this out - seems there are a couple variants?
   */
  links?: Array<DashboardLink>;
  /**
   * TODO docs
   */
  maxDataPoints?: number;
  /**
   * options is specified by the PanelOptions field in panel
   * plugin schemas.
   */
  options: Record<string, unknown>;
  /**
   * FIXME this almost certainly has to be changed in favor of scuemata versions
   */
  pluginVersion?: string;
  /**
   * Name of template variable to repeat for.
   */
  repeat?: string;
  /**
   * Direction to repeat in if 'repeat' is set.
   * "h" for horizontal, "v" for vertical.
   * TODO this is probably optional
   */
  repeatDirection: ('h' | 'v');
  /**
   * Id of the repeating panel.
   */
  repeatPanelId?: number;
  /**
   * TODO docs
   */
  tags?: Array<string>;
  /**
   * TODO docs
   */
  targets?: Array<Record<string, unknown>>;
  /**
   * TODO docs - seems to be an old field from old dashboard alerts?
   */
  thresholds?: Array<unknown>;
  /**
   * TODO docs
   * TODO tighter constraint
   */
  timeFrom?: string;
  /**
   * TODO docs
   */
  timeRegions?: Array<unknown>;
  /**
   * TODO docs
   * TODO tighter constraint
   */
  timeShift?: string;
  /**
   * Panel title.
   */
  title?: string;
  transformations: Array<DataTransformerConfig>;
  /**
   * Whether to display the panel without a background.
   */
  transparent: boolean;
  /**
   * The panel plugin type id. May not be empty.
   */
  type: string;
}

export const defaultPanel: Partial<Panel> = {
  links: [],
  repeatDirection: 'h',
  tags: [],
  targets: [],
  thresholds: [],
  timeRegions: [],
  transformations: [],
  transparent: false,
};

export interface FieldConfigSource {
  defaults: FieldConfig;
  overrides: Array<{
    matcher: MatcherConfig;
    properties: Array<{
      id: string;
      value?: unknown;
    }>;
  }>;
}

export const defaultFieldConfigSource: Partial<FieldConfigSource> = {
  overrides: [],
};

export interface LibraryPanelRef {
  name: string;
  uid: string;
}

export interface MatcherConfig {
  id: string;
  options?: unknown;
}

export const defaultMatcherConfig: Partial<MatcherConfig> = {
  id: '',
};

export interface FieldConfig {
  /**
   * Map values to a display color
   */
  color?: FieldColor;
  /**
   * custom is specified by the PanelFieldConfig field
   * in panel plugin schemas.
   */
  custom?: Record<string, unknown>;
  /**
   * Significant digits (for display)
   */
  decimals?: number;
  /**
   * Human readable field metadata
   */
  description?: string;
  /**
   * The display value for this field.  This supports template variables blank is auto
   */
  displayName?: string;
  /**
   * This can be used by data sources that return and explicit naming structure for values and labels
   * When this property is configured, this value is used rather than the default naming strategy.
   */
  displayNameFromDS?: string;
  /**
   * True if data source field supports ad-hoc filters
   */
  filterable?: boolean;
  /**
   * The behavior when clicking on a result
   */
  links?: Array<unknown>;
  /**
   * Convert input values into a display string
   */
  mappings?: Array<ValueMapping>;
  max?: number;
  min?: number;
  /**
   * Alternative to empty string
   */
  noValue?: string;
  /**
   * An explicit path to the field in the datasource.  When the frame meta includes a path,
   * This will default to `${frame.meta.path}/${field.name}
   * 
   * When defined, this value can be used as an identifier within the datasource scope, and
   * may be used to update the results
   */
  path?: string;
  /**
   * Map numeric values to states
   */
  thresholds?: ThresholdsConfig;
  /**
   * Numeric Options
   */
  unit?: string;
  /**
   * True if data source can write a value to the path.  Auth/authz are supported separately
   */
  writeable?: boolean;
}

export const defaultFieldConfig: Partial<FieldConfig> = {
  links: [],
  mappings: [],
};

/**
 * Row panel
 */
export interface RowPanel {
  collapsed: boolean;
  /**
   * Name of default datasource.
   */
  datasource?: {
    type?: string;
    uid?: string;
  };
  gridPos?: GridPos;
  id: number;
  panels: Array<(Panel | GraphPanel | HeatmapPanel)>;
  /**
   * Name of template variable to repeat for.
   */
  repeat?: string;
  title?: string;
  type: 'row';
}

export const defaultRowPanel: Partial<RowPanel> = {
  collapsed: false,
  panels: [],
};

/**
 * Support for legacy graph and heatmap panels.
 */
export interface GraphPanel {
  /**
   * @deprecated this is part of deprecated graph panel
   */
  legend?: {
    show: boolean;
    sort?: string;
    sortDesc?: boolean;
  };
  type: 'graph';
}

export interface HeatmapPanel {
  type: 'heatmap';
}

export interface Dashboard {
  /**
   * TODO docs
   */
  annotations?: {
    list?: Array<AnnotationQuery>;
  };
  /**
   * Description of dashboard.
   */
  description?: string;
  /**
   * Whether a dashboard is editable or not.
   */
  editable: boolean;
  /**
   * The month that the fiscal year starts on.  0 = January, 11 = December
   */
  fiscalYearStartMonth?: number;
  gnetId?: string;
  /**
   * Configuration of dashboard cursor sync behavior.
   */
  graphTooltip: DashboardCursorSync;
  /**
   * Unique numeric identifier for the dashboard.
   * TODO must isolate or remove identifiers local to a Grafana instance...?
   */
  id?: number;
  /**
   * TODO docs
   */
  links?: Array<DashboardLink>;
  /**
   * When set to true, the dashboard will redraw panels at an interval matching the pixel width.
   * This will keep data "moving left" regardless of the query refresh rate.  This setting helps
   * avoid dashboards presenting stale live data
   */
  liveNow?: boolean;
  panels?: Array<(Panel | RowPanel | GraphPanel | HeatmapPanel)>;
  /**
   * Refresh rate of dashboard. Represented via interval string, e.g. "5s", "1m", "1h", "1d".
   */
  refresh?: (string | false);
  /**
   * This property should only be used in dashboards defined by plugins.  It is a quick check
   * to see if the version has changed since the last time.  Unclear why using the version property
   * is insufficient.
   */
  revision?: number;
  /**
   * Version of the JSON schema, incremented each time a Grafana update brings
   * changes to said schema.
   * TODO this is the existing schema numbering system. It will be replaced by Thema's themaVersion
   */
  schemaVersion: number;
  snapshot?: {
    /**
     * TODO docs
     */
    created: string;
    /**
     * TODO docs
     */
    expires: string;
    /**
     * TODO docs
     */
    external: boolean;
    /**
     * TODO docs
     */
    externalUrl: string;
    /**
     * TODO docs
     */
    id: number;
    /**
     * TODO docs
     */
    key: string;
    /**
     * TODO docs
     */
    name: string;
    /**
     * TODO docs
     */
    orgId: number;
    /**
     * TODO docs
     */
    updated: string;
    /**
     * TODO docs
     */
    url?: string;
    /**
     * TODO docs
     */
    userId: number;
  };
  /**
   * Theme of dashboard.
   */
  style: ('light' | 'dark');
  /**
   * Tags associated with dashboard.
   */
  tags?: Array<string>;
  /**
   * TODO docs
   */
  templating?: {
    list?: Array<VariableModel>;
  };
  /**
   * Time range for dashboard, e.g. last 6 hours, last 7 days, etc
   */
  time?: {
    from: string;
    to: string;
  };
  /**
   * TODO docs
   * TODO this appears to be spread all over in the frontend. Concepts will likely need tidying in tandem with schema changes
   */
  timepicker?: {
    /**
     * Whether timepicker is collapsed or not.
     */
    collapse: boolean;
    /**
     * Whether timepicker is enabled or not.
     */
    enable: boolean;
    /**
     * Whether timepicker is visible or not.
     */
    hidden: boolean;
    /**
     * Selectable intervals for auto-refresh.
     */
    refresh_intervals: Array<string>;
    /**
     * TODO docs
     */
    time_options: Array<string>;
  };
  /**
   * Timezone of dashboard. Accepts IANA TZDB zone ID or "browser" or "utc".
   */
  timezone?: string;
  /**
   * Title of dashboard.
   */
  title?: string;
  /**
   * Unique dashboard identifier that can be generated by anyone. string (8-40)
   */
  uid?: string;
  /**
   * Version of the dashboard, incremented each time the dashboard is updated.
   */
  version?: number;
  /**
   * TODO docs
   */
  weekStart?: string;
}

export const defaultDashboard: Partial<Dashboard> = {
  editable: true,
  fiscalYearStartMonth: 0,
  graphTooltip: DashboardCursorSync.Off,
  links: [],
  panels: [],
  schemaVersion: 36,
  style: 'dark',
  tags: [],
  timezone: 'browser',
};
