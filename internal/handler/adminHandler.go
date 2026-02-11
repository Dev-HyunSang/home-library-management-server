package handler

import (
	"context"

	"github.com/dev-hyunsang/home-library-backend/internal/domain"
	"github.com/dev-hyunsang/home-library-backend/internal/infrastructure/kafka"
	"github.com/dev-hyunsang/home-library-backend/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AdminHandler struct {
	userRepo      domain.UserRepository
	kafkaProducer *kafka.Producer
	apiKeyUseCase domain.AdminAPIKeyUseCase
}

func NewAdminHandler(userRepo domain.UserRepository, kafkaProducer *kafka.Producer, apiKeyUseCase domain.AdminAPIKeyUseCase) *AdminHandler {
	return &AdminHandler{
		userRepo:      userRepo,
		kafkaProducer: kafkaProducer,
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
	req := new(BroadcastNotificationRequest)
	if err := ctx.BodyParser(req); err != nil {
		logger.Init().Sugar().Errorf("요청 바디 파싱 실패: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	if req.Title == "" || req.Message == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	users, err := h.userRepo.GetAllUsersWithFCM()
	if err != nil {
		logger.Init().Sugar().Errorf("사용자 목록 조회 실패: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
	}

	if len(users) == 0 {
		return ctx.Status(fiber.StatusOK).JSON(BroadcastNotificationResponse{
			Success:     true,
			Message:     "FCM 토큰이 등록된 사용자가 없습니다.",
			TotalUsers:  0,
			SentCount:   0,
			FailedCount: 0,
		})
	}

	sentCount := 0
	failedCount := 0

	for _, user := range users {
		err := h.kafkaProducer.ProduceNotification(
			context.Background(),
			user.ID.String(),
			req.Title,
			req.Message,
			"admin_broadcast",
		)
		if err != nil {
			logger.Init().Sugar().Errorf("사용자 %s에게 알림 발송 실패: %v", user.ID.String(), err)
			failedCount++
			continue
		}
		sentCount++
	}

	logger.Init().Sugar().Infof("관리자 일괄 알림 발송 완료 - 전체: %d, 성공: %d, 실패: %d",
		len(users), sentCount, failedCount)

	return ctx.Status(fiber.StatusOK).JSON(BroadcastNotificationResponse{
		Success:     true,
		Message:     "알림 발송이 완료되었습니다.",
		TotalUsers:  len(users),
		SentCount:   sentCount,
		FailedCount: failedCount,
	})
}

func (h *AdminHandler) CreateAPIKeyHandler(ctx *fiber.Ctx) error {
	req := new(domain.CreateAPIKeyRequest)
	if err := ctx.BodyParser(req); err != nil {
		logger.Init().Sugar().Errorf("요청 바디 파싱 실패: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	if req.Name == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	result, err := h.apiKeyUseCase.CreateAPIKey(req)
	if err != nil {
		logger.Init().Sugar().Errorf("API Key 생성 실패: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
	}

	logger.Init().Sugar().Infof("새 Admin API Key가 생성되었습니다. Name: %s, Prefix: %s", result.Name, result.KeyPrefix)

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "API Key가 생성되었습니다. 이 키는 다시 표시되지 않으니 안전하게 보관하세요.",
		"api_key": result,
	})
}

func (h *AdminHandler) GetAPIKeysHandler(ctx *fiber.Ctx) error {
	keys, err := h.apiKeyUseCase.GetAllAPIKeys()
	if err != nil {
		logger.Init().Sugar().Errorf("API Key 목록 조회 실패: %v", err)
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
		logger.Init().Sugar().Errorf("API Key 비활성화 실패: %v", err)
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
		logger.Init().Sugar().Errorf("API Key 삭제 실패: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
	}

	return ctx.Status(fiber.StatusNoContent).JSON(nil)
}
