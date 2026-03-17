package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	nethttp "net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	accountapp "github.com/Tokuchi61/Novascans/internal/modules/user/account/app"
	platformhttp "github.com/Tokuchi61/Novascans/internal/platform/http"
	"github.com/Tokuchi61/Novascans/internal/platform/validation"
)

type Handler struct {
	logger    *slog.Logger
	validator *validation.Validator
	service   *accountapp.Service
}

func NewHandler(logger *slog.Logger, validator *validation.Validator, service *accountapp.Service) *Handler {
	return &Handler{
		logger:    logger,
		validator: validator,
		service:   service,
	}
}

func (h *Handler) Me(w nethttp.ResponseWriter, r *nethttp.Request) {
	currentUser, ok := CurrentUser(r.Context())
	if !ok {
		platformhttp.WriteError(w, platformhttp.Unauthorized("authentication required"))
		return
	}

	account, err := h.service.GetAccount(r.Context(), currentUser.User)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, mapAccountResponse(account))
}

func (h *Handler) GetOwnProfile(w nethttp.ResponseWriter, r *nethttp.Request) {
	currentUser, ok := CurrentUser(r.Context())
	if !ok {
		platformhttp.WriteError(w, platformhttp.Unauthorized("authentication required"))
		return
	}

	profile, err := h.service.GetProfile(r.Context(), currentUser.User.ID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, mapProfileResponse(profile))
}

func (h *Handler) UpdateProfile(w nethttp.ResponseWriter, r *nethttp.Request) {
	currentUser, ok := CurrentUser(r.Context())
	if !ok {
		platformhttp.WriteError(w, platformhttp.Unauthorized("authentication required"))
		return
	}

	var request updateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		platformhttp.WriteError(w, platformhttp.BadRequest("invalid request body", err))
		return
	}

	if errs := request.Validate(h.validator); errs.HasAny() {
		platformhttp.WriteError(w, platformhttp.ValidationError(errs))
		return
	}

	profile, err := h.service.UpdateProfile(r.Context(), request.ToInput(currentUser.User.ID))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, mapProfileResponse(profile))
}

func (h *Handler) GetSettings(w nethttp.ResponseWriter, r *nethttp.Request) {
	currentUser, ok := CurrentUser(r.Context())
	if !ok {
		platformhttp.WriteError(w, platformhttp.Unauthorized("authentication required"))
		return
	}

	settings, err := h.service.GetSettings(r.Context(), currentUser.User.ID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, mapSettingsResponse(settings))
}

func (h *Handler) UpdateSettings(w nethttp.ResponseWriter, r *nethttp.Request) {
	currentUser, ok := CurrentUser(r.Context())
	if !ok {
		platformhttp.WriteError(w, platformhttp.Unauthorized("authentication required"))
		return
	}

	var request updateSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		platformhttp.WriteError(w, platformhttp.BadRequest("invalid request body", err))
		return
	}

	if errs := request.Validate(h.validator); errs.HasAny() {
		platformhttp.WriteError(w, platformhttp.ValidationError(errs))
		return
	}

	settings, err := h.service.UpdateSettings(r.Context(), request.ToInput(currentUser.User.ID))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, mapSettingsResponse(settings))
}

func (h *Handler) GetPrivacy(w nethttp.ResponseWriter, r *nethttp.Request) {
	currentUser, ok := CurrentUser(r.Context())
	if !ok {
		platformhttp.WriteError(w, platformhttp.Unauthorized("authentication required"))
		return
	}

	privacy, err := h.service.GetPrivacy(r.Context(), currentUser.User.ID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, mapPrivacyResponse(privacy))
}

func (h *Handler) UpdatePrivacy(w nethttp.ResponseWriter, r *nethttp.Request) {
	currentUser, ok := CurrentUser(r.Context())
	if !ok {
		platformhttp.WriteError(w, platformhttp.Unauthorized("authentication required"))
		return
	}

	var request updatePrivacyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		platformhttp.WriteError(w, platformhttp.BadRequest("invalid request body", err))
		return
	}

	if errs := request.Validate(h.validator); errs.HasAny() {
		platformhttp.WriteError(w, platformhttp.ValidationError(errs))
		return
	}

	privacy, err := h.service.UpdatePrivacy(r.Context(), request.ToInput(currentUser.User.ID))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, mapPrivacyResponse(privacy))
}

func (h *Handler) GetPublicProfile(w nethttp.ResponseWriter, r *nethttp.Request) {
	viewer := accountapp.Viewer{}
	if currentUser, ok := CurrentUser(r.Context()); ok {
		viewer.Authenticated = true
		viewer.UserID = &currentUser.User.ID
	}

	profile, err := h.service.GetPublicProfile(r.Context(), chi.URLParam(r, "username"), viewer)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, mapProfileResponse(profile))
}

func bearerToken(r *nethttp.Request) (string, error) {
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
