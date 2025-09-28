package usecase

import (
	"github.com/dev-hyunsang/home-library/internal/domain"
	"github.com/google/uuid"
)

type BookUseCase struct {
	bookRepo domain.BookRepository
}

func NewBookUseCase(repo domain.BookRepository) *BookUseCase {
	return &BookUseCase{bookRepo: repo}
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

	return bc.bookRepo.CreateReview(review)
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
