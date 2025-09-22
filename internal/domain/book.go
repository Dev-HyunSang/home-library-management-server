package domain

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID           uuid.UUID `json:"id"`
	OwnerID      uuid.UUID `json:"owner_id"`
	Title        string    `json:"title"`
	Author       string    `json:"author"`
	BookISBN     string    `json:"book_isbn"`
	RegisteredAt time.Time `json:"registered_at"`
	ComplatedAt  time.Time `json:"complated_at"`
}

type ReviewBook struct {
	ID         uuid.UUID `json:"id"`
	BookID     uuid.UUID `json:"book_id"`
	OwnerID    uuid.UUID `json:"owner_id"`
	BookTitle  string    `json:"book_title"`
	BookAuthor string    `json:"book_author"`
	Content    string    `json:"content"`
	Rating     int       `json:"rating"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type BookRepository interface {
	SaveByBookID(id uuid.UUID, book *Book) (*Book, error)
	GetBookByID(userID, id uuid.UUID) (*Book, error)
	GetBooksByUserID(id uuid.UUID) ([]*Book, error)
	Edit(id uuid.UUID, book *Book) error
	DeleteByID(userID, id uuid.UUID) error
	GetBooksByUserName(name string) ([]*Book, error)
	// Book Review
	CreateReview(review *ReviewBook) error
	GetReviewsByUserID(userID uuid.UUID) ([]*ReviewBook, error)
	GetReviewByID(id uuid.UUID) (*ReviewBook, error)
	UpdateReviewByID(review *ReviewBook) (ReviewBook, error)
}

type BookUseCase interface {
	SaveByBookID(userID uuid.UUID, book *Book) (*Book, error)
	GetBookByID(userID, id uuid.UUID) (*Book, error)
	GetBooksByUserID(userID uuid.UUID) ([]*Book, error)
	Edit(id uuid.UUID, book *Book) error
	DeleteByID(userID, id uuid.UUID) error
	GetBooksByUserName(name string) ([]*Book, error)
	// Book Review
	CreateReview(review *ReviewBook) error
	GetReviewsByUserID(userID uuid.UUID) ([]*ReviewBook, error)
	GetReviewByID(id uuid.UUID) (*ReviewBook, error)
	UpdateReviewByID(review *ReviewBook) (ReviewBook, error)
}
