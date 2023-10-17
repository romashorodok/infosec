package kanban

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/romashorodok/infosec/internal/kanban"
	"github.com/romashorodok/infosec/pkg/auth"
	"github.com/romashorodok/infosec/pkg/httputils"
	"github.com/romashorodok/infosec/pkg/openapi3utils"
	"go.uber.org/fx"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest --config=handler.cfg.yaml ./../../../../ent/openapi.json

type handler struct {
	Unimplemented

	handlerSpecValidator openapi3utils.HandlerSpecValidator
	kanbanService        kanban.Kanban
}

var _ ServerInterface = (*handler)(nil)
var _ httputils.HttpHandler = (*handler)(nil)

func (h *handler) ListUser(w http.ResponseWriter, r *http.Request, params ListUserParams) {
	result, err := h.kanbanService.GetUsers(r.Context())
	if err != nil {
		httputils.WriteErrorResponse(w, http.StatusInternalServerError, "Not found users")
		return
	}

	_ = json.NewEncoder(w).Encode(result)
}

func (h *handler) ReadBoard(w http.ResponseWriter, r *http.Request, id int32) {
	token, err := auth.WithTokenPayload(r.Context())
	if err != nil {
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Not found user token payload", err.Error())
		return
	}

	result, err := h.kanbanService.GetUesrBoard(r.Context(), token, id)
	if err != nil {
		httputils.WriteErrorResponse(w, http.StatusInternalServerError, "Unable find user board", err.Error())
		return
	}

	_ = json.NewEncoder(w).Encode(result)
}

func (h *handler) GetUserBoards(w http.ResponseWriter, r *http.Request) {
	token, err := auth.WithTokenPayload(r.Context())
	if err != nil {
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Not found user token payload", err.Error())
		return
	}

	result, err := h.kanbanService.GetUserBoards(r.Context(), token)
	if err != nil {
		httputils.WriteErrorResponse(w, http.StatusInternalServerError, "Unable find user boards", err.Error())
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (h *handler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	token, err := auth.WithTokenPayload(r.Context())
	if err != nil {
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Not found user token payload", err.Error())
		return
	}

	var req CreateBoardJSONBody

	_ = json.NewDecoder(r.Body).Decode(&req)

	if _, err := h.kanbanService.CreateUserBoard(r.Context(), token, kanban.CreateUserBoardParams{
		Title: req.Title,
	}); err != nil {
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Cannot create user board. Err:", err.Error())
		return
	}
}

func (h *handler) DeleteBoard(w http.ResponseWriter, r *http.Request, id int32) {
	token, err := auth.WithTokenPayload(r.Context())
	if err != nil {
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Not found user token payload", err.Error())
		return
	}

	if _, err := h.kanbanService.DeleteUserBoard(r.Context(), token, id); err != nil {
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Cannot delete user board. Err:", err.Error())
		return
	}
}

// Participant

func (h *handler) AddBoardParticipant(w http.ResponseWriter, r *http.Request, id int32) {
	token, err := auth.WithTokenPayload(r.Context())
	if err != nil {
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Not found user token payload", err.Error())
		return
	}

	var req ParticipantBoard

	_ = json.NewDecoder(r.Body).Decode(&req)

	if _, err = h.kanbanService.AddBoardParticipant(r.Context(), token, id, req.UserId); err != nil {
		handleErrorAddBoardParticipant(w, err)
		return
	}
}

func (h *handler) RemoveBoardParticipant(w http.ResponseWriter, r *http.Request, id int32, params RemoveBoardParticipantParams) {
	token, err := auth.WithTokenPayload(r.Context())
	if err != nil {
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Not found user token payload", err.Error())
		return
	}

	var req ParticipantBoard

	_ = json.NewDecoder(r.Body).Decode(&req)

	if _, err := h.kanbanService.DeleteBoardParticipant(r.Context(), token, id, params.UserId); err != nil {
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Cannot delete user board participant. Err:", err.Error())
		return
	}
}

// Pillar

func (h *handler) AddBoardPillar(w http.ResponseWriter, r *http.Request, id int32) {
	token, err := auth.WithTokenPayload(r.Context())
	if err != nil {
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Not found user token payload", err.Error())
		return
	}

	if _, err := h.kanbanService.AddBoardPillar(r.Context(), token, id); err != nil {
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Cannot add board pillar. Err:", err.Error())
		return
	}
}

func (h *handler) RemoveBoardPillar(w http.ResponseWriter, r *http.Request, id int32, params RemoveBoardPillarParams) {
	token, err := auth.WithTokenPayload(r.Context())
	if err != nil {
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Not found user token payload", err.Error())
		return
	}

	if _, err := h.kanbanService.RemoveBoardPillar(r.Context(), token, id, params.PillarId); err != nil {
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Cannot remove board pillar. Err:", err.Error())
		return
	}
}

// Task

func (h *handler) AddBoardTask(w http.ResponseWriter, r *http.Request, id int32, params AddBoardTaskParams) {
	token, err := auth.WithTokenPayload(r.Context())
	if err != nil {
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Not found user token payload", err.Error())
		return
	}

	var req AddBoardTaskJSONBody

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputils.WriteErrorResponse(w, http.StatusInternalServerError, "Unable deserialize body request. Err:", err.Error())
		return
	}

	if _, err := h.kanbanService.AddBoardTask(r.Context(), token, id, params.PillarId, &kanban.AddBoardTaskParams{
		Title:       *req.Title,
		Description: *req.Title,
	}); err != nil {
		httputils.WriteErrorResponse(w, http.StatusInternalServerError, "Unable add board task. Err:", err.Error())
		return
	}
}

func (h *handler) RemoveBoardTask(w http.ResponseWriter, r *http.Request, id int32, params RemoveBoardTaskParams) {
	token, err := auth.WithTokenPayload(r.Context())
	if err != nil {
		httputils.WriteErrorResponse(w, http.StatusPreconditionFailed, "Not found user token payload", err.Error())
		return
	}

	if _, err := h.kanbanService.DeleteBoardTask(r.Context(), token, id, params.TaskId); err != nil {
		httputils.WriteErrorResponse(w, http.StatusInternalServerError, "Unable remove board task. Err:", err.Error())
		return
	}
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

			options := ChiServerOptions{
				BaseRouter: mux,
				Middlewares: []MiddlewareFunc{
					h.handlerSpecValidator(spec),
					httputils.JsonMiddleware(),
				},
			}

			_ = HandlerWithOptions(h, options)
		}
	}
}

type HandlerParams struct {
	fx.In

	HandlerSpecValidator openapi3utils.HandlerSpecValidator
	KanbanService        *kanban.KanbanService
	Logger               *slog.Logger
}

func NewHandler(params HandlerParams) *handler {
	kanbanService := kanban.NewKanbanServiceLogging(params.KanbanService, kanban.KanbanServiceLoggingParams{
		Logger: params.Logger,
	})
	// kanbanService := params.KanbanService

	return &handler{
		handlerSpecValidator: params.HandlerSpecValidator,
		kanbanService:        kanbanService,
	}
}
