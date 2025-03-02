package server

import (
	"context"
	"errors"
	"net"
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
	logger *zap.SugaredLogger
	cfg    *config.Config
}

func (srv *httpServer) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", srv.cfg.Listen)
	if err != nil {
		srv.logger.Errorw("Error starting server", "error", err)
		return err
	}

	go func(srv *httpServer) {
		srv.logger.Infow("Server started serving new connections", "addr", srv.cfg.Listen)

		if err := srv.srv.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
			srv.logger.Errorw("Server error", "error", err)
		}

		srv.logger.Info("Server stopped serving new connections")
	}(srv)

	return nil
}

func (srv *httpServer) Shutdown(ctx context.Context) error {
	srv.logger.Infow("Server shutting down", "addr", srv.cfg.Listen)
	return srv.srv.Shutdown(ctx)
}

func NewHTTPServer(handlers []handler.HTTPHandler, lc fx.Lifecycle, cfg *config.Config, logger *zap.Logger) Server {
	mux := http.NewServeMux()

	for _, handler := range handlers {
		handler.GetMux(mux)
	}

	srv := &httpServer{
		srv: &http.Server{
			Handler:      mux,
			Addr:         cfg.Listen,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
		cfg:    cfg,
		logger: logger.Sugar().Named("HTTPServer"),
	}

	lc.Append(fx.Hook{
		OnStart: srv.Start,
		OnStop:  srv.Shutdown,
	})

	return srv
}
