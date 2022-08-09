package querydata_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	sdkhttpclient "github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana/pkg/infra/httpclient"
	"github.com/grafana/grafana/pkg/infra/tracing"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/tsdb/prometheus/buffered"
	"github.com/grafana/grafana/pkg/tsdb/prometheus/models"
	"github.com/grafana/grafana/pkg/tsdb/prometheus/querydata"
	apiv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	p "github.com/prometheus/common/model"
	"github.com/stretchr/testify/require"
)

func TestPrometheus_parseTimeSeriesResponse(t *testing.T) {
	t.Run("exemplars response should be sampled and parsed normally", func(t *testing.T) {
		t.Skip()
		exemplars := []apiv1.ExemplarQueryResult{
			{
				SeriesLabels: p.LabelSet{
					"__name__": "tns_request_duration_seconds_bucket",
					"instance": "app:80",
					"job":      "tns/app",
				},
				Exemplars: []apiv1.Exemplar{
					{
						Labels:    p.LabelSet{"traceID": "test1"},
						Value:     0.003535405,
						Timestamp: p.TimeFromUnixNano(time.Now().Add(-2 * time.Minute).UnixNano()),
					},
					{
						Labels:    p.LabelSet{"traceID": "test2"},
						Value:     0.005555605,
						Timestamp: p.TimeFromUnixNano(time.Now().Add(-4 * time.Minute).UnixNano()),
					},
					{
						Labels:    p.LabelSet{"traceID": "test3"},
						Value:     0.007545445,
						Timestamp: p.TimeFromUnixNano(time.Now().Add(-6 * time.Minute).UnixNano()),
					},
					{
						Labels:    p.LabelSet{"traceID": "test4"},
						Value:     0.009545445,
						Timestamp: p.TimeFromUnixNano(time.Now().Add(-7 * time.Minute).UnixNano()),
					},
				},
			},
		}

		tctx, err := setup(true)
		require.NoError(t, err)

		qm := models.QueryModel{
			LegendFormat:  "legend {{app}}",
			UtcOffsetSec:  0,
			ExemplarQuery: true,
		}
		b, err := json.Marshal(&qm)
		require.NoError(t, err)
		query := backend.DataQuery{
			RefID: "A",
			JSON:  b,
		}
		res, err := execute(tctx, query, exemplars)
		require.NoError(t, err)

		// Test fields
		require.Len(t, res, 1)
		//		require.Equal(t, res[0].Name, "exemplar")
		require.Equal(t, res[0].Fields[0].Name, "Time")
		require.Equal(t, res[0].Fields[1].Name, "Value")
		require.Len(t, res[0].Fields, 6)

		// Test correct values (sampled to 2)
		require.Equal(t, res[0].Fields[1].Len(), 2)
		require.Equal(t, res[0].Fields[1].At(0), 0.009545445)
		require.Equal(t, res[0].Fields[1].At(1), 0.003535405)
	})

	t.Run("matrix response should be parsed normally", func(t *testing.T) {
		values := []p.SamplePair{
			{Value: 1, Timestamp: 1000},
			{Value: 2, Timestamp: 2000},
			{Value: 3, Timestamp: 3000},
			{Value: 4, Timestamp: 4000},
			{Value: 5, Timestamp: 5000},
		}
		result := queryResult{
			Type: p.ValMatrix,
			Result: p.Matrix{
				&p.SampleStream{
					Metric: p.Metric{"app": "Application", "tag2": "tag2"},
					Values: values,
				},
			},
		}

		qm := models.QueryModel{
			LegendFormat: "legend {{app}}",
			UtcOffsetSec: 0,
			RangeQuery:   true,
		}
		b, err := json.Marshal(&qm)
		require.NoError(t, err)
		query := backend.DataQuery{
			TimeRange: backend.TimeRange{
				From: time.Unix(1, 0).UTC(),
				To:   time.Unix(5, 0).UTC(),
			},
			JSON: b,
		}
		tctx, err := setup(true)
		require.NoError(t, err)
		res, err := execute(tctx, query, result)
		require.NoError(t, err)

		require.Len(t, res, 1)
		require.Len(t, res[0].Fields, 2)
		require.Len(t, res[0].Fields[0].Labels, 0)
		require.Equal(t, "Time", res[0].Fields[0].Name)
		require.Len(t, res[0].Fields[1].Labels, 2)
		require.Equal(t, "app=Application, tag2=tag2", res[0].Fields[1].Labels.String())
		require.Equal(t, "legend Application", res[0].Fields[1].Name)

		// Ensure the timestamps are UTC zoned
		testValue := res[0].Fields[0].At(0)
		require.Equal(t, "UTC", testValue.(time.Time).Location().String())
	})

	t.Run("matrix response with missed data points should be parsed correctly", func(t *testing.T) {
		values := []p.SamplePair{
			{Value: 1, Timestamp: 1000},
			{Value: 4, Timestamp: 4000},
		}
		result := queryResult{
			Type: p.ValMatrix,
			Result: p.Matrix{
				&p.SampleStream{
					Metric: p.Metric{"app": "Application", "tag2": "tag2"},
					Values: values,
				},
			},
		}

		qm := models.QueryModel{
			LegendFormat: "",
			UtcOffsetSec: 0,
			RangeQuery:   true,
		}
		b, err := json.Marshal(&qm)
		require.NoError(t, err)
		query := backend.DataQuery{
			TimeRange: backend.TimeRange{
				From: time.Unix(1, 0).UTC(),
				To:   time.Unix(4, 0).UTC(),
			},
			JSON: b,
		}
		tctx, err := setup(true)
		require.NoError(t, err)
		res, err := execute(tctx, query, result)

		require.NoError(t, err)
		require.Len(t, res, 1)
		require.Equal(t, res[0].Fields[0].Len(), 2)
		require.Equal(t, time.Unix(1, 0).UTC(), res[0].Fields[0].At(0))
		require.Equal(t, time.Unix(4, 0).UTC(), res[0].Fields[0].At(1))
		require.Equal(t, res[0].Fields[1].Len(), 2)
		require.Equal(t, float64(1), *res[0].Fields[1].At(0).(*float64))
		require.Equal(t, float64(4), *res[0].Fields[1].At(1).(*float64))
	})

	t.Run("matrix response with from alerting missed data points should be parsed correctly", func(t *testing.T) {
		values := []p.SamplePair{
			{Value: 1, Timestamp: 1000},
			{Value: 4, Timestamp: 4000},
		}
		result := queryResult{
			Type: p.ValMatrix,
			Result: p.Matrix{
				&p.SampleStream{
					Metric: p.Metric{"app": "Application", "tag2": "tag2"},
					Values: values,
				},
			},
		}

		qm := models.QueryModel{
			LegendFormat: "",
			UtcOffsetSec: 0,
			RangeQuery:   true,
		}
		b, err := json.Marshal(&qm)
		require.NoError(t, err)
		query := backend.DataQuery{
			TimeRange: backend.TimeRange{
				From: time.Unix(1, 0).UTC(),
				To:   time.Unix(4, 0).UTC(),
			},
			JSON: b,
		}
		tctx, err := setup(true)
		require.NoError(t, err)
		res, err := execute(tctx, query, result)

		require.NoError(t, err)
		require.Len(t, res, 1)
		require.Len(t, res[0].Fields, 2)
		require.Len(t, res[0].Fields[0].Labels, 0)
		require.Equal(t, res[0].Fields[0].Name, "Time")
		require.Len(t, res[0].Fields[1].Labels, 2)
		require.Equal(t, res[0].Fields[1].Labels.String(), "app=Application, tag2=tag2")
		require.Equal(t, "{app=\"Application\", tag2=\"tag2\"}", res[0].Fields[1].Name)
	})

	t.Run("matrix response with NaN value should be changed to null", func(t *testing.T) {
		result := queryResult{
			Type: p.ValMatrix,
			Result: p.Matrix{
				&p.SampleStream{
					Metric: p.Metric{"app": "Application"},
					Values: []p.SamplePair{
						{Value: p.SampleValue(math.NaN()), Timestamp: 1000},
					},
				},
			},
		}

		qm := models.QueryModel{
			LegendFormat: "",
			UtcOffsetSec: 0,
			RangeQuery:   true,
		}
		b, err := json.Marshal(&qm)
		require.NoError(t, err)
		query := backend.DataQuery{
			TimeRange: backend.TimeRange{
				From: time.Unix(1, 0).UTC(),
				To:   time.Unix(4, 0).UTC(),
			},
			JSON: b,
		}

		tctx, err := setup(true)
		require.NoError(t, err)
		res, err := execute(tctx, query, result)
		require.NoError(t, err)

		require.Equal(t, "{app=\"Application\"}", res[0].Fields[1].Name)
		require.True(t, math.IsNaN(*res[0].Fields[1].At(0).(*float64)))
	})

	t.Run("vector response should be parsed normally", func(t *testing.T) {
		qr := queryResult{
			Type: p.ValVector,
			Result: p.Vector{
				&p.Sample{
					Metric:    p.Metric{"app": "Application", "tag2": "tag2"},
					Value:     1,
					Timestamp: 123,
				},
			},
		}
		qm := models.QueryModel{
			LegendFormat: "legend {{app}}",
			UtcOffsetSec: 0,
			InstantQuery: true,
		}
		b, err := json.Marshal(&qm)
		require.NoError(t, err)
		query := backend.DataQuery{
			JSON: b,
		}
		tctx, err := setup(true)
		require.NoError(t, err)
		res, err := execute(tctx, query, qr)
		require.NoError(t, err)

		require.Len(t, res, 1)
		require.Len(t, res[0].Fields, 2)
		require.Len(t, res[0].Fields[0].Labels, 0)
		require.Equal(t, res[0].Fields[0].Name, "Time")
		require.Equal(t, res[0].Fields[0].Name, "Time")
		require.Len(t, res[0].Fields[1].Labels, 2)
		require.Equal(t, res[0].Fields[1].Labels.String(), "app=Application, tag2=tag2")
		require.Equal(t, "legend Application", res[0].Fields[1].Name)

		// Ensure the timestamps are UTC zoned
		testValue := res[0].Fields[0].At(0)
		require.Equal(t, "UTC", testValue.(time.Time).Location().String())
		require.Equal(t, int64(123), testValue.(time.Time).UnixMilli())
	})

	t.Run("scalar response should be parsed normally", func(t *testing.T) {
		t.Skip("TODO: implement scalar responses")
		qr := queryResult{
			Type: p.ValScalar,
			Result: &p.Scalar{
				Value:     1,
				Timestamp: 123,
			},
		}
		qm := models.QueryModel{
			LegendFormat: "",
			UtcOffsetSec: 0,
			InstantQuery: true,
		}
		b, err := json.Marshal(&qm)
		require.NoError(t, err)
		query := backend.DataQuery{
			JSON: b,
		}
		tctx, err := setup(true)
		require.NoError(t, err)
		res, err := execute(tctx, query, qr)
		require.NoError(t, err)

		require.Len(t, res, 1)
		require.Len(t, res[0].Fields, 2)
		require.Len(t, res[0].Fields[0].Labels, 0)
		require.Equal(t, res[0].Fields[0].Name, "Time")
		require.Equal(t, "1", res[0].Fields[1].Name)

		// Ensure the timestamps are UTC zoned
		testValue := res[0].Fields[0].At(0)
		require.Equal(t, "UTC", testValue.(time.Time).Location().String())
		require.Equal(t, int64(123), testValue.(time.Time).UnixMilli())
	})
}

type queryResult struct {
	Type   p.ValueType `json:"resultType"`
	Result interface{} `json:"result"`
}

func execute(tctx *testContext, query backend.DataQuery, qr interface{}) (data.Frames, error) {
	req := backend.QueryDataRequest{
		Queries: []backend.DataQuery{query},
		Headers: map[string]string{},
	}

	promRes, err := toAPIResponse(qr)
	if err != nil {
		return nil, err
	}
	tctx.httpProvider.setResponse(promRes)

	res, err := tctx.queryData.Execute(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	return res.Responses[req.Queries[0].RefID].Frames, nil
}

type apiResponse struct {
	Status string          `json:"status"`
	Data   json.RawMessage `json:"data"`
}

func toAPIResponse(d interface{}) (*http.Response, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	res := apiResponse{
		Status: "success",
		Data:   json.RawMessage(b),
	}

	raw, err := json.Marshal(&res)
	if err != nil {
		return nil, err
	}

	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(raw)),
	}, nil
}

type testContext struct {
	httpProvider *fakeHttpClientProvider
	queryData    *querydata.QueryData
}

func setup(wideFrames bool) (*testContext, error) {
	tracer := tracing.InitializeTracerForTest()
	httpProvider := &fakeHttpClientProvider{
		opts: sdkhttpclient.Options{
			Timeouts: &sdkhttpclient.DefaultTimeoutOptions,
		},
		res: &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{}`))),
		},
	}
	settings := backend.DataSourceInstanceSettings{
		URL:      "http://localhost:9090",
		JSONData: json.RawMessage(`{"timeInterval": "15s"}`),
	}

	features := &fakeFeatureToggles{flags: map[string]bool{"prometheusStreamingJSONParser": true, "prometheusWideSeries": wideFrames}}

	opts, err := buffered.CreateTransportOptions(settings, &setting.Cfg{}, &fakeLogger{})
	if err != nil {
		return nil, err
	}

	httpClient, err := httpProvider.New(*opts)
	if err != nil {
		return nil, err
	}

	queryData, _ := querydata.New(httpClient, features, tracer, settings, &fakeLogger{})

	return &testContext{
		httpProvider: httpProvider,
		queryData:    queryData,
	}, nil
}

type fakeFeatureToggles struct {
	flags map[string]bool
}

func (f *fakeFeatureToggles) IsEnabled(feature string) bool {
	return f.flags[feature]
}

type fakeHttpClientProvider struct {
	httpclient.Provider
	opts sdkhttpclient.Options
	res  *http.Response
}

func (p *fakeHttpClientProvider) New(opts ...sdkhttpclient.Options) (*http.Client, error) {
	p.opts = opts[0]
	c, err := sdkhttpclient.New(opts[0])
	if err != nil {
		return nil, err
	}
	c.Transport = p
	return c, nil
}

func (p *fakeHttpClientProvider) GetTransport(opts ...sdkhttpclient.Options) (http.RoundTripper, error) {
	p.opts = opts[0]
	return http.DefaultTransport, nil
}

func (p *fakeHttpClientProvider) setResponse(res *http.Response) {
	p.res = res
}

func (p *fakeHttpClientProvider) RoundTrip(req *http.Request) (*http.Response, error) {
	return p.res, nil
}
