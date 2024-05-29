package ossaccesscontrol

import (
	"context"
	"errors"

	"github.com/grafana/grafana/pkg/api/routing"
	"github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/infra/metrics"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/accesscontrol/resourcepermissions"
	"github.com/grafana/grafana/pkg/services/dashboards"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/folder"
	"github.com/grafana/grafana/pkg/services/licensing"
	"github.com/grafana/grafana/pkg/services/team"
	"github.com/grafana/grafana/pkg/services/user"
	"github.com/grafana/grafana/pkg/setting"
)

type DashboardPermissionsService struct {
	*resourcepermissions.Service
}

var DashboardViewActions = []string{dashboards.ActionDashboardsRead}
var DashboardEditActions = append(DashboardViewActions, []string{dashboards.ActionDashboardsWrite, dashboards.ActionDashboardsDelete}...)
var DashboardAdminActions = append(DashboardEditActions, []string{dashboards.ActionDashboardsPermissionsRead, dashboards.ActionDashboardsPermissionsWrite}...)

func getDashboardViewActions(features featuremgmt.FeatureToggles) []string {
	if features.IsEnabled(context.Background(), featuremgmt.FlagAnnotationPermissionUpdate) {
		return append(DashboardViewActions, accesscontrol.ActionAnnotationsRead)
	}
	return DashboardViewActions
}

func getDashboardEditActions(features featuremgmt.FeatureToggles) []string {
	if features.IsEnabled(context.Background(), featuremgmt.FlagAnnotationPermissionUpdate) {
		return append(DashboardEditActions, []string{accesscontrol.ActionAnnotationsRead, accesscontrol.ActionAnnotationsWrite, accesscontrol.ActionAnnotationsDelete, accesscontrol.ActionAnnotationsCreate}...)
	}
	return DashboardEditActions
}

func getDashboardAdminActions(features featuremgmt.FeatureToggles) []string {
	if features.IsEnabled(context.Background(), featuremgmt.FlagAnnotationPermissionUpdate) {
		return append(DashboardAdminActions, []string{accesscontrol.ActionAnnotationsRead, accesscontrol.ActionAnnotationsWrite, accesscontrol.ActionAnnotationsDelete, accesscontrol.ActionAnnotationsCreate}...)
	}
	return DashboardAdminActions
}

func ProvideDashboardPermissions(
	cfg *setting.Cfg, features featuremgmt.FeatureToggles, router routing.RouteRegister, sql db.DB, ac accesscontrol.AccessControl,
	license licensing.Licensing, dashboardStore dashboards.Store, folderService folder.Service, service accesscontrol.Service,
	teamService team.Service, userService user.Service, actionSetService resourcepermissions.ActionSetService,
) (*DashboardPermissionsService, error) {
	getDashboard := func(ctx context.Context, orgID int64, resourceID string) (*dashboards.Dashboard, error) {
		query := &dashboards.GetDashboardQuery{UID: resourceID, OrgID: orgID}
		queryResult, err := dashboardStore.GetDashboard(ctx, query)
		if err != nil {
			return nil, err
		}
		return queryResult, nil
	}

	options := resourcepermissions.Options{
		Resource:          "dashboards",
		ResourceAttribute: "uid",
		ResourceValidator: func(ctx context.Context, orgID int64, resourceID string) error {
			dashboard, err := getDashboard(ctx, orgID, resourceID)
			if err != nil {
				return err
			}

			if dashboard.IsFolder {
				return errors.New("not found")
			}

			return nil
		},
		InheritedScopesSolver: func(ctx context.Context, orgID int64, resourceID string) ([]string, error) {
			dashboard, err := getDashboard(ctx, orgID, resourceID)
			if err != nil {
				return nil, err
			}
			metrics.MFolderIDsServiceCount.WithLabelValues(metrics.AccessControl).Inc()
			// nolint:staticcheck
			if dashboard.FolderID > 0 {
				query := &dashboards.GetDashboardQuery{ID: dashboard.FolderID, OrgID: orgID}
				queryResult, err := dashboardStore.GetDashboard(ctx, query)
				if err != nil {
					return nil, err
				}
				parentScope := dashboards.ScopeFoldersProvider.GetResourceScopeUID(queryResult.UID)

				nestedScopes, err := dashboards.GetInheritedScopes(ctx, orgID, queryResult.UID, folderService)
				if err != nil {
					return nil, err
				}
				return append([]string{parentScope}, nestedScopes...), nil
			}
			return []string{dashboards.ScopeFoldersProvider.GetResourceScopeUID(folder.GeneralFolderUID)}, nil
		},
		Assignments: resourcepermissions.Assignments{
			Users:           true,
			Teams:           true,
			BuiltInRoles:    true,
			ServiceAccounts: true,
		},
		PermissionsToActions: map[string][]string{
			"View":  getDashboardViewActions(features),
			"Edit":  getDashboardEditActions(features),
			"Admin": getDashboardAdminActions(features),
		},
		ReaderRoleName: "Dashboard permission reader",
		WriterRoleName: "Dashboard permission writer",
		RoleGroup:      "Dashboards",
	}

	srv, err := resourcepermissions.New(cfg, options, features, router, license, ac, service, sql, teamService, userService, actionSetService)
	if err != nil {
		return nil, err
	}
	return &DashboardPermissionsService{srv}, nil
}
