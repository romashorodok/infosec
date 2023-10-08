package schema

import "entgo.io/ent"

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
	return nil
}
