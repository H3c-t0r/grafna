package service

import (
	"context"

	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/dashboards"
)

// Gets public dashboard via generated Uid
func (dr *DashboardServiceImpl) GetPublicDashboard(ctx context.Context, dashboardUid string) (*models.Dashboard, error) {
	pubdash, d, err := dr.dashboardStore.GetPublicDashboard(dashboardUid)

	if err != nil {
		return nil, err
	}

	if pubdash == nil || d == nil {
		return nil, models.ErrPublicDashboardNotFound
	}

	if !pubdash.IsEnabled {
		return nil, models.ErrPublicDashboardNotFound
	}

	ts := pubdash.BuildTimeSettings(d)
	d.Data.SetPath([]string{"time", "from"}, ts.From)
	d.Data.SetPath([]string{"time", "to"}, ts.To)

	return d, nil
}

// GetPublicDashboardConfig is a helper method to retrieve the public dashboard configuration for a given dashboard from the database
func (dr *DashboardServiceImpl) GetPublicDashboardConfig(ctx context.Context, orgId int64, dashboardUid string) (*models.PublicDashboard, error) {
	pdc, err := dr.dashboardStore.GetPublicDashboardConfig(orgId, dashboardUid)
	if err != nil {
		return nil, err
	}

	return pdc, nil
}

// SavePublicDashboardConfig is a helper method to persist the sharing config
// to the database. It handles validations for sharing config and persistence
func (dr *DashboardServiceImpl) SavePublicDashboardConfig(ctx context.Context, dto *dashboards.SavePublicDashboardConfigDTO) (*models.PublicDashboard, error) {

	if len(dto.DashboardUid) == 0 {
		return nil, models.ErrDashboardIdentifierNotSet
	}

	cmd := models.SavePublicDashboardConfigCommand{
		DashboardUid:    dto.DashboardUid,
		OrgId:           dto.OrgId,
		PublicDashboard: *dto.PublicDashboard,
	}

	cmd.PublicDashboard.OrgId = dto.OrgId
	cmd.PublicDashboard.DashboardUid = dto.DashboardUid

	// populate additional fields
	if !cmd.PublicDashboard.IsPersisted() {
		cmd.PublicDashboard.CreatedBy = dto.UserId
		uid, err := dr.dashboardStore.GenerateNewPublicDashboardUid()
		if err != nil {
			return nil, err
		}
		cmd.PublicDashboard.Uid = uid
	}

	pubdash, err := dr.dashboardStore.SavePublicDashboardConfig(cmd)
	if err != nil {
		return nil, err
	}

	return pubdash, nil
}

func (dr *DashboardServiceImpl) BuildPublicDashboardMetricRequest(ctx context.Context, publicDashboardUid string, panelId int64) (dtos.MetricRequest, error) {
	publicDashboard, dashboard, err := dr.dashboardStore.GetPublicDashboard(publicDashboardUid)
	if err != nil {
		return dtos.MetricRequest{}, err
	}

	if !publicDashboard.IsEnabled {
		return dtos.MetricRequest{}, models.ErrPublicDashboardNotFound
	}

	queriesByPanel := models.GetQueriesFromDashboard(dashboard.Data)

	if _, ok := queriesByPanel[panelId]; !ok {
		return dtos.MetricRequest{}, models.ErrPublicDashboardPanelNotFound
	}

	ts := publicDashboard.BuildTimeSettings(dashboard)

	return dtos.MetricRequest{
		From:    ts.From,
		To:      ts.To,
		Queries: queriesByPanel[panelId],
	}, nil
}
