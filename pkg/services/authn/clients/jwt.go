package clients

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/authn"
	"github.com/grafana/grafana/pkg/services/org"
	"github.com/grafana/grafana/pkg/services/user"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/util/errutil"
	"github.com/jmespath/go-jmespath"
)

var _ authn.Client = new(JWT)

var (
	ErrJWTInvalid = errutil.NewBase(errutil.StatusUnauthorized,
		"jwt.invalid", errutil.WithPublicMessage("Failed to verify JWT"))
	ErrJWTMissingClaim = errutil.NewBase(errutil.StatusUnauthorized,
		"jwt.missing_claim", errutil.WithPublicMessage("Missing mandatory claim in JWT"))
	ErrJWTInvalidRole = errutil.NewBase(errutil.StatusForbidden,
		"jwt.invalid_role", errutil.WithPublicMessage("Invalid Role in claim"))
)

func ProvideJWT(userService user.Service, jwtService models.JWTService, cfg *setting.Cfg) *JWT {
	return &JWT{
		log:         log.New(authn.ClientAPIKey),
		jwtService:  jwtService,
		userService: userService,
	}
}

type JWT struct {
	cfg         *setting.Cfg
	log         log.Logger
	jwtService  models.JWTService
	userService user.Service
}

func (s *JWT) Authenticate(ctx context.Context, r *authn.Request) (*authn.Identity, error) {
	jwtToken := r.HTTPRequest.Header.Get(s.cfg.JWTAuthHeaderName)
	if jwtToken == "" && s.cfg.JWTAuthURLLogin {
		jwtToken = r.HTTPRequest.URL.Query().Get("auth_token")
	}

	// Strip the 'Bearer' prefix if it exists.
	jwtToken = strings.TrimPrefix(jwtToken, "Bearer ")

	claims, err := s.jwtService.Verify(ctx, jwtToken)
	if err != nil {
		s.log.Debug("Failed to verify JWT", "error", err)
		return nil, ErrJWTInvalid
	}

	sub, _ := claims["sub"].(string)
	if sub == "" {
		s.log.Warn("Got a JWT without the mandatory 'sub' claim", "error", err)
		return nil, ErrJWTMissingClaim.Errorf("missing mandatory 'sub' claim in JWT")
	}

	id := &authn.Identity{
		AuthModule: "jwt",
		AuthID:     sub,
		OrgRoles:   map[int64]org.RoleType{},
		ClientParams: authn.ClientParams{
			SyncUser:            true,
			SyncTeamMembers:     true,
			AllowSignUp:         false,
			EnableDisabledUsers: false,
		}}

	if key := s.cfg.JWTAuthUsernameClaim; key != "" {
		id.Login, _ = claims[key].(string)
		id.ClientParams.LookUpParams.Login = &id.Login
	}
	if key := s.cfg.JWTAuthEmailClaim; key != "" {
		id.Email, _ = claims[key].(string)
		id.ClientParams.LookUpParams.Email = &id.Email
	}

	if name, _ := claims["name"].(string); name != "" {
		id.Name = name
	}

	role, grafanaAdmin := s.extractRoleAndAdmin(claims)
	if s.cfg.JWTAuthRoleAttributeStrict && !role.IsValid() {
		s.log.Warn("extracted Role is invalid", "role", role, "auth_id", id.AuthID)
		return nil, ErrJWTInvalidRole.Errorf("invalid role claim in JWT: %s", role)
	}

	if role.IsValid() {
		var orgID int64
		// FIXME (jguer): GetIDForNewUser already has the auto assign information
		// just neeeds the org role. Find a meaningful way to pass this default
		// role to it (that doesn't involve id.OrgRoles[0] = role)
		if s.cfg.AutoAssignOrg && s.cfg.AutoAssignOrgId > 0 {
			orgID = int64(s.cfg.AutoAssignOrgId)
			s.log.Debug("The user has a role assignment and organization membership is auto-assigned",
				"role", role, "orgId", orgID)
		} else {
			orgID = int64(1)
			s.log.Debug("The user has a role assignment and organization membership is not auto-assigned",
				"role", role, "orgId", orgID)
		}

		id.OrgRoles[orgID] = role
		if s.cfg.JWTAuthAllowAssignGrafanaAdmin {
			id.IsGrafanaAdmin = &grafanaAdmin
		}
	}

	if id.Login == "" && id.Email == "" {
		s.log.Debug("Failed to get an authentication claim from JWT",
			"login", id.Login, "email", id.Email)
		return nil, ErrJWTMissingClaim
	}

	if s.cfg.JWTAuthAutoSignUp {
		id.ClientParams.AllowSignUp = true
	}

	return nil, nil
}

func (s *JWT) Test(ctx context.Context, r *authn.Request) bool {
	if !s.cfg.JWTAuthEnabled || s.cfg.JWTAuthHeaderName == "" {
		return false
	}

	jwtToken := r.HTTPRequest.Header.Get(s.cfg.JWTAuthHeaderName)
	if jwtToken == "" && s.cfg.JWTAuthURLLogin {
		jwtToken = r.HTTPRequest.URL.Query().Get("auth_token")
	}

	if jwtToken == "" {
		return false
	}

	// Strip the 'Bearer' prefix if it exists.
	jwtToken = strings.TrimPrefix(jwtToken, "Bearer ")

	// The header is Authorization and the token does not look like a JWT,
	// this is likely an API key. Pass it on.
	if s.cfg.JWTAuthHeaderName == "Authorization" && !looksLikeJWT(jwtToken) {
		return false
	}

	return true
}

func looksLikeJWT(token string) bool {
	// A JWT must have 3 parts separated by `.`.
	parts := strings.Split(token, ".")
	return len(parts) == 3
}

const roleGrafanaAdmin = "GrafanaAdmin"

func (s *JWT) extractRoleAndAdmin(claims map[string]interface{}) (org.RoleType, bool) {
	if s.cfg.JWTAuthRoleAttributePath == "" {
		return "", false
	}

	role, err := searchClaimsForStringAttr(s.cfg.JWTAuthRoleAttributePath, claims)
	if err != nil || role == "" {
		return "", false
	}

	if role == roleGrafanaAdmin {
		return org.RoleAdmin, true
	}
	return org.RoleType(role), false
}

func searchClaimsForStringAttr(attributePath string, claims map[string]interface{}) (string, error) {
	val, err := searchClaimsForAttr(attributePath, claims)
	if err != nil {
		return "", err
	}

	strVal, ok := val.(string)
	if ok {
		return strVal, nil
	}

	return "", nil
}

func searchClaimsForAttr(attributePath string, claims map[string]interface{}) (interface{}, error) {
	if attributePath == "" {
		return "", errors.New("no attribute path specified")
	}

	if len(claims) == 0 {
		return "", errors.New("empty claims provided")
	}

	val, err := jmespath.Search(attributePath, claims)
	if err != nil {
		return "", fmt.Errorf("failed to search claims with provided path: %q: %w", attributePath, err)
	}

	return val, nil
}
