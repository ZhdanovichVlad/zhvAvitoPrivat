package repository

import (
	"context"
	"log/slog"

	"github.com/ZhdanovichVlad/go_final_project/internal/entity"
	"github.com/ZhdanovichVlad/go_final_project/internal/pkg/errorsx"
	"github.com/jackc/pgx/v5"
)

const addTransactionPath = `INSERT INTO transactions(sender_uuid, recipient_uuid, quantity)
                             VALUES ($1, $2, $3)`

func (s *Storage) BuyMerch(ctx context.Context, userUUID *string, merchInfo *entity.Merch) error {

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		s.logger.Error("failed to begin tx",
			slog.String("method", "storage.BuyMerch"),
			slog.String("error", err.Error()))
		return errorsx.ErrDB
	}

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

	merch, err := s.GetMerchWithTh(ctx, tx, merchInfo)
	if err != nil {
		return err
	}

	user, err := s.FindUserBalanceWithTh(ctx, tx, userUUID)

	if user.Balance < merch.Price {
		return errorsx.ErrNotEnoughMoney
	}

	err = s.UpdateInventoryWithTx(ctx, tx, userUUID, &merch.UUID)
	if err != nil {
		return err
	}

	err = s.UpdateUserBalanceMinusWithTx(ctx, tx, userUUID, &merch.Price)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.logger.Error("failed to commit tx",
			slog.String("method", "service.BuyMerch"),
			slog.String("error", err.Error()))
		return errorsx.ErrService
	}
	err = nil
	return nil
}

func (s *Storage) TransferCoins(ctx context.Context, userUUID *string, receiverUUID *string, amount *int) error {

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		s.logger.Error("failed to begin tx",
			slog.String("method", "storage.TransferCoins"),
			slog.String("error", err.Error()))
		return errorsx.ErrDB
	}

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

	user, err := s.FindUserBalanceWithTh(ctx, tx, userUUID)
	if err != nil {
		return err
	}

	if user.Balance < *amount {
		return errorsx.ErrNotEnoughMoney
	}

	err = s.UpdateUserBalanceMinusWithTx(ctx, tx, userUUID, amount)
	if err != nil {
		return err
	}
	err = s.UpdateUserBalancePlusWithTx(ctx, tx, receiverUUID, amount)
	if err != nil {
		return err
	}
	err = s.AddTransactionWithTx(ctx, tx, userUUID, receiverUUID, amount)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.logger.Error("failed to Commit tx",
			slog.String("method", "storage.TransferCoins"),
			slog.String("error", err.Error()))
		return errorsx.ErrDB
	}
	err = nil
	return nil
}

func (s *Storage) AddTransactionWithTx(ctx context.Context,
	tx pgx.Tx,
	senderUUID *string,
	receiverUUID *string,
	amount *int) error {

	_, err := tx.Exec(ctx, addTransactionPath, senderUUID, receiverUUID, amount)
	if err != nil {
		s.logger.Error("failed to add transaction ",
			slog.String("method", "storage.AddTransactionWithTx"),
			slog.String("error", err.Error()))
		return errorsx.ErrDB
	}
	return nil
}
