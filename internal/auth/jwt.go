package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID       string `json:"user_id"`
	TokenID      string `json:"token_id"`
	TokenVersion int64  `json:"token_version"`
	TokenType    string `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

type JWTManager struct {
	secretKey         string
	accessTokenTTL    time.Duration
	refreshTokenTTL   time.Duration
	securityManager   *SecurityManager
}

func NewJWTManager(secretKey string, accessTTL, refreshTTL time.Duration, securityManager *SecurityManager) *JWTManager {
	return &JWTManager{
		secretKey:       secretKey,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
		securityManager: securityManager,
	}
}

func (manager *JWTManager) GenerateTokenPair(userID uuid.UUID) (accessToken, refreshToken string, err error) {
	// 사용자 토큰 버전 조회
	tokenVersion, err := manager.securityManager.GetUserTokenVersion(userID)
	if err != nil {
		return "", "", err
	}

	now := time.Now()
	accessTokenID := uuid.New().String()
	refreshTokenID := uuid.New().String()

	// Access Token 생성
	accessClaims := &JWTClaims{
		UserID:       userID.String(),
		TokenID:      accessTokenID,
		TokenVersion: tokenVersion,
		TokenType:    "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(manager.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString([]byte(manager.secretKey))
	if err != nil {
		return "", "", err
	}

	// Refresh Token 생성
	refreshClaims := &JWTClaims{
		UserID:       userID.String(),
		TokenID:      refreshTokenID,
		TokenVersion: tokenVersion,
		TokenType:    "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(manager.refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString([]byte(manager.secretKey))
	if err != nil {
		return "", "", err
	}

	// Redis에 세션 토큰과 리프레시 토큰 저장
	if err := manager.securityManager.StoreSessionToken(userID, accessTokenID, manager.accessTokenTTL); err != nil {
		return "", "", err
	}

	if err := manager.securityManager.StoreRefreshToken(userID, refreshTokenID, manager.refreshTokenTTL); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (manager *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(manager.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// 블랙리스트 확인
	isBlacklisted, err := manager.securityManager.IsTokenBlacklisted(claims.TokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to check token blacklist: %w", err)
	}
	if isBlacklisted {
		return nil, fmt.Errorf("token is blacklisted")
	}

	// 토큰 버전 확인
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	currentVersion, err := manager.securityManager.GetUserTokenVersion(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user token version: %w", err)
	}

	if claims.TokenVersion < currentVersion {
		return nil, fmt.Errorf("token version is outdated")
	}

	// 세션 토큰 유효성 확인 (Access Token의 경우)
	if claims.TokenType == "access" {
		_, err := manager.securityManager.ValidateSessionToken(claims.TokenID)
		if err != nil {
			return nil, fmt.Errorf("session token validation failed: %w", err)
		}
	}

	return claims, nil
}

func (manager *JWTManager) ExtractUserID(tokenString string) (uuid.UUID, error) {
	claims, err := manager.ValidateToken(tokenString)
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in token: %w", err)
	}

	return userID, nil
}

func (manager *JWTManager) RefreshToken(refreshTokenString string) (newAccessToken, newRefreshToken string, err error) {
	claims, err := manager.ValidateToken(refreshTokenString)
	if err != nil {
		return "", "", err
	}

	if claims.TokenType != "refresh" {
		return "", "", fmt.Errorf("invalid token type for refresh")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return "", "", fmt.Errorf("invalid user ID in refresh token")
	}

	// 리프레시 토큰 유효성 확인
	isValid, err := manager.securityManager.IsRefreshTokenValid(userID, claims.TokenID)
	if err != nil {
		return "", "", err
	}
	if !isValid {
		return "", "", fmt.Errorf("refresh token is not valid")
	}

	// 기존 리프레시 토큰 무효화
	if err := manager.securityManager.RevokeRefreshToken(userID, claims.TokenID); err != nil {
		return "", "", err
	}

	// 새 토큰 쌍 생성
	return manager.GenerateTokenPair(userID)
}

func (manager *JWTManager) InvalidateToken(tokenString string) error {
	claims, err := manager.ValidateToken(tokenString)
	if err != nil {
		return err
	}

	// 토큰을 블랙리스트에 추가
	expiration := time.Until(claims.ExpiresAt.Time)
	if expiration > 0 {
		return manager.securityManager.BlacklistToken(claims.TokenID, expiration)
	}

	return nil
}

func (manager *JWTManager) InvalidateAllUserTokens(userID uuid.UUID) error {
	return manager.securityManager.InvalidateAllUserTokens(userID)
}