package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dev-hyunsang/my-own-library-backend/internal/cache"
	"github.com/dev-hyunsang/my-own-library-backend/internal/domain"
	"github.com/dev-hyunsang/my-own-library-backend/logger"
	"github.com/google/uuid"
)

const (
	emailVerificationPrefix   = "email_verification:email:"
	emailVerificationIDPrefix = "email_verification:id:"
)

type EmailVerificationRepository struct {
	redisClient *cache.RedisClient
}

func NewEmailVerificationRepository(redisClient *cache.RedisClient) *EmailVerificationRepository {
	return &EmailVerificationRepository{
		redisClient: redisClient,
	}
}

func (r *EmailVerificationRepository) Save(verification *domain.EmailVerification) (*domain.EmailVerification, error) {
	verification.CreatedAt = time.Now()
	verification.IsVerified = false

	data, err := json.Marshal(verification)
	if err != nil {
		return nil, fmt.Errorf("인증 정보 직렬화 중 오류가 발생했습니다: %w", err)
	}

	ttl := time.Until(verification.ExpiresAt)
	if ttl <= 0 {
		return nil, fmt.Errorf("만료 시간이 유효하지 않습니다")
	}

	emailKey := emailVerificationPrefix + verification.Email
	if err := r.redisClient.Set(emailKey, string(data), ttl); err != nil {
		return nil, fmt.Errorf("이메일 인증 정보를 저장하는 도중 오류가 발생했습니다: %w", err)
	}

	idKey := emailVerificationIDPrefix + verification.ID.String()
	if err := r.redisClient.Set(idKey, verification.Email, ttl); err != nil {
		return nil, fmt.Errorf("인증 ID 인덱스 저장 중 오류가 발생했습니다: %w", err)
	}

	logger.Sugar().Infof("이메일 인증 정보를 생성했습니다. 이메일: %s", verification.Email)
	return verification, nil
}

func (r *EmailVerificationRepository) GetByEmailAndCode(email, code string) (*domain.EmailVerification, error) {
	emailKey := emailVerificationPrefix + email
	data, err := r.redisClient.Get(emailKey)
	if err != nil {
		return nil, fmt.Errorf("유효한 인증 정보를 찾을 수 없습니다: %w", err)
	}

	var verification domain.EmailVerification
	if err := json.Unmarshal([]byte(data), &verification); err != nil {
		return nil, fmt.Errorf("인증 정보 역직렬화 중 오류가 발생했습니다: %w", err)
	}

	if verification.Code != code {
		return nil, fmt.Errorf("유효한 인증 정보를 찾을 수 없습니다")
	}

	if verification.IsVerified {
		return nil, fmt.Errorf("이미 인증이 완료된 코드입니다")
	}

	if time.Now().After(verification.ExpiresAt) {
		return nil, fmt.Errorf("인증 코드가 만료되었습니다")
	}

	return &verification, nil
}

func (r *EmailVerificationRepository) GetLatestByEmail(email string) (*domain.EmailVerification, error) {
	emailKey := emailVerificationPrefix + email
	data, err := r.redisClient.Get(emailKey)
	if err != nil {
		return nil, fmt.Errorf("해당 이메일의 인증 정보를 찾을 수 없습니다: %w", err)
	}

	var verification domain.EmailVerification
	if err := json.Unmarshal([]byte(data), &verification); err != nil {
		return nil, fmt.Errorf("인증 정보 역직렬화 중 오류가 발생했습니다: %w", err)
	}

	return &verification, nil
}

func (r *EmailVerificationRepository) MarkAsVerified(id uuid.UUID) error {
	idKey := emailVerificationIDPrefix + id.String()
	email, err := r.redisClient.Get(idKey)
	if err != nil {
		return fmt.Errorf("해당하는 인증 정보를 찾을 수 없습니다: %w", err)
	}

	emailKey := emailVerificationPrefix + email
	data, err := r.redisClient.Get(emailKey)
	if err != nil {
		return fmt.Errorf("인증 정보를 찾을 수 없습니다: %w", err)
	}

	var verification domain.EmailVerification
	if err := json.Unmarshal([]byte(data), &verification); err != nil {
		return fmt.Errorf("인증 정보 역직렬화 중 오류가 발생했습니다: %w", err)
	}

	verification.IsVerified = true

	updatedData, err := json.Marshal(verification)
	if err != nil {
		return fmt.Errorf("인증 정보 직렬화 중 오류가 발생했습니다: %w", err)
	}

	ttl := time.Until(verification.ExpiresAt)
	if ttl <= 0 {
		ttl = time.Minute
	}

	if err := r.redisClient.Set(emailKey, string(updatedData), ttl); err != nil {
		return fmt.Errorf("인증 완료 처리 중 오류가 발생했습니다: %w", err)
	}

	logger.Sugar().Infof("이메일 인증이 완료되었습니다. ID: %s", id.String())
	return nil
}

func (r *EmailVerificationRepository) DeleteByEmail(email string) error {
	emailKey := emailVerificationPrefix + email
	data, err := r.redisClient.Get(emailKey)
	if err == nil {
		var verification domain.EmailVerification
		if err := json.Unmarshal([]byte(data), &verification); err == nil {
			idKey := emailVerificationIDPrefix + verification.ID.String()
			_ = r.redisClient.Delete(idKey)
		}
	}

	if err := r.redisClient.Delete(emailKey); err != nil {
		return fmt.Errorf("이메일 인증 정보 삭제 중 오류가 발생했습니다: %w", err)
	}

	logger.Sugar().Infof("이메일 인증 정보를 삭제했습니다. 이메일: %s", email)
	return nil
}

func (r *EmailVerificationRepository) DeleteExpired() error {
	// Redis TTL이 자동으로 만료된 키를 삭제하므로 별도 처리 불필요
	return nil
}
