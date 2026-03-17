package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	App     AppConfig
	HTTP    HTTPConfig
	DB      DBConfig
	Auth    AuthConfig
	Log     LogConfig
	Metrics MetricsConfig
}

type AppConfig struct {
	Name string
	Env  string
}

type HTTPConfig struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type DBConfig struct {
	Host            string
	Port            int
	Name            string
	TestName        string
	User            string
	Password        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type AuthConfig struct {
	AccessTokenSecret         string
	AccessTokenTTL            time.Duration
	RefreshTokenTTL           time.Duration
	EmailVerificationTokenTTL time.Duration
	PasswordResetTokenTTL     time.Duration
}

type LogConfig struct {
	Level string
}

type MetricsConfig struct {
	Enabled bool
}

func Default() Config {
	return Config{
		App: AppConfig{
			Name: "novascans",
			Env:  "development",
		},
		HTTP: HTTPConfig{
			Host:            "0.0.0.0",
			Port:            8080,
			ReadTimeout:     15 * time.Second,
			WriteTimeout:    15 * time.Second,
			IdleTimeout:     60 * time.Second,
			ShutdownTimeout: 10 * time.Second,
		},
		DB: DBConfig{
			Host:            "localhost",
			Port:            5432,
			Name:            "novascans",
			TestName:        "novascans_test",
			User:            "postgres",
			Password:        "postgres",
			SSLMode:         "disable",
			MaxOpenConns:    25,
			MaxIdleConns:    25,
			ConnMaxLifetime: 5 * time.Minute,
		},
		Auth: AuthConfig{
			AccessTokenSecret:         "change-me",
			AccessTokenTTL:            15 * time.Minute,
			RefreshTokenTTL:           30 * 24 * time.Hour,
			EmailVerificationTokenTTL: 24 * time.Hour,
			PasswordResetTokenTTL:     time.Hour,
		},
		Log: LogConfig{
			Level: "info",
		},
		Metrics: MetricsConfig{
			Enabled: true,
		},
	}
}

func LoadFromEnv() (Config, error) {
	if err := LoadDefaultEnvFile(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, err
	}

	return load(os.LookupEnv)
}

func load(lookup func(string) (string, bool)) (Config, error) {
	var errs []string
	cfg := Config{
		App: AppConfig{
			Name: "novascans",
		},
	}

	cfg.App.Env = requireString(lookup, "NOVASCANS_APP_ENV", &errs)
	cfg.HTTP.Host = requireString(lookup, "NOVASCANS_HTTP_HOST", &errs)
	cfg.HTTP.Port = requireInt(lookup, "NOVASCANS_HTTP_PORT", &errs)
	cfg.HTTP.ReadTimeout = requireDuration(lookup, "NOVASCANS_HTTP_READ_TIMEOUT", &errs)
	cfg.HTTP.WriteTimeout = requireDuration(lookup, "NOVASCANS_HTTP_WRITE_TIMEOUT", &errs)
	cfg.HTTP.IdleTimeout = requireDuration(lookup, "NOVASCANS_HTTP_IDLE_TIMEOUT", &errs)
	cfg.HTTP.ShutdownTimeout = requireDuration(lookup, "NOVASCANS_HTTP_SHUTDOWN_TIMEOUT", &errs)

	cfg.DB.Host = requireString(lookup, "NOVASCANS_DB_HOST", &errs)
	cfg.DB.Port = requireInt(lookup, "NOVASCANS_DB_PORT", &errs)
	cfg.DB.Name = requireString(lookup, "NOVASCANS_DB_NAME", &errs)
	cfg.DB.TestName = requireString(lookup, "NOVASCANS_DB_TEST_NAME", &errs)
	cfg.DB.User = requireString(lookup, "NOVASCANS_DB_USER", &errs)
	cfg.DB.Password = requireString(lookup, "NOVASCANS_DB_PASSWORD", &errs)
	cfg.DB.SSLMode = requireString(lookup, "NOVASCANS_DB_SSLMODE", &errs)
	cfg.DB.MaxOpenConns = requireInt(lookup, "NOVASCANS_DB_MAX_OPEN_CONNS", &errs)
	cfg.DB.MaxIdleConns = requireInt(lookup, "NOVASCANS_DB_MAX_IDLE_CONNS", &errs)
	cfg.DB.ConnMaxLifetime = requireDuration(lookup, "NOVASCANS_DB_CONN_MAX_LIFETIME", &errs)

	cfg.Auth.AccessTokenSecret = requireString(lookup, "NOVASCANS_AUTH_ACCESS_TOKEN_SECRET", &errs)
	cfg.Auth.AccessTokenTTL = requireDuration(lookup, "NOVASCANS_AUTH_ACCESS_TOKEN_TTL", &errs)
	cfg.Auth.RefreshTokenTTL = requireDuration(lookup, "NOVASCANS_AUTH_REFRESH_TOKEN_TTL", &errs)
	cfg.Auth.EmailVerificationTokenTTL = requireDuration(lookup, "NOVASCANS_AUTH_EMAIL_VERIFICATION_TOKEN_TTL", &errs)
	cfg.Auth.PasswordResetTokenTTL = requireDuration(lookup, "NOVASCANS_AUTH_PASSWORD_RESET_TOKEN_TTL", &errs)

	cfg.Log.Level = requireString(lookup, "NOVASCANS_LOG_LEVEL", &errs)
	cfg.Metrics.Enabled = requireBool(lookup, "NOVASCANS_METRICS_ENABLED", &errs)

	if len(errs) > 0 {
		return Config{}, errors.New(strings.Join(errs, "; "))
	}

	return cfg, nil
}

func (cfg HTTPConfig) Address() string {
	return net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
}

func (cfg Config) String() string {
	return fmt.Sprintf("%s(%s)", cfg.App.Name, cfg.App.Env)
}

func (cfg DBConfig) ConnectionString(databaseName string) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		databaseName,
		cfg.SSLMode,
	)
}

func requireString(lookup func(string) (string, bool), key string, errs *[]string) string {
	value, ok := lookup(key)
	if !ok || strings.TrimSpace(value) == "" {
		*errs = append(*errs, fmt.Sprintf("%s is required", key))
		return ""
	}

	return value
}

func requireInt(lookup func(string) (string, bool), key string, errs *[]string) int {
	value := requireString(lookup, key, errs)
	if value == "" {
		return 0
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		*errs = append(*errs, fmt.Sprintf("%s must be a valid integer", key))
		return 0
	}

	return parsed
}

func requireDuration(lookup func(string) (string, bool), key string, errs *[]string) time.Duration {
	value := requireString(lookup, key, errs)
	if value == "" {
		return 0
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		*errs = append(*errs, fmt.Sprintf("%s must be a valid duration", key))
		return 0
	}

	return parsed
}

func requireBool(lookup func(string) (string, bool), key string, errs *[]string) bool {
	value := requireString(lookup, key, errs)
	if value == "" {
		return false
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		*errs = append(*errs, fmt.Sprintf("%s must be a valid boolean", key))
		return false
	}

	return parsed
}
