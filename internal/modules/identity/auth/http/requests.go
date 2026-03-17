package http

import (
	"errors"
	"net/http"
	"strings"

	authapp "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/app"
	platformhttp "github.com/Tokuchi61/Novascans/internal/platform/http"
	"github.com/Tokuchi61/Novascans/internal/platform/validation"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type emailRequest struct {
	Email string `json:"email"`
}

type verifyEmailRequest struct {
	Token string `json:"token"`
}

type resetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

func (request registerRequest) Validate(v *validation.Validator) validation.FieldErrors {
	errs := validation.FieldErrors{}
	v.Email("email", request.Email, errs)
	v.MinLength("password", request.Password, 8, errs)
	return errs
}

func (request registerRequest) ToInput(r *http.Request) authapp.RegisterInput {
	return authapp.RegisterInput{
		Email:     strings.ToLower(strings.TrimSpace(request.Email)),
		Password:  request.Password,
		UserAgent: r.UserAgent(),
		IPAddress: requestIP(r),
	}
}

func (request loginRequest) Validate(v *validation.Validator) validation.FieldErrors {
	errs := validation.FieldErrors{}
	v.Email("email", request.Email, errs)
	v.RequiredString("password", request.Password, errs)
	return errs
}

func (request loginRequest) ToInput(r *http.Request) authapp.LoginInput {
	return authapp.LoginInput{
		Email:     strings.ToLower(strings.TrimSpace(request.Email)),
		Password:  request.Password,
		UserAgent: r.UserAgent(),
		IPAddress: requestIP(r),
	}
}

func (request refreshRequest) Validate(v *validation.Validator) validation.FieldErrors {
	errs := validation.FieldErrors{}
	v.RequiredString("refresh_token", request.RefreshToken, errs)
	return errs
}

func (request refreshRequest) ToInput(r *http.Request) authapp.RefreshSessionInput {
	return authapp.RefreshSessionInput{
		RefreshToken: strings.TrimSpace(request.RefreshToken),
		UserAgent:    r.UserAgent(),
		IPAddress:    requestIP(r),
	}
}

func (request emailRequest) Validate(v *validation.Validator) validation.FieldErrors {
	errs := validation.FieldErrors{}
	v.Email("email", request.Email, errs)
	return errs
}

func (request emailRequest) ToVerificationInput() authapp.RequestEmailVerificationInput {
	return authapp.RequestEmailVerificationInput{
		Email: strings.ToLower(strings.TrimSpace(request.Email)),
	}
}

func (request emailRequest) ToForgotPasswordInput() authapp.ForgotPasswordInput {
	return authapp.ForgotPasswordInput{
		Email: strings.ToLower(strings.TrimSpace(request.Email)),
	}
}

func (request verifyEmailRequest) Validate(v *validation.Validator) validation.FieldErrors {
	errs := validation.FieldErrors{}
	v.RequiredString("token", request.Token, errs)
	return errs
}

func (request verifyEmailRequest) ToInput() authapp.VerifyEmailInput {
	return authapp.VerifyEmailInput{Token: strings.TrimSpace(request.Token)}
}

func (request resetPasswordRequest) Validate(v *validation.Validator) validation.FieldErrors {
	errs := validation.FieldErrors{}
	v.RequiredString("token", request.Token, errs)
	v.MinLength("new_password", request.NewPassword, 8, errs)
	return errs
}

func (request resetPasswordRequest) ToInput() authapp.ResetPasswordInput {
	return authapp.ResetPasswordInput{
		Token:       strings.TrimSpace(request.Token),
		NewPassword: request.NewPassword,
	}
}

func requestIP(r *http.Request) string {
	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		return forwarded
	}

	return r.RemoteAddr
}

func bearerToken(r *http.Request) (string, error) {
	value := strings.TrimSpace(r.Header.Get("Authorization"))
	if value == "" {
		return "", platformhttp.Unauthorized("missing authorization header")
	}

	parts := strings.SplitN(value, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid authorization header")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errors.New("missing bearer token")
	}

	return token, nil
}
