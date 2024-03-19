// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     public/app/plugins/gen.go
// Using jennies:
//     PluginGoTypesJenny
//
// Run 'make gen-cue' from repository root to regenerate.

package dataquery

// Defines values for BucketAggregationType.
const (
	BucketAggregationTypeDateHistogram BucketAggregationType = "date_histogram"
	BucketAggregationTypeFilters       BucketAggregationType = "filters"
	BucketAggregationTypeGeohashGrid   BucketAggregationType = "geohash_grid"
	BucketAggregationTypeHistogram     BucketAggregationType = "histogram"
	BucketAggregationTypeNested        BucketAggregationType = "nested"
	BucketAggregationTypeTerms         BucketAggregationType = "terms"
)

// Defines values for BucketScriptType.
const (
	BucketScriptTypeBucketScript  BucketScriptType = "bucket_script"
	BucketScriptTypeCumulativeSum BucketScriptType = "cumulative_sum"
	BucketScriptTypeDerivative    BucketScriptType = "derivative"
	BucketScriptTypeMovingAvg     BucketScriptType = "moving_avg"
	BucketScriptTypeMovingFn      BucketScriptType = "moving_fn"
	BucketScriptTypeSerialDiff    BucketScriptType = "serial_diff"
)

// Defines values for ExtendedStatMetaType.
const (
	ExtendedStatMetaTypeAvg                     ExtendedStatMetaType = "avg"
	ExtendedStatMetaTypeCount                   ExtendedStatMetaType = "count"
	ExtendedStatMetaTypeMax                     ExtendedStatMetaType = "max"
	ExtendedStatMetaTypeMin                     ExtendedStatMetaType = "min"
	ExtendedStatMetaTypeStdDeviation            ExtendedStatMetaType = "std_deviation"
	ExtendedStatMetaTypeStdDeviationBoundsLower ExtendedStatMetaType = "std_deviation_bounds_lower"
	ExtendedStatMetaTypeStdDeviationBoundsUpper ExtendedStatMetaType = "std_deviation_bounds_upper"
	ExtendedStatMetaTypeSum                     ExtendedStatMetaType = "sum"
)

// Defines values for MaxType.
const (
	MaxTypeAvg           MaxType = "avg"
	MaxTypeBucketScript  MaxType = "bucket_script"
	MaxTypeCardinality   MaxType = "cardinality"
	MaxTypeCount         MaxType = "count"
	MaxTypeCumulativeSum MaxType = "cumulative_sum"
	MaxTypeDerivative    MaxType = "derivative"
	MaxTypeExtendedStats MaxType = "extended_stats"
	MaxTypeLogs          MaxType = "logs"
	MaxTypeMax           MaxType = "max"
	MaxTypeMin           MaxType = "min"
	MaxTypeMovingAvg     MaxType = "moving_avg"
	MaxTypeMovingFn      MaxType = "moving_fn"
	MaxTypePercentiles   MaxType = "percentiles"
	MaxTypeRate          MaxType = "rate"
	MaxTypeRawData       MaxType = "raw_data"
	MaxTypeRawDocument   MaxType = "raw_document"
	MaxTypeSerialDiff    MaxType = "serial_diff"
	MaxTypeSum           MaxType = "sum"
	MaxTypeTopMetrics    MaxType = "top_metrics"
)

// Defines values for MetricAggregationType.
const (
	MetricAggregationTypeAvg           MetricAggregationType = "avg"
	MetricAggregationTypeBucketScript  MetricAggregationType = "bucket_script"
	MetricAggregationTypeCardinality   MetricAggregationType = "cardinality"
	MetricAggregationTypeCount         MetricAggregationType = "count"
	MetricAggregationTypeCumulativeSum MetricAggregationType = "cumulative_sum"
	MetricAggregationTypeDerivative    MetricAggregationType = "derivative"
	MetricAggregationTypeExtendedStats MetricAggregationType = "extended_stats"
	MetricAggregationTypeLogs          MetricAggregationType = "logs"
	MetricAggregationTypeMax           MetricAggregationType = "max"
	MetricAggregationTypeMin           MetricAggregationType = "min"
	MetricAggregationTypeMovingAvg     MetricAggregationType = "moving_avg"
	MetricAggregationTypeMovingFn      MetricAggregationType = "moving_fn"
	MetricAggregationTypePercentiles   MetricAggregationType = "percentiles"
	MetricAggregationTypeRate          MetricAggregationType = "rate"
	MetricAggregationTypeRawData       MetricAggregationType = "raw_data"
	MetricAggregationTypeRawDocument   MetricAggregationType = "raw_document"
	MetricAggregationTypeSerialDiff    MetricAggregationType = "serial_diff"
	MetricAggregationTypeSum           MetricAggregationType = "sum"
	MetricAggregationTypeTopMetrics    MetricAggregationType = "top_metrics"
)

// Defines values for MovingAverageModel.
const (
	MovingAverageModelEwma        MovingAverageModel = "ewma"
	MovingAverageModelHolt        MovingAverageModel = "holt"
	MovingAverageModelHoltWinters MovingAverageModel = "holt_winters"
	MovingAverageModelLinear      MovingAverageModel = "linear"
	MovingAverageModelSimple      MovingAverageModel = "simple"
)

// Defines values for PipelineMetricAggregationType.
const (
	PipelineMetricAggregationTypeBucketScript  PipelineMetricAggregationType = "bucket_script"
	PipelineMetricAggregationTypeCumulativeSum PipelineMetricAggregationType = "cumulative_sum"
	PipelineMetricAggregationTypeDerivative    PipelineMetricAggregationType = "derivative"
	PipelineMetricAggregationTypeMovingAvg     PipelineMetricAggregationType = "moving_avg"
	PipelineMetricAggregationTypeMovingFn      PipelineMetricAggregationType = "moving_fn"
	PipelineMetricAggregationTypeSerialDiff    PipelineMetricAggregationType = "serial_diff"
)

// Defines values for RateType.
const (
	RateTypeAvg           RateType = "avg"
	RateTypeBucketScript  RateType = "bucket_script"
	RateTypeCardinality   RateType = "cardinality"
	RateTypeCount         RateType = "count"
	RateTypeCumulativeSum RateType = "cumulative_sum"
	RateTypeDerivative    RateType = "derivative"
	RateTypeExtendedStats RateType = "extended_stats"
	RateTypeLogs          RateType = "logs"
	RateTypeMax           RateType = "max"
	RateTypeMin           RateType = "min"
	RateTypeMovingAvg     RateType = "moving_avg"
	RateTypeMovingFn      RateType = "moving_fn"
	RateTypePercentiles   RateType = "percentiles"
	RateTypeRate          RateType = "rate"
	RateTypeRawData       RateType = "raw_data"
	RateTypeRawDocument   RateType = "raw_document"
	RateTypeSerialDiff    RateType = "serial_diff"
	RateTypeSum           RateType = "sum"
	RateTypeTopMetrics    RateType = "top_metrics"
)

// Defines values for TermsOrder.
const (
	TermsOrderAsc  TermsOrder = "asc"
	TermsOrderDesc TermsOrder = "desc"
)

// Average defines model for Average.
type Average struct {
	Field    *string `json:"field,omitempty"`
	Hide     *bool   `json:"hide,omitempty"`
	Id       *string `json:"id,omitempty"`
	Settings *struct {
		Script *any `json:"script,omitempty"`
	} `json:"settings,omitempty"`
	Type *MetricAggregationType `json:"type,omitempty"`
}

// BaseBucketAggregation defines model for BaseBucketAggregation.
type BaseBucketAggregation struct {
	Id       string                `json:"id"`
	Settings *any                  `json:"settings,omitempty"`
	Type     BucketAggregationType `json:"type"`
}

// BaseMetricAggregation defines model for BaseMetricAggregation.
type BaseMetricAggregation struct {
	Hide *bool                 `json:"hide,omitempty"`
	Id   string                `json:"id"`
	Type MetricAggregationType `json:"type"`
}

// BaseMovingAverageModelSettings defines model for BaseMovingAverageModelSettings.
type BaseMovingAverageModelSettings struct {
	Model   MovingAverageModel `json:"model"`
	Predict string             `json:"predict"`
	Window  string             `json:"window"`
}

// BasePipelineMetricAggregation defines model for BasePipelineMetricAggregation.
type BasePipelineMetricAggregation struct {
	Field       *string                `json:"field,omitempty"`
	Hide        *bool                  `json:"hide,omitempty"`
	Id          *string                `json:"id,omitempty"`
	PipelineAgg *string                `json:"pipelineAgg,omitempty"`
	Type        *MetricAggregationType `json:"type,omitempty"`
}

// BucketAggregationType defines model for BucketAggregationType.
type BucketAggregationType string

// BucketAggregationWithField defines model for BucketAggregationWithField.
type BucketAggregationWithField struct {
	Field    *string                `json:"field,omitempty"`
	Id       *string                `json:"id,omitempty"`
	Settings *any                   `json:"settings,omitempty"`
	Type     *BucketAggregationType `json:"type,omitempty"`
}

// BucketScript defines model for BucketScript.
type BucketScript struct {
	PipelineVariables []PipelineVariable `json:"pipelineVariables,omitempty"`
	Settings          *struct {
		Script *any `json:"script,omitempty"`
	} `json:"settings,omitempty"`
	Type *BucketScriptType `json:"type,omitempty"`
}

// BucketScriptType defines model for BucketScript.Type.
type BucketScriptType string

// Count defines model for Count.
type Count struct {
	Hide *bool                  `json:"hide,omitempty"`
	Id   *string                `json:"id,omitempty"`
	Type *MetricAggregationType `json:"type,omitempty"`
}

// CumulativeSum defines model for CumulativeSum.
type CumulativeSum struct {
	PipelineAgg *string `json:"pipelineAgg,omitempty"`
	Settings    *struct {
		Format *string `json:"format,omitempty"`
	} `json:"settings,omitempty"`
	Type *PipelineMetricAggregationType `json:"type,omitempty"`
}

// These are the common properties available to all queries in all datasources.
// Specific implementations will *extend* this interface, adding the required
// properties for the given context.
type DataQuery struct {
	// For mixed data sources the selected datasource is on the query level.
	// For non mixed scenarios this is undefined.
	// TODO find a better way to do this ^ that's friendly to schema
	// TODO this shouldn't be unknown but DataSourceRef | null
	Datasource *any `json:"datasource,omitempty"`

	// Hide true if query is disabled (ie should not be returned to the dashboard)
	// Note this does not always imply that the query should not be executed since
	// the results from a hidden query may be used as the input to other queries (SSE etc)
	Hide *bool `json:"hide,omitempty"`

	// Specify the query flavor
	// TODO make this required and give it a default
	QueryType *string `json:"queryType,omitempty"`

	// A unique identifier for the query within the list of targets.
	// In server side expressions, the refId is used as a variable name to identify results.
	// By default, the UI will assign A->Z; however setting meaningful names may be useful.
	RefId string `json:"refId"`
}

// DateHistogram defines model for DateHistogram.
type DateHistogram struct {
	Field    *string                `json:"field,omitempty"`
	Id       *string                `json:"id,omitempty"`
	Settings *any                   `json:"settings,omitempty"`
	Type     *BucketAggregationType `json:"type,omitempty"`
}

// DateHistogramSettings defines model for DateHistogramSettings.
type DateHistogramSettings struct {
	Interval    *string `json:"interval,omitempty"`
	MinDocCount *string `json:"min_doc_count,omitempty"`
	Offset      *string `json:"offset,omitempty"`
	TimeZone    *string `json:"timeZone,omitempty"`
	TrimEdges   *string `json:"trimEdges,omitempty"`
}

// Derivative defines model for Derivative.
type Derivative struct {
	PipelineAgg *string `json:"pipelineAgg,omitempty"`
	Settings    *struct {
		Unit *string `json:"unit,omitempty"`
	} `json:"settings,omitempty"`
	Type *PipelineMetricAggregationType `json:"type,omitempty"`
}

// ElasticsearchDataQuery defines model for ElasticsearchDataQuery.
type ElasticsearchDataQuery struct {
	// Alias pattern
	Alias *string `json:"alias,omitempty"`

	// List of bucket aggregations
	BucketAggs []any `json:"bucketAggs,omitempty"`

	// For mixed data sources the selected datasource is on the query level.
	// For non mixed scenarios this is undefined.
	// TODO find a better way to do this ^ that's friendly to schema
	// TODO this shouldn't be unknown but DataSourceRef | null
	Datasource *any `json:"datasource,omitempty"`

	// Hide true if query is disabled (ie should not be returned to the dashboard)
	// Note this does not always imply that the query should not be executed since
	// the results from a hidden query may be used as the input to other queries (SSE etc)
	Hide *bool `json:"hide,omitempty"`

	// List of metric aggregations
	Metrics []any `json:"metrics,omitempty"`

	// Lucene query
	Query *string `json:"query,omitempty"`

	// Specify the query flavor
	// TODO make this required and give it a default
	QueryType *string `json:"queryType,omitempty"`

	// A unique identifier for the query within the list of targets.
	// In server side expressions, the refId is used as a variable name to identify results.
	// By default, the UI will assign A->Z; however setting meaningful names may be useful.
	RefId *string `json:"refId,omitempty"`

	// Name of time field
	TimeField *string `json:"timeField,omitempty"`
}

// ExtendedStat defines model for ExtendedStat.
type ExtendedStat struct {
	Label string               `json:"label"`
	Value ExtendedStatMetaType `json:"value"`
}

// ExtendedStatMetaType defines model for ExtendedStatMetaType.
type ExtendedStatMetaType string

// ExtendedStats defines model for ExtendedStats.
type ExtendedStats struct {
	Field    *string        `json:"field,omitempty"`
	Hide     *bool          `json:"hide,omitempty"`
	Id       *string        `json:"id,omitempty"`
	Meta     map[string]any `json:"meta,omitempty"`
	Settings *struct {
		Script *any `json:"script,omitempty"`
	} `json:"settings,omitempty"`
	Type *MetricAggregationType `json:"type,omitempty"`
}

// Filter defines model for Filter.
type Filter struct {
	Label string `json:"label"`
	Query string `json:"query"`
}

// Filters defines model for Filters.
type Filters struct {
	Id       *string                `json:"id,omitempty"`
	Settings *any                   `json:"settings,omitempty"`
	Type     *BucketAggregationType `json:"type,omitempty"`
}

// FiltersSettings defines model for FiltersSettings.
type FiltersSettings struct {
	Filters []Filter `json:"filters,omitempty"`
}

// GeoHashGrid defines model for GeoHashGrid.
type GeoHashGrid struct {
	Field    *string                `json:"field,omitempty"`
	Id       *string                `json:"id,omitempty"`
	Settings *any                   `json:"settings,omitempty"`
	Type     *BucketAggregationType `json:"type,omitempty"`
}

// GeoHashGridSettings defines model for GeoHashGridSettings.
type GeoHashGridSettings struct {
	Precision *string `json:"precision,omitempty"`
}

// Histogram defines model for Histogram.
type Histogram struct {
	Field    *string                `json:"field,omitempty"`
	Id       *string                `json:"id,omitempty"`
	Settings *any                   `json:"settings,omitempty"`
	Type     *BucketAggregationType `json:"type,omitempty"`
}

// HistogramSettings defines model for HistogramSettings.
type HistogramSettings struct {
	Interval    *string `json:"interval,omitempty"`
	MinDocCount *string `json:"min_doc_count,omitempty"`
}

// Logs defines model for Logs.
type Logs struct {
	Hide     *bool   `json:"hide,omitempty"`
	Id       *string `json:"id,omitempty"`
	Settings *struct {
		Limit *string `json:"limit,omitempty"`
	} `json:"settings,omitempty"`
	Type *MetricAggregationType `json:"type,omitempty"`
}

// Max defines model for Max.
type Max struct {
	Field    *string `json:"field,omitempty"`
	Settings *struct {
		Script *any `json:"script,omitempty"`
	} `json:"settings,omitempty"`
	Type *MaxType `json:"type,omitempty"`
}

// MaxType defines model for Max.Type.
type MaxType string

// MetricAggregationType defines model for MetricAggregationType.
type MetricAggregationType string

// MetricAggregationWithField defines model for MetricAggregationWithField.
type MetricAggregationWithField struct {
	Field *string                `json:"field,omitempty"`
	Hide  *bool                  `json:"hide,omitempty"`
	Id    *string                `json:"id,omitempty"`
	Type  *MetricAggregationType `json:"type,omitempty"`
}

// MetricAggregationWithInlineScript defines model for MetricAggregationWithInlineScript.
type MetricAggregationWithInlineScript struct {
	Hide     *bool   `json:"hide,omitempty"`
	Id       *string `json:"id,omitempty"`
	Settings *struct {
		Script *any `json:"script,omitempty"`
	} `json:"settings,omitempty"`
	Type *MetricAggregationType `json:"type,omitempty"`
}

// MetricAggregationWithMissingSupport defines model for MetricAggregationWithMissingSupport.
type MetricAggregationWithMissingSupport struct {
	Hide     *bool   `json:"hide,omitempty"`
	Id       *string `json:"id,omitempty"`
	Settings *struct {
		Missing *string `json:"missing,omitempty"`
	} `json:"settings,omitempty"`
	Type *MetricAggregationType `json:"type,omitempty"`
}

// Min defines model for Min.
type Min struct {
	Field    *string `json:"field,omitempty"`
	Hide     *bool   `json:"hide,omitempty"`
	Id       *string `json:"id,omitempty"`
	Settings *struct {
		Script *any `json:"script,omitempty"`
	} `json:"settings,omitempty"`
	Type *MetricAggregationType `json:"type,omitempty"`
}

// #MovingAverage's settings are overridden in types.ts
type MovingAverage struct {
	PipelineAgg *string                        `json:"pipelineAgg,omitempty"`
	Settings    map[string]any                 `json:"settings,omitempty"`
	Type        *PipelineMetricAggregationType `json:"type,omitempty"`
}

// MovingAverageEWMAModelSettings defines model for MovingAverageEWMAModelSettings.
type MovingAverageEWMAModelSettings struct {
	Minimize *bool               `json:"minimize,omitempty"`
	Model    *MovingAverageModel `json:"model,omitempty"`
	Predict  *string             `json:"predict,omitempty"`
	Settings *struct {
		Alpha *string `json:"alpha,omitempty"`
	} `json:"settings,omitempty"`
	Window *string `json:"window,omitempty"`
}

// MovingAverageHoltModelSettings defines model for MovingAverageHoltModelSettings.
type MovingAverageHoltModelSettings struct {
	Minimize *bool               `json:"minimize,omitempty"`
	Model    *MovingAverageModel `json:"model,omitempty"`
	Predict  *string             `json:"predict,omitempty"`
	Settings *struct {
		Alpha *string `json:"alpha,omitempty"`
		Beta  *string `json:"beta,omitempty"`
	} `json:"settings,omitempty"`
	Window *string `json:"window,omitempty"`
}

// MovingAverageHoltWintersModelSettings defines model for MovingAverageHoltWintersModelSettings.
type MovingAverageHoltWintersModelSettings struct {
	Minimize *bool               `json:"minimize,omitempty"`
	Model    *MovingAverageModel `json:"model,omitempty"`
	Predict  *string             `json:"predict,omitempty"`
	Settings *struct {
		Alpha  *string `json:"alpha,omitempty"`
		Beta   *string `json:"beta,omitempty"`
		Gamma  *string `json:"gamma,omitempty"`
		Pad    *bool   `json:"pad,omitempty"`
		Period *string `json:"period,omitempty"`
	} `json:"settings,omitempty"`
	Window *string `json:"window,omitempty"`
}

// MovingAverageLinearModelSettings defines model for MovingAverageLinearModelSettings.
type MovingAverageLinearModelSettings struct {
	Model   *MovingAverageModel `json:"model,omitempty"`
	Predict *string             `json:"predict,omitempty"`
	Window  *string             `json:"window,omitempty"`
}

// MovingAverageModel defines model for MovingAverageModel.
type MovingAverageModel string

// MovingAverageModelOption defines model for MovingAverageModelOption.
type MovingAverageModelOption struct {
	Label string             `json:"label"`
	Value MovingAverageModel `json:"value"`
}

// MovingAverageSimpleModelSettings defines model for MovingAverageSimpleModelSettings.
type MovingAverageSimpleModelSettings struct {
	Model   *MovingAverageModel `json:"model,omitempty"`
	Predict *string             `json:"predict,omitempty"`
	Window  *string             `json:"window,omitempty"`
}

// MovingFunction defines model for MovingFunction.
type MovingFunction struct {
	PipelineAgg *string `json:"pipelineAgg,omitempty"`
	Settings    *struct {
		Script *any    `json:"script,omitempty"`
		Shift  *string `json:"shift,omitempty"`
		Window *string `json:"window,omitempty"`
	} `json:"settings,omitempty"`
	Type *PipelineMetricAggregationType `json:"type,omitempty"`
}

// Nested defines model for Nested.
type Nested struct {
	Field    *string                `json:"field,omitempty"`
	Id       *string                `json:"id,omitempty"`
	Settings *any                   `json:"settings,omitempty"`
	Type     *BucketAggregationType `json:"type,omitempty"`
}

// Percentiles defines model for Percentiles.
type Percentiles struct {
	Field    *string `json:"field,omitempty"`
	Hide     *bool   `json:"hide,omitempty"`
	Id       *string `json:"id,omitempty"`
	Settings *struct {
		Script *any `json:"script,omitempty"`
	} `json:"settings,omitempty"`
	Type *MetricAggregationType `json:"type,omitempty"`
}

// PipelineMetricAggregationType defines model for PipelineMetricAggregationType.
type PipelineMetricAggregationType string

// PipelineMetricAggregationWithMultipleBucketPaths defines model for PipelineMetricAggregationWithMultipleBucketPaths.
type PipelineMetricAggregationWithMultipleBucketPaths struct {
	Hide              *bool                  `json:"hide,omitempty"`
	Id                *string                `json:"id,omitempty"`
	PipelineVariables []PipelineVariable     `json:"pipelineVariables,omitempty"`
	Type              *MetricAggregationType `json:"type,omitempty"`
}

// PipelineVariable defines model for PipelineVariable.
type PipelineVariable struct {
	Name        string `json:"name"`
	PipelineAgg string `json:"pipelineAgg"`
}

// Rate defines model for Rate.
type Rate struct {
	Field    *string `json:"field,omitempty"`
	Settings *struct {
		Mode *string `json:"mode,omitempty"`
		Unit *string `json:"unit,omitempty"`
	} `json:"settings,omitempty"`
	Type *RateType `json:"type,omitempty"`
}

// RateType defines model for Rate.Type.
type RateType string

// RawData defines model for RawData.
type RawData struct {
	Hide     *bool   `json:"hide,omitempty"`
	Id       *string `json:"id,omitempty"`
	Settings *struct {
		Size *string `json:"size,omitempty"`
	} `json:"settings,omitempty"`
	Type *MetricAggregationType `json:"type,omitempty"`
}

// RawDocument defines model for RawDocument.
type RawDocument struct {
	Hide     *bool   `json:"hide,omitempty"`
	Id       *string `json:"id,omitempty"`
	Settings *struct {
		Size *string `json:"size,omitempty"`
	} `json:"settings,omitempty"`
	Type *MetricAggregationType `json:"type,omitempty"`
}

// SerialDiff defines model for SerialDiff.
type SerialDiff struct {
	Field       *string `json:"field,omitempty"`
	Hide        *bool   `json:"hide,omitempty"`
	Id          *string `json:"id,omitempty"`
	PipelineAgg *string `json:"pipelineAgg,omitempty"`
	Settings    *struct {
		Lag *string `json:"lag,omitempty"`
	} `json:"settings,omitempty"`
	Type *MetricAggregationType `json:"type,omitempty"`
}

// Sum defines model for Sum.
type Sum struct {
	Field    *string `json:"field,omitempty"`
	Hide     *bool   `json:"hide,omitempty"`
	Id       *string `json:"id,omitempty"`
	Settings *struct {
		Script *any `json:"script,omitempty"`
	} `json:"settings,omitempty"`
	Type *MetricAggregationType `json:"type,omitempty"`
}

// Terms defines model for Terms.
type Terms struct {
	Field    *string                `json:"field,omitempty"`
	Id       *string                `json:"id,omitempty"`
	Settings *any                   `json:"settings,omitempty"`
	Type     *BucketAggregationType `json:"type,omitempty"`
}

// TermsOrder defines model for TermsOrder.
type TermsOrder string

// TermsSettings defines model for TermsSettings.
type TermsSettings struct {
	MinDocCount *string     `json:"min_doc_count,omitempty"`
	Missing     *string     `json:"missing,omitempty"`
	Order       *TermsOrder `json:"order,omitempty"`
	OrderBy     *string     `json:"orderBy,omitempty"`
	Size        *string     `json:"size,omitempty"`
}

// TopMetrics defines model for TopMetrics.
type TopMetrics struct {
	Hide     *bool   `json:"hide,omitempty"`
	Id       *string `json:"id,omitempty"`
	Settings *struct {
		Metrics []string `json:"metrics,omitempty"`
		Order   *string  `json:"order,omitempty"`
		OrderBy *string  `json:"orderBy,omitempty"`
	} `json:"settings,omitempty"`
	Type *MetricAggregationType `json:"type,omitempty"`
}

// UniqueCount defines model for UniqueCount.
type UniqueCount struct {
	Field    *string `json:"field,omitempty"`
	Hide     *bool   `json:"hide,omitempty"`
	Id       *string `json:"id,omitempty"`
	Settings *struct {
		Missing            *string `json:"missing,omitempty"`
		PrecisionThreshold *string `json:"precision_threshold,omitempty"`
	} `json:"settings,omitempty"`
	Type *MetricAggregationType `json:"type,omitempty"`
}
