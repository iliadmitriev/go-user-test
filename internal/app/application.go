package app

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/iliadmitriev/go-user-test/internal/config"
	"github.com/iliadmitriev/go-user-test/internal/db"
	"github.com/iliadmitriev/go-user-test/internal/handler"
	"github.com/iliadmitriev/go-user-test/internal/repository"
	"github.com/iliadmitriev/go-user-test/internal/server"
	"github.com/iliadmitriev/go-user-test/internal/service"
)

func NewApplication() *fx.App {
	return fx.New(
		fx.Provide(fx.Annotate(
			server.NewHTTPServer,
			fx.ParamTags(`group:"http_routes"`),
			fx.ResultTags(`group:"servers"`),
			fx.As(new(server.Server)),
		)),

		fx.Provide(fx.Annotate(
			server.NewGRPCServer,
			fx.ParamTags(`group:"grpc_routes"`),
			fx.ResultTags(`group:"servers"`),
			fx.As(new(server.Server)),
		)),

		fx.Provide(fx.Annotate(
			handler.NewUserHandler,
			fx.ResultTags(`group:"http_routes"`),
			fx.As(new(handler.HTTPHandler)),
		)),

		fx.Provide(fx.Annotate(
			handler.NewGRPCUserHandler,
			fx.ResultTags(`group:"grpc_routes"`),
			fx.As(new(handler.GRPCHandler)),
		)),

		fx.Provide(config.NewConfig),
		fx.Provide(repository.NewUserDB),
		fx.Provide(service.NewUserService),
		fx.Provide(db.NewSqliteDB),
		fx.Provide(zap.NewProduction),

		fx.Invoke(fx.Annotate(
			func([]server.Server) {},
			fx.ParamTags(`group:"servers"`),
		)),

		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log.Named("fx")}
		}),
	)
}
