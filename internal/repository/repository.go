package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type Storage struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewStorage(db *pgxpool.Pool, logger *slog.Logger) *Storage {
	return &Storage{db: db, logger: logger}
}
