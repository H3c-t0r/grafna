package interceptors

import (
	"context"
	"testing"

	apikeygenprefix "github.com/grafana/grafana/pkg/components/apikeygenprefixed"
	"github.com/grafana/grafana/pkg/infra/tracing"
	"github.com/grafana/grafana/pkg/services/apikey"
	grpccontext "github.com/grafana/grafana/pkg/services/grpcserver/context"
	"github.com/grafana/grafana/pkg/services/org"
	"github.com/grafana/grafana/pkg/services/user"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestAuthenticator_Authenticate(t *testing.T) {
	tracer := tracing.InitializeTracerForTest()
	serviceAccountId := int64(1)
	t.Run("accepts service api key with admin role", func(t *testing.T) {
		s := newFakeAPIKey(&apikey.APIKey{
			Id:               1,
			OrgId:            1,
			Key:              "admin-api-key",
			Name:             "Admin API Key",
			ServiceAccountId: &serviceAccountId,
		}, nil)
		a := ProvideAuthenticator(s, &fakeUserService{OrgRole: org.RoleAdmin}, grpccontext.ProvideContextHandler(tracer))
		ctx, err := setupContext()
		require.NoError(t, err)
		_, err = a.Authenticate(ctx)
		require.NoError(t, err)
	})

	t.Run("rejects non-admin role", func(t *testing.T) {
		s := newFakeAPIKey(&apikey.APIKey{
			Id:               1,
			OrgId:            1,
			Key:              "admin-api-key",
			Name:             "Admin API Key",
			ServiceAccountId: &serviceAccountId,
		}, nil)
		a := ProvideAuthenticator(s, &fakeUserService{OrgRole: org.RoleEditor}, grpccontext.ProvideContextHandler(tracer))
		ctx, err := setupContext()
		require.NoError(t, err)
		_, err = a.Authenticate(ctx)
		require.NotNil(t, err)
	})

	t.Run("removes auth header from context", func(t *testing.T) {
		s := newFakeAPIKey(&apikey.APIKey{
			Id:               1,
			OrgId:            1,
			Key:              "admin-api-key",
			Name:             "Admin API Key",
			ServiceAccountId: &serviceAccountId,
		}, nil)
		a := ProvideAuthenticator(s, &fakeUserService{OrgRole: org.RoleAdmin}, grpccontext.ProvideContextHandler(tracer))
		ctx, err := setupContext()
		require.NoError(t, err)
		md, ok := metadata.FromIncomingContext(ctx)
		require.True(t, ok)
		require.NotEmpty(t, md["authorization"])
		ctx, err = a.Authenticate(ctx)
		require.NoError(t, err)
		md, ok = metadata.FromIncomingContext(ctx)
		require.True(t, ok)
		require.Empty(t, md["authorization"])
	})

	t.Run("sets SignInUser", func(t *testing.T) {
		s := newFakeAPIKey(&apikey.APIKey{
			Id:               1,
			OrgId:            1,
			Key:              "admin-api-key",
			Name:             "Admin API Key",
			ServiceAccountId: &serviceAccountId,
		}, nil)
		a := ProvideAuthenticator(s, &fakeUserService{OrgRole: org.RoleAdmin}, grpccontext.ProvideContextHandler(tracer))
		ctx, err := setupContext()
		require.NoError(t, err)
		ctx, err = a.Authenticate(ctx)
		require.NoError(t, err)
		signedInUser := grpccontext.FromContext(ctx).SignedInUser
		require.Equal(t, serviceAccountId, signedInUser.UserID)
	})
}

type fakeAPIKey struct {
	apikey.Service
	key *apikey.APIKey
	err error
}

func newFakeAPIKey(key *apikey.APIKey, err error) *fakeAPIKey {
	return &fakeAPIKey{
		key: key,
		err: err,
	}
}

func (f *fakeAPIKey) GetAPIKeyByHash(ctx context.Context, hash string) (*apikey.APIKey, error) {
	return f.key, f.err
}

type fakeUserService struct {
	user.Service
	OrgRole org.RoleType
}

func (f *fakeUserService) GetSignedInUserWithCacheCtx(ctx context.Context, query *user.GetSignedInUserQuery) (*user.SignedInUser, error) {
	return &user.SignedInUser{
		UserID:  1,
		OrgID:   1,
		OrgRole: f.OrgRole,
	}, nil
}

func setupContext() (context.Context, error) {
	ctx := context.Background()
	key, err := apikeygenprefix.New("sa")
	if err != nil {
		return ctx, err
	}
	md := metadata.New(map[string]string{})
	md["authorization"] = []string{"Bearer " + key.ClientSecret}
	return metadata.NewIncomingContext(ctx, md), nil
}
