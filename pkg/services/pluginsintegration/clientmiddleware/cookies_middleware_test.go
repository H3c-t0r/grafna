package clientmiddleware

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/grafana/grafana/pkg/infra/httpclient/httpclientprovider"
	"github.com/grafana/grafana/pkg/plugins/manager/client/clienttest"
	"github.com/grafana/grafana/pkg/services/user"
	"github.com/stretchr/testify/require"
)

func TestCookiesMiddleware(t *testing.T) {
	t.Run("When keepCookies not configured for a datasource", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/some/thing", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{
			Name: "cookie1",
		})
		req.AddCookie(&http.Cookie{
			Name: "cookie2",
		})
		req.AddCookie(&http.Cookie{
			Name: "cookie3",
		})

		cdt := clienttest.NewClientDecoratorTest(t,
			clienttest.WithReqContext(req, &user.SignedInUser{}),
			clienttest.WithMiddlewares(NewCookiesMiddleware([]string{"grafana_session"})),
		)

		jsonDataMap := map[string]interface{}{}
		jsonDataBytes, err := json.Marshal(&jsonDataMap)
		require.NoError(t, err)

		pluginCtx := backend.PluginContext{
			DataSourceInstanceSettings: &backend.DataSourceInstanceSettings{
				JSONData: jsonDataBytes,
			},
		}

		t.Run("Should not forward cookies when calling QueryData", func(t *testing.T) {
			_, err = cdt.Decorator.QueryData(req.Context(), &backend.QueryDataRequest{
				PluginContext: pluginCtx,
				Headers:       map[string]string{},
			})
			require.NoError(t, err)
			require.NotNil(t, cdt.QueryDataReq)
			require.Len(t, cdt.QueryDataReq.Headers, 0)

			middlewares := httpclient.ContextualMiddlewareFromContext(cdt.QueryDataCtx)
			require.Len(t, middlewares, 1)
			require.Equal(t, httpclientprovider.ForwardedCookiesMiddlewareName, middlewares[0].(httpclient.MiddlewareName).MiddlewareName())
		})

		t.Run("Should not forward cookies when calling CallResource", func(t *testing.T) {
			err = cdt.Decorator.CallResource(req.Context(), &backend.CallResourceRequest{
				PluginContext: pluginCtx,
				Headers:       map[string][]string{},
			}, nil)
			require.NoError(t, err)
			require.NotNil(t, cdt.CallResourceReq)
			require.Len(t, cdt.CallResourceReq.Headers, 0)

			middlewares := httpclient.ContextualMiddlewareFromContext(cdt.CallResourceCtx)
			require.Len(t, middlewares, 1)
			require.Equal(t, httpclientprovider.ForwardedCookiesMiddlewareName, middlewares[0].(httpclient.MiddlewareName).MiddlewareName())
		})

		t.Run("Should not forward cookies when calling CheckHealth", func(t *testing.T) {
			_, err = cdt.Decorator.CheckHealth(req.Context(), &backend.CheckHealthRequest{
				PluginContext: pluginCtx,
				Headers:       map[string]string{},
			})
			require.NoError(t, err)
			require.NotNil(t, cdt.CheckHealthReq)
			require.Len(t, cdt.CheckHealthReq.Headers, 0)

			middlewares := httpclient.ContextualMiddlewareFromContext(cdt.CheckHealthCtx)
			require.Len(t, middlewares, 1)
			require.Equal(t, httpclientprovider.ForwardedCookiesMiddlewareName, middlewares[0].(httpclient.MiddlewareName).MiddlewareName())
		})
	})

	t.Run("When keepCookies configured for a datasource", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/some/thing", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{
			Name: "cookie1",
		})
		req.AddCookie(&http.Cookie{
			Name: "cookie2",
		})
		req.AddCookie(&http.Cookie{
			Name: "cookie3",
		})
		req.AddCookie(&http.Cookie{
			Name: "grafana_session",
		})

		cdt := clienttest.NewClientDecoratorTest(t,
			clienttest.WithReqContext(req, &user.SignedInUser{}),
			clienttest.WithMiddlewares(NewCookiesMiddleware([]string{"grafana_session"})),
		)

		jsonDataMap := map[string]interface{}{
			"keepCookies": []string{"cookie2", "grafana_session"},
		}
		jsonDataBytes, err := json.Marshal(&jsonDataMap)
		require.NoError(t, err)

		pluginCtx := backend.PluginContext{
			DataSourceInstanceSettings: &backend.DataSourceInstanceSettings{
				JSONData: jsonDataBytes,
			},
		}

		t.Run("Should forward cookies when calling QueryData", func(t *testing.T) {
			_, err = cdt.Decorator.QueryData(req.Context(), &backend.QueryDataRequest{
				PluginContext: pluginCtx,
				Headers:       map[string]string{},
			})
			require.NoError(t, err)
			require.NotNil(t, cdt.QueryDataReq)
			require.Len(t, cdt.QueryDataReq.Headers, 1)
			require.EqualValues(t, "cookie2=", cdt.QueryDataReq.Headers["Cookie"])

			middlewares := httpclient.ContextualMiddlewareFromContext(cdt.QueryDataCtx)
			require.Len(t, middlewares, 1)
			require.Equal(t, httpclientprovider.ForwardedCookiesMiddlewareName, middlewares[0].(httpclient.MiddlewareName).MiddlewareName())
		})

		t.Run("Should forward cookies when calling CallResource", func(t *testing.T) {
			err = cdt.Decorator.CallResource(req.Context(), &backend.CallResourceRequest{
				PluginContext: pluginCtx,
				Headers:       map[string][]string{},
			}, nil)
			require.NoError(t, err)
			require.NotNil(t, cdt.CallResourceReq)
			require.Len(t, cdt.CallResourceReq.Headers, 1)
			require.Len(t, cdt.CallResourceReq.Headers["Cookie"], 1)
			require.EqualValues(t, "cookie2=", cdt.CallResourceReq.Headers["Cookie"][0])

			middlewares := httpclient.ContextualMiddlewareFromContext(cdt.CallResourceCtx)
			require.Len(t, middlewares, 1)
			require.Equal(t, httpclientprovider.ForwardedCookiesMiddlewareName, middlewares[0].(httpclient.MiddlewareName).MiddlewareName())
		})

		t.Run("Should forward cookies when calling CheckHealth", func(t *testing.T) {
			_, err = cdt.Decorator.CheckHealth(req.Context(), &backend.CheckHealthRequest{
				PluginContext: pluginCtx,
				Headers:       map[string]string{},
			})
			require.NoError(t, err)
			require.NotNil(t, cdt.CheckHealthReq)
			require.Len(t, cdt.CheckHealthReq.Headers, 1)
			require.EqualValues(t, "cookie2=", cdt.CheckHealthReq.Headers["Cookie"])

			middlewares := httpclient.ContextualMiddlewareFromContext(cdt.CheckHealthCtx)
			require.Len(t, middlewares, 1)
			require.Equal(t, httpclientprovider.ForwardedCookiesMiddlewareName, middlewares[0].(httpclient.MiddlewareName).MiddlewareName())
		})
	})
}
