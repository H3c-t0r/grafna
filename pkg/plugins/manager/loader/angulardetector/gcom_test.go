package angulardetector

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var mockGCOMResponse = []byte(`[{
	"name": "PanelCtrl",
	"type": "contains",
	"pattern": "PanelCtrl"
},
{
    "name": "QueryCtrl",
    "type": "regex",
    "pattern": "[\"']QueryCtrl[\"']"
}]`)

func mockGCOMHTTPHandlerFunc(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/api/plugins/angular_patterns" {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	_, _ = writer.Write(mockGCOMResponse)
}

func checkMockGCOMResponse(t *testing.T, detectors []detector) {
	require.Len(t, detectors, 2)
	d, ok := detectors[0].(*containsBytesDetector)
	require.True(t, ok)
	require.Equal(t, []byte(`PanelCtrl`), d.pattern)
	rd, ok := detectors[1].(*regexDetector)
	require.True(t, ok)
	require.Equal(t, `["']QueryCtrl["']`, rd.regex.String())
}

type gcomScenario struct {
	gcomHTTPHandlerFunc http.HandlerFunc
	gcomHTTPCalls       int
}

func (s *gcomScenario) newHTTPTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.gcomHTTPCalls++
		s.gcomHTTPHandlerFunc(w, r)
	}))
}

func newDefaultGCOMScenario() *gcomScenario {
	return &gcomScenario{gcomHTTPHandlerFunc: mockGCOMHTTPHandlerFunc}
}

func newError500GCOMScenario() *gcomScenario {
	return &gcomScenario{gcomHTTPHandlerFunc: func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
	}}
}

func TestGCOMDetectorsProvider(t *testing.T) {
	t.Run("returns value returned from gcom api", func(t *testing.T) {
		scenario := newDefaultGCOMScenario()
		srv := scenario.newHTTPTestServer()
		t.Cleanup(srv.Close)
		gcomProvider, err := newGCOMDetectorsProvider(srv.URL, defaultGCOMDetectorsProviderTTL)
		require.NoError(t, err)
		detectors := gcomProvider.provideDetectors(context.Background())
		require.Equal(t, 1, scenario.gcomHTTPCalls, "gcom api should be called")
		checkMockGCOMResponse(t, detectors)
	})

	t.Run("uses cache when called multiple times", func(t *testing.T) {
		scenario := newDefaultGCOMScenario()
		srv := scenario.newHTTPTestServer()
		t.Cleanup(srv.Close)
		gcomProvider, err := newGCOMDetectorsProvider(srv.URL, defaultGCOMDetectorsProviderTTL)
		require.NoError(t, err)
		detectors := gcomProvider.provideDetectors(context.Background())
		require.Equal(t, 1, scenario.gcomHTTPCalls, "gcom api should be called")
		checkMockGCOMResponse(t, detectors)

		secondDetectors := gcomProvider.provideDetectors(context.Background())
		require.Equal(t, 1, scenario.gcomHTTPCalls, "gcom api should be called only once")
		require.Equal(t, detectors, secondDetectors)
	})

	t.Run("calls gcom again when cache expires", func(t *testing.T) {
		scenario := newDefaultGCOMScenario()
		srv := scenario.newHTTPTestServer()
		t.Cleanup(srv.Close)
		gcomProvider, err := newGCOMDetectorsProvider(
			srv.URL,
			// Cache expires after 1 us
			time.Microsecond*1,
		)

		require.NoError(t, err)
		detectors := gcomProvider.provideDetectors(context.Background())
		checkMockGCOMResponse(t, detectors)
		require.Equal(t, 1, scenario.gcomHTTPCalls, "gcom api should be called")

		// Wait for cache to expire
		time.Sleep(time.Microsecond * 2)
		newDetectors := gcomProvider.provideDetectors(context.Background())
		checkMockGCOMResponse(t, newDetectors)
		require.Equal(t, 2, scenario.gcomHTTPCalls, "gcom api should be called again after cache expires")
	})

	t.Run("error handling", func(t *testing.T) {
		for _, tc := range []struct {
			*gcomScenario
			name string
		}{
			{name: "http error 500", gcomScenario: &gcomScenario{
				gcomHTTPHandlerFunc: func(writer http.ResponseWriter, request *http.Request) {
					writer.WriteHeader(http.StatusInternalServerError)
				},
			}},
			{name: "invalid json", gcomScenario: &gcomScenario{
				gcomHTTPHandlerFunc: func(writer http.ResponseWriter, request *http.Request) {
					_, _ = writer.Write([]byte(`not json`))
				},
			}},
			{name: "invalid regex", gcomScenario: &gcomScenario{
				gcomHTTPHandlerFunc: func(writer http.ResponseWriter, request *http.Request) {
					_, _ = writer.Write([]byte(`[{"name": "test", "type": "regex", "pattern": "((("}]`))
				},
			}},
		} {
			t.Run(tc.name, func(t *testing.T) {
				srv := tc.newHTTPTestServer()
				t.Cleanup(srv.Close)
				gcomProvider, err := newGCOMDetectorsProvider(srv.URL, defaultGCOMDetectorsProviderTTL)
				require.NoError(t, err)
				detectors := gcomProvider.provideDetectors(context.Background())
				require.Equal(t, 1, tc.gcomHTTPCalls, "gcom should be called")
				require.Empty(t, detectors, "returned detectors should be empty")
			})
		}
	})

	t.Run("handles gcom timeout", func(t *testing.T) {
		gcomScenario := &gcomScenario{
			gcomHTTPHandlerFunc: func(writer http.ResponseWriter, request *http.Request) {
				time.Sleep(time.Second * 1)
				_, _ = writer.Write([]byte(`[{"name": "test", "type": "regex", "pattern": "((("}]`))
			},
		}
		srv := gcomScenario.newHTTPTestServer()
		t.Cleanup(srv.Close)
		gcomProvider, err := newGCOMDetectorsProvider(srv.URL, defaultGCOMDetectorsProviderTTL)
		require.NoError(t, err)
		// Expired context
		ctx, canc := context.WithTimeout(context.Background(), time.Second*-1)
		defer canc()
		detectors := gcomProvider.provideDetectors(ctx)
		require.Zero(t, gcomScenario.gcomHTTPCalls, "gcom should be not called due to request timing out")
		require.Empty(t, detectors, "returned detectors should be empty")
	})

	t.Run("caches gcom error", func(t *testing.T) {
		scenario := newError500GCOMScenario()
		srv := scenario.newHTTPTestServer()
		t.Cleanup(srv.Close)
		gcomProvider, err := newGCOMDetectorsProvider(srv.URL, defaultGCOMDetectorsProviderTTL)
		require.NoError(t, err)
		detectors := gcomProvider.provideDetectors(context.Background())
		require.Equal(t, 1, scenario.gcomHTTPCalls, "gcom should be called")
		require.Empty(t, detectors, "returned detectors should be empty")

		newDetectors := gcomProvider.provideDetectors(context.Background())
		require.Equal(t, 1, scenario.gcomHTTPCalls, "gcom should not be called while cache is valid")
		require.Equal(t, detectors, newDetectors, "second call should return the same response")
	})

	t.Run("unknown pattern types do not break decoding", func(t *testing.T) {
		// Tests that we can introduce new pattern types in the future without breaking old Grafana versions.

		scenario := gcomScenario{gcomHTTPHandlerFunc: func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(`[
				{"name": "PanelCtrl", "type": "contains", "pattern": "PanelCtrl"},
				{"name": "Another", "type": "unknown", "pattern": "PanelCtrl"}
			]`))
		}}
		srv := scenario.newHTTPTestServer()
		t.Cleanup(srv.Close)
		gcomProvider, err := newGCOMDetectorsProvider(srv.URL, defaultGCOMDetectorsProviderTTL)
		require.NoError(t, err)
		detectors := gcomProvider.provideDetectors(context.Background())
		require.Equal(t, 1, scenario.gcomHTTPCalls, "gcom should be called")
		require.Len(t, detectors, 1, "should have decoded only 1 detector")
		d, ok := detectors[0].(*containsBytesDetector)
		require.True(t, ok, "decoded pattern should be of the correct type")
		require.Equal(t, []byte("PanelCtrl"), d.pattern, "decoded value for known pattern should be correct")
	})
}
