package oasimpl

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"

	ac "github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/auth/identity"
	"github.com/grafana/grafana/pkg/services/extsvcauth/oauthserver"
	"github.com/grafana/grafana/pkg/services/extsvcauth/oauthserver/utils"
	"github.com/grafana/grafana/pkg/services/team"
	"github.com/grafana/grafana/pkg/services/user"
)

// HandleTokenRequest handles the client's OAuth2 query to obtain an access_token by presenting its authorization
// grant (ex: client_credentials, jwtbearer)
func (s *OAuth2ServiceImpl) HandleTokenRequest(rw http.ResponseWriter, req *http.Request) {
	// This context will be passed to all methods.
	ctx := req.Context()

	// Create an empty session object which will be passed to the request handlers
	oauthSession := NewAuthSession()

	// This will create an access request object and iterate through the registered TokenEndpointHandlers to validate the request.
	accessRequest, err := s.oauthProvider.NewAccessRequest(ctx, req, oauthSession)
	if err != nil {
		s.writeAccessError(ctx, rw, accessRequest, err)
		return
	}

	client, err := s.GetExternalService(ctx, accessRequest.GetClient().GetID())
	if err != nil || client == nil {
		s.oauthProvider.WriteAccessError(ctx, rw, accessRequest, &fosite.RFC6749Error{
			DescriptionField: "Could not find the requested subject.",
			ErrorField:       "not_found",
			CodeField:        http.StatusBadRequest,
		})
		return
	}
	oauthSession.JWTClaims.Add("client_id", client.ClientID)

	errClientCred := s.handleClientCredentials(ctx, accessRequest, oauthSession, client)
	if errClientCred != nil {
		s.writeAccessError(ctx, rw, accessRequest, errClientCred)
		return
	}

	errJWTBearer := s.handleJWTBearer(ctx, accessRequest, oauthSession, client)
	if errJWTBearer != nil {
		s.writeAccessError(ctx, rw, accessRequest, errJWTBearer)
		return
	}

	// All tokens we generate in this service should target Grafana's API.
	accessRequest.GrantAudience(s.cfg.AppURL)

	// Prepare response, fosite handlers will populate the token.
	response, err := s.oauthProvider.NewAccessResponse(ctx, accessRequest)
	if err != nil {
		s.writeAccessError(ctx, rw, accessRequest, err)
		return
	}
	s.oauthProvider.WriteAccessResponse(ctx, rw, accessRequest, response)
}

// writeAccessError logs the error then uses fosite to write the error back to the user.
func (s *OAuth2ServiceImpl) writeAccessError(ctx context.Context, rw http.ResponseWriter, accessRequest fosite.AccessRequester, err error) {
	var fositeErr *fosite.RFC6749Error
	if errors.As(err, &fositeErr) {
		s.logger.Error("Description", fositeErr.DescriptionField, "hint", fositeErr.HintField, "error", fositeErr.ErrorField)
	} else {
		s.logger.Error("Error", err)
	}
	s.oauthProvider.WriteAccessError(ctx, rw, accessRequest, err)
}

// splitOAuthScopes sort scopes that are generic (profile, email, groups, entitlements) from scopes
// that are RBAC actions (used to further restrict the entitlements embedded in the access_token)
func splitOAuthScopes(requestedScopes fosite.Arguments) (map[string]bool, map[string]bool) {
	actionsFilter := map[string]bool{}
	claimsFilter := map[string]bool{}
	for _, scope := range requestedScopes {
		switch scope {
		case "profile", "email", "groups", "entitlements":
			claimsFilter[scope] = true
		default:
			actionsFilter[scope] = true
		}
	}
	return actionsFilter, claimsFilter
}

// handleJWTBearer populates the "impersonation" access_token generated by fosite to match the rfc9068 specifications (entitlements, groups).
// It ensures that the user can be impersonated, that the generated token audiences only contain Grafana's AppURL (and token endpoint)
// and that entitlements solely contain the user's permissions that the client is allowed to have.
func (s *OAuth2ServiceImpl) handleJWTBearer(ctx context.Context, accessRequest fosite.AccessRequester, oauthSession *oauth2.JWTSession, client *oauthserver.OAuthExternalService) error {
	if !accessRequest.GetGrantTypes().ExactOne(string(fosite.GrantTypeJWTBearer)) {
		return nil
	}

	userID, err := utils.ParseUserIDFromSubject(oauthSession.Subject)
	if err != nil {
		return &fosite.RFC6749Error{
			DescriptionField: "Could not find the requested subject.",
			ErrorField:       "not_found",
			CodeField:        http.StatusBadRequest,
		}
	}

	// Check audiences list only contains the AppURL and the token endpoint
	for _, aud := range accessRequest.GetGrantedAudience() {
		if aud != fmt.Sprintf("%voauth2/token", s.cfg.AppURL) && aud != s.cfg.AppURL {
			return &fosite.RFC6749Error{
				DescriptionField: "Client is not allowed to target this Audience.",
				HintField:        "The audience must be the AppURL or the token endpoint.",
				ErrorField:       "invalid_request",
				CodeField:        http.StatusForbidden,
			}
		}
	}

	// If the client was not allowed to impersonate the user we would not have reached this point given allowed scopes would have been empty
	// But just in case we check again
	ev := ac.EvalPermission(ac.ActionUsersImpersonate, ac.Scope("users", "id", strconv.FormatInt(userID, 10)))
	hasAccess, errAccess := s.accessControl.Evaluate(ctx, client.SignedInUser, ev)
	if errAccess != nil || !hasAccess {
		return &fosite.RFC6749Error{
			DescriptionField: "Client is not allowed to impersonate subject.",
			ErrorField:       "restricted_access",
			CodeField:        http.StatusForbidden,
		}
	}

	// Populate claims' suject from the session subject
	oauthSession.JWTClaims.Subject = oauthSession.Subject

	// Get the user
	query := user.GetUserByIDQuery{ID: userID}
	dbUser, err := s.userService.GetByID(ctx, &query)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return &fosite.RFC6749Error{
				DescriptionField: "Could not find the requested subject.",
				ErrorField:       "not_found",
				CodeField:        http.StatusBadRequest,
			}
		}
		return &fosite.RFC6749Error{
			DescriptionField: "The request subject could not be processed.",
			ErrorField:       "server_error",
			CodeField:        http.StatusInternalServerError,
		}
	}
	oauthSession.Username = dbUser.Login

	// Split scopes into actions and claims
	actionsFilter, claimsFilter := splitOAuthScopes(accessRequest.GetGrantedScopes())

	teams := []*team.TeamDTO{}
	// Fetch teams if the groups scope is requested or if we need to populate it in the entitlements
	if claimsFilter["groups"] ||
		(claimsFilter["entitlements"] && (len(actionsFilter) == 0 || actionsFilter["teams:read"])) {
		var errGetTeams error
		teams, errGetTeams = s.teamService.GetTeamsByUser(ctx, &team.GetTeamsByUserQuery{
			OrgID:  oauthserver.TmpOrgID,
			UserID: dbUser.ID,
			// Fetch teams without restriction on permissions
			SignedInUser: &user.SignedInUser{
				OrgID: oauthserver.TmpOrgID,
				Permissions: map[int64]map[string][]string{
					oauthserver.TmpOrgID: {
						ac.ActionTeamsRead: {ac.ScopeTeamsAll},
					},
				},
			},
		})
		if errGetTeams != nil {
			return &fosite.RFC6749Error{
				DescriptionField: "The teams scope could not be processed.",
				ErrorField:       "server_error",
				CodeField:        http.StatusInternalServerError,
			}
		}
	}
	if claimsFilter["profile"] {
		oauthSession.JWTClaims.Add("name", dbUser.Name)
		oauthSession.JWTClaims.Add("login", dbUser.Login)
		oauthSession.JWTClaims.Add("updated_at", dbUser.Updated.Unix())
	}
	if claimsFilter["email"] {
		oauthSession.JWTClaims.Add("email", dbUser.Email)
	}
	if claimsFilter["groups"] {
		teamNames := make([]string, 0, len(teams))
		for _, team := range teams {
			teamNames = append(teamNames, team.Name)
		}
		oauthSession.JWTClaims.Add("groups", teamNames)
	}

	if claimsFilter["entitlements"] {
		// Get the user permissions (apply the actions filter)
		permissions, errGetPermission := s.filteredUserPermissions(ctx, userID, actionsFilter)
		if errGetPermission != nil {
			return errGetPermission
		}

		// Compute the impersonated permissions (apply the actions filter, replace the scope self with the user id)
		impPerms := s.filteredImpersonatePermissions(client.ImpersonatePermissions, userID, teams, actionsFilter)

		// Intersect the permissions with the client permissions
		intesect := ac.Intersect(permissions, impPerms)

		oauthSession.JWTClaims.Add("entitlements", intesect)
	}

	return nil
}

// filteredUserPermissions gets the user permissions and applies the actions filter
func (s *OAuth2ServiceImpl) filteredUserPermissions(ctx context.Context, userID int64, actionsFilter map[string]bool) ([]ac.Permission, error) {
	permissions, err := s.acService.SearchUserPermissions(ctx, oauthserver.TmpOrgID,
		ac.SearchOptions{NamespacedID: fmt.Sprintf("%s:%d", identity.NamespaceUser, userID)})
	if err != nil {
		return nil, &fosite.RFC6749Error{
			DescriptionField: "The permissions scope could not be processed.",
			ErrorField:       "server_error",
			CodeField:        http.StatusInternalServerError,
		}
	}

	// Apply the actions filter
	if len(actionsFilter) > 0 {
		filtered := []ac.Permission{}
		for i := range permissions {
			if actionsFilter[permissions[i].Action] {
				filtered = append(filtered, permissions[i])
			}
		}
		permissions = filtered
	}
	return permissions, nil
}

// filteredImpersonatePermissions computes the impersonated permissions.
// It applies the actions filter and replaces the "self RBAC scopes" (~ scope templates) by the correct user id/team id.
func (*OAuth2ServiceImpl) filteredImpersonatePermissions(impersonatePermissions []ac.Permission, userID int64, teams []*team.TeamDTO, actionsFilter map[string]bool) []ac.Permission {
	// Compute the impersonated permissions
	impPerms := impersonatePermissions
	// Apply the actions filter
	if len(actionsFilter) > 0 {
		filtered := []ac.Permission{}
		for i := range impPerms {
			if actionsFilter[impPerms[i].Action] {
				filtered = append(filtered, impPerms[i])
			}
		}
		impPerms = filtered
	}

	// Replace the scope self with the user id
	correctScopes := []ac.Permission{}
	for i := range impPerms {
		switch impPerms[i].Scope {
		case oauthserver.ScopeGlobalUsersSelf:
			correctScopes = append(correctScopes, ac.Permission{
				Action: impPerms[i].Action,
				Scope:  ac.Scope("global.users", "id", strconv.FormatInt(userID, 10)),
			})
		case oauthserver.ScopeUsersSelf:
			correctScopes = append(correctScopes, ac.Permission{
				Action: impPerms[i].Action,
				Scope:  ac.Scope("users", "id", strconv.FormatInt(userID, 10)),
			})
		case oauthserver.ScopeTeamsSelf:
			for t := range teams {
				correctScopes = append(correctScopes, ac.Permission{
					Action: impPerms[i].Action,
					Scope:  ac.Scope("teams", "id", strconv.FormatInt(teams[t].ID, 10)),
				})
			}
		default:
			correctScopes = append(correctScopes, impPerms[i])
		}
		continue
	}
	return correctScopes
}

// handleClientCredentials populates the client's access_token generated by fosite to match the rfc9068 specifications (entitlements, groups)
func (s *OAuth2ServiceImpl) handleClientCredentials(ctx context.Context, accessRequest fosite.AccessRequester, oauthSession *oauth2.JWTSession, client *oauthserver.OAuthExternalService) error {
	if !accessRequest.GetGrantTypes().ExactOne("client_credentials") {
		return nil
	}
	// Set the subject to the service account associated to the client
	oauthSession.JWTClaims.Subject = fmt.Sprintf("user:id:%d", client.ServiceAccountID)

	sa := client.SignedInUser
	if sa == nil {
		return &fosite.RFC6749Error{
			DescriptionField: "Could not find the service account of the client",
			ErrorField:       "not_found",
			CodeField:        http.StatusNotFound,
		}
	}
	oauthSession.Username = sa.Login

	// For client credentials, scopes are not marked as granted by fosite but the request would have been rejected
	// already if the client was not allowed to request them
	for _, scope := range accessRequest.GetRequestedScopes() {
		accessRequest.GrantScope(scope)
	}

	// Split scopes into actions and claims
	actionsFilter, claimsFilter := splitOAuthScopes(accessRequest.GetGrantedScopes())

	if claimsFilter["profile"] {
		oauthSession.JWTClaims.Add("name", sa.Name)
		oauthSession.JWTClaims.Add("login", sa.Login)
	}
	if claimsFilter["email"] {
		s.logger.Debug("Service accounts have no emails")
	}
	if claimsFilter["groups"] {
		s.logger.Debug("Service accounts have no groups")
	}
	if claimsFilter["entitlements"] {
		s.logger.Debug("Processing client entitlements")
		if sa.Permissions != nil && sa.Permissions[oauthserver.TmpOrgID] != nil {
			perms := sa.Permissions[oauthserver.TmpOrgID]
			if len(actionsFilter) > 0 {
				filtered := map[string][]string{}
				for action := range actionsFilter {
					if _, ok := perms[action]; ok {
						filtered[action] = perms[action]
					}
				}
				perms = filtered
			}
			oauthSession.JWTClaims.Add("entitlements", perms)
		} else {
			s.logger.Debug("Client has no permissions")
		}
	}

	return nil
}
