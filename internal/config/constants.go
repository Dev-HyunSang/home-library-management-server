package config

import "time"

// Token configuration
const (
	AccessTokenExpirySeconds  = 3600          // 1 hour in seconds
	RefreshTokenExpirySeconds = 86400         // 24 hours in seconds
	AccessTokenExpiry         = 1 * time.Hour
	RefreshTokenExpiry        = 24 * time.Hour
)

// Rate limiting configuration
const (
	DefaultRateLimitRequests = 1000
	DefaultRateLimitWindow   = 1 * time.Hour
)

// Email configuration
const (
	DefaultSMTPHost = "smtp.gmail.com"
	DefaultSMTPPort = "587"
)

// Verification configuration
const (
	VerificationCodeExpiry = 5 * time.Minute
	TempPasswordLength     = 15
)
