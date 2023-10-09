package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"

	entlocal "github.com/romashorodok/infosec/ent"
)

// Board holds the schema definition for the Board entity.
type Board struct {
	ent.Schema
}

// Fields of the Board.
func (Board) Fields() []ent.Field {
	return []ent.Field{}
}

// Edges of the Board.
func (Board) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tasks", Task.Type),
		edge.To("participants", Participant.Type),
	}
}

func (Board) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entlocal.ElkSecurity,
	}
}
