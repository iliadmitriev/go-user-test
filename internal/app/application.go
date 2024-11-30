package app

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"

	"github.com/iliadmitriev/go-user-test/internal/config"
	"github.com/iliadmitriev/go-user-test/internal/handler"
	"github.com/iliadmitriev/go-user-test/internal/repository"
	"github.com/iliadmitriev/go-user-test/internal/server"
	"github.com/iliadmitriev/go-user-test/internal/service"
)

type Application struct {
	cfg     *config.Config
	logger  *zap.Logger
	servers []server.ServerInterface
}

func NewApplication() *Application {
	logger, _ := zap.NewDevelopment()

	logger.Info("Logger initialized")
	logger.Info("Starting application")

	appConfig := config.MustConfig(config.NewConfig())
	logger.Info("Config loaded")

	db := MustDB(sql.Open("sqlite3", appConfig.StoragePath))
	logger.Info("DB connection established")

	userRepository := repository.NewUserDB(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	httpServers := []server.ServerInterface{
		server.NewServer(userHandler, logger, appConfig),
	}

	return &Application{
		cfg:     appConfig,
		logger:  logger,
		servers: httpServers,
	}
}

func (app *Application) Run() error {
	defer func() {
		_ = app.logger.Sync()
	}()

	var wg sync.WaitGroup
	wg.Add(len(app.servers))

	for _, srv := range app.servers {
		srv.Start()
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownRelease()

	for _, srv := range app.servers {
		go func() {
			if err := srv.Shutdown(shutdownCtx); err != nil {
				app.logger.Error("Shutdown error", zap.Error(err))
			}

			wg.Done()
		}()
	}

	wg.Wait()

	app.logger.Info("Application shutdown")

	return nil
}

func MustDB(db *sql.DB, err error) *sql.DB {
	if err != nil {
		panic(err)
	}
	return db
}
