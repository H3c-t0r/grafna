package elasticsearch

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/grafana/grafana/pkg/components/simplejson"
	es "github.com/grafana/grafana/pkg/tsdb/elasticsearch/client"

	"github.com/grafana/grafana/pkg/tsdb"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	pplTSFormat   = "2006-01-02 15:04:05.000000"
	pplDateFormat = "2006-01-02"
)

func TestPPLResponseParser(t *testing.T) {
	Convey("PPL response parser test", t, func() {
		Convey("Simple time series query", func() {
			Convey("Time field as first field", func() {
				targets := map[string]string{
					"A": `{
						"timeField": "@timestamp"
					}`,
				}
				response := `{
					"schema": [
						{ "name": "timeName", "type": "timestamp" },
						{ "name": "testMetric", "type": "integer" }
					],
					"datarows": [
						["%s", 10],
						["%s", 15]
					],
					"total": 2,
					"size": 2
				}`
				response = fmt.Sprintf(response, formatUnixMs(100, pplTSFormat), formatUnixMs(200, pplTSFormat))
				rp, err := newPPLResponseParserForTest(targets, response)
				So(err, ShouldBeNil)
				queryRes, err := rp.parseTimeSeries()
				So(err, ShouldBeNil)
				So(queryRes, ShouldNotBeNil)
				So(queryRes.Series, ShouldHaveLength, 1)
				series := queryRes.Series[0]
				So(series.Name, ShouldEqual, "testMetric")
				So(series.Points, ShouldHaveLength, 2)
				So(series.Points[0][0].Float64, ShouldEqual, 10)
				So(series.Points[0][1].Float64, ShouldEqual, 100)
				So(series.Points[1][0].Float64, ShouldEqual, 15)
				So(series.Points[1][1].Float64, ShouldEqual, 200)
			})

			Convey("Time field as second field", func() {
				targets := map[string]string{
					"A": `{
						"timeField": "@timestamp"
					}`,
				}
				response := `{
					"schema": [
						{ "name": "testMetric", "type": "integer" },
						{ "name": "timeName", "type": "timestamp" }
					],
					"datarows": [
						[20, "%s"],
						[25, "%s"]
					],
					"total": 2,
					"size": 2
				}`
				response = fmt.Sprintf(response, formatUnixMs(100, pplTSFormat), formatUnixMs(200, pplTSFormat))
				rp, err := newPPLResponseParserForTest(targets, response)
				So(err, ShouldBeNil)
				queryRes, err := rp.parseTimeSeries()
				So(err, ShouldBeNil)
				So(queryRes, ShouldNotBeNil)
				So(queryRes.Series, ShouldHaveLength, 1)
				series := queryRes.Series[0]
				So(series.Name, ShouldEqual, "testMetric")
				So(series.Points, ShouldHaveLength, 2)
				So(series.Points[0][0].Float64, ShouldEqual, 20)
				So(series.Points[0][1].Float64, ShouldEqual, 100)
				So(series.Points[1][0].Float64, ShouldEqual, 25)
				So(series.Points[1][1].Float64, ShouldEqual, 200)
			})
		})

		Convey("Set series name to be value field name", func() {
			targets := map[string]string{
				"A": `{
					"timeField": "@timestamp"
				}`,
			}
			response := `{
				"schema": [
					{ "name": "valueField", "type": "integer" },
					{ "name": "timeName", "type": "timestamp" }
				],
				"datarows": [
					[20, "%s"]
				],
				"total": 1,
				"size": 1
			}`
			response = fmt.Sprintf(response, formatUnixMs(100, pplTSFormat))
			rp, err := newPPLResponseParserForTest(targets, response)
			So(err, ShouldBeNil)
			queryRes, err := rp.parseTimeSeries()
			So(err, ShouldBeNil)
			So(queryRes, ShouldNotBeNil)
			So(queryRes.Series, ShouldHaveLength, 1)
			series := queryRes.Series[0]
			So(series.Name, ShouldEqual, "valueField")
		})

		Convey("Different date formats", func() {
			targets := map[string]string{
				"A": `{
					"timeField": "@timestamp"
				}`,
			}
			response := `{
				"schema": [
					{ "name": "timeName", "type": "%s" },
					{ "name": "testMetric", "type": "integer" }
				],
				"datarows": [
					["%s", 10]
				],
				"total": 1,
				"size": 1
			}`

			Convey("Timestamp time field type", func() {
				formattedResponse := fmt.Sprintf(response, "timestamp", formatUnixMs(100, pplTSFormat))
				rp, err := newPPLResponseParserForTest(targets, formattedResponse)
				So(err, ShouldBeNil)
				queryRes, err := rp.parseTimeSeries()
				So(err, ShouldBeNil)
				So(queryRes, ShouldNotBeNil)
				So(queryRes.Series, ShouldHaveLength, 1)
				series := queryRes.Series[0]
				So(series.Name, ShouldEqual, "testMetric")
				So(series.Points, ShouldHaveLength, 1)
				So(series.Points[0][0].Float64, ShouldEqual, 10)
				So(series.Points[0][1].Float64, ShouldEqual, 100)
			})

			Convey("Datetime time field type", func() {
				formattedResponse := fmt.Sprintf(response, "datetime", formatUnixMs(100, pplTSFormat))
				rp, err := newPPLResponseParserForTest(targets, formattedResponse)
				So(err, ShouldBeNil)
				queryRes, err := rp.parseTimeSeries()
				So(err, ShouldBeNil)
				So(queryRes, ShouldNotBeNil)
				So(queryRes.Series, ShouldHaveLength, 1)
				series := queryRes.Series[0]
				So(series.Name, ShouldEqual, "testMetric")
				So(series.Points, ShouldHaveLength, 1)
				So(series.Points[0][0].Float64, ShouldEqual, 10)
				So(series.Points[0][1].Float64, ShouldEqual, 100)
			})

			Convey("Date time field type", func() {
				formattedResponse := fmt.Sprintf(response, "date", formatUnixMs(0, pplDateFormat))
				rp, err := newPPLResponseParserForTest(targets, formattedResponse)
				So(err, ShouldBeNil)
				queryRes, err := rp.parseTimeSeries()
				So(err, ShouldBeNil)
				So(queryRes, ShouldNotBeNil)
				So(queryRes.Series, ShouldHaveLength, 1)
				series := queryRes.Series[0]
				So(series.Name, ShouldEqual, "testMetric")
				So(series.Points, ShouldHaveLength, 1)
				So(series.Points[0][0].Float64, ShouldEqual, 10)
				So(series.Points[0][1].Float64, ShouldEqual, 0)
			})
		})

		Convey("Handle invalid schema for time series", func() {
			Convey("More than two fields", func() {
				targets := map[string]string{
					"A": `{
						"timeField": "@timestamp"
					}`,
				}
				response := `{
				"schema": [
					{ "name": "testMetric", "type": "integer" },
					{ "name": "extraMetric", "type": "integer" },
					{ "name": "timeName", "type": "timestamp" }
				],
				"datarows": [
					[20, 20, "%s"]
				],
				"total": 1,
				"size": 1
			}`
				response = fmt.Sprintf(response, formatUnixMs(100, pplTSFormat))
				rp, err := newPPLResponseParserForTest(targets, response)
				So(err, ShouldBeNil)
				_, err = rp.parseTimeSeries()
				So(err, ShouldNotBeNil)
			})

			Convey("Less than two fields", func() {
				targets := map[string]string{
					"A": `{
						"timeField": "@timestamp"
					}`,
				}
				response := `{
					"schema": [
						{ "name": "timeName", "type": "timestamp" }
					],
					"datarows": [
						["%s"]
					],
					"total": 1,
					"size": 1
				}`
				response = fmt.Sprintf(response, formatUnixMs(100, pplTSFormat))
				rp, err := newPPLResponseParserForTest(targets, response)
				So(err, ShouldBeNil)
				_, err = rp.parseTimeSeries()
				So(err, ShouldNotBeNil)
			})

			Convey("No valid time field type", func() {
				targets := map[string]string{
					"A": `{
						"timeField": "@timestamp"
					}`,
				}
				response := `{
					"schema": [
						{ "name": "timeName", "type": "string" },
						{ "name": "testMetric", "type": "integer" }
					],
					"datarows": [
						["%s", 10]
					],
					"total": 1,
					"size": 1
				}`
				response = fmt.Sprintf(response, formatUnixMs(100, pplTSFormat))
				rp, err := newPPLResponseParserForTest(targets, response)
				So(err, ShouldBeNil)
				_, err = rp.parseTimeSeries()
				So(err, ShouldNotBeNil)
			})

			Convey("Valid time field type with invalid value type", func() {
				targets := map[string]string{
					"A": `{
						"timeField": "@timestamp"
					}`,
				}
				response := `{
					"schema": [
						{ "name": "timeName", "type": "timestamp" },
						{ "name": "testMetric", "type": "string" }
					],
					"datarows": [
						["%s", "10"]
					],
					"total": 1,
					"size": 1
				}`
				response = fmt.Sprintf(response, formatUnixMs(100, pplTSFormat))
				rp, err := newPPLResponseParserForTest(targets, response)
				So(err, ShouldBeNil)
				_, err = rp.parseTimeSeries()
				So(err, ShouldNotBeNil)
			})

			Convey("Valid schema invalid time field type", func() {
				targets := map[string]string{
					"A": `{
						"timeField": "@timestamp"
					}`,
				}
				response := `{
					"schema": [
						{ "name": "timeName", "type": "timestamp" },
						{ "name": "testMetric", "type": "string" }
					],
					"datarows": [
						[10, "10"]
					],
					"total": 1,
					"size": 1
				}`
				rp, err := newPPLResponseParserForTest(targets, response)
				So(err, ShouldBeNil)
				_, err = rp.parseTimeSeries()
				So(err, ShouldNotBeNil)
			})

			Convey("Valid schema invalid time field value", func() {
				targets := map[string]string{
					"A": `{
						"timeField": "@timestamp"
					}`,
				}
				response := `{
					"schema": [
						{ "name": "timeName", "type": "timestamp" },
						{ "name": "testMetric", "type": "string" }
					],
					"datarows": [
						["foo", "10"]
					],
					"total": 1,
					"size": 1
				}`
				rp, err := newPPLResponseParserForTest(targets, response)
				So(err, ShouldBeNil)
				_, err = rp.parseTimeSeries()
				So(err, ShouldNotBeNil)
			})
		})

		Convey("Parses error response", func() {
			targets := map[string]string{
				"A": `{
					"timeField": "@timestamp"
				}`,
			}
			response := `{
				"error": {
					"reason": "Error occurred in Elasticsearch engine: no such index [unknown]",
					"details": "org.elasticsearch.index.IndexNotFoundException: no such index [unknown].",
					"type": "IndexNotFoundException"
				}
			}`
			rp, err := newPPLResponseParserForTest(targets, response)
			So(err, ShouldBeNil)
			queryRes, err := rp.parseTimeSeries()
			So(queryRes, ShouldNotBeNil)
			So(queryRes.ErrorString, ShouldEqual, "Error occurred in Elasticsearch engine: no such index [unknown]")
			So(queryRes.Series, ShouldHaveLength, 0)
			So(err, ShouldBeNil)
		})

		Convey("Query result meta field is set", func() {
			Convey("On successful response", func() {
				targets := map[string]string{
					"A": `{
						"timeField": "@timestamp"
					}`,
				}
				response := `{
					"schema": [
						{ "name": "valueField", "type": "integer" },
						{ "name": "timeName", "type": "timestamp" }
					],
					"datarows": [
						[20, "%s"]
					],
					"total": 1,
					"size": 1
				}`
				response = fmt.Sprintf(response, formatUnixMs(100, pplTSFormat))
				rp, err := newPPLResponseParserForTest(targets, response)
				So(err, ShouldBeNil)
				queryRes, err := rp.parseTimeSeries()
				So(err, ShouldBeNil)
				So(queryRes, ShouldNotBeNil)
				So(queryRes.Meta, ShouldNotBeNil)
			})
			Convey("On error response", func() {
				targets := map[string]string{
					"A": `{
						"timeField": "@timestamp"
					}`,
				}
				response := `{
					"error": {
						"reason": "Error occurred in Elasticsearch engine: no such index [unknown]",
						"details": "org.elasticsearch.index.IndexNotFoundException: no such index [unknown].",
						"type": "IndexNotFoundException"
					}
				}`
				rp, err := newPPLResponseParserForTest(targets, response)
				So(err, ShouldBeNil)
				queryRes, err := rp.parseTimeSeries()
				So(err, ShouldBeNil)
				So(queryRes, ShouldNotBeNil)
				So(queryRes.Meta, ShouldNotBeNil)
			})
		})
	})
}

func newPPLResponseParserForTest(tsdbQueries map[string]string, responseBody string) (*pplResponseParser, error) {
	from := time.Date(2018, 5, 15, 17, 50, 0, 0, time.UTC)
	to := time.Date(2018, 5, 15, 17, 55, 0, 0, time.UTC)
	fromStr := fmt.Sprintf("%d", from.UnixNano()/int64(time.Millisecond))
	toStr := fmt.Sprintf("%d", to.UnixNano()/int64(time.Millisecond))
	tsdbQuery := &tsdb.TsdbQuery{
		Queries:   []*tsdb.Query{},
		TimeRange: tsdb.NewTimeRange(fromStr, toStr),
	}

	for refID, tsdbQueryBody := range tsdbQueries {
		tsdbQueryJSON, err := simplejson.NewJson([]byte(tsdbQueryBody))
		if err != nil {
			return nil, err
		}

		tsdbQuery.Queries = append(tsdbQuery.Queries, &tsdb.Query{
			Model: tsdbQueryJSON,
			RefId: refID,
		})
	}

	var response es.PPLResponse
	err := json.Unmarshal([]byte(responseBody), &response)
	if err != nil {
		return nil, err
	}

	response.DebugInfo = &es.PPLDebugInfo{
		Response: &es.PPLResponseInfo{
			Status: 200,
		},
	}

	tsQueryParser := newTimeSeriesQueryParser()
	queries, err := tsQueryParser.parse(tsdbQuery)
	if err != nil {
		return nil, err
	}

	return newPPLResponseParser(&response, queries[0]), nil
}

func formatUnixMs(ms int64, format string) string {
	return time.Unix(0, ms*int64(time.Millisecond)).UTC().Format(format)
}
