package database

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
)

const (
	// userAssignsSQL is a query to select all users assignments.
	userAssignsSQL = `SELECT ur.user_id, ur.org_id, ur.role_id
	FROM user_role AS ur`

	// teamAssignsSQL is a query to select all users' team assignments.
	teamAssignsSQL = `SELECT tm.user_id, tr.org_id, tr.role_id
	FROM team_role AS tr
	INNER JOIN team_member AS tm ON tm.team_id = tr.team_id`

	// basicRoleAssignsSQL is a query to select all users basic role (Admin, Editor, Viewer, None) assignments.
	basicRoleAssignsSQL = `SELECT ou.user_id, ou.org_id, br.role_id
	FROM builtin_role AS br
	INNER JOIN org_user AS ou ON ou.role = br.role`

	// grafanaAdminAssignsSQL is a query to select all grafana admin users.
	// it has to be formatted with the quoted user table.
	grafanaAdminAssignsSQL = `SELECT sa.user_id, br.org_id, br.role_id
	FROM builtin_role AS br
	INNER JOIN (
		SELECT u.id AS user_id
	    FROM %s AS u WHERE u.is_admin
	) AS sa ON 1 = 1
	WHERE br.role = ?`
)

func ProvideService(sql db.DB) *AccessControlStore {
	return &AccessControlStore{sql}
}

type AccessControlStore struct {
	sql db.DB
}

func (s *AccessControlStore) GetUserPermissions(ctx context.Context, query accesscontrol.GetUserPermissionsQuery) ([]accesscontrol.Permission, error) {
	result := make([]accesscontrol.Permission, 0)
	err := s.sql.WithDbSession(ctx, func(sess *db.Session) error {
		if query.UserID == 0 && len(query.TeamIDs) == 0 && len(query.Roles) == 0 {
			// no permission to fetch
			return nil
		}

		filter, params := accesscontrol.UserRolesFilter(query.OrgID, query.UserID, query.TeamIDs, query.Roles)

		q := `
		SELECT
			permission.action,
			permission.scope
			FROM permission
			INNER JOIN role ON role.id = permission.role_id
		` + filter

		if len(query.RolePrefixes) > 0 {
			q += " WHERE ( " + strings.Repeat("role.name LIKE ? OR ", len(query.RolePrefixes)-1)
			q += "role.name LIKE ? )"
			for i := range query.RolePrefixes {
				params = append(params, query.RolePrefixes[i]+"%")
			}
		}

		if err := sess.SQL(q, params...).Find(&result); err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func (s *AccessControlStore) GetBasicRolesPermissions(ctx context.Context, roles []string, orgID int64) ([]accesscontrol.Permission, error) {
	result := make([]accesscontrol.Permission, 0)
	err := s.sql.WithDbSession(ctx, func(sess *db.Session) error {
		if len(roles) == 0 {
			// no permission to fetch
			return nil
		}

		q := `
		SELECT
			permission.action,
			permission.scope
		FROM permission
		INNER JOIN role ON role.id = permission.role_id
		INNER JOIN (
			SELECT br.role_id FROM builtin_role AS br
			WHERE br.role IN (?` + strings.Repeat(", ?", len(roles)-1) + `)
			  AND (br.org_id = ? OR br.org_id = ?)
		) as all_role ON role.id = all_role.role_id;
		`

		params := make([]any, 0)
		for _, role := range roles {
			params = append(params, role)
		}
		params = append(params, orgID, accesscontrol.GlobalOrgID)

		if err := sess.SQL(q, params...).Find(&result); err != nil {
			return err
		}

		return nil
	})
	return result, err
}

// SearchUsersPermissions returns the list of user permissions in specific organization indexed by UserID
func (s *AccessControlStore) SearchUsersPermissions(ctx context.Context, orgID int64, options accesscontrol.SearchOptions) (map[int64][]accesscontrol.Permission, error) {
	type UserRBACPermission struct {
		UserID int64  `xorm:"user_id"`
		Action string `xorm:"action"`
		Scope  string `xorm:"scope"`
	}
	dbPerms := make([]UserRBACPermission, 0)

	var userID int64
	if options.NamespacedID != "" {
		var err error
		userID, err = options.ComputeUserID()
		if err != nil {
			return nil, err
		}
	}

	if err := s.sql.WithDbSession(ctx, func(sess *db.Session) error {
		roleNameFilterJoin := ""
		if len(options.RolePrefixes) > 0 {
			roleNameFilterJoin = "INNER JOIN role AS r ON up.role_id = r.id"
		}

		params := []any{}

		direct := userAssignsSQL
		if options.NamespacedID != "" {
			direct += " WHERE ur.user_id = ?"
			params = append(params, userID)
		}

		team := teamAssignsSQL
		if options.NamespacedID != "" {
			team += " WHERE tm.user_id = ?"
			params = append(params, userID)
		}

		basic := basicRoleAssignsSQL
		if options.NamespacedID != "" {
			basic += " WHERE ou.user_id = ?"
			params = append(params, userID)
		}

		grafanaAdmin := fmt.Sprintf(grafanaAdminAssignsSQL, s.sql.Quote("user"))
		params = append(params, accesscontrol.RoleGrafanaAdmin)
		if options.NamespacedID != "" {
			grafanaAdmin += " AND sa.user_id = ?"
			params = append(params, userID)
		}

		// Find permissions
		q := `
		SELECT
			user_id,
			p.action,
			p.scope
		FROM (
			` + direct + `
			UNION ALL
			` + team + `
			UNION ALL
			` + basic + `
			UNION ALL
			` + grafanaAdmin + `
		) AS up ` + roleNameFilterJoin + `
		INNER JOIN permission AS p ON up.role_id = p.role_id
		WHERE (up.org_id = ? OR up.org_id = ?)
		`
		params = append(params, orgID, accesscontrol.GlobalOrgID)

		if options.ActionPrefix != "" {
			q += ` AND p.action LIKE ?`
			params = append(params, options.ActionPrefix+"%")
		}
		if options.Action != "" {
			q += ` AND p.action = ?`
			params = append(params, options.Action)
		}
		if options.Scope != "" {
			// Search for scope and wildcard that include the scope
			scopes := append(options.Wildcards(), options.Scope)
			q += ` AND p.scope IN ( ? ` + strings.Repeat(", ?", len(scopes)-1) + ")"
			for i := range scopes {
				params = append(params, scopes[i])
			}
		}
		if len(options.RolePrefixes) > 0 {
			q += " AND ( " + strings.Repeat("r.name LIKE ? OR ", len(options.RolePrefixes)-1)
			q += "r.name LIKE ? )"
			for _, prefix := range options.RolePrefixes {
				params = append(params, prefix+"%")
			}
		}

		return sess.SQL(q, params...).Find(&dbPerms)
	}); err != nil {
		return nil, err
	}

	mapped := map[int64][]accesscontrol.Permission{}
	for i := range dbPerms {
		mapped[dbPerms[i].UserID] = append(mapped[dbPerms[i].UserID], accesscontrol.Permission{Action: dbPerms[i].Action, Scope: dbPerms[i].Scope})
	}

	return mapped, nil
}

// GetUsersBasicRoles returns the list of user basic roles (Admin, Editor, Viewer, Grafana Admin) indexed by UserID
func (s *AccessControlStore) GetUsersBasicRoles(ctx context.Context, userFilter []int64, orgID int64) (map[int64][]string, error) {
	type UserOrgRole struct {
		UserID  int64  `xorm:"id"`
		OrgRole string `xorm:"role"`
		IsAdmin bool   `xorm:"is_admin"`
	}
	dbRoles := make([]UserOrgRole, 0)
	if err := s.sql.WithDbSession(ctx, func(sess *db.Session) error {
		// Find roles
		q := `
		SELECT u.id, ou.role, u.is_admin
		FROM ` + s.sql.GetDialect().Quote("user") + ` AS u
		LEFT JOIN org_user AS ou ON u.id = ou.user_id
		WHERE (u.is_admin OR ou.org_id = ?)
		`
		params := []any{orgID}
		if len(userFilter) > 0 {
			q += "AND u.id IN (?" + strings.Repeat(",?", len(userFilter)-1) + ")"
			for _, u := range userFilter {
				params = append(params, u)
			}
		}

		return sess.SQL(q, params...).Find(&dbRoles)
	}); err != nil {
		return nil, err
	}

	roles := map[int64][]string{}
	for i := range dbRoles {
		if dbRoles[i].OrgRole != "" {
			roles[dbRoles[i].UserID] = []string{dbRoles[i].OrgRole}
		}
		if dbRoles[i].IsAdmin {
			roles[dbRoles[i].UserID] = append(roles[dbRoles[i].UserID], accesscontrol.RoleGrafanaAdmin)
		}
	}
	return roles, nil
}

func (s *AccessControlStore) DeleteUserPermissions(ctx context.Context, orgID, userID int64) error {
	err := s.sql.WithDbSession(ctx, func(sess *db.Session) error {
		roleDeleteQuery := "DELETE FROM user_role WHERE user_id = ?"
		roleDeleteParams := []any{roleDeleteQuery, userID}
		if orgID != accesscontrol.GlobalOrgID {
			roleDeleteQuery += " AND org_id = ?"
			roleDeleteParams = []any{roleDeleteQuery, userID, orgID}
		}

		// Delete user role assignments
		if _, err := sess.Exec(roleDeleteParams...); err != nil {
			return err
		}

		// only delete scopes to user if all permissions is removed (i.e. user is removed)
		if orgID == accesscontrol.GlobalOrgID {
			// Delete permissions that are scoped to user
			if _, err := sess.Exec("DELETE FROM permission WHERE scope = ?", accesscontrol.Scope("users", "id", strconv.FormatInt(userID, 10))); err != nil {
				return err
			}
		}

		roleQuery := "SELECT id FROM role WHERE name = ?"
		roleParams := []any{accesscontrol.ManagedUserRoleName(userID)}
		if orgID != accesscontrol.GlobalOrgID {
			roleQuery += " AND org_id = ?"
			roleParams = []any{accesscontrol.ManagedUserRoleName(userID), orgID}
		}

		var roleIDs []int64
		if err := sess.SQL(roleQuery, roleParams...).Find(&roleIDs); err != nil {
			return err
		}

		if len(roleIDs) == 0 {
			return nil
		}

		permissionDeleteQuery := "DELETE FROM permission WHERE role_id IN(? " + strings.Repeat(",?", len(roleIDs)-1) + ")"
		permissionDeleteParams := make([]any, 0, len(roleIDs)+1)
		permissionDeleteParams = append(permissionDeleteParams, permissionDeleteQuery)
		for _, id := range roleIDs {
			permissionDeleteParams = append(permissionDeleteParams, id)
		}

		// Delete managed user permissions
		if _, err := sess.Exec(permissionDeleteParams...); err != nil {
			return err
		}

		managedRoleDeleteQuery := "DELETE FROM role WHERE id IN(? " + strings.Repeat(",?", len(roleIDs)-1) + ")"
		managedRoleDeleteParams := []any{managedRoleDeleteQuery}
		for _, id := range roleIDs {
			managedRoleDeleteParams = append(managedRoleDeleteParams, id)
		}
		// Delete managed user roles
		if _, err := sess.Exec(managedRoleDeleteParams...); err != nil {
			return err
		}

		return nil
	})
	return err
}

func (s *AccessControlStore) DeleteTeamPermissions(ctx context.Context, orgID, teamID int64) error {
	err := s.sql.WithDbSession(ctx, func(sess *db.Session) error {
		roleDeleteQuery := "DELETE FROM team_role WHERE team_id = ? AND org_id = ?"
		roleDeleteParams := []any{roleDeleteQuery, teamID, orgID}

		// Delete team role assignments
		if _, err := sess.Exec(roleDeleteParams...); err != nil {
			return err
		}

		// Delete permissions that are scoped to the team
		if _, err := sess.Exec("DELETE FROM permission WHERE scope = ?", accesscontrol.Scope("teams", "id", strconv.FormatInt(teamID, 10))); err != nil {
			return err
		}

		// Delete the team managed role
		roleQuery := "SELECT id FROM role WHERE name = ? AND org_id = ?"
		roleParams := []any{accesscontrol.ManagedTeamRoleName(teamID), orgID}

		var roleIDs []int64
		if err := sess.SQL(roleQuery, roleParams...).Find(&roleIDs); err != nil {
			return err
		}

		if len(roleIDs) == 0 {
			return nil
		}

		permissionDeleteQuery := "DELETE FROM permission WHERE role_id IN(? " + strings.Repeat(",?", len(roleIDs)-1) + ")"
		permissionDeleteParams := make([]any, 0, len(roleIDs)+1)
		permissionDeleteParams = append(permissionDeleteParams, permissionDeleteQuery)
		for _, id := range roleIDs {
			permissionDeleteParams = append(permissionDeleteParams, id)
		}

		// Delete managed team permissions
		if _, err := sess.Exec(permissionDeleteParams...); err != nil {
			return err
		}

		managedRoleDeleteQuery := "DELETE FROM role WHERE id IN(? " + strings.Repeat(",?", len(roleIDs)-1) + ")"
		managedRoleDeleteParams := []any{managedRoleDeleteQuery}
		for _, id := range roleIDs {
			managedRoleDeleteParams = append(managedRoleDeleteParams, id)
		}
		// Delete managed team role
		if _, err := sess.Exec(managedRoleDeleteParams...); err != nil {
			return err
		}

		return nil
	})
	return err
}
