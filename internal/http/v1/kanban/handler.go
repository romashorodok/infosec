package kanban

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/romashorodok/infosec/pkg/httputils"
	"github.com/romashorodok/infosec/pkg/openapi3utils"
	"go.uber.org/fx"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest --config=handler.cfg.yaml ./../../../../ent/openapi.json

type handler struct {
	Unimplemented

	handlerSpecValidator openapi3utils.HandlerSpecValidator
}

var _ ServerInterface = (*handler)(nil)
var _ httputils.HttpHandler = (*handler)(nil)

func (*handler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	var req CreateBoardJSONRequestBody
	_ = req

	log.Println("test")
}

func (h *handler) GetOption() httputils.HttpHandlerOption {
	return func(hand http.Handler) {
		switch hand.(type) {
		case *chi.Mux:
			mux := hand.(*chi.Mux)
			spec, err := GetSwagger()
			if err != nil {
				log.Panicf("unable get openapi spec for streamchannels.handler.Err: %s", err)
			}
			spec.Servers = nil

			HandlerWithOptions(h, ChiServerOptions{
				BaseRouter: mux,
				Middlewares: []MiddlewareFunc{
					h.handlerSpecValidator(spec),
					httputils.JsonMiddleware(),
				},
			})
		}
	}
}

type HandlerParams struct {
	fx.In

	HandlerSpecValidator openapi3utils.HandlerSpecValidator
}

func NewHandler(params HandlerParams) *handler {
	return &handler{
		handlerSpecValidator: params.HandlerSpecValidator,
	}
}
