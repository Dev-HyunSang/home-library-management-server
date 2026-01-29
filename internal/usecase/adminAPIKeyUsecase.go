package usecase

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/dev-hyunsang/home-library/internal/domain"
	repository "github.com/dev-hyunsang/home-library/internal/repository/mysql"
	"github.com/google/uuid"
)

type adminAPIKeyUseCase struct {
	apiKeyRepo domain.AdminAPIKeyRepository
}

func NewAdminAPIKeyUseCase(apiKeyRepo *repository.AdminAPIKeyRepository) *adminAPIKeyUseCase {
	return &adminAPIKeyUseCase{apiKeyRepo: apiKeyRepo}
}

const (
	keyLength = 20
	charset   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
)

func generateRandomKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i := range bytes {
		bytes[i] = charset[int(bytes[i])%len(charset)]
	}

	return string(bytes), nil
}

func hashKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

func (uc *adminAPIKeyUseCase) CreateAPIKey(req *domain.CreateAPIKeyRequest) (*domain.AdminAPIKeyWithRawKey, error) {
	if req.Name == "" {
		return nil, domain.ErrInvalidInput
	}

	rawKey, err := generateRandomKey(keyLength)
	if err != nil {
		return nil, err
	}

	keyHash := hashKey(rawKey)
	keyPrefix := rawKey[:8]

	apiKey := &domain.AdminAPIKey{
		ID:        uuid.New(),
		Name:      req.Name,
		KeyHash:   keyHash,
		KeyPrefix: keyPrefix,
		IsActive:  true,
		ExpiresAt: req.ExpiresAt,
	}

	created, err := uc.apiKeyRepo.Create(apiKey)
	if err != nil {
		return nil, err
	}

	return &domain.AdminAPIKeyWithRawKey{
		AdminAPIKey: *created,
		RawKey:      rawKey,
	}, nil
}

func (uc *adminAPIKeyUseCase) ValidateAPIKey(rawKey string) (*domain.AdminAPIKey, error) {
	if rawKey == "" {
		return nil, domain.ErrInvalidToken
	}

	keyHash := hashKey(rawKey)
	apiKey, err := uc.apiKeyRepo.GetByKeyHash(keyHash)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	if !apiKey.IsActive {
		return nil, domain.ErrInvalidToken
	}

	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		return nil, domain.ErrTokenExpired
	}

	go uc.apiKeyRepo.UpdateLastUsed(apiKey.ID)

	return apiKey, nil
}

func (uc *adminAPIKeyUseCase) GetAllAPIKeys() ([]*domain.AdminAPIKey, error) {
	return uc.apiKeyRepo.GetAll()
}

func (uc *adminAPIKeyUseCase) DeactivateAPIKey(id uuid.UUID) error {
	return uc.apiKeyRepo.Deactivate(id)
}

func (uc *adminAPIKeyUseCase) DeleteAPIKey(id uuid.UUID) error {
	return uc.apiKeyRepo.Delete(id)
}
