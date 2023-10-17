package kanban

import (
	"context"

	"github.com/google/uuid"
	"github.com/romashorodok/infosec/ent"
	"github.com/romashorodok/infosec/pkg/auth"
)

type Kanban interface {
	GetUsers(ctx context.Context) ([]*ent.User, error)
	GetUesrBoard(ctx context.Context, token *auth.TokenPayload, boardID int32) (*ent.Board, error)
	GetUserBoards(ctx context.Context, token *auth.TokenPayload) ([]*ent.Board, error)
	CreateUserBoard(ctx context.Context, token *auth.TokenPayload, params CreateUserBoardParams) (any, error)
	DeleteUserBoard(ctx context.Context, token *auth.TokenPayload, boardID int32) (any, error)
	AddBoardParticipant(ctx context.Context, token *auth.TokenPayload, boardID int32, userID uuid.UUID) (any, error)
	DeleteBoardParticipant(ctx context.Context, token *auth.TokenPayload, boardID int32, userID uuid.UUID) (any, error)
	AddBoardPillar(ctx context.Context, token *auth.TokenPayload, boardID int32) (any, error)
	RemoveBoardPillar(ctx context.Context, token *auth.TokenPayload, boardID int32, pillarID int32) (any, error)
	AddBoardTask(ctx context.Context, token *auth.TokenPayload, boardID int32, pillarID int32, params *AddBoardTaskParams) (any, error)
	DeleteBoardTask(ctx context.Context, token *auth.TokenPayload, boardID int32, taskID int32) (any, error)
}
