package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Book holds the schema definition for the Book entity.
type Book struct {
	ent.Schema
}

// Fields of the Book.
func (Book) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()),
		field.String("book_title").
			NotEmpty(),
		field.String("author").
			NotEmpty(),
		field.String("book_isbn").
			Optional(),
		field.Time("registered_at").
			Default(time.Now()),
		field.Time("complated_at").
			Default(time.Time{}),
	}
}

// Edges of the Book.
func (Book) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("books").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("reviews", Review.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
