package usecase

import (
	"fmt"

	"github.com/dev-hyunsang/home-library/internal/domain"
	"github.com/google/uuid"
)

type ReviewUseCase struct {
	reviewRepo domain.ReviewRepository
}

func NewReviewUseCase(repo domain.ReviewRepository) *ReviewUseCase {
	return &ReviewUseCase{
		reviewRepo: repo,
	}
}

func (uc *ReviewUseCase) CreateReview(userID uuid.UUID, isbn string, req *domain.CreateReviewRequest) (*domain.Review, error) {
	if isbn == "" {
		return nil, fmt.Errorf("ISBN은 필수입니다")
	}

	if req.Content == "" {
		return nil, fmt.Errorf("리뷰 내용은 필수입니다")
	}

	if req.Rating < 1 || req.Rating > 5 {
		return nil, fmt.Errorf("별점은 1~5 사이여야 합니다")
	}

	review := &domain.Review{
		ID:       uuid.New(),
		OwnerID:  userID,
		BookISBN: isbn,
		Content:  req.Content,
		Rating:   req.Rating,
		IsPublic: req.IsPublic,
	}

	return uc.reviewRepo.Create(review)
}

func (uc *ReviewUseCase) GetReviewByID(id uuid.UUID) (*domain.Review, error) {
	return uc.reviewRepo.GetByID(id)
}

func (uc *ReviewUseCase) GetReviewsByISBN(isbn string) ([]*domain.ReviewResponse, error) {
	if isbn == "" {
		return nil, fmt.Errorf("ISBN은 필수입니다")
	}

	return uc.reviewRepo.GetPublicByISBN(isbn)
}

func (uc *ReviewUseCase) GetUserReviews(userID uuid.UUID) ([]*domain.Review, error) {
	return uc.reviewRepo.GetByUserID(userID)
}

func (uc *ReviewUseCase) UpdateReview(userID, reviewID uuid.UUID, req *domain.UpdateReviewRequest) (*domain.Review, error) {
	existing, err := uc.reviewRepo.GetByID(reviewID)
	if err != nil {
		return nil, err
	}

	if existing.OwnerID != userID {
		return nil, fmt.Errorf("리뷰를 수정할 권한이 없습니다")
	}

	if req.Content != nil {
		if *req.Content == "" {
			return nil, fmt.Errorf("리뷰 내용은 필수입니다")
		}
		existing.Content = *req.Content
	}

	if req.Rating != nil {
		if *req.Rating < 1 || *req.Rating > 5 {
			return nil, fmt.Errorf("별점은 1~5 사이여야 합니다")
		}
		existing.Rating = *req.Rating
	}

	if req.IsPublic != nil {
		existing.IsPublic = *req.IsPublic
	}

	return uc.reviewRepo.Update(existing)
}

func (uc *ReviewUseCase) DeleteReview(userID, reviewID uuid.UUID) error {
	return uc.reviewRepo.Delete(userID, reviewID)
}
