package domain

import (
	"time"

	"github.com/dev-hyunsang/home-library/internal/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	NickName    string    `json:"nick_name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserRepository interface {
	Save(user *User) (*User, error)
	GetByID(id uuid.UUID) (*User, error)
	GetByEmail(email string) (*User, error)
	Edit(user *User) error
	Delete(id uuid.UUID) error
}

type UserUseCase interface {
	CreateUser(user *User) (*User, error)
	GetByID(id uuid.UUID) (*User, error)
	GetByEmail(email string) (*User, error)
	Edit(user *User) error
	Delete(id uuid.UUID) error
}

type AuthUseCase interface {
	GenerateToken(userID uuid.UUID) (string, error)
	GenerateTokenPair(userID uuid.UUID) (accessToken, refreshToken string, err error)
	ValidateToken(tokenString string) (*auth.JWTClaims, error)
	GetUserIDFromToken(ctx *fiber.Ctx) (uuid.UUID, error)
	ExtractTokenFromHeader(ctx *fiber.Ctx) (string, error)
	RefreshToken(refreshToken string) (newAccessToken, newRefreshToken string, err error)
	InvalidateToken(token string) error
	InvalidateAllUserTokens(userID uuid.UUID) error
	CheckRateLimit(userID uuid.UUID, action string, limit int, window time.Duration) (bool, error)
	// Legacy methods for backward compatibility
	SetSession(userID string, ctx *fiber.Ctx) error
	GetSessionByID(userID string, ctx *fiber.Ctx) (string, error)
	DeleteSession(ctx *fiber.Ctx) error
}

type AuthRepository interface {
	GenerateToken(userID uuid.UUID) (string, error)
	GenerateTokenPair(userID uuid.UUID) (accessToken, refreshToken string, err error)
	ValidateToken(tokenString string) (*auth.JWTClaims, error)
	GetUserIDFromToken(ctx *fiber.Ctx) (uuid.UUID, error)
	ExtractTokenFromHeader(ctx *fiber.Ctx) (string, error)
	RefreshToken(refreshToken string) (newAccessToken, newRefreshToken string, err error)
	InvalidateToken(token string) error
	InvalidateAllUserTokens(userID uuid.UUID) error
	CheckRateLimit(userID uuid.UUID, action string, limit int, window time.Duration) (bool, error)
	// Legacy methods for backward compatibility
	SetSession(userID string, ctx *fiber.Ctx) error
	GetSessionByID(userID string, ctx *fiber.Ctx) (string, error)
	DeleteSession(ctx *fiber.Ctx) error
}
