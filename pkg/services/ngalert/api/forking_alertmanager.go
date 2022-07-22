package api

import (
	"fmt"

	"github.com/grafana/grafana/pkg/api/response"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/datasources"
	apimodels "github.com/grafana/grafana/pkg/services/ngalert/api/tooling/definitions"
)

type AlertmanagerApiHandler struct {
	AMSvc           *LotexAM
	GrafanaSvc      *AlertmanagerSrv
	DatasourceCache datasources.CacheService
}

// NewForkingAM implements a set of routes that proxy to various Alertmanager-compatible backends.
func NewForkingAM(datasourceCache datasources.CacheService, proxy *LotexAM, grafana *AlertmanagerSrv) *AlertmanagerApiHandler {
	return &AlertmanagerApiHandler{
		AMSvc:           proxy,
		GrafanaSvc:      grafana,
		DatasourceCache: datasourceCache,
	}
}

func (f *AlertmanagerApiHandler) getService(ctx *models.ReqContext) (*LotexAM, error) {
	t, err := backendTypeByUID(ctx, f.DatasourceCache)
	if err != nil {
		return nil, err
	}

	switch t {
	case apimodels.AlertmanagerBackend:
		return f.AMSvc, nil
	default:
		return nil, fmt.Errorf("unexpected backend type (%v)", t)
	}
}

func (f *AlertmanagerApiHandler) handleRouteGetAMStatus(ctx *models.ReqContext, dsUID string) response.Response {
	s, err := f.getService(ctx)
	if err != nil {
		return response.Error(400, err.Error(), nil)
	}

	return s.RouteGetAMStatus(ctx)
}

func (f *AlertmanagerApiHandler) handleRouteCreateSilence(ctx *models.ReqContext, body apimodels.PostableSilence, dsUID string) response.Response {
	s, err := f.getService(ctx)
	if err != nil {
		return ErrResp(400, err, "")
	}

	return s.RouteCreateSilence(ctx, body)
}

func (f *AlertmanagerApiHandler) handleRouteDeleteAlertingConfig(ctx *models.ReqContext, dsUID string) response.Response {
	s, err := f.getService(ctx)
	if err != nil {
		return ErrResp(400, err, "")
	}

	return s.RouteDeleteAlertingConfig(ctx)
}

func (f *AlertmanagerApiHandler) handleRouteDeleteSilence(ctx *models.ReqContext, silenceID string, dsUID string) response.Response {
	s, err := f.getService(ctx)
	if err != nil {
		return ErrResp(400, err, "")
	}

	return s.RouteDeleteSilence(ctx, silenceID)
}

func (f *AlertmanagerApiHandler) handleRouteGetAlertingConfig(ctx *models.ReqContext, dsUID string) response.Response {
	s, err := f.getService(ctx)
	if err != nil {
		return ErrResp(400, err, "")
	}

	return s.RouteGetAlertingConfig(ctx)
}

func (f *AlertmanagerApiHandler) handleRouteGetAMAlertGroups(ctx *models.ReqContext, dsUID string) response.Response {
	s, err := f.getService(ctx)
	if err != nil {
		return ErrResp(400, err, "")
	}

	return s.RouteGetAMAlertGroups(ctx)
}

func (f *AlertmanagerApiHandler) handleRouteGetAMAlerts(ctx *models.ReqContext, dsUID string) response.Response {
	s, err := f.getService(ctx)
	if err != nil {
		return ErrResp(400, err, "")
	}

	return s.RouteGetAMAlerts(ctx)
}

func (f *AlertmanagerApiHandler) handleRouteGetSilence(ctx *models.ReqContext, silenceID string, dsUID string) response.Response {
	s, err := f.getService(ctx)
	if err != nil {
		return ErrResp(400, err, "")
	}

	return s.RouteGetSilence(ctx, silenceID)
}

func (f *AlertmanagerApiHandler) handleRouteGetSilences(ctx *models.ReqContext, dsUID string) response.Response {
	s, err := f.getService(ctx)
	if err != nil {
		return ErrResp(400, err, "")
	}

	return s.RouteGetSilences(ctx)
}

func (f *AlertmanagerApiHandler) handleRoutePostAlertingConfig(ctx *models.ReqContext, body apimodels.PostableUserConfig, dsUID string) response.Response {
	s, err := f.getService(ctx)
	if err != nil {
		return ErrResp(400, err, "")
	}

	b, err := backendTypeByUID(ctx, f.DatasourceCache)
	if err != nil {
		return ErrResp(400, err, "")
	}

	if err := body.AlertmanagerConfig.ReceiverType().MatchesBackend(b); err != nil {
		return ErrResp(400, err, "bad match")
	}

	return s.RoutePostAlertingConfig(ctx, body)
}

func (f *AlertmanagerApiHandler) handleRoutePostAMAlerts(ctx *models.ReqContext, body apimodels.PostableAlerts, dsUID string) response.Response {
	s, err := f.getService(ctx)
	if err != nil {
		return ErrResp(400, err, "")
	}

	return s.RoutePostAMAlerts(ctx, body)
}

func (f *AlertmanagerApiHandler) handleRoutePostTestReceivers(ctx *models.ReqContext, body apimodels.TestReceiversConfigBodyParams, dsUID string) response.Response {
	s, err := f.getService(ctx)
	if err != nil {
		return ErrResp(400, err, "")
	}

	return s.RoutePostTestReceivers(ctx, body)
}

func (f *AlertmanagerApiHandler) handleRouteDeleteGrafanaSilence(ctx *models.ReqContext, id string) response.Response {
	return f.GrafanaSvc.RouteDeleteSilence(ctx, id)
}

func (f *AlertmanagerApiHandler) handleRouteDeleteGrafanaAlertingConfig(ctx *models.ReqContext) response.Response {
	return f.GrafanaSvc.RouteDeleteAlertingConfig(ctx)
}

func (f *AlertmanagerApiHandler) handleRouteCreateGrafanaSilence(ctx *models.ReqContext, body apimodels.PostableSilence) response.Response {
	return f.GrafanaSvc.RouteCreateSilence(ctx, body)
}

func (f *AlertmanagerApiHandler) handleRouteGetGrafanaAMStatus(ctx *models.ReqContext) response.Response {
	return f.GrafanaSvc.RouteGetAMStatus(ctx)
}

func (f *AlertmanagerApiHandler) handleRouteGetGrafanaAMAlerts(ctx *models.ReqContext) response.Response {
	return f.GrafanaSvc.RouteGetAMAlerts(ctx)
}

func (f *AlertmanagerApiHandler) handleRouteGetGrafanaAMAlertGroups(ctx *models.ReqContext) response.Response {
	return f.GrafanaSvc.RouteGetAMAlertGroups(ctx)
}

func (f *AlertmanagerApiHandler) handleRouteGetGrafanaAlertingConfig(ctx *models.ReqContext) response.Response {
	return f.GrafanaSvc.RouteGetAlertingConfig(ctx)
}

func (f *AlertmanagerApiHandler) handleRouteGetGrafanaSilence(ctx *models.ReqContext, id string) response.Response {
	return f.GrafanaSvc.RouteGetSilence(ctx, id)
}

func (f *AlertmanagerApiHandler) handleRouteGetGrafanaSilences(ctx *models.ReqContext) response.Response {
	return f.GrafanaSvc.RouteGetSilences(ctx)
}

func (f *AlertmanagerApiHandler) handleRoutePostGrafanaAMAlerts(ctx *models.ReqContext, conf apimodels.PostableAlerts) response.Response {
	return f.GrafanaSvc.RoutePostAMAlerts(ctx, conf)
}

func (f *AlertmanagerApiHandler) handleRoutePostGrafanaAlertingConfig(ctx *models.ReqContext, conf apimodels.PostableUserConfig) response.Response {
	return f.GrafanaSvc.RoutePostAlertingConfig(ctx, conf)
}

func (f *AlertmanagerApiHandler) handleRoutePostTestGrafanaReceivers(ctx *models.ReqContext, conf apimodels.TestReceiversConfigBodyParams) response.Response {
	return f.GrafanaSvc.RoutePostTestReceivers(ctx, conf)
}
