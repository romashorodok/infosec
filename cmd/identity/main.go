package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/romashorodok/infosec/ent"
	"github.com/romashorodok/infosec/ent/migrate"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
	v1httpidentity "github.com/romashorodok/infosec/internal/http/v1/identity"
	v1httpkanban "github.com/romashorodok/infosec/internal/http/v1/kanban"
	"github.com/romashorodok/infosec/internal/identity"
	"github.com/romashorodok/infosec/internal/identity/security"
	"github.com/romashorodok/infosec/internal/kanban"
	"github.com/romashorodok/infosec/internal/storage/postgres/privatekey"
	"github.com/romashorodok/infosec/internal/storage/postgres/refreshtoken"
	"github.com/romashorodok/infosec/internal/storage/postgres/user"
	"github.com/romashorodok/infosec/pkg/auth"
	"github.com/romashorodok/infosec/pkg/envutils"
	"github.com/romashorodok/infosec/pkg/httputils"
	"github.com/romashorodok/infosec/pkg/openapi"
	"github.com/romashorodok/infosec/pkg/openapi3utils"
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

func (h *HTTPConfig) GetAddr() string {
	return net.JoinHostPort(h.Host, h.Port)
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

type NewSpecOptionsHandlerConstructorParams struct {
	fx.In

	Options openapi3filter.Options
}

func NewSpecOptionsHandlerConstructor(params NewSpecOptionsHandlerConstructorParams) openapi3utils.HandlerSpecValidator {
	return func(spec *openapi3utils.Spec) openapi3utils.HandlerFunc {
		return openapi.NewOpenAPIRequestMiddleware(spec, &openapi.Options{
			Options: params.Options,
		})
	}
}

type startServerParams struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    *HTTPConfig
	Handler   *chi.Mux
	Handlers  []httputils.HttpHandler `group:"http.handler"`
}

func startServer(params startServerParams) {
	server := &http.Server{
		Addr:    params.Config.GetAddr(),
		Handler: params.Handler,
	}

	for _, handler := range params.Handlers {
		handler.GetOption()(server.Handler)
	}

	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		panic(err)
	}

	go server.Serve(ln)

	params.Lifecycle.Append(
		fx.StopHook(func(ctx context.Context) error {
			return server.Shutdown(ctx)
		}),
	)
}

type EntClientParams struct {
	fx.In

	DB        *sql.DB
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
		if err := client.Schema.Create(ctx,
			migrate.WithDropColumn(false),
			migrate.WithDropIndex(false),
		); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}

		_, _ = params.DB.ExecContext(ctx, "ALTER TABLE users ALTER COLUMN id SET DEFAULT uuid_generate_v4();")

		// TODO: Why go fmt not work for that ?
		_, _ = params.DB.ExecContext(ctx, `
CREATE OR REPLACE FUNCTION delete_related_participants()
RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM participants
    WHERE id = OLD.participant_id;
    RETURN OLD;
END;

$$ LANGUAGE plpgsql;

CREATE TRIGGER after_board_participants_delete
AFTER DELETE ON board_participants
FOR EACH ROW
EXECUTE FUNCTION delete_related_participants();
`)
	}))

	return client, nil
}

func NewLogger() *slog.Logger {
	var level slog.LevelVar
	level.Set(slog.LevelDebug)

	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	return slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, logFile), &slog.HandlerOptions{
		AddSource: true,
		Level:     &level,
	}))
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
			httputils.AsHttpHandler(v1httpkanban.NewHandler),
			httputils.AsHttpHandler(v1httpidentity.NewHandler),
		),

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

			NewLogger,
			NewSpecOptionsHandlerConstructor,

			kanban.NewKanbanService,

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
		fx.Invoke(startServer),
		fx.Invoke(func(*ent.Client) {}),
	).Run()
}
