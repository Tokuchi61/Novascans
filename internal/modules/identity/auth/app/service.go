package app

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	accessdomain "github.com/Tokuchi61/Novascans/internal/modules/identity/access/domain"
	"github.com/Tokuchi61/Novascans/internal/modules/identity/auth/domain"
	"github.com/Tokuchi61/Novascans/internal/platform/events"
)

type Service struct {
	repo        Repository
	unitOfWork  UnitOfWork
	events      events.Bus
	config      ServiceConfig
	provisioner AccountProvisioner
}

func NewService(repo Repository, unitOfWork UnitOfWork, bus events.Bus, cfg ServiceConfig, provisioner AccountProvisioner) *Service {
	return &Service{
		repo:        repo,
		unitOfWork:  unitOfWork,
		events:      bus,
		config:      cfg,
		provisioner: provisioner,
	}
}

func (service *Service) SetAccountProvisioner(provisioner AccountProvisioner) {
	service.provisioner = provisioner
}

func (service *Service) Ping(_ context.Context) PingData {
	return NewPingData(time.Now())
}

func (service *Service) Register(ctx context.Context, input RegisterInput) (AuthResult, error) {
	email := normalizeEmail(input.Email)

	if _, err := service.repo.GetUserByEmail(ctx, email); err == nil {
		return AuthResult{}, Conflict("user already exists", nil)
	} else if !HasCode(err, CodeNotFound) {
		return AuthResult{}, Internal("failed to check existing user", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResult{}, Internal("failed to hash password", err)
	}

	now := time.Now().UTC()
	user := domain.User{
		ID:        uuid.New(),
		Email:     email,
		BaseRole:  accessdomain.BaseRoleUser,
		Status:    domain.StatusPendingVerification,
		CreatedAt: now,
		UpdatedAt: now,
	}

	credential := domain.PasswordCredential{
		UserID:       user.ID,
		PasswordHash: string(passwordHash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	var result AuthResult
	err = service.withWriteRepository(ctx, func(ctx context.Context, repo Repository) error {
		createdUser, createdSession, accessToken, refreshToken, err := service.createUserAndSession(ctx, repo, user, credential, input.UserAgent, input.IPAddress, now)
		if err != nil {
			return err
		}

		result = AuthResult{
			User:         createdUser,
			Session:      createdSession,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		return nil
	})
	if err != nil {
		return AuthResult{}, err
	}

	service.publishUserCreated(ctx, result.User.ID)
	return result, nil
}

func (service *Service) Login(ctx context.Context, input LoginInput) (AuthResult, error) {
	credential, err := service.repo.GetPasswordCredentialByEmail(ctx, normalizeEmail(input.Email))
	if err != nil {
		if HasCode(err, CodeNotFound) {
			return AuthResult{}, Unauthorized("invalid credentials", err)
		}

		return AuthResult{}, Internal("failed to fetch credential", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(credential.PasswordHash), []byte(input.Password)); err != nil {
		return AuthResult{}, Unauthorized("invalid credentials", err)
	}

	user, err := service.repo.GetUserByID(ctx, credential.UserID)
	if err != nil {
		if HasCode(err, CodeNotFound) {
			return AuthResult{}, Unauthorized("invalid credentials", err)
		}

		return AuthResult{}, Internal("failed to fetch user", err)
	}

	return service.issueAuthResult(ctx, user, input.UserAgent, input.IPAddress)
}

func (service *Service) Refresh(ctx context.Context, input RefreshSessionInput) (AuthResult, error) {
	session, err := service.repo.GetSessionByTokenHash(ctx, hashOpaqueToken(strings.TrimSpace(input.RefreshToken)))
	if err != nil {
		if HasCode(err, CodeNotFound) {
			return AuthResult{}, Unauthorized("invalid refresh token", err)
		}

		return AuthResult{}, Internal("failed to fetch session", err)
	}

	if err := validateSession(session); err != nil {
		return AuthResult{}, err
	}

	user, err := service.repo.GetUserByID(ctx, session.UserID)
	if err != nil {
		if HasCode(err, CodeNotFound) {
			return AuthResult{}, Unauthorized("session user not found", err)
		}

		return AuthResult{}, Internal("failed to fetch session user", err)
	}

	var result AuthResult
	now := time.Now().UTC()
	if err := service.withWriteRepository(ctx, func(ctx context.Context, repo Repository) error {
		rotatedSession, accessToken, refreshToken, err := service.createSessionForUser(ctx, repo, user, input.UserAgent, input.IPAddress, now)
		if err != nil {
			return err
		}

		if err := repo.RevokeSession(ctx, RevokeSessionInput{
			ID:                session.ID,
			RevokedAt:         now,
			ReplacedBySession: &rotatedSession.ID,
		}); err != nil {
			return Internal("failed to revoke refreshed session", err)
		}

		result = AuthResult{
			User:         user,
			Session:      rotatedSession,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		return nil
	}); err != nil {
		return AuthResult{}, err
	}

	return result, nil
}

func (service *Service) AuthenticateAccessToken(ctx context.Context, token string) (AuthenticatedUser, error) {
	claims, err := parseAccessToken(service.config.AccessTokenSecret, token)
	if err != nil {
		return AuthenticatedUser{}, err
	}

	session, err := service.repo.GetSessionByID(ctx, claims.SessionID)
	if err != nil {
		if HasCode(err, CodeNotFound) {
			return AuthenticatedUser{}, Unauthorized("session not found", err)
		}

		return AuthenticatedUser{}, Internal("failed to fetch session", err)
	}

	if session.UserID != claims.UserID {
		return AuthenticatedUser{}, Unauthorized("session does not match user", nil)
	}

	if err := validateSession(session); err != nil {
		return AuthenticatedUser{}, err
	}

	user, err := service.repo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		if HasCode(err, CodeNotFound) {
			return AuthenticatedUser{}, Unauthorized("user not found", err)
		}

		return AuthenticatedUser{}, Internal("failed to fetch user", err)
	}

	return AuthenticatedUser{
		User:    user,
		Session: session,
	}, nil
}

func (service *Service) LogoutCurrentSession(ctx context.Context, sessionID uuid.UUID) error {
	if err := service.repo.RevokeSession(ctx, RevokeSessionInput{
		ID:        sessionID,
		RevokedAt: time.Now().UTC(),
	}); err != nil {
		if HasCode(err, CodeNotFound) {
			return NotFound("session not found", err)
		}

		return Internal("failed to revoke session", err)
	}

	return nil
}

func (service *Service) LogoutAllSessions(ctx context.Context, userID uuid.UUID) error {
	if err := service.repo.RevokeAllSessionsForUser(ctx, userID, time.Now().UTC()); err != nil {
		return Internal("failed to revoke sessions", err)
	}

	return nil
}

func (service *Service) RequestEmailVerification(ctx context.Context, input RequestEmailVerificationInput) (RequestTokenResult, error) {
	user, err := service.repo.GetUserByEmail(ctx, normalizeEmail(input.Email))
	if err != nil {
		if HasCode(err, CodeNotFound) {
			return RequestTokenResult{}, nil
		}

		return RequestTokenResult{}, Internal("failed to fetch user", err)
	}

	rawToken := newOpaqueToken()
	now := time.Now().UTC()
	if err := service.withWriteRepository(ctx, func(ctx context.Context, repo Repository) error {
		if err := repo.InvalidateEmailVerificationTokensForUser(ctx, user.ID, now); err != nil {
			return Internal("failed to rotate verification tokens", err)
		}

		_, err := repo.CreateEmailVerificationToken(ctx, domain.EmailVerificationToken{
			ID:        uuid.New(),
			UserID:    user.ID,
			TokenHash: hashOpaqueToken(rawToken),
			ExpiresAt: now.Add(service.config.EmailVerificationTokenTTL),
			CreatedAt: now,
		})
		if err != nil {
			return Internal("failed to create verification token", err)
		}

		return nil
	}); err != nil {
		return RequestTokenResult{}, err
	}

	return service.tokenResult(rawToken), nil
}

func (service *Service) VerifyEmail(ctx context.Context, input VerifyEmailInput) error {
	token, err := service.repo.GetEmailVerificationTokenByHash(ctx, hashOpaqueToken(strings.TrimSpace(input.Token)))
	if err != nil {
		if HasCode(err, CodeNotFound) {
			return Unauthorized("invalid verification token", err)
		}

		return Internal("failed to fetch verification token", err)
	}

	if err := validateOneTimeToken(token.ExpiresAt, token.ConsumedAt); err != nil {
		return err
	}

	now := time.Now().UTC()
	return service.withWriteRepository(ctx, func(ctx context.Context, repo Repository) error {
		if err := repo.ConsumeEmailVerificationToken(ctx, token.ID, now); err != nil {
			return Internal("failed to consume verification token", err)
		}

		if err := repo.MarkUserEmailVerified(ctx, token.UserID, now); err != nil {
			return Internal("failed to mark email verified", err)
		}

		return nil
	})
}

func (service *Service) ForgotPassword(ctx context.Context, input ForgotPasswordInput) (RequestTokenResult, error) {
	user, err := service.repo.GetUserByEmail(ctx, normalizeEmail(input.Email))
	if err != nil {
		if HasCode(err, CodeNotFound) {
			return RequestTokenResult{}, nil
		}

		return RequestTokenResult{}, Internal("failed to fetch user", err)
	}

	rawToken := newOpaqueToken()
	now := time.Now().UTC()
	if err := service.withWriteRepository(ctx, func(ctx context.Context, repo Repository) error {
		if err := repo.InvalidatePasswordResetTokensForUser(ctx, user.ID, now); err != nil {
			return Internal("failed to rotate password reset tokens", err)
		}

		_, err := repo.CreatePasswordResetToken(ctx, domain.PasswordResetToken{
			ID:        uuid.New(),
			UserID:    user.ID,
			TokenHash: hashOpaqueToken(rawToken),
			ExpiresAt: now.Add(service.config.PasswordResetTokenTTL),
			CreatedAt: now,
		})
		if err != nil {
			return Internal("failed to create password reset token", err)
		}

		return nil
	}); err != nil {
		return RequestTokenResult{}, err
	}

	return service.tokenResult(rawToken), nil
}

func (service *Service) ResetPassword(ctx context.Context, input ResetPasswordInput) error {
	token, err := service.repo.GetPasswordResetTokenByHash(ctx, hashOpaqueToken(strings.TrimSpace(input.Token)))
	if err != nil {
		if HasCode(err, CodeNotFound) {
			return Unauthorized("invalid password reset token", err)
		}

		return Internal("failed to fetch password reset token", err)
	}

	if err := validateOneTimeToken(token.ExpiresAt, token.ConsumedAt); err != nil {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return Internal("failed to hash password", err)
	}

	now := time.Now().UTC()
	return service.withWriteRepository(ctx, func(ctx context.Context, repo Repository) error {
		if err := repo.UpdatePasswordHash(ctx, token.UserID, string(passwordHash), now); err != nil {
			return Internal("failed to update password hash", err)
		}

		if err := repo.ConsumePasswordResetToken(ctx, token.ID, now); err != nil {
			return Internal("failed to consume password reset token", err)
		}

		if err := repo.RevokeAllSessionsForUser(ctx, token.UserID, now); err != nil {
			return Internal("failed to revoke user sessions", err)
		}

		return nil
	})
}

func (service *Service) issueAuthResult(ctx context.Context, user domain.User, userAgent string, ipAddress string) (AuthResult, error) {
	now := time.Now().UTC()

	var result AuthResult
	if err := service.withWriteRepository(ctx, func(ctx context.Context, repo Repository) error {
		session, accessToken, refreshToken, err := service.createSessionForUser(ctx, repo, user, userAgent, ipAddress, now)
		if err != nil {
			return err
		}

		result = AuthResult{
			User:         user,
			Session:      session,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		return nil
	}); err != nil {
		return AuthResult{}, err
	}

	return result, nil
}

func (service *Service) createUserAndSession(ctx context.Context, repo Repository, user domain.User, credential domain.PasswordCredential, userAgent string, ipAddress string, now time.Time) (domain.User, domain.Session, string, string, error) {
	createdUser, err := repo.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, domain.Session{}, "", "", Internal("failed to create user", err)
	}

	if _, err := repo.CreatePasswordCredential(ctx, credential); err != nil {
		return domain.User{}, domain.Session{}, "", "", Internal("failed to create password credential", err)
	}

	if service.provisioner != nil {
		if err := service.provisioner.ProvisionDefaults(ctx, createdUser, now); err != nil {
			return domain.User{}, domain.Session{}, "", "", Internal("failed to provision account defaults", err)
		}
	}

	session, accessToken, refreshToken, err := service.createSessionForUser(ctx, repo, createdUser, userAgent, ipAddress, now)
	if err != nil {
		return domain.User{}, domain.Session{}, "", "", err
	}

	return createdUser, session, accessToken, refreshToken, nil
}

func (service *Service) createSessionForUser(ctx context.Context, repo Repository, user domain.User, userAgent string, ipAddress string, now time.Time) (domain.Session, string, string, error) {
	refreshToken := newOpaqueToken()
	session, err := repo.CreateSession(ctx, domain.Session{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: hashOpaqueToken(refreshToken),
		UserAgent: userAgent,
		IPAddress: ipAddress,
		ExpiresAt: now.Add(service.config.RefreshTokenTTL),
		CreatedAt: now,
	})
	if err != nil {
		return domain.Session{}, "", "", Internal("failed to create session", err)
	}

	accessToken, err := buildAccessToken(service.config.AccessTokenSecret, accessTokenClaims{
		UserID:    user.ID,
		SessionID: session.ID,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(service.config.AccessTokenTTL).Unix(),
	})
	if err != nil {
		return domain.Session{}, "", "", Internal("failed to sign access token", err)
	}

	return session, accessToken, refreshToken, nil
}

func (service *Service) publishUserCreated(ctx context.Context, userID uuid.UUID) {
	if service.events == nil {
		return
	}

	_ = service.events.Publish(ctx, events.Event{
		Name: "identity.auth.user_created",
		Payload: map[string]string{
			"user_id": userID.String(),
		},
	})
}

func (service *Service) withWriteRepository(ctx context.Context, fn func(ctx context.Context, repo Repository) error) error {
	if service.unitOfWork == nil {
		return fn(ctx, service.repo)
	}

	return service.unitOfWork.WithinTransaction(ctx, fn)
}

func (service *Service) tokenResult(rawToken string) RequestTokenResult {
	if !service.config.ExposeDebugTokens() {
		return RequestTokenResult{}
	}

	return RequestTokenResult{DebugToken: rawToken}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validateSession(session domain.Session) error {
	if session.RevokedAt != nil {
		return Unauthorized("session has been revoked", nil)
	}

	if time.Now().UTC().After(session.ExpiresAt.UTC()) {
		return Unauthorized("session has expired", nil)
	}

	return nil
}

func validateOneTimeToken(expiresAt time.Time, consumedAt *time.Time) error {
	if consumedAt != nil {
		return Unauthorized("token has already been consumed", nil)
	}

	if time.Now().UTC().After(expiresAt.UTC()) {
		return Unauthorized("token has expired", nil)
	}

	return nil
}
