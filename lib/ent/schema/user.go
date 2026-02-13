package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
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
			Optional().
			Comment("사용자 비밀번호 (OAuth 사용자의 경우 비어있을 수 있음)"),
		field.Bool("is_published").
			Default(false).
			Comment("사용자 계정의 책 공개 여부"),
		field.Bool("is_terms_agreed").
			Default(false).
			Comment("사용자 이용약관 동의 여부"),
		field.Bool("is_privacy_agreed").
			Default(false).
			Comment("사용자 개인정보 수집 이용 동의 여부"),
		field.String("fcm_token").
			Optional().
			Comment("FCM 디바이스 토큰"),
		field.String("timezone").
			Default("Asia/Seoul").
			Comment("사용자 타임존"),
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
	return []ent.Edge{
		edge.To("books", Book.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("reviews", Review.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("bookmarks", Bookmark.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("reading_reminders", ReadingReminder.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
