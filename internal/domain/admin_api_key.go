package domain

import (
	"time"

	"github.com/google/uuid"
)

type AdminAPIKey struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	KeyHash    string     `json:"-"`
	KeyPrefix  string     `json:"key_prefix"`
	IsActive   bool       `json:"is_active"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type AdminAPIKeyWithRawKey struct {
	AdminAPIKey
	RawKey string `json:"raw_key"`
}

type CreateAPIKeyRequest struct {
	Name      string     `json:"name"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type AdminAPIKeyRepository interface {
	Create(apiKey *AdminAPIKey) (*AdminAPIKey, error)
	GetByID(id uuid.UUID) (*AdminAPIKey, error)
	GetByKeyHash(keyHash string) (*AdminAPIKey, error)
	GetAll() ([]*AdminAPIKey, error)
	UpdateLastUsed(id uuid.UUID) error
	Deactivate(id uuid.UUID) error
	Delete(id uuid.UUID) error
}

type AdminAPIKeyUseCase interface {
	CreateAPIKey(req *CreateAPIKeyRequest) (*AdminAPIKeyWithRawKey, error)
	ValidateAPIKey(rawKey string) (*AdminAPIKey, error)
	GetAllAPIKeys() ([]*AdminAPIKey, error)
	DeactivateAPIKey(id uuid.UUID) error
	DeleteAPIKey(id uuid.UUID) error
}
