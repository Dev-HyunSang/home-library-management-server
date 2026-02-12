package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/dev-hyunsang/my-own-library-backend/internal/domain"
	"github.com/dev-hyunsang/my-own-library-backend/lib/ent"
	"github.com/dev-hyunsang/my-own-library-backend/lib/ent/review"
	"github.com/dev-hyunsang/my-own-library-backend/lib/ent/user"
	"github.com/dev-hyunsang/my-own-library-backend/logger"
	"github.com/google/uuid"
)

type ReviewRepository struct {
	client *ent.Client
}

func NewReviewRepository(client *ent.Client) *ReviewRepository {
	return &ReviewRepository{client: client}
}

func (r *ReviewRepository) Create(rev *domain.Review) (*domain.Review, error) {
	created, err := r.client.Review.Create().
		SetID(rev.ID).
		SetBookIsbn(rev.BookISBN).
		SetContent(rev.Content).
		SetRating(rev.Rating).
		SetIsPublic(rev.IsPublic).
		SetOwnerID(rev.OwnerID).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		Save(context.Background())

	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fmt.Errorf("리뷰 저장 중 제약조건 오류가 발생했습니다: %w", err)
		}
		return nil, fmt.Errorf("리뷰를 저장하는 도중 오류가 발생했습니다: %w", err)
	}

	logger.Sugar().Infof("리뷰가 생성되었습니다. ID: %s, ISBN: %s", created.ID.String(), created.BookIsbn)

	return &domain.Review{
		ID:        created.ID,
		OwnerID:   rev.OwnerID,
		BookISBN:  created.BookIsbn,
		Content:   created.Content,
		Rating:    created.Rating,
		IsPublic:  created.IsPublic,
		CreatedAt: created.CreatedAt,
		UpdatedAt: created.UpdatedAt,
	}, nil
}

func (r *ReviewRepository) GetByID(id uuid.UUID) (*domain.Review, error) {
	rev, err := r.client.Review.Query().
		Where(review.ID(id)).
		WithOwner().
		Only(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("해당 리뷰를 찾을 수 없습니다: %w", err)
		}
		return nil, fmt.Errorf("리뷰 조회 중 오류가 발생했습니다: %w", err)
	}

	return &domain.Review{
		ID:        rev.ID,
		OwnerID:   rev.Edges.Owner.ID,
		BookISBN:  rev.BookIsbn,
		Content:   rev.Content,
		Rating:    rev.Rating,
		IsPublic:  rev.IsPublic,
		CreatedAt: rev.CreatedAt,
		UpdatedAt: rev.UpdatedAt,
	}, nil
}

func (r *ReviewRepository) GetByISBN(isbn string) ([]*domain.Review, error) {
	reviews, err := r.client.Review.Query().
		Where(review.BookIsbn(isbn)).
		WithOwner().
		Order(ent.Desc(review.FieldCreatedAt)).
		All(context.Background())

	if err != nil {
		return nil, fmt.Errorf("ISBN으로 리뷰 조회 중 오류가 발생했습니다: %w", err)
	}

	result := make([]*domain.Review, len(reviews))
	for i, rev := range reviews {
		result[i] = &domain.Review{
			ID:        rev.ID,
			OwnerID:   rev.Edges.Owner.ID,
			BookISBN:  rev.BookIsbn,
			Content:   rev.Content,
			Rating:    rev.Rating,
			IsPublic:  rev.IsPublic,
			CreatedAt: rev.CreatedAt,
			UpdatedAt: rev.UpdatedAt,
		}
	}

	return result, nil
}

func (r *ReviewRepository) GetPublicByISBN(isbn string) ([]*domain.ReviewResponse, error) {
	reviews, err := r.client.Review.Query().
		Where(
			review.BookIsbn(isbn),
			review.IsPublic(true),
		).
		WithOwner().
		Order(ent.Desc(review.FieldCreatedAt)).
		All(context.Background())

	if err != nil {
		return nil, fmt.Errorf("공개 리뷰 조회 중 오류가 발생했습니다: %w", err)
	}

	result := make([]*domain.ReviewResponse, len(reviews))
	for i, rev := range reviews {
		nickname := ""
		if rev.Edges.Owner != nil {
			nickname = rev.Edges.Owner.NickName
		}

		result[i] = &domain.ReviewResponse{
			ID:            rev.ID,
			OwnerID:       rev.Edges.Owner.ID,
			OwnerNickname: nickname,
			BookISBN:      rev.BookIsbn,
			Content:       rev.Content,
			Rating:        rev.Rating,
			IsPublic:      rev.IsPublic,
			CreatedAt:     rev.CreatedAt,
			UpdatedAt:     rev.UpdatedAt,
		}
	}

	return result, nil
}

func (r *ReviewRepository) ExistsByUserAndISBN(userID uuid.UUID, isbn string) (bool, error) {
	exists, err := r.client.Review.Query().
		Where(
			review.HasOwnerWith(user.ID(userID)),
			review.BookIsbn(isbn),
		).
		Exist(context.Background())

	if err != nil {
		return false, fmt.Errorf("리뷰 존재 여부 확인 중 오류가 발생했습니다: %w", err)
	}

	return exists, nil
}

func (r *ReviewRepository) GetByUserID(userID uuid.UUID) ([]*domain.Review, error) {
	reviews, err := r.client.Review.Query().
		Where(review.HasOwnerWith(user.ID(userID))).
		Order(ent.Desc(review.FieldCreatedAt)).
		All(context.Background())

	if err != nil {
		return nil, fmt.Errorf("사용자 리뷰 조회 중 오류가 발생했습니다: %w", err)
	}

	result := make([]*domain.Review, len(reviews))
	for i, rev := range reviews {
		result[i] = &domain.Review{
			ID:        rev.ID,
			OwnerID:   userID,
			BookISBN:  rev.BookIsbn,
			Content:   rev.Content,
			Rating:    rev.Rating,
			IsPublic:  rev.IsPublic,
			CreatedAt: rev.CreatedAt,
			UpdatedAt: rev.UpdatedAt,
		}
	}

	return result, nil
}

func (r *ReviewRepository) Update(rev *domain.Review) (*domain.Review, error) {
	updated, err := r.client.Review.UpdateOneID(rev.ID).
		SetContent(rev.Content).
		SetRating(rev.Rating).
		SetIsPublic(rev.IsPublic).
		SetUpdatedAt(time.Now()).
		Save(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("해당 리뷰를 찾을 수 없습니다: %w", err)
		}
		return nil, fmt.Errorf("리뷰 수정 중 오류가 발생했습니다: %w", err)
	}

	logger.Sugar().Infof("리뷰가 수정되었습니다. ID: %s", updated.ID.String())

	return &domain.Review{
		ID:        updated.ID,
		OwnerID:   rev.OwnerID,
		BookISBN:  updated.BookIsbn,
		Content:   updated.Content,
		Rating:    updated.Rating,
		IsPublic:  updated.IsPublic,
		CreatedAt: updated.CreatedAt,
		UpdatedAt: updated.UpdatedAt,
	}, nil
}

func (r *ReviewRepository) Delete(userID, reviewID uuid.UUID) error {
	rev, err := r.client.Review.Query().
		Where(review.ID(reviewID)).
		WithOwner().
		Only(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("해당 리뷰를 찾을 수 없습니다: %w", err)
		}
		return fmt.Errorf("리뷰 조회 중 오류가 발생했습니다: %w", err)
	}

	if rev.Edges.Owner.ID != userID {
		return fmt.Errorf("리뷰를 삭제할 권한이 없습니다")
	}

	err = r.client.Review.DeleteOneID(reviewID).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("리뷰 삭제 중 오류가 발생했습니다: %w", err)
	}

	logger.Sugar().Infof("리뷰가 삭제되었습니다. ID: %s", reviewID.String())
	return nil
}
