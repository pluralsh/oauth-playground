package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Organization holds the schema definition for the Organization entity.
type Organization struct {
	ent.Schema
}

// Fields of the Organization.
func (Organization) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).Immutable().Unique(),
		field.String("name").Unique(),
	}
}

// Edges of the Organization.
func (Organization) Edges() []ent.Edge {
	return nil
}
