package api

import (
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
)

// API related actions
const (
	ActionProvisioningReload = "provisioning:reload"

	ActionDatasourcesRead   = "datasources:read"
	ActionDatasourcesQuery  = "datasources:query"
	ActionDatasourcesCreate = "datasources:create"
	ActionDatasourcesWrite  = "datasources:write"
	ActionDatasourcesDelete = "datasources:delete"
	ActionDatasourcesIDRead = "datasources.id:read"

	ActionOrgsRead             = "orgs:read"
	ActionOrgsPreferencesRead  = "orgs.preferences:read"
	ActionOrgsQuotasRead       = "orgs.quotas:read"
	ActionOrgsWrite            = "orgs:write"
	ActionOrgsPreferencesWrite = "orgs.preferences:write"
	ActionOrgsQuotasWrite      = "orgs.quotas:write"
	ActionOrgsDelete           = "orgs:delete"
	ActionOrgsCreate           = "orgs:create"
)

// API related scopes
var (
	ScopeProvisionersAll           = accesscontrol.Scope("provisioners", "*")
	ScopeProvisionersDashboards    = accesscontrol.Scope("provisioners", "dashboards")
	ScopeProvisionersPlugins       = accesscontrol.Scope("provisioners", "plugins")
	ScopeProvisionersDatasources   = accesscontrol.Scope("provisioners", "datasources")
	ScopeProvisionersNotifications = accesscontrol.Scope("provisioners", "notifications")

	ScopeDatasourcesAll = accesscontrol.Scope("datasources", "*")
	ScopeDatasourceID   = accesscontrol.Scope("datasources", "id", accesscontrol.Parameter(":id"))
	ScopeDatasourceUID  = accesscontrol.Scope("datasources", "uid", accesscontrol.Parameter(":uid"))
	ScopeDatasourceName = accesscontrol.Scope("datasources", "name", accesscontrol.Parameter(":name"))
)

// declareFixedRoles declares to the AccessControl service fixed roles and their
// grants to organization roles ("Viewer", "Editor", "Admin") or "Grafana Admin"
// that HTTPServer needs
func (hs *HTTPServer) declareFixedRoles() error {
	provisioningWriterRole := accesscontrol.RoleRegistration{
		Role: accesscontrol.RoleDTO{
			Version:     3,
			Name:        "fixed:provisioning:writer",
			DisplayName: "Provisioning writer",
			Description: "Reload provisioning.",
			Group:       "Provisioning",
			Permissions: []accesscontrol.Permission{
				{
					Action: ActionProvisioningReload,
					Scope:  ScopeProvisionersAll,
				},
			},
		},
		Grants: []string{accesscontrol.RoleGrafanaAdmin},
	}

	datasourcesReaderRole := accesscontrol.RoleRegistration{
		Role: accesscontrol.RoleDTO{
			Version:     3,
			Name:        "fixed:datasources:reader",
			DisplayName: "Data source reader",
			Description: "Read and query all data sources.",
			Group:       "Data sources",
			Permissions: []accesscontrol.Permission{
				{
					Action: ActionDatasourcesRead,
					Scope:  ScopeDatasourcesAll,
				},
				{
					Action: ActionDatasourcesQuery,
					Scope:  ScopeDatasourcesAll,
				},
			},
		},
		Grants: []string{string(models.ROLE_ADMIN)},
	}

	datasourcesWriterRole := accesscontrol.RoleRegistration{
		Role: accesscontrol.RoleDTO{
			Version:     3,
			Name:        "fixed:datasources:writer",
			DisplayName: "Data source writer",
			Description: "Create, update, delete, read, or query data sources.",
			Group:       "Data sources",
			Permissions: accesscontrol.ConcatPermissions(datasourcesReaderRole.Role.Permissions, []accesscontrol.Permission{
				{
					Action: ActionDatasourcesWrite,
					Scope:  ScopeDatasourcesAll,
				},
				{
					Action: ActionDatasourcesCreate,
				},
				{
					Action: ActionDatasourcesDelete,
					Scope:  ScopeDatasourcesAll,
				},
			}),
		},
		Grants: []string{string(models.ROLE_ADMIN)},
	}

	datasourcesIdReaderRole := accesscontrol.RoleRegistration{
		Role: accesscontrol.RoleDTO{
			Version:     4,
			Name:        "fixed:datasources.id:reader",
			DisplayName: "Data source ID reader",
			Description: "Read the ID of a data source based on its name.",
			Group:       "Infrequently used",
			Permissions: []accesscontrol.Permission{
				{
					Action: ActionDatasourcesIDRead,
					Scope:  ScopeDatasourcesAll,
				},
			},
		},
		Grants: []string{string(models.ROLE_VIEWER)},
	}

	datasourcesCompatibilityReaderRole := accesscontrol.RoleRegistration{
		Role: accesscontrol.RoleDTO{
			Version:     3,
			Name:        "fixed:datasources:compatibility:querier",
			DisplayName: "Data source compatibility querier",
			Description: "Only used for open source compatibility. Query data sources.",
			Group:       "Infrequently used",
			Permissions: []accesscontrol.Permission{
				{Action: ActionDatasourcesQuery},
			},
		},
		Grants: []string{string(models.ROLE_VIEWER)},
	}

	currentOrgReaderRole := accesscontrol.RoleRegistration{
		Role: accesscontrol.RoleDTO{
			Version:     4,
			Name:        "fixed:current.org:reader",
			DisplayName: "Current Organization reader",
			Description: "Read the current organization, such as its ID, name, address, or quotas.",
			Group:       "Organizations",
			Permissions: []accesscontrol.Permission{
				{Action: ActionOrgsRead},
				{Action: ActionOrgsQuotasRead},
			},
		},
		Grants: []string{string(models.ROLE_VIEWER)},
	}

	currentOrgWriterRole := accesscontrol.RoleRegistration{
		Role: accesscontrol.RoleDTO{
			Version:     4,
			Name:        "fixed:current.org:writer",
			DisplayName: "Current Organization writer",
			Description: "Read the current organization, its quotas, or its preferences. Update the current organization properties, or its preferences.",
			Group:       "Organizations",
			Permissions: accesscontrol.ConcatPermissions(currentOrgReaderRole.Role.Permissions, []accesscontrol.Permission{
				{Action: ActionOrgsPreferencesRead},
				{Action: ActionOrgsWrite},
				{Action: ActionOrgsPreferencesWrite},
			}),
		},
		Grants: []string{string(models.ROLE_ADMIN)},
	}

	orgReaderRole := accesscontrol.RoleRegistration{
		Role: accesscontrol.RoleDTO{
			Version:     2,
			Name:        "fixed:orgs:reader",
			DisplayName: "Organization reader",
			Description: "Read the organization and its quotas.",
			Group:       "Organizations",
			Permissions: []accesscontrol.Permission{
				{Action: ActionOrgsRead},
				{Action: ActionOrgsQuotasRead},
			},
		},
		Grants: []string{string(accesscontrol.RoleGrafanaAdmin)},
	}

	orgWriterRole := accesscontrol.RoleRegistration{
		Role: accesscontrol.RoleDTO{
			Version:     4,
			Name:        "fixed:orgs:writer",
			DisplayName: "Organization writer",
			Description: "Create, read, write, or delete an organization. Read or write an organization's quotas.",
			Group:       "Organizations",
			Permissions: accesscontrol.ConcatPermissions(orgReaderRole.Role.Permissions, []accesscontrol.Permission{
				{Action: ActionOrgsCreate},
				{Action: ActionOrgsWrite},
				{Action: ActionOrgsDelete},
				{Action: ActionOrgsQuotasWrite},
			}),
		},
		Grants: []string{string(accesscontrol.RoleGrafanaAdmin)},
	}

	return hs.AccessControl.DeclareFixedRoles(
		provisioningWriterRole, datasourcesReaderRole, datasourcesWriterRole, datasourcesIdReaderRole,
		datasourcesCompatibilityReaderRole, currentOrgReaderRole, currentOrgWriterRole, orgReaderRole, orgWriterRole,
	)
}

// Evaluators
// here is the list of complex evaluators we use in this package

// dataSourcesConfigurationAccessEvaluator is used to protect the "Configure > Data sources" tab access
var dataSourcesConfigurationAccessEvaluator = accesscontrol.EvalAll(
	accesscontrol.EvalPermission(ActionDatasourcesRead, ScopeDatasourcesAll),
	accesscontrol.EvalAny(
		accesscontrol.EvalPermission(ActionDatasourcesCreate),
		accesscontrol.EvalPermission(ActionDatasourcesDelete),
		accesscontrol.EvalPermission(ActionDatasourcesWrite),
	),
)

// dataSourcesNewAccessEvaluator is used to protect the "Configure > Data sources > New" page access
var dataSourcesNewAccessEvaluator = accesscontrol.EvalAll(
	accesscontrol.EvalPermission(ActionDatasourcesRead, ScopeDatasourcesAll),
	accesscontrol.EvalPermission(ActionDatasourcesCreate),
	accesscontrol.EvalPermission(ActionDatasourcesWrite),
)

// dataSourcesEditAccessEvaluator is used to protect the "Configure > Data sources > Edit" page access
var dataSourcesEditAccessEvaluator = accesscontrol.EvalAll(
	accesscontrol.EvalPermission(ActionDatasourcesRead, ScopeDatasourcesAll),
	accesscontrol.EvalPermission(ActionDatasourcesWrite),
)
