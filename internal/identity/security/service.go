package security

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"go.uber.org/fx"

	"github.com/romashorodok/infosec/internal/storage/postgres/privatekey"
	"github.com/romashorodok/infosec/internal/storage/postgres/refreshtoken"
	"github.com/romashorodok/infosec/pkg/tokenutils"
)

type securityService struct {
	refreshTokenRepository *refreshtoken.RefreshTokenRepository
	privateKeyRepository   *privatekey.PrivateKeyRepository
}

var _ SecurityService = (*securityService)(nil)

type SecurityService interface {
	CreateAccessToken(CreateAccessTokenParams) (string, error)
	CreateRefreshToken(CreateRefreshTokenParams) (string, *time.Time, error)

	GetPublicKeys(string) ([]byte, error)
	GetTokenKID(string) (uuid.UUID, error)
	GetPrivateKey(uuid.UUID) (jwsUnprotected string, err error)

	VerifyToken(string) ([]byte, error)
}

type CreateAccessTokenParams struct {
	PrivateKeyJwsMessage string
	Kid                  string
	Username             string
	Claims               []string
	UserID               uuid.UUID
}

func (s *securityService) createUnsecureToken(username string, claims []string, expiresAt time.Time) (jwt.Token, error) {
	return jwt.NewBuilder().
		Issuer("0.0.0.0").
		Audience(claims).
		Subject(username).
		Expiration(expiresAt).
		Build()
}

func (s *securityService) signToken(signKey string, headers jws.Headers, token jwt.Token) ([]byte, error) {
	sign, err := jwk.ParseKey([]byte(signKey))

	if err != nil {
		return nil, fmt.Errorf("unable serialize jws json message as jwk.Token. Error: %s\n", err)
	}

	return jwt.Sign(token, jwt.WithKey(jwa.RS256, sign, jws.WithProtectedHeaders(headers)))
}

func (s *securityService) CreateAccessToken(params CreateAccessTokenParams) (string, error) {
	expiresAt := time.Now().Add(time.Minute * 1)

	token, err := s.createUnsecureToken(params.Username, params.Claims, expiresAt)
	if err != nil {
		return "", fmt.Errorf("Unable create access token. Error: %s", err)
	}

	if err = token.Set("user:id", params.UserID); err != nil {
		return "", fmt.Errorf("Unable set `user:id` claim. Error: %s", err)
	}

	if err = token.Set("token:use", "access_token"); err != nil {
		return "", fmt.Errorf("unable set `token:use` claim. Error: %s", err)
	}

	headers := jws.NewHeaders()
	headers.Set(jws.KeyIDKey, params.Kid)

	signed, err := s.signToken(params.PrivateKeyJwsMessage, headers, token)
	if err != nil {
		return "", fmt.Errorf("Unable sign token. Error: %s", err)
	}

	return string(signed), nil
}

type CreateRefreshTokenParams struct {
	PrivateKeyJwsMessage string
	Kid                  string
	Username             string
	Claims               []string
	UserID               uuid.UUID
}

func (s *securityService) CreateRefreshToken(params CreateRefreshTokenParams) (string, *time.Time, error) {
	expiresAt := time.Now().AddDate(1, 0, 0)

	token, err := s.createUnsecureToken(params.Username, params.Claims, expiresAt)
	if err != nil {
		return "", nil, fmt.Errorf("Unable create refresh token. Error: %s", err)
	}

	if err = token.Set("user:id", params.UserID); err != nil {
		return "", nil, fmt.Errorf("Unable set `user:id` claim. Error: %s", err)
	}

	if err = token.Set("token:use", "refresh_token"); err != nil {
		return "", nil, fmt.Errorf("unable set `token:use` claim. Error: %s", err)
	}

	headers := jws.NewHeaders()
	headers.Set(jws.KeyIDKey, params.Kid)

	signed, err := s.signToken(params.PrivateKeyJwsMessage, headers, token)
	if err != nil {
		return "", nil, fmt.Errorf("Unable sign refresh token. Error: %s", err)
	}

	return string(signed), &expiresAt, nil
}

func (s *securityService) getKeySet(kid string, privateKey jwk.Key) jwk.Set {
	keyset := jwk.NewSet()

	pbkey, _ := jwk.PublicKeyOf(privateKey)
	_ = pbkey.Set(jwk.AlgorithmKey, jwa.RS256)
	_ = pbkey.Set(jwk.KeyIDKey, kid)
	_ = keyset.AddKey(pbkey)

	return keyset
}

func (s *securityService) getPrivateKey(kid uuid.UUID) (jwk.Key, error) {
	privateKeyRecord, err := s.privateKeyRepository.GetPrivateKeyById(kid)
	if err != nil {
		return nil, fmt.Errorf("Unable find private key record. Error: %s", err)
	}

	return tokenutils.JwkPrivateKey(privateKeyRecord.JwsMessage)
}

func (s *securityService) GetPublicKeys(rawToken string) ([]byte, error) {
	plainToken := tokenutils.TrimTokenBearer(rawToken)

	if plainToken == "" {
		return nil, errors.New("empty token")
	}

	tokenJws, err := tokenutils.InsecureJwsMessage(plainToken)
	if err != nil {
		return nil, fmt.Errorf("Unable get cast plain token to jws message. Error: %s", err)
	}

	kid := tokenJws.Signatures()[0].ProtectedHeaders().KeyID()
	kidUUID, err := uuid.Parse(kid)

	privateKey, err := s.getPrivateKey(kidUUID)
	if err != nil {
		return nil, fmt.Errorf("Unable get private key. Error: %s", err)
	}

	keyset := s.getKeySet(kid, privateKey)

	return tokenutils.SerializeKeyset(keyset)
}

func (s *securityService) GetTokenKID(rawToken string) (uuid.UUID, error) {
	plainToken := tokenutils.TrimTokenBearer(rawToken)
	tokenJws, err := tokenutils.InsecureJwsMessage(plainToken)
	if err != nil {
		return uuid.NullUUID{}.UUID, fmt.Errorf("Unable get cast plain token to jws message. Error: %s", err)
	}

	kid := tokenJws.Signatures()[0].ProtectedHeaders().KeyID()
	kidUUID, err := uuid.Parse(kid)
	if err != nil {
		return uuid.NullUUID{}.UUID, fmt.Errorf("Unable parse kid as uuid")
	}

	return kidUUID, nil
}

func (s *securityService) GetPrivateKey(kid uuid.UUID) (string, error) {
	privateKeyRecord, err := s.privateKeyRepository.GetPrivateKeyById(kid)
	if err != nil {
		return "", fmt.Errorf("Unable find private key record. Error: %s", err)
	}

	return privateKeyRecord.JwsMessage, nil
}

func (s *securityService) VerifyToken(rawToken string) ([]byte, error) {
	plainToken := tokenutils.TrimTokenBearer(rawToken)

	if plainToken == "" {
		return nil, errors.New("empty token")
	}

	tokenJws, err := tokenutils.InsecureJwsMessage(plainToken)
	if err != nil {
		return nil, fmt.Errorf("Unable get cast plain token to jws message. Error: %s", err)
	}

	kid := tokenJws.Signatures()[0].ProtectedHeaders().KeyID()
	kidUUID, err := uuid.Parse(kid)
	if err != nil {
		return nil, fmt.Errorf("Unable parse kid as uuid")
	}

	privateKey, err := s.getPrivateKey(kidUUID)
	if err != nil {
		return nil, fmt.Errorf("Unable get private key. Error: %s", err)
	}

	keyset := s.getKeySet(kid, privateKey)

	verifiedPayload, err := tokenutils.VerifyTokenByKeySet(keyset, plainToken)
	if err != nil {
		return nil, fmt.Errorf("Not verified token. Error: %s", err)
	}

	err = tokenutils.ValidateToken(string(verifiedPayload))
	if err != nil {
		return nil, fmt.Errorf("Invalid token. Error: %s", err)
	}

	return verifiedPayload, nil
}

type SecurityServiceParams struct {
	fx.In

	RefreshTokenRepository *refreshtoken.RefreshTokenRepository
	PrivateKeyRepositroy   *privatekey.PrivateKeyRepository
}

func NewSecurityService(params SecurityServiceParams) *securityService {
	return &securityService{
		refreshTokenRepository: params.RefreshTokenRepository,
		privateKeyRepository:   params.PrivateKeyRepositroy,
	}
}
