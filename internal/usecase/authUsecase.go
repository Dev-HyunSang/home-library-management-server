package usecase

import (
	"time"

	"github.com/dev-hyunsang/home-library/internal/auth"
	repository "github.com/dev-hyunsang/home-library/internal/repository/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthUseCase struct {
	authRepo *repository.AuthRepository
}

func NewAuthUseCase(repo *repository.AuthRepository) *AuthUseCase {
	return &AuthUseCase{authRepo: repo}
}

func (uc *AuthUseCase) GenerateToken(userID uuid.UUID) (string, error) {
	return uc.authRepo.GenerateToken(userID)
}

func (uc *AuthUseCase) GenerateTokenPair(userID uuid.UUID) (accessToken, refreshToken string, err error) {
	return uc.authRepo.GenerateTokenPair(userID)
}

func (uc *AuthUseCase) RefreshToken(refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	return uc.authRepo.RefreshToken(refreshToken)
}

func (uc *AuthUseCase) InvalidateToken(token string) error {
	return uc.authRepo.InvalidateToken(token)
}

func (uc *AuthUseCase) InvalidateAllUserTokens(userID uuid.UUID) error {
	return uc.authRepo.InvalidateAllUserTokens(userID)
}

func (uc *AuthUseCase) CheckRateLimit(userID uuid.UUID, action string, limit int, window time.Duration) (bool, error) {
	return uc.authRepo.CheckRateLimit(userID, action, limit, window)
}

func (uc *AuthUseCase) ValidateToken(tokenString string) (*auth.JWTClaims, error) {
	return uc.authRepo.ValidateToken(tokenString)
}

func (uc *AuthUseCase) GetUserIDFromToken(ctx *fiber.Ctx) (uuid.UUID, error) {
	return uc.authRepo.GetUserIDFromToken(ctx)
}

func (uc *AuthUseCase) ExtractTokenFromHeader(ctx *fiber.Ctx) (string, error) {
	return uc.authRepo.ExtractTokenFromHeader(ctx)
}

// Legacy methods for backward compatibility
func (uc *AuthUseCase) SetSession(userID string, ctx *fiber.Ctx) error {
	return uc.authRepo.SetSession(userID, ctx)
}

func (uc *AuthUseCase) GetSessionByID(userID string, ctx *fiber.Ctx) (string, error) {
	return uc.authRepo.GetSessionByID(userID, ctx)
}

func (uc *AuthUseCase) DeleteSession(ctx *fiber.Ctx) error {
	return uc.authRepo.DeleteSession(ctx)
}
