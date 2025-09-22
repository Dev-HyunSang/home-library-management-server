package memory

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dev-hyunsang/home-library/internal/domain"
	"github.com/dev-hyunsang/home-library/lib/ent"
	"github.com/dev-hyunsang/home-library/lib/ent/book"
	"github.com/dev-hyunsang/home-library/lib/ent/review"
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
func (rc *BookRepository) GetBookByID(userID, id uuid.UUID) (*domain.Book, error) {
	client := rc.client

	result, err := client.Book.
		Query().
		Where(
			book.ID(id),
			book.HasOwnerWith(user.ID(userID))). // UserID와 일치하는 조건을 찾습니다.
		Only(context.Background())

	logger.UserInfoLog(userID.String(), "해당 유저의 책 정보를 가지고 왔습니다.")

	if err != nil {
		return nil, fmt.Errorf("책 정보를 가져오는 도중 오류가 발생했습니다: %w", err)
	} else if ent.IsNotFound(err) {
		return nil, fmt.Errorf("해당 정보로 등록된 책을 찾을 수 없습니다: %w", err)
	}

	log.Println(result)

	return &domain.Book{
		ID:           result.ID,
		OwnerID:      result.QueryOwner().OnlyIDX(context.Background()),
		Title:        result.BookTitle,
		Author:       result.Author,
		BookISBN:     result.BookIsbn,
		RegisteredAt: result.RegisteredAt,
		ComplatedAt:  result.ComplatedAt,
	}, nil

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
		WithOwner(). // Owner 관계를 명시적으로 로드
		All(context.Background())
	if err != nil {
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

	logger.Init().Sugar().Infof("유저 닉네임을 통해 유저의 저장된 책 목록들을 가져왔습니다: %s", name)
	result := make([]*domain.Book, 0, len(books))
	for _, b := range books {
		// Owner가 로드되었는지 확인
		if b.Edges.Owner == nil {
			logger.Init().Sugar().Errorf("책 ID %s의 Owner가 로드되지 않았습니다", b.ID.String())
			continue
		}

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

func (bc *BookRepository) CreateReview(review *domain.ReviewBook) error {
	client := bc.client

	_, err := client.Review.Create().
		SetID(review.ID).
		SetBookID(review.BookID).
		SetOwnerID(review.OwnerID).
		SetContent(review.Content).
		SetRating(review.Rating).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		Save(context.Background())
	if err != nil {
		return fmt.Errorf("리뷰를 저장하는 도중 오류가 발생했습니다: %w", err)
	}

	logger.Init().Sugar().Infof("리뷰를 성공적으로 저장했습니다. 리뷰 ID: %s", review.ID.String())

	return nil
}

func (bc *BookRepository) GetReviewsByUserID(userID uuid.UUID) ([]*domain.ReviewBook, error) {
	var result []*domain.ReviewBook

	client := bc.client

	reviews, err := client.Review.Query().
		Where(review.HasOwnerWith(user.ID(userID))).
		All(context.Background())
	if err != nil {
		return nil, fmt.Errorf("책 리뷰 목록을 가져오는 도중 오류가 발생했습니다: %w", err)
	}

	for _, r := range reviews {
		result = append(result, &domain.ReviewBook{
			ID:         r.ID,
			BookID:     r.QueryBook().OnlyIDX(context.Background()),
			OwnerID:    r.QueryOwner().OnlyIDX(context.Background()),
			BookTitle:  r.QueryBook().OnlyX(context.Background()).BookTitle,
			BookAuthor: r.QueryBook().OnlyX(context.Background()).Author,
			Content:    r.Content,
			Rating:     r.Rating,
			CreatedAt:  r.CreatedAt,
			UpdatedAt:  r.UpdatedAt,
		})
		log.Println(result)
	}

	return result, nil
}

func (bc *BookRepository) GetReviewByID(reviewID uuid.UUID) (*domain.ReviewBook, error) {
	client := bc.client

	r, err := client.Review.Get(context.Background(), reviewID)
	if err != nil {
		return nil, fmt.Errorf("리뷰를 가져오는 도중 오류가 발생했습니다: %w", err)
	}

	return &domain.ReviewBook{
		ID:         r.ID,
		BookID:     r.QueryBook().OnlyIDX(context.Background()),
		OwnerID:    r.QueryOwner().OnlyIDX(context.Background()),
		BookTitle:  r.QueryBook().OnlyX(context.Background()).BookTitle,
		BookAuthor: r.QueryBook().OnlyX(context.Background()).Author,
		Content:    r.Content,
		Rating:     r.Rating,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}, nil
}

func (bc *BookRepository) UpdateReviewByID(review *domain.ReviewBook) (domain.ReviewBook, error) {
	client := bc.client

	result, err := client.Review.UpdateOneID(review.ID).
		SetContent(review.Content).
		SetRating(review.Rating).
		SetUpdatedAt(time.Now()).
		Save(context.Background())
	if err != nil {
		return domain.ReviewBook{}, fmt.Errorf("리뷰를 수정하는 도중 오류가 발생했습니다: %w", err)
	}

	return domain.ReviewBook{
		ID:        result.ID,
		BookID:    result.QueryBook().QueryOwner().OnlyIDX(context.Background()),
		OwnerID:   result.QueryOwner().OnlyIDX(context.Background()),
		Content:   result.Content,
		Rating:    result.Rating,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, nil
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
