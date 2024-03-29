// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/romashorodok/infosec/ent/board"
	"github.com/romashorodok/infosec/ent/participant"
	"github.com/romashorodok/infosec/ent/predicate"
	"github.com/romashorodok/infosec/ent/task"
	"github.com/romashorodok/infosec/ent/user"
)

// ParticipantUpdate is the builder for updating Participant entities.
type ParticipantUpdate struct {
	config
	hooks    []Hook
	mutation *ParticipantMutation
}

// Where appends a list predicates to the ParticipantUpdate builder.
func (pu *ParticipantUpdate) Where(ps ...predicate.Participant) *ParticipantUpdate {
	pu.mutation.Where(ps...)
	return pu
}

// AddBoardIDs adds the "boards" edge to the Board entity by IDs.
func (pu *ParticipantUpdate) AddBoardIDs(ids ...int) *ParticipantUpdate {
	pu.mutation.AddBoardIDs(ids...)
	return pu
}

// AddBoards adds the "boards" edges to the Board entity.
func (pu *ParticipantUpdate) AddBoards(b ...*Board) *ParticipantUpdate {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return pu.AddBoardIDs(ids...)
}

// AddTaskIDs adds the "tasks" edge to the Task entity by IDs.
func (pu *ParticipantUpdate) AddTaskIDs(ids ...int) *ParticipantUpdate {
	pu.mutation.AddTaskIDs(ids...)
	return pu
}

// AddTasks adds the "tasks" edges to the Task entity.
func (pu *ParticipantUpdate) AddTasks(t ...*Task) *ParticipantUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return pu.AddTaskIDs(ids...)
}

// SetUserID sets the "user" edge to the User entity by ID.
func (pu *ParticipantUpdate) SetUserID(id uuid.UUID) *ParticipantUpdate {
	pu.mutation.SetUserID(id)
	return pu
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (pu *ParticipantUpdate) SetNillableUserID(id *uuid.UUID) *ParticipantUpdate {
	if id != nil {
		pu = pu.SetUserID(*id)
	}
	return pu
}

// SetUser sets the "user" edge to the User entity.
func (pu *ParticipantUpdate) SetUser(u *User) *ParticipantUpdate {
	return pu.SetUserID(u.ID)
}

// Mutation returns the ParticipantMutation object of the builder.
func (pu *ParticipantUpdate) Mutation() *ParticipantMutation {
	return pu.mutation
}

// ClearBoards clears all "boards" edges to the Board entity.
func (pu *ParticipantUpdate) ClearBoards() *ParticipantUpdate {
	pu.mutation.ClearBoards()
	return pu
}

// RemoveBoardIDs removes the "boards" edge to Board entities by IDs.
func (pu *ParticipantUpdate) RemoveBoardIDs(ids ...int) *ParticipantUpdate {
	pu.mutation.RemoveBoardIDs(ids...)
	return pu
}

// RemoveBoards removes "boards" edges to Board entities.
func (pu *ParticipantUpdate) RemoveBoards(b ...*Board) *ParticipantUpdate {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return pu.RemoveBoardIDs(ids...)
}

// ClearTasks clears all "tasks" edges to the Task entity.
func (pu *ParticipantUpdate) ClearTasks() *ParticipantUpdate {
	pu.mutation.ClearTasks()
	return pu
}

// RemoveTaskIDs removes the "tasks" edge to Task entities by IDs.
func (pu *ParticipantUpdate) RemoveTaskIDs(ids ...int) *ParticipantUpdate {
	pu.mutation.RemoveTaskIDs(ids...)
	return pu
}

// RemoveTasks removes "tasks" edges to Task entities.
func (pu *ParticipantUpdate) RemoveTasks(t ...*Task) *ParticipantUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return pu.RemoveTaskIDs(ids...)
}

// ClearUser clears the "user" edge to the User entity.
func (pu *ParticipantUpdate) ClearUser() *ParticipantUpdate {
	pu.mutation.ClearUser()
	return pu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (pu *ParticipantUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, pu.sqlSave, pu.mutation, pu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (pu *ParticipantUpdate) SaveX(ctx context.Context) int {
	affected, err := pu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pu *ParticipantUpdate) Exec(ctx context.Context) error {
	_, err := pu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pu *ParticipantUpdate) ExecX(ctx context.Context) {
	if err := pu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (pu *ParticipantUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(participant.Table, participant.Columns, sqlgraph.NewFieldSpec(participant.FieldID, field.TypeInt))
	if ps := pu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if pu.mutation.BoardsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   participant.BoardsTable,
			Columns: participant.BoardsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(board.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.RemovedBoardsIDs(); len(nodes) > 0 && !pu.mutation.BoardsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   participant.BoardsTable,
			Columns: participant.BoardsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(board.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.BoardsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   participant.BoardsTable,
			Columns: participant.BoardsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(board.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pu.mutation.TasksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   participant.TasksTable,
			Columns: participant.TasksPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(task.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.RemovedTasksIDs(); len(nodes) > 0 && !pu.mutation.TasksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   participant.TasksTable,
			Columns: participant.TasksPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(task.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.TasksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   participant.TasksTable,
			Columns: participant.TasksPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(task.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pu.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   participant.UserTable,
			Columns: []string{participant.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pu.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   participant.UserTable,
			Columns: []string{participant.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, pu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{participant.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	pu.mutation.done = true
	return n, nil
}

// ParticipantUpdateOne is the builder for updating a single Participant entity.
type ParticipantUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ParticipantMutation
}

// AddBoardIDs adds the "boards" edge to the Board entity by IDs.
func (puo *ParticipantUpdateOne) AddBoardIDs(ids ...int) *ParticipantUpdateOne {
	puo.mutation.AddBoardIDs(ids...)
	return puo
}

// AddBoards adds the "boards" edges to the Board entity.
func (puo *ParticipantUpdateOne) AddBoards(b ...*Board) *ParticipantUpdateOne {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return puo.AddBoardIDs(ids...)
}

// AddTaskIDs adds the "tasks" edge to the Task entity by IDs.
func (puo *ParticipantUpdateOne) AddTaskIDs(ids ...int) *ParticipantUpdateOne {
	puo.mutation.AddTaskIDs(ids...)
	return puo
}

// AddTasks adds the "tasks" edges to the Task entity.
func (puo *ParticipantUpdateOne) AddTasks(t ...*Task) *ParticipantUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return puo.AddTaskIDs(ids...)
}

// SetUserID sets the "user" edge to the User entity by ID.
func (puo *ParticipantUpdateOne) SetUserID(id uuid.UUID) *ParticipantUpdateOne {
	puo.mutation.SetUserID(id)
	return puo
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (puo *ParticipantUpdateOne) SetNillableUserID(id *uuid.UUID) *ParticipantUpdateOne {
	if id != nil {
		puo = puo.SetUserID(*id)
	}
	return puo
}

// SetUser sets the "user" edge to the User entity.
func (puo *ParticipantUpdateOne) SetUser(u *User) *ParticipantUpdateOne {
	return puo.SetUserID(u.ID)
}

// Mutation returns the ParticipantMutation object of the builder.
func (puo *ParticipantUpdateOne) Mutation() *ParticipantMutation {
	return puo.mutation
}

// ClearBoards clears all "boards" edges to the Board entity.
func (puo *ParticipantUpdateOne) ClearBoards() *ParticipantUpdateOne {
	puo.mutation.ClearBoards()
	return puo
}

// RemoveBoardIDs removes the "boards" edge to Board entities by IDs.
func (puo *ParticipantUpdateOne) RemoveBoardIDs(ids ...int) *ParticipantUpdateOne {
	puo.mutation.RemoveBoardIDs(ids...)
	return puo
}

// RemoveBoards removes "boards" edges to Board entities.
func (puo *ParticipantUpdateOne) RemoveBoards(b ...*Board) *ParticipantUpdateOne {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return puo.RemoveBoardIDs(ids...)
}

// ClearTasks clears all "tasks" edges to the Task entity.
func (puo *ParticipantUpdateOne) ClearTasks() *ParticipantUpdateOne {
	puo.mutation.ClearTasks()
	return puo
}

// RemoveTaskIDs removes the "tasks" edge to Task entities by IDs.
func (puo *ParticipantUpdateOne) RemoveTaskIDs(ids ...int) *ParticipantUpdateOne {
	puo.mutation.RemoveTaskIDs(ids...)
	return puo
}

// RemoveTasks removes "tasks" edges to Task entities.
func (puo *ParticipantUpdateOne) RemoveTasks(t ...*Task) *ParticipantUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return puo.RemoveTaskIDs(ids...)
}

// ClearUser clears the "user" edge to the User entity.
func (puo *ParticipantUpdateOne) ClearUser() *ParticipantUpdateOne {
	puo.mutation.ClearUser()
	return puo
}

// Where appends a list predicates to the ParticipantUpdate builder.
func (puo *ParticipantUpdateOne) Where(ps ...predicate.Participant) *ParticipantUpdateOne {
	puo.mutation.Where(ps...)
	return puo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (puo *ParticipantUpdateOne) Select(field string, fields ...string) *ParticipantUpdateOne {
	puo.fields = append([]string{field}, fields...)
	return puo
}

// Save executes the query and returns the updated Participant entity.
func (puo *ParticipantUpdateOne) Save(ctx context.Context) (*Participant, error) {
	return withHooks(ctx, puo.sqlSave, puo.mutation, puo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (puo *ParticipantUpdateOne) SaveX(ctx context.Context) *Participant {
	node, err := puo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (puo *ParticipantUpdateOne) Exec(ctx context.Context) error {
	_, err := puo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (puo *ParticipantUpdateOne) ExecX(ctx context.Context) {
	if err := puo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (puo *ParticipantUpdateOne) sqlSave(ctx context.Context) (_node *Participant, err error) {
	_spec := sqlgraph.NewUpdateSpec(participant.Table, participant.Columns, sqlgraph.NewFieldSpec(participant.FieldID, field.TypeInt))
	id, ok := puo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Participant.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := puo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, participant.FieldID)
		for _, f := range fields {
			if !participant.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != participant.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := puo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if puo.mutation.BoardsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   participant.BoardsTable,
			Columns: participant.BoardsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(board.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.RemovedBoardsIDs(); len(nodes) > 0 && !puo.mutation.BoardsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   participant.BoardsTable,
			Columns: participant.BoardsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(board.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.BoardsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   participant.BoardsTable,
			Columns: participant.BoardsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(board.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if puo.mutation.TasksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   participant.TasksTable,
			Columns: participant.TasksPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(task.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.RemovedTasksIDs(); len(nodes) > 0 && !puo.mutation.TasksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   participant.TasksTable,
			Columns: participant.TasksPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(task.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.TasksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   participant.TasksTable,
			Columns: participant.TasksPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(task.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if puo.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   participant.UserTable,
			Columns: []string{participant.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := puo.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   participant.UserTable,
			Columns: []string{participant.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Participant{config: puo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, puo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{participant.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	puo.mutation.done = true
	return _node, nil
}
