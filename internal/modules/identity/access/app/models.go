package app

import "github.com/google/uuid"

type CreateSubRoleInput struct {
	Key            string
	Name           string
	Description    string
	PermissionKeys []string
}

type AssignSubRoleInput struct {
	UserID    uuid.UUID
	SubRoleID uuid.UUID
}

type UpdateBaseRoleInput struct {
	UserID   uuid.UUID
	BaseRole string
}
