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
	RouteGetPolicyTree(*models.ReqContext) response.Response
	RoutePostPolicyTree(*models.ReqContext) response.Response
}

func (f *ForkedProvisioningApi) RouteGetPolicyTree(ctx *models.ReqContext) response.Response {
	return f.forkRouteGetPolicyTree(ctx)
}

func (f *ForkedProvisioningApi) RoutePostPolicyTree(ctx *models.ReqContext) response.Response {
	conf := apimodels.Route{}
	if err := web.Bind(ctx.Req, &conf); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	return f.forkRoutePostPolicyTree(ctx, conf)
}

func (api *API) RegisterProvisioningApiEndpoints(srv ProvisioningApiForkingService, m *metrics.API) {
	api.RouteRegister.Group("", func(group routing.RouteRegister) {
		group.Get(
			toMacaronPath("/api/provisioning/policies"),
			api.authorize(http.MethodGet, "/api/provisioning/policies"),
			metrics.Instrument(
				http.MethodGet,
				"/api/provisioning/policies",
				srv.RouteGetPolicyTree,
				m,
			),
		)
		group.Post(
			toMacaronPath("/api/provisioning/policies"),
			api.authorize(http.MethodPost, "/api/provisioning/policies"),
			metrics.Instrument(
				http.MethodPost,
				"/api/provisioning/policies",
				srv.RoutePostPolicyTree,
				m,
			),
		)
	}, middleware.ReqSignedIn)
}
