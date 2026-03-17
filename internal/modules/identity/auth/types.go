package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/Tokuchi61/Novascans/internal/modules/identity/auth/store"
	"github.com/Tokuchi61/Novascans/internal/platform/validation"
)

type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type createSessionRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID              string `json:"id"`
	Email           string `json:"email"`
	Status          string `json:"status"`
	EmailVerifiedAt string `json:"email_verified_at,omitempty"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type sessionResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Token     string `json:"token,omitempty"`
	ExpiresAt string `json:"expires_at"`
	CreatedAt string `json:"created_at"`
	RevokedAt string `json:"revoked_at,omitempty"`
}

func (request createUserRequest) Validate(v *validation.Validator) validation.FieldErrors {
	errs := validation.FieldErrors{}
	v.Email("email", request.Email, errs)
	v.MinLength("password", request.Password, 8, errs)
	return errs
}

func (request createSessionRequest) Validate(v *validation.Validator) validation.FieldErrors {
	errs := validation.FieldErrors{}
	v.Email("email", request.Email, errs)
	v.RequiredString("password", request.Password, errs)
	return errs
}

func mapUserResponse(user store.User) userResponse {
	response := userResponse{
		ID:        user.ID,
		Email:     user.Email,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.UTC().Format(time.RFC3339),
	}

	if user.EmailVerifiedAt.Valid {
		response.EmailVerifiedAt = user.EmailVerifiedAt.Time.UTC().Format(time.RFC3339)
	}

	return response
}

func mapSessionResponse(session SessionOutput) sessionResponse {
	response := sessionResponse{
		ID:        session.ID,
		UserID:    session.UserID,
		Token:     session.Token,
		ExpiresAt: session.ExpiresAt.UTC().Format(time.RFC3339),
		CreatedAt: session.CreatedAt.UTC().Format(time.RFC3339),
	}

	if session.RevokedAt.Valid {
		response.RevokedAt = session.RevokedAt.Time.UTC().Format(time.RFC3339)
	}

	return response
}

func requestIP(r *http.Request) string {
	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		return forwarded
	}

	return r.RemoteAddr
}
