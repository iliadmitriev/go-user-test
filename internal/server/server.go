package server

import (
	"net/http"
	"time"

	"github.com/iliadmitriev/go-user-test/internal/handler"
)

type ServerInterface interface {
	ListenAndServe() error
}

type server struct {
	mux *http.ServeMux
}

func (server *server) ListenAndServe() error {
	srv := http.Server{
		Handler:      server.mux,
		Addr:         ":8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return srv.ListenAndServe()
}

func NewServer(handler handler.UserHandlerInterface) ServerInterface {
	return &server{handler.GetMux()}
}
