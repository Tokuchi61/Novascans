package http

import (
	"encoding/json"
	"log/slog"
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	accessapp "github.com/Tokuchi61/Novascans/internal/modules/identity/access/app"
	platformhttp "github.com/Tokuchi61/Novascans/internal/platform/http"
	"github.com/Tokuchi61/Novascans/internal/platform/validation"
)

type Handler struct {
	logger    *slog.Logger
	validator *validation.Validator
	service   *accessapp.Service
}

func NewHandler(logger *slog.Logger, validator *validation.Validator, service *accessapp.Service) *Handler {
	return &Handler{
		logger:    logger,
		validator: validator,
		service:   service,
	}
}

func (h *Handler) Me(w nethttp.ResponseWriter, r *nethttp.Request) {
	principal, ok := CurrentPrincipal(r.Context())
	if !ok {
		platformhttp.WriteError(w, platformhttp.Unauthorized("missing principal"))
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, mapPrincipalResponse(principal))
}

func (h *Handler) ListPermissions(w nethttp.ResponseWriter, r *nethttp.Request) {
	permissions, err := h.service.ListPermissions(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}

	response := make([]permissionResponse, 0, len(permissions))
	for _, permission := range permissions {
		response = append(response, mapPermissionResponse(permission))
	}

	platformhttp.WriteData(w, nethttp.StatusOK, response)
}

func (h *Handler) ListSubRoles(w nethttp.ResponseWriter, r *nethttp.Request) {
	subRoles, err := h.service.ListSubRoles(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}

	response := make([]subRoleResponse, 0, len(subRoles))
	for _, role := range subRoles {
		response = append(response, mapSubRoleResponse(role))
	}

	platformhttp.WriteData(w, nethttp.StatusOK, response)
}

func (h *Handler) CreateSubRole(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request createSubRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		platformhttp.WriteError(w, platformhttp.BadRequest("invalid request body", err))
		return
	}

	if errs := request.Validate(h.validator); errs.HasAny() {
		platformhttp.WriteError(w, platformhttp.ValidationError(errs))
		return
	}

	role, err := h.service.CreateSubRole(r.Context(), request.ToInput())
	if err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusCreated, mapSubRoleResponse(role))
}

func (h *Handler) UpdateBaseRole(w nethttp.ResponseWriter, r *nethttp.Request) {
	userID, err := parseUUIDParam(chi.URLParam(r, "userID"))
	if err != nil {
		platformhttp.WriteError(w, platformhttp.BadRequest("invalid user id", err))
		return
	}

	var request updateBaseRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		platformhttp.WriteError(w, platformhttp.BadRequest("invalid request body", err))
		return
	}

	if errs := request.Validate(h.validator); errs.HasAny() {
		platformhttp.WriteError(w, platformhttp.ValidationError(errs))
		return
	}

	if err := h.service.UpdateBaseRole(r.Context(), request.ToInput(userID)); err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, map[string]bool{"success": true})
}

func (h *Handler) AssignSubRole(w nethttp.ResponseWriter, r *nethttp.Request) {
	userID, err := parseUUIDParam(chi.URLParam(r, "userID"))
	if err != nil {
		platformhttp.WriteError(w, platformhttp.BadRequest("invalid user id", err))
		return
	}

	subRoleID, err := parseUUIDParam(chi.URLParam(r, "subRoleID"))
	if err != nil {
		platformhttp.WriteError(w, platformhttp.BadRequest("invalid sub role id", err))
		return
	}

	if err := h.service.AssignSubRole(r.Context(), accessapp.AssignSubRoleInput{
		UserID:    userID,
		SubRoleID: subRoleID,
	}); err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, map[string]bool{"success": true})
}

func (h *Handler) RemoveSubRole(w nethttp.ResponseWriter, r *nethttp.Request) {
	userID, err := parseUUIDParam(chi.URLParam(r, "userID"))
	if err != nil {
		platformhttp.WriteError(w, platformhttp.BadRequest("invalid user id", err))
		return
	}

	subRoleID, err := parseUUIDParam(chi.URLParam(r, "subRoleID"))
	if err != nil {
		platformhttp.WriteError(w, platformhttp.BadRequest("invalid sub role id", err))
		return
	}

	if err := h.service.RemoveSubRole(r.Context(), accessapp.AssignSubRoleInput{
		UserID:    userID,
		SubRoleID: subRoleID,
	}); err != nil {
		writeServiceError(w, err)
		return
	}

	platformhttp.WriteData(w, nethttp.StatusOK, map[string]bool{"success": true})
}

func parseUUIDParam(raw string) (uuid.UUID, error) {
	return uuid.Parse(raw)
}
