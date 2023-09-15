package identity

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	identitysvc "github.com/romashorodok/infosec/internal/identity"
	identitysecurtysvc "github.com/romashorodok/infosec/internal/identity/security"
	"github.com/romashorodok/infosec/pkg/httputils"
	"github.com/romashorodok/infosec/pkg/openapi"
	"github.com/romashorodok/infosec/pkg/tokenutils"

	"go.uber.org/fx"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest --config=handler.cfg.yaml ./identity.openapi.yaml

type IdentityHandler struct {
	Unimplemented

	identitySvc identitysvc.IdentityService
	securitySvc identitysecurtysvc.SecurityService
}

var _ ServerInterface = (*IdentityHandler)(nil)

func (h *IdentityHandler) IdentityServiceSignIn(w http.ResponseWriter, r *http.Request) {
	var request SignInRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Unable deserialize request body. Error:", err.Error())
		return
	}

	result, err := h.identitySvc.Login(r.Context(), request.Username, request.Password)
	if err != nil {
		switch err {
		case identitysvc.NotFoundUserError:
			httputils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		case identitysvc.InvalidCredentialsError:
			httputils.WriteErrorResponse(w, http.StatusUnprocessableEntity)
		default:
			httputils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	http.SetCookie(w, &result.RefreshToken)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: result.AccessToken,
	})
}

func (h *IdentityHandler) IdentityServiceSignUp(w http.ResponseWriter, r *http.Request) {
	var request SignUpRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Unable deserialize request body. Error:", err.Error())
		return
	}

	result, err := h.identitySvc.Register(r.Context(), request.Username, request.Password)
	if err != nil {
		switch err {
		case identitysvc.UnableStoreUser:
			httputils.WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		default:
			httputils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	http.SetCookie(w, &result.RefreshToken)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: result.AccessToken,
	})
}

func (h IdentityHandler) PublicKeyServicePublicKeyList(w http.ResponseWriter, r *http.Request) {
	plainToken := r.Header.Get("Authorization")

	if plainToken == "" {
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Empty token.")
		return
	}

	keyset, err := h.securitySvc.GetPublicKeys(plainToken)
	if err != nil {
		httputils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(keyset)
}

func (h IdentityHandler) TokenRevocationServiceVerifyTokenRevocation(w http.ResponseWriter, r *http.Request) {
	plainToken := r.Header.Get("Authorization")

	if plainToken == "" {
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Empty token.")
		return
	}

	verified, err := h.securitySvc.VerifyToken(plainToken)
	if err != nil {
		httputils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(verified)
}

func deleteRefreshTokenCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     tokenutils.REFRESH_TOKEN_COOKIE_NAME,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func (h IdentityHandler) TokenServiceExchangeToken(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie(tokenutils.REFRESH_TOKEN_COOKIE_NAME)
	if err != nil {
		deleteRefreshTokenCookie(w)
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Empty token.")
		return
	}

	rawRefreshToken := refreshTokenCookie.Value
	if rawRefreshToken == "" {
		deleteRefreshTokenCookie(w)
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Empty token.")
		return
	}

	result, err := h.identitySvc.ExchangeAccessToken(r.Context(), rawRefreshToken)
	if err != nil {
		deleteRefreshTokenCookie(w)
		httputils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	http.SetCookie(w, &result.RefreshToken)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: result.AccessToken,
	})
}

func (h *IdentityHandler) IdentityServiceSignOut(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie(tokenutils.REFRESH_TOKEN_COOKIE_NAME)
	if err != nil {
		deleteRefreshTokenCookie(w)
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Empty token.")
		return
	}

	rawRefreshToken := refreshTokenCookie.Value
	if rawRefreshToken == "" {
		deleteRefreshTokenCookie(w)
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Empty token.")
		return
	}

	err = h.identitySvc.DeleteRefreshToekn(r.Context(), rawRefreshToken)
	if err != nil {
		httputils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	deleteRefreshTokenCookie(w)
	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("Content-Type", "application/json")
}

type IdentityHandlerParams struct {
	fx.In

	IdentitySvc identitysvc.IdentityService
	SecuritySvc identitysecurtysvc.SecurityService

	Lifecycle     fx.Lifecycle
	Router        *chi.Mux
	FilterOptions openapi3filter.Options
}

func RegisterIdentityHandler(params IdentityHandlerParams) {
	spec, err := GetSwagger()
	spec.Servers = nil
	if err != nil {
		log.Panicf("Uanble get openapi spec. %s", err)
	}

	params.Router.Use(openapi.NewOpenAPIRequestMiddleware(spec, &openapi.Options{
		Options: params.FilterOptions,
	}))

	HandlerFromMux(&IdentityHandler{
		identitySvc: params.IdentitySvc,
		securitySvc: params.SecuritySvc,
	}, params.Router)
}
