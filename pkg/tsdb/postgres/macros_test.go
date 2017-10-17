package postgres

import (
	"testing"

	"github.com/grafana/grafana/pkg/tsdb"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMacroEngine(t *testing.T) {
	Convey("MacroEngine", t, func() {
		engine := &PostgresMacroEngine{}
		timeRange := &tsdb.TimeRange{From: "5m", To: "now"}

		Convey("interpolate __time function", func() {
			sql, err := engine.Interpolate(nil, "select $__time(time_column)")
			So(err, ShouldBeNil)

			So(sql, ShouldEqual, "select time_column AS \"time\"")
		})

		Convey("interpolate __time function wrapped in aggregation", func() {
			sql, err := engine.Interpolate(nil, "select min($__time(time_column))")
			So(err, ShouldBeNil)

			So(sql, ShouldEqual, "select min(time_column AS \"time\")")
		})

		Convey("interpolate __timeFilter function", func() {
			sql, err := engine.Interpolate(timeRange, "WHERE $__timeFilter(time_column)")
			So(err, ShouldBeNil)

			So(sql, ShouldEqual, "WHERE time_column >= to_timestamp(18446744066914186738) AND time_column <= to_timestamp(18446744066914187038)")
		})

		Convey("interpolate __timeFromRFC3339 function", func() {
			sql, err := engine.Interpolate(timeRange, "select '$__timeFromRFC3339()'")
			So(err, ShouldBeNil)

			So(sql, ShouldEqual, "select '0000-12-31T23:55:00Z'")
		})

		Convey("interpolate __timeFrom function", func() {
			sql, err := engine.Interpolate(timeRange, "select $__timeFrom(time_column)")
			So(err, ShouldBeNil)

			So(sql, ShouldEqual, "select to_timestamp(18446744066914186738)")
		})

		Convey("interpolate __timeGroup function", func() {

			sql, err := engine.Interpolate(timeRange, "GROUP BY $__timeGroup(time_column,'5m')")
			So(err, ShouldBeNil)

			So(sql, ShouldEqual, "GROUP BY (extract(epoch from \"time_column\")/extract(epoch from '5m'::interval))::int")
		})

		Convey("interpolate __timeTo function", func() {
			sql, err := engine.Interpolate(timeRange, "select $__timeTo(time_column)")
			So(err, ShouldBeNil)

			So(sql, ShouldEqual, "select to_timestamp(18446744066914187038)")
		})

		Convey("interpolate __timeToRFC3339 function", func() {
			sql, err := engine.Interpolate(timeRange, "select '$__timeToRFC3339()'")
			So(err, ShouldBeNil)

			So(sql, ShouldEqual, "select '0001-01-01T00:00:00Z'")
		})

		Convey("interpolate __unixEpochFilter function", func() {
			sql, err := engine.Interpolate(timeRange, "select $__unixEpochFilter(18446744066914186738)")
			So(err, ShouldBeNil)

			So(sql, ShouldEqual, "select 18446744066914186738 >= 18446744066914186738 AND 18446744066914186738 <= 18446744066914187038")
		})

		Convey("interpolate __unixEpochFrom function", func() {
			sql, err := engine.Interpolate(timeRange, "select $__unixEpochFrom()")
			So(err, ShouldBeNil)

			So(sql, ShouldEqual, "select 18446744066914186738")
		})

		Convey("interpolate __unixEpochTo function", func() {
			sql, err := engine.Interpolate(timeRange, "select $__unixEpochTo()")
			So(err, ShouldBeNil)

			So(sql, ShouldEqual, "select 18446744066914187038")
		})

	})
}
