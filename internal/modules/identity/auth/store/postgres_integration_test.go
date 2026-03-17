//go:build integration

package store

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"

	authapp "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/app"
	"github.com/Tokuchi61/Novascans/internal/modules/identity/auth/domain"
	"github.com/Tokuchi61/Novascans/internal/platform/config"
	platformdb "github.com/Tokuchi61/Novascans/internal/platform/db"
)

func TestPostgresRepositoryAndTxManager(t *testing.T) {
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

	applySchema(t, db)
	resetTables(t, db)

	repo := NewPostgresRepository(db)
	unitOfWork := NewPostgresUnitOfWork(db, platformdb.NewTxManager(db))

	now := time.Now().UTC()
	userID := uuid.New()

	err = unitOfWork.WithinTransaction(t.Context(), func(ctx context.Context, txRepo authapp.Repository) error {
		if _, err := txRepo.CreateUser(ctx, domain.User{
			ID:        userID,
			Email:     "integration@example.com",
			BaseRole:  "user",
			Status:    "active",
			CreatedAt: now,
			UpdatedAt: now,
		}); err != nil {
			return err
		}

		if _, err := txRepo.CreatePasswordCredential(ctx, domain.PasswordCredential{
			UserID:       userID,
			PasswordHash: "hash",
			CreatedAt:    now,
			UpdatedAt:    now,
		}); err != nil {
			return err
		}

		return errors.New("force rollback")
	})
	if err == nil {
		t.Fatal("expected transaction error, got nil")
	}

	if _, err := repo.GetUserByID(t.Context(), userID); !authapp.HasCode(err, authapp.CodeNotFound) {
		t.Fatalf("expected user to be rolled back, got %v", err)
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

func applySchema(t *testing.T, db *sql.DB) {
	t.Helper()

	schemaPath := filepath.Clean(filepath.Join("..", "..", "..", "..", "..", "db", "schema", "identity_auth.sql"))
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}

	if _, err := db.ExecContext(t.Context(), string(schema)); err != nil {
		t.Fatalf("apply schema: %v", err)
	}
}

func resetTables(t *testing.T, db *sql.DB) {
	t.Helper()

	statements := []string{
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
