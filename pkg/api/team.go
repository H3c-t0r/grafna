package api

import (
	"errors"

	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/api/response"
	"github.com/grafana/grafana/pkg/api/utils"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/teamguardian"
	"github.com/grafana/grafana/pkg/util"
)

// POST /api/teams
func (hs *HTTPServer) CreateTeam(c *models.ReqContext, cmd models.CreateTeamCommand) response.Response {
	cmd.OrgId = c.OrgId

	if c.OrgRole == models.ROLE_VIEWER {
		return utils.Error(403, "Not allowed to create team.", nil)
	}

	if err := hs.Bus.Dispatch(&cmd); err != nil {
		if errors.Is(err, models.ErrTeamNameTaken) {
			return utils.Error(409, "Team name taken", err)
		}
		return utils.Error(500, "Failed to create Team", err)
	}

	if c.OrgRole == models.ROLE_EDITOR && hs.Cfg.EditorsCanAdmin {
		// if the request is authenticated using API tokens
		// the SignedInUser is an empty struct therefore
		// an additional check whether it is an actual user is required
		if c.SignedInUser.IsRealUser() {
			addMemberCmd := models.AddTeamMemberCommand{
				UserId:     c.SignedInUser.UserId,
				OrgId:      cmd.OrgId,
				TeamId:     cmd.Result.Id,
				Permission: models.PERMISSION_ADMIN,
			}

			if err := hs.Bus.Dispatch(&addMemberCmd); err != nil {
				c.Logger.Error("Could not add creator to team.", "error", err)
			}
		} else {
			c.Logger.Warn("Could not add creator to team because is not a real user.")
		}
	}

	return utils.JSON(200, &util.DynMap{
		"teamId":  cmd.Result.Id,
		"message": "Team created",
	})
}

// PUT /api/teams/:teamId
func (hs *HTTPServer) UpdateTeam(c *models.ReqContext, cmd models.UpdateTeamCommand) response.Response {
	cmd.OrgId = c.OrgId
	cmd.Id = c.ParamsInt64(":teamId")

	if err := teamguardian.CanAdmin(hs.Bus, cmd.OrgId, cmd.Id, c.SignedInUser); err != nil {
		return utils.Error(403, "Not allowed to update team", err)
	}

	if err := hs.Bus.Dispatch(&cmd); err != nil {
		if errors.Is(err, models.ErrTeamNameTaken) {
			return utils.Error(400, "Team name taken", err)
		}
		return utils.Error(500, "Failed to update Team", err)
	}

	return utils.Success("Team updated")
}

// DELETE /api/teams/:teamId
func (hs *HTTPServer) DeleteTeamByID(c *models.ReqContext) response.Response {
	orgId := c.OrgId
	teamId := c.ParamsInt64(":teamId")
	user := c.SignedInUser

	if err := teamguardian.CanAdmin(hs.Bus, orgId, teamId, user); err != nil {
		return utils.Error(403, "Not allowed to delete team", err)
	}

	if err := hs.Bus.Dispatch(&models.DeleteTeamCommand{OrgId: orgId, Id: teamId}); err != nil {
		if errors.Is(err, models.ErrTeamNotFound) {
			return utils.Error(404, "Failed to delete Team. ID not found", nil)
		}
		return utils.Error(500, "Failed to delete Team", err)
	}
	return utils.Success("Team deleted")
}

// GET /api/teams/search
func (hs *HTTPServer) SearchTeams(c *models.ReqContext) response.Response {
	perPage := c.QueryInt("perpage")
	if perPage <= 0 {
		perPage = 1000
	}
	page := c.QueryInt("page")
	if page < 1 {
		page = 1
	}

	var userIdFilter int64
	if hs.Cfg.EditorsCanAdmin && c.OrgRole != models.ROLE_ADMIN {
		userIdFilter = c.SignedInUser.UserId
	}

	query := models.SearchTeamsQuery{
		OrgId:        c.OrgId,
		Query:        c.Query("query"),
		Name:         c.Query("name"),
		UserIdFilter: userIdFilter,
		Page:         page,
		Limit:        perPage,
		SignedInUser: c.SignedInUser,
		HiddenUsers:  hs.Cfg.HiddenUsers,
	}

	if err := bus.Dispatch(&query); err != nil {
		return utils.Error(500, "Failed to search Teams", err)
	}

	for _, team := range query.Result.Teams {
		team.AvatarUrl = dtos.GetGravatarUrlWithDefault(team.Email, team.Name)
	}

	query.Result.Page = page
	query.Result.PerPage = perPage

	return utils.JSON(200, query.Result)
}

// GET /api/teams/:teamId
func (hs *HTTPServer) GetTeamByID(c *models.ReqContext) response.Response {
	query := models.GetTeamByIdQuery{
		OrgId:        c.OrgId,
		Id:           c.ParamsInt64(":teamId"),
		SignedInUser: c.SignedInUser,
		HiddenUsers:  hs.Cfg.HiddenUsers,
	}

	if err := bus.Dispatch(&query); err != nil {
		if errors.Is(err, models.ErrTeamNotFound) {
			return utils.Error(404, "Team not found", err)
		}

		return utils.Error(500, "Failed to get Team", err)
	}

	query.Result.AvatarUrl = dtos.GetGravatarUrlWithDefault(query.Result.Email, query.Result.Name)
	return utils.JSON(200, &query.Result)
}

// GET /api/teams/:teamId/preferences
func (hs *HTTPServer) GetTeamPreferences(c *models.ReqContext) response.Response {
	teamId := c.ParamsInt64(":teamId")
	orgId := c.OrgId

	if err := teamguardian.CanAdmin(hs.Bus, orgId, teamId, c.SignedInUser); err != nil {
		return utils.Error(403, "Not allowed to view team preferences.", err)
	}

	return getPreferencesFor(orgId, 0, teamId)
}

// PUT /api/teams/:teamId/preferences
func (hs *HTTPServer) UpdateTeamPreferences(c *models.ReqContext, dtoCmd dtos.UpdatePrefsCmd) response.Response {
	teamId := c.ParamsInt64(":teamId")
	orgId := c.OrgId

	if err := teamguardian.CanAdmin(hs.Bus, orgId, teamId, c.SignedInUser); err != nil {
		return utils.Error(403, "Not allowed to update team preferences.", err)
	}

	return updatePreferencesFor(orgId, 0, teamId, &dtoCmd)
}
