package handler

import (
	"github.com/dev-hyunsang/home-library-backend/internal/domain"
	"github.com/dev-hyunsang/home-library-backend/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AdminHandler struct {
	userRepo      domain.UserRepository
	apiKeyUseCase domain.AdminAPIKeyUseCase
}

func NewAdminHandler(userRepo domain.UserRepository, apiKeyUseCase domain.AdminAPIKeyUseCase) *AdminHandler {
	return &AdminHandler{
		userRepo:      userRepo,
		apiKeyUseCase: apiKeyUseCase,
	}
}

type BroadcastNotificationRequest struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type BroadcastNotificationResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	TotalUsers  int    `json:"total_users"`
	SentCount   int    `json:"sent_count"`
	FailedCount int    `json:"failed_count"`
}

func (h *AdminHandler) BroadcastNotificationHandler(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"success": false,
		"message": "일괄 알림 기능은 현재 지원되지 않습니다.",
	})
}

func (h *AdminHandler) CreateAPIKeyHandler(ctx *fiber.Ctx) error {
	req := new(domain.CreateAPIKeyRequest)
	if err := ctx.BodyParser(req); err != nil {
		logger.Sugar().Errorf("요청 바디 파싱 실패: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	if req.Name == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	result, err := h.apiKeyUseCase.CreateAPIKey(req)
	if err != nil {
		logger.Sugar().Errorf("API Key 생성 실패: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
	}

	logger.Sugar().Infof("새 Admin API Key가 생성되었습니다. Name: %s, Prefix: %s", result.Name, result.KeyPrefix)

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "API Key가 생성되었습니다. 이 키는 다시 표시되지 않으니 안전하게 보관하세요.",
		"api_key": result,
	})
}

func (h *AdminHandler) GetAPIKeysHandler(ctx *fiber.Ctx) error {
	keys, err := h.apiKeyUseCase.GetAllAPIKeys()
	if err != nil {
		logger.Sugar().Errorf("API Key 목록 조회 실패: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":  true,
		"api_keys": keys,
		"count":    len(keys),
	})
}

func (h *AdminHandler) DeactivateAPIKeyHandler(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	if idParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	err = h.apiKeyUseCase.DeactivateAPIKey(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(ErrorHandler(domain.ErrNotFound))
		}
		logger.Sugar().Errorf("API Key 비활성화 실패: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "API Key가 비활성화되었습니다.",
	})
}

func (h *AdminHandler) DeleteAPIKeyHandler(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	if idParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	err = h.apiKeyUseCase.DeleteAPIKey(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(ErrorHandler(domain.ErrNotFound))
		}
		logger.Sugar().Errorf("API Key 삭제 실패: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
	}

	return ctx.Status(fiber.StatusNoContent).JSON(nil)
}
