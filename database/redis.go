package database

import (
	"context"
	"log"
	"time"

	"github.com/dev-hyunsang/home-library/config"
	"github.com/go-redis/redis/v8"
)

// Redis 클라이언트
var RedisClient *redis.Client

// InitRedis는 Redis 연결을 초기화합니다.
func InitRedis() {
	// Redis 클라이언트 생성
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	// Redis 연결 확인
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis 연결에 실패했습니다: %v", err)
	}

	log.Println("Redis 연결에 성공했습니다")
}

// SetJWTToken은 JWT 토큰을 Redis에 저장합니다.
func SetJWTToken(ctx context.Context, userID string, tokenType string, token string, expiration time.Duration) error {
	// 키 형식: "jwt:{userID}:{tokenType}" (예: "jwt:123:access" 또는 "jwt:123:refresh")
	key := "jwt:" + userID + ":" + tokenType

	// Redis에 토큰 저장
	return RedisClient.Set(ctx, key, token, expiration).Err()
}

// GetJWTToken은 Redis에서 JWT 토큰을 가져옵니다.
func GetJWTToken(ctx context.Context, userID string, tokenType string) (string, error) {
	// 키 형식: "jwt:{userID}:{tokenType}"
	key := "jwt:" + userID + ":" + tokenType

	// Redis에서 토큰 가져오기
	return RedisClient.Get(ctx, key).Result()
}

// DeleteJWTToken은 Redis에서 JWT 토큰을 삭제합니다.
func DeleteJWTToken(ctx context.Context, userID string, tokenType string) error {
	// 키 형식: "jwt:{userID}:{tokenType}"
	key := "jwt:" + userID + ":" + tokenType

	// Redis에서 토큰 삭제
	return RedisClient.Del(ctx, key).Err()
}

// DeleteAllUserTokens는 사용자의 모든 토큰을 삭제합니다.
func DeleteAllUserTokens(ctx context.Context, userID string) error {
	// 패턴: "jwt:{userID}:*"
	pattern := "jwt:" + userID + ":*"

	// 패턴과 일치하는 모든 키 검색
	keys, err := RedisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	// 키가 없으면 완료
	if len(keys) == 0 {
		return nil
	}

	// 모든 키 삭제
	return RedisClient.Del(ctx, keys...).Err()
}
