package authinfoservice

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/grafana/grafana/pkg/infra/usagestats"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/login/authinfoservice/database"
	secretstore "github.com/grafana/grafana/pkg/services/secrets/database"
	secretsManager "github.com/grafana/grafana/pkg/services/secrets/manager"
	"github.com/grafana/grafana/pkg/services/sqlstore"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

//nolint:goconst
func TestUserAuth(t *testing.T) {
	sqlStore := sqlstore.InitTestDB(t)
	secretsService := secretsManager.SetupTestService(t, secretstore.ProvideSecretsStore(sqlStore))
	authInfoStore := database.ProvideAuthInfoStore(sqlStore, secretsService)
	srv := ProvideAuthInfoService(
		&OSSUserProtectionImpl{},
		authInfoStore,
		&usagestats.UsageStatsMock{},
	)

	t.Run("Given 5 users", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			cmd := models.CreateUserCommand{
				Email: fmt.Sprint("user", i, "@test.com"),
				Name:  fmt.Sprint("user", i),
				Login: fmt.Sprint("loginuser", i),
			}
			_, err := sqlStore.CreateUser(context.Background(), cmd)
			require.Nil(t, err)
		}

		t.Run("Can find existing user", func(t *testing.T) {
			// By Login
			login := "loginuser0"

			query := &models.GetUserByAuthInfoQuery{UserLookupParams: models.UserLookupParams{Login: &login}}
			user, err := srv.LookupAndUpdate(context.Background(), query)

			require.Nil(t, err)
			require.Equal(t, user.Login, login)

			// By ID
			id := user.Id

			user, err = srv.LookupByOneOf(context.Background(), &models.UserLookupParams{
				UserID: &id,
			})

			require.Nil(t, err)
			require.Equal(t, user.Id, id)

			// By Email
			email := "user1@test.com"

			user, err = srv.LookupByOneOf(context.Background(), &models.UserLookupParams{
				Email: &email,
			})

			require.Nil(t, err)
			require.Equal(t, user.Email, email)

			// Don't find nonexistent user
			email = "nonexistent@test.com"

			user, err = srv.LookupByOneOf(context.Background(), &models.UserLookupParams{
				Email: &email,
			})

			require.Equal(t, models.ErrUserNotFound, err)
			require.Nil(t, user)
		})

		t.Run("Can set & locate by AuthModule and AuthId", func(t *testing.T) {
			// get nonexistent user_auth entry
			query := &models.GetUserByAuthInfoQuery{AuthModule: "test", AuthId: "test"}
			user, err := srv.LookupAndUpdate(context.Background(), query)

			require.Equal(t, models.ErrUserNotFound, err)
			require.Nil(t, user)

			// create user_auth entry
			login := "loginuser0"

			query.UserLookupParams.Login = &login
			user, err = srv.LookupAndUpdate(context.Background(), query)

			require.Nil(t, err)
			require.Equal(t, user.Login, login)

			// get via user_auth
			query = &models.GetUserByAuthInfoQuery{AuthModule: "test", AuthId: "test"}
			user, err = srv.LookupAndUpdate(context.Background(), query)

			require.Nil(t, err)
			require.Equal(t, user.Login, login)

			// get with non-matching id
			idPlusOne := user.Id + 1

			query.UserLookupParams.UserID = &idPlusOne
			user, err = srv.LookupAndUpdate(context.Background(), query)

			require.Nil(t, err)
			require.Equal(t, user.Login, "loginuser1")

			// get via user_auth
			query = &models.GetUserByAuthInfoQuery{AuthModule: "test", AuthId: "test"}
			user, err = srv.LookupAndUpdate(context.Background(), query)

			require.Nil(t, err)
			require.Equal(t, user.Login, "loginuser1")

			// remove user
			err = sqlStore.WithDbSession(context.Background(), func(sess *sqlstore.DBSession) error {
				_, err := sess.Exec("DELETE FROM "+sqlStore.Dialect.Quote("user")+" WHERE id=?", user.Id)
				return err
			})
			require.NoError(t, err)

			// get via user_auth for deleted user
			query = &models.GetUserByAuthInfoQuery{AuthModule: "test", AuthId: "test"}
			user, err = srv.LookupAndUpdate(context.Background(), query)

			require.Equal(t, err, models.ErrUserNotFound)
			require.Nil(t, user)
		})

		t.Run("Can set & retrieve oauth token information", func(t *testing.T) {
			token := &oauth2.Token{
				AccessToken:  "testaccess",
				RefreshToken: "testrefresh",
				Expiry:       time.Now(),
				TokenType:    "Bearer",
			}
			idToken := "testidtoken"
			token = token.WithExtra(map[string]interface{}{"id_token": idToken})

			// Find a user to set tokens on
			login := "loginuser0"

			// Calling GetUserByAuthInfoQuery on an existing user will populate an entry in the user_auth table
			query := &models.GetUserByAuthInfoQuery{AuthModule: "test", AuthId: "test", UserLookupParams: models.UserLookupParams{
				Login: &login,
			}}
			user, err := srv.LookupAndUpdate(context.Background(), query)

			require.Nil(t, err)
			require.Equal(t, user.Login, login)

			cmd := &models.UpdateAuthInfoCommand{
				UserId:     user.Id,
				AuthId:     query.AuthId,
				AuthModule: query.AuthModule,
				OAuthToken: token,
			}
			err = srv.authInfoStore.UpdateAuthInfo(context.Background(), cmd)

			require.Nil(t, err)

			getAuthQuery := &models.GetAuthInfoQuery{
				UserId: user.Id,
			}

			err = srv.authInfoStore.GetAuthInfo(context.Background(), getAuthQuery)

			require.Nil(t, err)
			require.Equal(t, token.AccessToken, getAuthQuery.Result.OAuthAccessToken)
			require.Equal(t, token.RefreshToken, getAuthQuery.Result.OAuthRefreshToken)
			require.Equal(t, token.TokenType, getAuthQuery.Result.OAuthTokenType)
			require.Equal(t, idToken, getAuthQuery.Result.OAuthIdToken)
		})

		t.Run("Always return the most recently used auth_module", func(t *testing.T) {
			// Restore after destructive operation
			sqlStore = sqlstore.InitTestDB(t)

			for i := 0; i < 5; i++ {
				cmd := models.CreateUserCommand{
					Email: fmt.Sprint("user", i, "@test.com"),
					Name:  fmt.Sprint("user", i),
					Login: fmt.Sprint("loginuser", i),
				}
				_, err := sqlStore.CreateUser(context.Background(), cmd)
				require.Nil(t, err)
			}

			// Find a user to set tokens on
			login := "loginuser0"

			// Calling srv.LookupAndUpdateQuery on an existing user will populate an entry in the user_auth table
			// Make the first log-in during the past
			database.GetTime = func() time.Time { return time.Now().AddDate(0, 0, -2) }
			query := &models.GetUserByAuthInfoQuery{AuthModule: "test1", AuthId: "test1", UserLookupParams: models.UserLookupParams{
				Login: &login,
			}}
			user, err := srv.LookupAndUpdate(context.Background(), query)
			database.GetTime = time.Now

			require.Nil(t, err)
			require.Equal(t, user.Login, login)

			// Add a second auth module for this user
			// Have this module's last log-in be more recent
			database.GetTime = func() time.Time { return time.Now().AddDate(0, 0, -1) }
			query = &models.GetUserByAuthInfoQuery{AuthModule: "test2", AuthId: "test2", UserLookupParams: models.UserLookupParams{
				Login: &login,
			}}
			user, err = srv.LookupAndUpdate(context.Background(), query)
			database.GetTime = time.Now

			require.Nil(t, err)
			require.Equal(t, user.Login, login)

			// Get the latest entry by not supply an authmodule or authid
			getAuthQuery := &models.GetAuthInfoQuery{
				UserId: user.Id,
			}

			err = authInfoStore.GetAuthInfo(context.Background(), getAuthQuery)

			require.Nil(t, err)
			require.Equal(t, getAuthQuery.Result.AuthModule, "test2")

			// "log in" again with the first auth module
			updateAuthCmd := &models.UpdateAuthInfoCommand{UserId: user.Id, AuthModule: "test1", AuthId: "test1"}
			err = authInfoStore.UpdateAuthInfo(context.Background(), updateAuthCmd)

			require.Nil(t, err)

			// Get the latest entry by not supply an authmodule or authid
			getAuthQuery = &models.GetAuthInfoQuery{
				UserId: user.Id,
			}

			err = authInfoStore.GetAuthInfo(context.Background(), getAuthQuery)

			require.Nil(t, err)
			require.Equal(t, getAuthQuery.Result.AuthModule, "test1")
		})

		t.Run("Keeps track of last used auth_module when not using oauth", func(t *testing.T) {
			// Restore after destructive operation
			sqlStore = sqlstore.InitTestDB(t)

			for i := 0; i < 5; i++ {
				cmd := models.CreateUserCommand{
					Email: fmt.Sprint("user", i, "@test.com"),
					Name:  fmt.Sprint("user", i),
					Login: fmt.Sprint("loginuser", i),
				}
				_, err := sqlStore.CreateUser(context.Background(), cmd)
				require.Nil(t, err)
			}

			// Find a user to set tokens on
			login := "loginuser0"

			fixedTime := time.Now()
			// Calling srv.LookupAndUpdateQuery on an existing user will populate an entry in the user_auth table
			// Make the first log-in during the past
			database.GetTime = func() time.Time { return fixedTime.AddDate(0, 0, -2) }
			queryOne := &models.GetUserByAuthInfoQuery{AuthModule: "test1", AuthId: "test1", UserLookupParams: models.UserLookupParams{
				Login: &login,
			}}
			user, err := srv.LookupAndUpdate(context.Background(), queryOne)
			database.GetTime = time.Now

			require.Nil(t, err)
			require.Equal(t, user.Login, login)

			// Add a second auth module for this user
			// Have this module's last log-in be more recent
			database.GetTime = func() time.Time { return fixedTime.AddDate(0, 0, -1) }
			queryTwo := &models.GetUserByAuthInfoQuery{AuthModule: "test2", AuthId: "test2", UserLookupParams: models.UserLookupParams{
				Login: &login,
			}}
			user, err = srv.LookupAndUpdate(context.Background(), queryTwo)
			require.Nil(t, err)
			require.Equal(t, user.Login, login)

			// Get the latest entry by not supply an authmodule or authid
			getAuthQuery := &models.GetAuthInfoQuery{
				UserId: user.Id,
			}

			err = authInfoStore.GetAuthInfo(context.Background(), getAuthQuery)

			require.Nil(t, err)
			require.Equal(t, "test2", getAuthQuery.Result.AuthModule)

			// Now reuse first auth module and make sure it's updated to the most recent
			database.GetTime = func() time.Time { return fixedTime }

			// add oauth info to auth_info to make sure update date does not overwrite it
			updateAuthCmd := &models.UpdateAuthInfoCommand{UserId: user.Id, AuthModule: "test1", AuthId: "test1", OAuthToken: &oauth2.Token{
				AccessToken:  "access_token",
				TokenType:    "token_type",
				RefreshToken: "refresh_token",
				Expiry:       fixedTime,
			}}
			err = authInfoStore.UpdateAuthInfo(context.Background(), updateAuthCmd)
			require.Nil(t, err)
			user, err = srv.LookupAndUpdate(context.Background(), queryOne)

			require.Nil(t, err)
			require.Equal(t, user.Login, login)

			err = authInfoStore.GetAuthInfo(context.Background(), getAuthQuery)

			require.Nil(t, err)
			require.Equal(t, "test1", getAuthQuery.Result.AuthModule)
			// make sure oauth info is not overwritten by update date
			require.Equal(t, "access_token", getAuthQuery.Result.OAuthAccessToken)

			// Now reuse second auth module and make sure it's updated to the most recent
			database.GetTime = func() time.Time { return fixedTime.AddDate(0, 0, 1) }
			user, err = srv.LookupAndUpdate(context.Background(), queryTwo)
			require.Nil(t, err)
			require.Equal(t, user.Login, login)

			err = authInfoStore.GetAuthInfo(context.Background(), getAuthQuery)
			require.Nil(t, err)
			require.Equal(t, "test2", getAuthQuery.Result.AuthModule)

			// Ensure test 1 did not have its entry modified
			getAuthQueryUnchanged := &models.GetAuthInfoQuery{
				UserId:     user.Id,
				AuthModule: "test1",
			}
			err = authInfoStore.GetAuthInfo(context.Background(), getAuthQueryUnchanged)
			require.Nil(t, err)
			require.Equal(t, "test1", getAuthQueryUnchanged.Result.AuthModule)
			require.Less(t, getAuthQueryUnchanged.Result.Created, getAuthQuery.Result.Created)
		})

		t.Run("Can set & locate by generic oauth auth module and user id", func(t *testing.T) {
			// Find a user to set tokens on
			login := "loginuser0"

			// Expect to pass since there's a matching login user
			database.GetTime = func() time.Time { return time.Now().AddDate(0, 0, -2) }
			query := &models.GetUserByAuthInfoQuery{AuthModule: genericOAuthModule, AuthId: "", UserLookupParams: models.UserLookupParams{
				Login: &login,
			}}
			user, err := srv.LookupAndUpdate(context.Background(), query)
			database.GetTime = time.Now

			require.Nil(t, err)
			require.Equal(t, user.Login, login)

			otherLoginUser := "aloginuser"
			// Should throw a "user not found" error since there's no matching login user
			database.GetTime = func() time.Time { return time.Now().AddDate(0, 0, -2) }
			query = &models.GetUserByAuthInfoQuery{AuthModule: genericOAuthModule, AuthId: "", UserLookupParams: models.UserLookupParams{
				Login: &otherLoginUser,
			}}
			user, err = srv.LookupAndUpdate(context.Background(), query)
			database.GetTime = time.Now

			require.NotNil(t, err)
			require.Nil(t, user)
		})

		t.Run("calculate metrics on duplicate userstats", func(t *testing.T) {
			// Restore after destructive operation
			sqlStore = sqlstore.InitTestDB(t)

			for i := 0; i < 5; i++ {
				cmd := models.CreateUserCommand{
					Email: fmt.Sprint("user", i, "@test.com"),
					Name:  fmt.Sprint("user", i),
					Login: fmt.Sprint("loginuser", i),
					OrgId: 1,
				}
				_, err := sqlStore.CreateUser(context.Background(), cmd)
				require.Nil(t, err)
			}

			// "Skipping duplicate users test for mysql as it does make unique constraint case insensitive by default
			if sqlStore.GetDialect().DriverName() != "mysql" {
				dupUserEmailcmd := models.CreateUserCommand{
					Email: "USERDUPLICATETEST1@TEST.COM",
					Name:  "user name 1",
					Login: "USER_DUPLICATE_TEST_1_LOGIN",
				}
				_, err := sqlStore.CreateUser(context.Background(), dupUserEmailcmd)
				require.NoError(t, err)

				// add additional user with duplicate login where DOMAIN is upper case
				dupUserLogincmd := models.CreateUserCommand{
					Email: "userduplicatetest1@test.com",
					Name:  "user name 1",
					Login: "user_duplicate_test_1_login",
				}
				_, err = sqlStore.CreateUser(context.Background(), dupUserLogincmd)
				require.NoError(t, err)
				// require metrics and statistics to be 2
				m, err := srv.authInfoStore.CollectLoginStats(context.Background())
				require.NoError(t, err)
				require.Equal(t, 2, m["stats.users.duplicate_user_entries"])
				require.Equal(t, 1, m["stats.users.has_duplicate_user_entries"])
			}
		})
	})
}
