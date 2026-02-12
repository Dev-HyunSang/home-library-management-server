package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/dev-hyunsang/my-own-library-backend/internal/domain"
	"github.com/dev-hyunsang/my-own-library-backend/lib/ent"
	"github.com/dev-hyunsang/my-own-library-backend/lib/ent/emailverification"
	"github.com/dev-hyunsang/my-own-library-backend/logger"
	"github.com/google/uuid"
)

type EmailVerificationRepository struct {
	client *ent.Client
}

func NewEmailVerificationRepository(client *ent.Client) *EmailVerificationRepository {
	return &EmailVerificationRepository{
		client: client,
	}
}

func (r *EmailVerificationRepository) Save(verification *domain.EmailVerification) (*domain.EmailVerification, error) {
	v, err := r.client.EmailVerification.Create().
		SetID(verification.ID).
		SetEmail(verification.Email).
		SetCode(verification.Code).
		SetExpiresAt(verification.ExpiresAt).
		SetIsVerified(false).
		SetCreatedAt(time.Now()).
		Save(context.Background())

	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fmt.Errorf("이메일 인증 정보 저장 중 제약조건 오류가 발생했습니다: %w", err)
		}
		return nil, fmt.Errorf("이메일 인증 정보를 저장하는 도중 오류가 발생했습니다: %w", err)
	}

	logger.Sugar().Infof("이메일 인증 정보를 생성했습니다. 이메일: %s", v.Email)
	return &domain.EmailVerification{
		ID:         v.ID,
		Email:      v.Email,
		Code:       v.Code,
		ExpiresAt:  v.ExpiresAt,
		IsVerified: v.IsVerified,
		CreatedAt:  v.CreatedAt,
	}, nil
}

func (r *EmailVerificationRepository) GetByEmailAndCode(email, code string) (*domain.EmailVerification, error) {
	v, err := r.client.EmailVerification.Query().
		Where(
			emailverification.Email(email),
			emailverification.Code(code),
			emailverification.ExpiresAtGT(time.Now()),
			emailverification.IsVerified(false),
		).
		Only(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("유효한 인증 정보를 찾을 수 없습니다: %w", err)
		}
		return nil, fmt.Errorf("인증 정보 조회 중 오류가 발생했습니다: %w", err)
	}

	return &domain.EmailVerification{
		ID:         v.ID,
		Email:      v.Email,
		Code:       v.Code,
		ExpiresAt:  v.ExpiresAt,
		IsVerified: v.IsVerified,
		CreatedAt:  v.CreatedAt,
	}, nil
}

func (r *EmailVerificationRepository) GetLatestByEmail(email string) (*domain.EmailVerification, error) {
	v, err := r.client.EmailVerification.Query().
		Where(emailverification.Email(email)).
		Order(ent.Desc(emailverification.FieldCreatedAt)).
		First(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("해당 이메일의 인증 정보를 찾을 수 없습니다: %w", err)
		}
		return nil, fmt.Errorf("인증 정보 조회 중 오류가 발생했습니다: %w", err)
	}

	return &domain.EmailVerification{
		ID:         v.ID,
		Email:      v.Email,
		Code:       v.Code,
		ExpiresAt:  v.ExpiresAt,
		IsVerified: v.IsVerified,
		CreatedAt:  v.CreatedAt,
	}, nil
}

func (r *EmailVerificationRepository) MarkAsVerified(id uuid.UUID) error {
	err := r.client.EmailVerification.UpdateOneID(id).
		SetIsVerified(true).
		Exec(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("해당하는 인증 정보를 찾을 수 없습니다: %w", err)
		}
		return fmt.Errorf("인증 완료 처리 중 오류가 발생했습니다: %w", err)
	}

	logger.Sugar().Infof("이메일 인증이 완료되었습니다. ID: %s", id.String())
	return nil
}

func (r *EmailVerificationRepository) DeleteByEmail(email string) error {
	_, err := r.client.EmailVerification.Delete().
		Where(emailverification.Email(email)).
		Exec(context.Background())

	if err != nil {
		return fmt.Errorf("이메일 인증 정보 삭제 중 오류가 발생했습니다: %w", err)
	}

	logger.Sugar().Infof("이메일 인증 정보를 삭제했습니다. 이메일: %s", email)
	return nil
}

func (r *EmailVerificationRepository) DeleteExpired() error {
	deleted, err := r.client.EmailVerification.Delete().
		Where(emailverification.ExpiresAtLT(time.Now())).
		Exec(context.Background())

	if err != nil {
		return fmt.Errorf("만료된 인증 정보 삭제 중 오류가 발생했습니다: %w", err)
	}

	if deleted > 0 {
		logger.Sugar().Infof("만료된 이메일 인증 정보 %d건을 삭제했습니다", deleted)
	}
	return nil
}
