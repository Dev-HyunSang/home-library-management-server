package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// AdminAPIKey holds the schema definition for the AdminAPIKey entity.
type AdminAPIKey struct {
	ent.Schema
}

// Fields of the AdminAPIKey.
func (AdminAPIKey) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()),
		field.String("name").
			NotEmpty().
			Comment("API Key 이름/설명"),
		field.String("key_hash").
			NotEmpty().
			Unique().
			Comment("해시된 API Key"),
		field.String("key_prefix").
			MaxLen(8).
			Comment("API Key 앞 8자리 (식별용)"),
		field.Bool("is_active").
			Default(true).
			Comment("활성화 여부"),
		field.Time("last_used_at").
			Optional().
			Nillable().
			Comment("마지막 사용 시간"),
		field.Time("expires_at").
			Optional().
			Nillable().
			Comment("만료 시간 (null이면 무제한)"),
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			Comment("생성 시간"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("수정 시간"),
	}
}

// Edges of the AdminAPIKey.
func (AdminAPIKey) Edges() []ent.Edge {
	return nil
}
