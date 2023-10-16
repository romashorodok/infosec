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
		// TODO:
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

func (s *KanbanService) CreateUserBoard(ctx context.Context, token *auth.TokenPayload, params CreateUserBoardParams) {

	tx, err := s.client.Tx(ctx)
	if err != nil {
		// TODO:
	}
	defer tx.Rollback()

	user, err := tx.User.Get(ctx, token.UserID)
	if err != nil {
		// TODO:
	}

	participant, err := tx.Participant.Create().Save(ctx)
	if err != nil {
		// TODO:
	}

	user, err = user.Update().AddParticipants(participant).Save(ctx)
	if err != nil {
		// TODO;
	}

	board, err := tx.Board.Create().SetTitle(params.Title).AddParticipants().AddParticipants(participant).Save(ctx)
	if err != nil {
		// TOOD:
	}

	_ = board

	_ = tx.Commit()
}

func (s KanbanService) DeleteUserBoard(ctx context.Context, token *auth.TokenPayload, boardID int32) {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		// TODO
	}
	defer tx.Rollback()

	boardq := tx.Board.Query().Where(board.IDEQ(int(boardID)))

	participantq := boardq.QueryParticipants().Where(participant.HasBoards())

	participantModel, err := participantq.First(ctx)
	if err != nil {
		// TODO
	}

	user, err := tx.User.Query().Where(
		user.HasParticipantsWith(participant.IDEQ(participantModel.ID)),
	).First(ctx)
	if err != nil {
		// TODO
	}

	// Request first board

	board, err := boardq.First(ctx)
	if err != nil {
		// TODO
	}

	result := tx.Board.DeleteOne(board).Exec(ctx)

	log.Println(result)
	log.Println(err)

	err = tx.Commit()
	if err != nil {
		// TODO
	}

	log.Println("--- Test ---")

	log.Printf("%+v", board)
	log.Printf("%+v", participantModel)
	log.Printf("%+v", user)
	log.Println("Same user", token.UserID == user.ID)
}

func (s *KanbanService) AddBoardParticipant(ctx context.Context, token *auth.TokenPayload, boardID int32, userID uuid.UUID) (any, error) {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return nil, err
		// TODO:
	}
	defer tx.Rollback()

	userBoardq := tx.User.Query().Where(user.IDEQ(token.UserID)).
		QueryParticipants().
		QueryBoards().
		Where(board.IDEQ(int(boardID)))

	board, err := userBoardq.First(ctx)
	if err != nil {
		// TODO: board not found
	}

	userModel, err := tx.User.Get(ctx, userID)
	if err != nil {
		log.Println("User not found", err)
		return nil, err
		// TODO: participant user not found
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

func (s *KanbanService) DeleteBoardParticipant(ctx context.Context, token *auth.TokenPayload, boardID int32, userID uuid.UUID) {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		// TODO:
	}
	defer tx.Rollback()

	userBoardq := tx.User.Query().Where(user.IDEQ(token.UserID)).
		QueryParticipants().
		QueryBoards().
		Where(board.IDEQ(int(boardID)))

	boardModel, err := userBoardq.First(ctx)
	if err != nil {
		// TODO: board not found
	}

	user, err := tx.User.Get(ctx, userID)
	if err != nil {
		log.Println("User not found", err)
		return
		// TODO: participant user not found
	}

	participantToDelete, err := user.QueryParticipants().
		Where(participant.HasBoardsWith(board.IDEQ(int(boardID)))).First(ctx)

	if exists, err := boardModel.QueryParticipants().Where(participant.ID(participantToDelete.ID)).Exist(ctx); !exists || err != nil {
		log.Println("Participant of user not found")
		return
	}

	if err = tx.Participant.DeleteOne(participantToDelete).Exec(ctx); err != nil {
		log.Println("Unable delete user participant")
		return
	}

	_ = tx.Commit()
}

func (s *KanbanService) AddBoardPillar(ctx context.Context, token *auth.TokenPayload, boardID int32) {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		// TODO:
		return
	}
	defer tx.Rollback()

	userBoardq := tx.User.Query().Where(user.IDEQ(token.UserID)).
		QueryParticipants().
		QueryBoards().
		Where(board.IDEQ(int(boardID)))

	board, err := userBoardq.First(ctx)
	if err != nil {
		// TODO:
		return
	}

	pillar, err := tx.Pillar.Create().Save(ctx)
	if err != nil {
		// TODO:
		return
	}

	if err := board.Update().AddPillars(pillar).Exec(ctx); err != nil {
		return
	}

	_ = tx.Commit()
}

func (s *KanbanService) RemoveBoardPillar(ctx context.Context, token *auth.TokenPayload, boardID int32, pillarID int32) {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		// TODO:
		return
	}
	defer tx.Rollback()

	userBoardq := tx.User.Query().Where(user.IDEQ(token.UserID)).
		QueryParticipants().
		QueryBoards().
		Where(board.IDEQ(int(boardID)))

	board, err := userBoardq.First(ctx)
	if err != nil {
		// TODO:
		return
	}

	pillar, err := board.QueryPillars().Where(pillar.IDEQ(int(pillarID))).First(ctx)
	if err != nil {
		return
	}

	if err := tx.Pillar.DeleteOne(pillar); err != nil {
		return
	}

	_ = tx.Commit()
}

type AddBoardTaskParams struct {
	Title       string
	Description string
}

func (s *KanbanService) AddBoardTask(ctx context.Context, token *auth.TokenPayload, boardID int32, pillarID int32, params *AddBoardTaskParams) {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return
	}
	defer tx.Rollback()

	userBoardq := tx.User.Query().Where(user.IDEQ(token.UserID)).
		QueryParticipants().
		QueryBoards().
		Where(board.IDEQ(int(boardID)))

	board, err := userBoardq.First(ctx)
	if err != nil {
		// TODO:
		return
	}

	pillar, err := board.QueryPillars().Where(pillar.IDEQ(int(pillarID))).First(ctx)
	if err != nil {
		return
	}

	task, err := tx.Task.Create().
		SetDescription(params.Description).
		SetTitle(params.Title).
		Save(ctx)
	if err != nil {
		return
	}

	if err = pillar.Update().AddTasks(task).Exec(ctx); err != nil {
		return
	}

	_ = tx.Commit()
}

func (s *KanbanService) DeleteBoardTask(ctx context.Context, token *auth.TokenPayload, boardID int32, taskID int32) {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return
	}
	defer tx.Rollback()

	userBoardq := tx.User.Query().Where(user.IDEQ(token.UserID)).
		QueryParticipants().
		QueryBoards().
		Where(board.IDEQ(int(boardID)))

	board, err := userBoardq.First(ctx)
	if err != nil {
		// TODO:
		return
	}

	taks, err := board.QueryTasks().Where(task.IDEQ(int(taskID))).First(ctx)
	if err != nil {
		return
	}

	if err := tx.Task.DeleteOne(taks); err != nil {
		return
	}

	_ = tx.Commit()
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
