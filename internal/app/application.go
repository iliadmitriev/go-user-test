package app

import (
	"go.uber.org/fx"
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
		fx.Provide(
			fx.Annotate(
				server.NewHTTPServer, fx.ParamTags(`group:"routes"`),
			),
		),

		fx.Provide(AsRoute(handler.NewUserHandler)),

		fx.Provide(config.NewConfig),
		fx.Provide(repository.NewUserDB),
		fx.Provide(service.NewUserService),
		fx.Provide(db.NewSqliteDB),
		fx.Provide(zap.NewProduction),

		fx.Invoke(func(server.Server) {}),
	)
}

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(handler.HandlerInterface)),
		fx.ResultTags(`group:"routes"`),
	)
}
