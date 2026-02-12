package request

// RegistrationRequest is the request body for user registration
type RegistrationRequest struct {
	NickName      string `json:"nick_name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	IsPublished   bool   `json:"is_published"`
	IsTermsAgreed bool   `json:"is_terms_agreed"`
}

// UpdateUserRequest is the request body for updating user information
type UpdateUserRequest struct {
	NickName    string  `json:"nick_name"`
	Email       string  `json:"email"`
	Password    *string `json:"password,omitempty"`
	IsPublished bool    `json:"is_published"`
}

// LoginRequest is the request body for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ForgotPasswordRequest is the request body for password reset
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// VerifyCodeRequest is the request body for email verification
type VerifyCodeRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// UpdateFCMTokenRequest is the request body for FCM token update
type UpdateFCMTokenRequest struct {
	FCMToken string `json:"fcm_token"`
}

// UpdateTimezoneRequest is the request body for timezone update
type UpdateTimezoneRequest struct {
	Timezone string `json:"timezone"`
}
