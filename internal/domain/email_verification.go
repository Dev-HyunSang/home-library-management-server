package domain

import (
	"time"

	"github.com/google/uuid"
)

type EmailVerification struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Code       string    `json:"code"`
	ExpiresAt  time.Time `json:"expires_at"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}

type EmailVerificationRepository interface {
	Save(verification *EmailVerification) (*EmailVerification, error)
	GetByEmailAndCode(email, code string) (*EmailVerification, error)
	GetLatestByEmail(email string) (*EmailVerification, error)
	MarkAsVerified(id uuid.UUID) error
	DeleteByEmail(email string) error
	DeleteExpired() error
}

type EmailVerificationUseCase interface {
	SendVerificationCode(email string) error
	VerifyCode(email, code string) error
	IsEmailVerified(email string) (bool, error)
}
