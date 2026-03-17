package http

import (
	"time"

	"github.com/Tokuchi61/Novascans/internal/modules/identity/access/domain"
)

type principalResponse struct {
	IsGuest        bool              `json:"is_guest"`
	UserID         string            `json:"user_id,omitempty"`
	SessionID      string            `json:"session_id,omitempty"`
	Email          string            `json:"email,omitempty"`
	BaseRole       string            `json:"base_role"`
	SubRoles       []subRoleResponse `json:"sub_roles"`
	PermissionKeys []string          `json:"permission_keys"`
}

type subRoleResponse struct {
	ID             string   `json:"id"`
	Key            string   `json:"key"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	PermissionKeys []string `json:"permission_keys"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

type permissionResponse struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func mapPrincipalResponse(principal domain.Principal) principalResponse {
	response := principalResponse{
		IsGuest:        principal.IsGuest,
		Email:          principal.Email,
		BaseRole:       principal.BaseRole,
		SubRoles:       make([]subRoleResponse, 0, len(principal.SubRoles)),
		PermissionKeys: principal.PermissionKeys,
	}

	if principal.UserID != nil {
		response.UserID = principal.UserID.String()
	}

	if principal.SessionID != nil {
		response.SessionID = principal.SessionID.String()
	}

	for _, role := range principal.SubRoles {
		response.SubRoles = append(response.SubRoles, mapSubRoleResponse(role))
	}

	return response
}

func mapSubRoleResponse(role domain.SubRole) subRoleResponse {
	return subRoleResponse{
		ID:             role.ID.String(),
		Key:            role.Key,
		Name:           role.Name,
		Description:    role.Description,
		PermissionKeys: role.PermissionKeys,
		CreatedAt:      role.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:      role.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func mapPermissionResponse(permission domain.Permission) permissionResponse {
	return permissionResponse{
		ID:          permission.ID.String(),
		Key:         permission.Key,
		Description: permission.Description,
		CreatedAt:   permission.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   permission.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
