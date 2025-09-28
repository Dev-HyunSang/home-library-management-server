package auth

import (
	"fmt"
	"time"

	"github.com/dev-hyunsang/home-library/internal/cache"
	"github.com/google/uuid"
)

type SecurityManager struct {
	redisClient *cache.RedisClient
}

func NewSecurityManager(redisClient *cache.RedisClient) *SecurityManager {
	return &SecurityManager{
		redisClient: redisClient,
	}
}

// 토큰 블랙리스트 관리
func (sm *SecurityManager) BlacklistToken(tokenID string, expiration time.Duration) error {
	key := fmt.Sprintf("blacklist:token:%s", tokenID)
	return sm.redisClient.Set(key, "blacklisted", expiration)
}

func (sm *SecurityManager) IsTokenBlacklisted(tokenID string) (bool, error) {
	key := fmt.Sprintf("blacklist:token:%s", tokenID)
	return sm.redisClient.Exists(key)
}

// 사용자별 모든 토큰 무효화
func (sm *SecurityManager) InvalidateAllUserTokens(userID uuid.UUID) error {
	key := fmt.Sprintf("user:token_version:%s", userID.String())
	_, err := sm.redisClient.Incr(key)
	if err != nil {
		return err
	}

	// 버전을 24시간 동안 유지 (토큰 만료 시간과 동일)
	return sm.redisClient.Expire(key, 24*time.Hour)
}

func (sm *SecurityManager) GetUserTokenVersion(userID uuid.UUID) (int64, error) {
	key := fmt.Sprintf("user:token_version:%s", userID.String())
	versionStr, err := sm.redisClient.Get(key)
	if err != nil {
		// 키가 없으면 버전 0으로 간주
		if err.Error() == "redis: nil" {
			return 0, nil
		}
		return 0, err
	}

	var version int64
	if _, err := fmt.Sscanf(versionStr, "%d", &version); err != nil {
		return 0, err
	}

	return version, nil
}

// 리프레시 토큰 관리
func (sm *SecurityManager) StoreRefreshToken(userID uuid.UUID, tokenID string, expiration time.Duration) error {
	key := fmt.Sprintf("refresh_token:%s:%s", userID.String(), tokenID)
	return sm.redisClient.Set(key, "valid", expiration)
}

func (sm *SecurityManager) IsRefreshTokenValid(userID uuid.UUID, tokenID string) (bool, error) {
	key := fmt.Sprintf("refresh_token:%s:%s", userID.String(), tokenID)
	return sm.redisClient.Exists(key)
}

func (sm *SecurityManager) RevokeRefreshToken(userID uuid.UUID, tokenID string) error {
	key := fmt.Sprintf("refresh_token:%s:%s", userID.String(), tokenID)
	return sm.redisClient.Delete(key)
}

// Rate Limiting
func (sm *SecurityManager) CheckRateLimit(userID uuid.UUID, action string, limit int, window time.Duration) (bool, error) {
	key := fmt.Sprintf("rate_limit:%s:%s", action, userID.String())

	count, err := sm.redisClient.Incr(key)
	if err != nil {
		return false, err
	}

	if count == 1 {
		// 첫 번째 요청이면 만료 시간 설정
		if err := sm.redisClient.Expire(key, window); err != nil {
			return false, err
		}
	}

	return count <= int64(limit), nil
}

// 세션 토큰 저장 (추가 보안층)
func (sm *SecurityManager) StoreSessionToken(userID uuid.UUID, tokenID string, expiration time.Duration) error {
	key := fmt.Sprintf("session:token:%s", tokenID)
	value := fmt.Sprintf("user:%s", userID.String())
	return sm.redisClient.Set(key, value, expiration)
}

func (sm *SecurityManager) ValidateSessionToken(tokenID string) (uuid.UUID, error) {
	key := fmt.Sprintf("session:token:%s", tokenID)
	value, err := sm.redisClient.Get(key)
	if err != nil {
		return uuid.Nil, fmt.Errorf("session token not found or expired")
	}

	var userIDStr string
	if _, err := fmt.Sscanf(value, "user:%s", &userIDStr); err != nil {
		return uuid.Nil, fmt.Errorf("invalid session token format")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in session token")
	}

	return userID, nil
}

func (sm *SecurityManager) RevokeSessionToken(tokenID string) error {
	key := fmt.Sprintf("session:token:%s", tokenID)
	return sm.redisClient.Delete(key)
}