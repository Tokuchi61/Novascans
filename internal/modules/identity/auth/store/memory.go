package store

import (
	"context"
	"database/sql"
	"sync"
	"time"
)

type MemoryRepository struct {
	mu          sync.RWMutex
	users       map[string]User
	usersByMail map[string]string
	credentials map[string]PasswordCredential
	sessions    map[string]Session
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users:       make(map[string]User),
		usersByMail: make(map[string]string),
		credentials: make(map[string]PasswordCredential),
		sessions:    make(map[string]Session),
	}
}

func (repo *MemoryRepository) CreateUser(_ context.Context, params CreateUserParams) (User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	user := User{
		ID:              params.ID,
		Email:           params.Email,
		Status:          params.Status,
		EmailVerifiedAt: params.EmailVerifiedAt,
		CreatedAt:       params.CreatedAt,
		UpdatedAt:       params.UpdatedAt,
	}

	repo.users[user.ID] = user
	repo.usersByMail[user.Email] = user.ID

	return user, nil
}

func (repo *MemoryRepository) GetUserByID(_ context.Context, id string) (User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	user, ok := repo.users[id]
	if !ok {
		return User{}, sql.ErrNoRows
	}

	return user, nil
}

func (repo *MemoryRepository) GetUserByEmail(_ context.Context, email string) (User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	userID, ok := repo.usersByMail[email]
	if !ok {
		return User{}, sql.ErrNoRows
	}

	user, ok := repo.users[userID]
	if !ok {
		return User{}, sql.ErrNoRows
	}

	return user, nil
}

func (repo *MemoryRepository) CreatePasswordCredential(_ context.Context, params CreatePasswordCredentialParams) (PasswordCredential, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	credential := PasswordCredential{
		UserID:       params.UserID,
		PasswordHash: params.PasswordHash,
		CreatedAt:    params.CreatedAt,
		UpdatedAt:    params.UpdatedAt,
	}

	repo.credentials[credential.UserID] = credential
	return credential, nil
}

func (repo *MemoryRepository) GetPasswordCredentialByEmail(ctx context.Context, email string) (PasswordCredential, error) {
	user, err := repo.GetUserByEmail(ctx, email)
	if err != nil {
		return PasswordCredential{}, err
	}

	repo.mu.RLock()
	defer repo.mu.RUnlock()

	credential, ok := repo.credentials[user.ID]
	if !ok {
		return PasswordCredential{}, sql.ErrNoRows
	}

	return credential, nil
}

func (repo *MemoryRepository) CreateSession(_ context.Context, params CreateSessionParams) (Session, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	session := Session{
		ID:        params.ID,
		UserID:    params.UserID,
		TokenHash: params.TokenHash,
		UserAgent: params.UserAgent,
		IPAddress: params.IPAddress,
		ExpiresAt: params.ExpiresAt,
		RevokedAt: sql.NullTime{},
		CreatedAt: params.CreatedAt,
	}

	repo.sessions[session.ID] = session
	return session, nil
}

func (repo *MemoryRepository) GetSessionByID(_ context.Context, id string) (Session, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	session, ok := repo.sessions[id]
	if !ok {
		return Session{}, sql.ErrNoRows
	}

	return session, nil
}

func (repo *MemoryRepository) RevokeSession(_ context.Context, id string, revokedAt time.Time) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	session, ok := repo.sessions[id]
	if !ok {
		return sql.ErrNoRows
	}

	session.RevokedAt = sql.NullTime{
		Time:  revokedAt,
		Valid: true,
	}
	repo.sessions[id] = session

	return nil
}

func (repo *MemoryRepository) WithTx(_ *sql.Tx) Repository {
	return repo
}
