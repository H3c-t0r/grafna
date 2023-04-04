package expr

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/grafana/dataplane/examples"
	"github.com/grafana/dataplane/sdata"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana/pkg/services/datasources"
	datafakes "github.com/grafana/grafana/pkg/services/datasources/fakes"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/util"
	"github.com/stretchr/testify/require"
)

func TestPassThroughDataplaneExamples(t *testing.T) {
	es, err := examples.GetExamples()
	require.NoError(t, err)

	validExamples, err := es.Filter(examples.FilterOptions{
		Version: data.FrameTypeVersion{0, 1},
		Valid:   util.Pointer(true),
	})
	require.NoError(t, err)

	for _, collection := range validExamples.Collections() {
		for _, example := range collection.ExampleSlice() {
			t.Run(example.Info().ID, func(t *testing.T) {
				_, err := framesPassThroughService(t, example.Frames("A"))
				require.NoError(t, err)
			})
		}
	}
}

func framesPassThroughService(t *testing.T, frames data.Frames) (data.Frames, error) {
	me := &mockEndpoint{
		Frames: frames,
	}

	cfg := setting.NewCfg()

	s := Service{
		cfg:               cfg,
		dataService:       me,
		dataSourceService: &datafakes.FakeDataSourceService{},
		features:          &featuremgmt.FeatureManager{},
	}
	queries := []Query{{
		RefID: "A",
		DataSource: &datasources.DataSource{
			OrgID: 1,
			UID:   "test",
			Type:  "test",
		},
		JSON: json.RawMessage(`{ "datasource": { "uid": "1" }, "intervalMs": 1000, "maxDataPoints": 1000 }`),
		TimeRange: AbsoluteTimeRange{
			From: time.Time{},
			To:   time.Time{},
		},
	}}

	req := &Request{Queries: queries}

	pl, err := s.BuildPipeline(req)
	require.NoError(t, err)

	res, err := s.ExecutePipeline(context.Background(), time.Now(), pl)
	require.NoError(t, err)

	require.Contains(t, res.Responses, "A")

	return res.Responses["A"].Frames, res.Responses["A"].Error
}

func TestShouldUseDataplane(t *testing.T) {
	t.Run("zero frames returns no data and is allowed", func(t *testing.T) {
		f := data.Frames{}
		k, use, err := shouldUseDataplane(f)
		require.NoError(t, err)
		require.True(t, use)
		require.Equal(t, data.KindUnknown, k)
	})

	t.Run("a frame with Type and TypeVersion 0.0 will not use dataplane", func(t *testing.T) {
		f := data.Frames{(&data.Frame{}).SetMeta(
			&data.FrameMeta{
				TypeVersion: data.FrameTypeVersion{},
				Type:        data.FrameTypeTimeSeriesMulti,
			},
		)}
		_, use, err := shouldUseDataplane(f)
		require.NoError(t, err)
		require.False(t, use)
	})

	t.Run("a frame without Type and TypeVersion > 0.0 will not use dataplane", func(t *testing.T) {
		f := data.Frames{(&data.Frame{}).SetMeta(
			&data.FrameMeta{
				TypeVersion: data.FrameTypeVersion{0, 1},
			},
		)}
		_, use, err := shouldUseDataplane(f)
		require.NoError(t, err)
		require.False(t, use)
	})

	t.Run("a frame with no metadata will not use dataplane", func(t *testing.T) {
		f := data.Frames{&data.Frame{}}
		_, use, err := shouldUseDataplane(f)
		require.NoError(t, err)
		require.False(t, use)
	})

	t.Run("a newer version that supported will return a warning but still use dataplane", func(t *testing.T) {
		ty := data.FrameTypeTimeSeriesMulti
		v := data.FrameTypeVersion{999, 999}
		f := data.Frames{(&data.Frame{}).SetMeta(
			&data.FrameMeta{
				Type:        ty,
				TypeVersion: v,
			},
		)}
		k, use, err := shouldUseDataplane(f)

		require.Error(t, err)

		var mockWarning *sdata.VersionWarning
		require.True(t, errors.As(err, &mockWarning))

		require.True(t, use)
		require.Equal(t, data.KindTimeSeries, k)
	})

	t.Run("all valid dataplane examples should use dataplane", func(t *testing.T) {
		es, err := examples.GetExamples()
		require.NoError(t, err)

		validExamples, err := es.Filter(examples.FilterOptions{
			Version: data.FrameTypeVersion{0, 1},
			Valid:   util.Pointer(true),
		})
		require.NoError(t, err)

		for _, collection := range validExamples.Collections() {
			for _, example := range collection.ExampleSlice() {
				t.Run(example.Info().ID, func(t *testing.T) {
					_, err := framesPassThroughService(t, example.Frames("A"))
					require.NoError(t, err)
				})
			}
		}
	})
}
