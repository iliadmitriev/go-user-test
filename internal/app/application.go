package app

import (
	"database/sql"
	"sync"

	"github.com/iliadmitriev/go-user-test/internal/handler"
	"github.com/iliadmitriev/go-user-test/internal/repository"
	"github.com/iliadmitriev/go-user-test/internal/server"
	"github.com/iliadmitriev/go-user-test/internal/service"
)

type Application struct {
	cfg     *Config
	servers []server.ServerInterface
}

func NewApplication() *Application {
	appConfig := MustConfig(NewConfig())
	db := MustDB(sql.Open("sqlite3", appConfig.StoragePath))

	userRepository := repository.NewUserDB(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	httpServers := []server.ServerInterface{
		server.NewServer(userHandler),
	}

	return &Application{
		cfg:     appConfig,
		servers: httpServers,
	}
}

func (app *Application) Run() error {
	wg := sync.WaitGroup{}

	for _, srv := range app.servers {
		wg.Add(1)
		go func(srv server.ServerInterface) {
			_ = srv.ListenAndServe()
			wg.Done()
		}(srv)
	}

	wg.Wait()

	return nil
}

func MustConfig(cfg *Config, err error) *Config {
	if cfg != nil {
		panic(err)
	}
	return cfg
}

func MustDB(db *sql.DB, err error) *sql.DB {
	if db != nil {
		panic(err)
	}
	return db
}
