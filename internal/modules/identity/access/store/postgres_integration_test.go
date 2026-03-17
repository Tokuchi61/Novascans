//go:build integration

package store

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Tokuchi61/Novascans/internal/modules/identity/access/domain"
	authapp "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/app"
	authdomain "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/domain"
	authstore "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/store"
	"github.com/Tokuchi61/Novascans/internal/platform/config"
)

func TestPostgresRepositoryResolvesUserAccess(t *testing.T) {
	databaseURL := os.Getenv("NOVASCANS_TEST_DATABASE_URL")
	if err := config.LoadDefaultEnvFile(); err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("load env file: %v", err)
	}

	if databaseURL == "" {
		cfg, err := config.LoadFromEnv()
		if err != nil {
			t.Skipf("integration config is not ready: %v", err)
		}

		if err := ensureTestDatabase(t, cfg.DB); err != nil {
			t.Fatalf("ensure test database: %v", err)
		}

		databaseURL = cfg.DB.ConnectionString(cfg.DB.TestName)
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := db.PingContext(t.Context()); err != nil {
		t.Fatalf("ping db: %v", err)
	}

	applySchemas(t, db)
	resetAccessTables(t, db)

	authRepo := authstore.NewPostgresRepository(db)
	repo := NewPostgresRepository(db)

	now := time.Now().UTC()
	userID := uuid.New()
	if _, err := authRepo.CreateUser(t.Context(), authdomain.User{
		ID:              userID,
		Email:           "moderator@example.com",
		BaseRole:        "moderator",
		Status:          authapp.StatusActive,
		EmailVerifiedAt: &now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}); err != nil {
		t.Fatalf("create user: %v", err)
	}

	permission, err := repo.CreatePermission(t.Context(), domain.Permission{
		ID:          uuid.New(),
		Key:         "comment.moderate",
		Description: "Moderate comments",
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		t.Fatalf("create permission: %v", err)
	}

	subRole, err := repo.CreateSubRole(t.Context(), domain.SubRole{
		ID:          uuid.New(),
		Key:         "comment_moderator",
		Name:        "Comment Moderator",
		Description: "Moderate comments",
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		t.Fatalf("create sub role: %v", err)
	}

	if err := repo.AttachPermissionToSubRole(t.Context(), subRole.ID, permission.ID, now); err != nil {
		t.Fatalf("attach permission to sub role: %v", err)
	}

	if err := repo.AssignSubRoleToUser(t.Context(), userID, subRole.ID, now); err != nil {
		t.Fatalf("assign sub role to user: %v", err)
	}

	permissionKeys, err := repo.ListPermissionKeysForUser(t.Context(), userID)
	if err != nil {
		t.Fatalf("list permission keys for user: %v", err)
	}

	if len(permissionKeys) != 1 || permissionKeys[0] != "comment.moderate" {
		t.Fatalf("expected comment.moderate permission, got %v", permissionKeys)
	}

	subRoles, err := repo.ListSubRolesForUser(t.Context(), userID)
	if err != nil {
		t.Fatalf("list sub roles for user: %v", err)
	}

	if len(subRoles) != 1 || subRoles[0].Key != "comment_moderator" {
		t.Fatalf("expected comment_moderator sub role, got %+v", subRoles)
	}

	if err := repo.UpdateUserBaseRole(t.Context(), userID, "admin", now); err != nil {
		t.Fatalf("update user base role: %v", err)
	}

	user, err := authRepo.GetUserByID(t.Context(), userID)
	if err != nil {
		t.Fatalf("get user by id: %v", err)
	}

	if user.BaseRole != "admin" {
		t.Fatalf("expected base role admin, got %q", user.BaseRole)
	}
}

func ensureTestDatabase(t *testing.T, dbCfg config.DBConfig) error {
	t.Helper()

	adminDB, err := sql.Open("pgx", dbCfg.ConnectionString("postgres"))
	if err != nil {
		return err
	}
	defer adminDB.Close()

	if err := adminDB.PingContext(t.Context()); err != nil {
		return err
	}

	var exists bool
	if err := adminDB.QueryRowContext(
		t.Context(),
		"SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)",
		dbCfg.TestName,
	).Scan(&exists); err != nil {
		return err
	}

	if exists {
		return nil
	}

	quotedName := `"` + strings.ReplaceAll(dbCfg.TestName, `"`, `""`) + `"`
	_, err = adminDB.ExecContext(t.Context(), "CREATE DATABASE "+quotedName)
	return err
}

func applySchemas(t *testing.T, db *sql.DB) {
	t.Helper()

	schemaDir := filepath.Clean(filepath.Join("..", "..", "..", "..", "..", "db", "schema"))
	files := []string{"identity_auth.sql", "identity_access.sql"}

	for _, file := range files {
		schema, err := os.ReadFile(filepath.Join(schemaDir, file))
		if err != nil {
			t.Fatalf("read schema %s: %v", file, err)
		}

		if _, err := db.ExecContext(t.Context(), string(schema)); err != nil {
			t.Fatalf("apply schema %s: %v", file, err)
		}
	}
}

func resetAccessTables(t *testing.T, db *sql.DB) {
	t.Helper()

	statements := []string{
		"DELETE FROM access_user_sub_roles",
		"DELETE FROM access_sub_role_permissions",
		"DELETE FROM access_sub_roles",
		"DELETE FROM access_permissions",
		"DELETE FROM auth_sessions",
		"DELETE FROM auth_password_credentials",
		"DELETE FROM users",
	}

	for _, statement := range statements {
		if _, err := db.ExecContext(t.Context(), statement); err != nil {
			t.Fatalf("reset tables with %q: %v", statement, err)
		}
	}
}
