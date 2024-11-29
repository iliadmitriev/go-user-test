package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iliadmitriev/go-user-test/internal/domain"
	"github.com/iliadmitriev/go-user-test/internal/repository"
)

var (
	ErrUserAlreadyExists = errors.New("user with login already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type UserServiceInterface interface {
	GetUser(ctx context.Context, login string) (*domain.UserOut, error)
	CreateUser(ctx context.Context, user *domain.UserIn) (*domain.UserOut, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func (userservice *userService) GetUser(ctx context.Context, login string) (*domain.UserOut, error) {
	user, err := userservice.userRepository.GetUser(ctx, login)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		} else {
			return nil, err
		}
	}

	return &domain.UserOut{
		ID:        user.ID,
		Login:     user.Login,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (userservice *userService) CreateUser(ctx context.Context, user *domain.UserIn) (*domain.UserOut, error) {
	if _, err := userservice.userRepository.GetUser(ctx, user.Login); err == nil {
		return nil, ErrUserAlreadyExists
	}

	id := uuid.New()

	userSave := &domain.User{
		ID:        id,
		Login:     user.Login,
		Password:  user.Password,
		Name:      user.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := userservice.userRepository.CreateUser(ctx, userSave); err != nil {
		return nil, err
	}

	return userservice.GetUser(ctx, userSave.Login)
}

func NewUserService(userRepository repository.UserRepository) UserServiceInterface {
	return &userService{
		userRepository,
	}
}
