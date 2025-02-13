package repository

import (
	"github.com/jackc/pgx/v4"
	"log/slog"
)

type storage struct {
	db     *pgx.Conn
	logger *slog.Logger
}

func NewStorage(db *pgx.Conn, logger *slog.Logger) *storage {
	return &storage{db: db, logger: logger}
}
