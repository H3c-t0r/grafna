// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     kinds/gen.go
// Using jennies:
//     GoResourceTypes
//
// Run 'make gen-cue' from repository root to regenerate.

package dashboard

import (
	"time"
)

// Defines values for CursorSync.
const (
	CursorSyncN0 CursorSync = 0
	CursorSyncN1 CursorSync = 1
	CursorSyncN2 CursorSync = 2
)

// Defines values for LinkType.
const (
	LinkTypeDashboards LinkType = "dashboards"
	LinkTypeLink       LinkType = "link"
)

// Defines values for FieldColorModeId.
const (
	FieldColorModeIdContinuousBlPu       FieldColorModeId = "continuous-BlPu"
	FieldColorModeIdContinuousBlYlRd     FieldColorModeId = "continuous-BlYlRd"
	FieldColorModeIdContinuousBlues      FieldColorModeId = "continuous-blues"
	FieldColorModeIdContinuousGrYlRd     FieldColorModeId = "continuous-GrYlRd"
	FieldColorModeIdContinuousGreens     FieldColorModeId = "continuous-greens"
	FieldColorModeIdContinuousPurples    FieldColorModeId = "continuous-purples"
	FieldColorModeIdContinuousRdYlGr     FieldColorModeId = "continuous-RdYlGr"
	FieldColorModeIdContinuousReds       FieldColorModeId = "continuous-reds"
	FieldColorModeIdContinuousYlBl       FieldColorModeId = "continuous-YlBl"
	FieldColorModeIdContinuousYlRd       FieldColorModeId = "continuous-YlRd"
	FieldColorModeIdFixed                FieldColorModeId = "fixed"
	FieldColorModeIdPaletteClassic       FieldColorModeId = "palette-classic"
	FieldColorModeIdPaletteClassicByName FieldColorModeId = "palette-classic-by-name"
	FieldColorModeIdShades               FieldColorModeId = "shades"
	FieldColorModeIdThresholds           FieldColorModeId = "thresholds"
)

// Defines values for FieldColorSeriesByMode.
const (
	FieldColorSeriesByModeLast FieldColorSeriesByMode = "last"
	FieldColorSeriesByModeMax  FieldColorSeriesByMode = "max"
	FieldColorSeriesByModeMin  FieldColorSeriesByMode = "min"
)

// Defines values for GraphPanelType.
const (
	GraphPanelTypeGraph GraphPanelType = "graph"
)

// Defines values for HeatmapPanelType.
const (
	HeatmapPanelTypeHeatmap HeatmapPanelType = "heatmap"
)

// Defines values for MappingType.
const (
	MappingTypeRange   MappingType = "range"
	MappingTypeRegex   MappingType = "regex"
	MappingTypeSpecial MappingType = "special"
	MappingTypeValue   MappingType = "value"
)

// Defines values for PanelRepeatDirection.
const (
	PanelRepeatDirectionH PanelRepeatDirection = "h"
	PanelRepeatDirectionV PanelRepeatDirection = "v"
)

// Defines values for RangeMapType.
const (
	RangeMapTypeRange   RangeMapType = "range"
	RangeMapTypeRegex   RangeMapType = "regex"
	RangeMapTypeSpecial RangeMapType = "special"
	RangeMapTypeValue   RangeMapType = "value"
)

// Defines values for RegexMapType.
const (
	RegexMapTypeRange   RegexMapType = "range"
	RegexMapTypeRegex   RegexMapType = "regex"
	RegexMapTypeSpecial RegexMapType = "special"
	RegexMapTypeValue   RegexMapType = "value"
)

// Defines values for RowPanelType.
const (
	RowPanelTypeRow RowPanelType = "row"
)

// Defines values for SpecStyle.
const (
	SpecStyleDark  SpecStyle = "dark"
	SpecStyleLight SpecStyle = "light"
)

// Defines values for SpecialValueMapType.
const (
	SpecialValueMapTypeRange   SpecialValueMapType = "range"
	SpecialValueMapTypeRegex   SpecialValueMapType = "regex"
	SpecialValueMapTypeSpecial SpecialValueMapType = "special"
	SpecialValueMapTypeValue   SpecialValueMapType = "value"
)

// Defines values for SpecialValueMatch.
const (
	SpecialValueMatchEmpty   SpecialValueMatch = "empty"
	SpecialValueMatchFalse   SpecialValueMatch = "false"
	SpecialValueMatchNan     SpecialValueMatch = "nan"
	SpecialValueMatchNull    SpecialValueMatch = "null"
	SpecialValueMatchNullNan SpecialValueMatch = "null+nan"
	SpecialValueMatchTrue    SpecialValueMatch = "true"
)

// Defines values for ThresholdsMode.
const (
	ThresholdsModeAbsolute   ThresholdsMode = "absolute"
	ThresholdsModePercentage ThresholdsMode = "percentage"
)

// Defines values for ValueMapType.
const (
	ValueMapTypeRange   ValueMapType = "range"
	ValueMapTypeRegex   ValueMapType = "regex"
	ValueMapTypeSpecial ValueMapType = "special"
	ValueMapTypeValue   ValueMapType = "value"
)

// Defines values for VariableHide.
const (
	VariableHideN0 VariableHide = 0
	VariableHideN1 VariableHide = 1
	VariableHideN2 VariableHide = 2
)

// Defines values for VariableRefresh.
const (
	VariableRefreshN0 VariableRefresh = 0
	VariableRefreshN1 VariableRefresh = 1
	VariableRefreshN2 VariableRefresh = 2
)

// Defines values for VariableType.
const (
	VariableTypeAdhoc      VariableType = "adhoc"
	VariableTypeConstant   VariableType = "constant"
	VariableTypeCustom     VariableType = "custom"
	VariableTypeDatasource VariableType = "datasource"
	VariableTypeInterval   VariableType = "interval"
	VariableTypeQuery      VariableType = "query"
	VariableTypeSystem     VariableType = "system"
	VariableTypeTextbox    VariableType = "textbox"
)

// TODO -- should not be a public interface on its own, but required for Veneer
type AnnotationContainer struct {
	// List of annotations
	List []AnnotationQuery `json:"list,omitempty"`
}

// AnnotationPanelFilter defines model for AnnotationPanelFilter.
type AnnotationPanelFilter struct {
	// Should the specified panels be included or excluded
	Exclude *bool `json:"exclude,omitempty"`

	// Panel IDs that should be included or excluded
	Ids []int `json:"ids"`
}

// TODO docs
// FROM: AnnotationQuery in grafana-data/src/types/annotations.ts
type AnnotationQuery struct {
	// Ref to a DataSource instance
	Datasource DataSourceRef `json:"datasource"`

	// When enabled the annotation query is issued with every dashboard refresh
	Enable bool                   `json:"enable"`
	Filter *AnnotationPanelFilter `json:"filter,omitempty"`

	// Annotation queries can be toggled on or off at the top of the dashboard.
	// When hide is true, the toggle is not shown in the dashboard.
	Hide *bool `json:"hide,omitempty"`

	// Color to use for the annotation event markers
	IconColor string `json:"iconColor"`

	// Name of annotation.
	Name string `json:"name"`

	// TODO: this should be a regular DataQuery that depends on the selected dashboard
	// these match the properties of the "grafana" datasouce that is default in most dashboards
	Target *AnnotationTarget `json:"target,omitempty"`

	// TODO -- this should not exist here, it is based on the --grafana-- datasource
	Type *string `json:"type,omitempty"`
}

// TODO: this should be a regular DataQuery that depends on the selected dashboard
// these match the properties of the "grafana" datasouce that is default in most dashboards
type AnnotationTarget struct {
	// Only required/valid for the grafana datasource...
	// but code+tests is already depending on it so hard to change
	Limit int64 `json:"limit"`

	// Only required/valid for the grafana datasource...
	// but code+tests is already depending on it so hard to change
	MatchAny bool `json:"matchAny"`

	// Only required/valid for the grafana datasource...
	// but code+tests is already depending on it so hard to change
	Tags []string `json:"tags"`

	// Only required/valid for the grafana datasource...
	// but code+tests is already depending on it so hard to change
	Type string `json:"type"`
}

// 0 for no shared crosshair or tooltip (default).
// 1 for shared crosshair.
// 2 for shared crosshair AND shared tooltip.
type CursorSync int

// Links with references to other dashboards or external resources
type Link struct {
	// If true, all dashboards links will be displayed in a dropdown. If false, all dashboards links will be displayed side by side. Only valid if the type is dashboards
	AsDropdown bool `json:"asDropdown"`

	// Icon name to be displayed with the link
	Icon string `json:"icon"`

	// If true, includes current template variables values in the link as query params
	IncludeVars bool `json:"includeVars"`

	// If true, includes current time range in the link as query params
	KeepTime bool `json:"keepTime"`

	// List of tags to limit the linked dashboards. If empty, all dashboards will be displayed. Only valid if the type is dashboards
	Tags []string `json:"tags"`

	// If true, the link will be opened in a new tab
	TargetBlank bool `json:"targetBlank"`

	// Title to display with the link
	Title string `json:"title"`

	// Tooltip to display when the user hovers their mouse over it
	Tooltip string `json:"tooltip"`

	// Dashboard Link type. Accepted values are dashboards (to refer to another dashboard) and link (to refer to an external resource)
	Type LinkType `json:"type"`

	// Link URL. Only required/valid if the type is link
	Url string `json:"url"`
}

// Dashboard Link type. Accepted values are dashboards (to refer to another dashboard) and link (to refer to an external resource)
type LinkType string

// Ref to a DataSource instance
type DataSourceRef struct {
	// The plugin type-id
	Type *string `json:"type,omitempty"`

	// Specific datasource instance
	Uid *string `json:"uid,omitempty"`
}

// Transformations allow to manipulate data returned by a query before the system applies a visualization.
// Using transformations you can: rename fields, join time series data, perform mathematical operations across queries,
// use the output of one transformation as the input to another transformation, etc.
type DataTransformerConfig struct {
	// Disabled transformations are skipped
	Disabled *bool `json:"disabled,omitempty"`

	// Optional frame matcher. When missing it will be applied to all results
	Filter *MatcherConfig `json:"filter,omitempty"`

	// Unique identifier of transformer
	Id string `json:"id"`

	// Options to be passed to the transformer
	// Valid options depend on the transformer id
	Options interface{} `json:"options"`
}

// DynamicConfigValue defines model for DynamicConfigValue.
type DynamicConfigValue struct {
	Id    string       `json:"id"`
	Value *interface{} `json:"value,omitempty"`
}

// Map a field to a color.
type FieldColor struct {
	// The fixed color value for fixed or shades color modes.
	FixedColor *string `json:"fixedColor,omitempty"`

	// Color mode for a field. You can specify a single color, or select a continuous (gradient) color schemes, based on a value.
	// Continuous color interpolates a color using the percentage of a value relative to min and max.
	// Accepted values are:
	// thresholds: From thresholds. Informs Grafana to take the color from the matching threshold
	// palette-classic: Classic palette. Grafana will assign color by looking up a color in a palette by series index. Useful for Graphs and pie charts and other categorical data visualizations
	// palette-classic-by-name: Classic palette (by name). Grafana will assign color by looking up a color in a palette by series name. Useful for Graphs and pie charts and other categorical data visualizations
	// continuous-GrYlRd: ontinuous Green-Yellow-Red palette mode
	// continuous-RdYlGr: Continuous Red-Yellow-Green palette mode
	// continuous-BlYlRd: Continuous Blue-Yellow-Red palette mode
	// continuous-YlRd: Continuous Yellow-Red palette mode
	// continuous-BlPu: Continuous Blue-Purple palette mode
	// continuous-YlBl: Continuous Yellow-Blue palette mode
	// continuous-blues: Continuous Blue palette mode
	// continuous-reds: Continuous Red palette mode
	// continuous-greens: Continuous Green palette mode
	// continuous-purples: Continuous Purple palette mode
	// shades: Shades of a single color. Specify a single color, useful in an override rule.
	// fixed: Fixed color mode. Specify a single color, useful in an override rule.
	Mode FieldColorModeId `json:"mode"`

	// Defines how to assign a series color from "by value" color schemes. For example for an aggregated data points like a timeseries, the color can be assigned by the min, max or last value.
	SeriesBy *FieldColorSeriesByMode `json:"seriesBy,omitempty"`
}

// Color mode for a field. You can specify a single color, or select a continuous (gradient) color schemes, based on a value.
// Continuous color interpolates a color using the percentage of a value relative to min and max.
// Accepted values are:
// thresholds: From thresholds. Informs Grafana to take the color from the matching threshold
// palette-classic: Classic palette. Grafana will assign color by looking up a color in a palette by series index. Useful for Graphs and pie charts and other categorical data visualizations
// palette-classic-by-name: Classic palette (by name). Grafana will assign color by looking up a color in a palette by series name. Useful for Graphs and pie charts and other categorical data visualizations
// continuous-GrYlRd: ontinuous Green-Yellow-Red palette mode
// continuous-RdYlGr: Continuous Red-Yellow-Green palette mode
// continuous-BlYlRd: Continuous Blue-Yellow-Red palette mode
// continuous-YlRd: Continuous Yellow-Red palette mode
// continuous-BlPu: Continuous Blue-Purple palette mode
// continuous-YlBl: Continuous Yellow-Blue palette mode
// continuous-blues: Continuous Blue palette mode
// continuous-reds: Continuous Red palette mode
// continuous-greens: Continuous Green palette mode
// continuous-purples: Continuous Purple palette mode
// shades: Shades of a single color. Specify a single color, useful in an override rule.
// fixed: Fixed color mode. Specify a single color, useful in an override rule.
type FieldColorModeId string

// Defines how to assign a series color from "by value" color schemes. For example for an aggregated data points like a timeseries, the color can be assigned by the min, max or last value.
type FieldColorSeriesByMode string

// FieldConfig defines model for FieldConfig.
type FieldConfig struct {
	// Map a field to a color.
	Color *FieldColor `json:"color,omitempty"`

	// custom is specified by the FieldConfig field
	// in panel plugin schemas.
	Custom map[string]interface{} `json:"custom,omitempty"`

	// Significant digits (for display)
	Decimals *float32 `json:"decimals,omitempty"`

	// Human readable field metadata
	Description *string `json:"description,omitempty"`

	// The display value for this field.  This supports template variables blank is auto
	DisplayName *string `json:"displayName,omitempty"`

	// This can be used by data sources that return and explicit naming structure for values and labels
	// When this property is configured, this value is used rather than the default naming strategy.
	DisplayNameFromDS *string `json:"displayNameFromDS,omitempty"`

	// True if data source field supports ad-hoc filters
	Filterable *bool `json:"filterable,omitempty"`

	// The behavior when clicking on a result
	Links []interface{} `json:"links,omitempty"`

	// Convert input values into a display string
	Mappings []interface{} `json:"mappings,omitempty"`
	Max      *float32      `json:"max,omitempty"`
	Min      *float32      `json:"min,omitempty"`

	// Alternative to empty string
	NoValue *string `json:"noValue,omitempty"`

	// An explicit path to the field in the datasource.  When the frame meta includes a path,
	// This will default to `${frame.meta.path}/${field.name}
	//
	// When defined, this value can be used as an identifier within the datasource scope, and
	// may be used to update the results
	Path *string `json:"path,omitempty"`

	// Thresholds configuration for the panel
	Thresholds *ThresholdsConfig `json:"thresholds,omitempty"`

	// Numeric Options
	Unit *string `json:"unit,omitempty"`

	// True if data source can write a value to the path.  Auth/authz are supported separately
	Writeable *bool `json:"writeable,omitempty"`
}

// FieldConfigSource defines model for FieldConfigSource.
type FieldConfigSource struct {
	Defaults  FieldConfig `json:"defaults"`
	Overrides []struct {
		// Optional frame matcher. When missing it will be applied to all results
		Matcher    MatcherConfig        `json:"matcher"`
		Properties []DynamicConfigValue `json:"properties"`
	} `json:"overrides"`
}

// Support for legacy graph and heatmap panels.
type GraphPanel struct {
	// @deprecated this is part of deprecated graph panel
	Legend *struct {
		Show     bool    `json:"show"`
		Sort     *string `json:"sort,omitempty"`
		SortDesc *bool   `json:"sortDesc,omitempty"`
	} `json:"legend,omitempty"`
	Type GraphPanelType `json:"type"`
}

// GraphPanelType defines model for GraphPanel.Type.
type GraphPanelType string

// Position and dimensions of a panel in the grid
type GridPos struct {
	// Panel height. The height is the number of rows from the top edge of the grid
	H int `json:"h"`

	// Whether the panel is fixed within the grid. If true, the panel will not be affected by other panels' interactions
	Static *bool `json:"static,omitempty"`

	// Panel width. The width is the number of columns from the left edge of the grid
	W int `json:"w"`

	// Panel x. The x coordinate is the number of columns from the left edge of the grid
	X int `json:"x"`

	// Panel y. The y coordinate is the number of rows from the top edge of the grid
	Y int `json:"y"`
}

// HeatmapPanel defines model for HeatmapPanel.
type HeatmapPanel struct {
	Type HeatmapPanelType `json:"type"`
}

// HeatmapPanelType defines model for HeatmapPanel.Type.
type HeatmapPanelType string

// LibraryPanelRef defines model for LibraryPanelRef.
type LibraryPanelRef struct {
	Name string `json:"name"`
	Uid  string `json:"uid"`
}

// Supported value mapping types
// ValueToText: Maps text values to a color or different display text and color. For example, you can configure a value mapping so that all instances of the value 10 appear as Perfection! rather than the number.
// RangeToText: Maps numerical ranges to a display text and color. For example, if a value is within a certain range, you can configure a range value mapping to display Low or High rather than the number.
// RegexToText: Maps regular expressions to replacement text and a color. For example, if a value is www.example.com, you can configure a regex value mapping so that Grafana displays www and truncates the domain.
// SpecialValue: Maps special values like Null, NaN (not a number), and boolean values like true and false to a display text and color. See SpecialValueMatch to see the list of special values. For example, you can configure a special value mapping so that null values appear as N/A.
type MappingType string

// Optional frame matcher. When missing it will be applied to all results
type MatcherConfig struct {
	Id      string       `json:"id"`
	Options *interface{} `json:"options,omitempty"`
}

// Dashboard panels are the basic visualization building blocks.
type Panel struct {
	// The datasource used in all targets.
	Datasource *struct {
		Type *string `json:"type,omitempty"`
		Uid  *string `json:"uid,omitempty"`
	} `json:"datasource,omitempty"`

	// Description Description.
	Description *string           `json:"description,omitempty"`
	FieldConfig FieldConfigSource `json:"fieldConfig"`

	// Position and dimensions of a panel in the grid
	GridPos *GridPos `json:"gridPos,omitempty"`

	// TODO docs
	Id *int `json:"id,omitempty"`

	// The min time interval setting defines a lower limit for the $__interval and $__interval_ms variables.
	// This value must be formatted as a number followed by a valid time
	// identifier like: "40s", "3d", etc.
	// See: https://grafana.com/docs/grafana/latest/panels-visualizations/query-transform-data/#query-options
	Interval     *string          `json:"interval,omitempty"`
	LibraryPanel *LibraryPanelRef `json:"libraryPanel,omitempty"`

	// Panel links.
	// TODO fill this out - seems there are a couple variants?
	Links []Link `json:"links,omitempty"`

	// The maximum number of data points that the panel queries are retrieving.
	MaxDataPoints *float32 `json:"maxDataPoints,omitempty"`

	// options is specified by the Options field in panel
	// plugin schemas.
	Options map[string]interface{} `json:"options"`

	// FIXME this almost certainly has to be changed in favor of scuemata versions
	PluginVersion *string `json:"pluginVersion,omitempty"`

	// Name of template variable to repeat for.
	Repeat *string `json:"repeat,omitempty"`

	// Direction to repeat in if 'repeat' is set.
	// "h" for horizontal, "v" for vertical.
	// TODO this is probably optional
	RepeatDirection PanelRepeatDirection `json:"repeatDirection"`

	// Id of the repeating panel.
	RepeatPanelId *int64 `json:"repeatPanelId,omitempty"`

	// TODO docs
	Tags []string `json:"tags,omitempty"`

	// TODO docs
	Targets []Target `json:"targets,omitempty"`

	// TODO docs - seems to be an old field from old dashboard alerts?
	Thresholds []interface{} `json:"thresholds,omitempty"`

	// Overrides the relative time range for individual panels,
	// which causes them to be different than what is selected in
	// the dashboard time picker in the top-right corner of the dashboard. You can use this to show metrics from different
	// time periods or days on the same dashboard.
	// The value is formatted as time operation like: `now-5m` (Last 5 minutes), `now/d` (the day so far),
	// `now-5d/d`(Last 5 days), `now/w` (This week so far), `now-2y/y` (Last 2 years).
	// Note: Panel time overrides have no effect when the dashboard’s time range is absolute.
	// See: https://grafana.com/docs/grafana/latest/panels-visualizations/query-transform-data/#query-options
	TimeFrom *string `json:"timeFrom,omitempty"`

	// TODO docs
	TimeRegions []interface{} `json:"timeRegions,omitempty"`

	// Overrides the time range for individual panels by shifting its start and end relative to the time picker.
	// For example, you can shift the time range for the panel to be two hours earlier than the dashboard time picker setting `2h`.
	// Note: Panel time overrides have no effect when the dashboard’s time range is absolute.
	// See: https://grafana.com/docs/grafana/latest/panels-visualizations/query-transform-data/#query-options
	TimeShift *string `json:"timeShift,omitempty"`

	// Panel title.
	Title *string `json:"title,omitempty"`

	// List of transformations that are applied to the panel data before rendering.
	// When there are multiple transformations, Grafana applies them in the order they are listed.
	// Each transformation creates a result set that then passes on to the next transformation in the processing pipeline.
	Transformations []DataTransformerConfig `json:"transformations"`

	// Whether to display the panel without a background.
	Transparent bool `json:"transparent"`

	// The panel plugin type id. May not be empty.
	Type string `json:"type"`
}

// Direction to repeat in if 'repeat' is set.
// "h" for horizontal, "v" for vertical.
// TODO this is probably optional
type PanelRepeatDirection string

// Maps numerical ranges to a display text and color.
// For example, if a value is within a certain range, you can configure a range value mapping to display Low or High rather than the number.
type RangeMap struct {
	// Range to match against and the result to apply when the value is within the range
	Options struct {
		// Min value of the range. It can be null which means -Infinity
		From float64 `json:"from"`

		// Result used as replacement with text and color when the value matches
		Result ValueMappingResult `json:"result"`

		// Max value of the range. It can be null which means +Infinity
		To float64 `json:"to"`
	} `json:"options"`
	Type RangeMapType `json:"type"`
}

// RangeMapType defines model for RangeMap.Type.
type RangeMapType string

// Maps regular expressions to replacement text and a color.
// For example, if a value is www.example.com, you can configure a regex value mapping so that Grafana displays www and truncates the domain.
type RegexMap struct {
	// Regular expression to match against and the result to apply when the value matches the regex
	Options struct {
		// Regular expression to match against
		Pattern string `json:"pattern"`

		// Result used as replacement with text and color when the value matches
		Result ValueMappingResult `json:"result"`
	} `json:"options"`
	Type RegexMapType `json:"type"`
}

// RegexMapType defines model for RegexMap.Type.
type RegexMapType string

// Row panel
type RowPanel struct {
	Collapsed bool `json:"collapsed"`

	// Name of default datasource.
	Datasource *struct {
		Type *string `json:"type,omitempty"`
		Uid  *string `json:"uid,omitempty"`
	} `json:"datasource,omitempty"`

	// Position and dimensions of a panel in the grid
	GridPos *GridPos      `json:"gridPos,omitempty"`
	Id      int           `json:"id"`
	Panels  []interface{} `json:"panels"`

	// Name of template variable to repeat for.
	Repeat *string      `json:"repeat,omitempty"`
	Title  *string      `json:"title,omitempty"`
	Type   RowPanelType `json:"type"`
}

// RowPanelType defines model for RowPanel.Type.
type RowPanelType string

// A dashboard snapshot shares an interactive dashboard publicly.
// It is a read-only version of a dashboard, and is not editable.
// It is possible to create a snapshot of a snapshot.
// Grafana strips away all sensitive information from the dashboard.
// Sensitive information stripped: queries (metric, template,annotation) and panel links.
type Snapshot struct {
	// Time when the snapshot was created
	Created time.Time `json:"created"`

	// Time when the snapshot expires, default is never to expire
	Expires string `json:"expires"`

	// Is the snapshot saved in an external grafana instance
	External bool `json:"external"`

	// ExternalUrl external url, if snapshot was shared in external grafana instance
	ExternalUrl string `json:"externalUrl"`

	// Unique identifier of the snapshot
	Id int `json:"id"`

	// Optional, defined the unique key of the snapshot, required if external is true
	Key string `json:"key"`

	// Optional, name of the snapshot
	Name string `json:"name"`

	// OrgId org id of the snapshot
	OrgId int `json:"orgId"`

	// Updated last time when the snapshot was updated
	Updated time.Time `json:"updated"`

	// url of the snapshot, if snapshot was shared internally
	Url *string `json:"url,omitempty"`

	// UserId user id of the snapshot creator
	UserId int `json:"userId"`
}

// Spec defines model for Spec.
type Spec struct {
	// TODO -- should not be a public interface on its own, but required for Veneer
	Annotations *AnnotationContainer `json:"annotations,omitempty"`

	// Description of dashboard.
	Description *string `json:"description,omitempty"`

	// Whether a dashboard is editable or not.
	Editable bool `json:"editable"`

	// The month that the fiscal year starts on.  0 = January, 11 = December
	FiscalYearStartMonth *int `json:"fiscalYearStartMonth,omitempty"`

	// ID of a dashboard imported from the https://grafana.com/grafana/dashboards/ portal
	GnetId *string `json:"gnetId,omitempty"`

	// 0 for no shared crosshair or tooltip (default).
	// 1 for shared crosshair.
	// 2 for shared crosshair AND shared tooltip.
	GraphTooltip CursorSync `json:"graphTooltip"`

	// Unique numeric identifier for the dashboard.
	// `id` is internal to a specific Grafana instance. `uid` should be used to identify a dashboard across Grafana instances.
	Id *int64 `json:"id,omitempty"`

	// Links with references to other dashboards or external websites.
	Links []Link `json:"links,omitempty"`

	// When set to true, the dashboard will redraw panels at an interval matching the pixel width.
	// This will keep data "moving left" regardless of the query refresh rate. This setting helps
	// avoid dashboards presenting stale live data
	LiveNow *bool `json:"liveNow,omitempty"`

	// List of dashboard panels
	Panels []interface{} `json:"panels,omitempty"`

	// Refresh rate of dashboard. Represented via interval string, e.g. "5s", "1m", "1h", "1d".
	Refresh *interface{} `json:"refresh,omitempty"`

	// This property should only be used in dashboards defined by plugins.  It is a quick check
	// to see if the version has changed since the last time.
	Revision *int64 `json:"revision,omitempty"`

	// Version of the JSON schema, incremented each time a Grafana update brings
	// changes to said schema.
	SchemaVersion int `json:"schemaVersion"`

	// A dashboard snapshot shares an interactive dashboard publicly.
	// It is a read-only version of a dashboard, and is not editable.
	// It is possible to create a snapshot of a snapshot.
	// Grafana strips away all sensitive information from the dashboard.
	// Sensitive information stripped: queries (metric, template,annotation) and panel links.
	Snapshot *Snapshot `json:"snapshot,omitempty"`

	// Theme of dashboard.
	// Default value: dark.
	Style SpecStyle `json:"style"`

	// Tags associated with dashboard.
	Tags []string `json:"tags,omitempty"`

	// Contains the list of configured template variables with their saved values along with some other metadata
	Templating *struct {
		List []VariableModel `json:"list,omitempty"`
	} `json:"templating,omitempty"`

	// Time range for dashboard.
	// Accepted values are relative time strings like {from: 'now-6h', to: 'now'} or absolute time strings like {from: '2020-07-10T08:00:00.000Z', to: '2020-07-10T14:00:00.000Z'}.
	Time *struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"time,omitempty"`

	// Configuration of the time picker shown at the top of a dashboard.
	Timepicker *struct {
		// Whether timepicker is collapsed or not. Has no effect on provisioned dashboard.
		Collapse bool `json:"collapse"`

		// Whether timepicker is enabled or not. Has no effect on provisioned dashboard.
		Enable bool `json:"enable"`

		// Whether timepicker is visible or not.
		Hidden bool `json:"hidden"`

		// Interval options available in the refresh picker dropdown.
		RefreshIntervals []string `json:"refresh_intervals"`

		// Selectable options available in the time picker dropdown. Has no effect on provisioned dashboard.
		TimeOptions []string `json:"time_options"`
	} `json:"timepicker,omitempty"`

	// Timezone of dashboard. Accepted values are IANA TZDB zone ID or "browser" or "utc".
	Timezone *string `json:"timezone,omitempty"`

	// Title of dashboard.
	Title *string `json:"title,omitempty"`

	// Unique dashboard identifier that can be generated by anyone. string (8-40)
	Uid *string `json:"uid,omitempty"`

	// Version of the dashboard, incremented each time the dashboard is updated.
	Version *int `json:"version,omitempty"`

	// Day when the week starts. Expressed by the name of the day in lowercase, e.g. "monday".
	WeekStart *string `json:"weekStart,omitempty"`
}

// Theme of dashboard.
// Default value: dark.
type SpecStyle string

// Maps special values like Null, NaN (not a number), and boolean values like true and false to a display text and color.
// See SpecialValueMatch to see the list of special values.
// For example, you can configure a special value mapping so that null values appear as N/A.
type SpecialValueMap struct {
	Options struct {
		// Special value types supported by the SpecialValueMap
		Match SpecialValueMatch `json:"match"`

		// Result used as replacement with text and color when the value matches
		Result ValueMappingResult `json:"result"`
	} `json:"options"`
	Type SpecialValueMapType `json:"type"`
}

// SpecialValueMapType defines model for SpecialValueMap.Type.
type SpecialValueMapType string

// Special value types supported by the SpecialValueMap
type SpecialValueMatch string

// Schema for panel targets is specified by datasource
// plugins. We use a placeholder definition, which the Go
// schema loader either left open/as-is with the Base
// variant of the Dashboard and Panel families, or filled
// with types derived from plugins in the Instance variant.
// When working directly from CUE, importers can extend this
// type directly to achieve the same effect.
type Target = map[string]interface{}

// User-defined value for a metric that triggers visual changes in a panel when this value is met or exceeded
// They are used to conditionally style and color visualizations based on query results , and can be applied to most visualizations.
type Threshold struct {
	// Color represents the color of the visual change that will occur in the dashboard when the threshold value is met or exceeded.
	Color string `json:"color"`

	// Value represents a specified metric for the threshold, which triggers a visual change in the dashboard when this value is met or exceeded.
	// Nulls currently appear here when serializing -Infinity to JSON.
	Value float32 `json:"value"`
}

// Thresholds configuration for the panel
type ThresholdsConfig struct {
	// Thresholds can either be absolute (specific number) or percentage (relative to min or max, it will be values between 0 and 1).
	Mode ThresholdsMode `json:"mode"`

	// Must be sorted by 'value', first value is always -Infinity
	Steps []Threshold `json:"steps"`
}

// Thresholds can either be absolute (specific number) or percentage (relative to min or max, it will be values between 0 and 1).
type ThresholdsMode string

// Maps text values to a color or different display text and color.
// For example, you can configure a value mapping so that all instances of the value 10 appear as Perfection! rather than the number.
type ValueMap struct {
	// Map with <value_to_match>: ValueMappingResult. For example: { "10": { text: "Perfection!", color: "green" } }
	Options map[string]ValueMappingResult `json:"options"`
	Type    ValueMapType                  `json:"type"`
}

// ValueMapType defines model for ValueMap.Type.
type ValueMapType string

// Result used as replacement with text and color when the value matches
type ValueMappingResult struct {
	// Text to use when the value matches
	Color *string `json:"color,omitempty"`

	// Icon to display when the value matches. Only specific visualizations.
	Icon *string `json:"icon,omitempty"`

	// Position in the mapping array. Only used internally.
	Index *int32 `json:"index,omitempty"`

	// Text to display when the value matches
	Text *string `json:"text,omitempty"`
}

// Determine if the variable shows on dashboard
// Accepted values are 0 (show label and value), 1 (show value only), 2 (show nothing).
type VariableHide int

// Generic variable model to be used for all variable types
type VariableModel struct {
	// Format to use while fetching all values from data source, eg: wildcard, glob, regex, pipe, etc.
	AllFormat *string `json:"allFormat,omitempty"`

	// Option to be selected in a variable.
	Current *VariableOption `json:"current,omitempty"`

	// Ref to a DataSource instance
	Datasource *DataSourceRef `json:"datasource,omitempty"`

	// Description of variable. It can be defined but `null`.
	Description *string `json:"description,omitempty"`

	// Determine if the variable shows on dashboard
	// Accepted values are 0 (show label and value), 1 (show value only), 2 (show nothing).
	Hide VariableHide `json:"hide"`

	// Unique numeric identifier for the dashboard.
	Id string `json:"id"`

	// Optional display name
	Label *string `json:"label,omitempty"`

	// Whether multiple values can be selected or not from variable value list
	Multi *bool `json:"multi,omitempty"`

	// Name of variable
	Name string `json:"name"`

	// Options that can be selected for a variable.
	Options []VariableOption `json:"options,omitempty"`

	// Query used to fetch values for a variable
	Query *interface{} `json:"query,omitempty"`

	// Options to config when to refresh a variable
	// 0: Never refresh the variable
	// 1: Queries the data source every time the dashboard loads.
	// 2: Queries the data source when the dashboard time range changes.
	Refresh *VariableRefresh `json:"refresh,omitempty"`

	// Whether the variable value should be managed by URL query params or not
	SkipUrlSync bool `json:"skipUrlSync"`

	// Dashboard variable type
	Type VariableType `json:"type"`
}

// Option to be selected in a variable.
type VariableOption struct {
	// Whether the option is selected or not
	Selected *bool `json:"selected,omitempty"`

	// Text to be displayed for the option
	Text interface{} `json:"text"`

	// Value of the option
	Value interface{} `json:"value"`
}

// Options to config when to refresh a variable
// 0: Never refresh the variable
// 1: Queries the data source every time the dashboard loads.
// 2: Queries the data source when the dashboard time range changes.
type VariableRefresh int

// Dashboard variable type
type VariableType string
