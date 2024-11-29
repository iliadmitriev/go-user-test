package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type UserOut struct {
	ID        uuid.UUID `json:"id"`
	Login     string    `json:"login"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        uuid.UUID
	Login     string
	Password  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
