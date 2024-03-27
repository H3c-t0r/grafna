package api

import (
	"net/http"
	"strconv"

	"github.com/grafana/grafana/pkg/api/response"
	"github.com/grafana/grafana/pkg/api/routing"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/infra/tracing"
	"github.com/grafana/grafana/pkg/middleware"
	"github.com/grafana/grafana/pkg/services/cloudmigration"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
	"github.com/grafana/grafana/pkg/web"
)

type CloudMigrationAPI struct {
	cloudMigrationsService cloudmigration.Service
	routeRegister          routing.RouteRegister
	log                    log.Logger
	tracer                 tracing.Tracer
}

func RegisterApi(
	rr routing.RouteRegister,
	cms cloudmigration.Service,
	tracer tracing.Tracer,
) *CloudMigrationAPI {
	api := &CloudMigrationAPI{
		log:                    log.New("cloudmigrations.api"),
		routeRegister:          rr,
		cloudMigrationsService: cms,
		tracer:                 tracer,
	}
	api.registerEndpoints()
	return api
}

// RegisterAPIEndpoints Registers Endpoints on Grafana Router
func (cma *CloudMigrationAPI) registerEndpoints() {
	cma.routeRegister.Group("/api/cloudmigration", func(cloudMigrationRoute routing.RouteRegister) {
		// migration
		cloudMigrationRoute.Get("/migration", routing.Wrap(cma.GetMigrationList))
		cloudMigrationRoute.Post("/migration", routing.Wrap(cma.CreateMigration))
		cloudMigrationRoute.Get("/migration/:id", routing.Wrap(cma.GetMigration))
		cloudMigrationRoute.Delete("migration/:id", routing.Wrap(cma.DeleteMigration))
		cloudMigrationRoute.Post("/migration/:id/run", routing.Wrap(cma.RunMigration))
		cloudMigrationRoute.Get("/migration/:id/run", routing.Wrap(cma.GetMigrationRunList))
		cloudMigrationRoute.Get("/migration/:id/run/:runID", routing.Wrap(cma.GetMigrationRun))
		cloudMigrationRoute.Post("/token", routing.Wrap(cma.CreateToken))
	}, middleware.ReqGrafanaAdmin)
}

// swagger:route POST /cloudmigration/token migrations
//
// Create gcom access token.
//
// Responses:
// 200: cloudMigrationCreateTokenResponse
// 401: unauthorisedError
// 403: forbiddenError
// 500: internalServerError
func (cma *CloudMigrationAPI) CreateToken(c *contextmodel.ReqContext) response.Response {
	ctx, span := cma.tracer.Start(c.Req.Context(), "MigrationAPI.CreateAccessToken")
	defer span.End()

	logger := cma.log.FromContext(ctx)

	resp, err := cma.cloudMigrationsService.CreateToken(ctx)
	if err != nil {
		logger.Error("creating gcom access token", "err", err.Error())
		return response.Error(http.StatusInternalServerError, "creating gcom access token", err)
	}

	return response.JSON(http.StatusOK, cloudmigration.CreateAccessTokenResponseDTO(resp))
}

// swagger:route GET /cloudmigration/migration migrations
//
// Get a list of all cloud migrations.
//
// Responses:
// 200: cloudMigrationListResponse
// 401: unauthorisedError
// 403: forbiddenError
// 500: internalServerError
func (cma *CloudMigrationAPI) GetMigrationList(c *contextmodel.ReqContext) response.Response {
	cloudMigrations, err := cma.cloudMigrationsService.GetMigrationList(c.Req.Context())
	if err != nil {
		return response.Error(http.StatusInternalServerError, "migration list error", err)
	}

	cloudMigrationListResponse := cloudmigration.CloudMigrationListResponse{Migrations: cloudMigrations}
	return response.JSON(http.StatusOK, cloudMigrationListResponse)
}

// swagger:route GET /cloudmigration/migration/{id} migrations migrationWithIDPathParam
//
// Get a cloud migration.
//
// It returns migrations that has been created.
//
// Responses:
// 200: cloudMigrationResponse
// 401: unauthorisedError
// 403: forbiddenError
// 500: internalServerError
func (cma *CloudMigrationAPI) GetMigration(c *contextmodel.ReqContext) response.Response {
	id, err := strconv.ParseInt(web.Params(c.Req)[":id"], 10, 64)
	if err != nil {
		return response.Error(http.StatusBadRequest, "id is invalid", err)
	}
	cloudMigration, err := cma.cloudMigrationsService.GetMigration(c.Req.Context(), id)
	if err != nil {
		return response.Error(http.StatusNotFound, "migration not found", err)
	}
	return response.JSON(http.StatusOK, cloudMigration)
}

// swagger:route POST /cloudmigration/migration migrations createMigration
//
// Create a migration.
//
// Responses:
// 200: cloudMigrationResponse
// 401: unauthorisedError
// 403: forbiddenError
// 500: internalServerError
func (cma *CloudMigrationAPI) CreateMigration(c *contextmodel.ReqContext) response.Response {
	cmd := cloudmigration.CloudMigrationRequest{}
	if err := web.Bind(c.Req, &cmd); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	cloudMigration, err := cma.cloudMigrationsService.CreateMigration(c.Req.Context(), cmd)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "migration creation error", err)
	}
	return response.JSON(http.StatusOK, cloudMigration)
}

// swagger:route GET /cloudmigration/migration/{id}/run migrations migrationWithIDPathParam
//
// Trigger the run of a migration to the Grafana Cloud.
//
// It returns migrations that has been created.
//
// Responses:
// 200: cloudMigrationRunResponse
// 401: unauthorisedError
// 403: forbiddenError
// 500: internalServerError
func (cma *CloudMigrationAPI) RunMigration(c *contextmodel.ReqContext) response.Response {
	cloudMigrationRun, err := cma.cloudMigrationsService.RunMigration(c.Req.Context(), web.Params(c.Req)[":id"])
	if err != nil {
		return response.Error(http.StatusInternalServerError, "migration run error", err)
	}
	return response.JSON(http.StatusOK, cloudMigrationRun)
}

// swagger:route GET /cloudmigration/migration/{id}/run/{runID} migrations migrationRunWithIDPathParam
//
// Get the result of a single migration run.
//
// Responses:
// 200: cloudMigrationRunResponse
// 401: unauthorisedError
// 403: forbiddenError
// 500: internalServerError
func (cma *CloudMigrationAPI) GetMigrationRun(c *contextmodel.ReqContext) response.Response {
	migrationStatus, err := cma.cloudMigrationsService.GetMigrationStatus(c.Req.Context(), web.Params(c.Req)[":id"], web.Params(c.Req)[":runID"])
	if err != nil {
		return response.Error(http.StatusInternalServerError, "migration status error", err)
	}
	return response.JSON(http.StatusOK, migrationStatus)
}

// swagger:route GET /cloudmigration/migration/{id}/run migrations migrationWithIDPathParam
//
// Get a list of migration runs for a migration.
//
// Responses:
// 200: cloudMigrationRunListResponse
// 401: unauthorisedError
// 403: forbiddenError
// 500: internalServerError
func (cma *CloudMigrationAPI) GetMigrationRunList(c *contextmodel.ReqContext) response.Response {
	migrationStatus, err := cma.cloudMigrationsService.GetMigrationStatusList(c.Req.Context(), web.Params(c.Req)[":id"])
	if err != nil {
		return response.Error(http.StatusInternalServerError, "migration status error", err)
	}

	runList := cloudmigration.CloudMigrationRunList{Runs: migrationStatus}
	return response.JSON(http.StatusOK, runList)
}

// swagger:route DELETE /cloudmigration/migration/{id} migrations migrationWithIDPathParam
//
// Delete a migration.
//
// Responses:
// 200: 200okResponse
// 401: unauthorisedError
// 403: forbiddenError
// 500: internalServerError
func (cma *CloudMigrationAPI) DeleteMigration(c *contextmodel.ReqContext) response.Response {
	err := cma.cloudMigrationsService.DeleteMigration(c.Req.Context(), web.Params(c.Req)[":id"])
	if err != nil {
		return response.Error(http.StatusInternalServerError, "migration delete error", err)
	}
	return response.Empty(http.StatusOK)
}
