package usecase

import (
	"github.com/dev-hyunsang/home-library/internal/domain"
	repository "github.com/dev-hyunsang/home-library/internal/repository/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type userUseCase struct {
	userRepo domain.UserRepository
	authRepo domain.AuthUseCase
}

func NewUserUseCase(userRepo *repository.UserRepository, authUseCase domain.AuthUseCase) *userUseCase {
	return &userUseCase{userRepo: userRepo, authRepo: authUseCase}
}

func ErrResponse(err error) map[string]string {
	return map[string]string{"error": err.Error()}
}

func (uc *userUseCase) CreateUser(user *domain.User) (*domain.User, error) {
	if user.NickName == "" || user.Email == "" || user.Password == "" {
		return nil, domain.ErrInvalidInput
	}

	return uc.userRepo.Save(user)
}

func (uc *userUseCase) GetByID(id uuid.UUID) (*domain.User, error) {
	return uc.userRepo.GetByID(id)
}

func (uc *userUseCase) GetByEmail(email string) (*domain.User, error) {
	if email == "" {
		return nil, domain.ErrInvalidInput
	}

	return uc.userRepo.GetByEmail(email)
}

func (uc *userUseCase) Edit(user *domain.User) error {
	return uc.userRepo.Edit(user)

}

func (uc *userUseCase) Delete(id uuid.UUID) error {
	return uc.userRepo.Delete(id)
}

func (uc *userUseCase) SetSession(userID string, ctx *fiber.Ctx) error {
	return uc.authRepo.SetSession(userID, ctx)
}

func (uc *userUseCase) GetSessionByID(userID string, ctx *fiber.Ctx) (string, error) {
	return uc.authRepo.GetSessionByID(userID, ctx)
}

func (uc *userUseCase) DeleteSession(ctx *fiber.Ctx) error {
	return uc.authRepo.DeleteSession(ctx)
}
