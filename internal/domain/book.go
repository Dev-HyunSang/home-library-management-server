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
	BookISBN     int       `json:"book_isbn"`
	RegisteredAt time.Time `json:"registered_at"`
	ComplatedAt  time.Time `json:"complated_at"`
}

type BookRepository interface {
	SaveByBookID(id uuid.UUID, book *Book) (*Book, error)
	GetByBookID(userID, id uuid.UUID) (*Book, error)
	GetBooksByUserID(id uuid.UUID) ([]*Book, error)
	Edit(id uuid.UUID, book *Book) error
	DeleteByID(userID, id uuid.UUID) error
	GetBooksByUserName(name string) ([]*Book, error)
}

type BookUseCase interface {
	SaveByBookID(userID uuid.UUID, book *Book) (*Book, error)
	GetByBookID(userID, id uuid.UUID) (*Book, error)
	GetBooksByUserID(userID uuid.UUID) ([]*Book, error)
	Edit(id uuid.UUID, book *Book) error
	DeleteByID(userID, id uuid.UUID) error
	GetBooksByUserName(name string) ([]*Book, error)
}
