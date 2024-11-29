package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/iliadmitriev/go-user-test/internal/domain"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserLoginExists = errors.New("user with login already exists")
)

//go:generate mockery --name=DB --output=../../internal/mocks/ --dry-run=false --with-expecter
type DB interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

func NewUserDB(db DB) *UserDB {
	return &UserDB{db}
}

type UserDB struct {
	db DB
}

var _ UserRepository = (*UserDB)(nil)

var (
	sqlGetUser    = `SELECT id, login, name, created_at, updated_at FROM users WHERE login = ?`
	sqlCreateUser = `INSERT INTO users (id, login, password, name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
)

func (u *UserDB) GetUser(ctx context.Context, login string) (*domain.User, error) {
	rows, err := u.db.QueryContext(ctx, sqlGetUser, login)
	if err != nil {
		return nil, ErrUserNotFound
	}

	rows.Next()
	var user domain.User
	if err := rows.Scan(&user.ID, &user.Login, &user.Name, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (u *UserDB) CreateUser(ctx context.Context, user *domain.User) error {
	_, err := u.db.ExecContext(ctx, sqlCreateUser, user.ID, user.Login, user.Password, user.Name, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return ErrUserLoginExists
	}
	return nil
}
