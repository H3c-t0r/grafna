package login

import (
	"context"
	"crypto/subtle"

	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/sqlstore"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/util"
)

var validatePassword = func(providedPassword string, userPassword string, userSalt string) error {
	passwordHashed, err := util.EncodePassword(providedPassword, userSalt)
	if err != nil {
		return err
	}
	if subtle.ConstantTimeCompare([]byte(passwordHashed), []byte(userPassword)) != 1 {
		return ErrInvalidCredentials
	}

	return nil
}

var loginUsingGrafanaDB = func(ctx context.Context, query *models.LoginUserQuery, store sqlstore.Store) (bool, error) {
	if setting.DisableLogin {
		return false, nil
	}

	userQuery := models.GetUserByLoginQuery{LoginOrEmail: query.Username}

	if err := store.GetUserByLogin(ctx, &userQuery); err != nil {
		return true, err
	}

	user := userQuery.Result

	if user.IsDisabled {
		return true, ErrUserDisabled
	}

	if err := validatePassword(query.Password, user.Password, user.Salt); err != nil {
		return true, err
	}

	query.User = user
	return true, nil
}
