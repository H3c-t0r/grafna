package prometheus

import (
	"context"
	"fmt"

	"github.com/grafana/grafana-azure-sdk-go/azsettings"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	sdkhttpclient "github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"

	"github.com/grafana/grafana/pkg/cmd/grafana-cli/logger"
	prometheuslibrary "github.com/grafana/grafana/pkg/prometheus-library"
	"github.com/grafana/grafana/pkg/tsdb/prometheus/azureauth"
)

type Service struct {
	lib *prometheuslibrary.Service
}

func ProvideService(httpClientProvider *sdkhttpclient.Provider) *Service {
	plog := backend.NewLoggerWith("logger", "tsdb.prometheus")
	plog.Debug("Initializing")
	return &Service{
		lib: prometheuslibrary.NewService(httpClientProvider, plog, extendClientOpts),
	}
}

func (s *Service) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return s.lib.QueryData(ctx, req)
}

func (s *Service) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	return s.lib.CallResource(ctx, req, sender)
}

func (s *Service) GetBuildInfo(ctx context.Context, req prometheuslibrary.BuildInfoRequest) (*prometheuslibrary.BuildInfoResponse, error) {
	return s.GetBuildInfo(ctx, req)
}

func (s *Service) GetHeuristics(ctx context.Context, req prometheuslibrary.HeuristicsRequest) (*prometheuslibrary.Heuristics, error) {
	return s.lib.GetHeuristics(ctx, req)
}

func (s *Service) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult,
	error) {
	return s.lib.CheckHealth(ctx, req)
}

func extendClientOpts(ctx context.Context, settings backend.DataSourceInstanceSettings, clientOpts *sdkhttpclient.Options) (*sdkhttpclient.Options, error) {
	// Set SigV4 service namespace
	if clientOpts.SigV4 != nil {
		clientOpts.SigV4.Service = "aps"
	}

	azureSettings, err := azsettings.ReadSettings(ctx)
	if err != nil {
		logger.Error("failed to read Azure settings from Grafana", "error", err.Error())
		return nil, fmt.Errorf("failed to read Azure settings from Grafana: %v", err)
	}

	// Set Azure authentication
	if azureSettings.AzureAuthEnabled {
		err = azureauth.ConfigureAzureAuthentication(settings, azureSettings, clientOpts)
		if err != nil {
			return nil, fmt.Errorf("error configuring Azure auth: %v", err)
		}
	}

	return clientOpts, nil
}
