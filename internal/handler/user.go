package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/iliadmitriev/go-user-test/internal/domain"
	"github.com/iliadmitriev/go-user-test/internal/service"
)

type HandlerInterface interface {
	GetMux(mux *http.ServeMux)
}

type userHandler struct {
	userService service.UserServiceInterface
	logger      *zap.Logger
}

type errorJSON struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (userhandler *userHandler) GetMux(mux *http.ServeMux) {
	mux.HandleFunc("/user/", userhandler.postUser)
	mux.HandleFunc("/user/{login}", userhandler.getUser)
}

func (userhandler *userHandler) postUser(w http.ResponseWriter, r *http.Request) {
	var userIn domain.UserIn
	body, errRead := io.ReadAll(r.Body)
	if errRead != nil {
		serveErrorJSON(w, http.StatusBadRequest, errRead)
		return
	}

	userhandler.logger.Info("Got request", zap.String("body", string(body)))

	if err := json.Unmarshal(body, &userIn); err != nil {
		serveErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	user, err := userhandler.userService.CreateUser(r.Context(), &userIn)
	if err != nil {
		serveErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	serveJSON(w, user, http.StatusCreated)
}

func (userhandler *userHandler) getUser(w http.ResponseWriter, r *http.Request) {
	login := r.PathValue("login")

	ctx := r.Context()

	user, err := userhandler.userService.GetUser(ctx, login)
	if errors.Is(err, service.ErrUserNotFound) {
		serveErrorJSON(w, http.StatusNotFound, err)
		return
	}

	userhandler.logger.Info("Got request", zap.String("login", login))

	if err != nil {
		serveErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	serveJSON(w, user, http.StatusOK)
}

func serveJSON(w http.ResponseWriter, v any, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	_ = encoder.Encode(v)
}

func serveErrorJSON(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	_ = encoder.Encode(errorJSON{Message: err.Error(), Code: code})
}

func NewUserHandler(userService service.UserServiceInterface, logger *zap.Logger) HandlerInterface {
	return &userHandler{
		userService,
		logger.Named("userHandler"),
	}
}
