package azuremonitor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana/pkg/api/pluginproxy"
	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/plugins"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/tsdb"
	"github.com/grafana/grafana/pkg/util/errutil"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context/ctxhttp"
)

// ApplicationInsightsDatasource calls the application insights query API's
type ApplicationInsightsDatasource struct {
	httpClient *http.Client
	dsInfo     *models.DataSource
}

type ApplicationInsightsQuery struct {
	RefID string

	IsRaw bool

	// Text based raw query options
	ApiURL string
	Params url.Values
	Alias  string
	Target string

	metricName  string
	dimensions  []string
	aggregation string
}

func (e *ApplicationInsightsDatasource) executeTimeSeriesQuery(ctx context.Context, originalQueries []*tsdb.Query, timeRange *tsdb.TimeRange) (*tsdb.Response, error) {
	result := &tsdb.Response{
		Results: map[string]*tsdb.QueryResult{},
	}

	queries, err := e.buildQueries(originalQueries, timeRange)
	if err != nil {
		return nil, err
	}

	for _, query := range queries {
		queryRes, err := e.executeQuery(ctx, query)
		if err != nil {
			return nil, err
		}
		result.Results[query.RefID] = queryRes
	}

	return result, nil
}

func (e *ApplicationInsightsDatasource) buildQueries(queries []*tsdb.Query, timeRange *tsdb.TimeRange) ([]*ApplicationInsightsQuery, error) {
	applicationInsightsQueries := []*ApplicationInsightsQuery{}
	startTime, err := timeRange.ParseFrom()
	if err != nil {
		return nil, err
	}

	endTime, err := timeRange.ParseTo()
	if err != nil {
		return nil, err
	}

	for _, query := range queries {
		queryBytes, err := query.Model.Encode()
		if err != nil {
			return nil, fmt.Errorf("failed to re-encode the Azure Application Insights query into JSON: %w", err)
		}
		queryJSONModel := insightsJSONQuery{}
		err = json.Unmarshal(queryBytes, &queryJSONModel)
		if err != nil {
			return nil, fmt.Errorf("failed to decode the Azure Application Insights query object from JSON: %w", err)
		}

		insightsJSONModel := queryJSONModel.AppInsights
		azlog.Debug("Application Insights", "target", insightsJSONModel)

		azureURL := fmt.Sprintf("metrics/%s", insightsJSONModel.MetricName)
		timeGrain := insightsJSONModel.TimeGrain
		timeGrains := insightsJSONModel.AllowedTimeGrainsMs
		if timeGrain == "auto" {
			timeGrain, err = setAutoTimeGrain(query.IntervalMs, timeGrains)
			if err != nil {
				return nil, err
			}
		}

		params := url.Values{}
		params.Add("timespan", fmt.Sprintf("%v/%v", startTime.UTC().Format(time.RFC3339), endTime.UTC().Format(time.RFC3339)))
		if timeGrain != "none" {
			params.Add("interval", timeGrain)
		}
		params.Add("aggregation", insightsJSONModel.Aggregation)

		dims := []string{}

		dimension := strings.TrimSpace(insightsJSONModel.Dimension)
		// Azure Monitor combines this and the following logic such that if dimensionFilter, must also Dimension, should that be done here as well?
		if dimension != "" && !strings.EqualFold(dimension, "none") {
			params.Add("segment", dimension)
			dims = append(dims, dimension)
		}

		dimensionFilter := strings.TrimSpace(insightsJSONModel.DimensionFilter)
		if dimensionFilter != "" {
			params.Add("filter", dimensionFilter)
		}

		applicationInsightsQueries = append(applicationInsightsQueries, &ApplicationInsightsQuery{
			RefID:       query.RefId,
			IsRaw:       false,
			ApiURL:      azureURL,
			Params:      params,
			Alias:       insightsJSONModel.Alias,
			Target:      params.Encode(),
			metricName:  insightsJSONModel.MetricName,
			aggregation: insightsJSONModel.Aggregation,
			dimensions:  dims,
		})

	}

	return applicationInsightsQueries, nil
}

func (e *ApplicationInsightsDatasource) executeQuery(ctx context.Context, query *ApplicationInsightsQuery) (*tsdb.QueryResult, error) {
	queryResult := &tsdb.QueryResult{Meta: simplejson.New(), RefId: query.RefID}

	req, err := e.createRequest(ctx, e.dsInfo)
	if err != nil {
		queryResult.Error = err
		return queryResult, nil
	}

	req.URL.Path = path.Join(req.URL.Path, query.ApiURL)
	req.URL.RawQuery = query.Params.Encode()

	span, ctx := opentracing.StartSpanFromContext(ctx, "application insights query")
	span.SetTag("target", query.Target)
	span.SetTag("datasource_id", e.dsInfo.Id)
	span.SetTag("org_id", e.dsInfo.OrgId)

	defer span.Finish()

	err = opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))

	if err != nil {
		azlog.Warn("failed to inject global tracer")
	}

	azlog.Debug("ApplicationInsights", "Request URL", req.URL.String())
	res, err := ctxhttp.Do(ctx, e.httpClient, req)
	if err != nil {
		queryResult.Error = err
		return queryResult, nil
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	if res.StatusCode/100 != 2 {
		azlog.Debug("Request failed", "status", res.Status, "body", string(body))
		return nil, fmt.Errorf("Request failed status: %v", res.Status)
	}

	mr := MetricsResult{}
	err = json.Unmarshal(body, &mr)
	if err != nil {
		return nil, err
	}

	frame, err := InsightsMetricsResultToFrame(mr, query.metricName, query.aggregation, query.dimensions)
	if err != nil {
		queryResult.Error = err
		return queryResult, nil
	}
	queryResult.Dataframes = tsdb.NewDecodedDataFrames(data.Frames{frame})
	return queryResult, nil
}

func (e *ApplicationInsightsDatasource) createRequest(ctx context.Context, dsInfo *models.DataSource) (*http.Request, error) {
	// find plugin
	plugin, ok := plugins.DataSources[dsInfo.Type]
	if !ok {
		return nil, errors.New("Unable to find datasource plugin Azure Application Insights")
	}

	cloudName := dsInfo.JsonData.Get("cloudName").MustString("azuremonitor")
	appInsightsRoute, pluginRouteName, err := e.getPluginRoute(plugin, cloudName)
	if err != nil {
		return nil, err
	}

	appInsightsAppID := dsInfo.JsonData.Get("appInsightsAppId").MustString()
	proxyPass := fmt.Sprintf("%s/v1/apps/%s", pluginRouteName, appInsightsAppID)

	u, err := url.Parse(dsInfo.Url)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, fmt.Sprintf("/v1/apps/%s", appInsightsAppID))

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		azlog.Debug("Failed to create request", "error", err)
		return nil, errutil.Wrap("Failed to create request", err)
	}

	req.Header.Set("User-Agent", fmt.Sprintf("Grafana/%s", setting.BuildVersion))

	pluginproxy.ApplyRoute(ctx, req, proxyPass, appInsightsRoute, dsInfo)

	return req, nil
}

func (e *ApplicationInsightsDatasource) getPluginRoute(plugin *plugins.DataSourcePlugin, cloudName string) (*plugins.AppPluginRoute, string, error) {
	pluginRouteName := "appinsights"

	if cloudName == "chinaazuremonitor" {
		pluginRouteName = "chinaappinsights"
	}

	var pluginRoute *plugins.AppPluginRoute

	for _, route := range plugin.Routes {
		if route.Path == pluginRouteName {
			pluginRoute = route
			break
		}
	}

	return pluginRoute, pluginRouteName, nil
}
