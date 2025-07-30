package config

import (
	"os"
	"time"
)

// 애플리케이션 설정
var (
	// JWT 설정
	JWTSecret       = getEnvOrDefault("JWT_SECRET", "your-secret-key-change-in-production")
	JWTExpiration   = 24 * time.Hour     // 토큰 만료 시간: 24시간
	RefreshDuration = 7 * 24 * time.Hour // 리프레시 토큰 만료 시간: 7일

	// Redis 설정
	RedisAddr     = getEnvOrDefault("REDIS_ADDR", "localhost:6379")
	RedisPassword = getEnvOrDefault("REDIS_PASSWORD", "")
	RedisDB       = 0
)

// getEnvOrDefault는 환경 변수 값을 가져오거나 기본값을 반환합니다.
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
