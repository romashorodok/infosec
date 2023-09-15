package identity

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/romashorodok/infosec/internal/identity/security"
	"github.com/romashorodok/infosec/internal/storage/postgres/privatekey"
	"github.com/romashorodok/infosec/internal/storage/postgres/refreshtoken"
	"github.com/romashorodok/infosec/internal/storage/postgres/user"
	"github.com/romashorodok/infosec/pkg/tokenutils"
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
)

var (
	DEFAULT_CLAIMS = []string{"can:read", "can:write"}
)

type identityService struct {
	refreshTokneRepository *refreshtoken.RefreshTokenRepository
	privateKeyRepository   *privatekey.PrivateKeyRepository

	userRepository  *user.UserRepository
	securityService security.SecurityService
	db              *sql.DB
}

var _ IdentityService = (*identityService)(nil)

type IdentityService interface {
	Register(ctx context.Context, username, password string) (*TokenResult, error)
	Login(ctx context.Context, username, password string) (*TokenResult, error)
	ExchangeAccessToken(ctx context.Context, rawRefreshToken string) (*TokenResult, error)
	DeleteRefreshToekn(ctx context.Context, rawRefreshToken string) error
}

func (s *identityService) userGenerateAccessToken(privateKeyJws, kid string, user userInfo) (string, error) {
	accessToken, err := s.securityService.CreateAccessToken(security.CreateAccessTokenParams{
		PrivateKeyJwsMessage: privateKeyJws,
		Kid:                  kid,
		Username:             user.username,
		UserID:               user.userID,
		Claims:               user.claims,
	})
	if err != nil {
		return "", fmt.Errorf("Unable generate access token. Error: %s", err)
	}

	return accessToken, nil
}

func (s *identityService) userGenerateRefreshToken(tx *sql.Tx, ctx context.Context, privateKeyJws string, kid uuid.UUID, user userInfo) (string, *time.Time, error) {

	refreshToken, expiresAt, err := s.securityService.CreateRefreshToken(security.CreateRefreshTokenParams{
		PrivateKeyJwsMessage: privateKeyJws,
		Kid:                  kid.String(),
		Username:             user.username,
		UserID:               user.userID,
		Claims:               user.claims,
	})

	if err != nil {
		return "", nil, fmt.Errorf("Unable generate refresh token. Error: %s", err)
	}

	refreshTokenStored, err := s.refreshTokneRepository.InsertRefreshToken(tx, ctx, kid, refreshToken, *expiresAt)
	if err != nil {
		return "", nil, fmt.Errorf("Unable store refresh token. Error: %s", err)
	}

	if err = s.userRepository.AttachRefreshToken(tx, ctx, user.userID, refreshTokenStored.ID); err != nil {
		return "", nil, fmt.Errorf("Unable store associated refresh token with user. Error.  Error: %s", err)
	}

	return refreshTokenStored.Plaintext, expiresAt, nil
}

func (s *identityService) userGenerateTokens(tx *sql.Tx, ctx context.Context, user userInfo) (*TokenResult, error) {

	privateKeyJws, err := tokenutils.CreatePrivateKeyAsJwsMessage()
	if err != nil {
		return nil, fmt.Errorf("Unable create private key. Error: %s", err)
	}

	privateKey, err := s.privateKeyRepository.InsertPrivateKey(tx, ctx, privateKeyJws)
	if err != nil {
		return nil, fmt.Errorf("Unable save private key record. Error: %s", err)
	}

	if err = s.userRepository.AttachPrivateKey(tx, ctx, user.userID, privateKey.ID); err != nil {
		return nil, fmt.Errorf("Unable store associated private key with user. Error: %s", err)
	}

	kid := privateKey.ID.String()

	refreshToken, expiresAt, err := s.userGenerateRefreshToken(tx, ctx, privateKey.JwsMessage, privateKey.ID, user)
	if err != nil {
		return nil, fmt.Errorf("Unable generate refresh token. Error: %s", err)
	}

	accessToken, err := s.userGenerateAccessToken(privateKey.JwsMessage, kid, user)
	if err != nil {
		return nil, err
	}

	return &TokenResult{
		AccessToken: accessToken,
		RefreshToken: http.Cookie{
			Name:     tokenutils.REFRESH_TOKEN_COOKIE_NAME,
			Value:    refreshToken,
			Expires:  *expiresAt,
			HttpOnly: true,
			Path:     "/",
		},
	}, nil
}

func (s *identityService) Register(ctx context.Context, username, password string) (*TokenResult, error) {

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Printf("%s. Err: %s", UnableHashPassword, err)
		return nil, UnableHashPassword
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("%s. Err: %s", TransactionFailedError, err)
		return nil, TransactionFailedError
	}
	defer tx.Rollback()

	user, err := s.userRepository.InsertUser(tx, ctx, username, string(hashed))
	if err != nil {
		log.Printf("%s. Err: %s", UnableStoreUser, err)
		return nil, UnableStoreUser
	}

	result, err := s.userGenerateTokens(tx, ctx, userInfo{
		userID:   user.ID,
		username: user.Username,
		claims:   DEFAULT_CLAIMS,
	})
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *identityService) Login(ctx context.Context, username, password string) (*TokenResult, error) {
	user, err := s.userRepository.FindUserByUsername(username)
	if err != nil {
		log.Printf("%s. Error: %s", NotFoundUserError, err)
		return nil, NotFoundUserError
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("%s. Error: %s", InvalidCredentialsError, err)
		return nil, InvalidCredentialsError
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("%s. Error: %s", TransactionFailedError, err)
		return nil, TransactionFailedError
	}
	defer tx.Rollback()

	result, err := s.userGenerateTokens(tx, ctx, userInfo{
		userID:   user.ID,
		username: user.Username,
		claims:   DEFAULT_CLAIMS,
	})
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}
func (s *identityService) userRefreshTokens(tx *sql.Tx, ctx context.Context, kid uuid.UUID) (*TokenResult, error) {

	user, err := s.userRepository.FindUserByPrivateKey(kid)
	if err != nil {
		return nil, fmt.Errorf("Unable find user by private key. Error: %s", err)
	}

	pkey, err := s.securityService.GetPrivateKey(kid)
	if err != nil {
		return nil, fmt.Errorf("Unable find private key. Error: %s", err)
	}

	accessToken, err := s.securityService.CreateAccessToken(security.CreateAccessTokenParams{
		PrivateKeyJwsMessage: pkey,
		Kid:                  kid.String(),
		Username:             user.Username,
		UserID:               user.ID,
		Claims:               DEFAULT_CLAIMS,
	})
	if err != nil {
		return nil, fmt.Errorf("Unable create access token. Error: %s", err)
	}

	if err = s.refreshTokneRepository.DeleteRefreshTokenByPrivateKey(tx, ctx, kid); err != nil {
		return nil, fmt.Errorf("Unable delete old access token. Error: %s", err)
	}

	refreshToke, expiresAt, err := s.userGenerateRefreshToken(
		tx,
		ctx,
		pkey,
		kid,
		userInfo{
			userID:   user.ID,
			username: user.Username,
			claims:   DEFAULT_CLAIMS,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Unable refresh token. Error: %s", err)
	}

	return &TokenResult{
		AccessToken: accessToken,
		RefreshToken: http.Cookie{
			Name:     tokenutils.REFRESH_TOKEN_COOKIE_NAME,
			Value:    refreshToke,
			Expires:  *expiresAt,
			HttpOnly: true,
			Path:     "/",
		},
	}, nil
}

func (s *identityService) ExchangeAccessToken(ctx context.Context, rawRefreshToken string) (*TokenResult, error) {
	_, err := s.securityService.VerifyToken(rawRefreshToken)
	if err != nil {
		return nil, err
	}

	kid, err := s.securityService.GetTokenKID(rawRefreshToken)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("Unable start transaction. Error: %s", err)
	}
	defer tx.Rollback()

	tokens, err := s.userRefreshTokens(tx, ctx, kid)
	if err != nil {
		return nil, fmt.Errorf("Unable refresh tokens. Error: %s", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *identityService) DeleteRefreshToekn(ctx context.Context, rawRefreshToken string) error {
	kid, err := s.securityService.GetTokenKID(rawRefreshToken)
	if err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Unable start transaction. Error: %s", err)
	}
	defer tx.Rollback()

	err = s.refreshTokneRepository.DeleteRefreshTokenByPrivateKey(tx, ctx, kid)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("Unable delete refresh token. Error: %s", err)
	}

	err = s.privateKeyRepository.DeletePrivateKey(tx, ctx, kid)
	if err != nil {
		return fmt.Errorf("Unable delete private key. Error: %s", err)
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

type IdentityServiceParams struct {
	fx.In

	RefreshTokenRepository *refreshtoken.RefreshTokenRepository
	PrivateKeyRepositroy   *privatekey.PrivateKeyRepository
	UserRepository         *user.UserRepository
	SecuritySvc            security.SecurityService
	DB                     *sql.DB
}

func NewIdentityService(params IdentityServiceParams) *identityService {
	return &identityService{
		refreshTokneRepository: params.RefreshTokenRepository,
		privateKeyRepository:   params.PrivateKeyRepositroy,
		userRepository:         params.UserRepository,
		securityService:        params.SecuritySvc,
		db:                     params.DB,
	}
}
