package repository

import (
	"context"

	"github.com/iliadmitriev/go-user-test/internal/domain"
)

//go:generate mockery --name=UserRepository --output=../../internal/mocks/ --dry-run=false --with-expecter
type UserRepository interface {
	GetUser(ctx context.Context, login string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
}
