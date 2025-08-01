package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/dev-hyunsang/home-library/internal/domain"
	"github.com/dev-hyunsang/home-library/lib/ent"
	"github.com/dev-hyunsang/home-library/lib/ent/book"
	"github.com/dev-hyunsang/home-library/lib/ent/user"
	"github.com/google/uuid"
)

type BookRepository struct {
	client *ent.Client
}

func NewBookRepository(client *ent.Client) *BookRepository {
	return &BookRepository{
		client: client,
	}
}

func (rc *BookRepository) SaveByBookID(userID uuid.UUID, book *domain.Book) (*domain.Book, error) {
	client := rc.client

	BookID, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate book id(uuid): %w", err)
	}

	b, err := client.Book.Create().
		SetOwnerID(userID).
		SetID(BookID).
		SetBookTitle(book.Title).
		SetAuthor(book.Author).
		SetBookIsbn(book.BookISBN).
		SetRegisteredAt(time.Now()).
		SetComplatedAt(time.Now()).
		Save(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to create book: %w", err)
	}

	return &domain.Book{
		ID:           b.ID,
		OwnerID:      b.QueryOwner().OnlyIDX(context.Background()),
		Title:        b.BookTitle,
		Author:       b.Author,
		BookISBN:     b.BookIsbn,
		RegisteredAt: b.RegisteredAt,
		ComplatedAt:  b.ComplatedAt,
	}, nil
}

// 책 정보를 가져옵니다. UserID와 (Book)ID가 일치하여만 책 정보를 가져올 수 있습니다.
func (rc *BookRepository) GetByBookID(userID, id uuid.UUID) (*domain.Book, error) {
	client := rc.client

	b, err := client.Book.
		Query().
		Where(
			book.ID(id),
			book.HasOwnerWith(user.ID(userID))). // UserID와 일치하는 조건을 찾습니다.
		Only(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get book by id: %w", err)
	}

	return &domain.Book{
		ID:           b.ID,
		OwnerID:      b.Edges.Owner.ID,
		Title:        b.BookTitle,
		Author:       b.Author,
		BookISBN:     b.BookIsbn,
		RegisteredAt: b.RegisteredAt,
		ComplatedAt:  b.ComplatedAt,
	}, nil
}

// 유저가 소유한 책의 목록을 가져옵니다. / UserID와 일치하는 경우에만 책의 목록을 가져올 수 있습니다.
func (rc *BookRepository) GetBooksByUserID(userID uuid.UUID) ([]*domain.Book, error) {
	client := rc.client

	books, err := client.Book.
		Query().
		Where(book.HasOwnerWith(user.ID(userID))).
		All(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get books by user id: %w", err)
	}

	var result []*domain.Book
	for _, b := range books {
		result = append(result, &domain.Book{
			ID:           b.ID,
			OwnerID:      b.Edges.Owner.ID,
			Title:        b.BookTitle,
			Author:       b.Author,
			BookISBN:     b.BookIsbn,
			RegisteredAt: b.RegisteredAt,
			ComplatedAt:  b.ComplatedAt,
		})
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no books found for user id: %s", userID)
	}

	return result, nil
}

// 사용자 이름을 통해 책의 목록을 가져옵니다.
// TODO: 추후에 공개계정과 비공개 계정을 나눌 예정
func (bc *BookRepository) GetBooksByUserName(name string) ([]*domain.Book, error) {
	client := bc.client

	books, err := client.Book.
		Query().
		Where(book.HasOwnerWith(user.NickName(name))).
		All(context.Background())
	if err != nil {
		return nil, domain.ErrNotFound
	}

	var result []*domain.Book
	for _, b := range books {
		result = append(result, &domain.Book{
			ID:           b.ID,
			OwnerID:      b.QueryOwner().OnlyIDX(context.Background()),
			Title:        b.BookTitle,
			Author:       b.Author,
			BookISBN:     b.BookIsbn,
			RegisteredAt: b.RegisteredAt,
			ComplatedAt:  b.ComplatedAt,
		})
	}

	return result, nil
}

func (bc *BookRepository) Edit(id uuid.UUID, book *domain.Book) error {
	client := bc.client

	_, err := client.Book.UpdateOneID(id).
		SetBookTitle(book.Title).
		SetAuthor(book.Author).
		SetBookIsbn(book.BookISBN).
		SetRegisteredAt(time.Now()).
		SetComplatedAt(book.ComplatedAt).
		Save(context.Background())
	if err != nil {
		return fmt.Errorf("failed to edit book by id: %w", err)
	}

	return nil
}

func (bc *BookRepository) DeleteByID(userID, id uuid.UUID) error {
	client := bc.client

	err := client.Book.DeleteOneID(id).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete book by id: %w", err)
	}

	return nil
}
