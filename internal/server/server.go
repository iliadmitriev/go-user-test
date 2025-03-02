package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/iliadmitriev/go-user-test/internal/config"
	"github.com/iliadmitriev/go-user-test/internal/handler"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Server interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type httpServer struct {
	srv    *http.Server
	logger *zap.Logger
	cfg    *config.Config
}

func (srv *httpServer) Start(ctx context.Context) error {
	go func(srv *httpServer) {
		srv.logger.Info("Server started serving new connections", zap.String("addr", srv.cfg.Listen))

		if err := srv.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			srv.logger.Error("Server error", zap.Error(err))
		}

		srv.logger.Info("Server stopped serving new connections")
	}(srv)

	return nil
}

func (srv *httpServer) Shutdown(ctx context.Context) error {
	srv.logger.Info("Server shutting down", zap.String("addr", srv.cfg.Listen))
	return srv.srv.Shutdown(ctx)
}

func NewHTTPServer(handlers []handler.HandlerInterface, lc fx.Lifecycle, cfg *config.Config, logger *zap.Logger) Server {
	mux := http.NewServeMux()

	for _, handler := range handlers {
		handler.GetMux(mux)
	}

	httpSrv := &http.Server{
		Handler:      mux,
		Addr:         cfg.Listen,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	srv := &httpServer{srv: httpSrv, logger: logger.Named("HTTPServer"), cfg: cfg}

	lc.Append(fx.Hook{
		OnStart: srv.Start,
		OnStop:  srv.Shutdown,
	})

	return srv
}
