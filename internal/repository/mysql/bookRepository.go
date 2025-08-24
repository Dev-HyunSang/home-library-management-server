package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/dev-hyunsang/home-library/internal/domain"
	"github.com/dev-hyunsang/home-library/lib/ent"
	"github.com/dev-hyunsang/home-library/lib/ent/book"
	"github.com/dev-hyunsang/home-library/lib/ent/user"
	"github.com/dev-hyunsang/home-library/logger"
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
		logger.Init().Sugar().Errorf("새로운 UUID를 생성하던 도중 오류가 발생했습니다. %w", err)
		return nil, fmt.Errorf("새로운 UUID를 생성하던 도중 오류가 발생했습니다. %w", err)
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
	if err == nil {
		logger.UserInfoLog(userID.String(), "해당 유저의 새로운 책을 저장했습니다.")
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

	switch {
	case ent.IsConstraintError(err):
		logger.Init().Sugar().Errorf("책을 저장하는 도중 제약 조건 오류가 발생했습니다: %w", err)
		return nil, fmt.Errorf("책을 저장하는 도중 제약 조건 오류가 발생했습니다: %w", err)
	default:
		return nil, fmt.Errorf("새로운 책을 저장하던 도중 오류가 발생했습니다: %w", err)
	}
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

	logger.UserInfoLog(userID.String(), "해당 유저의 책 정보를 가지고 왔습니다.")

	if err == nil {
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

	switch {
	case ent.IsNotFound(err):
		logger.Init().Sugar().Warnf("등록된 책을 찾을 수 없습니다: %w", err)
		return nil, fmt.Errorf("등록된 책을 찾을 수 없습니다: %w", err)
	case ent.IsConstraintError(err):
		logger.Init().Sugar().Errorf("책을 가져오는 도중 제약 조건 오류가 발생했습니다: %w", err)
		return nil, fmt.Errorf("책을 가져오는 도중 제약 조건 오류가 발생했습니다: %w", err)
	default:
		return nil, fmt.Errorf("책을 가져오는 도중 오류가 발생했습니다: %w", err)
	}

}

// 유저가 소유한 책의 목록을 가져옵니다. / UserID와 일치하는 경우에만 책의 목록을 가져올 수 있습니다.
func (rc *BookRepository) GetBooksByUserID(userID uuid.UUID) ([]*domain.Book, error) {
	var result []*domain.Book

	client := rc.client

	books, err := client.Book.
		Query().
		Where(book.HasOwnerWith(user.ID(userID))).
		All(context.Background())
	if err == nil {
		for _, b := range books {
			result = append(result, &domain.Book{
				ID:           b.ID,
				OwnerID:      userID,
				Title:        b.BookTitle,
				Author:       b.Author,
				BookISBN:     string(b.BookIsbn),
				RegisteredAt: b.RegisteredAt,
				ComplatedAt:  b.ComplatedAt,
			})
		}

		return result, nil
	}

	switch {
	case ent.IsNotFound(err):
		return nil, fmt.Errorf("등록된 책을 찾을 수 없습니다: %w", err)
	case ent.IsNotLoaded(err):
		return nil, fmt.Errorf("책의 목록을 불러오는 도중 데이터가 로드되지 않았습니다: %w", err)
	case ent.IsConstraintError(err):
		return nil, fmt.Errorf("책의 목록을 가져오는 도중 제약 조건 오류가 발생했습니다: %w", err)
	default:
		return nil, fmt.Errorf("책의 목록을 가져오는 도중 알 수 없는 오류가 발생했습니다: %w", err)
	}
}

// 사용자 이름을 통해 책의 목록을 가져옵니다.
func (bc *BookRepository) GetBooksByUserName(name string) ([]*domain.Book, error) {
	client := bc.client

	books, err := client.Book.
		Query().
		Where(book.HasOwnerWith(user.NickName(name))).
		Where(book.HasOwnerWith(user.IsPublished(true))).
		All(context.Background())
	if err == nil {
		logger.Init().Sugar().Infof("유저 닉네임을 통해 유저의 저장된 책 목록들을 가져왔습니다: %s", name)
		result := make([]*domain.Book, 0, len(books))
		for _, b := range books {
			result = append(result, &domain.Book{
				ID:           b.ID,
				OwnerID:      b.Edges.Owner.ID, // 이제 b.Edges.Owner는 nil이 아닙니다.
				Title:        b.BookTitle,
				Author:       b.Author,
				BookISBN:     b.BookIsbn,
				RegisteredAt: b.RegisteredAt,
				ComplatedAt:  b.ComplatedAt,
			})
		}

	}

	switch {
	case ent.IsNotFound(err):
		return nil, fmt.Errorf("등록된 책을 찾을 수 없습니다: %w", err)
	case ent.IsNotLoaded(err):
		return nil, fmt.Errorf("책의 목록을 불러오는 도중 데이터가 로드되지 않았습니다: %w", err)
	case ent.IsConstraintError(err):
		return nil, fmt.Errorf("책의 목록을 가져오는 도중 제약 조건 오류가 발생했습니다: %w", err)
	default:
		return nil, fmt.Errorf("책의 목록을 가져오는 도중 알 수 없는 오류가 발생했습니다: %w", err)
	}
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
	if err == nil {
		return nil
	}

	switch {
	case ent.IsNotFound(err):
		return fmt.Errorf("등록된 책을 찾을 수 없습니다: %w", err)
	case ent.IsConstraintError(err):
		return fmt.Errorf("책을 수정하는 도중 제약 조건 오류가 발생했습니다: %w", err)
	default:
		return fmt.Errorf("책을 수정하는 도중 오류가 발생했습니다: %w", err)
	}
}

func (bc *BookRepository) DeleteByID(userID, id uuid.UUID) error {
	client := bc.client

	err := client.Book.DeleteOneID(id).Exec(context.Background())
	if err == nil {
		return nil
	}

	switch {
	case ent.IsNotFound(err):
		return fmt.Errorf("등록된 책을 찾을 수 없습니다: %w", err)
	case ent.IsConstraintError(err):
		return fmt.Errorf("책을 삭제하는 도중 제약 조건 오류가 발생했습니다: %w", err)
	default:
		return fmt.Errorf("책을 삭제하는 도중 오류가 발생했습니다: %w", err)
	}
}
