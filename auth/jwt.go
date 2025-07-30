package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dev-hyunsang/home-library/config"
	"github.com/dev-hyunsang/home-library/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWT 관련 오류
var (
	ErrInvalidToken = errors.New("토큰이 유효하지 않습니다")
	ErrExpiredToken = errors.New("토큰이 만료되었습니다")
)

// JWTClaims는 JWT 토큰에 포함될 클레임(데이터)을 정의합니다.
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateAccessToken은 액세스 토큰을 생성합니다.
func GenerateAccessToken(userID uuid.UUID, email string) (string, error) {
	claims := JWTClaims{
		UserID: userID.String(),
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.JWTExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", err
	}

	// Redis에 토큰 저장
	ctx := context.Background()
	err = database.SetJWTToken(ctx, userID.String(), "access", tokenString, config.JWTExpiration)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken은 리프레시 토큰을 생성합니다.
func GenerateRefreshToken(userID uuid.UUID) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.RefreshDuration)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", err
	}

	// Redis에 토큰 저장
	ctx := context.Background()
	err = database.SetJWTToken(ctx, userID.String(), "refresh", tokenString, config.RefreshDuration)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyToken은 JWT 토큰을 검증합니다.
func VerifyToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 서명 방식이 HMAC인지 확인
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("예상치 못한 서명 방식: %v", token.Header["alg"])
		}
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Redis에서 토큰 확인 (토큰 무효화 체크)
	ctx := context.Background()
	storedToken, err := database.GetJWTToken(ctx, claims.UserID, "access")
	if err != nil || storedToken != tokenString {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// RefreshAccessToken은 리프레시 토큰을 사용하여 액세스 토큰을 갱신합니다.
func RefreshAccessToken(refreshToken string) (string, error) {
	// 리프레시 토큰 검증
	token, err := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("예상치 못한 서명 방식: %v", token.Header["alg"])
		}
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", ErrExpiredToken
		}
		return "", ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return "", ErrInvalidToken
	}

	// Redis에서 리프레시 토큰 확인
	ctx := context.Background()
	storedToken, err := database.GetJWTToken(ctx, claims.Subject, "refresh")
	if err != nil || storedToken != refreshToken {
		return "", ErrInvalidToken
	}

	// 사용자 정보 조회 (이메일 필드를 위한)
	var user database.User
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return "", err
	}

	if result := database.DB.Where("id = ?", userID).First(&user); result.Error != nil {
		return "", result.Error
	}

	// 새 액세스 토큰 생성
	return GenerateAccessToken(userID, user.Email)
}

// InvalidateTokens는 사용자의 모든 토큰을 무효화합니다.
func InvalidateTokens(userID uuid.UUID) error {
	ctx := context.Background()
	return database.DeleteAllUserTokens(ctx, userID.String())
}
