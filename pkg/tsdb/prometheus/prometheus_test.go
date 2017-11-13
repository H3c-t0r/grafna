package prometheus

import (
	"testing"
	"time"

	"github.com/grafana/grafana/pkg/tsdb"

	"github.com/grafana/grafana/pkg/components/simplejson"
	p "github.com/prometheus/common/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPrometheus(t *testing.T) {
	Convey("Prometheus", t, func() {
		Convey("converting metric name", func() {
			metric := map[p.LabelName]p.LabelValue{
				p.LabelName("app"):    p.LabelValue("backend"),
				p.LabelName("device"): p.LabelValue("mobile"),
			}

			query := &PrometheusQuery{
				LegendFormat: "legend {{app}} {{ device }} {{broken}}",
			}

			So(formatLegend(metric, query), ShouldEqual, "legend backend mobile {{broken}}")
		})

		Convey("build full serie name", func() {
			metric := map[p.LabelName]p.LabelValue{
				p.LabelName(p.MetricNameLabel): p.LabelValue("http_request_total"),
				p.LabelName("app"):             p.LabelValue("backend"),
				p.LabelName("device"):          p.LabelValue("mobile"),
			}

			query := &PrometheusQuery{
				LegendFormat: "",
			}

			So(formatLegend(metric, query), ShouldEqual, `http_request_total{app="backend", device="mobile"}`)
		})

		Convey("parsing query model with step", func() {
			json := `{
				"expr": "go_goroutines",
				"format": "time_series",
				"refId": "A"
			}`
			jsonModel, _ := simplejson.NewJson([]byte(json))
			queryContext := &tsdb.TsdbQuery{}
			queryModels := []*tsdb.Query{
				{Model: jsonModel},
			}

			Convey("with 48h time range", func() {
				queryContext.TimeRange = tsdb.NewTimeRange("12h", "now")

				model, err := parseQuery(queryModels, queryContext)

				So(err, ShouldBeNil)
				So(model.Step, ShouldEqual, time.Second*30)
			})
		})

		Convey("parsing query model without step parameter", func() {
			json := `{
				"expr": "go_goroutines",
				"format": "time_series",
				"intervalFactor": 1,
				"refId": "A"
			}`
			jsonModel, _ := simplejson.NewJson([]byte(json))
			queryContext := &tsdb.TsdbQuery{}
			queryModels := []*tsdb.Query{
				{Model: jsonModel},
			}

			Convey("with 48h time range", func() {
				queryContext.TimeRange = tsdb.NewTimeRange("48h", "now")

				model, err := parseQuery(queryModels, queryContext)

				So(err, ShouldBeNil)
				So(model.Step, ShouldEqual, time.Minute*2)
			})

			Convey("with 1h time range", func() {
				queryContext.TimeRange = tsdb.NewTimeRange("1h", "now")

				model, err := parseQuery(queryModels, queryContext)

				So(err, ShouldBeNil)
				So(model.Step, ShouldEqual, time.Second*2)
			})
		})

		Convey("parsing query model with intervalFactor", func() {
			Convey("high intervalFactor", func() {
				json := `{
					"expr": "go_goroutines",
					"format": "time_series",
					"intervalFactor": 10,
					"refId": "A"
				}`
				jsonModel, _ := simplejson.NewJson([]byte(json))
				queryContext := &tsdb.TsdbQuery{}
				queryModels := []*tsdb.Query{
					{Model: jsonModel},
				}

				Convey("with 48h time range", func() {
					queryContext.TimeRange = tsdb.NewTimeRange("48h", "now")

					model, err := parseQuery(queryModels, queryContext)

					So(err, ShouldBeNil)
					So(model.Step, ShouldEqual, time.Minute*20)
				})
			})

			Convey("low intervalFactor", func() {
				json := `{
					"expr": "go_goroutines",
					"format": "time_series",
					"intervalFactor": 1,
					"refId": "A"
				}`
				jsonModel, _ := simplejson.NewJson([]byte(json))
				queryContext := &tsdb.TsdbQuery{}
				queryModels := []*tsdb.Query{
					{Model: jsonModel},
				}

				Convey("with 48h time range", func() {
					queryContext.TimeRange = tsdb.NewTimeRange("48h", "now")

					model, err := parseQuery(queryModels, queryContext)

					So(err, ShouldBeNil)
					So(model.Step, ShouldEqual, time.Minute*2)
				})
			})
		})

	})
}
