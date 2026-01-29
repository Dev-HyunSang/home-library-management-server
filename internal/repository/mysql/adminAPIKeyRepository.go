package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/dev-hyunsang/home-library/internal/domain"
	"github.com/dev-hyunsang/home-library/lib/ent"
	"github.com/dev-hyunsang/home-library/lib/ent/adminapikey"
	"github.com/dev-hyunsang/home-library/logger"
	"github.com/google/uuid"
)

type AdminAPIKeyRepository struct {
	client *ent.Client
}

func NewAdminAPIKeyRepository(client *ent.Client) *AdminAPIKeyRepository {
	return &AdminAPIKeyRepository{
		client: client,
	}
}

func (r *AdminAPIKeyRepository) Create(apiKey *domain.AdminAPIKey) (*domain.AdminAPIKey, error) {
	builder := r.client.AdminAPIKey.Create().
		SetID(apiKey.ID).
		SetName(apiKey.Name).
		SetKeyHash(apiKey.KeyHash).
		SetKeyPrefix(apiKey.KeyPrefix).
		SetIsActive(apiKey.IsActive)

	if apiKey.ExpiresAt != nil {
		builder.SetExpiresAt(*apiKey.ExpiresAt)
	}

	key, err := builder.Save(context.Background())
	if err != nil {
		logger.Init().Sugar().Errorf("Admin API Key 생성 중 오류: %v", err)
		return nil, fmt.Errorf("API Key 생성 중 오류가 발생했습니다: %w", err)
	}

	logger.Init().Sugar().Infof("Admin API Key가 생성되었습니다. ID: %s, Name: %s", key.ID.String(), key.Name)

	return r.entToDomain(key), nil
}

func (r *AdminAPIKeyRepository) GetByID(id uuid.UUID) (*domain.AdminAPIKey, error) {
	key, err := r.client.AdminAPIKey.Query().
		Where(adminapikey.ID(id)).
		Only(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("API Key 조회 중 오류가 발생했습니다: %w", err)
	}

	return r.entToDomain(key), nil
}

func (r *AdminAPIKeyRepository) GetByKeyHash(keyHash string) (*domain.AdminAPIKey, error) {
	key, err := r.client.AdminAPIKey.Query().
		Where(
			adminapikey.KeyHash(keyHash),
			adminapikey.IsActive(true),
		).
		Only(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("API Key 조회 중 오류가 발생했습니다: %w", err)
	}

	return r.entToDomain(key), nil
}

func (r *AdminAPIKeyRepository) GetAll() ([]*domain.AdminAPIKey, error) {
	keys, err := r.client.AdminAPIKey.Query().
		Order(ent.Desc(adminapikey.FieldCreatedAt)).
		All(context.Background())

	if err != nil {
		return nil, fmt.Errorf("API Key 목록 조회 중 오류가 발생했습니다: %w", err)
	}

	result := make([]*domain.AdminAPIKey, len(keys))
	for i, key := range keys {
		result[i] = r.entToDomain(key)
	}

	return result, nil
}

func (r *AdminAPIKeyRepository) UpdateLastUsed(id uuid.UUID) error {
	err := r.client.AdminAPIKey.UpdateOneID(id).
		SetLastUsedAt(time.Now()).
		Exec(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("API Key 사용 시간 업데이트 중 오류가 발생했습니다: %w", err)
	}

	return nil
}

func (r *AdminAPIKeyRepository) Deactivate(id uuid.UUID) error {
	err := r.client.AdminAPIKey.UpdateOneID(id).
		SetIsActive(false).
		Exec(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("API Key 비활성화 중 오류가 발생했습니다: %w", err)
	}

	logger.Init().Sugar().Infof("Admin API Key가 비활성화되었습니다. ID: %s", id.String())
	return nil
}

func (r *AdminAPIKeyRepository) Delete(id uuid.UUID) error {
	err := r.client.AdminAPIKey.DeleteOneID(id).Exec(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("API Key 삭제 중 오류가 발생했습니다: %w", err)
	}

	logger.Init().Sugar().Infof("Admin API Key가 삭제되었습니다. ID: %s", id.String())
	return nil
}

func (r *AdminAPIKeyRepository) entToDomain(key *ent.AdminAPIKey) *domain.AdminAPIKey {
	result := &domain.AdminAPIKey{
		ID:        key.ID,
		Name:      key.Name,
		KeyHash:   key.KeyHash,
		KeyPrefix: key.KeyPrefix,
		IsActive:  key.IsActive,
		CreatedAt: key.CreatedAt,
		UpdatedAt: key.UpdatedAt,
	}

	if key.LastUsedAt != nil {
		result.LastUsedAt = key.LastUsedAt
	}
	if key.ExpiresAt != nil {
		result.ExpiresAt = key.ExpiresAt
	}

	return result
}
