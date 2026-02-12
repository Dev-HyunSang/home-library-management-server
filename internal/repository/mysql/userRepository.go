package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/dev-hyunsang/my-own-library-backend/internal/domain"
	"github.com/dev-hyunsang/my-own-library-backend/lib/ent"
	"github.com/dev-hyunsang/my-own-library-backend/lib/ent/book"
	"github.com/dev-hyunsang/my-own-library-backend/lib/ent/user"
	"github.com/dev-hyunsang/my-own-library-backend/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
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

	// User Create Builder 생성
	builder := client.User.Create().
		SetID(user.ID).
		SetNickName(user.NickName).
		SetEmail(user.Email).
		SetIsPublished(user.IsPublished). // 기본값은 비공개인 false로 설정
		SetIsTermsAgreed(user.IsTermsAgreed).
		SetUpdatedAt(time.Now()).
		SetCreatedAt(time.Now())

	if user.Password != "" {
		builder.SetPassword(user.Password)
	}

	u, err := builder.Save(context.Background())

	if err == nil {
		logger.Sugar().Infof("새로운 유저를 생성하였습니다. 새로운 유저: %s", u.ID.String())
		return &domain.User{
			ID:            u.ID,
			NickName:      u.NickName,
			Email:         u.Email,
			Password:      u.Password,
			IsPublished:   u.IsPublished,
			IsTermsAgreed: u.IsTermsAgreed,
			FCMToken:      u.FcmToken,
			Timezone:      u.Timezone,
			CreatedAt:     u.CreatedAt,
			UpdatedAt:     u.UpdatedAt,
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
		logger.Sugar().Infof("사용자 정보를 ID로 조회했습니다. 사용자ID: %s", u.ID.String())
		return &domain.User{
			ID:            u.ID,
			NickName:      u.NickName,
			Email:         u.Email,
			Password:      u.Password,
			IsPublished:   u.IsPublished,
			IsTermsAgreed: u.IsTermsAgreed,
			FCMToken:      u.FcmToken,
			Timezone:      u.Timezone,
			CreatedAt:     u.CreatedAt,
			UpdatedAt:     u.UpdatedAt,
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
		logger.Sugar().Infof("사용자 정보를 이메일로 조회했습니다. 사용자 이메일: %s", u.Email)
		return &domain.User{
			ID:            u.ID,
			NickName:      u.NickName,
			Email:         u.Email,
			Password:      u.Password,
			IsPublished:   u.IsPublished,
			IsTermsAgreed: u.IsTermsAgreed,
			FCMToken:      u.FcmToken,
			Timezone:      u.Timezone,
			CreatedAt:     u.CreatedAt,
			UpdatedAt:     u.UpdatedAt,
		}, nil
	}

	switch {
	case ent.IsNotFound(err):
		return nil, fmt.Errorf("해당하는 이메일로 사용자를 찾을 수 없습니다: %w", err)
	default:
		return nil, fmt.Errorf("사용자 정보를 이메일로 조회하는 도중 오류가 발생했습니다: %w", err)
	}
}

func (r *UserRepository) GetByNickname(nickname string) (*domain.User, error) {
	client := r.client

	u, err := client.User.Query().
		Where(user.NickName(nickname)).
		Only(context.Background())
	if err == nil {
		logger.Sugar().Infof("사용자 정보를 닉네임으로 조회했습니다. 사용자 닉네임: %s", u.NickName)
		return &domain.User{
			ID:            u.ID,
			NickName:      u.NickName,
			Email:         u.Email,
			Password:      u.Password,
			IsPublished:   u.IsPublished,
			IsTermsAgreed: u.IsTermsAgreed,
			FCMToken:      u.FcmToken,
			Timezone:      u.Timezone,
			CreatedAt:     u.CreatedAt,
			UpdatedAt:     u.UpdatedAt,
		}, nil
	}

	switch {
	case ent.IsNotFound(err):
		return nil, fmt.Errorf("해당하는 닉네임으로 사용자를 찾을 수 없습니다: %w", err)
	default:
		return nil, fmt.Errorf("사용자 정보를 닉네임으로 조회하는 도중 오류가 발생했습니다: %w", err)
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
		if ent.IsNotFound(err) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("사용자 정보를 업데이트하는 도중 오류가 발생했습니다: %w", err)
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

	logger.Sugar().Infof("해당하는 사용자를 삭제하였습니다: %s", id.String())

	return nil
}

func (r *UserRepository) UpdateFCMToken(userID uuid.UUID, fcmToken string) error {
	err := r.client.User.UpdateOneID(userID).
		SetFcmToken(fcmToken).
		SetUpdatedAt(time.Now()).
		Exec(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("해당하는 ID로 사용자를 찾을 수 없습니다: %w", err)
		}
		return fmt.Errorf("FCM 토큰 업데이트 중 오류가 발생했습니다: %w", err)
	}

	logger.Sugar().Infof("사용자 FCM 토큰이 업데이트되었습니다. 사용자ID: %s", userID.String())
	return nil
}

func (r *UserRepository) UpdateTimezone(userID uuid.UUID, timezone string) error {
	err := r.client.User.UpdateOneID(userID).
		SetTimezone(timezone).
		SetUpdatedAt(time.Now()).
		Exec(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("해당하는 ID로 사용자를 찾을 수 없습니다: %w", err)
		}
		return fmt.Errorf("타임존 업데이트 중 오류가 발생했습니다: %w", err)
	}

	logger.Sugar().Infof("사용자 타임존이 업데이트되었습니다. 사용자ID: %s, 타임존: %s", userID.String(), timezone)
	return nil
}

func (r *UserRepository) GetUserWithFCM(userID uuid.UUID) (*domain.User, error) {
	u, err := r.client.User.Query().
		Where(user.ID(userID)).
		Only(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("해당하는 ID로 사용자를 찾을 수 없습니다: %w", err)
		}
		return nil, fmt.Errorf("사용자 정보 조회 중 오류가 발생했습니다: %w", err)
	}

	return &domain.User{
		ID:            u.ID,
		NickName:      u.NickName,
		Email:         u.Email,
		Password:      u.Password,
		IsPublished:   u.IsPublished,
		IsTermsAgreed: u.IsTermsAgreed,
		FCMToken:      u.FcmToken,
		Timezone:      u.Timezone,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}, nil
}

func (r *UserRepository) GetAllUsersWithFCM() ([]*domain.User, error) {
	users, err := r.client.User.Query().
		Where(user.FcmTokenNEQ("")).
		All(context.Background())

	if err != nil {
		return nil, fmt.Errorf("FCM 토큰이 있는 사용자 목록 조회 중 오류가 발생했습니다: %w", err)
	}

	result := make([]*domain.User, len(users))
	for i, u := range users {
		result[i] = &domain.User{
			ID:            u.ID,
			NickName:      u.NickName,
			Email:         u.Email,
			IsPublished:   u.IsPublished,
			IsTermsAgreed: u.IsTermsAgreed,
			FCMToken:      u.FcmToken,
			Timezone:      u.Timezone,
			CreatedAt:     u.CreatedAt,
			UpdatedAt:     u.UpdatedAt,
		}
	}

	logger.Sugar().Infof("FCM 토큰이 있는 사용자 %d명 조회 완료", len(result))
	return result, nil
}
