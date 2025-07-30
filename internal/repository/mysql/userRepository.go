package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/dev-hyunsang/home-library/internal/domain"
	"github.com/dev-hyunsang/home-library/lib/ent"
	"github.com/dev-hyunsang/home-library/lib/ent/user"
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
		return nil, fmt.Errorf("failed to generate user id(uuid): %w", err)
	}

	// 평문 사용자 비밀번호를 해쉬화하여 데이터베이스에 저장합니다.
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	u, err := client.User.Create().
		SetID(userID).
		SetNickName(user.NickName).
		SetEmail(user.Email).
		SetPassword(string(hashedPw)). // 절대 평문 비밀번호를 저장하지 마시오.
		SetUpdatedAt(time.Now()).
		SetCreatedAt(time.Now()).
		Save(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return &domain.User{
		ID:        u.ID,
		NickName:  u.NickName,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

func (r *UserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	client := r.client

	u, err := client.User.Query().
		Where(user.ID(id)).
		Only(context.Background())
	if err == nil {
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
		return nil, fmt.Errorf("not found user: %w", err)
	default:
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	client := r.client

	u, err := client.User.Query().
		Where(user.Email(email)).
		Only(context.Background())
	if err == nil {
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
		return nil, fmt.Errorf("not found user with email %s: %w", email, err)
	default:
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
}

func (r *UserRepository) GetAll() ([]domain.User, error) {
	client := r.client

	users, err := client.User.Query().All(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	var userList []domain.User
	for _, u := range users {
		userList = append(userList, domain.User{
			ID:       u.ID,
			NickName: u.NickName,
			Email:    u.Email,
			Password: u.Password,
		})
	}

	return userList, nil
}

func (r *UserRepository) Edit(user *domain.User) (*domain.User, error) {
	client := r.client

	hasedPw, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	err = client.User.UpdateOneID(user.ID).
		SetEmail(user.Email).
		SetNickName(user.NickName).
		SetPassword(string(hasedPw)).
		SetUpdatedAt(time.Now()).
		Exec(context.Background())

	// 해당 되는 ID로 조회한 결과가 없는 경우
	switch {
	case ent.IsNotFound(err):
		return nil, fmt.Errorf("user not found")
	case err != nil:
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// 업데이트된 사용자 정보를 반환
	return &domain.User{
		ID:        user.ID,
		NickName:  user.NickName,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: time.Now(),
	}, nil
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	client := r.client

	err := client.User.DeleteOneID(id).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
