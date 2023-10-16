package kanban

import (
	"net/http"

	"github.com/romashorodok/infosec/internal/kanban"
	"github.com/romashorodok/infosec/pkg/httputils"
)

func handleErrorAddBoardParticipant(w http.ResponseWriter, err error) {
	switch err {
	case kanban.ParticipantAlreadyExistsError:
		httputils.WriteErrorResponse(w, http.StatusNotAcceptable, err.Error())
		return

	default:
		httputils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
