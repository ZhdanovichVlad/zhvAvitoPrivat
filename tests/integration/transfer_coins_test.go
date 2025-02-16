package integration

import (
	"context"
	"github.com/ZhdanovichVlad/go_final_project/internal/entity"
	"github.com/ZhdanovichVlad/go_final_project/internal/pkg/errorsx"
	"github.com/ZhdanovichVlad/go_final_project/internal/repository"
	"github.com/ZhdanovichVlad/go_final_project/internal/service"
	"github.com/ZhdanovichVlad/go_final_project/tests/integration/test_repository"
	"github.com/ZhdanovichVlad/go_final_project/tests/unit/servicetest/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestTransferCoins_Successful(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := test_repository.SetupTestDB(t)
	logger := slog.Default()
	repo := repository.NewStorage(db, logger)
	tokenGen := mocks.NewMocktokenGenerator(ctrl)
	service := service.NewService(repo, logger, tokenGen)
	ctx := context.Background()
	sender := "tester"
	var senderUserBeforeTransferCoins entity.User

	err := db.QueryRow(ctx, findUserPath, sender).Scan(&senderUserBeforeTransferCoins.UUID,
		&senderUserBeforeTransferCoins.Username,
		&senderUserBeforeTransferCoins.Balance)
	assert.NoError(t, err)

	recipient := "testerAvito"
	var recipientUserBeforeTransferCoins entity.User

	err = db.QueryRow(ctx, findUserPath, recipient).Scan(&recipientUserBeforeTransferCoins.UUID,
		&recipientUserBeforeTransferCoins.Username,
		&recipientUserBeforeTransferCoins.Balance)
	assert.NoError(t, err)

	transferAmount := 100
	sendingCoinsInfo := entity.SendingCoins{recipient, transferAmount}

	err = service.TransferCoins(ctx, &senderUserBeforeTransferCoins.UUID, &sendingCoinsInfo)
	assert.NoError(t, err)

	var senderUserAfterTransferCoins entity.User
	err = db.QueryRow(ctx, findUserPath, sender).Scan(&senderUserAfterTransferCoins.UUID,
		&senderUserAfterTransferCoins.Username,
		&senderUserAfterTransferCoins.Balance)
	assert.NoError(t, err)

	var recipientUserAfterTransferCoins entity.User
	err = db.QueryRow(ctx, findUserPath, recipient).Scan(&recipientUserAfterTransferCoins.UUID,
		&recipientUserAfterTransferCoins.Username,
		&recipientUserAfterTransferCoins.Balance)
	assert.NoError(t, err)

	balanceDifferenceSender := senderUserBeforeTransferCoins.Balance - senderUserAfterTransferCoins.Balance
	assert.Equal(t, transferAmount, transferAmount, balanceDifferenceSender)

	balanceDifferenceRecipient := recipientUserAfterTransferCoins.Balance - recipientUserBeforeTransferCoins.Balance
	assert.Equal(t, balanceDifferenceRecipient, transferAmount)

	transactions := []entity.Transaction{}

	rows, err := db.Query(ctx, findTransactionsBySender, senderUserBeforeTransferCoins.UUID)

	for rows.Next() {
		var transaction entity.Transaction
		err = rows.Scan(&transaction.SenderUUID, &transaction.ReceiverUUID, &transaction.Amount)
		assert.NoError(t, err)
		transactions = append(transactions, transaction)
	}

	assert.Equal(t, 1, len(transactions))
	assert.Equal(t, senderUserBeforeTransferCoins.UUID, transactions[0].SenderUUID)
	assert.Equal(t, recipientUserBeforeTransferCoins.UUID, transactions[0].ReceiverUUID)
	assert.Equal(t, transferAmount, transactions[0].Amount)
}

func TestTransferCoins_NotEnoughMoney(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := test_repository.SetupTestDB(t)
	logger := slog.Default()
	repo := repository.NewStorage(db, logger)
	tokenGen := mocks.NewMocktokenGenerator(ctrl)
	service := service.NewService(repo, logger, tokenGen)
	ctx := context.Background()

	sender := "tester"
	var senderUserBeforeTransferCoins entity.User

	err := db.QueryRow(ctx, findUserPath, sender).Scan(&senderUserBeforeTransferCoins.UUID,
		&senderUserBeforeTransferCoins.Username,
		&senderUserBeforeTransferCoins.Balance)
	assert.NoError(t, err)

	recipient := "testerAvito"
	var recipientUserBeforeTransferCoins entity.User

	err = db.QueryRow(ctx, findUserPath, recipient).Scan(&recipientUserBeforeTransferCoins.UUID,
		&recipientUserBeforeTransferCoins.Username,
		&recipientUserBeforeTransferCoins.Balance)
	assert.NoError(t, err)

	transferAmount := 3000
	sendingCoinsInfo := entity.SendingCoins{recipient, transferAmount}

	err = service.TransferCoins(ctx, &senderUserBeforeTransferCoins.UUID, &sendingCoinsInfo)
	assert.ErrorIs(t, err, errorsx.ErrNotEnoughMoney)

	var senderUserAfterTransferCoins entity.User
	err = db.QueryRow(ctx, findUserPath, sender).Scan(&senderUserAfterTransferCoins.UUID,
		&senderUserAfterTransferCoins.Username,
		&senderUserAfterTransferCoins.Balance)
	assert.NoError(t, err)

	var recipientUserAfterTransferCoins entity.User
	err = db.QueryRow(ctx, findUserPath, recipient).Scan(&recipientUserAfterTransferCoins.UUID,
		&recipientUserAfterTransferCoins.Username,
		&recipientUserAfterTransferCoins.Balance)
	assert.NoError(t, err)

	assert.Equal(t, senderUserBeforeTransferCoins.Balance, senderUserAfterTransferCoins.Balance)

	assert.Equal(t, recipientUserBeforeTransferCoins.Balance, recipientUserAfterTransferCoins.Balance)

	transactions := []entity.Transaction{}

	rows, err := db.Query(ctx, findTransactionsBySender, senderUserBeforeTransferCoins.UUID)

	for rows.Next() {
		var transaction entity.Transaction
		err = rows.Scan(&transaction.SenderUUID, &transaction.ReceiverUUID, &transaction.Amount)
		assert.NoError(t, err)
		transactions = append(transactions, transaction)
	}

	assert.Equal(t, 0, len(transactions))

}

func TestTransferCoins_TranslationToSelf(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := test_repository.SetupTestDB(t)
	logger := slog.Default()
	repo := repository.NewStorage(db, logger)
	tokenGen := mocks.NewMocktokenGenerator(ctrl)
	service := service.NewService(repo, logger, tokenGen)
	ctx := context.Background()

	sender := "tester"
	var senderUserBeforeTransferCoins entity.User

	err := db.QueryRow(ctx, findUserPath, sender).Scan(&senderUserBeforeTransferCoins.UUID,
		&senderUserBeforeTransferCoins.Username,
		&senderUserBeforeTransferCoins.Balance)
	assert.NoError(t, err)

	recipient := "tester"
	var recipientUserBeforeTransferCoins entity.User

	err = db.QueryRow(ctx, findUserPath, recipient).Scan(&recipientUserBeforeTransferCoins.UUID,
		&recipientUserBeforeTransferCoins.Username,
		&recipientUserBeforeTransferCoins.Balance)
	assert.NoError(t, err)

	transferAmount := 100
	sendingCoinsInfo := entity.SendingCoins{recipient, transferAmount}

	err = service.TransferCoins(ctx, &senderUserBeforeTransferCoins.UUID, &sendingCoinsInfo)
	assert.ErrorIs(t, err, errorsx.ErrForbiddenTransaction)

	var senderUserAfterTransferCoins entity.User
	err = db.QueryRow(ctx, findUserPath, sender).Scan(&senderUserAfterTransferCoins.UUID,
		&senderUserAfterTransferCoins.Username,
		&senderUserAfterTransferCoins.Balance)
	assert.NoError(t, err)

	assert.Equal(t, senderUserBeforeTransferCoins.Balance, senderUserAfterTransferCoins.Balance)

	transactions := []entity.Transaction{}
	rows, err := db.Query(ctx, findTransactionsBySender, senderUserBeforeTransferCoins.UUID)

	for rows.Next() {
		var transaction entity.Transaction
		err = rows.Scan(&transaction.SenderUUID, &transaction.ReceiverUUID, &transaction.Amount)
		assert.NoError(t, err)
		transactions = append(transactions, transaction)
	}

	assert.Equal(t, 0, len(transactions))

}

func TestTransferCoins_RecipientNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := test_repository.SetupTestDB(t)
	logger := slog.Default()
	repo := repository.NewStorage(db, logger)
	tokenGen := mocks.NewMocktokenGenerator(ctrl)
	service := service.NewService(repo, logger, tokenGen)
	ctx := context.Background()

	sender := "tester"
	var senderUserBeforeTransferCoins entity.User

	err := db.QueryRow(ctx, findUserPath, sender).Scan(&senderUserBeforeTransferCoins.UUID,
		&senderUserBeforeTransferCoins.Username,
		&senderUserBeforeTransferCoins.Balance)
	assert.NoError(t, err)

	recipient := "badTester"

	transferAmount := 100
	sendingCoinsInfo := entity.SendingCoins{recipient, transferAmount}

	err = service.TransferCoins(ctx, &senderUserBeforeTransferCoins.UUID, &sendingCoinsInfo)
	assert.ErrorIs(t, err, errorsx.ErrReceiverNotFound)

	var senderUserAfterTransferCoins entity.User
	err = db.QueryRow(ctx, findUserPath, sender).Scan(&senderUserAfterTransferCoins.UUID,
		&senderUserAfterTransferCoins.Username,
		&senderUserAfterTransferCoins.Balance)
	assert.NoError(t, err)

	assert.Equal(t, senderUserBeforeTransferCoins.Balance, senderUserAfterTransferCoins.Balance)

	transactions := []entity.Transaction{}
	rows, err := db.Query(ctx, findTransactionsBySender, senderUserBeforeTransferCoins.UUID)

	for rows.Next() {
		var transaction entity.Transaction
		err = rows.Scan(&transaction.SenderUUID, &transaction.ReceiverUUID, &transaction.Amount)
		assert.NoError(t, err)
		transactions = append(transactions, transaction)
	}

	assert.Equal(t, 0, len(transactions))
}

func TestTransferCoins_SenderNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := test_repository.SetupTestDB(t)
	logger := slog.Default()
	repo := repository.NewStorage(db, logger)
	tokenGen := mocks.NewMocktokenGenerator(ctrl)
	service := service.NewService(repo, logger, tokenGen)
	ctx := context.Background()

	senderUUID := "268c15c5-cee5-44b8-8729-38ddb2d4f682"
	var senderUserBeforeTransferCoins entity.User

	recipient := "testerAvito"
	var recipientUserBeforeTransferCoins entity.User

	err := db.QueryRow(ctx, findUserPath, recipient).Scan(&recipientUserBeforeTransferCoins.UUID,
		&recipientUserBeforeTransferCoins.Username,
		&recipientUserBeforeTransferCoins.Balance)
	assert.NoError(t, err)

	transferAmount := 100
	sendingCoinsInfo := entity.SendingCoins{recipient, transferAmount}

	err = service.TransferCoins(ctx, &senderUUID, &sendingCoinsInfo)
	assert.ErrorIs(t, err, errorsx.ErrUnknownUser)

	var recipientUserAfterTransferCoins entity.User
	err = db.QueryRow(ctx, findUserPath, recipient).Scan(&recipientUserAfterTransferCoins.UUID,
		&recipientUserAfterTransferCoins.Username,
		&recipientUserAfterTransferCoins.Balance)
	assert.NoError(t, err)

	assert.Equal(t, recipientUserAfterTransferCoins.Balance, recipientUserBeforeTransferCoins.Balance)

	transactions := []entity.Transaction{}

	rows, err := db.Query(ctx, findTransactionsByRecipient, senderUserBeforeTransferCoins.UUID)

	for rows.Next() {
		var transaction entity.Transaction
		err = rows.Scan(&transaction.SenderUUID, &transaction.ReceiverUUID, &transaction.Amount)
		assert.NoError(t, err)
		transactions = append(transactions, transaction)
	}

	assert.Equal(t, 0, len(transactions))
}

func TestTransferCoins_WrongUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := test_repository.SetupTestDB(t)
	logger := slog.Default()
	repo := repository.NewStorage(db, logger)
	tokenGen := mocks.NewMocktokenGenerator(ctrl)
	service := service.NewService(repo, logger, tokenGen)
	ctx := context.Background()

	senderUUID := "badTester"
	var senderUserBeforeTransferCoins entity.User

	recipient := "testerAvito"
	var recipientUserBeforeTransferCoins entity.User

	err := db.QueryRow(ctx, findUserPath, recipient).Scan(&recipientUserBeforeTransferCoins.UUID,
		&recipientUserBeforeTransferCoins.Username,
		&recipientUserBeforeTransferCoins.Balance)
	assert.NoError(t, err)

	transferAmount := 100
	sendingCoinsInfo := entity.SendingCoins{recipient, transferAmount}

	err = service.TransferCoins(ctx, &senderUUID, &sendingCoinsInfo)
	assert.ErrorIs(t, err, errorsx.ErrWrongUUID)

	var recipientUserAfterTransferCoins entity.User
	err = db.QueryRow(ctx, findUserPath, recipient).Scan(&recipientUserAfterTransferCoins.UUID,
		&recipientUserAfterTransferCoins.Username,
		&recipientUserAfterTransferCoins.Balance)
	assert.NoError(t, err)

	assert.Equal(t, recipientUserAfterTransferCoins.Balance, recipientUserBeforeTransferCoins.Balance)

	transactions := []entity.Transaction{}

	rows, err := db.Query(ctx, findTransactionsByRecipient, senderUserBeforeTransferCoins.UUID)

	for rows.Next() {
		var transaction entity.Transaction
		err = rows.Scan(&transaction.SenderUUID, &transaction.ReceiverUUID, &transaction.Amount)
		assert.NoError(t, err)
		transactions = append(transactions, transaction)
	}

	assert.Equal(t, 0, len(transactions))
}
