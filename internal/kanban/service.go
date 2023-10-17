package kanban

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/romashorodok/infosec/ent"
	"github.com/romashorodok/infosec/ent/board"
	"github.com/romashorodok/infosec/ent/participant"
	"github.com/romashorodok/infosec/ent/pillar"
	"github.com/romashorodok/infosec/ent/task"
	"github.com/romashorodok/infosec/ent/user"
	"github.com/romashorodok/infosec/pkg/auth"
	"go.uber.org/fx"
)

var (
	ParticipantAlreadyExistsError error = errors.New("Participant already exists")
)

type KanbanService struct {
	client *ent.Client
}

var _ Kanban = (*KanbanService)(nil)

func (s *KanbanService) GetUsers(ctx context.Context) ([]*ent.User, error) {
	return s.client.User.Query().Select("id", "username").All(ctx)
}

func (s *KanbanService) GetUesrBoard(ctx context.Context, token *auth.TokenPayload, boardID int32) (*ent.Board, error) {
	user, err := s.client.User.Get(ctx, token.UserID)
	if err != nil {
		return nil, err
	}

	result, err := user.QueryParticipants().QueryBoards().
		WithParticipants(func(q *ent.ParticipantQuery) {
			q.WithUser(func(q *ent.UserQuery) {
				q.Select("username", "id")
			})
		}).
		Where(board.IDEQ(int(boardID))).
		First(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *KanbanService) GetUserBoards(ctx context.Context, token *auth.TokenPayload) ([]*ent.Board, error) {

	user, err := s.client.User.Get(ctx, token.UserID)
	if err != nil {
		return nil, err
	}

	result, err := user.QueryParticipants().QueryBoards().All(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type CreateUserBoardParams struct {
	Title string
}

func (s *KanbanService) CreateUserBoard(ctx context.Context, token *auth.TokenPayload, params CreateUserBoardParams) (any, error) {

	tx, err := s.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user, err := tx.User.Get(ctx, token.UserID)
	if err != nil {
		return nil, err
	}

	participant, err := tx.Participant.Create().Save(ctx)
	if err != nil {
		return nil, err
	}

	user, err = user.Update().AddParticipants(participant).Save(ctx)
	if err != nil {
		return nil, err
	}

	_, err = tx.Board.Create().SetTitle(params.Title).AddParticipants().AddParticipants(participant).Save(ctx)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s KanbanService) DeleteUserBoard(ctx context.Context, token *auth.TokenPayload, boardID int32) (any, error) {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	boardq := tx.Board.Query().Where(board.IDEQ(int(boardID)))

	participantq := boardq.QueryParticipants().Where(participant.HasBoards())

	participantModel, err := participantq.First(ctx)
	if err != nil {
		return nil, err
	}

	if _, err = tx.User.Query().Where(
		user.HasParticipantsWith(participant.IDEQ(participantModel.ID)),
	).First(ctx); err != nil {
		return nil, err
	}

	board, err := boardq.First(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Board.DeleteOne(board).Exec(ctx); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *KanbanService) AddBoardParticipant(ctx context.Context, token *auth.TokenPayload, boardID int32, userID uuid.UUID) (any, error) {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	userBoardq := tx.User.Query().Where(user.IDEQ(token.UserID)).
		QueryParticipants().
		QueryBoards().
		Where(board.IDEQ(int(boardID)))

	board, err := userBoardq.First(ctx)
	if err != nil {
		return nil, err
	}

	userModel, err := tx.User.Get(ctx, userID)
	if err != nil {
		log.Println("User not found", err)
		return nil, err
	}

	exists, err := board.QueryParticipants().Where(participant.HasUserWith(user.IDEQ(userID))).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ParticipantAlreadyExistsError
	}

	participant, err := tx.Participant.Create().Save(ctx)
	if err != nil {
		return nil, err
	}

	if _, err = userModel.Update().AddParticipants(participant).Save(ctx); err != nil {
		return nil, err
	}

	if err = participant.Update().AddBoards(board).Exec(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *KanbanService) DeleteBoardParticipant(ctx context.Context, token *auth.TokenPayload, boardID int32, userID uuid.UUID) (any, error) {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	userBoardq := tx.User.Query().Where(user.IDEQ(token.UserID)).
		QueryParticipants().
		QueryBoards().
		Where(board.IDEQ(int(boardID)))

	boardModel, err := userBoardq.First(ctx)
	if err != nil {
		return nil, err
	}

	user, err := tx.User.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	participantToDelete, err := user.QueryParticipants().
		Where(participant.HasBoardsWith(board.IDEQ(int(boardID)))).First(ctx)

	if exists, err := boardModel.QueryParticipants().Where(participant.ID(participantToDelete.ID)).Exist(ctx); !exists || err != nil {
		return nil, err
	}

	if err = tx.Participant.DeleteOne(participantToDelete).Exec(ctx); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *KanbanService) AddBoardPillar(ctx context.Context, token *auth.TokenPayload, boardID int32) (any, error) {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	userBoardq := tx.User.Query().Where(user.IDEQ(token.UserID)).
		QueryParticipants().
		QueryBoards().
		Where(board.IDEQ(int(boardID)))

	board, err := userBoardq.First(ctx)
	if err != nil {
		return nil, err
	}

	pillar, err := tx.Pillar.Create().Save(ctx)
	if err != nil {
		return nil, err
	}

	if err := board.Update().AddPillars(pillar).Exec(ctx); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *KanbanService) RemoveBoardPillar(ctx context.Context, token *auth.TokenPayload, boardID int32, pillarID int32) (any, error) {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	userBoardq := tx.User.Query().Where(user.IDEQ(token.UserID)).
		QueryParticipants().
		QueryBoards().
		Where(board.IDEQ(int(boardID)))

	board, err := userBoardq.First(ctx)
	if err != nil {
		return nil, err
	}

	pillar, err := board.QueryPillars().Where(pillar.IDEQ(int(pillarID))).First(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Pillar.DeleteOne(pillar).Exec(ctx); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return nil, err
}

type AddBoardTaskParams struct {
	Title       string
	Description string
}

func (s *KanbanService) AddBoardTask(ctx context.Context, token *auth.TokenPayload, boardID int32, pillarID int32, params *AddBoardTaskParams) (any, error) {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	userBoardq := tx.User.Query().Where(user.IDEQ(token.UserID)).
		QueryParticipants().
		QueryBoards().
		Where(board.IDEQ(int(boardID)))

	board, err := userBoardq.First(ctx)
	if err != nil {
		return nil, err
	}

	pillar, err := board.QueryPillars().Where(pillar.IDEQ(int(pillarID))).First(ctx)
	if err != nil {
		return nil, err
	}

	task, err := tx.Task.Create().
		SetDescription(params.Description).
		SetTitle(params.Title).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	if err = pillar.Update().AddTasks(task).Exec(ctx); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *KanbanService) DeleteBoardTask(ctx context.Context, token *auth.TokenPayload, boardID int32, taskID int32) (any, error) {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	userBoardq := tx.User.Query().Where(user.IDEQ(token.UserID)).
		QueryParticipants().
		QueryBoards().
		Where(board.IDEQ(int(boardID)))

	board, err := userBoardq.First(ctx)
	if err != nil {
		return nil, err
	}

	taks, err := board.QueryTasks().Where(task.IDEQ(int(taskID))).First(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Task.DeleteOne(taks).Exec(ctx); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}

type KanbanServiceParams struct {
	fx.In

	Client *ent.Client
}

func NewKanbanService(params KanbanServiceParams) *KanbanService {
	return &KanbanService{
		client: params.Client,
	}
}
