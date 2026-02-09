package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dev-hyunsang/home-library/internal/config"
	"github.com/dev-hyunsang/home-library/internal/domain"
	"github.com/dev-hyunsang/home-library/lib/ent"
	"github.com/dev-hyunsang/home-library/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type BookHandler struct {
	bookUseCase domain.BookUseCase
	AuthHandler domain.AuthUseCase
}

type SaveBookRequest struct {
	Title        string `json:"title"`
	Author       string `json:"author"`
	BookISBN     string `json:"book_isbn"`
	ThumbnailURL string `json:"thumbnail_url"`
	Status       int    `json:"status"` // 0: 읽지 않음, 1: 읽는 중, 2: 읽음
}

type SearchBookRequest struct {
	BookISBN string `json:"book_isbn"`
}

type ApiResponse struct {
	LastBuildDate string `json:"lastBuildDate"`
	Total         int    `json:"total"`
	Start         int    `json:"start"`
	Display       int    `json:"display"`
	Items         []Book `json:"items"`
}

type Book struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Image       string `json:"image"`
	Author      string `json:"author"`
	Discount    string `json:"discount"`
	Publisher   string `json:"publisher"`
	PubDate     string `json:"pubdate"`
	ISBN        string `json:"isbn"`
	Description string `json:"description"`
}

type BookReviewRequest struct {
	BookISBN string `json:"book_isbn"`
	Content  string `json:"content"`
	Rating   int    `json:"rating"`
	IsPublic bool   `json:"is_public"`
}

func NewBookHandler(bookUseCase domain.BookUseCase, AuthHandler domain.AuthUseCase) *BookHandler {
	return &BookHandler{
		bookUseCase: bookUseCase,
		AuthHandler: AuthHandler,
	}
}

func (h *BookHandler) SaveBookHandler(ctx *fiber.Ctx) error {
	// JWT 토큰에서 사용자 ID 추출
	userID, err := h.AuthHandler.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("JWT 토큰을 통한 사용자 인증에 실패했습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	book := new(SaveBookRequest)
	if err := ctx.BodyParser(book); err != nil {
		logger.Init().Sugar().Errorf("요청 본문을 파싱하는 도중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	createdBook := &domain.Book{
		ID:           uuid.New(),
		OwnerID:      userID,
		Title:        book.Title,
		Author:       book.Author,
		BookISBN:     book.BookISBN,
		ThumbnailURL: book.ThumbnailURL,
		Status:       book.Status,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	result, err := h.bookUseCase.SaveByBookID(userID, createdBook)
	if err != nil {
		logger.Init().Sugar().Errorf("책을 저장하는 도중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	logger.Init().Sugar().Infof("책이 성공적으로 저장되었습니다 / 책ID: %s, 사용자ID: %s", result.ID.String(), userID.String())

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"is_success":   true,
		"data":         result,
		"responsed_at": time.Now(),
	})
}

func (h *BookHandler) GetBooksHandler(ctx *fiber.Ctx) error {
	// JWT 토큰에서 사용자 ID 추출
	userID, err := h.AuthHandler.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("JWT 토큰을 통한 사용자 인증에 실패했습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	books, err := h.bookUseCase.GetBooksByUserID(userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	logger.Init().Sugar().Infof("사용자의 책 목록을 성공적으로 조회했습니다. / 사용자ID: %s", userID.String())

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_success":   true,
		"data":         books,
		"responsed_at": time.Now(),
	})
}

func (h *BookHandler) GetBooksByUserNameHandler(ctx *fiber.Ctx) error {
	name := ctx.Params("name")
	if len(name) == 0 {
		logger.Init().Sugar().Error("사용자 이름이 입력되지 않았습니다.")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	books, err := h.bookUseCase.GetBooksByUserName(name)
	if err != nil {
		if ent.IsNotFound(err) || errors.Is(err, domain.ErrNotFound) {
			logger.Init().Sugar().Errorf("등록된 책을 찾을 수 없습니다: %v", err)
			return ctx.Status(fiber.StatusNotFound).JSON(ErrorHandler(domain.ErrNotFound))
		}
		logger.Init().Sugar().Errorf("책의 목록을 가져오는 도중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_success":   true,
		"data":         books,
		"responsed_at": time.Now(),
	})
}

func (h *BookHandler) BookDeleteHandler(ctx *fiber.Ctx) error {
	// JWT 토큰에서 사용자 ID 추출
	userID, err := h.AuthHandler.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("JWT 토큰을 통한 사용자 인증에 실패했습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	id := ctx.Params("id")
	if len(id) == 0 {
		logger.Init().Sugar().Error("책 ID가 입력되지 않았습니다.")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	err = h.bookUseCase.DeleteByID(userID, uuid.MustParse(id))
	if err != nil {
		if ent.IsNotFound(err) {
			logger.Init().Sugar().Errorf("등록된 책을 찾을 수 없습니다: %v", err)
			return ctx.Status(fiber.StatusNotFound).JSON(ErrorHandler(domain.ErrNotFound))
		}
		logger.Init().Sugar().Errorf("책을 삭제하는 도중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	logger.Init().Sugar().Infof("책이 성공적으로 삭제되었습니다 / 책ID: %s, 사용자ID: %s", id, userID.String())

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *BookHandler) SearchBookIsbnHandler(ctx *fiber.Ctx) error {
	// JWT 토큰에서 사용자 ID 추출
	userID, err := h.AuthHandler.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("JWT 토큰을 통한 사용자 인증에 실패했습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	req := new(SearchBookRequest)
	if err := ctx.BodyParser(req); err != nil {
		logger.Init().Sugar().Errorf("요청 본문을 파싱하는 도중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	logger.Init().Sugar().Infof("검색하는 책 ISBN: %s", req.BookISBN)

	// 네이버 OpenAI를 사용하여 ISBN과 동일한 서적을 검색합니다.
	// 한 개의 결과값만 나오며, 그 결과을 반환합니다.
	searchURL := fmt.Sprintf("https://openapi.naver.com/v1/search/book.json?query=%s&display=10&start=1&sort=sim", req.BookISBN)

	bookReq, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		logger.Init().Sugar().Errorf("책 검색 API 요청 생성 실패: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	bookReq.Header.Add("X-Naver-Client-Id", config.GetEnv("NAVER_API_CLIENT_ID"))
	bookReq.Header.Add("X-Naver-Client-Secret", config.GetEnv("NAVER_API_CLIENT_SECRET"))

	client := &http.Client{}
	resp, err := client.Do(bookReq)
	if err != nil {
		logger.Init().Sugar().Errorf("책 검색 API 요청 실패: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Init().Sugar().Errorf("책 검색 API 응답 본문 읽기 실패: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	var res ApiResponse
	err = json.Unmarshal([]byte(body), &res)
	if err != nil {
		logger.Init().Sugar().Errorf("책 검색 API 응답 본문을 언마샬하는 도중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	logger.Init().Sugar().Infof("네이버 책 검색 API 요청이 성공적으로 완료되었습니다. / 사용자ID: %s", userID.String())
	logger.Init().Sugar().Infof("네이버 책 검색 API 응답: %v", res)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_success":   true,
		"data":         res,
		"responsed_at": time.Now(),
	})
}

// Book Review

// SaveBookReviewHandler ISBN을 기반으로 책 리뷰를 작성합니다.
// 필수 필드: book_isbn (책 ISBN), content (리뷰 내용), rating (별점 1~5), is_public (공개 여부)
func (h *BookHandler) SaveBookReviewHandler(ctx *fiber.Ctx) error {
	req := new(BookReviewRequest)
	if err := ctx.BodyParser(req); err != nil {
		logger.Init().Sugar().Errorf("요청 본문을 파싱하는 도중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	// 필수 필드 검증
	if req.BookISBN == "" {
		logger.Init().Sugar().Error("책 ISBN이 입력되지 않았습니다.")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}
	if req.Rating < 1 || req.Rating > 5 {
		logger.Init().Sugar().Error("별점은 1~5 사이의 값이어야 합니다.")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}
	if req.Content == "" {
		logger.Init().Sugar().Error("리뷰 내용이 입력되지 않았습니다.")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	// JWT 토큰에서 사용자 ID 추출
	userID, err := h.AuthHandler.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("JWT 토큰을 통한 사용자 인증에 실패했습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	// ISBN으로 책 조회 (본인 책 우선, 없으면 다른 사용자의 책도 허용)
	book, err := h.bookUseCase.GetBookByISBN(userID, req.BookISBN)
	if err != nil {
		// 본인 책이 없으면 다른 사용자의 책 조회
		book, err = h.bookUseCase.GetAnyBookByISBN(req.BookISBN)
		if err != nil {
			logger.Init().Sugar().Errorf("해당 ISBN으로 등록된 책을 찾을 수 없습니다: %v", err)
			return ctx.Status(fiber.StatusNotFound).JSON(ErrorHandler(domain.ErrNotFound))
		}
	}

	reviewID := uuid.New()
	if err = h.bookUseCase.CreateReview(&domain.ReviewBook{
		ID:        reviewID,
		BookID:    book.ID,
		OwnerID:   userID,
		BookISBN:  req.BookISBN,
		Content:   req.Content,
		Rating:    req.Rating,
		IsPublic:  req.IsPublic,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}); err != nil {
		logger.Init().Sugar().Errorf("책 리뷰 생성 중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	logger.Init().Sugar().Infof("책 리뷰가 성공적으로 작성되었습니다 / 리뷰ID: %s, ISBN: %s", reviewID.String(), req.BookISBN)

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"is_success": true,
		"data": fiber.Map{
			"review_id": reviewID,
			"book_isbn": req.BookISBN,
			"rating":    req.Rating,
			"is_public": req.IsPublic,
		},
		"responsed_at": time.Now(),
	})
}

func (h *BookHandler) GetBookReviewByUserIDHandler(ctx *fiber.Ctx) error {
	// JWT 토큰에서 사용자 ID 추출
	userID, err := h.AuthHandler.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("JWT 토큰을 통한 사용자 인증에 실패했습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	results, err := h.bookUseCase.GetReviewsByUserID(userID)
	if err != nil {
		logger.Init().Sugar().Errorf("책 리뷰 목록을 가져오는 도중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_success":   true,
		"data":         results,
		"responsed_at": time.Now(),
	})
}

func (h *BookHandler) GetPublicReviewsByBookIDHandler(ctx *fiber.Ctx) error {
	bookID := ctx.Params("book_id")
	if len(bookID) == 0 {
		logger.Init().Sugar().Error("책 ID가 입력되지 않았습니다.")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	results, err := h.bookUseCase.GetPublicReviewsByBookID(uuid.MustParse(bookID))
	if err != nil {
		logger.Init().Sugar().Errorf("공개 리뷰 목록을 가져오는 도중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_success":   true,
		"data":         results,
		"responsed_at": time.Now(),
	})
}

// GetPublicReviewsByISBNHandler ISBN을 통해 공개 리뷰 목록을 조회합니다.
func (h *BookHandler) GetPublicReviewsByISBNHandler(ctx *fiber.Ctx) error {
	isbn := ctx.Params("isbn")
	if len(isbn) == 0 {
		logger.Init().Sugar().Error("책 ISBN이 입력되지 않았습니다.")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	results, err := h.bookUseCase.GetPublicReviewsByISBN(isbn)
	if err != nil {
		logger.Init().Sugar().Errorf("공개 리뷰 목록을 가져오는 도중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_success":   true,
		"data":         results,
		"responsed_at": time.Now(),
	})
}

// DeleteBookReviewHandler 리뷰를 삭제합니다. 본인의 리뷰만 삭제할 수 있습니다.
func (h *BookHandler) DeleteBookReviewHandler(ctx *fiber.Ctx) error {
	reviewID := ctx.Params("id")
	if len(reviewID) == 0 {
		logger.Init().Sugar().Error("리뷰 ID가 입력되지 않았습니다.")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	// JWT 토큰에서 사용자 ID 추출
	userID, err := h.AuthHandler.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("JWT 토큰을 통한 사용자 인증에 실패했습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	err = h.bookUseCase.DeleteReviewByID(userID, uuid.MustParse(reviewID))
	if err != nil {
		logger.Init().Sugar().Errorf("리뷰 삭제 중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	logger.Init().Sugar().Infof("리뷰가 성공적으로 삭제되었습니다 / 리뷰ID: %s", reviewID)

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *BookHandler) AddBookmarkHandler(ctx *fiber.Ctx) error {
	// JWT 토큰에서 사용자 ID 추출
	userID, err := h.AuthHandler.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("JWT 토큰을 통한 사용자 인증에 실패했습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	BookID := ctx.Params("id")

	result, err := h.bookUseCase.AddBookmarkByBookID(userID, uuid.MustParse(BookID))
	if err != nil {
		logger.Init().Sugar().Errorf("책 북마크 추가 중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_success":   true,
		"data":         result,
		"responsed_at": time.Now(),
	})
}

func (h *BookHandler) GetBookmarksByUserIDHandler(ctx *fiber.Ctx) error {
	// JWT 토큰에서 사용자 ID 추출
	userID, err := h.AuthHandler.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("JWT 토큰을 통한 사용자 인증에 실패했습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	results, err := h.bookUseCase.GetBookmarksByUserID(userID)
	if err != nil {
		logger.Init().Sugar().Errorf("책 북마크 목록을 가져오는 도중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_success":   true,
		"data":         results,
		"responsed_at": time.Now(),
	})
}

func (h *BookHandler) DeleteBookmarkHandler(ctx *fiber.Ctx) error {
	_, err := h.AuthHandler.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("JWT 토큰을 통한 사용자 인증에 실패했습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	BookID := ctx.Params("id")

	err = h.bookUseCase.DeleteBookmarkByID(uuid.MustParse(BookID))
	if err != nil {
		logger.Init().Sugar().Errorf("서적의 북마크 제거하던 도중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
