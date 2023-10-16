package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/romashorodok/infosec/pkg/entutils"
)

// Board holds the schema definition for the Board entity.
type Board struct {
	ent.Schema
}

// Fields of the Board.
func (Board) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
	}
}

// Edges of the Board.
func (Board) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tasks", Task.Type),
		edge.To("participants", Participant.Type),
		edge.To("pillars", Pillar.Type),
	}
}

func (Board) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entutils.ElkSecurity,
	}
}
