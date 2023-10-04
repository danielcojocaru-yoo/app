package models

import "time"

type User struct {
	Email        string `json:"email"`
	AvatarURL    string `json:"avatar_url"`
	FullName     string `json:"name"`
	ReferralCode string `json:"referral_code"`
}

type Response struct {
	Message string         `json:"message"`
	Error   string         `json:"error,omitempty"`
	Tokens  *TokenResponse `json:"tokens,omitempty"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiry       int64  `json:"expiry"`
}

type CodeWithExpiry struct {
	Code int
	Exp  time.Time
}
type CodeInfo struct {
	Code int
	Exp  time.Time
}

type codeExpiry struct {
	Code string
	Exp  time.Time
}

// Define global variables
var (
	expectedCode int
	loginCodes   = map[string]CodeWithExpiry{}
)

var emailCodes = make(map[string]codeExpiry)
