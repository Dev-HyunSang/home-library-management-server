package redis

import (
	"fmt"
	"strings"
	"time"

	"github.com/dev-hyunsang/home-library/internal/auth"
	"github.com/dev-hyunsang/home-library/internal/cache"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthRepository struct {
	jwtManager      *auth.JWTManager
	securityManager *auth.SecurityManager
}

func NewAuthRepository(secretKey string, accessTTL, refreshTTL time.Duration, redisClient *cache.RedisClient, issuer, audience string) *AuthRepository {
	securityManager := auth.NewSecurityManager(redisClient)
	return &AuthRepository{
		jwtManager:      auth.NewJWTManager(secretKey, accessTTL, refreshTTL, securityManager, issuer, audience),
		securityManager: securityManager,
	}
}

func (repo *AuthRepository) GenerateTokenPair(userID uuid.UUID) (accessToken, refreshToken string, err error) {
	return repo.jwtManager.GenerateTokenPair(userID)
}

func (repo *AuthRepository) GenerateToken(userID uuid.UUID) (string, error) {
	accessToken, _, err := repo.jwtManager.GenerateTokenPair(userID)
	return accessToken, err
}

func (repo *AuthRepository) ValidateToken(tokenString string) (*auth.JWTClaims, error) {
	return repo.jwtManager.ValidateToken(tokenString)
}

func (repo *AuthRepository) ExtractTokenFromHeader(ctx *fiber.Ctx) (string, error) {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header not found")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

func (repo *AuthRepository) GetUserIDFromToken(ctx *fiber.Ctx) (uuid.UUID, error) {
	token, err := repo.ExtractTokenFromHeader(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	return repo.jwtManager.ExtractUserID(token)
}

func (repo *AuthRepository) RefreshToken(refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	return repo.jwtManager.RefreshToken(refreshToken)
}

func (repo *AuthRepository) InvalidateToken(token string) error {
	return repo.jwtManager.InvalidateToken(token)
}

func (repo *AuthRepository) InvalidateAllUserTokens(userID uuid.UUID) error {
	return repo.jwtManager.InvalidateAllUserTokens(userID)
}

func (repo *AuthRepository) CheckRateLimit(userID uuid.UUID, action string, limit int, window time.Duration) (bool, error) {
	return repo.securityManager.CheckRateLimit(userID, action, limit, window)
}

// Legacy methods for backward compatibility - these will be removed
func (repo *AuthRepository) SetSession(userID string, ctx *fiber.Ctx) error {
	return fmt.Errorf("session-based auth is deprecated, use JWT tokens")
}

func (repo *AuthRepository) GetSessionByID(userID string, ctx *fiber.Ctx) (string, error) {
	return "", fmt.Errorf("session-based auth is deprecated, use JWT tokens")
}

func (repo *AuthRepository) GetAllSession(ctx *fiber.Ctx) ([]string, error) {
	return nil, fmt.Errorf("session-based auth is deprecated, use JWT tokens")
}

func (repo *AuthRepository) DeleteSession(ctx *fiber.Ctx) error {
	return fmt.Errorf("session-based auth is deprecated, use JWT tokens")
}
