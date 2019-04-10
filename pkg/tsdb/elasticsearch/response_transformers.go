package elasticsearch

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/grafana/grafana/pkg/components/null"
	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/tsdb"
	"github.com/grafana/grafana/pkg/tsdb/elasticsearch/client"
)

const (
	// Metric types
	countType         = "count"
	percentilesType   = "percentiles"
	extendedStatsType = "extended_stats"
	// Bucket types
	dateHistType    = "date_histogram"
	histogramType   = "histogram"
	filtersType     = "filters"
	termsType       = "terms"
	geohashGridType = "geohash_grid"
)

type responseTransformer interface {
	transform() (*tsdb.Response, error)
}

type timeSeriesQueryResponseTransformer struct {
	Responses []*es.SearchResponse
	Targets   []*Query
}

var newTimeSeriesQueryResponseTransformer = func(responses []*es.SearchResponse, targets []*Query) responseTransformer {
	return &timeSeriesQueryResponseTransformer{
		Responses: responses,
		Targets:   targets,
	}
}

func (rp *timeSeriesQueryResponseTransformer) transform() (*tsdb.Response, error) {
	result := &tsdb.Response{}
	result.Results = make(map[string]*tsdb.QueryResult)

	if rp.Responses == nil {
		return result, nil
	}

	for i, res := range rp.Responses {
		target := rp.Targets[i]

		if res.Error != nil {
			result.Results[target.RefID] = getErrorFromElasticResponse(res.Error)
			continue
		}

		queryRes := tsdb.NewQueryResult()
		props := make(map[string]string)
		table := tsdb.Table{
			Columns: make([]tsdb.TableColumn, 0),
			Rows:    make([]tsdb.RowValues, 0),
		}

		if res.Hits != nil && len(res.Hits.Hits) > 0 {
			err := rp.processHits(res.Hits, &table)
			if err != nil {
				return nil, err
			}
			queryRes.Tables = append(queryRes.Tables, &table)
		}

		if res.Aggregations != nil && len(res.Aggregations) > 0 {
			err := rp.processBuckets(res.Aggregations, target, &queryRes.Series, &table, props, 0)
			if err != nil {
				return nil, err
			}
			rp.nameSeries(&queryRes.Series, target)
			rp.trimDatapoints(&queryRes.Series, target)

			if len(table.Rows) > 0 {
				queryRes.Tables = append(queryRes.Tables, &table)
			}
		}

		result.Results[target.RefID] = queryRes
	}
	return result, nil
}

func (rp *timeSeriesQueryResponseTransformer) processHits(hits *es.SearchResponseHits, table *tsdb.Table) error {
	table.Columns = append(table.Columns, tsdb.TableColumn{Text: "hits"})
	table.Columns = append(table.Columns, tsdb.TableColumn{Text: "total"})

	docs := []interface{}{}

	for _, v := range hits.Hits {
		hit := simplejson.NewFromAny(v)
		doc := map[string]interface{}{}

		doc["_id"] = hit.Get("_id").Interface()
		doc["_type"] = hit.Get("_type").Interface()
		doc["_index"] = hit.Get("_index").Interface()

		if sourceProp, ok := hit.CheckGet("_source"); ok {
			propNames := []string{}
			props := sourceProp.MustMap()
			for prop := range props {
				propNames = append(propNames, prop)
			}
			sort.Strings(propNames)
			for _, propName := range propNames {
				doc[propName] = props[propName]
			}
		}

		if fieldsProp, ok := hit.CheckGet("fields"); ok {
			propNames := []string{}
			props := fieldsProp.MustMap()
			for prop := range props {
				propNames = append(propNames, prop)
			}
			sort.Strings(propNames)
			for _, propName := range propNames {
				doc[propName] = props[propName]
			}
		}

		docs = append(docs, doc)
	}

	table.Rows = append(table.Rows, tsdb.RowValues{docs, hits.Total})

	return nil
}

func (rp *timeSeriesQueryResponseTransformer) processBuckets(aggs map[string]interface{}, target *Query, series *tsdb.TimeSeriesSlice, table *tsdb.Table, props map[string]string, depth int) error {
	var err error
	maxDepth := len(target.BucketAggs) - 1

	aggIDs := make([]string, 0)
	for k := range aggs {
		aggIDs = append(aggIDs, k)
	}
	sort.Strings(aggIDs)
	for _, aggID := range aggIDs {
		v := aggs[aggID]
		aggDef, _ := findAgg(target, aggID)
		esAgg := simplejson.NewFromAny(v)
		if aggDef == nil {
			continue
		}

		if depth == maxDepth {
			if aggDef.Type == dateHistType {
				err = rp.processMetrics(esAgg, target, series, props)
			} else {
				err = rp.processAggregationDocs(esAgg, aggDef, target, table, props)
			}
			if err != nil {
				return err
			}
		} else {
			for _, b := range esAgg.Get("buckets").MustArray() {
				bucket := simplejson.NewFromAny(b)
				newProps := make(map[string]string)

				for k, v := range props {
					newProps[k] = v
				}

				if key, err := bucket.Get("key").String(); err == nil {
					newProps[aggDef.Field] = key
				} else if key, err := bucket.Get("key").Int64(); err == nil {
					newProps[aggDef.Field] = strconv.FormatInt(key, 10)
				}

				if key, err := bucket.Get("key_as_string").String(); err == nil {
					newProps[aggDef.Field] = key
				}
				err = rp.processBuckets(bucket.MustMap(), target, series, table, newProps, depth+1)
				if err != nil {
					return err
				}
			}

			buckets := esAgg.Get("buckets").MustMap()
			bucketKeys := make([]string, 0)
			for k := range buckets {
				bucketKeys = append(bucketKeys, k)
			}
			sort.Strings(bucketKeys)

			for _, bucketKey := range bucketKeys {
				bucket := simplejson.NewFromAny(buckets[bucketKey])
				newProps := make(map[string]string)

				for k, v := range props {
					newProps[k] = v
				}

				newProps["filter"] = bucketKey

				err = rp.processBuckets(bucket.MustMap(), target, series, table, newProps, depth+1)
				if err != nil {
					return err
				}
			}
		}

	}
	return nil

}

func (rp *timeSeriesQueryResponseTransformer) processMetrics(esAgg *simplejson.Json, target *Query, series *tsdb.TimeSeriesSlice, props map[string]string) error {
	for _, metric := range target.Metrics {
		if metric.Hide {
			continue
		}

		switch metric.Type {
		case countType:
			newSeries := tsdb.TimeSeries{
				Tags: make(map[string]string),
			}

			for _, v := range esAgg.Get("buckets").MustArray() {
				bucket := simplejson.NewFromAny(v)
				value := castToNullFloat(bucket.Get("doc_count"))
				key := castToNullFloat(bucket.Get("key"))
				newSeries.Points = append(newSeries.Points, tsdb.TimePoint{value, key})
			}

			for k, v := range props {
				newSeries.Tags[k] = v
			}
			newSeries.Tags["metric"] = countType
			*series = append(*series, &newSeries)

		case percentilesType:
			buckets := esAgg.Get("buckets").MustArray()
			if len(buckets) == 0 {
				break
			}

			firstBucket := simplejson.NewFromAny(buckets[0])
			percentiles := firstBucket.GetPath(metric.ID, "values").MustMap()

			percentileKeys := make([]string, 0)
			for k := range percentiles {
				percentileKeys = append(percentileKeys, k)
			}
			sort.Strings(percentileKeys)
			for _, percentileName := range percentileKeys {
				newSeries := tsdb.TimeSeries{
					Tags: make(map[string]string),
				}
				for k, v := range props {
					newSeries.Tags[k] = v
				}
				newSeries.Tags["metric"] = "p" + percentileName
				newSeries.Tags["field"] = metric.Field
				for _, v := range buckets {
					bucket := simplejson.NewFromAny(v)
					value := castToNullFloat(bucket.GetPath(metric.ID, "values", percentileName))
					key := castToNullFloat(bucket.Get("key"))
					newSeries.Points = append(newSeries.Points, tsdb.TimePoint{value, key})
				}
				*series = append(*series, &newSeries)
			}
		case extendedStatsType:
			buckets := esAgg.Get("buckets").MustArray()

			metaKeys := make([]string, 0)
			meta := metric.Meta.MustMap()
			for k := range meta {
				metaKeys = append(metaKeys, k)
			}
			sort.Strings(metaKeys)
			for _, statName := range metaKeys {
				v := meta[statName]
				if enabled, ok := v.(bool); !ok || !enabled {
					continue
				}

				newSeries := tsdb.TimeSeries{
					Tags: make(map[string]string),
				}
				for k, v := range props {
					newSeries.Tags[k] = v
				}
				newSeries.Tags["metric"] = statName
				newSeries.Tags["field"] = metric.Field

				for _, v := range buckets {
					bucket := simplejson.NewFromAny(v)
					key := castToNullFloat(bucket.Get("key"))
					var value null.Float
					if statName == "std_deviation_bounds_upper" {
						value = castToNullFloat(bucket.GetPath(metric.ID, "std_deviation_bounds", "upper"))
					} else if statName == "std_deviation_bounds_lower" {
						value = castToNullFloat(bucket.GetPath(metric.ID, "std_deviation_bounds", "lower"))
					} else {
						value = castToNullFloat(bucket.GetPath(metric.ID, statName))
					}
					newSeries.Points = append(newSeries.Points, tsdb.TimePoint{value, key})
				}
				*series = append(*series, &newSeries)
			}
		default:
			newSeries := tsdb.TimeSeries{
				Tags: make(map[string]string),
			}
			for k, v := range props {
				newSeries.Tags[k] = v
			}

			newSeries.Tags["metric"] = metric.Type
			newSeries.Tags["field"] = metric.Field
			newSeries.Tags["metricId"] = metric.ID
			for _, v := range esAgg.Get("buckets").MustArray() {
				bucket := simplejson.NewFromAny(v)
				key := castToNullFloat(bucket.Get("key"))
				valueObj, err := bucket.Get(metric.ID).Map()
				if err != nil {
					continue
				}
				var value null.Float
				if _, ok := valueObj["normalized_value"]; ok {
					value = castToNullFloat(bucket.GetPath(metric.ID, "normalized_value"))
				} else {
					value = castToNullFloat(bucket.GetPath(metric.ID, "value"))
				}
				newSeries.Points = append(newSeries.Points, tsdb.TimePoint{value, key})
			}
			*series = append(*series, &newSeries)
		}
	}
	return nil
}

func (rp *timeSeriesQueryResponseTransformer) processAggregationDocs(esAgg *simplejson.Json, aggDef *BucketAgg, target *Query, table *tsdb.Table, props map[string]string) error {
	propKeys := make([]string, 0)
	for k := range props {
		propKeys = append(propKeys, k)
	}
	sort.Strings(propKeys)

	if len(table.Columns) == 0 {
		for _, propKey := range propKeys {
			table.Columns = append(table.Columns, tsdb.TableColumn{Text: propKey})
		}
		table.Columns = append(table.Columns, tsdb.TableColumn{Text: aggDef.Field})
	}

	addMetricValue := func(values *tsdb.RowValues, metricName string, value null.Float) {
		found := false
		for _, c := range table.Columns {
			if c.Text == metricName {
				found = true
				break
			}
		}
		if !found {
			table.Columns = append(table.Columns, tsdb.TableColumn{Text: metricName})
		}
		*values = append(*values, value)
	}

	for _, v := range esAgg.Get("buckets").MustArray() {
		bucket := simplejson.NewFromAny(v)
		values := make(tsdb.RowValues, 0)

		for _, propKey := range propKeys {
			values = append(values, props[propKey])
		}

		if key, err := bucket.Get("key").String(); err == nil {
			values = append(values, key)
		} else {
			values = append(values, castToNullFloat(bucket.Get("key")))
		}

		for _, metric := range target.Metrics {
			switch metric.Type {
			case countType:
				addMetricValue(&values, rp.getMetricName(metric.Type), castToNullFloat(bucket.Get("doc_count")))
			case extendedStatsType:
				metaKeys := make([]string, 0)
				meta := metric.Meta.MustMap()
				for k := range meta {
					metaKeys = append(metaKeys, k)
				}
				sort.Strings(metaKeys)
				for _, statName := range metaKeys {
					v := meta[statName]
					if enabled, ok := v.(bool); !ok || !enabled {
						continue
					}

					var value null.Float
					if statName == "std_deviation_bounds_upper" {
						value = castToNullFloat(bucket.GetPath(metric.ID, "std_deviation_bounds", "upper"))
					} else if statName == "std_deviation_bounds_lower" {
						value = castToNullFloat(bucket.GetPath(metric.ID, "std_deviation_bounds", "lower"))
					} else {
						value = castToNullFloat(bucket.GetPath(metric.ID, statName))
					}

					addMetricValue(&values, rp.getMetricName(metric.Type), value)
					break
				}
			default:
				metricName := rp.getMetricName(metric.Type)
				otherMetrics := make([]*MetricAgg, 0)

				for _, m := range target.Metrics {
					if m.Type == metric.Type {
						otherMetrics = append(otherMetrics, m)
					}
				}

				if len(otherMetrics) > 1 {
					metricName += " " + metric.Field
				}

				addMetricValue(&values, metricName, castToNullFloat(bucket.GetPath(metric.ID, "value")))
			}
		}

		table.Rows = append(table.Rows, values)
	}

	return nil
}

func (rp *timeSeriesQueryResponseTransformer) trimDatapoints(series *tsdb.TimeSeriesSlice, target *Query) {
	var histogram *BucketAgg
	for _, bucketAgg := range target.BucketAggs {
		if bucketAgg.Type == dateHistType {
			histogram = bucketAgg
			break
		}
	}

	if histogram == nil {
		return
	}

	trimEdges, err := histogram.Settings.Get("trimEdges").Int()
	if err != nil {
		return
	}

	for _, s := range *series {
		if len(s.Points) > trimEdges*2 {
			s.Points = s.Points[trimEdges : len(s.Points)-trimEdges]
		}
	}
}

func (rp *timeSeriesQueryResponseTransformer) nameSeries(seriesList *tsdb.TimeSeriesSlice, target *Query) {
	set := make(map[string]string)
	for _, v := range *seriesList {
		if metricType, exists := v.Tags["metric"]; exists {
			if _, ok := set[metricType]; !ok {
				set[metricType] = ""
			}
		}
	}
	metricTypeCount := len(set)
	for _, series := range *seriesList {
		series.Name = rp.getSeriesName(series, target, metricTypeCount)
	}

}

var aliasPatternRegex = regexp.MustCompile(`\{\{([\s\S]+?)\}\}`)

func (rp *timeSeriesQueryResponseTransformer) getSeriesName(series *tsdb.TimeSeries, target *Query, metricTypeCount int) string {
	metricType := series.Tags["metric"]
	metricName := rp.getMetricName(metricType)
	delete(series.Tags, "metric")

	field := ""
	if v, ok := series.Tags["field"]; ok {
		field = v
		delete(series.Tags, "field")
	}

	if target.Alias != "" {
		seriesName := target.Alias

		subMatches := aliasPatternRegex.FindAllStringSubmatch(target.Alias, -1)
		for _, subMatch := range subMatches {
			group := subMatch[0]

			if len(subMatch) > 1 {
				group = subMatch[1]
			}

			if strings.Index(group, "term ") == 0 {
				seriesName = strings.Replace(seriesName, subMatch[0], series.Tags[group[5:]], 1)
			}
			if v, ok := series.Tags[group]; ok {
				seriesName = strings.Replace(seriesName, subMatch[0], v, 1)
			}
			if group == "metric" {
				seriesName = strings.Replace(seriesName, subMatch[0], metricName, 1)
			}
			if group == "field" {
				seriesName = strings.Replace(seriesName, subMatch[0], field, 1)
			}
		}

		return seriesName
	}
	// todo, if field and pipelineAgg
	if field != "" && isPipelineAgg(metricType) {
		if isPipelineAggWithMultipleBucketPaths(metricType) {
			metricID := ""
			if v, ok := series.Tags["metricId"]; ok {
				metricID = v
			}

			for _, metric := range target.Metrics {
				if metric.ID == metricID {
					metricName = metric.Settings.Get("script").MustString()
					for name, pipelineAgg := range metric.PipelineVariables {
						for _, m := range target.Metrics {
							if m.ID == pipelineAgg {
								metricName = strings.Replace(metricName, "params."+name, describeMetric(m.Type, m.Field), -1)
							}
						}
					}
				}
			}
		} else {
			found := false
			for _, metric := range target.Metrics {
				if metric.ID == field {
					metricName += " " + describeMetric(metric.Type, field)
					found = true
				}
			}
			if !found {
				metricName = "Unset"
			}
		}
	} else if field != "" {
		metricName += " " + field
	}

	delete(series.Tags, "metricId")

	if len(series.Tags) == 0 {
		return metricName
	}

	name := ""
	for _, v := range series.Tags {
		name += v + " "
	}

	if metricTypeCount == 1 {
		return strings.TrimSpace(name)
	}

	return strings.TrimSpace(name) + " " + metricName

}

func (rp *timeSeriesQueryResponseTransformer) getMetricName(metric string) string {
	if text, ok := metricAggType[metric]; ok {
		return text
	}

	if text, ok := extendedStats[metric]; ok {
		return text
	}

	return metric
}

var fieldTypeMap = map[string]string{
	"float":        "number",
	"double":       "number",
	"integer":      "number",
	"long":         "number",
	"date":         "date",
	"string":       "string",
	"text":         "string",
	"scaled_float": "number",
	"nested":       "nested",
}

type fieldsQueryResponseTransformer struct {
	Response        *es.IndexMappingResponse
	FieldTypeFilter string
	RefID           string
}

var newFieldsQueryResponseTransformer = func(response *es.IndexMappingResponse, fieldTypeFilter, refID string) responseTransformer {
	return &fieldsQueryResponseTransformer{
		Response:        response,
		FieldTypeFilter: fieldTypeFilter,
		RefID:           refID,
	}
}

func (t *fieldsQueryResponseTransformer) transform() (*tsdb.Response, error) {
	res := t.Response

	if res.Error != nil {
		return &tsdb.Response{
			Results: map[string]*tsdb.QueryResult{
				t.RefID: getErrorFromElasticResponse(res.Error),
			},
		}, nil
	}

	fields := map[string]string{}

	walkFunc := func(node simplejson.JsonNode, path []string, err error) error {
		switch nt := node.(type) {
		case simplejson.JsonObjectProperty:
			if strings.HasPrefix(nt.Key, "_") {
				return simplejson.SkipNode
			}
		case simplejson.JsonValue:
			lastElm := path[len(path)-1]
			if lastElm != "type" {
				return nil
			}

			fieldType, ok := nt.Value.(string)
			if !ok {
				return nil
			}

			if t.FieldTypeFilter == "" || t.FieldTypeFilter == fieldType || t.FieldTypeFilter == fieldTypeMap[fieldType] {
				path = path[:len(path)-1]
				fieldNameParts := []string{}
				for _, p := range path {
					if p != "properties" && p != "fields" {
						fieldNameParts = append(fieldNameParts, p)
					}
				}

				mappedFieldType, ok := fieldTypeMap[fieldType]
				if !ok {
					mappedFieldType = fieldType
				}

				fields[strings.Join(fieldNameParts, ".")] = mappedFieldType
			}
		}

		return nil
	}

	for indexName := range res.Mappings {
		index := simplejson.NewFromAny(res.Mappings[indexName])
		typeNames := index.Get("mappings").MustMap()

		for _, v := range typeNames {
			simplejson.Walk(simplejson.NewFromAny(v), walkFunc)
		}
	}

	table := tsdb.Table{
		Columns: make([]tsdb.TableColumn, 0),
		Rows:    make([]tsdb.RowValues, 0),
	}

	table.Columns = append(table.Columns, tsdb.TableColumn{Text: "name"})
	table.Columns = append(table.Columns, tsdb.TableColumn{Text: "type"})

	fieldNames := []string{}
	for fieldName := range fields {
		fieldNames = append(fieldNames, fieldName)
	}
	sort.Strings(fieldNames)

	for _, fieldName := range fieldNames {
		table.Rows = append(table.Rows, tsdb.RowValues{fieldName, fields[fieldName]})
	}

	result := tsdb.Response{
		Results: map[string]*tsdb.QueryResult{
			t.RefID: {
				RefId:  t.RefID,
				Tables: []*tsdb.Table{&table},
			},
		},
	}

	return &result, nil
}

type termsQueryResponseTransformer struct {
	Response   *es.SearchResponse
	TermsAggID string
	RefID      string
}

var newTermsQueryResponseTransformer = func(response *es.SearchResponse, termsAggID, refID string) responseTransformer {
	return &termsQueryResponseTransformer{
		Response:   response,
		TermsAggID: termsAggID,
		RefID:      refID,
	}
}

func (rp *termsQueryResponseTransformer) transform() (*tsdb.Response, error) {
	res := rp.Response

	if res.Error != nil {
		return &tsdb.Response{
			Results: map[string]*tsdb.QueryResult{
				rp.RefID: getErrorFromElasticResponse(res.Error),
			},
		}, nil
	}

	termsAgg, ok := res.Aggregations[rp.TermsAggID]
	if !ok {
		return nil, fmt.Errorf("terms aggregation with id 1 not found in response")
	}

	table := tsdb.Table{
		Columns: make([]tsdb.TableColumn, 0),
		Rows:    make([]tsdb.RowValues, 0),
	}

	table.Columns = append(table.Columns, tsdb.TableColumn{Text: "term"})
	table.Columns = append(table.Columns, tsdb.TableColumn{Text: "doc_count"})

	agg := simplejson.NewFromAny(termsAgg)
	for _, v := range agg.Get("buckets").MustArray() {
		bucket := simplejson.NewFromAny(v)
		term := ""
		docCount := castToNullFloat(bucket.Get("doc_count"))

		if key, err := bucket.Get("key").String(); err == nil {
			term = key
		} else if key, err := bucket.Get("key").Int64(); err == nil {
			term = strconv.FormatInt(key, 10)
		}

		if key, err := bucket.Get("key_as_string").String(); err == nil {
			term = key
		}

		table.Rows = append(table.Rows, tsdb.RowValues{term, docCount})
	}

	result := tsdb.Response{
		Results: map[string]*tsdb.QueryResult{
			rp.RefID: {
				RefId:  rp.RefID,
				Tables: []*tsdb.Table{&table},
			},
		},
	}

	return &result, nil
}

func castToNullFloat(j *simplejson.Json) null.Float {
	f, err := j.Float64()
	if err == nil {
		return null.FloatFrom(f)
	}

	if s, err := j.String(); err == nil {
		if strings.ToLower(s) == "nan" {
			return null.NewFloat(0, false)
		}

		if v, err := strconv.ParseFloat(s, 64); err == nil {
			return null.FloatFromPtr(&v)
		}
	}

	return null.NewFloat(0, false)
}

func findAgg(target *Query, aggID string) (*BucketAgg, error) {
	for _, v := range target.BucketAggs {
		if aggID == v.ID {
			return v, nil
		}
	}
	return nil, errors.New("can't found aggDef, aggID:" + aggID)
}

func getErrorFromElasticResponse(err map[string]interface{}) *tsdb.QueryResult {
	result := tsdb.NewQueryResult()
	json := simplejson.NewFromAny(err)
	reason := json.Get("reason").MustString()
	rootCauseReason := json.Get("root_cause").GetIndex(0).Get("reason").MustString()

	if rootCauseReason != "" {
		result.ErrorString = rootCauseReason
	} else if reason != "" {
		result.ErrorString = reason
	} else {
		result.ErrorString = "Unknown elasticsearch error response"
	}

	return result
}
