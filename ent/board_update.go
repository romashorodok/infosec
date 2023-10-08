// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/romashorodok/infosec/ent/board"
	"github.com/romashorodok/infosec/ent/participant"
	"github.com/romashorodok/infosec/ent/predicate"
	"github.com/romashorodok/infosec/ent/task"
)

// BoardUpdate is the builder for updating Board entities.
type BoardUpdate struct {
	config
	hooks    []Hook
	mutation *BoardMutation
}

// Where appends a list predicates to the BoardUpdate builder.
func (bu *BoardUpdate) Where(ps ...predicate.Board) *BoardUpdate {
	bu.mutation.Where(ps...)
	return bu
}

// AddTaskIDs adds the "tasks" edge to the Task entity by IDs.
func (bu *BoardUpdate) AddTaskIDs(ids ...int) *BoardUpdate {
	bu.mutation.AddTaskIDs(ids...)
	return bu
}

// AddTasks adds the "tasks" edges to the Task entity.
func (bu *BoardUpdate) AddTasks(t ...*Task) *BoardUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return bu.AddTaskIDs(ids...)
}

// AddParticipantIDs adds the "participants" edge to the Participant entity by IDs.
func (bu *BoardUpdate) AddParticipantIDs(ids ...int) *BoardUpdate {
	bu.mutation.AddParticipantIDs(ids...)
	return bu
}

// AddParticipants adds the "participants" edges to the Participant entity.
func (bu *BoardUpdate) AddParticipants(p ...*Participant) *BoardUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return bu.AddParticipantIDs(ids...)
}

// Mutation returns the BoardMutation object of the builder.
func (bu *BoardUpdate) Mutation() *BoardMutation {
	return bu.mutation
}

// ClearTasks clears all "tasks" edges to the Task entity.
func (bu *BoardUpdate) ClearTasks() *BoardUpdate {
	bu.mutation.ClearTasks()
	return bu
}

// RemoveTaskIDs removes the "tasks" edge to Task entities by IDs.
func (bu *BoardUpdate) RemoveTaskIDs(ids ...int) *BoardUpdate {
	bu.mutation.RemoveTaskIDs(ids...)
	return bu
}

// RemoveTasks removes "tasks" edges to Task entities.
func (bu *BoardUpdate) RemoveTasks(t ...*Task) *BoardUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return bu.RemoveTaskIDs(ids...)
}

// ClearParticipants clears all "participants" edges to the Participant entity.
func (bu *BoardUpdate) ClearParticipants() *BoardUpdate {
	bu.mutation.ClearParticipants()
	return bu
}

// RemoveParticipantIDs removes the "participants" edge to Participant entities by IDs.
func (bu *BoardUpdate) RemoveParticipantIDs(ids ...int) *BoardUpdate {
	bu.mutation.RemoveParticipantIDs(ids...)
	return bu
}

// RemoveParticipants removes "participants" edges to Participant entities.
func (bu *BoardUpdate) RemoveParticipants(p ...*Participant) *BoardUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return bu.RemoveParticipantIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (bu *BoardUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, bu.sqlSave, bu.mutation, bu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (bu *BoardUpdate) SaveX(ctx context.Context) int {
	affected, err := bu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (bu *BoardUpdate) Exec(ctx context.Context) error {
	_, err := bu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bu *BoardUpdate) ExecX(ctx context.Context) {
	if err := bu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (bu *BoardUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(board.Table, board.Columns, sqlgraph.NewFieldSpec(board.FieldID, field.TypeInt))
	if ps := bu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if bu.mutation.TasksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   board.TasksTable,
			Columns: []string{board.TasksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(task.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bu.mutation.RemovedTasksIDs(); len(nodes) > 0 && !bu.mutation.TasksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   board.TasksTable,
			Columns: []string{board.TasksColumn},
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
	if nodes := bu.mutation.TasksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   board.TasksTable,
			Columns: []string{board.TasksColumn},
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
	if bu.mutation.ParticipantsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   board.ParticipantsTable,
			Columns: board.ParticipantsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(participant.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bu.mutation.RemovedParticipantsIDs(); len(nodes) > 0 && !bu.mutation.ParticipantsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   board.ParticipantsTable,
			Columns: board.ParticipantsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(participant.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bu.mutation.ParticipantsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   board.ParticipantsTable,
			Columns: board.ParticipantsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(participant.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, bu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{board.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	bu.mutation.done = true
	return n, nil
}

// BoardUpdateOne is the builder for updating a single Board entity.
type BoardUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *BoardMutation
}

// AddTaskIDs adds the "tasks" edge to the Task entity by IDs.
func (buo *BoardUpdateOne) AddTaskIDs(ids ...int) *BoardUpdateOne {
	buo.mutation.AddTaskIDs(ids...)
	return buo
}

// AddTasks adds the "tasks" edges to the Task entity.
func (buo *BoardUpdateOne) AddTasks(t ...*Task) *BoardUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return buo.AddTaskIDs(ids...)
}

// AddParticipantIDs adds the "participants" edge to the Participant entity by IDs.
func (buo *BoardUpdateOne) AddParticipantIDs(ids ...int) *BoardUpdateOne {
	buo.mutation.AddParticipantIDs(ids...)
	return buo
}

// AddParticipants adds the "participants" edges to the Participant entity.
func (buo *BoardUpdateOne) AddParticipants(p ...*Participant) *BoardUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return buo.AddParticipantIDs(ids...)
}

// Mutation returns the BoardMutation object of the builder.
func (buo *BoardUpdateOne) Mutation() *BoardMutation {
	return buo.mutation
}

// ClearTasks clears all "tasks" edges to the Task entity.
func (buo *BoardUpdateOne) ClearTasks() *BoardUpdateOne {
	buo.mutation.ClearTasks()
	return buo
}

// RemoveTaskIDs removes the "tasks" edge to Task entities by IDs.
func (buo *BoardUpdateOne) RemoveTaskIDs(ids ...int) *BoardUpdateOne {
	buo.mutation.RemoveTaskIDs(ids...)
	return buo
}

// RemoveTasks removes "tasks" edges to Task entities.
func (buo *BoardUpdateOne) RemoveTasks(t ...*Task) *BoardUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return buo.RemoveTaskIDs(ids...)
}

// ClearParticipants clears all "participants" edges to the Participant entity.
func (buo *BoardUpdateOne) ClearParticipants() *BoardUpdateOne {
	buo.mutation.ClearParticipants()
	return buo
}

// RemoveParticipantIDs removes the "participants" edge to Participant entities by IDs.
func (buo *BoardUpdateOne) RemoveParticipantIDs(ids ...int) *BoardUpdateOne {
	buo.mutation.RemoveParticipantIDs(ids...)
	return buo
}

// RemoveParticipants removes "participants" edges to Participant entities.
func (buo *BoardUpdateOne) RemoveParticipants(p ...*Participant) *BoardUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return buo.RemoveParticipantIDs(ids...)
}

// Where appends a list predicates to the BoardUpdate builder.
func (buo *BoardUpdateOne) Where(ps ...predicate.Board) *BoardUpdateOne {
	buo.mutation.Where(ps...)
	return buo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (buo *BoardUpdateOne) Select(field string, fields ...string) *BoardUpdateOne {
	buo.fields = append([]string{field}, fields...)
	return buo
}

// Save executes the query and returns the updated Board entity.
func (buo *BoardUpdateOne) Save(ctx context.Context) (*Board, error) {
	return withHooks(ctx, buo.sqlSave, buo.mutation, buo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (buo *BoardUpdateOne) SaveX(ctx context.Context) *Board {
	node, err := buo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (buo *BoardUpdateOne) Exec(ctx context.Context) error {
	_, err := buo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (buo *BoardUpdateOne) ExecX(ctx context.Context) {
	if err := buo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (buo *BoardUpdateOne) sqlSave(ctx context.Context) (_node *Board, err error) {
	_spec := sqlgraph.NewUpdateSpec(board.Table, board.Columns, sqlgraph.NewFieldSpec(board.FieldID, field.TypeInt))
	id, ok := buo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Board.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := buo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, board.FieldID)
		for _, f := range fields {
			if !board.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != board.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := buo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if buo.mutation.TasksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   board.TasksTable,
			Columns: []string{board.TasksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(task.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := buo.mutation.RemovedTasksIDs(); len(nodes) > 0 && !buo.mutation.TasksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   board.TasksTable,
			Columns: []string{board.TasksColumn},
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
	if nodes := buo.mutation.TasksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   board.TasksTable,
			Columns: []string{board.TasksColumn},
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
	if buo.mutation.ParticipantsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   board.ParticipantsTable,
			Columns: board.ParticipantsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(participant.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := buo.mutation.RemovedParticipantsIDs(); len(nodes) > 0 && !buo.mutation.ParticipantsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   board.ParticipantsTable,
			Columns: board.ParticipantsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(participant.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := buo.mutation.ParticipantsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   board.ParticipantsTable,
			Columns: board.ParticipantsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(participant.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Board{config: buo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, buo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{board.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	buo.mutation.done = true
	return _node, nil
}