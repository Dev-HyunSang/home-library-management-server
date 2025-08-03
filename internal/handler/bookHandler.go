package handler

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dev-hyunsang/home-library/internal/domain"
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
	BookISBN int    `json:"book_isbn"`
}

func NewBookHandler(bookUseCase domain.BookUseCase, AuthHandler domain.AuthUseCase) *BookHandler {
	return &BookHandler{
		bookUseCase: bookUseCase,
		AuthHandler: AuthHandler,
	}
}

func (h *BookHandler) SaveBookHandler(ctx *fiber.Ctx) error {
	sessionID := ctx.Cookies("user")

	userID, err := h.AuthHandler.GetSessionByID(sessionID, ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrResponse(err))
	}

	book := new(SaveBookRequest)
	if err := ctx.BodyParser(book); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrResponse(domain.ErrInvalidInput))
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
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrResponse(err))
	}

	return ctx.Status(fiber.StatusCreated).JSON(result)
}

func (h *BookHandler) GetBooksHandler(ctx *fiber.Ctx) error {
	sessionID := ctx.Cookies("user")

	result, err := h.AuthHandler.GetSessionByID(sessionID, ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrResponse(err))
	}

	books, err := h.bookUseCase.GetBooksByUserID(uuid.MustParse(result))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrResponse(err))
	}

	if len(books) == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(ErrResponse(fmt.Errorf("%s", "등록되어 있는 책을 찾을 수 없습니다.")))
	}

	return ctx.JSON(books)
}

func (h *BookHandler) GetBooksByUserNameHandler(ctx *fiber.Ctx) error {
	name := ctx.Params("name")
	if len(name) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrResponse(domain.ErrInvalidInput))
	}

	books, err := h.bookUseCase.GetBooksByUserName(name)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrResponse(err))
	}
	if len(books) == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(ErrResponse(fmt.Errorf("%s", "등록되어 있는 책을 찾을 수 없습니다.")))
	}

	switch {
	case errors.Is(err, domain.ErrPrivateAccount):
		return ctx.Status(fiber.StatusForbidden).JSON(ErrResponse(domain.ErrPrivateAccount))
	}

	return ctx.Status(fiber.StatusOK).JSON(books)
}

func (h *BookHandler) BookDeleteHandler(ctx *fiber.Ctx) error {
	sessionID := ctx.Cookies("user")

	userID, err := h.AuthHandler.GetSessionByID(sessionID, ctx)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrResponse(err))
	}

	id := ctx.Params("id")
	if len(id) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrResponse(domain.ErrInvalidInput))
	}

	err = h.bookUseCase.DeleteByID(uuid.MustParse(userID), uuid.MustParse(id))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrResponse(err))
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
