package usecase

import (
	repository "github.com/dev-hyunsang/home-library/internal/repository/redis"
	"github.com/gofiber/fiber/v2"
)

type AuthUseCase struct {
	authRepo repository.AuthRepository
}

func NewAuthUseCase(repo *repository.AuthRepository) *AuthUseCase {
	return &AuthUseCase{authRepo: *repo}
}

func (uc *AuthUseCase) SetSession(userID string, ctx *fiber.Ctx) error {
	return uc.authRepo.SetSession(userID, ctx)
}

func (uc *AuthUseCase) GetSessionByID(userID string, ctx *fiber.Ctx) (string, error) {
	return uc.authRepo.GetSessionByID(userID, ctx)
}

func (uc *AuthUseCase) DeleteSession(ctx *fiber.Ctx) error {
	return uc.authRepo.DeleteSession(ctx)
}
