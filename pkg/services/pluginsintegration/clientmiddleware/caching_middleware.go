package clientmiddleware

import (
	"context"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/plugins"
	"github.com/grafana/grafana/pkg/services/caching"
	"github.com/grafana/grafana/pkg/services/contexthandler"
	"github.com/prometheus/client_golang/prometheus"
)

// NewCachingMiddleware creates a new plugins.ClientMiddleware that will
// attempt to read and write query results to the cache
func NewCachingMiddleware(cachingService caching.CachingService) plugins.ClientMiddleware {
	return plugins.ClientMiddlewareFunc(func(next plugins.Client) plugins.Client {
		log := log.New("caching_middleware")
		if err := prometheus.Register(QueryRequestHistogram); err != nil {
			log.Error("error registering prometheus collector", "error", err)
		}
		if err := prometheus.Register(ResourceRequestHistogram); err != nil {
			log.Error("error registering prometheus collector", "error", err)
		}
		return &CachingMiddleware{
			next:    next,
			caching: cachingService,
			log:     log,
		}
	})
}

type CachingMiddleware struct {
	next    plugins.Client
	caching caching.CachingService
	log     log.Logger
}

// QueryData receives a data request and attempts to access results already stored in the cache for that request.
// If data is found, it will return it immediately. Otherwise, it will perform the queries as usual, then write the response to the cache.
// If the cache service is implemented, we capture the request duration as a metric. The service is expected to write any response headers.
func (m *CachingMiddleware) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	if req == nil {
		return m.next.QueryData(ctx, req)
	}

	reqCtx := contexthandler.FromContext(ctx)
	if reqCtx == nil {
		return m.next.QueryData(ctx, req)
	}

	// time how long this request takes
	start := time.Now()

	// First look in the query cache if enabled
	cr := m.caching.HandleQueryRequest(ctx, req)
	ch := reqCtx.Resp.Header().Get(caching.XCacheHeader)

	defer func() {
		// record request duration if caching was used
		if ch != "" {
			QueryRequestHistogram.With(prometheus.Labels{
				"datasource_type": req.PluginContext.DataSourceInstanceSettings.Type,
				"cache":           ch,
				"query_type":      getQueryType(reqCtx),
			}).Observe(time.Since(start).Seconds())
		}
	}()

	// Cache hit; return the response
	if ch == caching.StatusHit {
		return cr.Response, nil
	}

	// Cache miss; do the actual queries
	resp, err := m.next.QueryData(ctx, req)

	// Update the query cache with the result for this metrics request
	if err == nil && cr.UpdateCacheFn != nil {
		cr.UpdateCacheFn(ctx, resp)
	}

	return resp, err
}

// CallResource receives a resource request and attempts to access results already stored in the cache for that request.
// If data is found, it will return it immediately. Otherwise, it will perform the request as usual. The caller of CallResource is expected to explicitly update the cache with any responses.
// If the cache service is implemented, we capture the request duration as a metric. The service is expected to write any response headers.
func (m *CachingMiddleware) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	if req == nil {
		return m.next.CallResource(ctx, req, sender)
	}

	reqCtx := contexthandler.FromContext(ctx)
	if reqCtx == nil {
		return m.next.CallResource(ctx, req, sender)
	}

	// time how long this request takes
	start := time.Now()

	// First look in the resource cache if enabled
	resp := m.caching.HandleResourceRequest(ctx, req)
	ch := reqCtx.Resp.Header().Get(caching.XCacheHeader)

	defer func() {
		// record request duration if caching was used
		if ch != "" {
			ResourceRequestHistogram.With(prometheus.Labels{
				"datasource_type": req.PluginContext.DataSourceInstanceSettings.Type,
				"cache":           resp.Headers[caching.XCacheHeader][0],
			}).Observe(time.Since(start).Seconds())
		}
	}()

	// Cache hit; send the response and return
	if ch == caching.StatusHit {
		return sender.Send(resp)
	}

	// Cache miss; do the actual request
	return m.next.CallResource(ctx, req, sender)
}

func (m *CachingMiddleware) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	return m.next.CheckHealth(ctx, req)
}

func (m *CachingMiddleware) CollectMetrics(ctx context.Context, req *backend.CollectMetricsRequest) (*backend.CollectMetricsResult, error) {
	return m.next.CollectMetrics(ctx, req)
}

func (m *CachingMiddleware) SubscribeStream(ctx context.Context, req *backend.SubscribeStreamRequest) (*backend.SubscribeStreamResponse, error) {
	return m.next.SubscribeStream(ctx, req)
}

func (m *CachingMiddleware) PublishStream(ctx context.Context, req *backend.PublishStreamRequest) (*backend.PublishStreamResponse, error) {
	return m.next.PublishStream(ctx, req)
}

func (m *CachingMiddleware) RunStream(ctx context.Context, req *backend.RunStreamRequest, sender *backend.StreamSender) error {
	return m.next.RunStream(ctx, req, sender)
}
