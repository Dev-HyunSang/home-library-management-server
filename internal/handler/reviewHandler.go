package handler

import (
	"time"

	"github.com/dev-hyunsang/my-own-library-backend/internal/domain"
	"github.com/dev-hyunsang/my-own-library-backend/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReviewHandler struct {
	reviewUseCase domain.ReviewUseCase
	authUseCase   domain.AuthUseCase
}

func NewReviewHandler(reviewUseCase domain.ReviewUseCase, authUseCase domain.AuthUseCase) *ReviewHandler {
	return &ReviewHandler{
		reviewUseCase: reviewUseCase,
		authUseCase:   authUseCase,
	}
}

// POST /api/reviews/isbn/:isbn
func (h *ReviewHandler) CreateReviewHandler(ctx *fiber.Ctx) error {
	isbn := ctx.Params("isbn")
	if isbn == "" {
		logger.Sugar().Error("ISBN이 입력되지 않았습니다.")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"is_success": false,
			"message":    "ISBN은 필수입니다.",
			"time":       time.Now().String(),
		})
	}

	userID, err := h.authUseCase.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Sugar().Errorf("JWT 토큰 인증 실패: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"is_success": false,
			"message":    "인증이 필요합니다.",
			"time":       time.Now().String(),
		})
	}

	req := new(domain.CreateReviewRequest)
	if err := ctx.BodyParser(req); err != nil {
		logger.Sugar().Errorf("요청 파싱 실패: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"is_success": false,
			"message":    "올바르지 않은 요청입니다.",
			"time":       time.Now().String(),
		})
	}

	review, err := h.reviewUseCase.CreateReview(userID, isbn, req)
	if err != nil {
		logger.Sugar().Errorf("리뷰 생성 실패: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"is_success": false,
			"message":    err.Error(),
			"time":       time.Now().String(),
		})
	}

	logger.Sugar().Infof("리뷰가 생성되었습니다. ID: %s, ISBN: %s", review.ID.String(), isbn)

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"is_success": true,
		"message":    "리뷰가 생성되었습니다.",
		"data":       review,
	})
}

// GET /api/reviews/isbn/:isbn
func (h *ReviewHandler) GetReviewsByISBNHandler(ctx *fiber.Ctx) error {
	isbn := ctx.Params("isbn")
	if isbn == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"is_success": false,
			"message":    "ISBN은 필수입니다.",
			"time":       time.Now().String(),
		})
	}

	reviews, err := h.reviewUseCase.GetReviewsByISBN(isbn)
	if err != nil {
		logger.Sugar().Errorf("리뷰 조회 실패: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"is_success": false,
			"message":    "리뷰 조회 중 오류가 발생했습니다.",
			"time":       time.Now().String(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_success": true,
		"data":       reviews,
		"count":      len(reviews),
	})
}

// GET /api/reviews/isbn/:isbn/:id
func (h *ReviewHandler) GetReviewByIDHandler(ctx *fiber.Ctx) error {
	reviewID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"is_success": false,
			"message":    "올바르지 않은 리뷰 ID입니다.",
			"time":       time.Now().String(),
		})
	}

	review, err := h.reviewUseCase.GetReviewByID(reviewID)
	if err != nil {
		logger.Sugar().Errorf("리뷰 조회 실패: %v", err)
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"is_success": false,
			"message":    "리뷰를 찾을 수 없습니다.",
			"time":       time.Now().String(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_success": true,
		"data":       review,
	})
}

// PUT /api/reviews/isbn/:isbn/:id
func (h *ReviewHandler) UpdateReviewHandler(ctx *fiber.Ctx) error {
	reviewID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"is_success": false,
			"message":    "올바르지 않은 리뷰 ID입니다.",
			"time":       time.Now().String(),
		})
	}

	userID, err := h.authUseCase.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Sugar().Errorf("JWT 토큰 인증 실패: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"is_success": false,
			"message":    "인증이 필요합니다.",
			"time":       time.Now().String(),
		})
	}

	req := new(domain.UpdateReviewRequest)
	if err := ctx.BodyParser(req); err != nil {
		logger.Sugar().Errorf("요청 파싱 실패: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"is_success": false,
			"message":    "올바르지 않은 요청입니다.",
			"time":       time.Now().String(),
		})
	}

	review, err := h.reviewUseCase.UpdateReview(userID, reviewID, req)
	if err != nil {
		logger.Sugar().Errorf("리뷰 수정 실패: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"is_success": false,
			"message":    err.Error(),
			"time":       time.Now().String(),
		})
	}

	logger.Sugar().Infof("리뷰가 수정되었습니다. ID: %s", review.ID.String())

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_success": true,
		"message":    "리뷰가 수정되었습니다.",
		"data":       review,
	})
}

// DELETE /api/reviews/isbn/:isbn/:id
func (h *ReviewHandler) DeleteReviewHandler(ctx *fiber.Ctx) error {
	reviewID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"is_success": false,
			"message":    "올바르지 않은 리뷰 ID입니다.",
			"time":       time.Now().String(),
		})
	}

	userID, err := h.authUseCase.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Sugar().Errorf("JWT 토큰 인증 실패: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"is_success": false,
			"message":    "인증이 필요합니다.",
			"time":       time.Now().String(),
		})
	}

	err = h.reviewUseCase.DeleteReview(userID, reviewID)
	if err != nil {
		logger.Sugar().Errorf("리뷰 삭제 실패: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"is_success": false,
			"message":    err.Error(),
			"time":       time.Now().String(),
		})
	}

	logger.Sugar().Infof("리뷰가 삭제되었습니다. ID: %s", reviewID.String())

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_success": true,
		"message":    "리뷰가 삭제되었습니다.",
	})
}

// GET /api/reviews/me
func (h *ReviewHandler) GetMyReviewsHandler(ctx *fiber.Ctx) error {
	userID, err := h.authUseCase.GetUserIDFromToken(ctx)
	if err != nil {
		logger.Sugar().Errorf("JWT 토큰 인증 실패: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"is_success": false,
			"message":    "인증이 필요합니다.",
			"time":       time.Now().String(),
		})
	}

	reviews, err := h.reviewUseCase.GetUserReviews(userID)
	if err != nil {
		logger.Sugar().Errorf("리뷰 조회 실패: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"is_success": false,
			"message":    "리뷰 조회 중 오류가 발생했습니다.",
			"time":       time.Now().String(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_success": true,
		"data":       reviews,
		"count":      len(reviews),
	})
}
