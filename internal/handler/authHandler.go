package handler

import (
	"time"

	"github.com/dev-hyunsang/my-own-library-backend/internal/config"
	"github.com/dev-hyunsang/my-own-library-backend/internal/domain"
	"github.com/dev-hyunsang/my-own-library-backend/logger"
	"github.com/gofiber/fiber/v2"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthHandler struct {
	authUseCase domain.AuthUseCase
}

func NewAuthHandler(authUseCase domain.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

func (h *AuthHandler) RefreshTokenHandler(ctx *fiber.Ctx) error {
	req := new(RefreshTokenRequest)
	if err := ctx.BodyParser(req); err != nil {
		logger.Sugar().Errorf("요청 본문을 파싱하는 도중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	if req.RefreshToken == "" {
		logger.Sugar().Error("리프레시 토큰이 제공되지 않았습니다.")
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	newAccessToken, newRefreshToken, err := h.authUseCase.RefreshToken(req.RefreshToken)
	if err != nil {
		logger.Sugar().Errorf("토큰 갱신 중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrInvalidToken))
	}

	logger.Sugar().Info("토큰이 성공적으로 갱신되었습니다.")

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
		"token_type":    "Bearer",
		"expires_in":    config.AccessTokenExpirySeconds,
	})
}

func (h *AuthHandler) RevokeAllTokensHandler(ctx *fiber.Ctx) error {
	// JWT 토큰에서 사용자 ID 추출
	userID, err := h.authUseCase.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Sugar().Errorf("JWT 토큰을 통한 사용자 인증에 실패했습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	// 사용자의 모든 토큰 무효화
	if err := h.authUseCase.InvalidateAllUserTokens(userID); err != nil {
		logger.Sugar().Errorf("모든 토큰 무효화에 실패했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
	}

	logger.Sugar().Infof("사용자의 모든 토큰이 무효화되었습니다 / 사용자ID: %s", userID.String())

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "모든 기기에서 로그아웃되었습니다.",
	})
}

func (h *AuthHandler) CheckRateLimitHandler(ctx *fiber.Ctx) error {
	// JWT 토큰에서 사용자 ID 추출
	userID, err := h.authUseCase.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Sugar().Errorf("JWT 토큰을 통한 사용자 인증에 실패했습니다: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	// Rate limit 체크
	allowed, err := h.authUseCase.CheckRateLimit(userID, "api_call", config.DefaultRateLimitRequests, config.DefaultRateLimitWindow)
	if err != nil {
		logger.Sugar().Errorf("Rate limit 체크 중 오류가 발생했습니다: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"allowed":   allowed,
		"user_id":   userID.String(),
		"action":    "api_call",
		"limit":     config.DefaultRateLimitRequests,
		"window":    config.DefaultRateLimitWindow.String(),
		"timestamp": time.Now(),
	})
}