package module

import (
	"database/sql"
	"log/slog"

	"github.com/go-chi/chi/v5"

	"github.com/Tokuchi61/Novascans/internal/platform/config"
	platformdb "github.com/Tokuchi61/Novascans/internal/platform/db"
	"github.com/Tokuchi61/Novascans/internal/platform/events"
	"github.com/Tokuchi61/Novascans/internal/platform/metrics"
	"github.com/Tokuchi61/Novascans/internal/platform/validation"
)

type Dependencies struct {
	Config    config.Config
	Logger    *slog.Logger
	Events    events.Bus
	Metrics   *metrics.Registry
	Validator *validation.Validator
	DB        *sql.DB
	TxManager platformdb.TxManager
}

type Module interface {
	Key() string
	RegisterRoutes(r chi.Router)
}
