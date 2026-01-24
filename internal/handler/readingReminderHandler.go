package handler

import (
	"github.com/dev-hyunsang/home-library/internal/domain"
	"github.com/dev-hyunsang/home-library/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReadingReminderHandler struct {
	reminderUseCase domain.ReadingReminderUseCase
	authUseCase     domain.AuthUseCase
}

func NewReadingReminderHandler(reminderUseCase domain.ReadingReminderUseCase, authUseCase domain.AuthUseCase) *ReadingReminderHandler {
	return &ReadingReminderHandler{
		reminderUseCase: reminderUseCase,
		authUseCase:     authUseCase,
	}
}

func (h *ReadingReminderHandler) CreateReminderHandler(ctx *fiber.Ctx) error {
	userID, err := h.authUseCase.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("JWT 토큰에서 사용자 ID 추출 실패: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	req := new(domain.CreateReminderRequest)
	if err := ctx.BodyParser(req); err != nil {
		logger.Init().Sugar().Errorf("요청 바디 파싱 실패: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	if req.ReminderTime == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidReminderTime))
	}

	reminder, err := h.reminderUseCase.CreateReminder(userID, req)
	if err != nil {
		logger.Init().Sugar().Errorf("알림 생성 실패: %v", err)
		if err == domain.ErrInvalidReminderTime || err == domain.ErrInvalidDayOfWeek {
			return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(err))
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
	}

	logger.Init().Sugar().Infof("알림이 생성되었습니다. 알림ID: %s, 사용자ID: %s", reminder.ID.String(), userID.String())
	return ctx.Status(fiber.StatusCreated).JSON(reminder)
}

func (h *ReadingReminderHandler) GetRemindersHandler(ctx *fiber.Ctx) error {
	userID, err := h.authUseCase.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("JWT 토큰에서 사용자 ID 추출 실패: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	reminders, err := h.reminderUseCase.GetUserReminders(userID)
	if err != nil {
		logger.Init().Sugar().Errorf("알림 목록 조회 실패: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
	}

	logger.Init().Sugar().Infof("알림 목록을 조회했습니다. 사용자ID: %s, 알림 수: %d", userID.String(), len(reminders))
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"reminders": reminders,
		"count":     len(reminders),
	})
}

func (h *ReadingReminderHandler) UpdateReminderHandler(ctx *fiber.Ctx) error {
	userID, err := h.authUseCase.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("JWT 토큰에서 사용자 ID 추출 실패: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	idParam := ctx.Params("id")
	if idParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	reminderID, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	req := new(domain.UpdateReminderRequest)
	if err := ctx.BodyParser(req); err != nil {
		logger.Init().Sugar().Errorf("요청 바디 파싱 실패: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	reminder, err := h.reminderUseCase.UpdateReminder(reminderID, userID, req)
	if err != nil {
		logger.Init().Sugar().Errorf("알림 수정 실패: %v", err)
		switch err {
		case domain.ErrReminderNotFound:
			return ctx.Status(fiber.StatusNotFound).JSON(ErrorHandler(err))
		case domain.ErrReminderOwnerMismatch:
			return ctx.Status(fiber.StatusForbidden).JSON(ErrorHandler(domain.ErrPermissionDenied))
		case domain.ErrInvalidReminderTime, domain.ErrInvalidDayOfWeek:
			return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(err))
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
		}
	}

	logger.Init().Sugar().Infof("알림이 수정되었습니다. 알림ID: %s", reminder.ID.String())
	return ctx.Status(fiber.StatusOK).JSON(reminder)
}

func (h *ReadingReminderHandler) ToggleReminderHandler(ctx *fiber.Ctx) error {
	userID, err := h.authUseCase.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("JWT 토큰에서 사용자 ID 추출 실패: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	idParam := ctx.Params("id")
	if idParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	reminderID, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	reminder, err := h.reminderUseCase.ToggleReminder(reminderID, userID)
	if err != nil {
		logger.Init().Sugar().Errorf("알림 토글 실패: %v", err)
		switch err {
		case domain.ErrReminderNotFound:
			return ctx.Status(fiber.StatusNotFound).JSON(ErrorHandler(err))
		case domain.ErrReminderOwnerMismatch:
			return ctx.Status(fiber.StatusForbidden).JSON(ErrorHandler(domain.ErrPermissionDenied))
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
		}
	}

	logger.Init().Sugar().Infof("알림이 토글되었습니다. 알림ID: %s, 활성화: %v", reminder.ID.String(), reminder.IsEnabled)
	return ctx.Status(fiber.StatusOK).JSON(reminder)
}

func (h *ReadingReminderHandler) DeleteReminderHandler(ctx *fiber.Ctx) error {
	userID, err := h.authUseCase.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Init().Sugar().Errorf("JWT 토큰에서 사용자 ID 추출 실패: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorHandler(domain.ErrUserNotLoggedIn))
	}

	idParam := ctx.Params("id")
	if idParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	reminderID, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorHandler(domain.ErrInvalidInput))
	}

	err = h.reminderUseCase.DeleteReminder(reminderID, userID)
	if err != nil {
		logger.Init().Sugar().Errorf("알림 삭제 실패: %v", err)
		switch err {
		case domain.ErrReminderNotFound:
			return ctx.Status(fiber.StatusNotFound).JSON(ErrorHandler(err))
		case domain.ErrReminderOwnerMismatch:
			return ctx.Status(fiber.StatusForbidden).JSON(ErrorHandler(domain.ErrPermissionDenied))
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(ErrorHandler(domain.ErrInternal))
		}
	}

	logger.Init().Sugar().Infof("알림이 삭제되었습니다. 알림ID: %s", reminderID.String())
	return ctx.Status(fiber.StatusNoContent).JSON(nil)
}
