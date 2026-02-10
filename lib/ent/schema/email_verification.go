package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// EmailVerification holds the schema definition for email verification codes.
type EmailVerification struct {
	ent.Schema
}

// Fields of the EmailVerification.
func (EmailVerification) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()),
		field.String("email").
			NotEmpty().
			Comment("인증 대상 이메일 주소"),
		field.String("code").
			MaxLen(6).
			MinLen(6).
			NotEmpty().
			Comment("6자리 인증번호"),
		field.Time("expires_at").
			Comment("인증번호 만료 시간"),
		field.Bool("is_verified").
			Default(false).
			Comment("인증 완료 여부"),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			Comment("생성 시간"),
	}
}

// Edges of the EmailVerification.
func (EmailVerification) Edges() []ent.Edge {
	return nil
}
