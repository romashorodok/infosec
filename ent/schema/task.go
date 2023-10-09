package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	entlocal "github.com/romashorodok/infosec/ent"
)

// Task holds the schema definition for the Task entity.
type Task struct {
	ent.Schema
}

// Fields of the Task.
func (Task) Fields() []ent.Field {
	return []ent.Field{
		field.Text("Description"),
	}
}

// Edges of the Task.
func (Task) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("participants", Participant.Type),
	}
}

func (Task) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entlocal.ElkSecurity,
	}
}
