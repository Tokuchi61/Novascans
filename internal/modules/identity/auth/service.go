package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Tokuchi61/Novascans/internal/modules/identity/auth/store"
	platformdb "github.com/Tokuchi61/Novascans/internal/platform/db"
	"github.com/Tokuchi61/Novascans/internal/platform/events"
	platformhttp "github.com/Tokuchi61/Novascans/internal/platform/http"
)

type Service struct {
	repo      store.Repository
	txManager platformdb.TxManager
	events    events.Bus
}

type PingData struct {
	Module      string `json:"module"`
	Status      string `json:"status"`
	GeneratedAt string `json:"generated_at"`
}

type CreateUserInput struct {
	Email    string
	Password string
}

type CreateSessionInput struct {
	Email     string
	Password  string
	UserAgent string
	IPAddress string
}

type SessionOutput struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt sql.NullTime
}

func NewService(repo store.Repository, txManager platformdb.TxManager, bus events.Bus) *Service {
	return &Service{
		repo:      repo,
		txManager: txManager,
		events:    bus,
	}
}

func (service *Service) Ping(_ context.Context) PingData {
	return PingData{
		Module:      "identity.auth",
		Status:      "ok",
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func (service *Service) CreateUser(ctx context.Context, input CreateUserInput) (store.User, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	if _, err := service.repo.GetUserByEmail(ctx, email); err == nil {
		return store.User{}, platformhttp.Conflict("user already exists", nil)
	} else if !errors.Is(err, sql.ErrNoRows) {
		return store.User{}, platformhttp.Internal("failed to check existing user", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return store.User{}, platformhttp.Internal("failed to hash password", err)
	}

	now := time.Now().UTC()
	userID := newIdentifier()
	var createdUser store.User

	err = service.withWriteStore(ctx, func(repo store.Repository) error {
		user, err := repo.CreateUser(ctx, store.CreateUserParams{
			ID:              userID,
			Email:           email,
			Status:          "active",
			EmailVerifiedAt: sql.NullTime{},
			CreatedAt:       now,
			UpdatedAt:       now,
		})
		if err != nil {
			return platformhttp.Internal("failed to create user", err)
		}

		if _, err := repo.CreatePasswordCredential(ctx, store.CreatePasswordCredentialParams{
			UserID:       userID,
			PasswordHash: string(passwordHash),
			CreatedAt:    now,
			UpdatedAt:    now,
		}); err != nil {
			return platformhttp.Internal("failed to create password credential", err)
		}

		createdUser = user
		return nil
	})
	if err != nil {
		return store.User{}, err
	}

	if service.events != nil {
		_ = service.events.Publish(ctx, events.Event{
			Name:    "identity.auth.user_created",
			Payload: map[string]string{"user_id": createdUser.ID},
		})
	}

	return createdUser, nil
}

func (service *Service) GetUserByID(ctx context.Context, id string) (store.User, error) {
	user, err := service.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return store.User{}, platformhttp.NotFound("user not found", err)
		}

		return store.User{}, platformhttp.Internal("failed to fetch user", err)
	}

	return user, nil
}

func (service *Service) CreateSession(ctx context.Context, input CreateSessionInput) (SessionOutput, error) {
	credential, err := service.repo.GetPasswordCredentialByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return SessionOutput{}, platformhttp.Unauthorized("invalid credentials")
		}

		return SessionOutput{}, platformhttp.Internal("failed to fetch credential", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(credential.PasswordHash), []byte(input.Password)); err != nil {
		return SessionOutput{}, platformhttp.Unauthorized("invalid credentials")
	}

	now := time.Now().UTC()
	token := newIdentifier()
	tokenHash := sha256.Sum256([]byte(token))

	session, err := service.repo.CreateSession(ctx, store.CreateSessionParams{
		ID:        newIdentifier(),
		UserID:    credential.UserID,
		TokenHash: hex.EncodeToString(tokenHash[:]),
		UserAgent: input.UserAgent,
		IPAddress: input.IPAddress,
		ExpiresAt: now.Add(24 * time.Hour),
		CreatedAt: now,
	})
	if err != nil {
		return SessionOutput{}, platformhttp.Internal("failed to create session", err)
	}

	return SessionOutput{
		ID:        session.ID,
		UserID:    session.UserID,
		Token:     token,
		ExpiresAt: session.ExpiresAt,
		CreatedAt: session.CreatedAt,
		RevokedAt: session.RevokedAt,
	}, nil
}

func (service *Service) RevokeSession(ctx context.Context, id string) error {
	if _, err := service.repo.GetSessionByID(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return platformhttp.NotFound("session not found", err)
		}

		return platformhttp.Internal("failed to fetch session", err)
	}

	if err := service.repo.RevokeSession(ctx, id, time.Now().UTC()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return platformhttp.NotFound("session not found", err)
		}

		return platformhttp.Internal("failed to revoke session", err)
	}

	return nil
}

func (service *Service) withWriteStore(ctx context.Context, fn func(repo store.Repository) error) error {
	if service.txManager == nil {
		return fn(service.repo)
	}

	return service.txManager.WithinTransaction(ctx, func(_ context.Context, tx *sql.Tx) error {
		return fn(service.repo.WithTx(tx))
	})
}

func newIdentifier() string {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}

	return hex.EncodeToString(raw[:])
}
