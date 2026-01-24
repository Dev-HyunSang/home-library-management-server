package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ReadingReminder holds the schema definition for the ReadingReminder entity.
type ReadingReminder struct {
	ent.Schema
}

// Fields of the ReadingReminder.
func (ReadingReminder) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()),
		field.String("reminder_time").
			NotEmpty().
			Comment("알림 시간 (HH:MM 형식)"),
		field.Enum("day_of_week").
			Values("everyday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday").
			Default("everyday").
			Comment("알림 요일"),
		field.Bool("is_enabled").
			Default(true).
			Comment("알림 활성화 여부"),
		field.String("message").
			Default("책 읽을 시간이에요!").
			Comment("알림 메시지"),
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

// Edges of the ReadingReminder.
func (ReadingReminder) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("reading_reminders").
			Unique().
			Required().
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
