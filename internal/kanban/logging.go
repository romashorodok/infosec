package kanban

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/romashorodok/infosec/ent"
	"github.com/romashorodok/infosec/pkg/auth"
)

type KanbanServiceLogging struct {
	log  *slog.Logger
	next Kanban
}

func (s *KanbanServiceLogging) GetUesrBoard(ctx context.Context, token *auth.TokenPayload, boardID int32) (_ *ent.Board, err error) {
	defer func() {
		if err != nil {
			switch err {
			default:
				s.log.Error("Catch errror", "msg", err.Error())
			}
		}
	}()
	return s.next.GetUesrBoard(ctx, token, boardID)
}

func (s *KanbanServiceLogging) AddBoardParticipant(ctx context.Context, token *auth.TokenPayload, boardID int32, userID uuid.UUID) (_ any, err error) {
	defer func(start time.Time) {
		if err != nil {
			switch err {
			case ParticipantAlreadyExistsError:
				s.log.Error(err.Error(), slog.Group("params",
					slog.Any("user_id", userID),
					slog.Any("board_id", boardID),
				))
			default:
				s.log.Error("Catch errror", "msg", err.Error())
			}
		}
		elapsed := time.Since(start)
		s.log.Info("Time elapsed", "elapsed", elapsed)
	}(time.Now())
	return s.next.AddBoardParticipant(ctx, token, boardID, userID)
}

func (s *KanbanServiceLogging) AddBoardPillar(ctx context.Context, token *auth.TokenPayload, boardID int32) (_ any, err error) {
	defer func() {
		if err != nil {
			switch err {
			default:
				s.log.Error("Catch errror", "msg", err.Error())
			}
		}
	}()
	return s.AddBoardPillar(ctx, token, boardID)
}

func (s *KanbanServiceLogging) AddBoardTask(ctx context.Context, token *auth.TokenPayload, boardID int32, pillarID int32, params *AddBoardTaskParams) (_ any, err error) {
	defer func() {
		if err != nil {
			switch err {
			default:
				s.log.Error("Catch errror", "msg", err.Error())
			}
		}
	}()
	return s.AddBoardTask(ctx, token, boardID, pillarID, params)
}

func (s *KanbanServiceLogging) CreateUserBoard(ctx context.Context, token *auth.TokenPayload, params CreateUserBoardParams) (_ any, err error) {
	defer func() {
		if err != nil {
			switch err {
			default:
				s.log.Error("Catch errror", "msg", err.Error())
			}
		}
	}()
	return s.next.CreateUserBoard(ctx, token, params)
}

func (s *KanbanServiceLogging) DeleteBoardParticipant(ctx context.Context, token *auth.TokenPayload, boardID int32, userID uuid.UUID) (_ any, err error) {
	defer func() {
		if err != nil {
			switch err {
			default:
				s.log.Error("Catch errror", "msg", err.Error())
			}
		}
	}()
	return s.next.DeleteBoardParticipant(ctx, token, boardID, userID)
}

func (s *KanbanServiceLogging) DeleteBoardTask(ctx context.Context, token *auth.TokenPayload, boardID int32, taskID int32) (_ any, err error) {
	defer func() {
		if err != nil {
			switch err {
			default:
				s.log.Error("Catch errror", "msg", err.Error())
			}
		}
	}()
	return s.next.DeleteBoardTask(ctx, token, boardID, taskID)
}

func (s *KanbanServiceLogging) DeleteUserBoard(ctx context.Context, token *auth.TokenPayload, boardID int32) (_ any, err error) {
	defer func() {
		if err != nil {
			switch err {
			default:
				s.log.Error("Catch errror", "msg", err.Error())
			}
		}
	}()
	return s.next.DeleteUserBoard(ctx, token, boardID)
}

func (s *KanbanServiceLogging) GetUserBoards(ctx context.Context, token *auth.TokenPayload) (_ []*ent.Board, err error) {
	defer func() {
		if err != nil {
			switch err {
			default:
				s.log.Error("Catch errror", "msg", err.Error())
			}
		}
	}()
	return s.next.GetUserBoards(ctx, token)
}

func (s *KanbanServiceLogging) GetUsers(ctx context.Context) (_ []*ent.User, err error) {
	defer func() {
		if err != nil {
			switch err {
			default:
				s.log.Error("Catch errror", "msg", err.Error())
			}
		}
	}()
	return s.next.GetUsers(ctx)
}

func (s *KanbanServiceLogging) RemoveBoardPillar(ctx context.Context, token *auth.TokenPayload, boardID int32, pillarID int32) (_ any, err error) {
	defer func() {
		if err != nil {
			switch err {
			default:
				s.log.Error("Catch errror", "msg", err.Error())
			}
		}
	}()
	return s.next.RemoveBoardPillar(ctx, token, boardID, pillarID)
}

var _ Kanban = (*KanbanServiceLogging)(nil)

type KanbanServiceLoggingParams struct {
	Logger *slog.Logger
}

func NewKanbanServiceLogging(service Kanban, params KanbanServiceLoggingParams) *KanbanServiceLogging {
	return &KanbanServiceLogging{
		log:  params.Logger,
		next: service,
	}
}
