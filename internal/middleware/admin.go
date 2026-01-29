package middleware

import (
	"github.com/dev-hyunsang/home-library/internal/domain"
	"github.com/dev-hyunsang/home-library/internal/handler"
	"github.com/dev-hyunsang/home-library/logger"
	"github.com/gofiber/fiber/v2"
)

func AdminAPIKeyMiddleware(apiKeyUseCase domain.AdminAPIKeyUseCase) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		providedKey := ctx.Get("X-Admin-API-Key")
		if providedKey == "" {
			providedKey = ctx.Query("api_key")
		}

		if providedKey == "" {
			logger.Init().Sugar().Warn("Admin API Key가 제공되지 않았습니다")
			return ctx.Status(fiber.StatusUnauthorized).JSON(handler.ErrorHandler(domain.ErrInvalidToken))
		}

		apiKey, err := apiKeyUseCase.ValidateAPIKey(providedKey)
		if err != nil {
			if err == domain.ErrTokenExpired {
				logger.Init().Sugar().Warn("만료된 Admin API Key")
				return ctx.Status(fiber.StatusUnauthorized).JSON(handler.ErrorHandler(domain.ErrTokenExpired))
			}
			logger.Init().Sugar().Warn("유효하지 않은 Admin API Key")
			return ctx.Status(fiber.StatusForbidden).JSON(handler.ErrorHandler(domain.ErrPermissionDenied))
		}

		ctx.Locals("adminAPIKeyID", apiKey.ID.String())
		ctx.Locals("adminAPIKeyName", apiKey.Name)

		logger.Init().Sugar().Infof("Admin API Key 인증 성공: %s (%s)", apiKey.Name, apiKey.KeyPrefix)
		return ctx.Next()
	}
}

func AdminBootstrapMiddleware(bootstrapKey string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if bootstrapKey == "" {
			logger.Init().Sugar().Error("ADMIN_BOOTSTRAP_KEY가 설정되지 않았습니다")
			return ctx.Status(fiber.StatusInternalServerError).JSON(handler.ErrorHandler(domain.ErrInternal))
		}

		providedKey := ctx.Get("X-Admin-Bootstrap-Key")
		if providedKey == "" {
			providedKey = ctx.Query("bootstrap_key")
		}

		if providedKey == "" {
			logger.Init().Sugar().Warn("Admin Bootstrap Key가 제공되지 않았습니다")
			return ctx.Status(fiber.StatusUnauthorized).JSON(handler.ErrorHandler(domain.ErrInvalidToken))
		}

		if providedKey != bootstrapKey {
			logger.Init().Sugar().Warn("유효하지 않은 Admin Bootstrap Key")
			return ctx.Status(fiber.StatusForbidden).JSON(handler.ErrorHandler(domain.ErrPermissionDenied))
		}

		logger.Init().Sugar().Info("Admin Bootstrap Key 인증 성공")
		return ctx.Next()
	}
}
