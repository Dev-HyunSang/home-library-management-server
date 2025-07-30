package domain

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	NickName  string    `json:"nick_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	Save(user *User) (*User, error)
	GetByID(id uuid.UUID) (*User, error)
	GetAll() ([]User, error)
	Edit(user *User) (*User, error)
	Delete(id uuid.UUID) error
}

type UserUseCase interface {
	CreateUser(user *User) (*User, error)
	GetByID(id uuid.UUID) (*User, error)
	GetAll() ([]User, error)
	Edit(user *User) (*User, error)
	Delete(id uuid.UUID) error
}

type AuthUseCase interface {
	SetSession(userID string, ctx *fiber.Ctx) error
	GetSessionByID(userID string, ctx *fiber.Ctx) (string, error)
	DeleteSession(ctx *fiber.Ctx) error
}

type AuthRepository interface {
	SetSession(userID string, ctx *fiber.Ctx) error
	GetSessionByID(userID string, ctx *fiber.Ctx) (string, error)
	DeleteSession(ctx *fiber.Ctx) error
}
