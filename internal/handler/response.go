package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// ErrResponse is a standard error response structure
type ErrResponse struct {
	IsSuccess bool   `json:"is_success"`
	Message   string `json:"message"`
	Time      string `json:"time"`
}

// ErrorHandler creates a standardized error response
func ErrorHandler(err error) ErrResponse {
	return ErrResponse{
		IsSuccess: false,
		Message:   err.Error(),
		Time:      time.Now().Format(time.RFC3339),
	}
}

// SuccessResponse creates a standardized success response with data
func SuccessResponse(data interface{}) fiber.Map {
	return fiber.Map{
		"is_success": true,
		"data":       data,
	}
}

// SuccessMessageResponse creates a success response with a message
func SuccessMessageResponse(message string) fiber.Map {
	return fiber.Map{
		"is_success": true,
		"message":    message,
	}
}

// SuccessListResponse creates a success response with list data and count
func SuccessListResponse(data interface{}, count int) fiber.Map {
	return fiber.Map{
		"is_success": true,
		"data":       data,
		"count":      count,
	}
}
