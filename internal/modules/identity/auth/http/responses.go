package http

import (
	"time"

	authapp "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/app"
	"github.com/Tokuchi61/Novascans/internal/modules/identity/auth/domain"
)

type authResponse struct {
	User         userResponse    `json:"user"`
	Session      sessionResponse `json:"session"`
	AccessToken  string          `json:"access_token,omitempty"`
	RefreshToken string          `json:"refresh_token,omitempty"`
}

type currentUserResponse struct {
	User    userResponse    `json:"user"`
	Session sessionResponse `json:"session"`
}

type userResponse struct {
	ID              string `json:"id"`
	Email           string `json:"email"`
	BaseRole        string `json:"base_role"`
	Status          string `json:"status"`
	EmailVerifiedAt string `json:"email_verified_at,omitempty"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type sessionResponse struct {
	ID                string `json:"id"`
	UserID            string `json:"user_id"`
	ExpiresAt         string `json:"expires_at"`
	CreatedAt         string `json:"created_at"`
	RevokedAt         string `json:"revoked_at,omitempty"`
	ReplacedBySession string `json:"replaced_by_session_id,omitempty"`
}

type tokenRequestResponse struct {
	DebugToken string `json:"debug_token,omitempty"`
}

type successResponse struct {
	Success bool `json:"success"`
}

func mapAuthResponse(result authapp.AuthResult) authResponse {
	return authResponse{
		User:         mapUserResponse(result.User),
		Session:      mapSessionResponse(result.Session),
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}
}

func mapCurrentUserResponse(result authapp.AuthenticatedUser) currentUserResponse {
	return currentUserResponse{
		User:    mapUserResponse(result.User),
		Session: mapSessionResponse(result.Session),
	}
}

func mapUserResponse(user domain.User) userResponse {
	response := userResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		BaseRole:  user.BaseRole,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.UTC().Format(time.RFC3339),
	}

	if user.EmailVerifiedAt != nil {
		response.EmailVerifiedAt = user.EmailVerifiedAt.UTC().Format(time.RFC3339)
	}

	return response
}

func mapSessionResponse(session domain.Session) sessionResponse {
	response := sessionResponse{
		ID:        session.ID.String(),
		UserID:    session.UserID.String(),
		ExpiresAt: session.ExpiresAt.UTC().Format(time.RFC3339),
		CreatedAt: session.CreatedAt.UTC().Format(time.RFC3339),
	}

	if session.RevokedAt != nil {
		response.RevokedAt = session.RevokedAt.UTC().Format(time.RFC3339)
	}

	if session.ReplacedBySession != nil {
		response.ReplacedBySession = session.ReplacedBySession.String()
	}

	return response
}

func mapTokenRequestResponse(result authapp.RequestTokenResult) tokenRequestResponse {
	return tokenRequestResponse{
		DebugToken: result.DebugToken,
	}
}
