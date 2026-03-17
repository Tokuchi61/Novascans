package app

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/Tokuchi61/Novascans/internal/modules/identity/access/domain"
	authapp "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/app"
	authdomain "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/domain"
)

var assignableBaseRoles = map[string]struct{}{
	"user":      {},
	"moderator": {},
	"admin":     {},
}

type Service struct {
	repo       Repository
	unitOfWork UnitOfWork
}

func NewService(repo Repository, unitOfWork UnitOfWork) *Service {
	return &Service{
		repo:       repo,
		unitOfWork: unitOfWork,
	}
}

func (service *Service) ResolvePrincipal(ctx context.Context, user authdomain.User, sessionID uuid.UUID) (domain.Principal, error) {
	subRoles, err := service.repo.ListSubRolesForUser(ctx, user.ID)
	if err != nil {
		return domain.Principal{}, Internal("failed to load user sub roles", err)
	}

	permissions, err := service.repo.ListPermissionKeysForUser(ctx, user.ID)
	if err != nil {
		return domain.Principal{}, Internal("failed to load user permissions", err)
	}

	return domain.Principal{
		UserID:          &user.ID,
		SessionID:       &sessionID,
		Email:           user.Email,
		BaseRole:        user.BaseRole,
		Status:          user.Status,
		EmailVerifiedAt: user.EmailVerifiedAt,
		SubRoles:        subRoles,
		PermissionKeys:  permissions,
	}, nil
}

func (service *Service) GuestPrincipal() domain.Principal {
	return domain.GuestPrincipal()
}

func (service *Service) ListPermissions(ctx context.Context) ([]domain.Permission, error) {
	permissions, err := service.repo.ListPermissions(ctx)
	if err != nil {
		return nil, Internal("failed to list permissions", err)
	}

	return permissions, nil
}

func (service *Service) ListSubRoles(ctx context.Context) ([]domain.SubRole, error) {
	links, err := service.repo.ListSubRolePermissionLinks(ctx)
	if err != nil {
		return nil, Internal("failed to list sub roles", err)
	}

	return aggregateSubRoles(links), nil
}

func (service *Service) CreateSubRole(ctx context.Context, input CreateSubRoleInput) (domain.SubRole, error) {
	key := normalizeKey(input.Key)
	if key == "" {
		return domain.SubRole{}, BadRequest("sub role key is required", nil)
	}

	if strings.TrimSpace(input.Name) == "" {
		return domain.SubRole{}, BadRequest("sub role name is required", nil)
	}

	permissionKeys := normalizeKeys(input.PermissionKeys)
	if len(permissionKeys) == 0 {
		return domain.SubRole{}, BadRequest("at least one permission key is required", nil)
	}

	if _, err := service.repo.GetSubRoleByKey(ctx, key); err == nil {
		return domain.SubRole{}, Conflict("sub role already exists", nil)
	} else if !HasCode(err, CodeNotFound) {
		return domain.SubRole{}, Internal("failed to check existing sub role", err)
	}

	permissions, err := service.repo.GetPermissionsByKeys(ctx, permissionKeys)
	if err != nil {
		return domain.SubRole{}, Internal("failed to resolve permission keys", err)
	}

	if len(permissions) != len(permissionKeys) {
		return domain.SubRole{}, BadRequest("one or more permission keys are invalid", nil)
	}

	now := time.Now().UTC()
	role := domain.SubRole{
		ID:             uuid.New(),
		Key:            key,
		Name:           strings.TrimSpace(input.Name),
		Description:    strings.TrimSpace(input.Description),
		PermissionKeys: permissionKeys,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := service.withWriteRepository(ctx, func(ctx context.Context, repo Repository) error {
		createdRole, err := repo.CreateSubRole(ctx, role)
		if err != nil {
			return Internal("failed to create sub role", err)
		}

		for _, permission := range permissions {
			if err := repo.AttachPermissionToSubRole(ctx, createdRole.ID, permission.ID, now); err != nil {
				return Internal("failed to attach permission to sub role", err)
			}
		}

		role = createdRole
		return nil
	}); err != nil {
		return domain.SubRole{}, err
	}

	role.PermissionKeys = permissionKeys
	return role, nil
}

func (service *Service) AssignSubRole(ctx context.Context, input AssignSubRoleInput) error {
	if _, err := service.repo.GetSubRoleByID(ctx, input.SubRoleID); err != nil {
		if HasCode(err, CodeNotFound) {
			return NotFound("sub role not found", err)
		}

		return Internal("failed to fetch sub role", err)
	}

	if err := service.repo.AssignSubRoleToUser(ctx, input.UserID, input.SubRoleID, time.Now().UTC()); err != nil {
		return Internal("failed to assign sub role", err)
	}

	return nil
}

func (service *Service) RemoveSubRole(ctx context.Context, input AssignSubRoleInput) error {
	if err := service.repo.RemoveSubRoleFromUser(ctx, input.UserID, input.SubRoleID); err != nil {
		return Internal("failed to remove sub role", err)
	}

	return nil
}

func (service *Service) UpdateBaseRole(ctx context.Context, input UpdateBaseRoleInput) error {
	baseRole := normalizeKey(input.BaseRole)
	if _, ok := assignableBaseRoles[baseRole]; !ok {
		return BadRequest("invalid base role", nil)
	}

	if err := service.repo.UpdateUserBaseRole(ctx, input.UserID, baseRole, time.Now().UTC()); err != nil {
		return Internal("failed to update base role", err)
	}

	return nil
}

func (service *Service) AuthenticatePrincipal(ctx context.Context, authenticator authapp.Authenticator, bearerToken string) (domain.Principal, error) {
	currentUser, err := authenticator.AuthenticateAccessToken(ctx, bearerToken)
	if err != nil {
		return domain.Principal{}, err
	}

	return service.ResolvePrincipal(ctx, currentUser.User, currentUser.Session.ID)
}

func (service *Service) withWriteRepository(ctx context.Context, fn func(ctx context.Context, repo Repository) error) error {
	if service.unitOfWork == nil {
		return fn(ctx, service.repo)
	}

	return service.unitOfWork.WithinTransaction(ctx, fn)
}

func aggregateSubRoles(links []SubRolePermissionLink) []domain.SubRole {
	if len(links) == 0 {
		return []domain.SubRole{}
	}

	roles := make(map[uuid.UUID]domain.SubRole)
	order := make([]uuid.UUID, 0, len(links))

	for _, link := range links {
		role, exists := roles[link.SubRole.ID]
		if !exists {
			role = link.SubRole
			role.PermissionKeys = []string{}
			roles[link.SubRole.ID] = role
			order = append(order, link.SubRole.ID)
		}

		if link.Permission != nil {
			role.PermissionKeys = append(role.PermissionKeys, link.Permission.Key)
			roles[link.SubRole.ID] = role
		}
	}

	result := make([]domain.SubRole, 0, len(order))
	for _, id := range order {
		result = append(result, roles[id])
	}

	return result
}

func normalizeKey(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func normalizeKeys(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	keys := make([]string, 0, len(values))

	for _, value := range values {
		normalized := normalizeKey(value)
		if normalized == "" {
			continue
		}

		if _, exists := seen[normalized]; exists {
			continue
		}

		seen[normalized] = struct{}{}
		keys = append(keys, normalized)
	}

	return keys
}
