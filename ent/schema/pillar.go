package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"

	entlocal "github.com/romashorodok/infosec/ent"
)

// Pillar holds the schema definition for the Pillar entity.
type Pillar struct {
	ent.Schema
}

// Fields of the Pillar.
func (Pillar) Fields() []ent.Field {
	return nil
}

// Edges of the Pillar.
func (Pillar) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tasks", Task.Type),
	}
}

func (Pillar) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entlocal.ElkSecurity,
	}
}
