package response

import "time"

// ErrorResponse is a standard error response structure
type ErrorResponse struct {
	IsSuccess bool   `json:"is_success"`
	Message   string `json:"message"`
	Time      string `json:"time"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		IsSuccess: false,
		Message:   err.Error(),
		Time:      time.Now().Format(time.RFC3339),
	}
}

// TokenResponse is a standard token response structure
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// NewTokenResponse creates a new token response
func NewTokenResponse(accessToken, refreshToken string, expiresIn int) TokenResponse {
	return TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}
}

// SuccessResponse is a generic success response
type SuccessResponse struct {
	IsSuccess bool        `json:"is_success"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
}

// NewSuccessResponse creates a new success response with data
func NewSuccessResponse(data interface{}) SuccessResponse {
	return SuccessResponse{
		IsSuccess: true,
		Data:      data,
	}
}

// NewSuccessMessageResponse creates a new success response with message
func NewSuccessMessageResponse(message string) SuccessResponse {
	return SuccessResponse{
		IsSuccess: true,
		Message:   message,
	}
}
