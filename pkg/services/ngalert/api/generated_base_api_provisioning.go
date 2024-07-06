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
	"github.com/grafana/grafana/pkg/middleware/requestmeta"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
	apimodels "github.com/grafana/grafana/pkg/services/ngalert/api/tooling/definitions"
	"github.com/grafana/grafana/pkg/services/ngalert/metrics"
	"github.com/grafana/grafana/pkg/web"
)

type ProvisioningApi interface {
	RouteDeleteAlertRule(*contextmodel.ReqContext) response.Response
	RouteDeleteAlertRuleGroup(*contextmodel.ReqContext) response.Response
	RouteDeleteContactpoints(*contextmodel.ReqContext) response.Response
	RouteDeleteMuteTiming(*contextmodel.ReqContext) response.Response
	RouteDeleteTemplate(*contextmodel.ReqContext) response.Response
	RouteExportMuteTiming(*contextmodel.ReqContext) response.Response
	RouteExportMuteTimings(*contextmodel.ReqContext) response.Response
	RouteGetAlertRule(*contextmodel.ReqContext) response.Response
	RouteGetAlertRuleExport(*contextmodel.ReqContext) response.Response
	RouteGetAlertRuleGroup(*contextmodel.ReqContext) response.Response
	RouteGetAlertRuleGroupExport(*contextmodel.ReqContext) response.Response
	RouteGetAlertRules(*contextmodel.ReqContext) response.Response
	RouteGetAlertRulesExport(*contextmodel.ReqContext) response.Response
	RouteGetContactpoints(*contextmodel.ReqContext) response.Response
	RouteGetContactpointsExport(*contextmodel.ReqContext) response.Response
	RouteGetMuteTiming(*contextmodel.ReqContext) response.Response
	RouteGetMuteTimings(*contextmodel.ReqContext) response.Response
	RouteGetPolicyTree(*contextmodel.ReqContext) response.Response
	RouteGetPolicyTreeExport(*contextmodel.ReqContext) response.Response
	RouteGetTemplate(*contextmodel.ReqContext) response.Response
	RouteGetTemplates(*contextmodel.ReqContext) response.Response
	RoutePostAlertRule(*contextmodel.ReqContext) response.Response
	RoutePostAlertRuleGroup(*contextmodel.ReqContext) response.Response
	RoutePostContactpoints(*contextmodel.ReqContext) response.Response
	RoutePostMuteTiming(*contextmodel.ReqContext) response.Response
	RoutePutAlertRule(*contextmodel.ReqContext) response.Response
	RoutePutAlertRuleGroup(*contextmodel.ReqContext) response.Response
	RoutePutContactpoint(*contextmodel.ReqContext) response.Response
	RoutePutMuteTiming(*contextmodel.ReqContext) response.Response
	RoutePutPolicyTree(*contextmodel.ReqContext) response.Response
	RoutePutTemplate(*contextmodel.ReqContext) response.Response
	RouteResetPolicyTree(*contextmodel.ReqContext) response.Response
}

func (f *ProvisioningApiHandler) RouteDeleteAlertRule(ctx *contextmodel.ReqContext) response.Response {
	// Parse Path Parameters
	uIDParam := web.Params(ctx.Req)[":UID"]
	return f.handleRouteDeleteAlertRule(ctx, uIDParam)
}
func (f *ProvisioningApiHandler) RouteDeleteAlertRuleGroup(ctx *contextmodel.ReqContext) response.Response {
	// Parse Path Parameters
	folderUIDParam := web.Params(ctx.Req)[":FolderUID"]
	groupParam := web.Params(ctx.Req)[":Group"]
	return f.handleRouteDeleteAlertRuleGroup(ctx, folderUIDParam, groupParam)
}
func (f *ProvisioningApiHandler) RouteDeleteContactpoints(ctx *contextmodel.ReqContext) response.Response {
	// Parse Path Parameters
	uIDParam := web.Params(ctx.Req)[":UID"]
	return f.handleRouteDeleteContactpoints(ctx, uIDParam)
}
func (f *ProvisioningApiHandler) RouteDeleteMuteTiming(ctx *contextmodel.ReqContext) response.Response {
	// Parse Path Parameters
	nameParam := web.Params(ctx.Req)[":name"]
	return f.handleRouteDeleteMuteTiming(ctx, nameParam)
}
func (f *ProvisioningApiHandler) RouteDeleteTemplate(ctx *contextmodel.ReqContext) response.Response {
	// Parse Path Parameters
	nameParam := web.Params(ctx.Req)[":name"]
	return f.handleRouteDeleteTemplate(ctx, nameParam)
}
func (f *ProvisioningApiHandler) RouteExportMuteTiming(ctx *contextmodel.ReqContext) response.Response {
	// Parse Path Parameters
	nameParam := web.Params(ctx.Req)[":name"]
	return f.handleRouteExportMuteTiming(ctx, nameParam)
}
func (f *ProvisioningApiHandler) RouteExportMuteTimings(ctx *contextmodel.ReqContext) response.Response {
	return f.handleRouteExportMuteTimings(ctx)
}
func (f *ProvisioningApiHandler) RouteGetAlertRule(ctx *contextmodel.ReqContext) response.Response {
	// Parse Path Parameters
	uIDParam := web.Params(ctx.Req)[":UID"]
	return f.handleRouteGetAlertRule(ctx, uIDParam)
}
func (f *ProvisioningApiHandler) RouteGetAlertRuleExport(ctx *contextmodel.ReqContext) response.Response {
	// Parse Path Parameters
	uIDParam := web.Params(ctx.Req)[":UID"]
	return f.handleRouteGetAlertRuleExport(ctx, uIDParam)
}
func (f *ProvisioningApiHandler) RouteGetAlertRuleGroup(ctx *contextmodel.ReqContext) response.Response {
	// Parse Path Parameters
	folderUIDParam := web.Params(ctx.Req)[":FolderUID"]
	groupParam := web.Params(ctx.Req)[":Group"]
	return f.handleRouteGetAlertRuleGroup(ctx, folderUIDParam, groupParam)
}
func (f *ProvisioningApiHandler) RouteGetAlertRuleGroupExport(ctx *contextmodel.ReqContext) response.Response {
	// Parse Path Parameters
	folderUIDParam := web.Params(ctx.Req)[":FolderUID"]
	groupParam := web.Params(ctx.Req)[":Group"]
	return f.handleRouteGetAlertRuleGroupExport(ctx, folderUIDParam, groupParam)
}
func (f *ProvisioningApiHandler) RouteGetAlertRules(ctx *contextmodel.ReqContext) response.Response {
	return f.handleRouteGetAlertRules(ctx)
}
func (f *ProvisioningApiHandler) RouteGetAlertRulesExport(ctx *contextmodel.ReqContext) response.Response {
	return f.handleRouteGetAlertRulesExport(ctx)
}
func (f *ProvisioningApiHandler) RouteGetContactpoints(ctx *contextmodel.ReqContext) response.Response {
	return f.handleRouteGetContactpoints(ctx)
}
func (f *ProvisioningApiHandler) RouteGetContactpointsExport(ctx *contextmodel.ReqContext) response.Response {
	return f.handleRouteGetContactpointsExport(ctx)
}
func (f *ProvisioningApiHandler) RouteGetMuteTiming(ctx *contextmodel.ReqContext) response.Response {
	// Parse Path Parameters
	nameParam := web.Params(ctx.Req)[":name"]
	return f.handleRouteGetMuteTiming(ctx, nameParam)
}
func (f *ProvisioningApiHandler) RouteGetMuteTimings(ctx *contextmodel.ReqContext) response.Response {
	return f.handleRouteGetMuteTimings(ctx)
}
func (f *ProvisioningApiHandler) RouteGetPolicyTree(ctx *contextmodel.ReqContext) response.Response {
	return f.handleRouteGetPolicyTree(ctx)
}
func (f *ProvisioningApiHandler) RouteGetPolicyTreeExport(ctx *contextmodel.ReqContext) response.Response {
	return f.handleRouteGetPolicyTreeExport(ctx)
}
func (f *ProvisioningApiHandler) RouteGetTemplate(ctx *contextmodel.ReqContext) response.Response {
	// Parse Path Parameters
	nameParam := web.Params(ctx.Req)[":name"]
	return f.handleRouteGetTemplate(ctx, nameParam)
}
func (f *ProvisioningApiHandler) RouteGetTemplates(ctx *contextmodel.ReqContext) response.Response {
	return f.handleRouteGetTemplates(ctx)
}
func (f *ProvisioningApiHandler) RoutePostAlertRule(ctx *contextmodel.ReqContext) response.Response {
	// Parse Request Body
	conf := apimodels.ProvisionedAlertRule{}
	if err := web.Bind(ctx.Req, &conf); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	return f.handleRoutePostAlertRule(ctx, conf)
}
func (f *ProvisioningApiHandler) RoutePostAlertRuleGroup(ctx *contextmodel.ReqContext) response.Response {
	// Parse Path Parameters
	folderUIDParam := web.Params(ctx.Req)[":FolderUID"]
	// Parse Request Body
	conf := apimodels.AlertRuleGroup{}
	if err := web.Bind(ctx.Req, &conf); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	return f.handleRoutePostAlertRuleGroup(ctx, conf, folderUIDParam)
}
func (f *ProvisioningApiHandler) RoutePostContactpoints(ctx *contextmodel.ReqContext) response.Response {
	// Parse Request Body
	conf := apimodels.EmbeddedContactPoint{}
	if err := web.Bind(ctx.Req, &conf); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	return f.handleRoutePostContactpoints(ctx, conf)
}
func (f *ProvisioningApiHandler) RoutePostMuteTiming(ctx *contextmodel.ReqContext) response.Response {
	// Parse Request Body
	conf := apimodels.MuteTimeInterval{}
	if err := web.Bind(ctx.Req, &conf); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	return f.handleRoutePostMuteTiming(ctx, conf)
}
func (f *ProvisioningApiHandler) RoutePutAlertRule(ctx *contextmodel.ReqContext) response.Response {
	// Parse Path Parameters
	uIDParam := web.Params(ctx.Req)[":UID"]
	// Parse Request Body
	conf := apimodels.ProvisionedAlertRule{}
	if err := web.Bind(ctx.Req, &conf); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	return f.handleRoutePutAlertRule(ctx, conf, uIDParam)
}
func (f *ProvisioningApiHandler) RoutePutAlertRuleGroup(ctx *contextmodel.ReqContext) response.Response {
	// Parse Path Parameters
	folderUIDParam := web.Params(ctx.Req)[":FolderUID"]
	groupParam := web.Params(ctx.Req)[":Group"]
	// Parse Request Body
	conf := apimodels.AlertRuleGroup{}
	if err := web.Bind(ctx.Req, &conf); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	return f.handleRoutePutAlertRuleGroup(ctx, conf, folderUIDParam, groupParam)
}
func (f *ProvisioningApiHandler) RoutePutContactpoint(ctx *contextmodel.ReqContext) response.Response {
	// Parse Path Parameters
	uIDParam := web.Params(ctx.Req)[":UID"]
	// Parse Request Body
	conf := apimodels.EmbeddedContactPoint{}
	if err := web.Bind(ctx.Req, &conf); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	return f.handleRoutePutContactpoint(ctx, conf, uIDParam)
}
func (f *ProvisioningApiHandler) RoutePutMuteTiming(ctx *contextmodel.ReqContext) response.Response {
	// Parse Path Parameters
	nameParam := web.Params(ctx.Req)[":name"]
	// Parse Request Body
	conf := apimodels.MuteTimeInterval{}
	if err := web.Bind(ctx.Req, &conf); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	return f.handleRoutePutMuteTiming(ctx, conf, nameParam)
}
func (f *ProvisioningApiHandler) RoutePutPolicyTree(ctx *contextmodel.ReqContext) response.Response {
	// Parse Request Body
	conf := apimodels.Route{}
	if err := web.Bind(ctx.Req, &conf); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	return f.handleRoutePutPolicyTree(ctx, conf)
}
func (f *ProvisioningApiHandler) RoutePutTemplate(ctx *contextmodel.ReqContext) response.Response {
	// Parse Path Parameters
	nameParam := web.Params(ctx.Req)[":name"]
	// Parse Request Body
	conf := apimodels.NotificationTemplateContent{}
	if err := web.Bind(ctx.Req, &conf); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	return f.handleRoutePutTemplate(ctx, conf, nameParam)
}
func (f *ProvisioningApiHandler) RouteResetPolicyTree(ctx *contextmodel.ReqContext) response.Response {
	return f.handleRouteResetPolicyTree(ctx)
}

func (api *API) RegisterProvisioningApiEndpoints(srv ProvisioningApi, m *metrics.API) {
	api.RouteRegister.Group("", func(group routing.RouteRegister) {
		group.Delete(
			toMacaronPath("/api/v1/provisioning/alert-rules/{UID}"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodDelete, "/api/v1/provisioning/alert-rules/{UID}"),
			metrics.Instrument(
				http.MethodDelete,
				"/api/v1/provisioning/alert-rules/{UID}",
				api.Hooks.Wrap(srv.RouteDeleteAlertRule),
				m,
			),
		)
		group.Delete(
			toMacaronPath("/api/v1/provisioning/folder/{FolderUID}/rule-groups/{Group}"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodDelete, "/api/v1/provisioning/folder/{FolderUID}/rule-groups/{Group}"),
			metrics.Instrument(
				http.MethodDelete,
				"/api/v1/provisioning/folder/{FolderUID}/rule-groups/{Group}",
				api.Hooks.Wrap(srv.RouteDeleteAlertRuleGroup),
				m,
			),
		)
		group.Delete(
			toMacaronPath("/api/v1/provisioning/contact-points/{UID}"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodDelete, "/api/v1/provisioning/contact-points/{UID}"),
			metrics.Instrument(
				http.MethodDelete,
				"/api/v1/provisioning/contact-points/{UID}",
				api.Hooks.Wrap(srv.RouteDeleteContactpoints),
				m,
			),
		)
		group.Delete(
			toMacaronPath("/api/v1/provisioning/mute-timings/{name}"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodDelete, "/api/v1/provisioning/mute-timings/{name}"),
			metrics.Instrument(
				http.MethodDelete,
				"/api/v1/provisioning/mute-timings/{name}",
				api.Hooks.Wrap(srv.RouteDeleteMuteTiming),
				m,
			),
		)
		group.Delete(
			toMacaronPath("/api/v1/provisioning/templates/{name}"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodDelete, "/api/v1/provisioning/templates/{name}"),
			metrics.Instrument(
				http.MethodDelete,
				"/api/v1/provisioning/templates/{name}",
				api.Hooks.Wrap(srv.RouteDeleteTemplate),
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/mute-timings/{name}/export"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodGet, "/api/v1/provisioning/mute-timings/{name}/export"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/mute-timings/{name}/export",
				api.Hooks.Wrap(srv.RouteExportMuteTiming),
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/mute-timings/export"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodGet, "/api/v1/provisioning/mute-timings/export"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/mute-timings/export",
				api.Hooks.Wrap(srv.RouteExportMuteTimings),
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/alert-rules/{UID}"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodGet, "/api/v1/provisioning/alert-rules/{UID}"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/alert-rules/{UID}",
				api.Hooks.Wrap(srv.RouteGetAlertRule),
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/alert-rules/{UID}/export"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodGet, "/api/v1/provisioning/alert-rules/{UID}/export"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/alert-rules/{UID}/export",
				api.Hooks.Wrap(srv.RouteGetAlertRuleExport),
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/folder/{FolderUID}/rule-groups/{Group}"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodGet, "/api/v1/provisioning/folder/{FolderUID}/rule-groups/{Group}"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/folder/{FolderUID}/rule-groups/{Group}",
				api.Hooks.Wrap(srv.RouteGetAlertRuleGroup),
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/folder/{FolderUID}/rule-groups/{Group}/export"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodGet, "/api/v1/provisioning/folder/{FolderUID}/rule-groups/{Group}/export"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/folder/{FolderUID}/rule-groups/{Group}/export",
				api.Hooks.Wrap(srv.RouteGetAlertRuleGroupExport),
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/alert-rules"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodGet, "/api/v1/provisioning/alert-rules"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/alert-rules",
				api.Hooks.Wrap(srv.RouteGetAlertRules),
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/alert-rules/export"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodGet, "/api/v1/provisioning/alert-rules/export"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/alert-rules/export",
				api.Hooks.Wrap(srv.RouteGetAlertRulesExport),
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/contact-points"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodGet, "/api/v1/provisioning/contact-points"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/contact-points",
				api.Hooks.Wrap(srv.RouteGetContactpoints),
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/contact-points/export"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodGet, "/api/v1/provisioning/contact-points/export"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/contact-points/export",
				api.Hooks.Wrap(srv.RouteGetContactpointsExport),
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/mute-timings/{name}"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodGet, "/api/v1/provisioning/mute-timings/{name}"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/mute-timings/{name}",
				api.Hooks.Wrap(srv.RouteGetMuteTiming),
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/mute-timings"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodGet, "/api/v1/provisioning/mute-timings"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/mute-timings",
				api.Hooks.Wrap(srv.RouteGetMuteTimings),
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/policies"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodGet, "/api/v1/provisioning/policies"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/policies",
				api.Hooks.Wrap(srv.RouteGetPolicyTree),
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/policies/export"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodGet, "/api/v1/provisioning/policies/export"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/policies/export",
				api.Hooks.Wrap(srv.RouteGetPolicyTreeExport),
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/templates/{name}"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodGet, "/api/v1/provisioning/templates/{name}"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/templates/{name}",
				api.Hooks.Wrap(srv.RouteGetTemplate),
				m,
			),
		)
		group.Get(
			toMacaronPath("/api/v1/provisioning/templates"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodGet, "/api/v1/provisioning/templates"),
			metrics.Instrument(
				http.MethodGet,
				"/api/v1/provisioning/templates",
				api.Hooks.Wrap(srv.RouteGetTemplates),
				m,
			),
		)
		group.Post(
			toMacaronPath("/api/v1/provisioning/alert-rules"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodPost, "/api/v1/provisioning/alert-rules"),
			metrics.Instrument(
				http.MethodPost,
				"/api/v1/provisioning/alert-rules",
				api.Hooks.Wrap(srv.RoutePostAlertRule),
				m,
			),
		)
		group.Post(
			toMacaronPath("/api/v1/provisioning/folder/{FolderUID}/rule-groups"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodPost, "/api/v1/provisioning/folder/{FolderUID}/rule-groups"),
			metrics.Instrument(
				http.MethodPost,
				"/api/v1/provisioning/folder/{FolderUID}/rule-groups",
				api.Hooks.Wrap(srv.RoutePostAlertRuleGroup),
				m,
			),
		)
		group.Post(
			toMacaronPath("/api/v1/provisioning/contact-points"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodPost, "/api/v1/provisioning/contact-points"),
			metrics.Instrument(
				http.MethodPost,
				"/api/v1/provisioning/contact-points",
				api.Hooks.Wrap(srv.RoutePostContactpoints),
				m,
			),
		)
		group.Post(
			toMacaronPath("/api/v1/provisioning/mute-timings"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodPost, "/api/v1/provisioning/mute-timings"),
			metrics.Instrument(
				http.MethodPost,
				"/api/v1/provisioning/mute-timings",
				api.Hooks.Wrap(srv.RoutePostMuteTiming),
				m,
			),
		)
		group.Put(
			toMacaronPath("/api/v1/provisioning/alert-rules/{UID}"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodPut, "/api/v1/provisioning/alert-rules/{UID}"),
			metrics.Instrument(
				http.MethodPut,
				"/api/v1/provisioning/alert-rules/{UID}",
				api.Hooks.Wrap(srv.RoutePutAlertRule),
				m,
			),
		)
		group.Put(
			toMacaronPath("/api/v1/provisioning/folder/{FolderUID}/rule-groups/{Group}"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodPut, "/api/v1/provisioning/folder/{FolderUID}/rule-groups/{Group}"),
			metrics.Instrument(
				http.MethodPut,
				"/api/v1/provisioning/folder/{FolderUID}/rule-groups/{Group}",
				api.Hooks.Wrap(srv.RoutePutAlertRuleGroup),
				m,
			),
		)
		group.Put(
			toMacaronPath("/api/v1/provisioning/contact-points/{UID}"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodPut, "/api/v1/provisioning/contact-points/{UID}"),
			metrics.Instrument(
				http.MethodPut,
				"/api/v1/provisioning/contact-points/{UID}",
				api.Hooks.Wrap(srv.RoutePutContactpoint),
				m,
			),
		)
		group.Put(
			toMacaronPath("/api/v1/provisioning/mute-timings/{name}"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodPut, "/api/v1/provisioning/mute-timings/{name}"),
			metrics.Instrument(
				http.MethodPut,
				"/api/v1/provisioning/mute-timings/{name}",
				api.Hooks.Wrap(srv.RoutePutMuteTiming),
				m,
			),
		)
		group.Put(
			toMacaronPath("/api/v1/provisioning/policies"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodPut, "/api/v1/provisioning/policies"),
			metrics.Instrument(
				http.MethodPut,
				"/api/v1/provisioning/policies",
				api.Hooks.Wrap(srv.RoutePutPolicyTree),
				m,
			),
		)
		group.Put(
			toMacaronPath("/api/v1/provisioning/templates/{name}"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodPut, "/api/v1/provisioning/templates/{name}"),
			metrics.Instrument(
				http.MethodPut,
				"/api/v1/provisioning/templates/{name}",
				api.Hooks.Wrap(srv.RoutePutTemplate),
				m,
			),
		)
		group.Delete(
			toMacaronPath("/api/v1/provisioning/policies"),
			requestmeta.SetOwner(requestmeta.TeamAlerting),
			requestmeta.SetSLOGroup(requestmeta.SLOGroupHighSlow),
			api.authorize(http.MethodDelete, "/api/v1/provisioning/policies"),
			metrics.Instrument(
				http.MethodDelete,
				"/api/v1/provisioning/policies",
				api.Hooks.Wrap(srv.RouteResetPolicyTree),
				m,
			),
		)
	}, middleware.ReqSignedIn)
}
