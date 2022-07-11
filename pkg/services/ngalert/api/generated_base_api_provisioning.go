/*Package api contains base API implementation of unified alerting
 *
 *Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 *
 *Do not manually edit these files, please find ngalert/api/swagger-codegen/ for commands on how to generate them.
 */
package api

import (
	"net/http"

	"github.com/grafana/grafana/pkg/api/response"
	"github.com/grafana/grafana/pkg/api/routing"
	"github.com/grafana/grafana/pkg/middleware"
	"github.com/grafana/grafana/pkg/models"
	apimodels "github.com/grafana/grafana/pkg/services/ngalert/api/tooling/definitions"
	"github.com/grafana/grafana/pkg/services/ngalert/metrics"
	"github.com/grafana/grafana/pkg/web"
)

type ProvisioningApiForkingService interface {
	RouteDeleteAlertRule(*models.ReqContext) response.Response
	RouteDeleteContactpoints(*models.ReqContext) response.Response
	RouteDeleteMuteTiming(*models.ReqContext) response.Response
	RouteDeleteTemplate(*models.ReqContext) response.Response
	RouteGetAlertRule(*models.ReqContext) response.Response
	RouteGetAlertRuleGroup(*models.ReqContext) response.Response
	RouteGetContactpoints(*models.ReqContext) response.Response
	RouteGetMuteTiming(*models.ReqContext) response.Response
	RouteGetMuteTimings(*models.ReqContext) response.Response
	RouteGetPolicyTree(*models.ReqContext) response.Response
	RouteGetTemplate(*models.ReqContext) response.Response
	RouteGetTemplates(*models.ReqContext) response.Response
	RoutePostAlertRule(*models.ReqContext) response.Response
	RoutePostContactpoints(*models.ReqContext) response.Response
	RoutePostMuteTiming(*models.ReqContext) response.Response
	RoutePutAlertRule(*models.ReqContext) response.Response
	RoutePutAlertRuleGroup(*models.ReqContext) response.Response
	RoutePutContactpoint(*models.ReqContext) response.Response
	RoutePutMuteTiming(*models.ReqContext) response.Response
	RoutePutPolicyTree(*models.ReqContext) response.Response
	RoutePutTemplate(*models.ReqContext) response.Response
	RouteResetPolicyTree(*models.ReqContext) response.Response
}

func (f *ForkedProvisioningApi) RouteDeleteAlertRule(ctx *models.ReqContext) response.Response {
	uIDParam := web.Params(ctx.Req)[":UID"]
	return f.forkRouteDeleteAlertRule(ctx, uIDParam)
}
func (f *ForkedProvisioningApi) RouteDeleteContactpoints(ctx *models.ReqContext) response.Response {
	uIDParam := web.Params(ctx.Req)[":UID"]
	return f.forkRouteDeleteContactpoints(ctx, uIDParam)
}
func (f *ForkedProvisioningApi) RouteDeleteMuteTiming(ctx *models.ReqContext) response.Response {
	nameParam := web.Params(ctx.Req)[":name"]
	return f.forkRouteDeleteMuteTiming(ctx, nameParam)
}
func (f *ForkedProvisioningApi) RouteDeleteTemplate(ctx *models.ReqContext) response.Response {
	nameParam := web.Params(ctx.Req)[":name"]
	return f.forkRouteDeleteTemplate(ctx, nameParam)
}
func (f *ForkedProvisioningApi) RouteGetAlertRule(ctx *models.ReqContext) response.Response {
	uIDParam := web.Params(ctx.Req)[":UID"]
	return f.forkRouteGetAlertRule(ctx, uIDParam)
}
func (f *ForkedProvisioningApi) RouteGetAlertRuleGroup(ctx *models.ReqContext) response.Response {
	folderUIDParam := web.Params(ctx.Req)[":FolderUID"]
	groupParam := web.Params(ctx.Req)[":Group"]
	return f.forkRouteGetAlertRuleGroup(ctx, folderUIDParam, groupParam)
}
func (f *ForkedProvisioningApi) RouteGetContactpoints(ctx *models.ReqContext) response.Response {
	return f.forkRouteGetContactpoints(ctx)
}
func (f *ForkedProvisioningApi) RouteGetMuteTiming(ctx *models.ReqContext) response.Response {
	nameParam := web.Params(ctx.Req)[":name"]
	return f.forkRouteGetMuteTiming(ctx, nameParam)
}
func (f *ForkedProvisioningApi) RouteGetMuteTimings(ctx *models.ReqContext) response.Response {
	return f.forkRouteGetMuteTimings(ctx)
}
func (f *ForkedProvisioningApi) RouteGetPolicyTree(ctx *models.ReqContext) response.Response {
	return f.forkRouteGetPolicyTree(ctx)
}
func (f *ForkedProvisioningApi) RouteGetTemplate(ctx *models.ReqContext) response.Response {
	nameParam := web.Params(ctx.Req)[":name"]
	return f.forkRouteGetTemplate(ctx, nameParam)
}
func (f *ForkedProvisioningApi) RouteGetTemplates(ctx *models.ReqContext) response.Response {
	return f.forkRouteGetTemplates(ctx)
}
func (f *ForkedProvisioningApi) RoutePostAlertRule(ctx *models.ReqContext) response.Response {
	conf := apimodels.AlertRule{}
	if err := web.Bind(ctx.Req, &conf); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	return f.forkRoutePostAlertRule(ctx, conf)
}
func (f *ForkedProvisioningApi) RoutePostContactpoints(ctx *models.ReqContext) response.Response {
	conf := apimodels.EmbeddedContactPoint{}
	if err := web.Bind(ctx.Req, &conf); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	return f.forkRoutePostContactpoints(ctx, conf)
}
func (f *ForkedProvisioningApi) RoutePostMuteTiming(ctx *models.ReqContext) response.Response {
	conf := apimodels.MuteTimeInterval{}
	if err := web.Bind(ctx.Req, &conf); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	return f.forkRoutePostMuteTiming(ctx, conf)
}
func (f *ForkedProvisioningApi) RoutePutAlertRule(ctx *models.ReqContext) response.Response {
	uIDParam := web.Params(ctx.Req)[":UID"]
	conf := apimodels.AlertRule{}
	if err := web.Bind(ctx.Req, &conf); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	return f.forkRoutePutAlertRule(ctx, conf, uIDParam)
}
func (f *ForkedProvisioningApi) RoutePutAlertRuleGroup(ctx *models.ReqContext) response.Response {
	folderUIDParam := web.Params(ctx.Req)[":FolderUID"]
	groupParam := web.Params(ctx.Req)[":Group"]
	conf := apimodels.AlertRuleGroupMetadata{}
	if err := web.Bind(ctx.Req, &conf); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	return f.forkRoutePutAlertRuleGroup(ctx, conf, folderUIDParam, groupParam)
}
func (f *ForkedProvisioningApi) RoutePutContactpoint(ctx *models.ReqContext) response.Response {
	uIDParam := web.Params(ctx.Req)[":UID"]
	conf := apimodels.EmbeddedContactPoint{}
	if err := web.Bind(ctx.Req, &conf); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	return f.forkRoutePutContactpoint(ctx, conf, uIDParam)
}
func (f *ForkedProvisioningApi) RoutePutMuteTiming(ctx *models.ReqContext) response.Response {
	nameParam := web.Params(ctx.Req)[":name"]
	conf := apimodels.MuteTimeInterval{}
	if err := web.Bind(ctx.Req, &conf); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	return f.forkRoutePutMuteTiming(ctx, conf, nameParam)
}
func (f *ForkedProvisioningApi) RoutePutPolicyTree(ctx *models.ReqContext) response.Response {
	conf := apimodels.Route{}
	if err := web.Bind(ctx.Req, &conf); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	return f.forkRoutePutPolicyTree(ctx, conf)
}
func (f *ForkedProvisioningApi) RoutePutTemplate(ctx *models.ReqContext) response.Response {
	nameParam := web.Params(ctx.Req)[":name"]
	conf := apimodels.MessageTemplateContent{}
	if err := web.Bind(ctx.Req, &conf); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	return f.forkRoutePutTemplate(ctx, conf, nameParam)
}
func (f *ForkedProvisioningApi) RouteResetPolicyTree(ctx *models.ReqContext) response.Response {
	return f.forkRouteResetPolicyTree(ctx)
}

func (api *API) RegisterProvisioningApiEndpoints(srv ProvisioningApiForkingService, m *metrics.API) {
	api.RouteRegister.Group("", func(group routing.RouteRegister) {
		group.Delete(
			toMacaronPath("/api/v1/provisioning/alert-rules/{UID}"),
			api.authorize(http.MethodDelete, "/api/v1/provisioning/alert-rules/{UID}"),
			metrics.Instrument(
				http.MethodDelete,
				"/api/v1/provisioning/alert-rules/{UID}",
				srv.RouteDeleteAlertRule,
				m,
			),
		)
		group.Delete(
			toMacaronPath("/api/v1/provisioning/contact-points/{UID}"),
			api.authorize(http.MethodDelete, "/api/v1/provisioning/contact-points/{UID}"),
			metrics.Instrument(
				http.MethodDelete,
				"/api/v1/provisioning/contact-points/{UID}",
				srv.RouteDeleteContactpoints,
				m,
			),
		)
		group.Delete(
			toMacaronPath("/api/v1/provisioning/mute-timings/{name}"),
			api.authorize(http.MethodDelete, "/api/v1/provisioning/mute-timings/{name}"),
			metrics.Instrument(
				http.MethodDelete,
				"/api/v1/provisioning/mute-timings/{name}",
				srv.RouteDeleteMuteTiming,
				m,
			),
		)
		group.Delete(
			toMacaronPath("/api/v1/provisioning/templates/{name}"),
			api.authorize(http.MethodDelete, "/api/v1/provisioning/templates/{name}"),
			metrics.Instrument(
				http.MethodDelete,
				"/api/v1/provisioning/templates/{name}",
				srv.RouteDeleteTemplate,
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/alert-rules/{UID}"),
			api.authorize(http.MethodGet, "/api/v1/provisioning/alert-rules/{UID}"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/alert-rules/{UID}",
				srv.RouteGetAlertRule,
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/folder/{FolderUID}/rule-groups/{Group}"),
			api.authorize(http.MethodGet, "/api/v1/provisioning/folder/{FolderUID}/rule-groups/{Group}"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/folder/{FolderUID}/rule-groups/{Group}",
				srv.RouteGetAlertRuleGroup,
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/contact-points"),
			api.authorize(http.MethodGet, "/api/v1/provisioning/contact-points"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/contact-points",
				srv.RouteGetContactpoints,
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/mute-timings/{name}"),
			api.authorize(http.MethodGet, "/api/v1/provisioning/mute-timings/{name}"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/mute-timings/{name}",
				srv.RouteGetMuteTiming,
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/mute-timings"),
			api.authorize(http.MethodGet, "/api/v1/provisioning/mute-timings"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/mute-timings",
				srv.RouteGetMuteTimings,
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/policies"),
			api.authorize(http.MethodGet, "/api/v1/provisioning/policies"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/policies",
				srv.RouteGetPolicyTree,
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/templates/{name}"),
			api.authorize(http.MethodGet, "/api/v1/provisioning/templates/{name}"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/templates/{name}",
				srv.RouteGetTemplate,
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/templates"),
			api.authorize(http.MethodGet, "/api/v1/provisioning/templates"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/templates",
				srv.RouteGetTemplates,
				m,
			),
		)
		group.Post(
			toMacaronPath("/api/v1/provisioning/alert-rules"),
			api.authorize(http.MethodPost, "/api/v1/provisioning/alert-rules"),
			metrics.Instrument(
				http.MethodPost,
				"/api/v1/provisioning/alert-rules",
				srv.RoutePostAlertRule,
				m,
			),
		)
		group.Post(
			toMacaronPath("/api/v1/provisioning/contact-points"),
			api.authorize(http.MethodPost, "/api/v1/provisioning/contact-points"),
			metrics.Instrument(
				http.MethodPost,
				"/api/v1/provisioning/contact-points",
				srv.RoutePostContactpoints,
				m,
			),
		)
		group.Post(
			toMacaronPath("/api/v1/provisioning/mute-timings"),
			api.authorize(http.MethodPost, "/api/v1/provisioning/mute-timings"),
			metrics.Instrument(
				http.MethodPost,
				"/api/v1/provisioning/mute-timings",
				srv.RoutePostMuteTiming,
				m,
			),
		)
		group.Put(
			toMacaronPath("/api/v1/provisioning/alert-rules/{UID}"),
			api.authorize(http.MethodPut, "/api/v1/provisioning/alert-rules/{UID}"),
			metrics.Instrument(
				http.MethodPut,
				"/api/v1/provisioning/alert-rules/{UID}",
				srv.RoutePutAlertRule,
				m,
			),
		)
		group.Put(
			toMacaronPath("/api/v1/provisioning/folder/{FolderUID}/rule-groups/{Group}"),
			api.authorize(http.MethodPut, "/api/v1/provisioning/folder/{FolderUID}/rule-groups/{Group}"),
			metrics.Instrument(
				http.MethodPut,
				"/api/v1/provisioning/folder/{FolderUID}/rule-groups/{Group}",
				srv.RoutePutAlertRuleGroup,
				m,
			),
		)
		group.Put(
			toMacaronPath("/api/v1/provisioning/contact-points/{UID}"),
			api.authorize(http.MethodPut, "/api/v1/provisioning/contact-points/{UID}"),
			metrics.Instrument(
				http.MethodPut,
				"/api/v1/provisioning/contact-points/{UID}",
				srv.RoutePutContactpoint,
				m,
			),
		)
		group.Put(
			toMacaronPath("/api/v1/provisioning/mute-timings/{name}"),
			api.authorize(http.MethodPut, "/api/v1/provisioning/mute-timings/{name}"),
			metrics.Instrument(
				http.MethodPut,
				"/api/v1/provisioning/mute-timings/{name}",
				srv.RoutePutMuteTiming,
				m,
			),
		)
		group.Put(
			toMacaronPath("/api/v1/provisioning/policies"),
			api.authorize(http.MethodPut, "/api/v1/provisioning/policies"),
			metrics.Instrument(
				http.MethodPut,
				"/api/v1/provisioning/policies",
				srv.RoutePutPolicyTree,
				m,
			),
		)
		group.Put(
			toMacaronPath("/api/v1/provisioning/templates/{name}"),
			api.authorize(http.MethodPut, "/api/v1/provisioning/templates/{name}"),
			metrics.Instrument(
				http.MethodPut,
				"/api/v1/provisioning/templates/{name}",
				srv.RoutePutTemplate,
				m,
			),
		)
		group.Delete(
			toMacaronPath("/api/v1/provisioning/policies"),
			api.authorize(http.MethodDelete, "/api/v1/provisioning/policies"),
			metrics.Instrument(
				http.MethodDelete,
				"/api/v1/provisioning/policies",
				srv.RouteResetPolicyTree,
				m,
			),
		)
	}, middleware.ReqSignedIn)
}
