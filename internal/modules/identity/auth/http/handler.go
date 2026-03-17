package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	nethttp "net/http"

	authapp "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/app"
	platformhttp "github.com/Tokuchi61/Novascans/internal/platform/http"
	"github.com/Tokuchi61/Novascans/internal/platform/validation"
)

type Handler struct {
	logger    *slog.Logger
	validator *validation.Validator
	service   *authapp.Service
}

func NewHandler(logger *slog.Logger, validator *validation.Validator, service *authapp.Service) *Handler {
	return &Handler{
		logger:    logger,
		validator: validator,
		service:   service,
	}
}

func (h *Handler) Ping(w nethttp.ResponseWriter, r *nethttp.Request) {
	platformhttp.WriteData(w, nethttp.StatusOK, h.service.Ping(r.Context()))
}

func (h *Handler) Register(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request registerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		platformhttp.WriteError(w, platformhttp.BadRequest("invalid request body", err))
		return
	}

	if errs := request.Validate(h.validator); errs.HasAny() {
		platformhttp.WriteError(w, platformhttp.ValidationError(errs))
		return
	}

	result, err := h.service.Register(r.Context(), request.ToInput(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusCreated, mapAuthResponse(result))
}

func (h *Handler) Login(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request loginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		platformhttp.WriteError(w, platformhttp.BadRequest("invalid request body", err))
		return
	}

	if errs := request.Validate(h.validator); errs.HasAny() {
		platformhttp.WriteError(w, platformhttp.ValidationError(errs))
		return
	}

	result, err := h.service.Login(r.Context(), request.ToInput(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, mapAuthResponse(result))
}

func (h *Handler) Refresh(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		platformhttp.WriteError(w, platformhttp.BadRequest("invalid request body", err))
		return
	}

	if errs := request.Validate(h.validator); errs.HasAny() {
		platformhttp.WriteError(w, platformhttp.ValidationError(errs))
		return
	}

	result, err := h.service.Refresh(r.Context(), request.ToInput(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, mapAuthResponse(result))
}

func (h *Handler) Me(w nethttp.ResponseWriter, r *nethttp.Request) {
	currentUser, err := h.currentUser(r)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, mapCurrentUserResponse(currentUser))
}

func (h *Handler) LogoutCurrentSession(w nethttp.ResponseWriter, r *nethttp.Request) {
	currentUser, err := h.currentUser(r)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	if err := h.service.LogoutCurrentSession(r.Context(), currentUser.Session.ID); err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteNoContent(w)
}

func (h *Handler) LogoutAllSessions(w nethttp.ResponseWriter, r *nethttp.Request) {
	currentUser, err := h.currentUser(r)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	if err := h.service.LogoutAllSessions(r.Context(), currentUser.User.ID); err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteNoContent(w)
}

func (h *Handler) RequestEmailVerification(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request emailRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		platformhttp.WriteError(w, platformhttp.BadRequest("invalid request body", err))
		return
	}

	if errs := request.Validate(h.validator); errs.HasAny() {
		platformhttp.WriteError(w, platformhttp.ValidationError(errs))
		return
	}

	result, err := h.service.RequestEmailVerification(r.Context(), request.ToVerificationInput())
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, mapTokenRequestResponse(result))
}

func (h *Handler) VerifyEmail(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request verifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		platformhttp.WriteError(w, platformhttp.BadRequest("invalid request body", err))
		return
	}

	if errs := request.Validate(h.validator); errs.HasAny() {
		platformhttp.WriteError(w, platformhttp.ValidationError(errs))
		return
	}

	if err := h.service.VerifyEmail(r.Context(), request.ToInput()); err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, successResponse{Success: true})
}

func (h *Handler) ForgotPassword(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request emailRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		platformhttp.WriteError(w, platformhttp.BadRequest("invalid request body", err))
		return
	}

	if errs := request.Validate(h.validator); errs.HasAny() {
		platformhttp.WriteError(w, platformhttp.ValidationError(errs))
		return
	}

	result, err := h.service.ForgotPassword(r.Context(), request.ToForgotPasswordInput())
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, mapTokenRequestResponse(result))
}

func (h *Handler) ResetPassword(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request resetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		platformhttp.WriteError(w, platformhttp.BadRequest("invalid request body", err))
		return
	}

	if errs := request.Validate(h.validator); errs.HasAny() {
		platformhttp.WriteError(w, platformhttp.ValidationError(errs))
		return
	}

	if err := h.service.ResetPassword(r.Context(), request.ToInput()); err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, successResponse{Success: true})
}

func (h *Handler) currentUser(r *nethttp.Request) (authapp.AuthenticatedUser, error) {
	token, err := bearerToken(r)
	if err != nil {
		var appErr *platformhttp.AppError
		if errors.As(err, &appErr) {
			return authapp.AuthenticatedUser{}, authapp.Unauthorized(appErr.Message, err)
		}

		return authapp.AuthenticatedUser{}, authapp.Unauthorized("invalid authorization header", err)
	}

	return h.service.AuthenticateAccessToken(r.Context(), token)
}
