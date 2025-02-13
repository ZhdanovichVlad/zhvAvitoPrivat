package repository

import (
	"context"
	"errors"
	"github.com/ZhdanovichVlad/go_final_project/internal/entity"
	"github.com/ZhdanovichVlad/go_final_project/internal/pkg/errorsx"
	"log/slog"

	"github.com/jackc/pgx/v4"
)

const (
	saveUserQuery               = "INSERT INTO users (username,password_hash) VALUES ($1, $2) RETURNING uuid;"
	findUserQuery               = "SELECT * FROM users WHERE username = $1"
	findBalanceQuery            = "SELECT balance FROM users WHERE uuid = $1"
	existsUserQuery             = "SELECT EXISTS(SELECT 1 FROM users WHERE uuid=$1)"
	findUserBalanceWithThQuery  = "SELECT uuid, balance FROM users WHERE uuid = $1"
	updateUserBalanceMinusQuery = "UPDATE users SET balance = balance - $1 WHERE uuid = $2"
	updateUserBalancePlusQuery  = "UPDATE users SET balance = balance + $1 WHERE uuid = $2"
	getItemsFromTheInventory    = `SELECT m.name, oi.quantity FROM owned_inventory oi
                                  JOIN merchandise m ON oi.merchandise_uuid = m.uuid
                                  WHERE oi.user_uuid = $1;`
	updateOwnedInventoryQuery = `INSERT INTO owned_inventory (user_uuid, merchandise_uuid, quantity)
                                  VALUES ($1, $2, 1)
                                  ON CONFLICT (user_uuid, merchandise_uuid)
                                  DO UPDATE SET quantity = owned_inventory.quantity + 1;`
	GetTransactionsHistoryQuery = `SELECT t.sender_uuid, sender.username, t.recipient_uuid, recipient.username, t.quantity
                                   FROM  transactions t
                                   JOIN  users sender ON t.sender_uuid = sender.uuid
                                   JOIN  users recipient ON t.recipient_uuid = recipient.uuid 
                                   WHERE t.sender_uuid = $1 OR t.recipient_uuid = $1;`
)

func (s *storage) SaveUser(ctx context.Context, user *entity.User) (*string, error) {
	var uuid string

	err := s.db.QueryRow(ctx, saveUserQuery, user.Username, user.PasswordHash).Scan(&uuid)
	if err != nil {
		s.logger.Error("failed to add user",
			slog.String("method", "storage.saveUser"),
			slog.String("err", err.Error()))
		return nil, errorsx.DBError
	}

	return &uuid, nil
}

func (s *storage) FindUser(ctx context.Context, userRequest *entity.User) (*entity.User, error) {
	user := entity.User{}

	err := s.db.QueryRow(ctx, findUserQuery, userRequest.Username).Scan(&user.Uuid, &user.Username, &user.PasswordHash, &user.Balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		s.logger.Error("failed to find user.",
			slog.String("method", "storage.FindUser"),
			slog.String("err", err.Error()))
		return nil, errorsx.DBError
	}
	return &user, nil
}

func (s *storage) FindUserBalanceWithTh(ctx context.Context, tx pgx.Tx, userUuid *string) (*entity.User, error) {
	user := entity.User{}

	err := tx.QueryRow(ctx, findUserBalanceWithThQuery, userUuid).Scan(&user.Uuid, &user.Balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		s.logger.Error("failed to find user.",
			slog.String("method", "storage.FindUser"),
			slog.String("err", err.Error()))
		return nil, errorsx.DBError
	}

	return &user, nil
}

func (s *storage) GetBalanceWithTh(ctx context.Context, tx pgx.Tx, userUuid *string) (int, error) {
	var balance int

	err := tx.QueryRow(ctx, findBalanceQuery, userUuid).Scan(&balance)
	if err != nil {
		s.logger.Error("failed to get balance",
			slog.String("method", "storage.GetBalanceWithTh"),
			slog.String("err", err.Error()))
		return 0, errorsx.DBError
	}

	return balance, nil
}

func (s *storage) ExistsUser(ctx context.Context, userUuid *string) (*bool, error) {
	var exists bool
	err := s.db.QueryRow(ctx, existsUserQuery, userUuid).Scan(&exists)
	if err != nil {
		s.logger.Error("failed to verify user exists.",
			slog.String("method", "storage.ExistsUser"),
			slog.String("err", err.Error()))
		return nil, errorsx.DBError
	}
	return &exists, nil
}

func (s *storage) UpdateUserBalanceMinusWithTx(ctx context.Context, tx pgx.Tx, userUuid *string, merchPrice *int) error {

	_, err := tx.Exec(ctx, updateUserBalanceMinusQuery, merchPrice, userUuid)
	if err != nil {
		s.logger.Error("failed to get balance",
			slog.String("method", "storage.GetBalanceWithTh"),
			slog.String("err", err.Error()))
		return errorsx.DBError
	}

	return nil
}

func (s *storage) UpdateUserBalancePlusWithTx(ctx context.Context, tx pgx.Tx, userUuid *string, merchPrice *int) error {

	_, err := tx.Exec(ctx, updateUserBalancePlusQuery, merchPrice, userUuid)
	if err != nil {
		s.logger.Error("failed to get balance",
			slog.String("method", "storage.GetBalanceWithTh"),
			slog.String("err", err.Error()))
		return errorsx.DBError
	}

	return nil
}

func (s *storage) UpdateInventoryWithTx(ctx context.Context, tx pgx.Tx, userUuid *string, merchUuid *string) error {

	_, err := tx.Exec(ctx, updateOwnedInventoryQuery, userUuid, merchUuid)
	if err != nil {
		s.logger.Error("failed to update owned inventory query",
			slog.String("method", "storage.BuyMerch"),
			slog.String("err", err.Error()))
		return errorsx.DBError
	}

	return nil
}

func (s *storage) GetInventory(ctx context.Context, tx pgx.Tx, userUuid *string) (*[]entity.Item, error) {

	var inventory []entity.Item
	rows, err := tx.Query(ctx, getItemsFromTheInventory, userUuid)
	defer rows.Close()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &inventory, nil
		}
		s.logger.Error("failed get Inventory",
			slog.String("method", "storage.GetInventory"),
			slog.String("err", err.Error()))
		return nil, errorsx.DBError
	}

	for rows.Next() {
		var item entity.Item
		if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
			s.logger.Error("scan rows error",
				slog.String("method", "storage.GetInventory"),
				slog.String("err", err.Error()))
			return nil, errorsx.DBError
		}
		inventory = append(inventory, item)
	}
	if err := rows.Err(); err != nil {
		s.logger.Error("rows.Err() return error",
			slog.String("method", "storage.GetInventory"),
			slog.String("err", err.Error()))
		return nil, errorsx.DBError
	}
	return &inventory, nil
}

func (s *storage) GetTransactionsHistory(ctx context.Context, tx pgx.Tx, userUuid *string) (*entity.CoinHistory, error) {

	rows, err := tx.Query(ctx, GetTransactionsHistoryQuery, userUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var received []entity.UserTransfer
	var sent []entity.UserTransfer
	for rows.Next() {
		var transaction entity.Transaction
		err = rows.Scan(&transaction.SenderUuid, &transaction.Sender, &transaction.ReceiverUuid, &transaction.Receiver, &transaction.Amount)
		if err != nil {
			s.logger.Error("scan rows error",
				slog.String("method", "storage.GetTransactionsHistory"),
				slog.String("err", err.Error()))
			return nil, errorsx.DBError
		}

		if transaction.ReceiverUuid == *userUuid {
			received = append(received, entity.UserTransfer{transaction.Sender, transaction.Amount})
		} else if transaction.SenderUuid == *userUuid {
			sent = append(sent, entity.UserTransfer{transaction.Receiver, transaction.Amount})
		}
	}

	coinHistory := entity.CoinHistory{received, sent}
	return &coinHistory, nil
}
