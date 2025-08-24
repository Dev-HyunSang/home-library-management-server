package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	Title    string `json:"title"`
	Author   string `json:"author"`
	BookISBN string `json:"book_isbn"`
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

var userToken string

func TokenSelection(sessionID string) string {
	userPart := strings.Split(sessionID, ";")[0]
	keyValue := strings.Split(userPart, "=")
	if len(keyValue) > 1 {
		return keyValue[1]
	}
	return ""
}

func NewBookHandler(bookUseCase domain.BookUseCase, AuthHandler domain.AuthUseCase) *BookHandler {
	return &BookHandler{
		bookUseCase: bookUseCase,
		AuthHandler: AuthHandler,
	}
}

func (h *BookHandler) SaveBookHandler(ctx *fiber.Ctx) error {
	sessionID := ctx.Cookies("user")
	if len(sessionID) == 0 {
		logger.Init().Sugar().Error("클라이언트측 세션 쿠키가 존재하지 않습니다.")
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	userToken = TokenSelection(sessionID)

	userID, err := h.AuthHandler.GetSessionByID(userToken, ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("세션에 해당하는 쿠키 정보를 찾을 수 없습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	book := new(SaveBookRequest)
	if err := ctx.BodyParser(book); err != nil {
		logger.Init().Sugar().Errorf("요청 본문을 파싱하는 도중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	createdBook := &domain.Book{
		ID:           uuid.New(),
		OwnerID:      uuid.MustParse(userID),
		Title:        book.Title,
		Author:       book.Author,
		BookISBN:     book.BookISBN,
		RegisteredAt: time.Now(),
		ComplatedAt:  time.Time{},
	}

	result, err := h.bookUseCase.SaveByBookID(uuid.MustParse(userID), createdBook)
	if err != nil {
		logger.Init().Sugar().Errorf("책을 저장하는 도중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	logger.Init().Sugar().Infof("책이 성공적으로 저장되었습니다 / 책ID: %s, 사용자ID: %s", result.ID.String(), userID)

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"is_success":   true,
		"data":         result,
		"responsed_at": time.Now(),
	})
}

func (h *BookHandler) GetBooksHandler(ctx *fiber.Ctx) error {
	sessionID := ctx.Cookies("user")
	if len(sessionID) == 0 {
		logger.Init().Sugar().Error("클라이언트측 세션 쿠키가 존재하지 않습니다.")
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	result, err := h.AuthHandler.GetSessionByID(sessionID, ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("세션에 해당하는 쿠키 정보를 찾을 수 없습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	books, err := h.bookUseCase.GetBooksByUserID(uuid.MustParse(result))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	logger.Init().Sugar().Infof("사용자의 책 목록을 성공적으로 조회했습니다. / 사용자ID: %s", result)

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
	sessionID := ctx.Cookies("user")
	if len(sessionID) == 0 {
		logger.Init().Sugar().Error("클라이언트측 세션 쿠키가 존재하지 않습니다.")
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	userID, err := h.AuthHandler.GetSessionByID(sessionID, ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("세션에 해당하는 쿠키 정보를 찾을 수 없습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(err))
	}

	id := ctx.Params("id")
	if len(id) == 0 {
		logger.Init().Sugar().Error("책 ID가 입력되지 않았습니다.")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	err = h.bookUseCase.DeleteByID(uuid.MustParse(userID), uuid.MustParse(id))
	if err != nil {
		if ent.IsNotFound(err) {
			logger.Init().Sugar().Errorf("등록된 책을 찾을 수 없습니다: %v", err)
			return ctx.Status(fiber.StatusNotFound).JSON(ErrorHandler(domain.ErrNotFound))
		}
		logger.Init().Sugar().Errorf("책을 삭제하는 도중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	logger.Init().Sugar().Infof("책이 성공적으로 삭제되었습니다 / 책ID: %s, 사용자ID: %s", id, userID)

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *BookHandler) SearchBookIsbnHandler(ctx *fiber.Ctx) error {
	sessionID := ctx.Cookies("user")
	if len(sessionID) == 0 {
		logger.Init().Sugar().Error("클라이언트측 세션 쿠키가 존재하지 않습니다.")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	result, err := h.AuthHandler.GetSessionByID(userToken, ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("세션에 해당하는 쿠키 정보를 찾을 수 없습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	if len(result) == 0 {
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

	bookReq.Header.Add("X-Naver-Client-Id", config.LoadEnv("NAVER_API_CLIENT_ID"))
	bookReq.Header.Add("X-Naver-Client-Secret", config.LoadEnv("NAVER_API_CLIENT_SECRET"))

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

	logger.Init().Sugar().Infof("네이버 책 검색 API 요청이 성공적으로 완료되었습니다. / 사용자ID: %s", result)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_success":   true,
		"data":         res,
		"responsed_at": time.Now(),
	})
}
