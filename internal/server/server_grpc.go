package server

import (
	"context"
	"net"

	"github.com/iliadmitriev/go-user-test/internal/config"
	"github.com/iliadmitriev/go-user-test/internal/handler"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	srv    *grpc.Server
	logger *zap.SugaredLogger
	cfg    *config.Config
}

func NewGRPCServer(handler []handler.GRPCHandler, lc fx.Lifecycle, cfg *config.Config, logger *zap.Logger) Server {
	server := grpc.NewServer()

	for _, handler := range handler {
		handler.RegisterGRPC(server)
	}

	reflection.Register(server)

	srv := &grpcServer{
		srv:    server,
		cfg:    cfg,
		logger: logger.Sugar().Named("GRPCServer"),
	}

	lc.Append(fx.Hook{
		OnStart: srv.Start,
		OnStop:  srv.Shutdown,
	})

	return srv
}

func (srv *grpcServer) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", srv.cfg.ListenGRPC)
	if err != nil {
		srv.logger.Errorw("Error starting server", "error", err)
		return err
	}

	go func(srv *grpcServer) {
		srv.logger.Infow("Server started serving new connections", "addr", srv.cfg.ListenGRPC)
		if err := srv.srv.Serve(listener); err != nil {
			srv.logger.Errorw("Server error", "error", err)
		}
	}(srv)

	return nil
}

func (srv *grpcServer) Shutdown(ctx context.Context) error {
	srv.logger.Infow("Server shutting down", "addr", srv.cfg.ListenGRPC)
	srv.srv.GracefulStop()
	return nil
}
