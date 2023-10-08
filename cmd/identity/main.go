package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/romashorodok/infosec/ent"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
	v1httpidentity "github.com/romashorodok/infosec/internal/http/v1/identity"
	"github.com/romashorodok/infosec/internal/identity"
	"github.com/romashorodok/infosec/internal/identity/security"
	"github.com/romashorodok/infosec/internal/storage/postgres/privatekey"
	"github.com/romashorodok/infosec/internal/storage/postgres/refreshtoken"
	"github.com/romashorodok/infosec/internal/storage/postgres/user"
	"github.com/romashorodok/infosec/pkg/auth"
	"github.com/romashorodok/infosec/pkg/envutils"
	"go.uber.org/fx"

	_ "github.com/lib/pq"
)

const (
	HTTP_HOST_DEFAULT = "0.0.0.0"
	HTTP_PORT_DEFAULT = "8080"

	DATABASE_HOST_DEFAULT     = "0.0.0.0"
	DATABASE_PORT_DEFAULT     = "5432"
	DATABASE_USERNAME_DEFAULT = "user"
	DATABASE_PASSWORD_DEFAULT = "password"
	DATABASE_NAME_DEFAULT     = "postgres"
)

const (
	HTTP_HOST_VAR = "HTTP_HOST"
	HTTP_PORT_VAR = "HTTP_PORT"

	DATABASE_HOST_VAR     = "DATABASE_HOST"
	DATABASE_PORT_VAR     = "DATABASE_PORT"
	DATABASE_USERNAME_VAR = "DATABASE_USERNAME"
	DATABASE_PASSWORD_VAR = "DATABASE_PASSWORD"
	DATABASE_NAME_VAR     = "DATABASE_NAME"
)

type HTTPConfig struct {
	Port string
	Host string
}

func NewHTTPConfig() *HTTPConfig {
	return &HTTPConfig{
		Port: envutils.Env(HTTP_PORT_VAR, HTTP_PORT_DEFAULT),
		Host: envutils.Env(HTTP_HOST_VAR, HTTP_HOST_DEFAULT),
	}
}

var router = chi.NewRouter()

func NewRouter() *chi.Mux {
	return router
}

type DatabaseConfig struct {
	Username string
	Password string
	Database string
	Host     string
	Port     string
	Driver   string
}

func (dconf *DatabaseConfig) GetURI() string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		dconf.Driver,
		dconf.Username,
		dconf.Password,
		dconf.Host,
		dconf.Port,
		dconf.Database,
	)
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Driver:   "postgres",
		Username: envutils.Env(DATABASE_USERNAME_VAR, DATABASE_USERNAME_DEFAULT),
		Password: envutils.Env(DATABASE_PASSWORD_VAR, DATABASE_PASSWORD_DEFAULT),
		Host:     envutils.Env(DATABASE_HOST_VAR, DATABASE_HOST_DEFAULT),
		Port:     envutils.Env(DATABASE_PORT_VAR, DATABASE_PORT_DEFAULT),
		Database: envutils.Env(DATABASE_NAME_VAR, DATABASE_NAME_DEFAULT),
	}
}

type DatabaseConnectionParams struct {
	fx.In

	Dconf     *DatabaseConfig
	Lifecycle fx.Lifecycle
}

func NewDatabaseConnection(params DatabaseConnectionParams) *sql.DB {
	uri := params.Dconf.GetURI()

	db, err := sql.Open(params.Dconf.Driver, uri+"?sslmode=disable&connect_timeout=5")

	if err != nil {
		log.Panicf("Unable connect to database %s. Error: %s \n", uri, err)
	}

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(context.Context) error {
			db.Close()
			return nil
		},
	})

	return db
}

type OpenAPI3FilterOptions struct {
	fx.In

	IdentityPublicKeyResolver auth.IdentityPublicKeyResolver
}

func NewOpenAPI3FilterOptions(params OpenAPI3FilterOptions) openapi3filter.Options {
	return openapi3filter.Options{
		AuthenticationFunc: auth.NewAsymmetricEncryptionAuthenticator(params.IdentityPublicKeyResolver),
		MultiError:         true,
	}
}

func startServer(lifecycle fx.Lifecycle, config *HTTPConfig, handler *chi.Mux) {
	server := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: handler,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			ln, err := net.Listen("tcp", server.Addr)
			if err != nil {
				return err
			}
			go server.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}

type EntClientParams struct {
	fx.In

	Dconf     *DatabaseConfig
	Lifecycle fx.Lifecycle
}

func NewEntClient(params EntClientParams) (*ent.Client, error) {
	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		params.Dconf.Host,
		params.Dconf.Port,
		params.Dconf.Username,
		params.Dconf.Database,
		params.Dconf.Password,
	)
	client, err := ent.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	params.Lifecycle.Append(fx.StopHook(func() {
		client.Close()
	}))

	params.Lifecycle.Append(fx.StartHook(func(ctx context.Context) {
		if err := client.Schema.Create(ctx); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}

		// id, _ := uuid.FromString("bc3a2ad3-cc1c-4f46-8be4-62c6dc3872ff")
		// user, _ := client.User.Get(ctx, id)

		// result, err := user.QueryParticipants().First(ctx)

		// log.Printf("%+v", user)

	}))

	return client, nil
}

func main() {
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))


	fx.New(
		fx.Provide(
			fx.Annotate(
				security.NewInternalServicePublicKeyResolver,
				fx.As(new(auth.IdentityPublicKeyResolver)),
			),
			NewOpenAPI3FilterOptions,
			NewDatabaseConfig,
			NewDatabaseConnection,
			NewHTTPConfig,
			NewRouter,
			NewEntClient,

			user.NewUserRepository,
			privatekey.NewPrivateKeyRepositroy,
			refreshtoken.NewRefreshTokenRepository,

			fx.Annotate(
				identity.NewIdentityService,
				fx.As(new(identity.IdentityService)),
			),
			fx.Annotate(
				security.NewSecurityService,
				fx.As(new(security.SecurityService)),
			),
		),
		fx.Invoke(v1httpidentity.RegisterIdentityHandler),
		fx.Invoke(startServer),
		fx.Invoke(func(*ent.Client) {}),
	).Run()
}
