package app

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Tokuchi61/Novascans/internal/modules/identity/auth/domain"
)

type Repository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	MarkUserEmailVerified(ctx context.Context, userID uuid.UUID, verifiedAt time.Time) error
	CreatePasswordCredential(ctx context.Context, credential domain.PasswordCredential) (domain.PasswordCredential, error)
	GetPasswordCredentialByUserID(ctx context.Context, userID uuid.UUID) (domain.PasswordCredential, error)
	GetPasswordCredentialByEmail(ctx context.Context, email string) (domain.PasswordCredential, error)
	UpdatePasswordHash(ctx context.Context, userID uuid.UUID, passwordHash string, updatedAt time.Time) error
	CreateSession(ctx context.Context, session domain.Session) (domain.Session, error)
	GetSessionByID(ctx context.Context, id uuid.UUID) (domain.Session, error)
	GetSessionByTokenHash(ctx context.Context, tokenHash string) (domain.Session, error)
	RevokeSession(ctx context.Context, input RevokeSessionInput) error
	RevokeAllSessionsForUser(ctx context.Context, userID uuid.UUID, revokedAt time.Time) error
	CreateEmailVerificationToken(ctx context.Context, token domain.EmailVerificationToken) (domain.EmailVerificationToken, error)
	GetEmailVerificationTokenByHash(ctx context.Context, tokenHash string) (domain.EmailVerificationToken, error)
	ConsumeEmailVerificationToken(ctx context.Context, tokenID uuid.UUID, consumedAt time.Time) error
	InvalidateEmailVerificationTokensForUser(ctx context.Context, userID uuid.UUID, consumedAt time.Time) error
	CreatePasswordResetToken(ctx context.Context, token domain.PasswordResetToken) (domain.PasswordResetToken, error)
	GetPasswordResetTokenByHash(ctx context.Context, tokenHash string) (domain.PasswordResetToken, error)
	ConsumePasswordResetToken(ctx context.Context, tokenID uuid.UUID, consumedAt time.Time) error
	InvalidatePasswordResetTokensForUser(ctx context.Context, userID uuid.UUID, consumedAt time.Time) error
}

type UnitOfWork interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context, repo Repository) error) error
}

type AccountProvisioner interface {
	ProvisionDefaults(ctx context.Context, user domain.User, now time.Time) error
}

type Authenticator interface {
	AuthenticateAccessToken(ctx context.Context, token string) (AuthenticatedUser, error)
}

type RevokeSessionInput struct {
	ID                uuid.UUID
	RevokedAt         time.Time
	ReplacedBySession *uuid.UUID
}
