package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"github.com/nava1525/bilio-backend/internal/app/transport"
	"github.com/nava1525/bilio-backend/internal/config"
)

type HTTPServer struct {
	server          *http.Server
	log             zerolog.Logger
	shutdownTimeout time.Duration
}

func NewHTTPServer(cfg *config.Config, log zerolog.Logger, db *sql.DB) (*HTTPServer, error) {
	router, err := transport.NewRouter(cfg, log, db)
	if err != nil {
		return nil, err
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return &HTTPServer{
		server:          srv,
		log:             log,
		shutdownTimeout: cfg.Server.ShutdownTimeout,
	}, nil
}

func (s *HTTPServer) Start() error {
	s.log.Info().Str("addr", s.server.Addr).Msg("http server starting")

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *HTTPServer) Stop() error {
	s.log.Info().Msg("http server shutting down")
	s.server.SetKeepAlivesEnabled(false)

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	s.log.Info().Msg("http server shutdown complete")
	return nil
}
