package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User 모델 정의
type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Email     string         `gorm:"size:100;not null;unique" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"` // 패스워드는 JSON 응답에서 제외
	FirstName string         `gorm:"size:100" json:"first_name"`
	LastName  string         `gorm:"size:100" json:"last_name"`
}

// Book 모델 정의
type Book struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	Author      string         `gorm:"size:255;not null" json:"author"`
	ISBN        string         `gorm:"size:20;unique" json:"isbn"`
	PublishedAt time.Time      `json:"published_at"`
	UserID      uuid.UUID      `gorm:"type:uuid" json:"user_id"` // User의 UUID를 참조
}

// BeforeCreate는 레코드 생성 전에 UUID를 자동으로 생성합니다.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}

// BeforeCreate는 레코드 생성 전에 UUID를 자동으로 생성합니다.
func (b *Book) BeforeCreate(tx *gorm.DB) error {
	b.ID = uuid.New()
	return nil
}

// MigrateModels는 모델을 데이터베이스에 마이그레이션합니다.
func MigrateModels(db *gorm.DB) {
	// 자동 마이그레이션 실행
	err := db.AutoMigrate(&User{}, &Book{})
	if err != nil {
		panic("데이터베이스 마이그레이션에 실패했습니다: " + err.Error())
	}
}
