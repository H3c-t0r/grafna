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

// Defines values for TermsOrder.
const (
	TermsOrderAsc  TermsOrder = "asc"
	TermsOrderDesc TermsOrder = "desc"
)

// Average defines model for Average.
type Average struct {
	MetricAggregationWithField
	MetricAggregationWithInlineScript
	MetricAggregationWithMissingSupport
	Hide     *bool  `json:"hide,omitempty"`
	Id       string `json:"id"`
	Settings *struct {
		Missing *string `json:"missing,omitempty"`
		Script  *any    `json:"script,omitempty"`
	} `json:"settings,omitempty"`
	Type MetricAggregationType `json:"type"`
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
	MetricAggregationWithField
	Hide        *bool                 `json:"hide,omitempty"`
	Id          string                `json:"id"`
	PipelineAgg *string               `json:"pipelineAgg,omitempty"`
	Type        MetricAggregationType `json:"type"`
}

// BucketAggregationType defines model for BucketAggregationType.
type BucketAggregationType string

// BucketAggregationWithField defines model for BucketAggregationWithField.
type BucketAggregationWithField struct {
	BaseBucketAggregation
	Field *string `json:"field,omitempty"`
}

// BucketScript defines model for BucketScript.
type BucketScript struct {
	PipelineMetricAggregationWithMultipleBucketPaths
	Hide     *bool  `json:"hide,omitempty"`
	Id       string `json:"id"`
	Settings *struct {
		Script *any `json:"script,omitempty"`
	} `json:"settings,omitempty"`
	Type MetricAggregationType `json:"type"`
}

// Count defines model for Count.
type Count struct {
	BaseMetricAggregation
	Hide *bool                 `json:"hide,omitempty"`
	Id   string                `json:"id"`
	Type MetricAggregationType `json:"type"`
}

// CumulativeSum defines model for CumulativeSum.
type CumulativeSum struct {
	BasePipelineMetricAggregation
	Hide     *bool  `json:"hide,omitempty"`
	Id       string `json:"id"`
	Settings *struct {
		Format *string `json:"format,omitempty"`
	} `json:"settings,omitempty"`
	Type MetricAggregationType `json:"type"`
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

	// If hide is set to true, Grafana will filter out the response(s) associated with this query before returning it to the panel.
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
	BucketAggregationWithField
	Id       string                `json:"id"`
	Settings *any                  `json:"settings,omitempty"`
	Type     BucketAggregationType `json:"type"`
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
	BasePipelineMetricAggregation
	Hide     *bool  `json:"hide,omitempty"`
	Id       string `json:"id"`
	Settings *struct {
		Unit *string `json:"unit,omitempty"`
	} `json:"settings,omitempty"`
	Type MetricAggregationType `json:"type"`
}

// ElasticsearchDataQuery defines model for ElasticsearchDataQuery.
type ElasticsearchDataQuery struct {
	// DataQuery These are the common properties available to all queries in all datasources.
	// Specific implementations will *extend* this interface, adding the required
	// properties for the given context.
	DataQuery

	// Alias pattern
	Alias *string `json:"alias,omitempty"`

	// List of bucket aggregations
	BucketAggs []any `json:"bucketAggs,omitempty"`

	// List of metric aggregations
	Metrics []any `json:"metrics,omitempty"`

	// Lucene query
	Query *string `json:"query,omitempty"`

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
	MetricAggregationWithField
	MetricAggregationWithInlineScript
	Hide     *bool          `json:"hide,omitempty"`
	Id       string         `json:"id"`
	Meta     map[string]any `json:"meta,omitempty"`
	Settings *struct {
		Missing *string `json:"missing,omitempty"`
		Script  *any    `json:"script,omitempty"`
		Sigma   *string `json:"sigma,omitempty"`
	} `json:"settings,omitempty"`
	Type MetricAggregationType `json:"type"`
}

// Filter defines model for Filter.
type Filter struct {
	Label string `json:"label"`
	Query string `json:"query"`
}

// Filters defines model for Filters.
type Filters struct {
	BaseBucketAggregation
	Id       string                `json:"id"`
	Settings *any                  `json:"settings,omitempty"`
	Type     BucketAggregationType `json:"type"`
}

// FiltersSettings defines model for FiltersSettings.
type FiltersSettings struct {
	Filters []Filter `json:"filters,omitempty"`
}

// GeoHashGrid defines model for GeoHashGrid.
type GeoHashGrid struct {
	BucketAggregationWithField
	Id       string                `json:"id"`
	Settings *any                  `json:"settings,omitempty"`
	Type     BucketAggregationType `json:"type"`
}

// GeoHashGridSettings defines model for GeoHashGridSettings.
type GeoHashGridSettings struct {
	Precision *string `json:"precision,omitempty"`
}

// Histogram defines model for Histogram.
type Histogram struct {
	BucketAggregationWithField
	Id       string                `json:"id"`
	Settings *any                  `json:"settings,omitempty"`
	Type     BucketAggregationType `json:"type"`
}

// HistogramSettings defines model for HistogramSettings.
type HistogramSettings struct {
	Interval    *string `json:"interval,omitempty"`
	MinDocCount *string `json:"min_doc_count,omitempty"`
}

// Logs defines model for Logs.
type Logs struct {
	BaseMetricAggregation
	Hide     *bool  `json:"hide,omitempty"`
	Id       string `json:"id"`
	Settings *struct {
		Limit *string `json:"limit,omitempty"`
	} `json:"settings,omitempty"`
	Type MetricAggregationType `json:"type"`
}

// Max defines model for Max.
type Max struct {
	MetricAggregationWithField
	MetricAggregationWithInlineScript
	Hide     *bool  `json:"hide,omitempty"`
	Id       string `json:"id"`
	Settings *struct {
		Missing *string `json:"missing,omitempty"`
		Script  *any    `json:"script,omitempty"`
	} `json:"settings,omitempty"`
	Type MetricAggregationType `json:"type"`
}

// MetricAggregationType defines model for MetricAggregationType.
type MetricAggregationType string

// MetricAggregationWithField defines model for MetricAggregationWithField.
type MetricAggregationWithField struct {
	BaseMetricAggregation
	Field *string `json:"field,omitempty"`
}

// MetricAggregationWithInlineScript defines model for MetricAggregationWithInlineScript.
type MetricAggregationWithInlineScript struct {
	BaseMetricAggregation
	Settings *struct {
		Script *any `json:"script,omitempty"`
	} `json:"settings,omitempty"`
}

// MetricAggregationWithMissingSupport defines model for MetricAggregationWithMissingSupport.
type MetricAggregationWithMissingSupport struct {
	BaseMetricAggregation
	Settings *struct {
		Missing *string `json:"missing,omitempty"`
	} `json:"settings,omitempty"`
}

// Min defines model for Min.
type Min struct {
	MetricAggregationWithField
	MetricAggregationWithInlineScript
	Hide     *bool  `json:"hide,omitempty"`
	Id       string `json:"id"`
	Settings *struct {
		Missing *string `json:"missing,omitempty"`
		Script  *any    `json:"script,omitempty"`
	} `json:"settings,omitempty"`
	Type MetricAggregationType `json:"type"`
}

// MovingAverage defines model for MovingAverage.
type MovingAverage struct {
	BasePipelineMetricAggregation
	Hide     *bool                 `json:"hide,omitempty"`
	Id       string                `json:"id"`
	Settings map[string]any        `json:"settings,omitempty"`
	Type     MetricAggregationType `json:"type"`
}

// MovingAverageEWMAModelSettings defines model for MovingAverageEWMAModelSettings.
type MovingAverageEWMAModelSettings struct {
	BaseMovingAverageModelSettings
	Minimize bool               `json:"minimize"`
	Model    MovingAverageModel `json:"model"`
	Predict  string             `json:"predict"`
	Settings *struct {
		Alpha *string `json:"alpha,omitempty"`
	} `json:"settings,omitempty"`
	Window string `json:"window"`
}

// MovingAverageHoltModelSettings defines model for MovingAverageHoltModelSettings.
type MovingAverageHoltModelSettings struct {
	BaseMovingAverageModelSettings
	Minimize bool               `json:"minimize"`
	Model    MovingAverageModel `json:"model"`
	Predict  string             `json:"predict"`
	Settings struct {
		Alpha *string `json:"alpha,omitempty"`
		Beta  *string `json:"beta,omitempty"`
	} `json:"settings"`
	Window string `json:"window"`
}

// MovingAverageHoltWintersModelSettings defines model for MovingAverageHoltWintersModelSettings.
type MovingAverageHoltWintersModelSettings struct {
	BaseMovingAverageModelSettings
	Minimize bool               `json:"minimize"`
	Model    MovingAverageModel `json:"model"`
	Predict  string             `json:"predict"`
	Settings struct {
		Alpha  *string `json:"alpha,omitempty"`
		Beta   *string `json:"beta,omitempty"`
		Gamma  *string `json:"gamma,omitempty"`
		Pad    *bool   `json:"pad,omitempty"`
		Period *string `json:"period,omitempty"`
	} `json:"settings"`
	Window string `json:"window"`
}

// MovingAverageLinearModelSettings defines model for MovingAverageLinearModelSettings.
type MovingAverageLinearModelSettings struct {
	BaseMovingAverageModelSettings
	Model   MovingAverageModel `json:"model"`
	Predict string             `json:"predict"`
	Window  string             `json:"window"`
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
	BaseMovingAverageModelSettings
	Model   MovingAverageModel `json:"model"`
	Predict string             `json:"predict"`
	Window  string             `json:"window"`
}

// MovingFunction defines model for MovingFunction.
type MovingFunction struct {
	BasePipelineMetricAggregation
	Hide     *bool  `json:"hide,omitempty"`
	Id       string `json:"id"`
	Settings *struct {
		Script *any    `json:"script,omitempty"`
		Shift  *string `json:"shift,omitempty"`
		Window *string `json:"window,omitempty"`
	} `json:"settings,omitempty"`
	Type MetricAggregationType `json:"type"`
}

// Nested defines model for Nested.
type Nested struct {
	BucketAggregationWithField
	Id       string                `json:"id"`
	Settings *any                  `json:"settings,omitempty"`
	Type     BucketAggregationType `json:"type"`
}

// Percentiles defines model for Percentiles.
type Percentiles struct {
	MetricAggregationWithField
	MetricAggregationWithInlineScript
	Hide     *bool  `json:"hide,omitempty"`
	Id       string `json:"id"`
	Settings *struct {
		Missing  *string  `json:"missing,omitempty"`
		Percents []string `json:"percents,omitempty"`
		Script   *any     `json:"script,omitempty"`
	} `json:"settings,omitempty"`
	Type MetricAggregationType `json:"type"`
}

// PipelineMetricAggregationType defines model for PipelineMetricAggregationType.
type PipelineMetricAggregationType string

// PipelineMetricAggregationWithMultipleBucketPaths defines model for PipelineMetricAggregationWithMultipleBucketPaths.
type PipelineMetricAggregationWithMultipleBucketPaths struct {
	BaseMetricAggregation
	PipelineVariables []PipelineVariable `json:"pipelineVariables,omitempty"`
}

// PipelineVariable defines model for PipelineVariable.
type PipelineVariable struct {
	Name        string `json:"name"`
	PipelineAgg string `json:"pipelineAgg"`
}

// Rate defines model for Rate.
type Rate struct {
	MetricAggregationWithField
	Hide     *bool  `json:"hide,omitempty"`
	Id       string `json:"id"`
	Settings *struct {
		Mode *string `json:"mode,omitempty"`
		Unit *string `json:"unit,omitempty"`
	} `json:"settings,omitempty"`
	Type MetricAggregationType `json:"type"`
}

// RawData defines model for RawData.
type RawData struct {
	BaseMetricAggregation
	Hide     *bool  `json:"hide,omitempty"`
	Id       string `json:"id"`
	Settings *struct {
		Size *string `json:"size,omitempty"`
	} `json:"settings,omitempty"`
	Type MetricAggregationType `json:"type"`
}

// RawDocument defines model for RawDocument.
type RawDocument struct {
	BaseMetricAggregation
	Hide     *bool  `json:"hide,omitempty"`
	Id       string `json:"id"`
	Settings *struct {
		Size *string `json:"size,omitempty"`
	} `json:"settings,omitempty"`
	Type MetricAggregationType `json:"type"`
}

// SerialDiff defines model for SerialDiff.
type SerialDiff struct {
	BasePipelineMetricAggregation
	Hide     *bool  `json:"hide,omitempty"`
	Id       string `json:"id"`
	Settings *struct {
		Lag *string `json:"lag,omitempty"`
	} `json:"settings,omitempty"`
	Type MetricAggregationType `json:"type"`
}

// Sum defines model for Sum.
type Sum struct {
	MetricAggregationWithField
	MetricAggregationWithInlineScript
	Hide     *bool  `json:"hide,omitempty"`
	Id       string `json:"id"`
	Settings *struct {
		Missing *string `json:"missing,omitempty"`
		Script  *any    `json:"script,omitempty"`
	} `json:"settings,omitempty"`
	Type MetricAggregationType `json:"type"`
}

// Terms defines model for Terms.
type Terms struct {
	BucketAggregationWithField
	Id       string                `json:"id"`
	Settings *any                  `json:"settings,omitempty"`
	Type     BucketAggregationType `json:"type"`
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
	BaseMetricAggregation
	Hide     *bool  `json:"hide,omitempty"`
	Id       string `json:"id"`
	Settings *struct {
		Metrics []string `json:"metrics,omitempty"`
		Order   *string  `json:"order,omitempty"`
		OrderBy *string  `json:"orderBy,omitempty"`
	} `json:"settings,omitempty"`
	Type MetricAggregationType `json:"type"`
}

// UniqueCount defines model for UniqueCount.
type UniqueCount struct {
	MetricAggregationWithField
	Hide     *bool  `json:"hide,omitempty"`
	Id       string `json:"id"`
	Settings *struct {
		Missing            *string `json:"missing,omitempty"`
		PrecisionThreshold *string `json:"precision_threshold,omitempty"`
	} `json:"settings,omitempty"`
	Type MetricAggregationType `json:"type"`
}
