package handler

import (
	"errors"

	"github.com/dev-hyunsang/home-library/internal/domain"
	"github.com/dev-hyunsang/home-library/lib/ent"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	userUseCase domain.UserUseCase
}

type RegisterationRequest struct {
	NickName string `json:"nick_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserHandler(userUseCase domain.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func ErrResponse(err error) map[string]string {
	return map[string]string{
		"error": err.Error(),
	}
}

func (h *UserHandler) Register(ctx *fiber.Ctx) error {
	user := new(RegisterationRequest)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrResponse(domain.ErrInvalidInput))
	}

	if user.NickName == "" || user.Email == "" || user.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrResponse(domain.ErrInvalidInput))
	}
	result, err := h.userUseCase.CreateUser(&domain.User{
		NickName: user.NickName,
		Email:    user.Email,
		Password: user.Password,
	})
	if err == nil {
		return ctx.Status(fiber.StatusCreated).JSON(result)
	}

	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrResponse(err))
	case errors.Is(err, domain.ErrAlreadyExists):
		return ctx.Status(fiber.StatusConflict).JSON(ErrResponse(err))
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrResponse(domain.ErrInternal))
	}
}

func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	return nil
}

func (h *UserHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrResponse(domain.ErrInvalidInput))
	}

	err := h.userUseCase.Delete(uuid.MustParse(id))
	if err == nil {
		return ctx.Status(fiber.StatusNoContent).JSON("successfully deleted")
	}

	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrResponse(err))
	case errors.Is(err, domain.ErrNotFound):
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrResponse(err))
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrResponse(domain.ErrInternal))
	}
}

func (h *UserHandler) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrResponse(domain.ErrInvalidInput))
	}

	user, err := h.userUseCase.GetByID(uuid.MustParse(id))
	if err == nil {
		return ctx.Status(fiber.StatusOK).JSON(user)
	}

	switch {
	case errors.Is(err, domain.ErrNotFound):
		return ctx.Status(fiber.StatusNotFound).JSON(ErrResponse(err))
	case ent.IsNotFound(err):
		return ctx.Status(fiber.StatusNotFound).JSON(ErrResponse(err))
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrResponse(domain.ErrInternal))
	}
}

func (h *UserHandler) GetAll(ctx *fiber.Ctx) error {
	users, err := h.userUseCase.GetAll()
	if err == nil {
		return ctx.Status(fiber.StatusOK).JSON(users)
	}

	switch {
	case errors.Is(err, domain.ErrNotFound):
		return ctx.Status(fiber.StatusNotFound).JSON(ErrResponse(err))
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrResponse(domain.ErrInternal))
	}
}

func (h *UserHandler) Edit(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrResponse(domain.ErrInvalidInput))
	}

	user := new(domain.User)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrResponse(domain.ErrInvalidInput))
	}

	user, err := h.userUseCase.Edit(user)
	if err == nil {
		return ctx.Status(fiber.StatusOK).JSON(user)
	}

	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrResponse(err))
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrResponse(domain.ErrInternal))
	}
}
