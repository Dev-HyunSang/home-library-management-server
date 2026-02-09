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
	ThumbnailURL string    `json:"thumbnail_url"`
	Status       int       `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ReviewBook struct {
	ID         uuid.UUID `json:"id"`
	BookID     uuid.UUID `json:"book_id"`
	OwnerID    uuid.UUID `json:"owner_id"`
	BookISBN   string    `json:"book_isbn"`
	BookTitle  string    `json:"book_title"`
	BookAuthor string    `json:"book_author"`
	Content    string    `json:"content"`
	Rating     int       `json:"rating"`
	IsPublic   bool      `json:"is_public"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Bookmark struct {
	ID        uuid.UUID `json:"id"`
	OwnerID   uuid.UUID `json:"owner_id"`
	BookID    uuid.UUID `json:"book_id"`
	CreatedAt time.Time `json:"created_at"`
}

type BookRepository interface {
	SaveByBookID(id uuid.UUID, book *Book) (*Book, error)
	GetBookByID(userID, id uuid.UUID) (*Book, error)
	GetBookByISBN(userID uuid.UUID, isbn string) (*Book, error)
	GetAnyBookByISBN(isbn string) (*Book, error)
	GetBooksByUserID(id uuid.UUID) ([]*Book, error)
	Edit(id uuid.UUID, book *Book) error
	DeleteByID(userID, id uuid.UUID) error
	GetBooksByUserName(name string) ([]*Book, error)
	// Book Review
	CreateReview(review *ReviewBook) error
	GetReviewsByUserID(userID uuid.UUID) ([]*ReviewBook, error)
	GetReviewByID(id uuid.UUID) (*ReviewBook, error)
	UpdateReviewByID(review *ReviewBook) (ReviewBook, error)
	DeleteReviewByID(userID, reviewID uuid.UUID) error
	GetPublicReviewsByBookID(bookID uuid.UUID) ([]*ReviewBook, error)
	GetPublicReviewsByISBN(isbn string) ([]*ReviewBook, error)
	// Book Bookmark
	AddBookmarkByBookID(userID, bookID uuid.UUID) (*Bookmark, error)
	GetBookmarksByUserID(userID uuid.UUID) ([]*Bookmark, error)
	DeleteBookmarkByID(id uuid.UUID) error
}

type BookUseCase interface {
	SaveByBookID(userID uuid.UUID, book *Book) (*Book, error)
	GetBookByID(userID, id uuid.UUID) (*Book, error)
	GetBookByISBN(userID uuid.UUID, isbn string) (*Book, error)
	GetAnyBookByISBN(isbn string) (*Book, error)
	GetBooksByUserID(userID uuid.UUID) ([]*Book, error)
	Edit(id uuid.UUID, book *Book) error
	DeleteByID(userID, id uuid.UUID) error
	GetBooksByUserName(name string) ([]*Book, error)
	// Book Review
	CreateReview(review *ReviewBook) error
	GetReviewsByUserID(userID uuid.UUID) ([]*ReviewBook, error)
	GetReviewByID(id uuid.UUID) (*ReviewBook, error)
	UpdateReviewByID(review *ReviewBook) (ReviewBook, error)
	DeleteReviewByID(userID, reviewID uuid.UUID) error
	GetPublicReviewsByBookID(bookID uuid.UUID) ([]*ReviewBook, error)
	GetPublicReviewsByISBN(isbn string) ([]*ReviewBook, error)
	// Book Bookmark
	AddBookmarkByBookID(userID, bookID uuid.UUID) (*Bookmark, error)
	GetBookmarksByUserID(userID uuid.UUID) ([]*Bookmark, error)
	DeleteBookmarkByID(id uuid.UUID) error
}
