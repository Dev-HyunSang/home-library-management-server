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
	NickName      string `json:"nick_name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	IsPublished   bool   `json:"is_published"`
	IsTermsAgreed bool   `json:"is_terms_agreed"`
}

type UpdateUserRequest struct {
	NickName    string  `json:"nick_name"`
	Email       string  `json:"email"`
	Password    *string `json:"password,omitempty"` // 포인터로 설정하여 nil일 때 비밀번호 변경 안함
	IsPublished bool    `json:"is_published"`
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

	// 이메일 중복 확인
	_, err := h.userUseCase.GetByEmail(user.Email)
	if err == nil {
		logger.Init().Sugar().Warnf("이미 존재하는 이메일로 가입 시도: %s", user.Email)
		return ctx.Status(fiber.StatusConflict).JSON(ErrorHandler(domain.ErrAlreadyExists))
	}

	if !user.IsTermsAgreed {
		logger.Init().Sugar().Warn("회원가입시 이용약관에 동의하지 않았습니다.")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrTermsNotAgreed))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Init().Sugar().Errorf("사용자 비밀번호 해시화 중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
	}

	log.Println(string(hashedPassword))

	result, err := h.userUseCase.Save(&domain.User{
		ID:            uuid.New(),
		NickName:      user.NickName,
		Email:         user.Email,
		Password:      string(hashedPassword),
		IsPublished:   user.IsPublished,
		IsTermsAgreed: user.IsTermsAgreed,
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

	log.Println(result)
	log.Println(user.Password)

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
	if err != nil {
		logger.Init().Sugar().Errorf("비밀번호 비교 중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrInvalidCredentials))
	}

	// JWT 토큰 쌍 생성 (Access + Refresh)
	accessToken, refreshToken, err := h.AuthHandler.GenerateTokenPair(result.ID)
	if err != nil {
		logger.Init().Sugar().Errorf("JWT 토큰 생성 중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
	}

	logger.Init().Sugar().Infof("사용자가 성공적으로 로그인했습니다 / 사용자ID: %s", result.ID.String())

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"user":          result,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    3600, // 1시간 (seconds)
	})
}

func (h *UserHandler) UserGetByIdHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if len(id) == 0 {
		logger.Init().Sugar().Error("사용자 ID가 입력되지 않았습니다.")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	// JWT 토큰에서 사용자 ID 추출
	userIDFromToken, err := h.AuthHandler.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("JWT 토큰을 통한 사용자 인증에 실패했습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	logger.Init().Sugar().Infof("JWT 토큰에서 사용자 ID를 성공적으로 추출했습니다: %s", userIDFromToken.String())

	user, err := h.userUseCase.GetByID(uuid.MustParse(id))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(ErrorHandler(domain.ErrNotFound))
	}

	logger.Init().Sugar().Infof("사용자 정보를 성공적으로 조회했습니다 / 사용자ID: %s", user.ID.String())

	return ctx.Status(fiber.StatusOK).JSON(user)
}
func (h *UserHandler) UserSignOutHandler(ctx *fiber.Ctx) error {
	// JWT 토큰에서 사용자 ID 추출
	userID, err := h.AuthHandler.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("JWT 토큰을 통한 사용자 인증에 실패했습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	// 토큰 무효화
	token, err := h.AuthHandler.ExtractTokenFromHeader(ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("토큰 추출에 실패했습니다: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	if err := h.AuthHandler.InvalidateToken(token); err != nil {
		logger.Init().Sugar().Errorf("토큰 무효화에 실패했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
	}

	logger.Init().Sugar().Infof("사용자가 성공적으로 로그아웃했습니다 / 사용자ID: %s", userID.String())

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "성공적으로 로그아웃되었습니다.",
	})
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
	err = h.userUseCase.Update(&domain.User{
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

func (h *UserHandler) UserVerifyByEmailHandler(ctx *fiber.Ctx) error {
	email := ctx.Params("email")
	if len(email) == 0 {
		logger.Init().Sugar().Error("이메일이 입력되지 않았습니다.")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	result, err := h.userUseCase.GetByEmail(email)
	if err != nil {
		logger.Init().Sugar().Errorf("사용자 정보를 이메일로 조회하는 도중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusNotFound).JSON(ErrorHandler(domain.ErrNotFound))
	}

	if result.Email == email {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"is_success": false,
			"message":    "이미 사용 중인 이메일입니다.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_success": true,
		"message":    "사용 가능한 이메일이며, 해당 이메일로 인증번호를 발송했습니다.",
	})
}

func (h *UserHandler) UserVerifyByNicknameHandler(ctx *fiber.Ctx) error {
	nickname := ctx.Query("nickname")
	if len(nickname) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	if !IsValidNickname(nickname) {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidNickname))
	}

	_, err := h.userUseCase.GetByNickname(nickname)
	if err == nil {
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"is_success": false,
			"message":    "이미 사용 중인 닉네임입니다.",
		})
	}

	if ent.IsNotFound(err) {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"is_success": true,
			"message":    "사용 가능한 닉네임입니다.",
		})
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(err))
}

func (h *UserHandler) UserEditHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if len(id) == 0 {
		logger.Init().Sugar().Error("사용자 ID가 입력되지 않았습니다.")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	// JWT 토큰에서 사용자 ID 추출
	userIDFromToken, err := h.AuthHandler.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("JWT 토큰을 통한 사용자 인증에 실패했습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	// 요청한 ID와 토큰의 사용자 ID가 일치하는지 확인
	if userIDFromToken.String() != id {
		logger.Init().Sugar().Errorf("권한이 없는 사용자 정보 수정 시도: 토큰 사용자 ID %s, 요청 사용자 ID %s", userIDFromToken.String(), id)
		return ctx.Status(fiber.StatusForbidden).JSON(ErrorHandler(domain.ErrPermissionDenied))
	}

	// 기존 사용자 정보 조회
	existingUser, err := h.userUseCase.GetByID(uuid.MustParse(id))
	if err != nil {
		logger.Init().Sugar().Errorf("기존 사용자 정보를 조회하는 도중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
	}

	updateReq := new(UpdateUserRequest)
	if err := ctx.BodyParser(updateReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	// 업데이트할 사용자 정보 구성
	updatedUser := &domain.User{
		ID:          existingUser.ID,
		NickName:    updateReq.NickName,
		Email:       updateReq.Email,
		Password:    existingUser.Password, // 기본값으로 기존 비밀번호 유지
		IsPublished: updateReq.IsPublished,
		CreatedAt:   existingUser.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	// 비밀번호가 제공된 경우에만 해시화 후 업데이트
	if updateReq.Password != nil && *updateReq.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*updateReq.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Init().Sugar().Errorf("비밀번호 해시화 중 오류가 발생했습니다: %v", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
		}
		updatedUser.Password = string(hashedPassword)
		logger.Init().Sugar().Infof("사용자 비밀번호가 변경되었습니다 / 사용자ID: %s", id)
	}

	if err = h.userUseCase.Update(updatedUser); err != nil {
		logger.Init().Sugar().Errorf("사용자 정보를 업데이트하는 도중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
	}

	logger.Init().Sugar().Infof("사용자 정보가 성공적으로 업데이트되었습니다 / 사용자ID: %s", updatedUser.ID.String())

	// 응답에서 비밀번호 제거
	responseUser := &domain.User{
		ID:          updatedUser.ID,
		NickName:    updatedUser.NickName,
		Email:       updatedUser.Email,
		IsPublished: updatedUser.IsPublished,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
	}

	return ctx.Status(fiber.StatusOK).JSON(responseUser)
}

func (h *UserHandler) UserVerifyHandler(ctx *fiber.Ctx) error {
	// JWT 토큰에서 사용자 ID 추출
	userID, err := h.AuthHandler.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("JWT 토큰을 통한 사용자 인증에 실패했습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	user, err := h.userUseCase.GetByID(userID)
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
