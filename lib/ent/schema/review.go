package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Review holds the schema definition for the Review entity.
type Review struct {
	ent.Schema
}

// Fields of the Review.
func (Review) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("id"),
		field.Text("content").
			NotEmpty(),
		field.Int("rating").
			Min(1).
			Max(5),
		field.Bool("is_public").
			Default(false),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Review.
func (Review) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("reviews").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.From("book", Book.Type).
			Ref("reviews").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
