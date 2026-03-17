package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	platformhttp "github.com/Tokuchi61/Novascans/internal/platform/http"
	"github.com/Tokuchi61/Novascans/internal/platform/validation"
)

type Handler struct {
	logger    *slog.Logger
	validator *validation.Validator
	service   *Service
}

func NewHandler(logger *slog.Logger, validator *validation.Validator, service *Service) *Handler {
	return &Handler{
		logger:    logger,
		validator: validator,
		service:   service,
	}
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	platformhttp.WriteData(w, http.StatusOK, h.service.Ping(r.Context()))
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var request createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		platformhttp.WriteError(w, platformhttp.BadRequest("invalid request body", err))
		return
	}

	if errs := request.Validate(h.validator); errs.HasAny() {
		platformhttp.WriteError(w, platformhttp.ValidationError(errs))
		return
	}

	user, err := h.service.CreateUser(r.Context(), CreateUserInput{
		Email:    strings.ToLower(strings.TrimSpace(request.Email)),
		Password: request.Password,
	})
	if err != nil {
		platformhttp.WriteError(w, err)
		return
	}

	platformhttp.WriteData(w, http.StatusCreated, mapUserResponse(user))
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	user, err := h.service.GetUserByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		platformhttp.WriteError(w, err)
		return
	}

	platformhttp.WriteData(w, http.StatusOK, mapUserResponse(user))
}

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var request createSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		platformhttp.WriteError(w, platformhttp.BadRequest("invalid request body", err))
		return
	}

	if errs := request.Validate(h.validator); errs.HasAny() {
		platformhttp.WriteError(w, platformhttp.ValidationError(errs))
		return
	}

	session, err := h.service.CreateSession(r.Context(), CreateSessionInput{
		Email:     strings.ToLower(strings.TrimSpace(request.Email)),
		Password:  request.Password,
		UserAgent: r.UserAgent(),
		IPAddress: requestIP(r),
	})
	if err != nil {
		platformhttp.WriteError(w, err)
		return
	}

	platformhttp.WriteData(w, http.StatusCreated, mapSessionResponse(session))
}

func (h *Handler) RevokeSession(w http.ResponseWriter, r *http.Request) {
	if err := h.service.RevokeSession(r.Context(), chi.URLParam(r, "id")); err != nil {
		platformhttp.WriteError(w, err)
		return
	}

	platformhttp.WriteNoContent(w)
}
