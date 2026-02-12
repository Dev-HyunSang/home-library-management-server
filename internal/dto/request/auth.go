package request

// RefreshTokenRequest is the request body for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
