package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"

	entlocal "github.com/romashorodok/infosec/ent"
)

// Grant holds the schema definition for the Grant entity.
type Grant struct {
	ent.Schema
}

// Fields of the Grant.
func (Grant) Fields() []ent.Field {
	return nil
}

// Edges of the Grant.
func (Grant) Edges() []ent.Edge {
	return nil
}

func (Grant) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entlocal.ElkSecurity,
	}
}
