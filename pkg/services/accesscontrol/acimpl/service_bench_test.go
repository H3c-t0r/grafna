package acimpl

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/grafana/grafana/pkg/extensions/accesscontrol/database"
	"github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/org"
	"github.com/grafana/grafana/pkg/services/user"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/stretchr/testify/require"
)

const batchSize = 500

func batch(count, size int, eachFn func(start, end int) error) error {
	for i := 0; i < count; {
		end := i + size
		if end > count {
			end = count
		}

		if err := eachFn(i, end); err != nil {
			return err
		}

		i = end
	}

	return nil
}

func setupBenchEnv(b *testing.B, usersCount, resourceCount int) (accesscontrol.Service, *user.SignedInUser) {
	now := time.Now()
	sqlStore := db.InitTestDB(b)
	store := database.ProvideService(sqlStore)
	acService := &Service{
		cfg:           setting.NewCfg(),
		log:           log.New("accesscontrol-test"),
		registrations: accesscontrol.RegistrationList{},
		store:         store,
		roles:         accesscontrol.BuildBasicRoleDefinitions(),
	}

	// Prepare users and permissions
	action := "resources:action"
	permissions := make([]accesscontrol.Permission, 0, resourceCount*usersCount)
	roles := make([]accesscontrol.Role, 0, usersCount)
	orgUsers := make([]org.OrgUser, 0, usersCount)
	users := make([]user.User, 0, usersCount)
	userRoles := make([]accesscontrol.UserRole, 0, usersCount)
	for u := 1; u < usersCount+1; u++ {
		users = append(users, user.User{
			ID:      int64(u),
			Name:    fmt.Sprintf("user%v", u),
			Login:   fmt.Sprintf("user%v", u),
			Email:   fmt.Sprintf("user%v@example.org", u),
			Created: now,
			Updated: now,
		})
		orgUsers = append(orgUsers, org.OrgUser{
			ID:      int64(u),
			UserID:  int64(u),
			OrgID:   1,
			Role:    org.RoleViewer,
			Created: now,
			Updated: now,
		})
		roles = append(roles, accesscontrol.Role{
			ID:      int64(u),
			UID:     fmt.Sprintf("managed_users_%v_permissions", u),
			Name:    fmt.Sprintf("managed:users:%v:permissions", u),
			Version: 1,
			Created: now,
			Updated: now,
		})
		userRoles = append(userRoles, accesscontrol.UserRole{
			ID:      int64(u),
			OrgID:   1,
			RoleID:  int64(u),
			UserID:  int64(u),
			Created: now,
		})

		for r := 1; r < resourceCount+1; r++ {
			permissions = append(permissions, accesscontrol.Permission{
				RoleID:  int64(u),
				Action:  action,
				Scope:   fmt.Sprintf("resources:id:%v", r),
				Created: now,
				Updated: now,
			})
		}
	}

	// Populate store
	if err := batch(len(roles), batchSize, func(start, end int) error {
		err := sqlStore.WithDbSession(context.Background(), func(sess *db.Session) error {
			if _, err := sess.Insert(users[start:end]); err != nil {
				return err
			}
			if _, err := sess.Insert(orgUsers[start:end]); err != nil {
				return err
			}
			if _, err := sess.Insert(roles[start:end]); err != nil {
				return err
			}
			_, err := sess.Insert(userRoles[start:end])
			return err
		})
		return err
	}); err != nil {
		require.NoError(b, err, "could not insert users and roles")
		return nil, nil
	}
	if err := batch(len(permissions), batchSize, func(start, end int) error {
		err := sqlStore.WithDbSession(context.Background(), func(sess *db.Session) error {
			_, err := sess.Insert(permissions[start:end])
			return err
		})
		return err
	}); err != nil {
		require.NoError(b, err, "could not insert permissions")
		return nil, nil
	}

	// Allow signed in user to view all users permissions in the worst way
	userPermissions := map[string][]string{}
	for u := 1; u < usersCount+1; u++ {
		userPermissions[accesscontrol.ActionUsersPermissionsRead] =
			append(userPermissions[accesscontrol.ActionUsersPermissionsRead], fmt.Sprintf("users:id:%v", u))
	}
	return acService, &user.SignedInUser{OrgID: 1, Permissions: map[int64]map[string][]string{1: userPermissions}}
}

func benchGetUsersPermissions(b *testing.B, usersCount, resourceCount int) {
	acService, siu := setupBenchEnv(b, usersCount, resourceCount)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		usersPermissions, err := acService.GetUsersPermissions(context.Background(), siu, 1, "resources:")
		require.NoError(b, err)
		require.Len(b, usersPermissions, usersCount)
		for _, permissions := range usersPermissions {
			require.Len(b, permissions, resourceCount)
		}
	}
}

// Lots of resources
func BenchmarkGetUsersPermissions_10_1K(b *testing.B)   { benchGetUsersPermissions(b, 10, 1000) }   // ~0.047s/op
func BenchmarkGetUsersPermissions_10_10K(b *testing.B)  { benchGetUsersPermissions(b, 10, 10000) }  // ~0.5s/op
func BenchmarkGetUsersPermissions_10_100K(b *testing.B) { benchGetUsersPermissions(b, 10, 100000) } // ~4.6s/op
func BenchmarkGetUsersPermissions_10_1M(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}
	benchGetUsersPermissions(b, 10, 1000000)
} // ~55.36s/op

// Lots of users (most probable case)
func BenchmarkGetUsersPermissions_1K_10(b *testing.B)   { benchGetUsersPermissions(b, 1000, 10) }   // ~0.056s/op
func BenchmarkGetUsersPermissions_10K_10(b *testing.B)  { benchGetUsersPermissions(b, 10000, 10) }  // ~0.58s/op
func BenchmarkGetUsersPermissions_100K_10(b *testing.B) { benchGetUsersPermissions(b, 100000, 10) } // ~6.21s/op
func BenchmarkGetUsersPermissions_1M_10(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}
	benchGetUsersPermissions(b, 1000000, 10)
} // ~57s/op
