package db

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/iliadmitriev/go-user-test/internal/config"
)

type DB interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

func NewSqliteDB(cfg *config.Config) (DB, error) {
	return sql.Open("sqlite3", cfg.StoragePath)
}
