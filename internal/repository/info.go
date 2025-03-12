package repository

import (
	"context"
	"log/slog"

	"github.com/ZhdanovichVlad/go_final_project/internal/entity"
	"github.com/ZhdanovichVlad/go_final_project/internal/pkg/errorsx"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) GetUserInfo(ctx context.Context, userUUID *string) (*entity.UserInfo, error) {

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				s.logger.Error("failed to rollback tx",
					slog.String("method", "storage.BuyMerch"),
					slog.String("error", rollbackErr.Error()))
			}
		}
	}()

	if err != nil {
		s.logger.Error("failed to begin tx",
			slog.String("method", "storage.GetUserInfo"),
			slog.String("error", err.Error()))
		return nil, errorsx.ErrDB
	}

	balance, err := s.GetBalanceWithTh(ctx, tx, userUUID)
	if err != nil {
		return nil, errorsx.ErrDB
	}

	inventory, err := s.GetInventory(ctx, tx, userUUID)
	if err != nil {
		return nil, err
	}

	coinHistory, err := s.GetTransactionsHistory(ctx, tx, userUUID)
	if err != nil {
		return nil, err
	}

	userInfo := &entity.UserInfo{
		Coins:       balance,
		Inventory:   *inventory,
		CoinHistory: *coinHistory,
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, errorsx.ErrDB
	}

	err = nil
	return userInfo, nil
}
