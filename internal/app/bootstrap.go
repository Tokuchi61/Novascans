package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Tokuchi61/Novascans/internal/platform/config"
	platformdb "github.com/Tokuchi61/Novascans/internal/platform/db"
	"github.com/Tokuchi61/Novascans/internal/platform/events"
	"github.com/Tokuchi61/Novascans/internal/platform/logger"
	"github.com/Tokuchi61/Novascans/internal/platform/metrics"
	moduleshared "github.com/Tokuchi61/Novascans/internal/platform/module"
	"github.com/Tokuchi61/Novascans/internal/platform/validation"
)

type Runtime struct {
	Config    config.Config
	Logger    *slog.Logger
	Events    events.Bus
	Metrics   *metrics.Registry
	Validator *validation.Validator
	Database  *sql.DB
	Modules   []moduleshared.Module
	Router    http.Handler
	ReadyFunc func(context.Context) error
}

type Option func(*bootstrapOptions)

type bootstrapOptions struct {
	db     *sql.DB
	skipDB bool
}

func WithSkipDB() Option {
	return func(options *bootstrapOptions) {
		options.skipDB = true
	}
}

func WithDB(db *sql.DB) Option {
	return func(options *bootstrapOptions) {
		options.db = db
	}
}

func Bootstrap(ctx context.Context, cfg config.Config, opts ...Option) (*Runtime, error) {
	log := logger.New(cfg.Log)
	bus := events.NewInMemoryBus()
	metricRegistry := metrics.NewRegistry()
	inputValidator := validation.New()

	options := bootstrapOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	database := options.db
	if database == nil && !options.skipDB {
		var err error
		database, err = platformdb.OpenPostgres(ctx, cfg.DB)
		if err != nil {
			return nil, fmt.Errorf("bootstrap database: %w", err)
		}
	}

	var txManager platformdb.TxManager
	if database != nil {
		txManager = platformdb.NewTxManager(database)
	}

	deps := moduleshared.Dependencies{
		Config:    cfg,
		Logger:    log,
		Events:    bus,
		Metrics:   metricRegistry,
		Validator: inputValidator,
		DB:        database,
		TxManager: txManager,
	}

	modules := buildModules(deps)
	if len(modules) == 0 {
		return nil, fmt.Errorf("bootstrap modules: no modules registered")
	}

	readyFunc := func(ctx context.Context) error {
		if options.skipDB || database == nil {
			return nil
		}

		readyCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		if err := database.PingContext(readyCtx); err != nil {
			return fmt.Errorf("database not ready: %w", err)
		}

		return nil
	}

	router := buildRouter(routeOptions{
		Config:    cfg,
		Logger:    log,
		Metrics:   metricRegistry,
		Modules:   modules,
		ReadyFunc: readyFunc,
	})

	return &Runtime{
		Config:    cfg,
		Logger:    log,
		Events:    bus,
		Metrics:   metricRegistry,
		Validator: inputValidator,
		Database:  database,
		Modules:   modules,
		Router:    router,
		ReadyFunc: readyFunc,
	}, nil
}

func (runtime *Runtime) Close() error {
	if runtime.Database == nil {
		return nil
	}

	return runtime.Database.Close()
}
