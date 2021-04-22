package accesscontrol

import (
	"time"
)

type Role struct {
	Version     int64  `json:"version"`
	UID         string `json:"uid"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Updated time.Time `json:"updated"`
	Created time.Time `json:"created"`
}

type RoleDTO struct {
	Version     int64        `json:"version"`
	UID         string       `json:"uid"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions,omitempty"`
}

type Permission struct {
	Action string `json:"action"`
	Scope  string `json:"scope"`
}

type EvaluationResult struct {
	HasAccess bool
	Meta      interface{}
}

func (p RoleDTO) Role() Role {
	return Role{
		Name:        p.Name,
		Description: p.Description,
	}
}

const (
	// Permission actions

	// Users actions
	ActionUsersRead     = "users:read"
	ActionUsersWrite    = "users:write"
	ActionUsersTeamRead = "users.teams:read"
	// We can ignore gosec G101 since this does not contain any credentials
	// nolint:gosec
	ActionUsersAuthTokenList = "users.authtoken:list"
	// We can ignore gosec G101 since this does not contain any credentials
	// nolint:gosec
	ActionUsersAuthTokenUpdate = "users.authtoken:update"
	// We can ignore gosec G101 since this does not contain any credentials
	// nolint:gosec
	ActionUsersPasswordUpdate    = "users.password:update"
	ActionUsersDelete            = "users:delete"
	ActionUsersCreate            = "users:create"
	ActionUsersEnable            = "users:enable"
	ActionUsersDisable           = "users:disable"
	ActionUsersPermissionsUpdate = "users.permissions:update"
	ActionUsersLogout            = "users:logout"
	ActionUsersQuotasList        = "users.quotas:list"
	ActionUsersQuotasUpdate      = "users.quotas:update"

	// Org actions
	ActionOrgUsersRead       = "org.users:read"
	ActionOrgUsersAdd        = "org.users:add"
	ActionOrgUsersRemove     = "org.users:remove"
	ActionOrgUsersRoleUpdate = "org.users.role:update"

	// LDAP actions
	ActionLDAPUsersRead  = "ldap.user:read"
	ActionLDAPUsersSync  = "ldap.user:sync"
	ActionLDAPStatusRead = "ldap.status:read"

	// Access control actions
	ActionAccessControlRolesList   = "accesscontrol.roles:list"
	ActionAccessControlRolesRead   = "accesscontrol.roles:read"
	ActionAccessControlRolesCreate = "accesscontrol.roles:create"
	ActionAccessControlRolesUpdate = "accesscontrol.roles:update"
	ActionAccessControlRolesDelete = "accesscontrol.roles:delete"

	ActionAccessControlTeamGrantList   = "accesscontrol.team.grants:list"
	ActionAccessControlTeamGrantCreate = "accesscontrol.team.grants:create"
	ActionAccessControlTeamGrantDelete = "accesscontrol.team.grants:delete"

	ActionAccessControlUserGrantList   = "accesscontrol.user.grants:list"
	ActionAccessControlUserGrantCreate = "accesscontrol.user.grants:create"
	ActionAccessControlUserGrantDelete = "accesscontrol.user.grants:delete"

	ActionAccessControlBuiltInRolesList   = "accesscontrol.builtin.grants:list"
	ActionAccessControlBuiltInRolesCreate = "accesscontrol.builtin.grants:create"
	ActionAccessControlBuiltInRolesDelete = "accesscontrol.builtin.grants:delete"

	// Global Scopes
	ScopeUsersAll  = "users:*"
	ScopeUsersSelf = "users:self"

	ScopeOrgAllUsersAll     = "org:*/users:*"
	ScopeOrgCurrentUsersAll = "org:current/users:*"

	ScopeOrgAllRolesAll = "org:*/roles:*"
)

const RoleGrafanaAdmin = "Grafana Admin"
