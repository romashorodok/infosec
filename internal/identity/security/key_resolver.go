package security

import (
	"context"
	"fmt"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/romashorodok/infosec/pkg/auth"
	"go.uber.org/fx"
)

type InternalServicePublicKeyResolver struct {
	securitySvc SecurityService
}

var _ auth.IdentityPublicKeyResolver = (*InternalServicePublicKeyResolver)(nil)

func (r *InternalServicePublicKeyResolver) GetKeys(context context.Context, token string) (jwk.Set, error) {
	rawKeyset, err := r.securitySvc.GetPublicKeys(token)
	if err != nil {
		return nil, err
	}

	keyset, err := jwk.Parse(rawKeyset)
	if err != nil {
		return nil, fmt.Errorf("Unable parse keyset. Error: %s", err)
	}

	return keyset, nil
}

type InternalServicePublicKeyResolverParams struct {
	fx.In

	SecuritySvc SecurityService
}

func NewInternalServicePublicKeyResolver(params InternalServicePublicKeyResolverParams) *InternalServicePublicKeyResolver {
	return &InternalServicePublicKeyResolver{
		securitySvc: params.SecuritySvc,
	}
}
