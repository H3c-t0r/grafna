package mathexp

import (
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/stretchr/testify/assert"
)

func TestSeriesSort(t *testing.T) {
	var tests = []struct {
		name           string
		descending     bool
		series         Series
		sortedSeriesIs assert.ComparisonAssertionFunc
		sortedSeries   Series
		panics         assert.PanicTestFunc
	}{
		{
			name:       "unordered series should sort by time ascending",
			descending: false,
			series: makeSeries("", nil, tp{
				time.Unix(3, 0), float64Pointer(3),
			}, tp{
				time.Unix(1, 0), float64Pointer(1),
			}, tp{
				time.Unix(2, 0), float64Pointer(2),
			}),
			sortedSeriesIs: assert.Equal,
			sortedSeries: makeSeries("", nil, tp{
				time.Unix(1, 0), float64Pointer(1),
			}, tp{
				time.Unix(2, 0), float64Pointer(2),
			}, tp{
				time.Unix(3, 0), float64Pointer(3),
			}),
		},
		{
			name:       "unordered series should sort by time descending",
			descending: true,
			series: makeSeries("", nil, tp{
				time.Unix(3, 0), float64Pointer(3),
			}, tp{
				time.Unix(1, 0), float64Pointer(1),
			}, tp{
				time.Unix(2, 0), float64Pointer(2),
			}),
			sortedSeriesIs: assert.Equal,
			sortedSeries: makeSeries("", nil, tp{
				time.Unix(3, 0), float64Pointer(3),
			}, tp{
				time.Unix(2, 0), float64Pointer(2),
			}, tp{
				time.Unix(1, 0), float64Pointer(1),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.series.SortByTime(tt.descending)
			tt.sortedSeriesIs(t, tt.series, tt.sortedSeries)
		})
	}
}

func TestSeriesFromFrame(t *testing.T) {
	var tests = []struct {
		name   string
		frame  *data.Frame
		errIs  assert.ErrorAssertionFunc
		Is     assert.ComparisonAssertionFunc
		Series Series
	}{
		{
			name: "[]time, []float frame should convert",
			frame: &data.Frame{
				Name: "test",
				Fields: []*data.Field{
					data.NewField("time", nil, []time.Time{}),
					data.NewField("value", nil, []float64{}),
				},
			},
			errIs: assert.NoError,
			Is:    assert.Equal,
			Series: Series{
				Frame: &data.Frame{
					Name: "test",
					Fields: []*data.Field{
						data.NewField("time", nil, []time.Time{}),
						data.NewField("value", nil, []*float64{}),
					},
				},
			},
		},
		{
			name: "[]time, []*float frame should convert",
			frame: &data.Frame{
				Name: "test",
				Fields: []*data.Field{
					data.NewField("time", nil, []time.Time{time.Unix(5, 0)}),
					data.NewField("value", nil, []*float64{float64Pointer(5)}),
				},
			},
			errIs: assert.NoError,
			Is:    assert.Equal,
			Series: Series{
				Frame: &data.Frame{
					Name: "test",
					Fields: []*data.Field{
						data.NewField("time", nil, []time.Time{time.Unix(5, 0)}),
						data.NewField("value", nil, []*float64{float64Pointer(5)}),
					},
				},
			},
		},
		{
			name: "[]*float, []time frame should convert",
			frame: &data.Frame{
				Name: "test",
				Fields: []*data.Field{
					data.NewField("value", nil, []*float64{float64Pointer(5)}),
					data.NewField("time", nil, []time.Time{time.Unix(5, 0)}),
				},
			},
			errIs: assert.NoError,
			Is:    assert.Equal,
			Series: Series{
				Frame: &data.Frame{
					Name: "test",
					Fields: []*data.Field{
						data.NewField("time", nil, []time.Time{time.Unix(5, 0)}),
						data.NewField("value", nil, []*float64{float64Pointer(5)}),
					},
				},
			},
		},
		{
			name: "[]*int, []*time frame should convert",
			frame: &data.Frame{
				Name: "test",
				Fields: []*data.Field{
					data.NewField("value", nil, []*int64{int64Pointer(5)}),
					data.NewField("time", nil, []*time.Time{unixTimePointer(5, 0)}),
				},
			},
			errIs: assert.NoError,
			Is:    assert.Equal,
			Series: Series{
				Frame: &data.Frame{
					Name: "test",
					Fields: []*data.Field{
						data.NewField("value", nil, []*float64{float64Pointer(5)}),
						data.NewField("time", nil, []*time.Time{unixTimePointer(5, 0)}),
					},
				},
				TimeIdx:         1,
				TimeIsNullable:  true,
				ValueIdx:        0,
				ValueIsNullable: true,
			},
		},
		{
			name: "[]int, []*time frame should convert",
			frame: &data.Frame{
				Name: "test",
				Fields: []*data.Field{
					data.NewField("value", nil, []int64{5}),
					data.NewField("time", nil, []*time.Time{unixTimePointer(5, 0)}),
				},
			},
			errIs: assert.NoError,
			Is:    assert.Equal,
			Series: Series{
				Frame: &data.Frame{
					Name: "test",
					Fields: []*data.Field{
						data.NewField("value", nil, []*float64{float64Pointer(5)}),
						data.NewField("time", nil, []*time.Time{unixTimePointer(5, 0)}),
					},
				},
				TimeIdx:         1,
				TimeIsNullable:  true,
				ValueIdx:        0,
				ValueIsNullable: true,
			},
		},
		{
			name: "[]string, []*time frame should convert",
			frame: &data.Frame{
				Name: "test",
				Fields: []*data.Field{
					data.NewField("value", nil, []string{"5"}),
					data.NewField("time", nil, []*time.Time{unixTimePointer(5, 0)}),
				},
			},
			errIs: assert.NoError,
			Is:    assert.Equal,
			Series: Series{
				Frame: &data.Frame{
					Name: "test",
					Fields: []*data.Field{
						data.NewField("value", nil, []*float64{float64Pointer(5)}),
						data.NewField("time", nil, []*time.Time{unixTimePointer(5, 0)}),
					},
				},
				TimeIdx:         1,
				TimeIsNullable:  true,
				ValueIdx:        0,
				ValueIsNullable: true,
			},
		},
		{
			name: "[]*string, []*time frame should convert",
			frame: &data.Frame{
				Name: "test",
				Fields: []*data.Field{
					data.NewField("value", nil, []*string{strPointer("5")}),
					data.NewField("time", nil, []*time.Time{unixTimePointer(5, 0)}),
				},
			},
			errIs: assert.NoError,
			Is:    assert.Equal,
			Series: Series{
				Frame: &data.Frame{
					Name: "test",
					Fields: []*data.Field{
						data.NewField("value", nil, []*float64{float64Pointer(5)}),
						data.NewField("time", nil, []*time.Time{unixTimePointer(5, 0)}),
					},
				},
				TimeIdx:         1,
				TimeIsNullable:  true,
				ValueIdx:        0,
				ValueIsNullable: true,
			},
		},
		{
			name: "[]bool, []*time frame should convert",
			frame: &data.Frame{
				Name: "test",
				Fields: []*data.Field{
					data.NewField("value", nil, []bool{true}),
					data.NewField("time", nil, []*time.Time{unixTimePointer(5, 0)}),
				},
			},
			errIs: assert.NoError,
			Is:    assert.Equal,
			Series: Series{
				Frame: &data.Frame{
					Name: "test",
					Fields: []*data.Field{
						data.NewField("value", nil, []*float64{float64Pointer(1)}),
						data.NewField("time", nil, []*time.Time{unixTimePointer(5, 0)}),
					},
				},
				TimeIdx:         1,
				TimeIsNullable:  true,
				ValueIdx:        0,
				ValueIsNullable: true,
			},
		},
		{
			name: "[]*bool, []*time frame should convert",
			frame: &data.Frame{
				Name: "test",
				Fields: []*data.Field{
					data.NewField("value", nil, []*bool{boolPointer(true)}),
					data.NewField("time", nil, []*time.Time{unixTimePointer(5, 0)}),
				},
			},
			errIs: assert.NoError,
			Is:    assert.Equal,
			Series: Series{
				Frame: &data.Frame{
					Name: "test",
					Fields: []*data.Field{
						data.NewField("value", nil, []*float64{float64Pointer(1)}),
						data.NewField("time", nil, []*time.Time{unixTimePointer(5, 0)}),
					},
				},
				TimeIdx:         1,
				TimeIsNullable:  true,
				ValueIdx:        0,
				ValueIsNullable: true,
			},
		},
		{
			name: "[]*time, []*time frame should error",
			frame: &data.Frame{
				Name: "test",
				Fields: []*data.Field{
					data.NewField("time", nil, []*time.Time{}),
					data.NewField("time", nil, []*time.Time{}),
				},
			},
			errIs: assert.Error,
		},
		{
			name: "[]*float64, []float64 frame should error",
			frame: &data.Frame{
				Name: "test",
				Fields: []*data.Field{
					data.NewField("value", nil, []*float64{}),
					data.NewField("value", nil, []*float64{}),
				},
			},
			errIs: assert.Error,
		},
		{
			name: "[]*float64 frame should error",
			frame: &data.Frame{
				Name: "test",
				Fields: []*data.Field{
					data.NewField("value", nil, []*float64{}),
				},
			},
			errIs: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := SeriesFromFrame(tt.frame)
			tt.errIs(t, err)
			if err == nil {
				tt.Is(t, s, tt.Series)
			}
		})
	}
}
