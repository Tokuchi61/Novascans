package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	accesssqlc "github.com/Tokuchi61/Novascans/internal/gen/sqlc/identity/access"
	accessapp "github.com/Tokuchi61/Novascans/internal/modules/identity/access/app"
	"github.com/Tokuchi61/Novascans/internal/modules/identity/access/domain"
	platformdb "github.com/Tokuchi61/Novascans/internal/platform/db"
)

type PostgresRepository struct {
	queries *accesssqlc.Queries
}

type PostgresUnitOfWork struct {
	db        *sql.DB
	txManager platformdb.TxManager
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return newPostgresRepository(accesssqlc.New(db))
}

func NewPostgresUnitOfWork(db *sql.DB, txManager platformdb.TxManager) *PostgresUnitOfWork {
	return &PostgresUnitOfWork{
		db:        db,
		txManager: txManager,
	}
}

func newPostgresRepository(queries *accesssqlc.Queries) *PostgresRepository {
	return &PostgresRepository{queries: queries}
}

func (repo *PostgresRepository) CreatePermission(ctx context.Context, permission domain.Permission) (domain.Permission, error) {
	row, err := repo.queries.CreatePermission(ctx, accesssqlc.CreatePermissionParams{
		ID:          permission.ID,
		Key:         permission.Key,
		Description: permission.Description,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	})
	if err != nil {
		return domain.Permission{}, err
	}

	return mapPermission(row), nil
}

func (repo *PostgresRepository) GetPermissionByKey(ctx context.Context, key string) (domain.Permission, error) {
	row, err := repo.queries.GetPermissionByKey(ctx, key)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Permission{}, accessapp.NotFound("permission not found", err)
		}

		return domain.Permission{}, err
	}

	return mapPermission(row), nil
}

func (repo *PostgresRepository) GetPermissionsByKeys(ctx context.Context, keys []string) ([]domain.Permission, error) {
	rows, err := repo.queries.GetPermissionsByKeys(ctx, keys)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Permission, 0, len(rows))
	for _, row := range rows {
		result = append(result, mapPermission(row))
	}

	return result, nil
}

func (repo *PostgresRepository) ListPermissions(ctx context.Context) ([]domain.Permission, error) {
	rows, err := repo.queries.ListPermissions(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Permission, 0, len(rows))
	for _, row := range rows {
		result = append(result, mapPermission(row))
	}

	return result, nil
}

func (repo *PostgresRepository) CreateSubRole(ctx context.Context, role domain.SubRole) (domain.SubRole, error) {
	row, err := repo.queries.CreateSubRole(ctx, accesssqlc.CreateSubRoleParams{
		ID:          role.ID,
		Key:         role.Key,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	})
	if err != nil {
		return domain.SubRole{}, err
	}

	return mapSubRole(row), nil
}

func (repo *PostgresRepository) GetSubRoleByID(ctx context.Context, id uuid.UUID) (domain.SubRole, error) {
	row, err := repo.queries.GetSubRoleByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.SubRole{}, accessapp.NotFound("sub role not found", err)
		}

		return domain.SubRole{}, err
	}

	return mapSubRole(row), nil
}

func (repo *PostgresRepository) GetSubRoleByKey(ctx context.Context, key string) (domain.SubRole, error) {
	row, err := repo.queries.GetSubRoleByKey(ctx, key)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.SubRole{}, accessapp.NotFound("sub role not found", err)
		}

		return domain.SubRole{}, err
	}

	return mapSubRole(row), nil
}

func (repo *PostgresRepository) ListSubRoles(ctx context.Context) ([]domain.SubRole, error) {
	rows, err := repo.queries.ListSubRoles(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.SubRole, 0, len(rows))
	for _, row := range rows {
		result = append(result, mapSubRole(row))
	}

	return result, nil
}

func (repo *PostgresRepository) AttachPermissionToSubRole(ctx context.Context, subRoleID uuid.UUID, permissionID uuid.UUID, createdAt time.Time) error {
	return repo.queries.AttachPermissionToSubRole(ctx, accesssqlc.AttachPermissionToSubRoleParams{
		SubRoleID:    subRoleID,
		PermissionID: permissionID,
		CreatedAt:    createdAt,
	})
}

func (repo *PostgresRepository) ListSubRolePermissionLinks(ctx context.Context) ([]accessapp.SubRolePermissionLink, error) {
	rows, err := repo.queries.ListSubRolePermissionLinks(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]accessapp.SubRolePermissionLink, 0, len(rows))
	for _, row := range rows {
		link := accessapp.SubRolePermissionLink{
			SubRole: domain.SubRole{
				ID:          row.SubRoleID,
				Key:         row.SubRoleKey,
				Name:        row.SubRoleName,
				Description: row.SubRoleDescription,
				CreatedAt:   row.SubRoleCreatedAt,
				UpdatedAt:   row.SubRoleUpdatedAt,
			},
		}

		if row.PermissionID.Valid {
			link.Permission = &domain.Permission{
				ID:          row.PermissionID.UUID,
				Key:         row.PermissionKey.String,
				Description: row.PermissionDescription.String,
				CreatedAt:   row.PermissionCreatedAt.Time,
				UpdatedAt:   row.PermissionUpdatedAt.Time,
			}
		}

		result = append(result, link)
	}

	return result, nil
}

func (repo *PostgresRepository) AssignSubRoleToUser(ctx context.Context, userID uuid.UUID, subRoleID uuid.UUID, assignedAt time.Time) error {
	return repo.queries.AssignSubRoleToUser(ctx, accesssqlc.AssignSubRoleToUserParams{
		UserID:     userID,
		SubRoleID:  subRoleID,
		AssignedAt: assignedAt,
	})
}

func (repo *PostgresRepository) RemoveSubRoleFromUser(ctx context.Context, userID uuid.UUID, subRoleID uuid.UUID) error {
	return repo.queries.RemoveSubRoleFromUser(ctx, accesssqlc.RemoveSubRoleFromUserParams{
		UserID:    userID,
		SubRoleID: subRoleID,
	})
}

func (repo *PostgresRepository) ListSubRolesForUser(ctx context.Context, userID uuid.UUID) ([]domain.SubRole, error) {
	rows, err := repo.queries.ListSubRolesForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]domain.SubRole, 0, len(rows))
	for _, row := range rows {
		result = append(result, mapSubRole(row))
	}

	return result, nil
}

func (repo *PostgresRepository) ListPermissionKeysForUser(ctx context.Context, userID uuid.UUID) ([]string, error) {
	return repo.queries.ListPermissionKeysForUser(ctx, userID)
}

func (repo *PostgresRepository) UpdateUserBaseRole(ctx context.Context, userID uuid.UUID, baseRole string, updatedAt time.Time) error {
	return repo.queries.UpdateUserBaseRole(ctx, accesssqlc.UpdateUserBaseRoleParams{
		ID:        userID,
		BaseRole:  baseRole,
		UpdatedAt: updatedAt,
	})
}

func (uow *PostgresUnitOfWork) WithinTransaction(ctx context.Context, fn func(ctx context.Context, repo accessapp.Repository) error) error {
	if uow.txManager == nil {
		return fn(ctx, newPostgresRepository(accesssqlc.New(uow.db)))
	}

	return uow.txManager.WithinTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		return fn(ctx, newPostgresRepository(accesssqlc.New(tx)))
	})
}

func mapPermission(row accesssqlc.AccessPermission) domain.Permission {
	return domain.Permission{
		ID:          row.ID,
		Key:         row.Key,
		Description: row.Description,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func mapSubRole(row accesssqlc.AccessSubRole) domain.SubRole {
	return domain.SubRole{
		ID:          row.ID,
		Key:         row.Key,
		Name:        row.Name,
		Description: row.Description,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}
