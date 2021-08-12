package grafanaschema

TableCellDisplayMode: {
  Auto: "auto",
  ColorText: "color-text",
  ColorBackground: "color-background",
  GradientGauge: "gradient-gauge",
  LcdGauge: "lcd-gauge",
  JSONView: "json-view",
  BasicGauge: "basic",
  Image: "image",
} @cuetsy(targetType="enum")

TableFieldOptions: {
  width?: number
  align: FieldTextAlignment | *"auto"
  displayMode: TableCellDisplayMode | *"auto"
  hidden?: bool  // ?? default is missing or false ??
} @cuetsy(targetType="interface")

TableSortByFieldState: {
  displayName: string
  desc?: bool
} @cuetsy(targetType="interface")

TooltipDisplayMode:          "single" | "multi" | "none"          @cuetsy(targetType="enum")
FieldTextAlignment:   "auto" | "left" | "right" | "center" @cuetsy(targetType="type")
AxisPlacement:        "auto" | "top" | "right" | "bottom" | "left" | "hidden" @cuetsy(targetType="enum")
PointVisibility:      "auto" | "never" | "always"                             @cuetsy(targetType="enum")
DrawStyle:            "line" | "bars" | "points"                              @cuetsy(targetType="enum")
LineInterpolation:    "linear" | "smooth" | "stepBefore" | "stepAfter"        @cuetsy(targetType="enum")
ScaleDistribution:    "linear" | "log"                                        @cuetsy(targetType="enum")
GraphGradientMode:    "none" | "opacity" | "hue" | "scheme"                   @cuetsy(targetType="enum")
LineStyle: {
	fill?: "solid" | "dash" | "dot" | "square"
	dash?: [...number]
} @cuetsy(targetType="interface")
LineConfig: {
	lineColor?:         string
	lineWidth?:         number
	lineInterpolation?: LineInterpolation
	lineStyle?:         LineStyle
	spanNulls?:         bool | number
} @cuetsy(targetType="interface")
FillConfig: {
	fillColor?:   string
	fillOpacity?: number
	fillBelowTo?: string
} @cuetsy(targetType="interface")
PointsConfig: {
	showPoints?:  PointVisibility
	pointSize?:   number
	pointColor?:  string
	pointSymbol?: string
} @cuetsy(targetType="interface")
ScaleDistributionConfig: {
	type: ScaleDistribution
	log?: number
} @cuetsy(targetType="interface")
AxisConfig: {
	axisPlacement?:     AxisPlacement
	axisLabel?:         string
	axisWidth?:         number
	axisSoftMin?:       number
	axisSoftMax?:       number
	scaleDistribution?: ScaleDistributionConfig
} @cuetsy(targetType="interface")
HideSeriesConfig: {
	tooltip: bool
	legend:  bool
	viz:   bool
} @cuetsy(targetType="interface")
LegendPlacement:   "bottom" | "right"          @cuetsy(targetType="type")
LegendDisplayMode: "list" | "table" | "hidden" @cuetsy(targetType="enum")
TableFieldOptions: {
	width?:      number
	align:       FieldTextAlignment | *"auto"
	displayMode: TableCellDisplayMode | *"auto"
	hidden?:     bool
} @cuetsy(targetType="interface")
GraphFieldConfig: LineConfig & FillConfig & PointsConfig & AxisConfig & {
	drawStyle?:    DrawStyle
	gradientMode?: GraphGradientMode
	hideFrom?:     HideSeriesConfig
} @cuetsy(targetType="interface")
VizLegendOptions: {
	displayMode: LegendDisplayMode
	placement:   LegendPlacement
	asTable: bool | *false
	isVisible: bool | *false
	calcs: [...string]
} @cuetsy(targetType="interface")
VizTooltipOptions: {
	mode: TooltipDisplayMode
} @cuetsy(targetType="interface")
// TODO copy back to appropriate place
SingleStatBaseOptions: {
	OptionsWithTextFormatting
	reduceOptions: ReduceDataOptions
	orientation: VizOrientation
} @cuetsy(targetType="interface")
// TODO copy back to appropriate place
ReduceDataOptions: {
  // If true show each row value
  values?: bool
  // if showing all values limit
  limit?: number
  // When !values, pick one value for the whole field
  calcs: [...string]
  // Which fields to show.  By default this is only numeric fields
  fields?: string
} @cuetsy(targetType="interface")
// TODO copy back to appropriate place
VizOrientation: "auto" | "vertical" | "horizontal" @cuetsy(targetType="enum")
// TODO copy back to appropriate place
OptionsWithTextFormatting: {
	text?: VizTextDisplayOptions
} @cuetsy(targetType="interface")
// TODO copy back to appropriate place
VizTextDisplayOptions: {
	// Explicit title text size
	titleSize?: number
	// Explicit value text size
	valueSize?: number
} @cuetsy(targetType="interface")
// TODO copy back to appropriate place
BigValueColorMode: "value" | "background" | "none" @cuetsy(targetType="enum")
// TODO copy back to appropriate place
BigValueGraphMode: "none" | "line" | "area" @cuetsy(targetType="enum")
// TODO copy back to appropriate place
BigValueJustifyMode: "auto" | "center" @cuetsy(targetType="enum")
// TODO copy back to appropriate place
// TODO does cuetsy handle underscores the expected way?
BigValueTextMode: "auto" | "value" | "value_and_name" | "name" | "none"  @cuetsy(targetType="enum")
// TODO copy back to appropriate place
BarGaugeDisplayMode: "basic" | "lcd" | "gradient" @cuetsy(targetType="enum") 