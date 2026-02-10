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
		field.String("book_isbn").
			NotEmpty().
			Comment("ISBN of the reviewed book"),
		field.Text("content").
			NotEmpty().
			Comment("Review content"),
		field.Int("rating").
			Min(1).
			Max(5).
			Comment("Rating 1-5"),
		field.Bool("is_public").
			Default(false).
			Comment("Whether the review is public"),
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
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
