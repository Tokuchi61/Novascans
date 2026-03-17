package app

import (
	"time"

	"github.com/google/uuid"

	"github.com/Tokuchi61/Novascans/internal/modules/identity/auth/domain"
)

type PingData struct {
	Module      string `json:"module"`
	Status      string `json:"status"`
	GeneratedAt string `json:"generated_at"`
}

type RegisterInput struct {
	Email     string
	Password  string
	UserAgent string
	IPAddress string
}

type LoginInput struct {
	Email     string
	Password  string
	UserAgent string
	IPAddress string
}

type RefreshSessionInput struct {
	RefreshToken string
	UserAgent    string
	IPAddress    string
}

type RequestEmailVerificationInput struct {
	Email string
}

type VerifyEmailInput struct {
	Token string
}

type ForgotPasswordInput struct {
	Email string
}

type ResetPasswordInput struct {
	Token       string
	NewPassword string
}

type AuthResult struct {
	User         domain.User
	Session      domain.Session
	AccessToken  string
	RefreshToken string
}

type AuthenticatedUser struct {
	User    domain.User
	Session domain.Session
}

type RequestTokenResult struct {
	DebugToken string
}

type SessionEventPayload struct {
	UserID uuid.UUID
}

func NewPingData(now time.Time) PingData {
	return PingData{
		Module:      "identity.auth",
		Status:      "ok",
		GeneratedAt: now.UTC().Format(time.RFC3339),
	}
}
