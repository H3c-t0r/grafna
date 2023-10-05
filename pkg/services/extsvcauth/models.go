package extsvcauth

import (
	"context"

	"github.com/grafana/grafana/pkg/services/accesscontrol"
)

const (
	OAuth2Server AuthProvider = "OAuth2Server"
)

type AuthProvider string

type ExternalServiceRegistry interface {
	// SaveExternalService creates or updates an external service in the database. Based on the requested auth provider,
	// it generates client_id, secrets and any additional provider specificities (ex: rsa keys). It also ensures that the
	// associated service account has the correct permissions.
	SaveExternalService(ctx context.Context, cmd *ExternalServiceRegistration) (*ExternalServiceDTO, error)
}

type SelfCfg struct {
	// Enabled allows the service to request access tokens for itself
	Enabled bool `json:"enabled"`
	// Permissions are the permissions that the external service needs its associated service account to have.
	Permissions []accesscontrol.Permission `json:"permissions,omitempty"`
}

type ImpersonationCfg struct {
	// Enabled allows the service to request access tokens to impersonate users
	Enabled bool `json:"enabled"`
	// Groups allows the service to list the impersonated user's teams
	Groups bool `json:"groups"`
	// Permissions are the permissions that the external service needs when impersonating a user.
	// The intersection of this set with the impersonated user's permission guarantees that the client will not
	// gain more privileges than the impersonated user has and vice versa.
	Permissions []accesscontrol.Permission `json:"permissions,omitempty"`
}

// ExternalServiceRegistration represents the registration form to save new client.
type ExternalServiceRegistration struct {
	Name string `json:"name"`
	// Impersonation access configuration
	// (this is not available on all auth providers)
	Impersonation ImpersonationCfg `json:"impersonation"`
	// Self access configuration
	Self SelfCfg `json:"self"`
	// Auth Provider that the client will use to connect to Grafana
	AuthProvider AuthProvider `json:"authProvider"`
	// Auth Provider specific config
	AuthProviderCfg interface{} `json:"authProviderCfg"`
}

// ExternalServiceDTO represents the credentials that the ExternalService can use to connect to Grafana.
type ExternalServiceDTO struct {
	Name   string      `json:"name"`
	ID     string      `json:"id"`
	Secret string      `json:"secret"`
	Extra  interface{} `json:"extra"` // Auth Provider specificities (ex: ecdsa key pair)
}
