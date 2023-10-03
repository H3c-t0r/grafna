package signingkeysimpl

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v3"

	"github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/infra/remotecache"
	"github.com/grafana/grafana/pkg/services/secrets"
	"github.com/grafana/grafana/pkg/services/signingkeys"
	"github.com/grafana/grafana/pkg/services/signingkeys/signingkeystore"
)

var _ signingkeys.Service = new(Service)

func ProvideEmbeddedSigningKeysService(dbStore db.DB, secretsService secrets.Service,
	remoteCache remotecache.CacheStorage,
) (*Service, error) {
	s := &Service{
		log:         log.New("auth.key_service"),
		store:       signingkeystore.NewSigningKeyStore(dbStore, secretsService),
		remoteCache: remoteCache,
	}

	return s, nil
}

// Service provides functionality for managing signing keys used to sign and verify JWT tokens for
// the OSS version of Grafana.
//
// The service is under active development and is not yet ready for production use.
type Service struct {
	log         log.Logger
	store       signingkeystore.SigningStore
	remoteCache remotecache.CacheStorage
}

// GetJWKS returns the JSON Web Key Set (JWKS) with all the keys that can be used to verify tokens (public keys)
func (s *Service) GetJWKS(ctx context.Context) (jose.JSONWebKeySet, error) {
	jwks, err := s.store.GetJWKS(ctx)
	return jwks, err
}

// GetPublicKey returns the public key with the specified key ID
func (s *Service) GetPublicKey(ctx context.Context, keyID string) (crypto.PublicKey, error) {
	privateKey, err := s.store.GetPrivateKey(ctx, keyID)
	if err != nil {
		s.log.Error("The specified key was not found", "keyID", keyID, "err", err)
		return nil, signingkeys.ErrSigningKeyNotFound.Errorf("The specified key was not found: %s", keyID)
	}

	return privateKey.Public(), nil
}

// GetPrivateKey returns the private key with the specified key ID
func (s *Service) GetPrivateKey(ctx context.Context, keyID string) (crypto.PrivateKey, error) {
	privateKey, err := s.store.GetPrivateKey(ctx, keyID)
	if err != nil {
		s.log.Error("The specified key was not found", "keyID", keyID, "err", err)
		return nil, signingkeys.ErrSigningKeyNotFound.Errorf("The specified key was not found: %s", keyID)
	}

	return privateKey, nil
}

// AddPrivateKey adds a private key to the service
func (s *Service) AddPrivateKey(ctx context.Context, keyID string,
	privateKey crypto.Signer, alg jose.SignatureAlgorithm, expiresAt *time.Time, force bool) error {
	if _, err := s.store.AddPrivateKey(ctx, keyID, alg, privateKey, expiresAt, force); err != nil {
		s.log.Error("Failed to add private key", "keyID", keyID, "err", err)
		return signingkeys.ErrKeyGenerationFailed.Errorf("Failed to add private key: %w", err)
	}

	return nil
}

// GetOrCreatePrivateKey returns the private key with the specified key ID. If the key does not exist, it will be
// created with the specified algorithm.
// The key will be automatically rotated at the beginning of each month. The previous key will be kept for 30 days.
func (s *Service) GetOrCreatePrivateKey(ctx context.Context, keyPrefix string, alg jose.SignatureAlgorithm) (crypto.Signer, error) {
	if alg != jose.ES256 {
		s.log.Error("Only ES256 is supported", "alg", alg)
		return nil, signingkeys.ErrKeyGenerationFailed.Errorf("Only ES256 is supported: %v", alg)
	}

	keyID := KeyMonthScopedID(keyPrefix, alg)
	signer, err := s.store.GetPrivateKey(ctx, keyID)
	if err == nil {
		return signer, nil
	}
	s.log.Debug("Private key not found, generating new key", "keyID", keyID, "err", err)

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		s.log.Error("Error generating private key", "err", err)
		return nil, signingkeys.ErrKeyGenerationFailed.Errorf("Error generating private key: %v", err)
	}

	expiry := time.Now().Add(30 * 24 * time.Hour)
	if signer, err = s.store.AddPrivateKey(ctx, keyID, alg, privateKey, &expiry, false); err != nil && !errors.Is(err, signingkeys.ErrSigningKeyAlreadyExists) {
		return nil, err
	}

	return signer, nil
}

func KeyMonthScopedID(keyPrefix string, alg jose.SignatureAlgorithm) string {
	keyID := keyPrefix + "-" + time.Now().UTC().Format("2006-01") + "-" + strings.ToLower(string(alg))
	return keyID
}
