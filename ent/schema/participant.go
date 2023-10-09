package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"

	entlocal "github.com/romashorodok/infosec/ent"
)

// Participant holds the schema definition for the Participant entity.
type Participant struct {
	ent.Schema
}

// Fields of the Participant.
func (Participant) Fields() []ent.Field {
	return nil
}

// Edges of the Participant.
func (Participant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("boards", Board.Type).Ref("participants"),
		edge.From("tasks", Task.Type).Ref("participants"),
	}
}

func (Participant) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entlocal.ElkSecurity,
	}
}
