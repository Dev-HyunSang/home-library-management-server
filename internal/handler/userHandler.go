package handler

import (
	"errors"
	"log"

	"github.com/dev-hyunsang/home-library/internal/domain"
	"github.com/dev-hyunsang/home-library/lib/ent"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userUseCase domain.UserUseCase
	AuthHandler domain.AuthUseCase
}

type RegisterationRequest struct {
	NickName    string `json:"nick_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsPublished bool   `json:"is_published"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserHandler(userUseCase domain.UserUseCase, authUseCase domain.AuthUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		AuthHandler: authUseCase,
	}
}

func ErrResponse(err error) map[string]string {
	return map[string]string{
		"status":  "error",
		"message": err.Error(),
	}
}

func (h *UserHandler) UserRegisterHandler(ctx *fiber.Ctx) error {
	user := new(RegisterationRequest)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrResponse(domain.ErrInvalidInput))
	}

	if user.NickName == "" || user.Email == "" || user.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrResponse(domain.ErrInvalidInput))
	}

	result, err := h.userUseCase.CreateUser(&domain.User{
		NickName:    user.NickName,
		Email:       user.Email,
		Password:    user.Password,
		IsPublished: user.IsPublished,
	})
	log.Println(result)

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

func (h *UserHandler) UserLoginHandler(ctx *fiber.Ctx) error {
	user := new(LoginRequest)
	log.Println(user)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrResponse(domain.ErrInvalidInput))
	}

	result, err := h.userUseCase.GetByEmail(user.Email)
	if err != nil {
		switch {
		case ent.IsNotFound(err):
			return ctx.Status(fiber.StatusNotFound).JSON(ErrResponse(err))
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(ErrResponse(domain.ErrInternal))
		}
	}

	if result == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(ErrResponse(domain.ErrNotFound))
	}

	log.Println(result)

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrResponse(domain.ErrInvalidCredentials))
	}

	// 로그인 성공 시 세션 생성
	err = h.AuthHandler.SetSession(result.ID.String(), ctx)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrResponse(domain.ErrInternal))
	}

	// User ID를 세션에 저장하고, 쿠키로도 보냄.
	c := &fiber.Cookie{
		Name:   "user",
		Value:  result.ID.String(),
		Secure: true,
	}

	ctx.Cookie(c)
	return ctx.Status(fiber.StatusOK).JSON(result)
}

func (h *UserHandler) UserGetByIdHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if len(id) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrResponse(domain.ErrInvalidInput))
	}

	sessionID := ctx.Cookies("user")

	result, err := h.AuthHandler.GetSessionByID(sessionID, ctx)
	if err != nil {
		log.Println(err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrResponse(domain.ErrUserNotLoggedIn))
	}

	log.Printf("Get Session ID : %s", result)

	user, err := h.userUseCase.GetByID(uuid.MustParse(id))
	if err == nil {
		log.Println(err)
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

func (h *UserHandler) GetAllHandler(ctx *fiber.Ctx) error {
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

func (h *UserHandler) UserEditHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if len(id) == 0 {
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

func (h *UserHandler) UserDeleteHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if len(id) == 0 {
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
