package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
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
