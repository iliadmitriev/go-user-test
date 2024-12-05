package server

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/iliadmitriev/go-user-test/internal/config"
	"github.com/iliadmitriev/go-user-test/internal/handler"
	"go.uber.org/zap"
)

type ServerInterface interface {
	Start()
	Shutdown(ctx context.Context) error
}

type server struct {
	srv    *http.Server
	logger *zap.Logger
	cfg    *config.Config
}

func (srv *server) Start() {
	go func(srv *server) {
		srv.logger.Info("Server started serving new connections", zap.String("addr", srv.cfg.Listen))

		if err := srv.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			srv.logger.Error("Server error", zap.Error(err))
			os.Exit(1)
		}

		srv.logger.Info("Server stopped serving new connections")
	}(srv)
}

func (srv *server) Shutdown(ctx context.Context) error {
	srv.logger.Info("Server shutting down", zap.String("addr", srv.cfg.Listen))
	return srv.srv.Shutdown(ctx)
}

func NewServer(handler handler.UserHandlerInterface, logger *zap.Logger, cfg *config.Config) ServerInterface {
	srv := http.Server{
		Handler:      handler.GetMux(),
		Addr:         cfg.Listen,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return &server{srv: &srv, logger: logger, cfg: cfg}
}
