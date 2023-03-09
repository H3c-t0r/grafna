package clients

import (
	"context"
	"errors"
	"net/url"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/infra/network"
	"github.com/grafana/grafana/pkg/services/auth"
	"github.com/grafana/grafana/pkg/services/authn"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/user"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/web"
)

var _ authn.HookClient = new(Session)
var _ authn.ContextAwareClient = new(Session)

func ProvideSession(cfg *setting.Cfg, sessionService auth.UserTokenService, features *featuremgmt.FeatureManager) *Session {
	return &Session{
		cfg:            cfg,
		features:       features,
		sessionService: sessionService,
		log:            log.New(authn.ClientSession),
	}
}

type Session struct {
	cfg            *setting.Cfg
	features       *featuremgmt.FeatureManager
	sessionService auth.UserTokenService
	userService    user.Service
	log            log.Logger
}

func (s *Session) Name() string {
	return authn.ClientSession
}

func (s *Session) Authenticate(ctx context.Context, r *authn.Request) (*authn.Identity, error) {
	unescapedCookie, err := r.HTTPRequest.Cookie(s.cfg.LoginCookieName)
	if err != nil {
		return nil, err
	}

	rawSessionToken, err := url.QueryUnescape(unescapedCookie.Value)
	if err != nil {
		return nil, err
	}

	var token *auth.UserToken

	if s.features.IsEnabled(featuremgmt.FlagFrontendTokenRotation) {
		token, err = s.sessionService.GetToken(ctx, auth.GetTokenQuery{UnHashedToken: rawSessionToken})
	} else {
		token, err = s.sessionService.LookupToken(ctx, rawSessionToken)
	}

	if err != nil {
		s.log.FromContext(ctx).Warn("Failed to look up session from cookie", "error", err)
		return nil, err
	}

	return &authn.Identity{
		ID:           authn.NamespacedID(authn.NamespaceUser, token.UserId),
		SessionToken: token,
		ClientParams: authn.ClientParams{
			FetchSyncedUser: true,
			SyncPermissions: true,
		},
	}, nil
}

func (s *Session) Test(ctx context.Context, r *authn.Request) bool {
	if s.cfg.LoginCookieName == "" {
		return false
	}

	if _, err := r.HTTPRequest.Cookie(s.cfg.LoginCookieName); err != nil {
		return false
	}

	return true
}

func (s *Session) Priority() uint {
	return 60
}

func (s *Session) Hook(ctx context.Context, identity *authn.Identity, r *authn.Request) error {
	if identity.SessionToken == nil || s.features.IsEnabled(featuremgmt.FlagFrontendTokenRotation) {
		return nil
	}

	r.Resp.Before(func(w web.ResponseWriter) {
		if w.Written() || errors.Is(ctx.Err(), context.Canceled) {
			return
		}

		// FIXME (jguer): get real values
		addr := web.RemoteAddr(r.HTTPRequest)
		userAgent := r.HTTPRequest.UserAgent()

		// addr := reqContext.RemoteAddr()
		ip, err := network.GetIPFromAddress(addr)
		if err != nil {
			s.log.Debug("failed to get client IP address", "addr", addr, "err", err)
			ip = nil
		}
		rotated, newToken, err := s.sessionService.TryRotateToken(ctx, identity.SessionToken, ip, userAgent)
		if err != nil {
			s.log.Error("failed to rotate token", "error", err)
			return
		}

		if rotated {
			identity.SessionToken = newToken
			s.log.Debug("rotated session token", "user", identity.ID)

			authn.WriteSessionCookie(w, s.cfg, identity)
		}
	})

	return nil
}
