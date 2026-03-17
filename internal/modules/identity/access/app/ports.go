package app

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Tokuchi61/Novascans/internal/modules/identity/access/domain"
)

type Repository interface {
	CreatePermission(ctx context.Context, permission domain.Permission) (domain.Permission, error)
	GetPermissionByKey(ctx context.Context, key string) (domain.Permission, error)
	GetPermissionsByKeys(ctx context.Context, keys []string) ([]domain.Permission, error)
	ListPermissions(ctx context.Context) ([]domain.Permission, error)
	CreateSubRole(ctx context.Context, role domain.SubRole) (domain.SubRole, error)
	GetSubRoleByID(ctx context.Context, id uuid.UUID) (domain.SubRole, error)
	GetSubRoleByKey(ctx context.Context, key string) (domain.SubRole, error)
	ListSubRoles(ctx context.Context) ([]domain.SubRole, error)
	AttachPermissionToSubRole(ctx context.Context, subRoleID uuid.UUID, permissionID uuid.UUID, createdAt time.Time) error
	ListSubRolePermissionLinks(ctx context.Context) ([]SubRolePermissionLink, error)
	AssignSubRoleToUser(ctx context.Context, userID uuid.UUID, subRoleID uuid.UUID, assignedAt time.Time) error
	RemoveSubRoleFromUser(ctx context.Context, userID uuid.UUID, subRoleID uuid.UUID) error
	ListSubRolesForUser(ctx context.Context, userID uuid.UUID) ([]domain.SubRole, error)
	ListPermissionKeysForUser(ctx context.Context, userID uuid.UUID) ([]string, error)
	UpdateUserBaseRole(ctx context.Context, userID uuid.UUID, baseRole string, updatedAt time.Time) error
}

type UnitOfWork interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context, repo Repository) error) error
}

type SubRolePermissionLink struct {
	SubRole    domain.SubRole
	Permission *domain.Permission
}
