package handler

import (
	"github.com/dev-hyunsang/home-library/internal/config"
	"github.com/dev-hyunsang/home-library/internal/domain"
)

type AuthHandler struct {
	authUseCase domain.AuthUseCase
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(authUseCase *domain.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: *authUseCase,
	}

}

type TokenService struct {
	JwtConfig *config.JwtConfig
}

func NewTokenService(cfg *config.JwtConfig) *TokenService {
	return &TokenService{
		JwtConfig: cfg,
	}
}
