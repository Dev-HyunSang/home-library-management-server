package middleware

import (
	"time"

	"github.com/dev-hyunsang/home-library-backend/internal/domain"
	"github.com/dev-hyunsang/home-library-backend/internal/handler"
	"github.com/dev-hyunsang/home-library-backend/logger"
	"github.com/gofiber/fiber/v2"
)

func JWTAuthMiddleware(authUseCase domain.AuthUseCase) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// JWT 토큰에서 사용자 ID 추출
		userID, err := authUseCase.GetUserIDFromToken(ctx)
		if err != nil {
			logger.Init().Sugar().Errorf("JWT 토큰 인증에 실패했습니다: %v", err)
			return ctx.Status(fiber.StatusUnauthorized).JSON(handler.ErrorHandler(domain.ErrUserNotLoggedIn))
		}

		// Rate limiting 체크 (API 호출 제한)
		allowed, err := authUseCase.CheckRateLimit(userID, "api_call", 1000, 1*time.Hour)
		if err != nil {
			logger.Init().Sugar().Errorf("Rate limit 체크 중 오류가 발생했습니다: %v", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(handler.ErrorHandler(domain.ErrInternal))
		}
		if !allowed {
			logger.Init().Sugar().Warnf("Rate limit 초과: 사용자ID %s", userID.String())
			return ctx.Status(fiber.StatusTooManyRequests).JSON(handler.ErrorHandler(domain.ErrTooManyRequests))
		}

		// 컨텍스트에 사용자 ID 저장
		ctx.Locals("userID", userID.String())

		logger.Init().Sugar().Infof("JWT 미들웨어를 통해 사용자 인증이 완료되었습니다 / 사용자ID: %s", userID.String())

		return ctx.Next()
	}
}

func LoginRateLimitMiddleware(authUseCase domain.AuthUseCase) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// IP 기반 로그인 시도 제한
		clientIP := ctx.IP()

		// 임시로 고정 UUID 사용 (실제로는 IP를 UUID로 변환하거나 다른 방식 사용)
		// 여기서는 간단히 건너뛰고 다음 핸들러로 진행
		logger.Init().Sugar().Infof("로그인 요청 - IP: %s", clientIP)

		return ctx.Next()
	}
}
