package store

import (
	"context"
	"database/sql"
	"time"

	authsqlc "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/store/sqlc"
)

type PostgresRepository struct {
	queries *authsqlc.Queries
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		queries: authsqlc.New(db),
	}
}

func (repo *PostgresRepository) CreateUser(ctx context.Context, params CreateUserParams) (User, error) {
	row, err := repo.queries.CreateUser(ctx, authsqlc.CreateUserParams{
		ID:              params.ID,
		Email:           params.Email,
		Status:          params.Status,
		EmailVerifiedAt: params.EmailVerifiedAt,
		CreatedAt:       params.CreatedAt,
		UpdatedAt:       params.UpdatedAt,
	})
	if err != nil {
		return User{}, err
	}

	return mapUser(row), nil
}

func (repo *PostgresRepository) GetUserByID(ctx context.Context, id string) (User, error) {
	row, err := repo.queries.GetUserByID(ctx, id)
	if err != nil {
		return User{}, err
	}

	return mapUser(row), nil
}

func (repo *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row, err := repo.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return User{}, err
	}

	return mapUser(row), nil
}

func (repo *PostgresRepository) CreatePasswordCredential(ctx context.Context, params CreatePasswordCredentialParams) (PasswordCredential, error) {
	row, err := repo.queries.CreatePasswordCredential(ctx, authsqlc.CreatePasswordCredentialParams{
		UserID:       params.UserID,
		PasswordHash: params.PasswordHash,
		CreatedAt:    params.CreatedAt,
		UpdatedAt:    params.UpdatedAt,
	})
	if err != nil {
		return PasswordCredential{}, err
	}

	return mapPasswordCredential(row), nil
}

func (repo *PostgresRepository) GetPasswordCredentialByEmail(ctx context.Context, email string) (PasswordCredential, error) {
	row, err := repo.queries.GetPasswordCredentialByEmail(ctx, email)
	if err != nil {
		return PasswordCredential{}, err
	}

	return mapPasswordCredential(row), nil
}

func (repo *PostgresRepository) CreateSession(ctx context.Context, params CreateSessionParams) (Session, error) {
	row, err := repo.queries.CreateSession(ctx, authsqlc.CreateSessionParams{
		ID:        params.ID,
		UserID:    params.UserID,
		TokenHash: params.TokenHash,
		UserAgent: params.UserAgent,
		IpAddress: params.IPAddress,
		ExpiresAt: params.ExpiresAt,
		CreatedAt: params.CreatedAt,
	})
	if err != nil {
		return Session{}, err
	}

	return mapSession(row), nil
}

func (repo *PostgresRepository) GetSessionByID(ctx context.Context, id string) (Session, error) {
	row, err := repo.queries.GetSessionByID(ctx, id)
	if err != nil {
		return Session{}, err
	}

	return mapSession(row), nil
}

func (repo *PostgresRepository) RevokeSession(ctx context.Context, id string, revokedAt time.Time) error {
	return repo.queries.RevokeSession(ctx, authsqlc.RevokeSessionParams{
		ID:        id,
		RevokedAt: sql.NullTime{Time: revokedAt, Valid: true},
	})
}

func (repo *PostgresRepository) WithTx(tx *sql.Tx) Repository {
	return &PostgresRepository{
		queries: repo.queries.WithTx(tx),
	}
}

func mapUser(row authsqlc.User) User {
	return User{
		ID:              row.ID,
		Email:           row.Email,
		Status:          row.Status,
		EmailVerifiedAt: row.EmailVerifiedAt,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}
}

func mapPasswordCredential(row authsqlc.AuthPasswordCredential) PasswordCredential {
	return PasswordCredential{
		UserID:       row.UserID,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}

func mapSession(row authsqlc.AuthSession) Session {
	return Session{
		ID:        row.ID,
		UserID:    row.UserID,
		TokenHash: row.TokenHash,
		UserAgent: row.UserAgent,
		IPAddress: row.IpAddress,
		ExpiresAt: row.ExpiresAt,
		RevokedAt: row.RevokedAt,
		CreatedAt: row.CreatedAt,
	}
}
