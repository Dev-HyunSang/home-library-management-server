package memory

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dev-hyunsang/home-library-backend/internal/domain"
	"github.com/dev-hyunsang/home-library-backend/lib/ent"
	"github.com/dev-hyunsang/home-library-backend/lib/ent/book"
	"github.com/dev-hyunsang/home-library-backend/lib/ent/bookmark"
	"github.com/dev-hyunsang/home-library-backend/lib/ent/review"
	"github.com/dev-hyunsang/home-library-backend/lib/ent/user"
	"github.com/dev-hyunsang/home-library-backend/logger"
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

	// Check if user exists before creating book
	exists, err := client.User.Query().Where(user.ID(userID)).Exist(context.Background())
	if err != nil {
		logger.Init().Sugar().Errorf("사용자 존재 확인 중 오류가 발생했습니다: %w", err)
		return nil, fmt.Errorf("사용자 존재 확인 중 오류가 발생했습니다: %w", err)
	}
	if !exists {
		logger.Init().Sugar().Errorf("존재하지 않는 사용자입니다: %s", userID.String())
		return nil, fmt.Errorf("존재하지 않는 사용자입니다: %s", userID.String())
	}

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
		SetThumbnailURL(book.ThumbnailURL).
		SetStatus(book.Status).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		Save(context.Background())
	if err == nil {
		logger.UserInfoLog(userID.String(), "해당 유저의 새로운 책을 저장했습니다.")
		return &domain.Book{
			ID:           b.ID,
			OwnerID:      b.QueryOwner().OnlyIDX(context.Background()),
			Title:        b.BookTitle,
			Author:       b.Author,
			BookISBN:     b.BookIsbn,
			ThumbnailURL: b.ThumbnailURL,
			Status:       b.Status,
			CreatedAt:    b.CreatedAt,
			UpdatedAt:    b.UpdatedAt,
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
		ThumbnailURL: result.ThumbnailURL,
		Status:       result.Status,
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}, nil
}

// GetBookByISBN ISBN을 통해 사용자의 책을 조회합니다.
func (rc *BookRepository) GetBookByISBN(userID uuid.UUID, isbn string) (*domain.Book, error) {
	client := rc.client

	result, err := client.Book.
		Query().
		Where(
			book.BookIsbn(isbn),
			book.HasOwnerWith(user.ID(userID)),
		).
		Only(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("해당 ISBN으로 등록된 책을 찾을 수 없습니다: %w", err)
		}
		return nil, fmt.Errorf("책 정보를 가져오는 도중 오류가 발생했습니다: %w", err)
	}

	logger.UserInfoLog(userID.String(), fmt.Sprintf("ISBN을 통해 책 정보를 조회했습니다. ISBN: %s", isbn))

	return &domain.Book{
		ID:           result.ID,
		OwnerID:      result.QueryOwner().OnlyIDX(context.Background()),
		Title:        result.BookTitle,
		Author:       result.Author,
		BookISBN:     result.BookIsbn,
		ThumbnailURL: result.ThumbnailURL,
		Status:       result.Status,
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}, nil
}

// GetAnyBookByISBN ISBN으로 등록된 책을 조회합니다 (소유자 무관).
// 다른 사용자가 등록한 책에도 리뷰를 작성할 수 있도록 지원합니다.
func (rc *BookRepository) GetAnyBookByISBN(isbn string) (*domain.Book, error) {
	client := rc.client

	result, err := client.Book.
		Query().
		Where(book.BookIsbn(isbn)).
		WithOwner().
		First(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("해당 ISBN으로 등록된 책을 찾을 수 없습니다: %w", err)
		}
		return nil, fmt.Errorf("책 정보를 가져오는 도중 오류가 발생했습니다: %w", err)
	}

	ownerID := uuid.Nil
	if result.Edges.Owner != nil {
		ownerID = result.Edges.Owner.ID
	}

	return &domain.Book{
		ID:           result.ID,
		OwnerID:      ownerID,
		Title:        result.BookTitle,
		Author:       result.Author,
		BookISBN:     result.BookIsbn,
		ThumbnailURL: result.ThumbnailURL,
		Status:       result.Status,
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
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
				ThumbnailURL: b.ThumbnailURL,
				Status:       b.Status,
				CreatedAt:    b.CreatedAt,
				UpdatedAt:    b.UpdatedAt,
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
			CreatedAt:    b.CreatedAt,
			UpdatedAt:    b.UpdatedAt,
			ThumbnailURL: b.ThumbnailURL,
			Status:       b.Status,
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
		SetThumbnailURL(book.ThumbnailURL).
		SetStatus(book.Status).
		SetUpdatedAt(time.Now()).
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
		SetIsPublic(review.IsPublic).
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
		WithBook().
		All(context.Background())
	if err != nil {
		return nil, fmt.Errorf("책 리뷰 목록을 가져오는 도중 오류가 발생했습니다: %w", err)
	}

	for _, r := range reviews {
		b := r.Edges.Book
		if b == nil {
			continue
		}
		result = append(result, &domain.ReviewBook{
			ID:         r.ID,
			BookID:     b.ID,
			OwnerID:    userID,
			BookISBN:   b.BookIsbn,
			BookTitle:  b.BookTitle,
			BookAuthor: b.Author,
			Content:    r.Content,
			Rating:     r.Rating,
			IsPublic:   r.IsPublic,
			CreatedAt:  r.CreatedAt,
			UpdatedAt:  r.UpdatedAt,
		})
	}

	return result, nil
}

func (bc *BookRepository) GetReviewByID(reviewID uuid.UUID) (*domain.ReviewBook, error) {
	client := bc.client

	r, err := client.Review.Query().
		Where(review.ID(reviewID)).
		WithBook().
		WithOwner().
		Only(context.Background())
	if err != nil {
		return nil, fmt.Errorf("리뷰를 가져오는 도중 오류가 발생했습니다: %w", err)
	}

	b := r.Edges.Book
	o := r.Edges.Owner
	if b == nil || o == nil {
		return nil, fmt.Errorf("리뷰의 책 또는 작성자 정보를 불러올 수 없습니다")
	}

	return &domain.ReviewBook{
		ID:         r.ID,
		BookID:     b.ID,
		OwnerID:    o.ID,
		BookISBN:   b.BookIsbn,
		BookTitle:  b.BookTitle,
		BookAuthor: b.Author,
		Content:    r.Content,
		Rating:     r.Rating,
		IsPublic:   r.IsPublic,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}, nil
}

func (bc *BookRepository) UpdateReviewByID(review *domain.ReviewBook) (domain.ReviewBook, error) {
	client := bc.client

	result, err := client.Review.UpdateOneID(review.ID).
		SetContent(review.Content).
		SetRating(review.Rating).
		SetIsPublic(review.IsPublic).
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
		IsPublic:  result.IsPublic,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

func (bc *BookRepository) GetPublicReviewsByBookID(bookID uuid.UUID) ([]*domain.ReviewBook, error) {
	var result []*domain.ReviewBook

	client := bc.client

	reviews, err := client.Review.Query().
		Where(review.HasBookWith(book.ID(bookID))).
		Where(review.IsPublic(true)).
		WithBook().
		WithOwner().
		All(context.Background())
	if err != nil {
		return nil, fmt.Errorf("공개 리뷰 목록을 가져오는 도중 오류가 발생했습니다: %w", err)
	}

	for _, r := range reviews {
		b := r.Edges.Book
		o := r.Edges.Owner
		if b == nil || o == nil {
			continue
		}
		result = append(result, &domain.ReviewBook{
			ID:         r.ID,
			BookID:     b.ID,
			OwnerID:    o.ID,
			BookISBN:   b.BookIsbn,
			BookTitle:  b.BookTitle,
			BookAuthor: b.Author,
			Content:    r.Content,
			Rating:     r.Rating,
			IsPublic:   r.IsPublic,
			CreatedAt:  r.CreatedAt,
			UpdatedAt:  r.UpdatedAt,
		})
	}

	return result, nil
}

// DeleteReviewByID 리뷰를 삭제합니다. 본인의 리뷰만 삭제할 수 있습니다.
func (bc *BookRepository) DeleteReviewByID(userID, reviewID uuid.UUID) error {
	client := bc.client

	// 리뷰가 해당 사용자의 것인지 확인
	r, err := client.Review.Query().
		Where(review.ID(reviewID)).
		Where(review.HasOwnerWith(user.ID(userID))).
		Only(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("삭제할 리뷰를 찾을 수 없거나 권한이 없습니다: %w", err)
		}
		return fmt.Errorf("리뷰 조회 중 오류가 발생했습니다: %w", err)
	}

	err = client.Review.DeleteOne(r).Exec(context.Background())
	if err != nil {
		return fmt.Errorf("리뷰를 삭제하는 도중 오류가 발생했습니다: %w", err)
	}

	logger.Init().Sugar().Infof("리뷰가 성공적으로 삭제되었습니다. 리뷰 ID: %s", reviewID.String())
	return nil
}

// GetPublicReviewsByISBN ISBN을 통해 공개 리뷰 목록을 조회합니다.
func (bc *BookRepository) GetPublicReviewsByISBN(isbn string) ([]*domain.ReviewBook, error) {
	var result []*domain.ReviewBook

	client := bc.client

	reviews, err := client.Review.Query().
		Where(review.HasBookWith(book.BookIsbn(isbn))).
		Where(review.IsPublic(true)).
		WithBook().
		WithOwner().
		All(context.Background())
	if err != nil {
		return nil, fmt.Errorf("공개 리뷰 목록을 가져오는 도중 오류가 발생했습니다: %w", err)
	}

	for _, r := range reviews {
		b := r.Edges.Book
		o := r.Edges.Owner
		if b == nil || o == nil {
			continue
		}
		result = append(result, &domain.ReviewBook{
			ID:         r.ID,
			BookID:     b.ID,
			OwnerID:    o.ID,
			BookISBN:   b.BookIsbn,
			BookTitle:  b.BookTitle,
			BookAuthor: b.Author,
			Content:    r.Content,
			Rating:     r.Rating,
			IsPublic:   r.IsPublic,
			CreatedAt:  r.CreatedAt,
			UpdatedAt:  r.UpdatedAt,
		})
	}

	return result, nil
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

func (bc *BookRepository) AddBookmarkByBookID(ownerID, bookID uuid.UUID) (*domain.Bookmark, error) {
	client := bc.client

	result, err := client.Bookmark.Create().
		SetOwnerID(ownerID).
		SetBookID(bookID).
		SetCreatedAt(time.Now()).
		Save(context.Background())
	if err != nil {
		return nil, fmt.Errorf("북마크를 추가하는 도중 오류가 발생했습니다: %w", err)
	}

	return &domain.Bookmark{
		ID:        result.ID,
		OwnerID:   result.QueryOwner().OnlyIDX(context.Background()),
		BookID:    result.QueryBook().OnlyIDX(context.Background()),
		CreatedAt: result.CreatedAt,
	}, nil
}

func (bc *BookRepository) GetBookmarksByUserID(userID uuid.UUID) ([]*domain.Bookmark, error) {
	var result []*domain.Bookmark

	client := bc.client

	bookmarks, err := client.Bookmark.Query().
		Where(bookmark.HasOwnerWith(user.ID(userID))).
		All(context.Background())
	if err != nil {
		return nil, fmt.Errorf("북마크 목록을 가져오는 도중 오류가 발생했습니다: %w", err)
	}

	for _, b := range bookmarks {
		result = append(result, &domain.Bookmark{
			ID:        b.ID,
			OwnerID:   b.QueryOwner().OnlyIDX(context.Background()),
			BookID:    b.QueryBook().OnlyIDX(context.Background()),
			CreatedAt: b.CreatedAt,
		})
		log.Println(result)
	}

	return result, nil
}

func (bc *BookRepository) DeleteBookmarkByID(id uuid.UUID) error {
	client := bc.client

	err := client.Bookmark.DeleteOneID(id).Exec(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("등록된 북마크를 찾을 수 없습니다: %w", err)
		}
		return fmt.Errorf("북마크를 삭제하는 도중 오류가 발생했습니다: %w", err)
	}

	return nil
}
