package store

import (
	"context"
	"database/sql"
	"time"
)

type Repository interface {
	CreateUser(ctx context.Context, params CreateUserParams) (User, error)
	GetUserByID(ctx context.Context, id string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	CreatePasswordCredential(ctx context.Context, params CreatePasswordCredentialParams) (PasswordCredential, error)
	GetPasswordCredentialByEmail(ctx context.Context, email string) (PasswordCredential, error)
	CreateSession(ctx context.Context, params CreateSessionParams) (Session, error)
	GetSessionByID(ctx context.Context, id string) (Session, error)
	RevokeSession(ctx context.Context, id string, revokedAt time.Time) error
	WithTx(tx *sql.Tx) Repository
}

type User struct {
	ID              string
	Email           string
	Status          string
	EmailVerifiedAt sql.NullTime
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type PasswordCredential struct {
	UserID       string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Session struct {
	ID        string
	UserID    string
	TokenHash string
	UserAgent string
	IPAddress string
	ExpiresAt time.Time
	RevokedAt sql.NullTime
	CreatedAt time.Time
}

type CreateUserParams struct {
	ID              string
	Email           string
	Status          string
	EmailVerifiedAt sql.NullTime
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type CreatePasswordCredentialParams struct {
	UserID       string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateSessionParams struct {
	ID        string
	UserID    string
	TokenHash string
	UserAgent string
	IPAddress string
	ExpiresAt time.Time
	CreatedAt time.Time
}
