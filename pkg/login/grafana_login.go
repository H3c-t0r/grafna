package login

import (
	"context"
	"crypto/subtle"

	"github.com/grafana/grafana/pkg/bus"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/util"
)

var validatePassword = func(providedPassword string, userPassword string, userSalt string) error {
	passwordHashed := util.EncodePassword(providedPassword, userSalt)
	if subtle.ConstantTimeCompare([]byte(passwordHashed), []byte(userPassword)) != 1 {
		return ErrInvalidCredentials
	}

	return nil
}

var loginUsingGrafanaDB = func(query *m.LoginUserQuery) error {
	userQuery := m.GetUserByLoginQuery{LoginOrEmail: query.Username}

	if err := bus.DispatchCtx(context.TODO(), &userQuery); err != nil {
		return err
	}

	user := userQuery.Result

	if user.IsDisabled {
		return ErrUserDisabled
	}

	if err := validatePassword(query.Password, user.Password, user.Salt); err != nil {
		return err
	}

	query.User = user
	return nil
}
