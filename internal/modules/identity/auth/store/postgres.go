package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	authsqlc "github.com/Tokuchi61/Novascans/internal/gen/sqlc/identity/auth"
	authapp "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/app"
	"github.com/Tokuchi61/Novascans/internal/modules/identity/auth/domain"
	platformdb "github.com/Tokuchi61/Novascans/internal/platform/db"
)

type PostgresRepository struct {
	queries *authsqlc.Queries
}

type PostgresUnitOfWork struct {
	db        *sql.DB
	txManager platformdb.TxManager
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return newPostgresRepository(authsqlc.New(db))
}

func NewPostgresUnitOfWork(db *sql.DB, txManager platformdb.TxManager) *PostgresUnitOfWork {
	return &PostgresUnitOfWork{
		db:        db,
		txManager: txManager,
	}
}

func newPostgresRepository(queries *authsqlc.Queries) *PostgresRepository {
	return &PostgresRepository{queries: queries}
}

func (repo *PostgresRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	row, err := repo.queries.CreateUser(ctx, authsqlc.CreateUserParams{
		ID:              user.ID,
		Email:           user.Email,
		BaseRole:        user.BaseRole,
		Status:          user.Status,
		EmailVerifiedAt: nullableTime(user.EmailVerifiedAt),
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	})
	if err != nil {
		return domain.User{}, err
	}

	return mapUser(row), nil
}

func (repo *PostgresRepository) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	row, err := repo.queries.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, authapp.NotFound("user not found", err)
		}

		return domain.User{}, err
	}

	return mapUser(row), nil
}

func (repo *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	row, err := repo.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, authapp.NotFound("user not found", err)
		}

		return domain.User{}, err
	}

	return mapUser(row), nil
}

func (repo *PostgresRepository) MarkUserEmailVerified(ctx context.Context, userID uuid.UUID, verifiedAt time.Time) error {
	return repo.queries.MarkUserEmailVerified(ctx, authsqlc.MarkUserEmailVerifiedParams{
		ID:              userID,
		EmailVerifiedAt: nullableTime(&verifiedAt),
	})
}

func (repo *PostgresRepository) CreatePasswordCredential(ctx context.Context, credential domain.PasswordCredential) (domain.PasswordCredential, error) {
	row, err := repo.queries.CreatePasswordCredential(ctx, authsqlc.CreatePasswordCredentialParams{
		UserID:       credential.UserID,
		PasswordHash: credential.PasswordHash,
		CreatedAt:    credential.CreatedAt,
		UpdatedAt:    credential.UpdatedAt,
	})
	if err != nil {
		return domain.PasswordCredential{}, err
	}

	return mapPasswordCredential(row), nil
}

func (repo *PostgresRepository) GetPasswordCredentialByUserID(ctx context.Context, userID uuid.UUID) (domain.PasswordCredential, error) {
	row, err := repo.queries.GetPasswordCredentialByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.PasswordCredential{}, authapp.NotFound("password credential not found", err)
		}

		return domain.PasswordCredential{}, err
	}

	return mapPasswordCredential(row), nil
}

func (repo *PostgresRepository) GetPasswordCredentialByEmail(ctx context.Context, email string) (domain.PasswordCredential, error) {
	row, err := repo.queries.GetPasswordCredentialByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.PasswordCredential{}, authapp.NotFound("password credential not found", err)
		}

		return domain.PasswordCredential{}, err
	}

	return mapPasswordCredential(row), nil
}

func (repo *PostgresRepository) UpdatePasswordHash(ctx context.Context, userID uuid.UUID, passwordHash string, updatedAt time.Time) error {
	return repo.queries.UpdatePasswordHash(ctx, authsqlc.UpdatePasswordHashParams{
		UserID:       userID,
		PasswordHash: passwordHash,
		UpdatedAt:    updatedAt,
	})
}

func (repo *PostgresRepository) CreateSession(ctx context.Context, session domain.Session) (domain.Session, error) {
	row, err := repo.queries.CreateSession(ctx, authsqlc.CreateSessionParams{
		ID:        session.ID,
		UserID:    session.UserID,
		TokenHash: session.TokenHash,
		UserAgent: session.UserAgent,
		IpAddress: session.IPAddress,
		ExpiresAt: session.ExpiresAt,
		CreatedAt: session.CreatedAt,
	})
	if err != nil {
		return domain.Session{}, err
	}

	return mapSession(row), nil
}

func (repo *PostgresRepository) GetSessionByID(ctx context.Context, id uuid.UUID) (domain.Session, error) {
	row, err := repo.queries.GetSessionByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Session{}, authapp.NotFound("session not found", err)
		}

		return domain.Session{}, err
	}

	return mapSession(row), nil
}

func (repo *PostgresRepository) GetSessionByTokenHash(ctx context.Context, tokenHash string) (domain.Session, error) {
	row, err := repo.queries.GetSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Session{}, authapp.NotFound("session not found", err)
		}

		return domain.Session{}, err
	}

	return mapSession(row), nil
}

func (repo *PostgresRepository) RevokeSession(ctx context.Context, input authapp.RevokeSessionInput) error {
	return repo.queries.RevokeSession(ctx, authsqlc.RevokeSessionParams{
		ID:                  input.ID,
		RevokedAt:           nullableTime(&input.RevokedAt),
		ReplacedBySessionID: nullableUUID(input.ReplacedBySession),
	})
}

func (repo *PostgresRepository) RevokeAllSessionsForUser(ctx context.Context, userID uuid.UUID, revokedAt time.Time) error {
	return repo.queries.RevokeAllSessionsForUser(ctx, authsqlc.RevokeAllSessionsForUserParams{
		UserID:    userID,
		RevokedAt: nullableTime(&revokedAt),
	})
}

func (repo *PostgresRepository) CreateEmailVerificationToken(ctx context.Context, token domain.EmailVerificationToken) (domain.EmailVerificationToken, error) {
	row, err := repo.queries.CreateEmailVerificationToken(ctx, authsqlc.CreateEmailVerificationTokenParams{
		ID:         token.ID,
		UserID:     token.UserID,
		TokenHash:  token.TokenHash,
		ExpiresAt:  token.ExpiresAt,
		ConsumedAt: nullableTime(token.ConsumedAt),
		CreatedAt:  token.CreatedAt,
	})
	if err != nil {
		return domain.EmailVerificationToken{}, err
	}

	return mapEmailVerificationToken(row), nil
}

func (repo *PostgresRepository) GetEmailVerificationTokenByHash(ctx context.Context, tokenHash string) (domain.EmailVerificationToken, error) {
	row, err := repo.queries.GetEmailVerificationTokenByHash(ctx, tokenHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.EmailVerificationToken{}, authapp.NotFound("email verification token not found", err)
		}

		return domain.EmailVerificationToken{}, err
	}

	return mapEmailVerificationToken(row), nil
}

func (repo *PostgresRepository) ConsumeEmailVerificationToken(ctx context.Context, tokenID uuid.UUID, consumedAt time.Time) error {
	return repo.queries.ConsumeEmailVerificationToken(ctx, authsqlc.ConsumeEmailVerificationTokenParams{
		ID:         tokenID,
		ConsumedAt: nullableTime(&consumedAt),
	})
}

func (repo *PostgresRepository) InvalidateEmailVerificationTokensForUser(ctx context.Context, userID uuid.UUID, consumedAt time.Time) error {
	return repo.queries.InvalidateEmailVerificationTokensForUser(ctx, authsqlc.InvalidateEmailVerificationTokensForUserParams{
		UserID:     userID,
		ConsumedAt: nullableTime(&consumedAt),
	})
}

func (repo *PostgresRepository) CreatePasswordResetToken(ctx context.Context, token domain.PasswordResetToken) (domain.PasswordResetToken, error) {
	row, err := repo.queries.CreatePasswordResetToken(ctx, authsqlc.CreatePasswordResetTokenParams{
		ID:         token.ID,
		UserID:     token.UserID,
		TokenHash:  token.TokenHash,
		ExpiresAt:  token.ExpiresAt,
		ConsumedAt: nullableTime(token.ConsumedAt),
		CreatedAt:  token.CreatedAt,
	})
	if err != nil {
		return domain.PasswordResetToken{}, err
	}

	return mapPasswordResetToken(row), nil
}

func (repo *PostgresRepository) GetPasswordResetTokenByHash(ctx context.Context, tokenHash string) (domain.PasswordResetToken, error) {
	row, err := repo.queries.GetPasswordResetTokenByHash(ctx, tokenHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.PasswordResetToken{}, authapp.NotFound("password reset token not found", err)
		}

		return domain.PasswordResetToken{}, err
	}

	return mapPasswordResetToken(row), nil
}

func (repo *PostgresRepository) ConsumePasswordResetToken(ctx context.Context, tokenID uuid.UUID, consumedAt time.Time) error {
	return repo.queries.ConsumePasswordResetToken(ctx, authsqlc.ConsumePasswordResetTokenParams{
		ID:         tokenID,
		ConsumedAt: nullableTime(&consumedAt),
	})
}

func (repo *PostgresRepository) InvalidatePasswordResetTokensForUser(ctx context.Context, userID uuid.UUID, consumedAt time.Time) error {
	return repo.queries.InvalidatePasswordResetTokensForUser(ctx, authsqlc.InvalidatePasswordResetTokensForUserParams{
		UserID:     userID,
		ConsumedAt: nullableTime(&consumedAt),
	})
}

func (uow *PostgresUnitOfWork) WithinTransaction(ctx context.Context, fn func(ctx context.Context, repo authapp.Repository) error) error {
	if uow.txManager == nil {
		return fn(ctx, newPostgresRepository(authsqlc.New(uow.db)))
	}

	return uow.txManager.WithinTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		return fn(ctx, newPostgresRepository(authsqlc.New(tx)))
	})
}

func mapUser(row authsqlc.User) domain.User {
	return domain.User{
		ID:              row.ID,
		Email:           row.Email,
		BaseRole:        row.BaseRole,
		Status:          row.Status,
		EmailVerifiedAt: pointerTime(row.EmailVerifiedAt),
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}
}

func mapPasswordCredential(row authsqlc.AuthPasswordCredential) domain.PasswordCredential {
	return domain.PasswordCredential{
		UserID:       row.UserID,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}

func mapSession(row authsqlc.AuthSession) domain.Session {
	return domain.Session{
		ID:                row.ID,
		UserID:            row.UserID,
		TokenHash:         row.TokenHash,
		UserAgent:         row.UserAgent,
		IPAddress:         row.IpAddress,
		ExpiresAt:         row.ExpiresAt,
		RevokedAt:         pointerTime(row.RevokedAt),
		ReplacedBySession: pointerUUID(row.ReplacedBySessionID),
		LastUsedAt:        pointerTime(row.LastUsedAt),
		CreatedAt:         row.CreatedAt,
	}
}

func mapEmailVerificationToken(row authsqlc.AuthEmailVerificationToken) domain.EmailVerificationToken {
	return domain.EmailVerificationToken{
		ID:         row.ID,
		UserID:     row.UserID,
		TokenHash:  row.TokenHash,
		ExpiresAt:  row.ExpiresAt,
		ConsumedAt: pointerTime(row.ConsumedAt),
		CreatedAt:  row.CreatedAt,
	}
}

func mapPasswordResetToken(row authsqlc.AuthPasswordResetToken) domain.PasswordResetToken {
	return domain.PasswordResetToken{
		ID:         row.ID,
		UserID:     row.UserID,
		TokenHash:  row.TokenHash,
		ExpiresAt:  row.ExpiresAt,
		ConsumedAt: pointerTime(row.ConsumedAt),
		CreatedAt:  row.CreatedAt,
	}
}

func nullableTime(value *time.Time) sql.NullTime {
	if value == nil {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Time:  value.UTC(),
		Valid: true,
	}
}

func pointerTime(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}

	t := value.Time.UTC()
	return &t
}

func nullableUUID(value *uuid.UUID) uuid.NullUUID {
	if value == nil {
		return uuid.NullUUID{}
	}

	return uuid.NullUUID{
		UUID:  *value,
		Valid: true,
	}
}

func pointerUUID(value uuid.NullUUID) *uuid.UUID {
	if !value.Valid {
		return nil
	}

	id := value.UUID
	return &id
}
