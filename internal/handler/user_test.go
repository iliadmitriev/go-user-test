package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/iliadmitriev/go-user-test/internal/domain"
	"github.com/iliadmitriev/go-user-test/internal/mocks"
	"github.com/iliadmitriev/go-user-test/internal/repository"
	"github.com/iliadmitriev/go-user-test/internal/service"
)

type fakeUser struct {
	ID        string    `fake:"{uuid}" json:"id"`
	Login     string    `fake:"{username}" json:"login"`
	Password  string    `fake:"{password}" json:"-"`
	Name      string    `fake:"{firstname}" json:"name"`
	CreatedAt time.Time `fake:"{date}" json:"created_at"`
	UpdatedAt time.Time `fake:"{date}" json:"updated_at"`
}

func (u *fakeUser) toJSON(t *testing.T) string {
	t.Helper()

	data, err := json.Marshal(u)
	if err != nil {
		t.Fatal(err)
	}

	return string(data)
}

func newFakeUser(t *testing.T) *fakeUser {
	t.Helper()

	u := fakeUser{}
	if err := gofakeit.Struct(&u); err != nil {
		t.Fatal(err)
	}

	return &u
}

func Test_userHandler_getUser_SQL_level(t *testing.T) {
	tests := []struct {
		name       string
		closeError error
		rowError   error
		wantCode   int
		wantResp   string
	}{
		{
			name:     "get user by login OK",
			wantCode: http.StatusOK,
		},
		{
			name:     "get user by login not found",
			rowError: sql.ErrNoRows,
			wantCode: http.StatusNotFound,
			wantResp: `{"code":404, "message":"user not found"}`,
		},
		{
			name:       "get user by login connection error",
			closeError: sql.ErrConnDone,
			wantCode:   http.StatusInternalServerError,
			wantResp:   `{"code":500,"message":"sql: connection is already closed"}`,
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

			user := newFakeUser(t)
			url := fmt.Sprintf("http://example.com/user/%s", user.Login)

			// prepare query response

			if tt.closeError != nil {
				dbMock.ExpectQuery(repository.SQLGetUser).WillReturnError(tt.closeError)
			} else {
				rows := sqlmock.NewRows([]string{"id", "loging", "name", "created_at", "updated_at"})
				if tt.rowError != nil {
					// add data to return rows
					rows = rows.CloseError(tt.rowError)
				} else {
					rows = rows.AddRow(user.ID, user.Login, user.Name, user.CreatedAt, user.UpdatedAt)
				}
				// set rows reding errors sql.ErrNoRows
				dbMock.ExpectQuery(repository.SQLGetUser).WithArgs(user.Login).WillReturnRows(rows)
			}

			// create new request
			r, err := http.NewRequest("GET", url, nil)
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
			if tt.wantCode == http.StatusOK {
				assert.JSONEqf(t, user.toJSON(t), w.Body.String(), "response body not match")
			} else {
				assert.JSONEqf(t, tt.wantResp, w.Body.String(), "response body not match")
			}
		})
	}
}

func Test_userHandler_getUser_Repo_level(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		login     string
		data      *domain.User
		dataError error
		wantCode  int
		wantResp  string
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
			dataError: nil,
			wantCode:  http.StatusOK,
			wantResp: `{"id":"70868a75-adbb-4b4d-b482-93915ee11777","login":"b","name":"b",` +
				`"created_at":"2024-11-29T18:33:55.0000001Z","updated_at":"2024-11-29T18:33:55.0000001Z"}`,
		},
		{
			name:      "get user by login not found",
			url:       "http://example.com/user/eee",
			login:     "eee",
			dataError: repository.ErrUserNotFound,
			wantCode:  http.StatusNotFound,
			wantResp:  `{"code":404,"message":"user not found"}`,
		},
		{
			name:      "get user by login connection error",
			url:       "http://example.com/user/eee",
			login:     "eee",
			dataError: sql.ErrConnDone,
			wantCode:  http.StatusInternalServerError,
			wantResp:  `{"code":500,"message":"sql: connection is already closed"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// build whole stack mockRepo -> userService -> userHandler
			mockUserRepo := mocks.NewUserRepository(t)
			userService := service.NewUserService(mockUserRepo)
			userHandler := NewUserHandler(userService)
			handleFunc := userHandler.GetMux()

			// set user repo mock
			// response once on method `getUser`
			mockUserRepo.On("GetUser", mock.Anything, tt.login).Return(tt.data, tt.dataError).Once()

			// create new request
			r, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			// create test response writer
			w := httptest.NewRecorder()

			// run request
			handleFunc.ServeHTTP(w, r)

			// assert
			require.Equal(t, tt.wantCode, w.Code, "status code not match")
			require.JSONEqf(t, tt.wantResp, w.Body.String(), "response body not match")
		})
	}
}
