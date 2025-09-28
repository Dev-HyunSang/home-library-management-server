package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/dev-hyunsang/home-library/internal/domain"
	"github.com/dev-hyunsang/home-library/lib/ent"
	"github.com/dev-hyunsang/home-library/lib/ent/book"
	"github.com/dev-hyunsang/home-library/lib/ent/user"
	"github.com/dev-hyunsang/home-library/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	client *ent.Client
	store  *session.Store
}

func NewUserRepository(client *ent.Client, store *session.Store) *UserRepository {
	return &UserRepository{
		client: client,
		store:  store,
	}
}

func (r *UserRepository) Save(user *domain.User) (*domain.User, error) {
	client := r.client

	// Create User ID(UUID)
	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("사용자의 UUID를 생성하던 도중 오류가 발생했습니다: %w", err)
	}

	// 평문 사용자 비밀번호를 해쉬화하여 데이터베이스에 저장합니다.
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("사용자의 암호를 안전하게 해쉬하던 도중 오류가 발생했습니다: %w", err)
	}

	u, err := client.User.Create().
		SetID(userID).
		SetNickName(user.NickName).
		SetEmail(user.Email).
		SetPassword(string(hashedPw)).    // 절대 평문 비밀번호를 저장하지 마시오.
		SetIsPublished(user.IsPublished). // 기본값은 비공개인 false로 설정
		SetUpdatedAt(time.Now()).
		SetCreatedAt(time.Now()).
		Save(context.Background())
	if err == nil {
		logger.Init().Sugar().Infof("새로운 유저를 생성하였습니다. 새로운 유저: %s", u.ID.String())
		return &domain.User{
			ID:          u.ID,
			NickName:    u.NickName,
			Email:       u.Email,
			Password:    u.Password,
			IsPublished: u.IsPublished,
			CreatedAt:   u.CreatedAt,
			UpdatedAt:   u.UpdatedAt,
		}, nil
	}

	switch {
	case ent.IsConstraintError(err):
		return nil, fmt.Errorf("저장 도중 제약조건 관련 오류가 발생 했습니다.: %w", err)
	default:
		return nil, fmt.Errorf("사용자 정보를 저장하는 도중 알 수 없는 오류가 발생했습니다.: %w", err)
	}
}

func (r *UserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	client := r.client

	u, err := client.User.Query().
		Where(user.ID(id)).
		Only(context.Background())
	if err == nil {
		logger.Init().Sugar().Infof("사용자 정보를 ID로 조회했습니다. 사용자ID: %s", u.ID.String())
		return &domain.User{
			ID:        u.ID,
			NickName:  u.NickName,
			Email:     u.Email,
			Password:  u.Password,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}, nil
	}

	switch {
	case ent.IsNotFound(err):
		return nil, fmt.Errorf("해당하는 ID로 사용자를 찾을 수 없습니다: %w", err)
	default:
		return nil, fmt.Errorf("사용자 정보를 ID로 조회하는 도중 오류가 발생했습니다: %w", err)
	}
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	client := r.client

	u, err := client.User.Query().
		Where(user.Email(email)).
		Only(context.Background())
	if err == nil {
		logger.Init().Sugar().Infof("사용자 정보를 이메일로 조회했습니다. 사용자 이메일: %s", u.Email)
		return &domain.User{
			ID:        u.ID,
			NickName:  u.NickName,
			Email:     u.Email,
			Password:  u.Password,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}, nil
	}

	switch {
	case ent.IsNotFound(err):
		return nil, fmt.Errorf("해당하는 이메일로 사용자를 찾을 수 없습니다: %w", err)
	default:
		return nil, fmt.Errorf("사용자 정보를 이메일로 조회하는 도중 오류가 발생했습니다: %w", err)
	}
}

func (r *UserRepository) Update(user *domain.User) error {
	client := r.client

	err := client.User.UpdateOneID(user.ID).
		SetEmail(user.Email).
		SetNickName(user.NickName).
		SetPassword(user.Password).
		SetIsPublished(user.IsPublished).
		SetUpdatedAt(time.Now()).
		Exec(context.Background())

	if err != nil {
		return fmt.Errorf("사용자 정보를 업데이트하는 도중 오류가 발생했습니다: %w", err)
	} else if ent.IsNotFound(err) {
		return fmt.Errorf("해당하는 ID로 사용자를 찾을 수 없습니다: %w", err)
	}

	return nil
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	client := r.client

	// 사용자가 등록한 책을 먼저 삭제합니다.
	_, err := client.Book.Delete().Where(book.HasOwnerWith(user.ID(id))).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("해당하는 사용자의 등록된 책을 삭제하던 도중 오류가 발생했습니다: %w", err)
	}

	// 책이 성공적으로 삭제되었다면, 사용자를 삭제합니다.
	err = client.User.DeleteOneID(id).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("사용자를 삭제하는 도중 오류가 발생했습니다: %w", err)
	}

	logger.Init().Sugar().Infof("해당하는 사용자를 삭제하였습니다: %s", id.String())

	return nil
}
