package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/iliadmitriev/go-user-test/internal/domain"
	"github.com/iliadmitriev/go-user-test/internal/repository"
	"github.com/iliadmitriev/go-user-test/internal/service"
)

func Test_userHandler_getUser_SQL_level(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		login      string
		data       *domain.User
		closeError error
		wantError  string
		wantCode   int
		wantResp   string
	}{
		{
			name:  "get user by login OK",
			url:   "http://example.com/user/b",
			login: "b",
			data: &domain.User{
				ID:        uuid.MustParse("70868a75-adbb-4b4d-b482-93915ee11777"),
				Login:     "b",
				Name:      "b",
				CreatedAt: time.Date(2024, 11, 29, 18, 33, 55, 100, time.UTC),
				UpdatedAt: time.Date(2024, 11, 29, 18, 33, 55, 100, time.UTC),
			},
			wantCode: http.StatusOK,
			wantResp: `{"id":"70868a75-adbb-4b4d-b482-93915ee11777","login":"b","name":"b",` +
				`"created_at":"2024-11-29T18:33:55.0000001Z","updated_at":"2024-11-29T18:33:55.0000001Z"}`,
		},
		{
			name:       "get user by login not found",
			url:        "http://example.com/user/eee",
			login:      "eee",
			closeError: sql.ErrNoRows,
			wantCode:   http.StatusNotFound,
			wantResp:   `{"code":404,"message":"user not found"}`,
		},
		{
			name:       "get user by login connection error",
			url:        "http://example.com/user/eee",
			login:      "eee",
			closeError: sql.ErrConnDone,
			wantCode:   http.StatusInternalServerError,
			wantResp:   `{"code":500,"message":"sql: connection is already closed"}`,
		},
		{
			name:      "get user by login SQL error",
			login:     "eee",
			url:       "http://example.com/user/eee",
			wantCode:  http.StatusInternalServerError,
			wantError: "some error",
			wantResp:  `{"code":500,"message":"some error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// build testing stack
			// db -> repository -> service -> handler
			db, dbMock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("err not expected: %v", err)
			}
			userRepository := repository.NewUserDB(db)
			userService := service.NewUserService(userRepository)
			userHandler := NewUserHandler(userService)
			handleFunc := userHandler.GetMux()

			// prepare query response

			if tt.wantError != "" {
				dbMock.ExpectQuery(repository.SQLGetUser).WillReturnError(errors.New(tt.wantError))
			} else {
				rows := sqlmock.NewRows([]string{"id", "loging", "name", "created_at", "updated_at"})
				if tt.data != nil {
					// add data to return rows
					rows = rows.AddRow(tt.data.ID, tt.data.Login, tt.data.Name, tt.data.CreatedAt, tt.data.UpdatedAt)
				}
				// set rows reding errors sql.ErrNoRows
				rows = rows.CloseError(tt.closeError)
				dbMock.ExpectQuery(repository.SQLGetUser).WithArgs(tt.login).WillReturnRows(rows)
			}

			// create new request
			r, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			// create test response writer
			w := httptest.NewRecorder()

			// run request
			handleFunc.ServeHTTP(w, r)

			// assert if db worked properly
			if err := dbMock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			// assert status and data
			assert.Equal(t, tt.wantCode, w.Code, "status code not match")
			assert.JSONEqf(t, tt.wantResp, w.Body.String(), "response body not match")
		})
	}
}
