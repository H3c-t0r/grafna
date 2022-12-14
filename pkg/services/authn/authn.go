package authn

import (
	"context"
	"net/http"

	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/org"
	"github.com/grafana/grafana/pkg/services/user"
	"golang.org/x/oauth2"
)

const (
	ClientAnonymous = "auth.anonymous"
)

type ClientParams struct {
	SyncUser    bool
	AllowSignUp bool
}

type PostAuthHookFn func(ctx context.Context, clientParams *ClientParams, identity *Identity) error

type Service interface {
	// RegisterPostAuthHook registers a hook that is called after a successful authentication.
	RegisterPostAuthHook(hook PostAuthHookFn)
	// Authenticate authenticates a request using the specified client.
	Authenticate(ctx context.Context, client string, r *Request) (*Identity, error)
}

type Client interface {
	Authenticate(ctx context.Context, r *Request) (*Identity, error)
}

type Request struct {
	HTTPRequest *http.Request
}

type Identity struct {
	OrgID       int64
	OrgName     string
	IsAnonymous bool
	OrgRoles    map[int64]org.RoleType

	Login        string
	Name         string
	Email        string
	AuthModule   string // AuthModule is the name of the external system
	AuthID       string // AuthId is the unique identifier for the user in the external system
	OAuthToken   *oauth2.Token
	LookUpParams models.UserLookupParams
}

func (i *Identity) Role() org.RoleType {
	return i.OrgRoles[i.OrgID]
}

func (i *Identity) SignedInUser() *user.SignedInUser {
	return &user.SignedInUser{
		OrgID:       i.OrgID,
		OrgName:     i.OrgName,
		OrgRole:     i.Role(),
		IsAnonymous: i.IsAnonymous,
	}
}
