package handler

import (
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"regexp"
	"time"

	"github.com/dev-hyunsang/home-library/internal/config"
	"github.com/dev-hyunsang/home-library/internal/domain"
	"github.com/dev-hyunsang/home-library/lib/ent"
	"github.com/dev-hyunsang/home-library/logger"
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

type ErrResponse struct {
	IsSuccess bool   `json:"is_success"`
	Message   string `json:"message"`
	Time      string `json:"time"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

func NewUserHandler(userUseCase domain.UserUseCase, authUseCase domain.AuthUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		AuthHandler: authUseCase,
	}
}
func IsValidNickname(nickname string) bool {
	matched, _ := regexp.MatchString(`^[a-z._]+$`, nickname)
	return matched
}

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

func ErrorHandler(err error) ErrResponse {
	return ErrResponse{
		IsSuccess: false,
		Message:   err.Error(),
		Time:      time.Now().String(),
	}

}

func (h *UserHandler) UserSignUpHandler(ctx *fiber.Ctx) error {
	user := new(RegisterationRequest)
	if err := ctx.BodyParser(user); err != nil {
		logger.Init().Sugar().Errorf("올바르지 않은 요청입니다: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	if len(user.NickName) == 0 || len(user.Email) == 0 || len(user.Password) == 0 {
		logger.Init().Sugar().Warn("회원가입에 필수적인 필드에 입력값이 없습니다.")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	// 닉네임 유효성 검사
	if !IsValidNickname(user.NickName) {
		logger.Init().Sugar().Errorf("사용자가 유효하지 않은 닉네임을 입력했습니다: %s", user.NickName)
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidNickname))
	}

	result, err := h.userUseCase.CreateUser(&domain.User{
		NickName:    user.NickName,
		Email:       user.Email,
		Password:    user.Password,
		IsPublished: user.IsPublished,
	})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	logger.Init().Sugar().Infof("새로운 유저가 데이터베이스 상에 정상적으로 생성되었습니다 / 사용자ID: %s", result.ID.String())

	return ctx.Status(fiber.StatusCreated).JSON(result)
}

func (h *UserHandler) UserSignInHandler(ctx *fiber.Ctx) error {
	user := new(LoginRequest)
	if err := ctx.BodyParser(user); err != nil {
		logger.Init().Sugar().Errorf("올바르지 않은 요청입니다: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	result, err := h.userUseCase.GetByEmail(user.Email)
	if err != nil {
		logger.Init().Sugar().Errorf("사용자 이메일로 사용자를 조회하는 도중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	log.Println(result.Password)

	log.Println(user.Password)

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
	if err != nil {
		logger.Init().Sugar().Errorf("비밀번호 비교 중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrInvalidCredentials))
	}

	// 로그인 성공 시 세션 생성
	err = h.AuthHandler.SetSession(result.ID.String(), ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("세션 생성 중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
	}

	// User ID를 세션에 저장하고, 쿠키로도 보냄.
	c := &fiber.Cookie{
		Name:   "auth_token",
		Value:  result.ID.String(),
		Secure: true,
	}

	ctx.Cookie(c)

	logger.Init().Sugar().Infof("사용자가 성공적으로 로그인했습니다 / 사용자ID: %s", result.ID.String())

	return ctx.Status(fiber.StatusOK).JSON(result)
}

func (h *UserHandler) UserGetByIdHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if len(id) == 0 {
		logger.Init().Sugar().Error("사용자 ID가 입력되지 않았습니다.")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	sessionID := ctx.Cookies("auth_token")

	result, err := h.AuthHandler.GetSessionByID(sessionID, ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("세션에 해당하는 쿠키 정보를 찾을 수 없습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	logger.Init().Sugar().Infof("세션 ID를 성공적으로 가져왔습니다: %s", result)

	user, err := h.userUseCase.GetByID(uuid.MustParse(id))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(ErrorHandler(domain.ErrNotFound))
	}

	logger.Init().Sugar().Infof("사용자 정보를 성공적으로 조회했습니다 / 사용자ID: %s", user.ID.String())

	return ctx.Status(fiber.StatusOK).JSON(user)
}
func (h *UserHandler) UserSignOutHandler(ctx *fiber.Ctx) error {
	err := h.AuthHandler.DeleteSession(ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("세션 삭제 중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "successfully logged out"})
}

// 가입된 메일을 통해서 사용자를 찾고, 임시 비밀번호를 생성하여 이메일로 발송하는 핸들러
func (h *UserHandler) UserRestPasswordHandler(ctx *fiber.Ctx) error {
	user := new(ForgotPasswordRequest)
	if err := ctx.BodyParser(user); err != nil {
		logger.Init().Sugar().Errorf("올바르지 않은 요청입니다: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	if len(user.Email) == 0 {
		logger.Init().Sugar().Warn("이메일이 입력되지 않았습니다.")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	// 사용자 존재 확인
	existingUser, err := h.userUseCase.GetByEmail(user.Email)
	if err != nil {
		logger.Init().Sugar().Errorf("사용자를 찾을 수 없습니다: %v", err)
		return ctx.Status(fiber.StatusNotFound).JSON(ErrorHandler(domain.ErrNotFound))
	}

	if ent.IsNotFound(err) {
		logger.Init().Sugar().Warnf("이메일에 해당하는 사용자가 존재하지 않습니다: %s", user.Email)
		return ctx.Status(fiber.StatusNotFound).JSON(ErrorHandler(domain.ErrNotFound))
	}

	// 임시 비밀번호 생성
	tempPassword := stringWithCharset(15)

	// 비밀번호 해시화
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Init().Sugar().Errorf("임시 비밀번호 해시화 중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
	}

	// 기존 사용자 정보를 유지하면서 비밀번호만 업데이트
	err = h.userUseCase.Edit(&domain.User{
		ID:          existingUser.ID,          // 기존 ID 유지
		NickName:    existingUser.NickName,    // 기존 닉네임 유지
		Email:       existingUser.Email,       // 기존 이메일 유지
		Password:    string(hashedPassword),   // 새 비밀번호
		IsPublished: existingUser.IsPublished, // 기존 설정 유지
		CreatedAt:   existingUser.CreatedAt,
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		logger.Init().Sugar().Errorf("사용자 비밀번호 업데이트 중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
	}

	logger.Init().Sugar().Infof("사용자 비밀번호가 성공적으로 업데이트되었습니다 / 사용자ID: %s", existingUser.ID.String())

	mailAuth := smtp.PlainAuth("", config.GetEnv("GOOGLE_MAIL_ADDRESS"), config.GetEnv("GOOGLE_MAIL_PASSWORD"), config.GetEnv("GOOGLE_MAIL_SMTP"))

	from := config.GetEnv("GOOGLE_MAIL_ADDRESS")
	to := []string{user.Email}

	headerSubject := "Subject: 비밀번호 초기화\r\n"
	headerBlank := "\r\n"
	body := fmt.Sprintf("비밀번호 초기화 메일 테스트입니다. %s\r\n", tempPassword)
	msg := []byte(headerSubject + headerBlank + body)

	err = smtp.SendMail("smtp.gmail.com:587", mailAuth, from, to, msg)
	if err != nil {
		logger.Init().Sugar().Errorf("임시 비밀번호 이메일 발송 중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
	}

	logger.Init().Sugar().Infof("임시 비밀번호가 이메일로 발송되었습니다. 이메일: %s", user.Email)

	logger.Init().Sugar().Infof("비밀번호 재설정 요청이 처리되었습니다 / 사용자ID: %s", existingUser.ID.String())

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "비밀번호 재설정 이메일이 발송되었습니다",
		"email":   user.Email,
	})
}

func (h *UserHandler) UserEditHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if len(id) == 0 {
		logger.Init().Sugar().Error("사용자 ID가 입력되지 않았습니다.")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	sessionID := ctx.Cookies("auth_token")
	if len(sessionID) == 0 {
		logger.Init().Sugar().Error("클라이언트측 세션 쿠키가 존재하지 않습니다.")
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	_, err := h.AuthHandler.GetSessionByID(sessionID, ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("세션에 해당하는 쿠키 정보를 찾을 수 없습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(err))
	}

	user := new(domain.User)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	logger.Init().Sugar().Infof("사용자 정보가 성공적으로 업데이트되었습니다 / 사용자ID: %s", user.ID.String())

	return ctx.Status(fiber.StatusOK).JSON(user)
}

func (h *UserHandler) UserVerifyHandler(ctx *fiber.Ctx) error {
	sessionID := ctx.Cookies("auth_token")

	if len(sessionID) == 0 {
		logger.Init().Sugar().Error("클라이언트측 세션 쿠키가 존재하지 않습니다.")
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	userID, err := h.AuthHandler.GetSessionByID(sessionID, ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("세션에 해당하는 쿠키 정보를 찾을 수 없습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	user, err := h.userUseCase.GetByID(uuid.MustParse(userID))
	if err != nil {
		logger.Init().Sugar().Errorf("사용자 정보를 조회하는 도중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
	}

	logger.Init().Sugar().Infof("사용자 인증이 성공적으로 완료되었습니다 / 사용자ID: %s", user.ID.String())

	return ctx.Status(fiber.StatusOK).JSON(user)
}

func (h *UserHandler) UserDeleteHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if len(id) == 0 {
		logger.Init().Sugar().Error("사용자 ID가 입력되지 않았습니다.")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	err := h.userUseCase.Delete(uuid.MustParse(id))
	if err != nil {
		logger.Init().Sugar().Errorf("사용자 정보를 삭제하는 도중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
	}

	logger.Init().Sugar().Infof("사용자 정보가 성공적으로 삭제되었습니다 / 사용자ID: %s", id)
	return ctx.Status(fiber.StatusNoContent).JSON("successfully deleted")
}
