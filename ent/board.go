// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/romashorodok/infosec/ent/board"
)

// Board is the model entity for the Board schema.
type Board struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Title holds the value of the "title" field.
	Title string `json:"title,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the BoardQuery when eager-loading is set.
	Edges        BoardEdges `json:"edges"`
	selectValues sql.SelectValues
}

// BoardEdges holds the relations/edges for other nodes in the graph.
type BoardEdges struct {
	// Tasks holds the value of the tasks edge.
	Tasks []*Task `json:"tasks,omitempty"`
	// Participants holds the value of the participants edge.
	Participants []*Participant `json:"participants,omitempty"`
	// Pillars holds the value of the pillars edge.
	Pillars []*Pillar `json:"pillars,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// TasksOrErr returns the Tasks value or an error if the edge
// was not loaded in eager-loading.
func (e BoardEdges) TasksOrErr() ([]*Task, error) {
	if e.loadedTypes[0] {
		return e.Tasks, nil
	}
	return nil, &NotLoadedError{edge: "tasks"}
}

// ParticipantsOrErr returns the Participants value or an error if the edge
// was not loaded in eager-loading.
func (e BoardEdges) ParticipantsOrErr() ([]*Participant, error) {
	if e.loadedTypes[1] {
		return e.Participants, nil
	}
	return nil, &NotLoadedError{edge: "participants"}
}

// PillarsOrErr returns the Pillars value or an error if the edge
// was not loaded in eager-loading.
func (e BoardEdges) PillarsOrErr() ([]*Pillar, error) {
	if e.loadedTypes[2] {
		return e.Pillars, nil
	}
	return nil, &NotLoadedError{edge: "pillars"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Board) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case board.FieldID:
			values[i] = new(sql.NullInt64)
		case board.FieldTitle:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Board fields.
func (b *Board) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case board.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			b.ID = int(value.Int64)
		case board.FieldTitle:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field title", values[i])
			} else if value.Valid {
				b.Title = value.String
			}
		default:
			b.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Board.
// This includes values selected through modifiers, order, etc.
func (b *Board) Value(name string) (ent.Value, error) {
	return b.selectValues.Get(name)
}

// QueryTasks queries the "tasks" edge of the Board entity.
func (b *Board) QueryTasks() *TaskQuery {
	return NewBoardClient(b.config).QueryTasks(b)
}

// QueryParticipants queries the "participants" edge of the Board entity.
func (b *Board) QueryParticipants() *ParticipantQuery {
	return NewBoardClient(b.config).QueryParticipants(b)
}

// QueryPillars queries the "pillars" edge of the Board entity.
func (b *Board) QueryPillars() *PillarQuery {
	return NewBoardClient(b.config).QueryPillars(b)
}

// Update returns a builder for updating this Board.
// Note that you need to call Board.Unwrap() before calling this method if this Board
// was returned from a transaction, and the transaction was committed or rolled back.
func (b *Board) Update() *BoardUpdateOne {
	return NewBoardClient(b.config).UpdateOne(b)
}

// Unwrap unwraps the Board entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (b *Board) Unwrap() *Board {
	_tx, ok := b.config.driver.(*txDriver)
	if !ok {
		panic("ent: Board is not a transactional entity")
	}
	b.config.driver = _tx.drv
	return b
}

// String implements the fmt.Stringer.
func (b *Board) String() string {
	var builder strings.Builder
	builder.WriteString("Board(")
	builder.WriteString(fmt.Sprintf("id=%v, ", b.ID))
	builder.WriteString("title=")
	builder.WriteString(b.Title)
	builder.WriteByte(')')
	return builder.String()
}

// Boards is a parsable slice of Board.
type Boards []*Board
