package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/romashorodok/infosec/pkg/tokenutils"
)

type SecurityError struct {
	Err error
}

func GetSecurityErrorPrefix() string {
	return ": Security Error"
}

func (s *SecurityError) Error() string {
	return fmt.Sprintf("%s %s", s.Err.Error(), GetSecurityErrorPrefix())
}

func NewSecurityError(err error) *SecurityError {
	return &SecurityError{Err: err}
}

func authenticator(ctx context.Context, input *openapi3filter.AuthenticationInput, resolver IdentityPublicKeyResolver) error {
	if input.SecuritySchemeName != "BearerAuth" {
		return fmt.Errorf("security scheme %s != 'BearerAuth'", input.SecuritySchemeName)
	}

	request := input.RequestValidationInput.Request
	rawToken := request.Header.Get("Authorization")
	if rawToken == "" {
		return errors.New("empty bearer token")
	}

	plainToken := tokenutils.TrimTokenBearer(rawToken)

	keyset, err := resolver.GetKeys(ctx, plainToken)
	if err != nil {
		return fmt.Errorf("Not found keyset for token. Error: %s", err)
	}

	verifiedPayload, err := tokenutils.VerifyTokenByKeySet(keyset, plainToken)
	if err != nil {
		return fmt.Errorf("Not verified token. Error: %s", err)
	}

	err = tokenutils.ValidateToken(string(verifiedPayload))
	if err != nil {
		return fmt.Errorf("Invalid token. Error: %s", err)
	}

	authContext := context.WithValue(
		request.Context(),
		TOKEN_CONTEXT_VALUE,
		verifiedPayload,
	)

	newReq := request.Clone(authContext)

	input.RequestValidationInput.Request = newReq.WithContext(authContext)

	return nil
}

type IdentityPublicKeyResolver interface {
	GetKeys(context.Context, string) (jwk.Set, error)
}

func NewAsymmetricEncryptionAuthenticator(resolver IdentityPublicKeyResolver) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		if err := authenticator(ctx, input, resolver); err != nil {
			return NewSecurityError(err)
		}
		return nil
	}
}
