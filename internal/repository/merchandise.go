package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ZhdanovichVlad/go_final_project/internal/entity"
	"github.com/ZhdanovichVlad/go_final_project/internal/pkg/errorsx"
	"github.com/jackc/pgx/v5"
)

const (
	findMerchQuery = "SELECT * FROM merchandise WHERE name = $1"
)

func (s *Storage) GetMerchWithTh(ctx context.Context, tx pgx.Tx, merchInfo *entity.Merch) (*entity.Merch, error) {
	merch := &entity.Merch{}
	err := tx.QueryRow(ctx, findMerchQuery, merchInfo.Name).Scan(&merch.UUID, &merch.Name, &merch.Price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errorsx.ErrItemNotFound
		}
		s.logger.Error("failed to get balance",
			slog.String("method", "storage.GetBalanceWithTh"),
			slog.String("error", err.Error()))
		return nil, errorsx.ErrDB
	}

	return merch, nil
}
