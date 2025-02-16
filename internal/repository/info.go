package repository

import (
	"context"
	"github.com/ZhdanovichVlad/go_final_project/internal/entity"
	"github.com/ZhdanovichVlad/go_final_project/internal/pkg/errorsx"
	"github.com/jackc/pgx/v4"
	"log/slog"
)

const ()

func (s *storage) GetUserInfo(ctx context.Context, userUuid *string) (*entity.UserInfo, error) {

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	defer tx.Rollback(ctx)
	if err != nil {
		s.logger.Error("failed to begin tx",
			slog.String("method", "storage.BuyMerch"),
			slog.String("error", err.Error()))
		return nil, errorsx.DBError
	}
	s.logger.Info("ms1")
	balance, err := s.GetBalanceWithTh(ctx, tx, userUuid)
	if err != nil {
		return nil, errorsx.DBError
	}
	s.logger.Info("ms2")
	inventory, err := s.GetInventory(ctx, tx, userUuid)
	if err != nil {
		return nil, err
	}

	s.logger.Info("ms3")
	coinHistory, err := s.GetTransactionsHistory(ctx, tx, userUuid)
	if err != nil {
		return nil, err
	}
	s.logger.Info("ms4")
	userInfo := &entity.UserInfo{
		Coins:       balance,
		Inventory:   *inventory,
		CoinHistory: *coinHistory,
	}

	return userInfo, nil
}
