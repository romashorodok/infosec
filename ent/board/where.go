// Code generated by ent, DO NOT EDIT.

package board

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/romashorodok/infosec/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Board {
	return predicate.Board(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Board {
	return predicate.Board(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Board {
	return predicate.Board(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Board {
	return predicate.Board(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Board {
	return predicate.Board(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Board {
	return predicate.Board(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Board {
	return predicate.Board(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Board {
	return predicate.Board(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Board {
	return predicate.Board(sql.FieldLTE(FieldID, id))
}

// Title applies equality check predicate on the "title" field. It's identical to TitleEQ.
func Title(v string) predicate.Board {
	return predicate.Board(sql.FieldEQ(FieldTitle, v))
}

// TitleEQ applies the EQ predicate on the "title" field.
func TitleEQ(v string) predicate.Board {
	return predicate.Board(sql.FieldEQ(FieldTitle, v))
}

// TitleNEQ applies the NEQ predicate on the "title" field.
func TitleNEQ(v string) predicate.Board {
	return predicate.Board(sql.FieldNEQ(FieldTitle, v))
}

// TitleIn applies the In predicate on the "title" field.
func TitleIn(vs ...string) predicate.Board {
	return predicate.Board(sql.FieldIn(FieldTitle, vs...))
}

// TitleNotIn applies the NotIn predicate on the "title" field.
func TitleNotIn(vs ...string) predicate.Board {
	return predicate.Board(sql.FieldNotIn(FieldTitle, vs...))
}

// TitleGT applies the GT predicate on the "title" field.
func TitleGT(v string) predicate.Board {
	return predicate.Board(sql.FieldGT(FieldTitle, v))
}

// TitleGTE applies the GTE predicate on the "title" field.
func TitleGTE(v string) predicate.Board {
	return predicate.Board(sql.FieldGTE(FieldTitle, v))
}

// TitleLT applies the LT predicate on the "title" field.
func TitleLT(v string) predicate.Board {
	return predicate.Board(sql.FieldLT(FieldTitle, v))
}

// TitleLTE applies the LTE predicate on the "title" field.
func TitleLTE(v string) predicate.Board {
	return predicate.Board(sql.FieldLTE(FieldTitle, v))
}

// TitleContains applies the Contains predicate on the "title" field.
func TitleContains(v string) predicate.Board {
	return predicate.Board(sql.FieldContains(FieldTitle, v))
}

// TitleHasPrefix applies the HasPrefix predicate on the "title" field.
func TitleHasPrefix(v string) predicate.Board {
	return predicate.Board(sql.FieldHasPrefix(FieldTitle, v))
}

// TitleHasSuffix applies the HasSuffix predicate on the "title" field.
func TitleHasSuffix(v string) predicate.Board {
	return predicate.Board(sql.FieldHasSuffix(FieldTitle, v))
}

// TitleEqualFold applies the EqualFold predicate on the "title" field.
func TitleEqualFold(v string) predicate.Board {
	return predicate.Board(sql.FieldEqualFold(FieldTitle, v))
}

// TitleContainsFold applies the ContainsFold predicate on the "title" field.
func TitleContainsFold(v string) predicate.Board {
	return predicate.Board(sql.FieldContainsFold(FieldTitle, v))
}

// HasTasks applies the HasEdge predicate on the "tasks" edge.
func HasTasks() predicate.Board {
	return predicate.Board(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, TasksTable, TasksColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTasksWith applies the HasEdge predicate on the "tasks" edge with a given conditions (other predicates).
func HasTasksWith(preds ...predicate.Task) predicate.Board {
	return predicate.Board(func(s *sql.Selector) {
		step := newTasksStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasParticipants applies the HasEdge predicate on the "participants" edge.
func HasParticipants() predicate.Board {
	return predicate.Board(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, ParticipantsTable, ParticipantsPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasParticipantsWith applies the HasEdge predicate on the "participants" edge with a given conditions (other predicates).
func HasParticipantsWith(preds ...predicate.Participant) predicate.Board {
	return predicate.Board(func(s *sql.Selector) {
		step := newParticipantsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasPillars applies the HasEdge predicate on the "pillars" edge.
func HasPillars() predicate.Board {
	return predicate.Board(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PillarsTable, PillarsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPillarsWith applies the HasEdge predicate on the "pillars" edge with a given conditions (other predicates).
func HasPillarsWith(preds ...predicate.Pillar) predicate.Board {
	return predicate.Board(func(s *sql.Selector) {
		step := newPillarsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Board) predicate.Board {
	return predicate.Board(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Board) predicate.Board {
	return predicate.Board(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Board) predicate.Board {
	return predicate.Board(sql.NotPredicates(p))
}
