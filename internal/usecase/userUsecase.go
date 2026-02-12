package usecase

import (
	"github.com/dev-hyunsang/my-own-library-backend/internal/domain"
	repository "github.com/dev-hyunsang/my-own-library-backend/internal/repository/mysql"
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

func (uc *userUseCase) Save(user *domain.User) (*domain.User, error) {
	// 필수 필드 검사 (비밀번호는 OAuth 로그인 시 없을 수 있음)
	if user.NickName == "" || user.Email == "" {
		return nil, domain.ErrInvalidInput
	}

	// Provider가 없는 경우(일반 회원가입)에는 비밀번호가 필수
	if user.Password == "" {
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

func (uc *userUseCase) GetByNickname(nickname string) (*domain.User, error) {
	if nickname == "" {
		return nil, domain.ErrInvalidInput
	}
	return uc.userRepo.GetByNickname(nickname)
}

func (uc *userUseCase) Update(user *domain.User) error {
	return uc.userRepo.Update(user)

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

func (uc *userUseCase) UpdateFCMToken(userID uuid.UUID, fcmToken string) error {
	if fcmToken == "" {
		return domain.ErrInvalidInput
	}
	return uc.userRepo.UpdateFCMToken(userID, fcmToken)
}

func (uc *userUseCase) UpdateTimezone(userID uuid.UUID, timezone string) error {
	if timezone == "" {
		return domain.ErrInvalidInput
	}
	return uc.userRepo.UpdateTimezone(userID, timezone)
}
