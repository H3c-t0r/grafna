package signingkeysimpl

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"

	"github.com/go-jose/go-jose/v3"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/services/signingkeys"
)

const (
	serverPrivateKeyID = "default"
)

var _ signingkeys.Service = new(EmbeddedKeyService)

func ProvideEmbeddedKeyService() *EmbeddedKeyService {
	s := &EmbeddedKeyService{
		log: log.New("auth.key_service"),
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		s.log.Error("Error generating private key", "err", err)
	}

	s.keys = map[string]crypto.Signer{}
	s.keys[serverPrivateKeyID] = privateKey

	return s
}

type EmbeddedKeyService struct {
	log  log.Logger
	keys map[string]crypto.Signer
}

func (s *EmbeddedKeyService) GetJWKS() jose.JSONWebKeySet {
	result := jose.JSONWebKeySet{}

	for keyID := range s.keys {
		// Skip error check because keyID must be a valid key ID
		jwk, _ := s.GetJWK(keyID)
		result.Keys = append(result.Keys, jwk)
	}

	return result
}

func (s *EmbeddedKeyService) GetJWK(keyID string) (jose.JSONWebKey, error) {
	privateKey, ok := s.keys[keyID]
	if !ok {
		s.log.Error("The specified key was not found", "keyID", keyID)
		return jose.JSONWebKey{}, signingkeys.ErrSigningKeyNotFound.Errorf("The specified key was not found: %s", keyID)
	}

	result := jose.JSONWebKey{
		Key: privateKey.Public(),
		Use: "sig",
	}

	return result, nil
}

func (s *EmbeddedKeyService) GetPublicKey(keyID string) (crypto.PublicKey, error) {
	privateKey, ok := s.keys[keyID]
	if !ok {
		s.log.Error("The specified key was not found", "keyID", keyID)
		return nil, signingkeys.ErrSigningKeyNotFound.Errorf("The specified key was not found: %s", keyID)
	}

	return privateKey.Public(), nil
}

func (s *EmbeddedKeyService) GetPrivateKey(keyID string) (crypto.PrivateKey, error) {
	privateKey, ok := s.keys[keyID]
	if !ok {
		s.log.Error("The specified key was not found", "keyID", keyID)
		return nil, signingkeys.ErrSigningKeyNotFound.Errorf("The specified key was not found: %s", keyID)
	}

	return privateKey, nil
}

func (s *EmbeddedKeyService) GetServerPrivateKey() (crypto.PrivateKey, error) {
	return s.GetPrivateKey(serverPrivateKeyID)
}
