package repository

import (
	"context"
	"github.com/ZhdanovichVlad/go_final_project/internal/entity"
	"github.com/ZhdanovichVlad/go_final_project/internal/pkg/errorsx"
	"github.com/jackc/pgx/v4"
	"log/slog"
)

const addTransactionPath = `INSERT INTO transactions(sender_uuid, recipient_uuid, quantity)
                             VALUES ($1, $2, $3)`

func (s *storage) BuyMerch(ctx context.Context, userUuid *string, merchInfo *entity.Merch) error {

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	defer tx.Rollback(ctx)
	if err != nil {
		s.logger.Error("failed to begin tx",
			slog.String("method", "storage.BuyMerch"),
			slog.String("error", err.Error()))
		return errorsx.DBError
	}

	merch, err := s.GetMerchWithTh(ctx, tx, merchInfo)
	if err != nil {
		return err
	}

	user, err := s.FindUserBalanceWithTh(ctx, tx, userUuid)

	if user.Balance < merch.Price {
		return errorsx.NotEnoughMoney
	}

	err = s.UpdateInventoryWithTx(ctx, tx, userUuid, &merch.Uuid)
	if err != nil {
		return err
	}

	err = s.UpdateUserBalanceMinusWithTx(ctx, tx, userUuid, &merch.Price)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.logger.Error("failed to commit tx",
			slog.String("method", "service.BuyMerch"),
			slog.String("error", err.Error()))
		return errorsx.ServiceError
	}
	return nil
}

func (s *storage) TransferCoins(ctx context.Context, userUuid *string, receiverUuid *string, amount *int) error {

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	defer tx.Rollback(ctx)
	if err != nil {
		s.logger.Error("failed to begin tx",
			slog.String("method", "storage.BuyMerch"),
			slog.String("error", err.Error()))
		return errorsx.DBError
	}

	user, err := s.FindUserBalanceWithTh(ctx, tx, userUuid)

	if user.Balance < *amount {
		return errorsx.NotEnoughMoney
	}

	err = s.UpdateUserBalanceMinusWithTx(ctx, tx, userUuid, amount)
	if err != nil {
		return err
	}
	err = s.UpdateUserBalancePlusWithTx(ctx, tx, receiverUuid, amount)
	if err != nil {
		return err
	}
	err = s.AddTransactionWithTx(ctx, tx, userUuid, receiverUuid, amount)
	if err != nil {
		return err
	}

	tx.Commit(ctx)
	return nil
}

func (s *storage) AddTransactionWithTx(ctx context.Context,
	tx pgx.Tx,
	senderUuid *string,
	receiverUuid *string,
	amount *int) error {

	_, err := tx.Exec(ctx, addTransactionPath, senderUuid, receiverUuid, amount)
	if err != nil {
		s.logger.Error("failed to add transaction ",
			slog.String("method", "storage.AddTransactionWithTx"),
			slog.String("error", err.Error()))
		return errorsx.DBError
	}
	return nil
}
