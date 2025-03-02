package repository

import (
	"context"
	"errors"

	"github.com/iliadmitriev/go-user-test/internal/db"
	"github.com/iliadmitriev/go-user-test/internal/domain"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserLoginExists = errors.New("user with login already exists")
)

//go:generate mockery --name=UserRepository --output=../../internal/mocks/ --dry-run=false --with-expecter
type UserRepository interface {
	GetUser(ctx context.Context, login string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
}

func NewUserDB(db db.DB) UserRepository {
	return &UserDB{db}
}

type UserDB struct {
	db db.DB
}

var _ UserRepository = (*UserDB)(nil)

const (
	SQLGetUser    = `SELECT id, login, name, created_at, updated_at FROM users WHERE login = ?`
	SQLCreateUser = `INSERT INTO users (id, login, password, name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
)

func (u *UserDB) GetUser(ctx context.Context, login string) (*domain.User, error) {
	rows, err := u.db.QueryContext(ctx, SQLGetUser, login)
	if err != nil {
		return nil, err
	}

	rows.Next()
	var user domain.User
	if err := rows.Scan(&user.ID, &user.Login, &user.Name, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, ErrUserNotFound
	}

	return &user, nil
}

func (u *UserDB) CreateUser(ctx context.Context, user *domain.User) error {
	_, err := u.db.ExecContext(ctx, SQLCreateUser, user.ID, user.Login, user.Password, user.Name, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return ErrUserLoginExists
	}
	return nil
}
