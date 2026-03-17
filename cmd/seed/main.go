package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"

	accessapp "github.com/Tokuchi61/Novascans/internal/modules/identity/access/app"
	accessdomain "github.com/Tokuchi61/Novascans/internal/modules/identity/access/domain"
	accessstore "github.com/Tokuchi61/Novascans/internal/modules/identity/access/store"
	authapp "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/app"
	authdomain "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/domain"
	authstore "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/store"
	"github.com/Tokuchi61/Novascans/internal/platform/config"
)

const seedPassword = "SeedPass123!"

type seededUser struct {
	email      string
	baseRole   string
	subRoleKey string
}

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		slog.New(slog.NewTextHandler(os.Stderr, nil)).Error("config load failed", "error", err)
		os.Exit(1)
	}

	if err := run(context.Background(), cfg); err != nil {
		slog.New(slog.NewTextHandler(os.Stderr, nil)).Error("seed failed", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, cfg config.Config) error {
	db, err := sql.Open("pgx", cfg.DB.ConnectionString(cfg.DB.Name))
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	authRepo := authstore.NewPostgresRepository(db)
	accessRepo := accessstore.NewPostgresRepository(db)

	permissions := []accessdomain.Permission{
		newPermission("manga.create", "Create manga entries"),
		newPermission("manga.update", "Update manga entries"),
		newPermission("manga.delete", "Delete manga entries"),
		newPermission("comment.moderate", "Moderate comment entries"),
		newPermission("chapter.create", "Create chapter entries"),
		newPermission("chapter.update", "Update chapter entries"),
		newPermission("chapter.delete", "Delete chapter entries"),
	}

	for _, permission := range permissions {
		if err := ensurePermission(ctx, accessRepo, permission); err != nil {
			return err
		}
	}

	subRoles := []accessdomain.SubRole{
		newSubRole("manga_moderator", "Manga Moderator", "Manage manga catalog", []string{"manga.create", "manga.update", "manga.delete"}),
		newSubRole("comment_moderator", "Comment Moderator", "Moderate comments", []string{"comment.moderate"}),
		newSubRole("chapter_moderator", "Chapter Moderator", "Manage chapter catalog", []string{"chapter.create", "chapter.update", "chapter.delete"}),
	}

	createdSubRoles := make(map[string]accessdomain.SubRole, len(subRoles))
	for _, subRole := range subRoles {
		role, err := ensureSubRole(ctx, accessRepo, subRole)
		if err != nil {
			return err
		}

		createdSubRoles[subRole.Key] = role
	}

	users := []seededUser{
		{email: "user.seed@novascans.local", baseRole: "user"},
		{email: "manga.mod@novascans.local", baseRole: "moderator", subRoleKey: "manga_moderator"},
		{email: "comment.mod@novascans.local", baseRole: "moderator", subRoleKey: "comment_moderator"},
		{email: "chapter.mod@novascans.local", baseRole: "moderator", subRoleKey: "chapter_moderator"},
		{email: "admin.seed@novascans.local", baseRole: "admin"},
	}

	for _, userSeed := range users {
		user, err := ensureUser(ctx, authRepo, accessRepo, userSeed.email, seedPassword, userSeed.baseRole)
		if err != nil {
			return err
		}

		if userSeed.subRoleKey != "" {
			role := createdSubRoles[userSeed.subRoleKey]
			if err := accessRepo.AssignSubRoleToUser(ctx, user.ID, role.ID, time.Now().UTC()); err != nil {
				return fmt.Errorf("assign sub role %s to %s: %w", role.Key, user.Email, err)
			}
		}
	}

	slog.Info("seed completed",
		"seed_password", seedPassword,
		"user_count", len(users),
		"sub_role_count", len(subRoles),
	)

	return nil
}

func ensureUser(ctx context.Context, authRepo *authstore.PostgresRepository, accessRepo *accessstore.PostgresRepository, email string, password string, baseRole string) (authdomain.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	now := time.Now().UTC()

	user, err := authRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if !authapp.HasCode(err, authapp.CodeNotFound) {
			return authdomain.User{}, fmt.Errorf("get user %s: %w", email, err)
		}

		verifiedAt := now
		user = authdomain.User{
			ID:              uuid.New(),
			Email:           email,
			BaseRole:        baseRole,
			Status:          authapp.StatusActive,
			EmailVerifiedAt: &verifiedAt,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		if _, err := authRepo.CreateUser(ctx, user); err != nil {
			return authdomain.User{}, fmt.Errorf("create user %s: %w", email, err)
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return authdomain.User{}, fmt.Errorf("hash password for %s: %w", email, err)
	}

	if _, err := authRepo.GetPasswordCredentialByUserID(ctx, user.ID); err != nil {
		if !authapp.HasCode(err, authapp.CodeNotFound) {
			return authdomain.User{}, fmt.Errorf("get credential for %s: %w", email, err)
		}

		if _, err := authRepo.CreatePasswordCredential(ctx, authdomain.PasswordCredential{
			UserID:       user.ID,
			PasswordHash: string(passwordHash),
			CreatedAt:    now,
			UpdatedAt:    now,
		}); err != nil {
			return authdomain.User{}, fmt.Errorf("create credential for %s: %w", email, err)
		}
	} else if err := authRepo.UpdatePasswordHash(ctx, user.ID, string(passwordHash), now); err != nil {
		return authdomain.User{}, fmt.Errorf("update password hash for %s: %w", email, err)
	}

	if user.BaseRole != baseRole {
		if err := accessRepo.UpdateUserBaseRole(ctx, user.ID, baseRole, now); err != nil {
			return authdomain.User{}, fmt.Errorf("update base role for %s: %w", email, err)
		}
		user.BaseRole = baseRole
	}

	if user.EmailVerifiedAt == nil || user.Status != authapp.StatusActive {
		if err := authRepo.MarkUserEmailVerified(ctx, user.ID, now); err != nil {
			return authdomain.User{}, fmt.Errorf("mark user verified for %s: %w", email, err)
		}
		user.Status = authapp.StatusActive
		user.EmailVerifiedAt = &now
	}

	return user, nil
}

func ensurePermission(ctx context.Context, repo *accessstore.PostgresRepository, permission accessdomain.Permission) error {
	if _, err := repo.GetPermissionByKey(ctx, permission.Key); err == nil {
		return nil
	} else if !accessapp.HasCode(err, accessapp.CodeNotFound) {
		return fmt.Errorf("get permission %s: %w", permission.Key, err)
	}

	if _, err := repo.CreatePermission(ctx, permission); err != nil {
		return fmt.Errorf("create permission %s: %w", permission.Key, err)
	}

	return nil
}

func ensureSubRole(ctx context.Context, repo *accessstore.PostgresRepository, subRole accessdomain.SubRole) (accessdomain.SubRole, error) {
	role, err := repo.GetSubRoleByKey(ctx, subRole.Key)
	if err != nil {
		if !accessapp.HasCode(err, accessapp.CodeNotFound) {
			return accessdomain.SubRole{}, fmt.Errorf("get sub role %s: %w", subRole.Key, err)
		}

		role, err = repo.CreateSubRole(ctx, subRole)
		if err != nil {
			return accessdomain.SubRole{}, fmt.Errorf("create sub role %s: %w", subRole.Key, err)
		}
	}

	for _, permissionKey := range subRole.PermissionKeys {
		permission, err := repo.GetPermissionByKey(ctx, permissionKey)
		if err != nil {
			return accessdomain.SubRole{}, fmt.Errorf("get permission %s for sub role %s: %w", permissionKey, subRole.Key, err)
		}

		if err := repo.AttachPermissionToSubRole(ctx, role.ID, permission.ID, time.Now().UTC()); err != nil {
			return accessdomain.SubRole{}, fmt.Errorf("attach permission %s to sub role %s: %w", permissionKey, subRole.Key, err)
		}
	}

	role.PermissionKeys = subRole.PermissionKeys
	return role, nil
}

func newPermission(key string, description string) accessdomain.Permission {
	now := time.Now().UTC()
	return accessdomain.Permission{
		ID:          uuid.New(),
		Key:         key,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func newSubRole(key string, name string, description string, permissionKeys []string) accessdomain.SubRole {
	now := time.Now().UTC()
	return accessdomain.SubRole{
		ID:             uuid.New(),
		Key:            key,
		Name:           name,
		Description:    description,
		PermissionKeys: permissionKeys,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
