package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Category holds the schema definition for the Category entity.
type Category struct {
	ent.Schema
}

// Fields of the Category.
func (Category) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.String("code").Unique().NotEmpty(),
		field.Int("parent_id").Optional(),
	}
}

// Edges of the Category.
func (Category) Edges() []ent.Edge {
	// TODO: O2O relation?
	return []ent.Edge{
		edge.To("children", Category.Type).From("parent").Unique().Field("parent_id"),
	}
}
