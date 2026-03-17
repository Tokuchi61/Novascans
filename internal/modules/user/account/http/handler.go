package http

import (
	"context"
	"encoding/json"
	"log/slog"
	nethttp "net/http"

	"github.com/go-chi/chi/v5"

	accessdomain "github.com/Tokuchi61/Novascans/internal/modules/identity/access/domain"
	accesshttp "github.com/Tokuchi61/Novascans/internal/modules/identity/access/http"
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
	principal, ok := currentPrincipal(r.Context())
	if !ok {
		platformhttp.WriteError(w, platformhttp.Unauthorized("authentication required"))
		return
	}

	account, err := h.service.GetAccount(r.Context(), *principal.UserID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, mapAccountResponse(principal, account))
}

func (h *Handler) GetOwnProfile(w nethttp.ResponseWriter, r *nethttp.Request) {
	principal, ok := currentPrincipal(r.Context())
	if !ok {
		platformhttp.WriteError(w, platformhttp.Unauthorized("authentication required"))
		return
	}

	profile, err := h.service.GetProfile(r.Context(), *principal.UserID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, mapProfileResponse(profile))
}

func (h *Handler) UpdateProfile(w nethttp.ResponseWriter, r *nethttp.Request) {
	principal, ok := currentPrincipal(r.Context())
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

	profile, err := h.service.UpdateProfile(r.Context(), request.ToInput(*principal.UserID))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, mapProfileResponse(profile))
}

func (h *Handler) GetSettings(w nethttp.ResponseWriter, r *nethttp.Request) {
	principal, ok := currentPrincipal(r.Context())
	if !ok {
		platformhttp.WriteError(w, platformhttp.Unauthorized("authentication required"))
		return
	}

	settings, err := h.service.GetSettings(r.Context(), *principal.UserID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, mapSettingsResponse(settings))
}

func (h *Handler) UpdateSettings(w nethttp.ResponseWriter, r *nethttp.Request) {
	principal, ok := currentPrincipal(r.Context())
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

	settings, err := h.service.UpdateSettings(r.Context(), request.ToInput(*principal.UserID))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, mapSettingsResponse(settings))
}

func (h *Handler) GetPrivacy(w nethttp.ResponseWriter, r *nethttp.Request) {
	principal, ok := currentPrincipal(r.Context())
	if !ok {
		platformhttp.WriteError(w, platformhttp.Unauthorized("authentication required"))
		return
	}

	privacy, err := h.service.GetPrivacy(r.Context(), *principal.UserID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, mapPrivacyResponse(privacy))
}

func (h *Handler) UpdatePrivacy(w nethttp.ResponseWriter, r *nethttp.Request) {
	principal, ok := currentPrincipal(r.Context())
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

	privacy, err := h.service.UpdatePrivacy(r.Context(), request.ToInput(*principal.UserID))
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, mapPrivacyResponse(privacy))
}

func (h *Handler) GetPublicProfile(w nethttp.ResponseWriter, r *nethttp.Request) {
	viewer := accountapp.Viewer{}
	if principal, ok := accesshttp.CurrentPrincipal(r.Context()); ok && principal.IsAuthenticated() {
		viewer.Authenticated = true
		viewer.UserID = principal.UserID
	}

	profile, err := h.service.GetPublicProfile(r.Context(), chi.URLParam(r, "username"), viewer)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, mapProfileResponse(profile))
}

func currentPrincipal(ctx context.Context) (accessdomain.Principal, bool) {
	principal, ok := accesshttp.CurrentPrincipal(ctx)
	if !ok || !principal.IsAuthenticated() || principal.UserID == nil {
		return accessdomain.Principal{}, false
	}

	return principal, true
}
