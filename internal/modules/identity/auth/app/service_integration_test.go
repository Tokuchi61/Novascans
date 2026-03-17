//go:build integration

package app_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	authapp "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/app"
	"github.com/Tokuchi61/Novascans/internal/modules/identity/auth/domain"
	authstore "github.com/Tokuchi61/Novascans/internal/modules/identity/auth/store"
	accountapp "github.com/Tokuchi61/Novascans/internal/modules/user/account/app"
	accountstore "github.com/Tokuchi61/Novascans/internal/modules/user/account/store"
	"github.com/Tokuchi61/Novascans/internal/platform/config"
	platformdb "github.com/Tokuchi61/Novascans/internal/platform/db"
)

type failingProvisioner struct{}

func (failingProvisioner) ProvisionDefaults(context.Context, domain.User, time.Time) error {
	return errors.New("boom")
}

func TestRegisterCreatesAccountDefaults(t *testing.T) {
	db := openIntegrationDB(t)
	defer db.Close()

	resetTables(t, db)

	txManager := platformdb.NewTxManager(db)
	authRepo := authstore.NewPostgresRepository(db)
	authUow := authstore.NewPostgresUnitOfWork(db, txManager)
	accountRepo := accountstore.NewPostgresRepository(db)
	accountUow := accountstore.NewPostgresUnitOfWork(accountRepo, txManager)
	accountService := accountapp.NewService(accountRepo, accountUow)

	service := authapp.NewService(authRepo, authUow, nil, testServiceConfig(), accountService)

	result, err := service.Register(t.Context(), authapp.RegisterInput{
		Email:     "register@example.com",
		Password:  "Password123!",
		UserAgent: "integration-test",
		IPAddress: "127.0.0.1",
	})
	if err != nil {
		t.Fatalf("register user: %v", err)
	}

	profile, err := accountRepo.GetProfileByUserID(t.Context(), result.User.ID)
	if err != nil {
		t.Fatalf("get provisioned profile: %v", err)
	}

	if profile.Username == "" {
		t.Fatal("expected provisioned username")
	}

	if _, err := accountRepo.GetSettingsByUserID(t.Context(), result.User.ID); err != nil {
		t.Fatalf("get provisioned settings: %v", err)
	}

	if _, err := accountRepo.GetPrivacySettingsByUserID(t.Context(), result.User.ID); err != nil {
		t.Fatalf("get provisioned privacy settings: %v", err)
	}
}

func TestRegisterRollsBackWhenAccountProvisionFails(t *testing.T) {
	db := openIntegrationDB(t)
	defer db.Close()

	resetTables(t, db)

	txManager := platformdb.NewTxManager(db)
	authRepo := authstore.NewPostgresRepository(db)
	authUow := authstore.NewPostgresUnitOfWork(db, txManager)
	service := authapp.NewService(authRepo, authUow, nil, testServiceConfig(), failingProvisioner{})

	if _, err := service.Register(t.Context(), authapp.RegisterInput{
		Email:     "rollback@example.com",
		Password:  "Password123!",
		UserAgent: "integration-test",
		IPAddress: "127.0.0.1",
	}); err == nil {
		t.Fatal("expected register to fail when account provisioning fails")
	}

	if _, err := authRepo.GetUserByEmail(t.Context(), "rollback@example.com"); !authapp.HasCode(err, authapp.CodeNotFound) {
		t.Fatalf("expected user insert to roll back, got %v", err)
	}
}

func openIntegrationDB(t *testing.T) *sql.DB {
	t.Helper()

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

	if err := db.PingContext(t.Context()); err != nil {
		db.Close()
		t.Fatalf("ping db: %v", err)
	}

	applySchemas(t, db)
	return db
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
	files := []string{"identity_auth.sql", "user_account.sql"}

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

func resetTables(t *testing.T, db *sql.DB) {
	t.Helper()

	statements := []string{
		"DELETE FROM account_privacy_settings",
		"DELETE FROM account_settings",
		"DELETE FROM account_profiles",
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

func testServiceConfig() authapp.ServiceConfig {
	return authapp.ServiceConfig{
		AppEnv:                    "development",
		AccessTokenSecret:         "integration-secret",
		AccessTokenTTL:            15 * time.Minute,
		RefreshTokenTTL:           24 * time.Hour,
		EmailVerificationTokenTTL: time.Hour,
		PasswordResetTokenTTL:     time.Hour,
	}
}
