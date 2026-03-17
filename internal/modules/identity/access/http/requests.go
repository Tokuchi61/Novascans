package http

import (
	"strings"

	"github.com/google/uuid"

	accessapp "github.com/Tokuchi61/Novascans/internal/modules/identity/access/app"
	"github.com/Tokuchi61/Novascans/internal/platform/validation"
)

type createSubRoleRequest struct {
	Key            string   `json:"key"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	PermissionKeys []string `json:"permission_keys"`
}

type updateBaseRoleRequest struct {
	BaseRole string `json:"base_role"`
}

func (request createSubRoleRequest) Validate(v *validation.Validator) validation.FieldErrors {
	errs := validation.FieldErrors{}
	v.RequiredString("key", request.Key, errs)
	v.RequiredString("name", request.Name, errs)
	if len(request.PermissionKeys) == 0 {
		errs.Add("permission_keys", "required")
	}
	return errs
}

func (request createSubRoleRequest) ToInput() accessapp.CreateSubRoleInput {
	return accessapp.CreateSubRoleInput{
		Key:            strings.TrimSpace(request.Key),
		Name:           strings.TrimSpace(request.Name),
		Description:    strings.TrimSpace(request.Description),
		PermissionKeys: request.PermissionKeys,
	}
}

func (request updateBaseRoleRequest) Validate(v *validation.Validator) validation.FieldErrors {
	errs := validation.FieldErrors{}
	v.RequiredString("base_role", request.BaseRole, errs)
	return errs
}

func (request updateBaseRoleRequest) ToInput(userID uuid.UUID) accessapp.UpdateBaseRoleInput {
	return accessapp.UpdateBaseRoleInput{
		UserID:   userID,
		BaseRole: strings.TrimSpace(request.BaseRole),
	}
}
