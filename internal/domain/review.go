package domain

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID        uuid.UUID `json:"id"`
	OwnerID   uuid.UUID `json:"owner_id"`
	BookISBN  string    `json:"book_isbn"`
	Content   string    `json:"content"`
	Rating    int       `json:"rating"`
	IsPublic  bool      `json:"is_public"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ReviewResponse struct {
	ID            uuid.UUID `json:"id"`
	OwnerID       uuid.UUID `json:"owner_id"`
	OwnerNickname string    `json:"owner_nickname,omitempty"`
	BookISBN      string    `json:"book_isbn"`
	Content       string    `json:"content"`
	Rating        int       `json:"rating"`
	IsPublic      bool      `json:"is_public"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ReviewWithBook struct {
	ID        uuid.UUID `json:"id"`
	OwnerID   uuid.UUID `json:"owner_id"`
	BookISBN  string    `json:"book_isbn"`
	Content   string    `json:"content"`
	Rating    int       `json:"rating"`
	IsPublic  bool      `json:"is_public"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Book      *BookInfo `json:"book,omitempty"`
}

type BookInfo struct {
	Title        string `json:"title"`
	Author       string `json:"author"`
	ThumbnailURL string `json:"thumbnail_url"`
}

type CreateReviewRequest struct {
	Content  string `json:"content"`
	Rating   int    `json:"rating"`
	IsPublic bool   `json:"is_public"`
}

type UpdateReviewRequest struct {
	Content  *string `json:"content,omitempty"`
	Rating   *int    `json:"rating,omitempty"`
	IsPublic *bool   `json:"is_public,omitempty"`
}

type ReviewRepository interface {
	Create(review *Review) (*Review, error)
	GetByID(id uuid.UUID) (*Review, error)
	GetByISBN(isbn string) ([]*Review, error)
	GetPublicByISBN(isbn string) ([]*ReviewResponse, error)
	GetByUserID(userID uuid.UUID) ([]*Review, error)
	ExistsByUserAndISBN(userID uuid.UUID, isbn string) (bool, error)
	Update(review *Review) (*Review, error)
	Delete(userID, reviewID uuid.UUID) error
}

type ReviewUseCase interface {
	CreateReview(userID uuid.UUID, isbn string, req *CreateReviewRequest) (*Review, error)
	GetReviewByID(id uuid.UUID) (*Review, error)
	GetReviewsByISBN(isbn string) ([]*ReviewResponse, error)
	GetUserReviews(userID uuid.UUID) ([]*Review, error)
	UpdateReview(userID, reviewID uuid.UUID, req *UpdateReviewRequest) (*Review, error)
	DeleteReview(userID, reviewID uuid.UUID) error
}
