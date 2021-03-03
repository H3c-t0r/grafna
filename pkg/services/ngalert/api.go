package ngalert

import (
	"fmt"

	"github.com/go-macaron/binding"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana/pkg/api/response"
	"github.com/grafana/grafana/pkg/api/routing"
	"github.com/grafana/grafana/pkg/middleware"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/ngalert/eval"
	"github.com/grafana/grafana/pkg/tsdb"
	"github.com/grafana/grafana/pkg/util"
)

type apiImpl struct {
	ngalert         *AlertNG
	definitionStore definitionStore
	instanceStore   instanceStore
}

func (api *apiImpl) registerAPIEndpoints() {
	api.ngalert.RouteRegister.Group("/api/alert-definitions", func(alertDefinitions routing.RouteRegister) {
		alertDefinitions.Get("", middleware.ReqSignedIn, routing.Wrap(api.listAlertDefinitions))
		alertDefinitions.Get("/eval/:alertDefinitionUID", api.validateOrgAlertDefinition, routing.Wrap(api.alertDefinitionEvalEndpoint))
		alertDefinitions.Post("/eval", middleware.ReqSignedIn, binding.Bind(evalAlertConditionCommand{}), routing.Wrap(api.conditionEvalEndpoint))
		alertDefinitions.Get("/:alertDefinitionUID", api.validateOrgAlertDefinition, routing.Wrap(api.getAlertDefinitionEndpoint))
		alertDefinitions.Delete("/:alertDefinitionUID", api.validateOrgAlertDefinition, routing.Wrap(api.deleteAlertDefinitionEndpoint))
		alertDefinitions.Post("/", middleware.ReqEditorRole, binding.Bind(saveAlertDefinitionCommand{}), routing.Wrap(api.createAlertDefinitionEndpoint))
		alertDefinitions.Put("/:alertDefinitionUID", api.validateOrgAlertDefinition, binding.Bind(updateAlertDefinitionCommand{}), routing.Wrap(api.updateAlertDefinitionEndpoint))
		alertDefinitions.Post("/pause", middleware.ReqEditorRole, binding.Bind(updateAlertDefinitionPausedCommand{}), routing.Wrap(api.alertDefinitionPauseEndpoint))
		alertDefinitions.Post("/unpause", middleware.ReqEditorRole, binding.Bind(updateAlertDefinitionPausedCommand{}), routing.Wrap(api.alertDefinitionUnpauseEndpoint))
	})

	api.ngalert.RouteRegister.Group("/api/ngalert/", func(schedulerRouter routing.RouteRegister) {
		schedulerRouter.Post("/pause", routing.Wrap(api.pauseScheduler))
		schedulerRouter.Post("/unpause", routing.Wrap(api.unpauseScheduler))
	}, middleware.ReqOrgAdmin)

	api.ngalert.RouteRegister.Group("/api/alert-instances", func(alertInstances routing.RouteRegister) {
		alertInstances.Get("", middleware.ReqSignedIn, routing.Wrap(api.listAlertInstancesEndpoint))
	})
}

// conditionEvalEndpoint handles POST /api/alert-definitions/eval.
func (api *apiImpl) conditionEvalEndpoint(c *models.ReqContext, dto evalAlertConditionCommand) response.Response {
	if err := api.validateCondition(dto.Condition, c.SignedInUser, c.SkipCache); err != nil {
		return response.Error(400, "invalid condition", err)
	}

	evaluator := eval.Evaluator{Cfg: api.ngalert.Cfg}
	evalResults, err := evaluator.ConditionEval(&dto.Condition, timeNow())
	if err != nil {
		return response.Error(400, "Failed to evaluate conditions", err)
	}

	frame := evalResults.AsDataFrame()
	df := tsdb.NewDecodedDataFrames([]*data.Frame{&frame})
	instances, err := df.Encoded()
	if err != nil {
		return response.Error(400, "Failed to encode result dataframes", err)
	}

	return response.JSON(200, util.DynMap{
		"instances": instances,
	})
}

// alertDefinitionEvalEndpoint handles GET /api/alert-definitions/eval/:alertDefinitionUID.
func (api *apiImpl) alertDefinitionEvalEndpoint(c *models.ReqContext) response.Response {
	alertDefinitionUID := c.Params(":alertDefinitionUID")

	condition, err := api.LoadAlertCondition(alertDefinitionUID, c.SignedInUser.OrgId)
	if err != nil {
		return response.Error(400, "Failed to load alert definition conditions", err)
	}

	if err := api.validateCondition(*condition, c.SignedInUser, c.SkipCache); err != nil {
		return response.Error(400, "invalid condition", err)
	}

	evaluator := eval.Evaluator{Cfg: api.ngalert.Cfg}
	evalResults, err := evaluator.ConditionEval(condition, timeNow())
	if err != nil {
		return response.Error(400, "Failed to evaluate alert", err)
	}
	frame := evalResults.AsDataFrame()

	df := tsdb.NewDecodedDataFrames([]*data.Frame{&frame})
	if err != nil {
		return response.Error(400, "Failed to instantiate Dataframes from the decoded frames", err)
	}

	instances, err := df.Encoded()
	if err != nil {
		return response.Error(400, "Failed to encode result dataframes", err)
	}
	return response.JSON(200, util.DynMap{
		"instances": instances,
	})
}

// getAlertDefinitionEndpoint handles GET /api/alert-definitions/:alertDefinitionUID.
func (api *apiImpl) getAlertDefinitionEndpoint(c *models.ReqContext) response.Response {
	alertDefinitionUID := c.Params(":alertDefinitionUID")

	query := getAlertDefinitionByUIDQuery{
		UID:   alertDefinitionUID,
		OrgID: c.SignedInUser.OrgId,
	}

	if err := api.definitionStore.getAlertDefinitionByUID(&query); err != nil {
		return response.Error(500, "Failed to get alert definition", err)
	}

	return response.JSON(200, &query.Result)
}

// deleteAlertDefinitionEndpoint handles DELETE /api/alert-definitions/:alertDefinitionUID.
func (api *apiImpl) deleteAlertDefinitionEndpoint(c *models.ReqContext) response.Response {
	alertDefinitionUID := c.Params(":alertDefinitionUID")

	cmd := deleteAlertDefinitionByUIDCommand{
		UID:   alertDefinitionUID,
		OrgID: c.SignedInUser.OrgId,
	}

	if err := api.definitionStore.deleteAlertDefinitionByUID(&cmd); err != nil {
		return response.Error(500, "Failed to delete alert definition", err)
	}

	return response.Success("Alert definition deleted")
}

// updateAlertDefinitionEndpoint handles PUT /api/alert-definitions/:alertDefinitionUID.
func (api *apiImpl) updateAlertDefinitionEndpoint(c *models.ReqContext, cmd updateAlertDefinitionCommand) response.Response {
	cmd.UID = c.Params(":alertDefinitionUID")
	cmd.OrgID = c.SignedInUser.OrgId

	if err := api.validateCondition(cmd.Condition, c.SignedInUser, c.SkipCache); err != nil {
		return response.Error(400, "invalid condition", err)
	}

	if err := api.definitionStore.updateAlertDefinition(&cmd); err != nil {
		return response.Error(500, "Failed to update alert definition", err)
	}

	return response.JSON(200, cmd.Result)
}

// createAlertDefinitionEndpoint handles POST /api/alert-definitions.
func (api *apiImpl) createAlertDefinitionEndpoint(c *models.ReqContext, cmd saveAlertDefinitionCommand) response.Response {
	cmd.OrgID = c.SignedInUser.OrgId

	if err := api.validateCondition(cmd.Condition, c.SignedInUser, c.SkipCache); err != nil {
		return response.Error(400, "invalid condition", err)
	}

	if err := api.definitionStore.saveAlertDefinition(&cmd); err != nil {
		return response.Error(500, "Failed to create alert definition", err)
	}

	return response.JSON(200, cmd.Result)
}

// listAlertDefinitions handles GET /api/alert-definitions.
func (api *apiImpl) listAlertDefinitions(c *models.ReqContext) response.Response {
	query := listAlertDefinitionsQuery{OrgID: c.SignedInUser.OrgId}

	if err := api.definitionStore.getOrgAlertDefinitions(&query); err != nil {
		return response.Error(500, "Failed to list alert definitions", err)
	}

	return response.JSON(200, util.DynMap{"results": query.Result})
}

func (api *apiImpl) pauseScheduler() response.Response {
	err := api.ngalert.schedule.Pause()
	if err != nil {
		return response.Error(500, "Failed to pause scheduler", err)
	}
	return response.JSON(200, util.DynMap{"message": "alert definition scheduler paused"})
}

func (api *apiImpl) unpauseScheduler() response.Response {
	err := api.ngalert.schedule.Unpause()
	if err != nil {
		return response.Error(500, "Failed to unpause scheduler", err)
	}
	return response.JSON(200, util.DynMap{"message": "alert definition scheduler unpaused"})
}

// alertDefinitionPauseEndpoint handles POST /api/alert-definitions/pause.
func (api *apiImpl) alertDefinitionPauseEndpoint(c *models.ReqContext, cmd updateAlertDefinitionPausedCommand) response.Response {
	cmd.OrgID = c.SignedInUser.OrgId
	cmd.Paused = true

	err := api.definitionStore.updateAlertDefinitionPaused(&cmd)
	if err != nil {
		return response.Error(500, "Failed to pause alert definition", err)
	}
	return response.JSON(200, util.DynMap{"message": fmt.Sprintf("%d alert definitions paused", cmd.ResultCount)})
}

// alertDefinitionUnpauseEndpoint handles POST /api/alert-definitions/unpause.
func (api *apiImpl) alertDefinitionUnpauseEndpoint(c *models.ReqContext, cmd updateAlertDefinitionPausedCommand) response.Response {
	cmd.OrgID = c.SignedInUser.OrgId
	cmd.Paused = false

	err := api.definitionStore.updateAlertDefinitionPaused(&cmd)
	if err != nil {
		return response.Error(500, "Failed to unpause alert definition", err)
	}
	return response.JSON(200, util.DynMap{"message": fmt.Sprintf("%d alert definitions unpaused", cmd.ResultCount)})
}
