package usecase

import (
	"context"

	"github.com/dev-hyunsang/home-library/internal/domain"
	"github.com/google/uuid"
)

type BookUseCase struct {
	bookRepo domain.BookRepository
	notifier domain.NotificationProducer
}

func NewBookUseCase(repo domain.BookRepository, notifier domain.NotificationProducer) *BookUseCase {
	return &BookUseCase{
		bookRepo: repo,
		notifier: notifier,
	}
}

func (bc *BookUseCase) SaveByBookID(userID uuid.UUID, book *domain.Book) (*domain.Book, error) {
	if book.Title == "" || book.Author == "" {
		return nil, domain.ErrInvalidInput
	}

	return bc.bookRepo.SaveByBookID(userID, book)
}

func (bc *BookUseCase) GetBookByID(userID, id uuid.UUID) (*domain.Book, error) {
	if id == uuid.Nil {
		return nil, domain.ErrInvalidInput
	}

	return bc.bookRepo.GetBookByID(userID, id)
}

func (bc *BookUseCase) GetBooksByUserID(userID uuid.UUID) ([]*domain.Book, error) {
	if userID == uuid.Nil {
		return nil, domain.ErrInvalidInput
	}

	return bc.bookRepo.GetBooksByUserID(userID)
}

func (bc *BookUseCase) GetBooksByUserName(name string) ([]*domain.Book, error) {
	if len(name) == 0 {
		return nil, domain.ErrInvalidInput
	}

	return bc.bookRepo.GetBooksByUserName(name)
}

func (bc *BookUseCase) Edit(id uuid.UUID, book *domain.Book) error {
	if id == uuid.Nil || book == nil {
		return domain.ErrInvalidInput
	}

	return bc.bookRepo.Edit(id, book)
}

func (bc *BookUseCase) DeleteByID(userID, id uuid.UUID) error {
	if userID == uuid.Nil || id == uuid.Nil {
		return domain.ErrInvalidInput
	}

	return bc.bookRepo.DeleteByID(userID, id)
}

// Book Review

func (bc *BookUseCase) CreateReview(review *domain.ReviewBook) error {
	if review == nil || review.BookID == uuid.Nil || review.OwnerID == uuid.Nil || review.Rating < 1 || review.Rating > 5 || review.Content == "" {
		return domain.ErrInvalidInput
	}

	err := bc.bookRepo.CreateReview(review)
	if err != nil {
		return err
	}

	// 알림 발송 (비동기 처리 가능하나 여기선 동기 호출, Producer 자체는 빠름)
	// TODO: Context 전달 혹은 Background 사용
	if bc.notifier != nil {
		// 자신의 책에 리뷰를 남기면 알림 X? 하지만 여기선 테스트 목적
		// 실제로는 책 주인의 ID를 찾아야 함. ReviewBook 구조체에 BookOwnerID가 있다면 좋겠지만...
		// 지금은 리뷰 작성자에게 "리뷰가 작성되었습니다" 알림을 보내거나, 책 주인에게 보냄.
		// ReviewBook에 OwnerID는 "리뷰 작성자"임.
		// 책 정보를 조회해서 책 주인을 찾아야 함.
		// 간단하게 리뷰 작성자에게 알림 테스트.
		_ = bc.notifier.ProduceNotification(context.Background(), review.OwnerID.String(), "리뷰 작성 완료", "리뷰가 성공적으로 작성되었습니다.", "review")
	}

	return nil
}

func (bc *BookUseCase) GetReviewsByUserID(userID uuid.UUID) ([]*domain.ReviewBook, error) {
	if userID == uuid.Nil {
		return nil, domain.ErrInvalidInput
	}

	return bc.bookRepo.GetReviewsByUserID(userID)
}

func (bc *BookUseCase) GetReviewByID(id uuid.UUID) (*domain.ReviewBook, error) {
	if id == uuid.Nil {
		return nil, domain.ErrInvalidInput
	}

	return bc.bookRepo.GetReviewByID(id)
}

func (bc *BookUseCase) UpdateReviewByID(review *domain.ReviewBook) (domain.ReviewBook, error) {
	if review == nil || review.ID == uuid.Nil {
		return domain.ReviewBook{}, domain.ErrInvalidInput
	}

	return bc.bookRepo.UpdateReviewByID(review)
}

func (bc *BookUseCase) AddBookmarkByBookID(userID, bookID uuid.UUID) (*domain.Bookmark, error) {
	if userID == uuid.Nil || bookID == uuid.Nil {
		return nil, domain.ErrInvalidInput
	}

	return bc.bookRepo.AddBookmarkByBookID(userID, bookID)
}

func (bc *BookUseCase) GetBookmarksByUserID(userID uuid.UUID) ([]*domain.Bookmark, error) {
	if userID == uuid.Nil {
		return nil, domain.ErrInvalidInput
	}

	return bc.bookRepo.GetBookmarksByUserID(userID)
}

func (bc *BookUseCase) DeleteBookmarkByID(id uuid.UUID) error {
	if id == uuid.Nil {
		return domain.ErrInvalidInput
	}

	return bc.bookRepo.DeleteBookmarkByID(id)
}
