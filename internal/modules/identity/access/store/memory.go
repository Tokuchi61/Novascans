package store

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	accessapp "github.com/Tokuchi61/Novascans/internal/modules/identity/access/app"
	"github.com/Tokuchi61/Novascans/internal/modules/identity/access/domain"
)

type MemoryRepository struct {
	mu                sync.RWMutex
	permissions       map[uuid.UUID]domain.Permission
	permissionsByKey  map[string]uuid.UUID
	subRoles          map[uuid.UUID]domain.SubRole
	subRolesByKey     map[string]uuid.UUID
	subRolePermission map[uuid.UUID]map[uuid.UUID]struct{}
	userSubRoles      map[uuid.UUID]map[uuid.UUID]struct{}
	userBaseRoles     map[uuid.UUID]string
}

type MemoryUnitOfWork struct {
	repo *MemoryRepository
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		permissions:       make(map[uuid.UUID]domain.Permission),
		permissionsByKey:  make(map[string]uuid.UUID),
		subRoles:          make(map[uuid.UUID]domain.SubRole),
		subRolesByKey:     make(map[string]uuid.UUID),
		subRolePermission: make(map[uuid.UUID]map[uuid.UUID]struct{}),
		userSubRoles:      make(map[uuid.UUID]map[uuid.UUID]struct{}),
		userBaseRoles:     make(map[uuid.UUID]string),
	}
}

func NewMemoryUnitOfWork(repo *MemoryRepository) *MemoryUnitOfWork {
	return &MemoryUnitOfWork{repo: repo}
}

func (repo *MemoryRepository) CreatePermission(_ context.Context, permission domain.Permission) (domain.Permission, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.permissions[permission.ID] = permission
	repo.permissionsByKey[permission.Key] = permission.ID
	return permission, nil
}

func (repo *MemoryRepository) GetPermissionByKey(_ context.Context, key string) (domain.Permission, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	id, ok := repo.permissionsByKey[key]
	if !ok {
		return domain.Permission{}, accessapp.NotFound("permission not found", nil)
	}

	return repo.permissions[id], nil
}

func (repo *MemoryRepository) GetPermissionsByKeys(_ context.Context, keys []string) ([]domain.Permission, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	result := make([]domain.Permission, 0, len(keys))
	for _, key := range keys {
		id, ok := repo.permissionsByKey[key]
		if !ok {
			continue
		}

		result = append(result, repo.permissions[id])
	}

	return result, nil
}

func (repo *MemoryRepository) ListPermissions(_ context.Context) ([]domain.Permission, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	result := make([]domain.Permission, 0, len(repo.permissions))
	for _, permission := range repo.permissions {
		result = append(result, permission)
	}

	return result, nil
}

func (repo *MemoryRepository) CreateSubRole(_ context.Context, role domain.SubRole) (domain.SubRole, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.subRoles[role.ID] = role
	repo.subRolesByKey[role.Key] = role.ID
	return role, nil
}

func (repo *MemoryRepository) GetSubRoleByID(_ context.Context, id uuid.UUID) (domain.SubRole, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	role, ok := repo.subRoles[id]
	if !ok {
		return domain.SubRole{}, accessapp.NotFound("sub role not found", nil)
	}

	return role, nil
}

func (repo *MemoryRepository) GetSubRoleByKey(_ context.Context, key string) (domain.SubRole, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	id, ok := repo.subRolesByKey[key]
	if !ok {
		return domain.SubRole{}, accessapp.NotFound("sub role not found", nil)
	}

	return repo.subRoles[id], nil
}

func (repo *MemoryRepository) ListSubRoles(_ context.Context) ([]domain.SubRole, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	result := make([]domain.SubRole, 0, len(repo.subRoles))
	for _, role := range repo.subRoles {
		result = append(result, role)
	}

	return result, nil
}

func (repo *MemoryRepository) AttachPermissionToSubRole(_ context.Context, subRoleID uuid.UUID, permissionID uuid.UUID, _ time.Time) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, ok := repo.subRolePermission[subRoleID]; !ok {
		repo.subRolePermission[subRoleID] = make(map[uuid.UUID]struct{})
	}

	repo.subRolePermission[subRoleID][permissionID] = struct{}{}
	return nil
}

func (repo *MemoryRepository) ListSubRolePermissionLinks(_ context.Context) ([]accessapp.SubRolePermissionLink, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	result := make([]accessapp.SubRolePermissionLink, 0, len(repo.subRoles))
	for _, role := range repo.subRoles {
		permissionIDs := repo.subRolePermission[role.ID]
		if len(permissionIDs) == 0 {
			result = append(result, accessapp.SubRolePermissionLink{SubRole: role})
			continue
		}

		for permissionID := range permissionIDs {
			permission := repo.permissions[permissionID]
			permissionCopy := permission
			result = append(result, accessapp.SubRolePermissionLink{
				SubRole:    role,
				Permission: &permissionCopy,
			})
		}
	}

	return result, nil
}

func (repo *MemoryRepository) AssignSubRoleToUser(_ context.Context, userID uuid.UUID, subRoleID uuid.UUID, _ time.Time) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, ok := repo.userSubRoles[userID]; !ok {
		repo.userSubRoles[userID] = make(map[uuid.UUID]struct{})
	}

	repo.userSubRoles[userID][subRoleID] = struct{}{}
	return nil
}

func (repo *MemoryRepository) RemoveSubRoleFromUser(_ context.Context, userID uuid.UUID, subRoleID uuid.UUID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if roles, ok := repo.userSubRoles[userID]; ok {
		delete(roles, subRoleID)
	}

	return nil
}

func (repo *MemoryRepository) ListSubRolesForUser(_ context.Context, userID uuid.UUID) ([]domain.SubRole, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	roleIDs := repo.userSubRoles[userID]
	result := make([]domain.SubRole, 0, len(roleIDs))
	for roleID := range roleIDs {
		if role, ok := repo.subRoles[roleID]; ok {
			result = append(result, role)
		}
	}

	return result, nil
}

func (repo *MemoryRepository) ListPermissionKeysForUser(_ context.Context, userID uuid.UUID) ([]string, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	seen := make(map[string]struct{})
	result := []string{}
	for roleID := range repo.userSubRoles[userID] {
		for permissionID := range repo.subRolePermission[roleID] {
			permission, ok := repo.permissions[permissionID]
			if !ok {
				continue
			}

			if _, exists := seen[permission.Key]; exists {
				continue
			}

			seen[permission.Key] = struct{}{}
			result = append(result, permission.Key)
		}
	}

	return result, nil
}

func (repo *MemoryRepository) UpdateUserBaseRole(_ context.Context, userID uuid.UUID, baseRole string, _ time.Time) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.userBaseRoles[userID] = baseRole
	return nil
}

func (uow *MemoryUnitOfWork) WithinTransaction(ctx context.Context, fn func(ctx context.Context, repo accessapp.Repository) error) error {
	return fn(ctx, uow.repo)
}
