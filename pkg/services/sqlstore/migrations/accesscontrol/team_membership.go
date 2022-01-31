package accesscontrol

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"xorm.io/xorm"

	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/sqlstore/migrator"
	"github.com/grafana/grafana/pkg/util"
)

const batchSize = 500

func AddTeamMembershipMigrations(mg *migrator.Migrator) {
	mg.AddMigration("teams permissions migration", &teamPermissionMigrator{editorsCanAdmin: mg.Cfg.EditorsCanAdmin})
}

var _ migrator.CodeMigration = new(teamPermissionMigrator)

type teamPermissionMigrator struct {
	migrator.MigrationBase
	editorsCanAdmin bool
	sess            *xorm.Session
}

func (p *teamPermissionMigrator) getAssignmentKey(orgID int64, name string) string {
	return fmt.Sprint(orgID, "-", name)
}

func (p *teamPermissionMigrator) SQL(dialect migrator.Dialect) string {
	return "code migration"
}

func (p *teamPermissionMigrator) Exec(sess *xorm.Session, migrator *migrator.Migrator) error {
	p.sess = sess
	return p.migrateMemberships()
}

func generateNewRoleUID(sess *xorm.Session, orgID int64) (string, error) {
	for i := 0; i < 3; i++ {
		uid := util.GenerateShortUID()

		exists, err := sess.Where("org_id=? AND uid=?", orgID, uid).Get(&accesscontrol.Role{})
		if err != nil {
			return "", err
		}

		if !exists {
			return uid, nil
		}
	}

	return "", fmt.Errorf("failed to generate uid")
}

func (p *teamPermissionMigrator) findRole(orgID int64, name string) (*accesscontrol.Role, error) {
	// check if role exists
	var role accesscontrol.Role
	_, err := p.sess.Table("role").Select("id").Where("org_id = ? AND name = ?", orgID, name).Get(&role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func batch(count, batchSize int, eachFn func(start, end int) error) error {
	for i := 0; i < count; {
		end := i + batchSize - 1
		if end > count-1 {
			end = count - 1
		}

		if err := eachFn(i, end); err != nil {
			return err
		}

		i = end + 1
	}

	return nil
}

func (p *teamPermissionMigrator) bulkCreateRoles(allRoles []*accesscontrol.Role) ([]*accesscontrol.Role, error) {
	ts := time.Now()
	allCreatedRoles := make([]*accesscontrol.Role, len(allRoles))

	// bulk role creations
	err := batch(len(allRoles), batchSize, func(start, end int) error {
		roles := allRoles[start:end]
		createdRoles := make([]*accesscontrol.Role, len(roles))
		valueStrings := make([]string, len(roles))
		args := make([]interface{}, 0, len(roles)*5)

		for i, r := range roles {
			uid, err := generateNewRoleUID(p.sess, r.OrgID)
			if err != nil {
				return err
			}

			valueStrings[i] = "(?, ?, ?, 1, ?, ?)"
			args = append(args, r.OrgID, uid, r.Name, ts, ts)
		}

		if len(valueStrings) == 0 {
			return nil
		}

		valueString := strings.Join(valueStrings, ",")
		sql := fmt.Sprintf("INSERT INTO role (org_id, uid, name, version, created, updated) VALUES %s RETURNING id, org_id, name", valueString)
		if errCreate := p.sess.SQL(sql, args...).Find(&createdRoles); errCreate != nil {
			return errCreate
		}

		allCreatedRoles = append(allCreatedRoles, createdRoles...)
		return nil
	})

	return allCreatedRoles, err
}

func (p *teamPermissionMigrator) bulkAssignRoles(rolesMap map[string]*accesscontrol.Role, assignments map[int64]map[string]struct{}) error {
	ts := time.Now()

	roleAssignments := make([]accesscontrol.UserRole, len(assignments))
	for userID, rolesByRoleKey := range assignments {
		for key := range rolesByRoleKey {
			role, ok := rolesMap[key]
			if !ok {
				return &ErrUnknownRole{key}
			}

			roleAssignments = append(roleAssignments, accesscontrol.UserRole{
				OrgID:   role.OrgID,
				RoleID:  role.ID,
				UserID:  userID,
				Created: ts,
			})
		}
	}

	if len(roleAssignments) == 0 {
		return nil
	}

	return batch(len(roleAssignments), batchSize, func(start, end int) error {
		roleAssignmentsChunk := roleAssignments[start:end]
		_, err := p.sess.Table("user_role").InsertMulti(roleAssignmentsChunk)
		return err
	})
}

// setRolePermissions sets the role permissions deleting any team related ones before inserting any.
func (p *teamPermissionMigrator) setRolePermissions(roleID int64, permissions []accesscontrol.Permission) error {
	// First drop existing permissions
	if _, errDeletingPerms := p.sess.SQL("DELETE FROM permission WHERE role_id = ? AND (action LIKE ? OR action LIKE ?)", roleID, "teams:%", "teams.permissions:%").Exec(); errDeletingPerms != nil {
		return errDeletingPerms
	}

	// Then insert new permissions
	var newPermissions []accesscontrol.Permission
	now := time.Now()
	for _, permission := range permissions {
		permission.RoleID = roleID
		permission.Created = now
		permission.Updated = now
		newPermissions = append(newPermissions, permission)
	}

	if _, errInsertPerms := p.sess.InsertMulti(&newPermissions); errInsertPerms != nil {
		return errInsertPerms
	}

	return nil
}

// mapPermissionToFGAC translates the legacy membership (Member or Admin) into FGAC permissions
func (p *teamPermissionMigrator) mapPermissionToFGAC(permission models.PermissionType, teamID int64) []accesscontrol.Permission {
	teamIDScope := accesscontrol.Scope("teams", "id", strconv.FormatInt(teamID, 10))
	switch permission {
	case 0:
		return []accesscontrol.Permission{{Action: "teams:read", Scope: teamIDScope}}
	case models.PERMISSION_ADMIN:
		return []accesscontrol.Permission{
			{Action: "teams:delete", Scope: teamIDScope},
			{Action: "teams:read", Scope: teamIDScope},
			{Action: "teams:write", Scope: teamIDScope},
			{Action: "teams.permissions:read", Scope: teamIDScope},
			{Action: "teams.permissions:write", Scope: teamIDScope},
		}
	default:
		return []accesscontrol.Permission{}
	}
}

func (p *teamPermissionMigrator) getUserRoleByOrgMapping() (map[int64]map[int64]string, error) {
	var orgUsers []*models.OrgUserDTO
	if err := p.sess.SQL(`SELECT * FROM org_user`).Cols("org_user.org_id", "org_user.user_id", "org_user.role").Find(&orgUsers); err != nil {
		return nil, err
	}

	userRolesByOrg := map[int64]map[int64]string{}

	// Loop through users and organise them by organization ID
	for _, orgUser := range orgUsers {
		orgRoles, initialized := userRolesByOrg[orgUser.OrgId]
		if !initialized {
			orgRoles = map[int64]string{}
		}

		orgRoles[orgUser.UserId] = orgUser.Role
		userRolesByOrg[orgUser.OrgId] = orgRoles
	}

	return userRolesByOrg, nil
}

// TODO SPLIT
// migrateMemberships generate managed permissions for users based on their memberships to teams
func (p *teamPermissionMigrator) migrateMemberships() error {
	userRolesByOrg, err := p.getUserRoleByOrgMapping()
	if err != nil {
		return err
	}

	var teamMemberships []models.TeamMember
	if err := p.sess.SQL(`SELECT * FROM team_member`).Find(&teamMemberships); err != nil {
		return err
	}

	// Loop through memberships and generate associated permissions
	// Downgrade team permissions if needed - only organisation admins or organisation editors (when editorsCanAdmin feature is enabled)
	// can access team administration endpoints
	// log?
	userPermissionsByOrg, errGen := p.generateAssociatedPermissions(teamMemberships, userRolesByOrg)
	if errGen != nil {
		return errGen
	}

	// Create a map of roles to create
	// userID orgID-RoleName -> Role
	// orgID-RoleName -> Role
	rolesToCreate, assignments, rolesByOrg, errOrganizeRoles := p.sortRolesToAssign(userPermissionsByOrg)
	if errOrganizeRoles != nil {
		return errOrganizeRoles
	}

	// Create missing roles
	createdRoles, errCreate := p.bulkCreateRoles(rolesToCreate)
	if errCreate != nil {
		return errCreate
	}

	// Populate rolesMap with the newly created roles
	for i := range createdRoles {
		roleKey := p.getAssignmentKey(createdRoles[i].OrgID, createdRoles[i].Name)
		rolesByOrg[roleKey] = createdRoles[i]
	}

	// Assign missing roles
	if errAssign := p.bulkAssignRoles(rolesByOrg, assignments); errAssign != nil {
		return errAssign
	}

	// Set roles permissions
	return p.setRolePermissionsForOrgs(userPermissionsByOrg, rolesByOrg)
}

func (p *teamPermissionMigrator) setRolePermissionsForOrgs(userPermissionsByOrg map[int64]map[int64][]accesscontrol.Permission, rolesByOrg map[string]*accesscontrol.Role) error {
	for orgID, userPermissions := range userPermissionsByOrg {
		for userID, permissions := range userPermissions {
			key := p.getAssignmentKey(orgID, fmt.Sprintf("managed:users:%d:permissions", userID))

			role, ok := rolesByOrg[key]
			if !ok {
				return &ErrUnknownRole{key}
			}

			if errSettingPerms := p.setRolePermissions(role.ID, permissions); errSettingPerms != nil {
				return errSettingPerms
			}
		}
	}
	return nil
}

func (p *teamPermissionMigrator) sortRolesToAssign(userPermissionsByOrg map[int64]map[int64][]accesscontrol.Permission) ([]*accesscontrol.Role, map[int64]map[string]struct{}, map[string]*accesscontrol.Role, error) {
	var rolesToCreate []*accesscontrol.Role

	assignments := map[int64]map[string]struct{}{}

	rolesByOrg := map[string]*accesscontrol.Role{}
	for orgID, userPermissions := range userPermissionsByOrg {
		for userID := range userPermissions {
			roleName := fmt.Sprintf("managed:users:%d:permissions", userID)
			role, errFindingRoles := p.findRole(orgID, roleName)
			if errFindingRoles != nil {
				return nil, nil, nil, errFindingRoles
			}

			roleKey := p.getAssignmentKey(orgID, roleName)

			if role != nil {
				rolesByOrg[roleKey] = role
			} else {
				roleToCreate := &accesscontrol.Role{
					Name:  roleName,
					OrgID: orgID,
				}
				rolesToCreate = append(rolesToCreate, roleToCreate)

				userAssignments, initialized := assignments[userID]
				if !initialized {
					userAssignments = map[string]struct{}{}
				}

				userAssignments[roleKey] = struct{}{}
				assignments[userID] = userAssignments
			}
		}
	}

	return rolesToCreate, assignments, rolesByOrg, nil
}

func (p *teamPermissionMigrator) generateAssociatedPermissions(teamMemberships []models.TeamMember,
	userRolesByOrg map[int64]map[int64]string) (map[int64]map[int64][]accesscontrol.Permission, error) {
	userPermissionsByOrg := map[int64]map[int64][]accesscontrol.Permission{}

	for _, m := range teamMemberships {
		if m.Permission == models.PERMISSION_ADMIN {
			if userRolesByOrg[m.OrgId][m.UserId] == string(models.ROLE_VIEWER) || (userRolesByOrg[m.OrgId][m.UserId] == string(models.ROLE_EDITOR) && !p.editorsCanAdmin) {
				m.Permission = 0

				if _, err := p.sess.Cols("permission").Where("org_id=? and team_id=? and user_id=?", m.OrgId, m.TeamId, m.UserId).Update(m); err != nil {
					return nil, err
				}
			}
		}

		userPermissions, initialized := userPermissionsByOrg[m.OrgId]
		if !initialized {
			userPermissions = map[int64][]accesscontrol.Permission{}
		}
		userPermissions[m.UserId] = append(userPermissions[m.UserId], p.mapPermissionToFGAC(m.Permission, m.TeamId)...)
		userPermissionsByOrg[m.OrgId] = userPermissions
	}

	return userPermissionsByOrg, nil
}
