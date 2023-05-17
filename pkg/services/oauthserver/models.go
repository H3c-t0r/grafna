package oauthserver

import (
	"context"
	"net/http"

	"github.com/grafana/grafana/pkg/services/accesscontrol"
	"gopkg.in/square/go-jose.v2"
)

const (
	// TmpOrgID is the orgID we use while global service accounts are not supported.
	TmpOrgID int64 = 1
	// NoServiceAccountID is the ID we use for client that have no service account associated.
	NoServiceAccountID int64 = 0

	// List of scopes used to identify the impersonated user.
	ScopeUsersSelf       = "users:self"
	ScopeGlobalUsersSelf = "global.users:self"
	ScopeTeamsSelf       = "teams:self"
)

type OAuth2Service interface {
	SaveExternalService(ctx context.Context, cmd *ExternalServiceRegistration) (*ClientDTO, error)
	GetExternalService(ctx context.Context, id string) (*Client, error)
	HandleTokenRequest(rw http.ResponseWriter, req *http.Request)
	HandleIntrospectionRequest(rw http.ResponseWriter, req *http.Request)
}

//go:generate mockery --name Store --structname MockStore --outpkg oauthtest --filename store_mock.go --output ./oauthtest/

type Store interface {
	RegisterExternalService(ctx context.Context, client *Client) error
	SaveExternalService(ctx context.Context, client *Client) error
	GetExternalService(ctx context.Context, id string) (*Client, error)
	GetExternalServiceByName(ctx context.Context, app string) (*Client, error)

	GetExternalServicePublicKey(ctx context.Context, clientID string) (*jose.JSONWebKey, error)
}

type KeyOption struct {
	// URL       string `json:"url,omitempty"` // TODO allow specifying a URL (to a .jwks file) to fetch the key from
	// PublicPEM contains the Base64 encoded public key in PEM format
	PublicPEM string `json:"public_pem,omitempty"`
	Generate  bool   `json:"generate,omitempty"`
}

type SelfCfg struct {
	Enabled     bool                       `json:"enabled"`
	Permissions []accesscontrol.Permission `json:"permissions,omitempty"`
}
type ImpersonationCfg struct {
	Enabled     bool                       `json:"enabled"`
	Groups      bool                       `json:"groups"`
	Permissions []accesscontrol.Permission `json:"permissions,omitempty"`
}
type ExternalServiceRegistration struct {
	ExternalServiceName string           `json:"name"`
	RedirectURI         *string          `json:"redirectUri,omitempty"`
	Impersonation       ImpersonationCfg `json:"impersonation"`
	Self                SelfCfg          `json:"self"`
	Key                 *KeyOption       `json:"key,omitempty"`
}

const (
	RS256 = "RS256"
	ES256 = "ES256"
)
