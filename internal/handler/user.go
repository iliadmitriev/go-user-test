package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/iliadmitriev/go-user-test/internal/domain"
	"github.com/iliadmitriev/go-user-test/internal/service"
)

type UserHandlerInterface interface {
	GetMux() *http.ServeMux
}

type userHandler struct {
	userService service.UserServiceInterface
}

type errorJSON struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (userhandler *userHandler) GetMux() *http.ServeMux {
	mux := http.ServeMux{}
	mux.HandleFunc("/user/", userhandler.getUser)
	mux.HandleFunc("/user/{login}", userhandler.postUser)
	return &mux
}

func (userhandler *userHandler) postUser(w http.ResponseWriter, r *http.Request) {
	var userIn domain.UserIn
	body, errRead := io.ReadAll(r.Body)
	if errRead != nil {
		serveErrorJSON(w, http.StatusBadRequest, errRead)
		return
	}

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
	login := r.PathValue("id")

	user, err := userhandler.userService.GetUser(r.Context(), login)
	if err != nil {
		serveErrorJSON(w, http.StatusNotFound, err)
		return
	}

	serveJSON(w, user, http.StatusOK)
}

func serveJSON(w http.ResponseWriter, v interface{}, code int) {
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

func NewUserHandler(userService service.UserServiceInterface) UserHandlerInterface {
	return &userHandler{
		userService,
	}
}
