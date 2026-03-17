package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadFromLookupParsesAllRequiredValues(t *testing.T) {
	t.Setenv("NOVASCANS_APP_ENV", "test")

	values := map[string]string{
		"NOVASCANS_APP_ENV":                           "development",
		"NOVASCANS_HTTP_HOST":                         "0.0.0.0",
		"NOVASCANS_HTTP_PORT":                         "8080",
		"NOVASCANS_HTTP_READ_TIMEOUT":                 "15s",
		"NOVASCANS_HTTP_WRITE_TIMEOUT":                "15s",
		"NOVASCANS_HTTP_IDLE_TIMEOUT":                 "60s",
		"NOVASCANS_HTTP_SHUTDOWN_TIMEOUT":             "10s",
		"NOVASCANS_DB_HOST":                           "postgres",
		"NOVASCANS_DB_PORT":                           "5432",
		"NOVASCANS_DB_NAME":                           "novascans",
		"NOVASCANS_DB_TEST_NAME":                      "novascans_test",
		"NOVASCANS_DB_USER":                           "postgres",
		"NOVASCANS_DB_PASSWORD":                       "postgres",
		"NOVASCANS_DB_SSLMODE":                        "disable",
		"NOVASCANS_DB_MAX_OPEN_CONNS":                 "25",
		"NOVASCANS_DB_MAX_IDLE_CONNS":                 "25",
		"NOVASCANS_DB_CONN_MAX_LIFETIME":              "5m",
		"NOVASCANS_AUTH_ACCESS_TOKEN_SECRET":          "development-secret",
		"NOVASCANS_AUTH_ACCESS_TOKEN_TTL":             "15m",
		"NOVASCANS_AUTH_REFRESH_TOKEN_TTL":            "720h",
		"NOVASCANS_AUTH_EMAIL_VERIFICATION_TOKEN_TTL": "24h",
		"NOVASCANS_AUTH_PASSWORD_RESET_TOKEN_TTL":     "1h",
		"NOVASCANS_LOG_LEVEL":                         "debug",
		"NOVASCANS_METRICS_ENABLED":                   "true",
	}

	cfg, err := load(func(key string) (string, bool) {
		value, ok := values[key]
		return value, ok
	})
	if err != nil {
		t.Fatalf("load returned error: %v", err)
	}

	if cfg.HTTP.Port != 8080 {
		t.Fatalf("expected HTTP port 8080, got %d", cfg.HTTP.Port)
	}

	if cfg.DB.ConnMaxLifetime != 5*time.Minute {
		t.Fatalf("expected DB conn max lifetime 5m, got %s", cfg.DB.ConnMaxLifetime)
	}

	if cfg.Auth.AccessTokenTTL != 15*time.Minute {
		t.Fatalf("expected access token TTL 15m, got %s", cfg.Auth.AccessTokenTTL)
	}

	if !cfg.Metrics.Enabled {
		t.Fatal("expected metrics to be enabled")
	}
}

func TestLoadFromLookupFailsFastOnMissingValues(t *testing.T) {
	_, err := load(func(key string) (string, bool) {
		return "", false
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "NOVASCANS_APP_ENV is required") {
		t.Fatalf("expected missing env error, got %q", err.Error())
	}
}

func TestLoadEnvFileLoadsValuesWithoutOverridingExistingEnvironment(t *testing.T) {
	tempDir := t.TempDir()
	envPath := filepath.Join(tempDir, ".env")
	content := strings.Join([]string{
		"# comment",
		"NOVASCANS_APP_ENV=development",
		"NOVASCANS_DB_HOST=localhost",
		"NOVASCANS_DB_PASSWORD=\"postgres\"",
		"",
	}, "\n")

	if err := os.WriteFile(envPath, []byte(content), 0o600); err != nil {
		t.Fatalf("write env file: %v", err)
	}

	t.Setenv("NOVASCANS_APP_ENV", "test")

	if err := LoadEnvFile(envPath); err != nil {
		t.Fatalf("load env file: %v", err)
	}

	if got := os.Getenv("NOVASCANS_APP_ENV"); got != "test" {
		t.Fatalf("expected existing env to stay untouched, got %q", got)
	}

	if got := os.Getenv("NOVASCANS_DB_HOST"); got != "localhost" {
		t.Fatalf("expected DB host to be loaded, got %q", got)
	}

	if got := os.Getenv("NOVASCANS_DB_PASSWORD"); got != "postgres" {
		t.Fatalf("expected quoted value to be unwrapped, got %q", got)
	}
}

func TestLoadEnvFileReturnsNotExistForMissingFile(t *testing.T) {
	err := LoadEnvFile(filepath.Join(t.TempDir(), ".env"))
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected not exist error, got %v", err)
	}
}
