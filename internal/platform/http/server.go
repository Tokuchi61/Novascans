package http

import (
	stdhttp "net/http"

	"github.com/Tokuchi61/Novascans/internal/platform/config"
)

func NewServer(cfg config.HTTPConfig, handler stdhttp.Handler) *stdhttp.Server {
	return &stdhttp.Server{
		Addr:         cfg.Address(),
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}
