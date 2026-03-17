package store

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	authapp "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/app"
	"github.com/Tokuchi61/Novascans/internal/modules/identity/auth/domain"
)

type MemoryRepository struct {
	mu                        sync.RWMutex
	users                     map[uuid.UUID]domain.User
	usersByMail               map[string]uuid.UUID
	credentials               map[uuid.UUID]domain.PasswordCredential
	sessions                  map[uuid.UUID]domain.Session
	sessionsByTokenHash       map[string]uuid.UUID
	emailVerificationTokens   map[uuid.UUID]domain.EmailVerificationToken
	emailVerificationByHash   map[string]uuid.UUID
	passwordResetTokens       map[uuid.UUID]domain.PasswordResetToken
	passwordResetTokensByHash map[string]uuid.UUID
}

type MemoryUnitOfWork struct {
	repo *MemoryRepository
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users:                     make(map[uuid.UUID]domain.User),
		usersByMail:               make(map[string]uuid.UUID),
		credentials:               make(map[uuid.UUID]domain.PasswordCredential),
		sessions:                  make(map[uuid.UUID]domain.Session),
		sessionsByTokenHash:       make(map[string]uuid.UUID),
		emailVerificationTokens:   make(map[uuid.UUID]domain.EmailVerificationToken),
		emailVerificationByHash:   make(map[string]uuid.UUID),
		passwordResetTokens:       make(map[uuid.UUID]domain.PasswordResetToken),
		passwordResetTokensByHash: make(map[string]uuid.UUID),
	}
}

func NewMemoryUnitOfWork(repo *MemoryRepository) *MemoryUnitOfWork {
	return &MemoryUnitOfWork{repo: repo}
}

func (repo *MemoryRepository) CreateUser(_ context.Context, user domain.User) (domain.User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.users[user.ID] = user
	repo.usersByMail[user.Email] = user.ID
	return user, nil
}

func (repo *MemoryRepository) GetUserByID(_ context.Context, id uuid.UUID) (domain.User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	user, ok := repo.users[id]
	if !ok {
		return domain.User{}, authapp.NotFound("user not found", nil)
	}

	return user, nil
}

func (repo *MemoryRepository) GetUserByEmail(_ context.Context, email string) (domain.User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	userID, ok := repo.usersByMail[email]
	if !ok {
		return domain.User{}, authapp.NotFound("user not found", nil)
	}

	user, ok := repo.users[userID]
	if !ok {
		return domain.User{}, authapp.NotFound("user not found", nil)
	}

	return user, nil
}

func (repo *MemoryRepository) MarkUserEmailVerified(_ context.Context, userID uuid.UUID, verifiedAt time.Time) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	user, ok := repo.users[userID]
	if !ok {
		return authapp.NotFound("user not found", nil)
	}

	verifiedAt = verifiedAt.UTC()
	user.Status = domain.StatusActive
	user.EmailVerifiedAt = &verifiedAt
	user.UpdatedAt = verifiedAt
	repo.users[userID] = user
	return nil
}

func (repo *MemoryRepository) CreatePasswordCredential(_ context.Context, credential domain.PasswordCredential) (domain.PasswordCredential, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.credentials[credential.UserID] = credential
	return credential, nil
}

func (repo *MemoryRepository) GetPasswordCredentialByUserID(_ context.Context, userID uuid.UUID) (domain.PasswordCredential, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	credential, ok := repo.credentials[userID]
	if !ok {
		return domain.PasswordCredential{}, authapp.NotFound("password credential not found", nil)
	}

	return credential, nil
}

func (repo *MemoryRepository) GetPasswordCredentialByEmail(ctx context.Context, email string) (domain.PasswordCredential, error) {
	user, err := repo.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.PasswordCredential{}, err
	}

	return repo.GetPasswordCredentialByUserID(ctx, user.ID)
}

func (repo *MemoryRepository) UpdatePasswordHash(_ context.Context, userID uuid.UUID, passwordHash string, updatedAt time.Time) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	credential, ok := repo.credentials[userID]
	if !ok {
		return authapp.NotFound("password credential not found", nil)
	}

	credential.PasswordHash = passwordHash
	credential.UpdatedAt = updatedAt.UTC()
	repo.credentials[userID] = credential
	return nil
}

func (repo *MemoryRepository) CreateSession(_ context.Context, session domain.Session) (domain.Session, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.sessions[session.ID] = session
	repo.sessionsByTokenHash[session.TokenHash] = session.ID
	return session, nil
}

func (repo *MemoryRepository) GetSessionByID(_ context.Context, id uuid.UUID) (domain.Session, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	session, ok := repo.sessions[id]
	if !ok {
		return domain.Session{}, authapp.NotFound("session not found", nil)
	}

	return session, nil
}

func (repo *MemoryRepository) GetSessionByTokenHash(_ context.Context, tokenHash string) (domain.Session, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	sessionID, ok := repo.sessionsByTokenHash[tokenHash]
	if !ok {
		return domain.Session{}, authapp.NotFound("session not found", nil)
	}

	session, ok := repo.sessions[sessionID]
	if !ok {
		return domain.Session{}, authapp.NotFound("session not found", nil)
	}

	return session, nil
}

func (repo *MemoryRepository) RevokeSession(_ context.Context, input authapp.RevokeSessionInput) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	session, ok := repo.sessions[input.ID]
	if !ok {
		return authapp.NotFound("session not found", nil)
	}

	revokedAt := input.RevokedAt.UTC()
	session.RevokedAt = &revokedAt
	session.ReplacedBySession = input.ReplacedBySession
	repo.sessions[input.ID] = session
	return nil
}

func (repo *MemoryRepository) RevokeAllSessionsForUser(_ context.Context, userID uuid.UUID, revokedAt time.Time) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	revokedAt = revokedAt.UTC()
	for id, session := range repo.sessions {
		if session.UserID != userID || session.RevokedAt != nil {
			continue
		}

		session.RevokedAt = &revokedAt
		repo.sessions[id] = session
	}

	return nil
}

func (repo *MemoryRepository) CreateEmailVerificationToken(_ context.Context, token domain.EmailVerificationToken) (domain.EmailVerificationToken, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.emailVerificationTokens[token.ID] = token
	repo.emailVerificationByHash[token.TokenHash] = token.ID
	return token, nil
}

func (repo *MemoryRepository) GetEmailVerificationTokenByHash(_ context.Context, tokenHash string) (domain.EmailVerificationToken, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	tokenID, ok := repo.emailVerificationByHash[tokenHash]
	if !ok {
		return domain.EmailVerificationToken{}, authapp.NotFound("email verification token not found", nil)
	}

	token, ok := repo.emailVerificationTokens[tokenID]
	if !ok {
		return domain.EmailVerificationToken{}, authapp.NotFound("email verification token not found", nil)
	}

	return token, nil
}

func (repo *MemoryRepository) ConsumeEmailVerificationToken(_ context.Context, tokenID uuid.UUID, consumedAt time.Time) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	token, ok := repo.emailVerificationTokens[tokenID]
	if !ok {
		return authapp.NotFound("email verification token not found", nil)
	}

	consumedAt = consumedAt.UTC()
	token.ConsumedAt = &consumedAt
	repo.emailVerificationTokens[tokenID] = token
	return nil
}

func (repo *MemoryRepository) InvalidateEmailVerificationTokensForUser(_ context.Context, userID uuid.UUID, consumedAt time.Time) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	consumedAt = consumedAt.UTC()
	for id, token := range repo.emailVerificationTokens {
		if token.UserID != userID || token.ConsumedAt != nil {
			continue
		}

		token.ConsumedAt = &consumedAt
		repo.emailVerificationTokens[id] = token
	}

	return nil
}

func (repo *MemoryRepository) CreatePasswordResetToken(_ context.Context, token domain.PasswordResetToken) (domain.PasswordResetToken, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.passwordResetTokens[token.ID] = token
	repo.passwordResetTokensByHash[token.TokenHash] = token.ID
	return token, nil
}

func (repo *MemoryRepository) GetPasswordResetTokenByHash(_ context.Context, tokenHash string) (domain.PasswordResetToken, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	tokenID, ok := repo.passwordResetTokensByHash[tokenHash]
	if !ok {
		return domain.PasswordResetToken{}, authapp.NotFound("password reset token not found", nil)
	}

	token, ok := repo.passwordResetTokens[tokenID]
	if !ok {
		return domain.PasswordResetToken{}, authapp.NotFound("password reset token not found", nil)
	}

	return token, nil
}

func (repo *MemoryRepository) ConsumePasswordResetToken(_ context.Context, tokenID uuid.UUID, consumedAt time.Time) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	token, ok := repo.passwordResetTokens[tokenID]
	if !ok {
		return authapp.NotFound("password reset token not found", nil)
	}

	consumedAt = consumedAt.UTC()
	token.ConsumedAt = &consumedAt
	repo.passwordResetTokens[tokenID] = token
	return nil
}

func (repo *MemoryRepository) InvalidatePasswordResetTokensForUser(_ context.Context, userID uuid.UUID, consumedAt time.Time) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	consumedAt = consumedAt.UTC()
	for id, token := range repo.passwordResetTokens {
		if token.UserID != userID || token.ConsumedAt != nil {
			continue
		}

		token.ConsumedAt = &consumedAt
		repo.passwordResetTokens[id] = token
	}

	return nil
}

func (uow *MemoryUnitOfWork) WithinTransaction(ctx context.Context, fn func(ctx context.Context, repo authapp.Repository) error) error {
	return fn(ctx, uow.repo)
}
