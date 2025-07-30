package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()),
		field.String("nick_name").
			NotEmpty().
			Comment("사용자 닉네임"),
		field.String("email").
			NotEmpty().
			Unique().
			Comment("사용자 이메일"),
		field.String("password").
			NotEmpty().
			Comment("사용자 비밀번호"),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			Comment("사용자 생성 시간"),
		field.Time("updated_at").
			Default(time.Now).
			Comment("사용자 수정 시간"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
