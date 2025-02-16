package servicetest_test

import (
	"context"
	"errors"
	"github.com/ZhdanovichVlad/go_final_project/internal/service"
	"github.com/ZhdanovichVlad/go_final_project/tests/unit/servicetest/mocks"
	"log/slog"
	"testing"

	"github.com/ZhdanovichVlad/go_final_project/internal/entity"
	"github.com/ZhdanovichVlad/go_final_project/internal/pkg/errorsx"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTransferCoins_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	mockTokenGenerator := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(mockRepo, slog.Default(), mockTokenGenerator)

	userUUID := "a831f52d-9de2-4af1-8677-4f3d1226fed2"
	receiverUUID := "a831f52d-9de2-4af1-8677-4f3d1226fyf5"
	sendingCoinsInfo := &entity.SendingCoins{User: "receiver", Amount: 100}
	receiver := &entity.User{UUID: receiverUUID, Username: "receiver"}

	answer := true
	mockRepo.EXPECT().
		ExistsUser(gomock.Any(), &userUUID).
		Return(&answer, nil)

	mockRepo.EXPECT().
		FindUser(gomock.Any(), gomock.Any()).
		Return(receiver, nil)

	mockRepo.EXPECT().
		TransferCoins(gomock.Any(), &userUUID, &receiverUUID, &sendingCoinsInfo.Amount).
		Return(nil)

	err := service.TransferCoins(context.Background(), &userUUID, sendingCoinsInfo)

	assert.NoError(t, err)
}

func TestTransferCoins_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	mockTokenGenerator := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(mockRepo, slog.Default(), mockTokenGenerator)

	userUUID := "a831f52d-9de2-4af1-8677-4f3d1226fed2"
	sendingCoinsInfo := &entity.SendingCoins{User: "receiver", Amount: 100}

	answer := false
	mockRepo.EXPECT().
		ExistsUser(gomock.Any(), &userUUID).
		Return(&answer, nil)

	err := service.TransferCoins(context.Background(), &userUUID, sendingCoinsInfo)

	assert.ErrorIs(t, err, errorsx.ErrUnknownUser)
}

func TestTransferCoins_ReceiverNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	mockTokenGenerator := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(mockRepo, slog.Default(), mockTokenGenerator)

	userUUID := "a831f52d-9de2-4af1-8677-4f3d1226fed2"
	sendingCoinsInfo := &entity.SendingCoins{User: "receiver", Amount: 100}

	answer := true
	mockRepo.EXPECT().
		ExistsUser(gomock.Any(), &userUUID).
		Return(&answer, nil)

	mockRepo.EXPECT().
		FindUser(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	err := service.TransferCoins(context.Background(), &userUUID, sendingCoinsInfo)

	assert.ErrorIs(t, err, errorsx.ErrReceiverNotFound)
}

func TestTransferCoins_ForbiddenTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	mockTokenGenerator := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(mockRepo, slog.Default(), mockTokenGenerator)

	userUUID := "a831f52d-9de2-4af1-8677-4f3d1226fed2"
	sendingCoinsInfo := &entity.SendingCoins{User: "receiver", Amount: 100}
	receiver := &entity.User{UUID: userUUID, Username: "receiver"} // Получатель == отправитель

	answer := true
	mockRepo.EXPECT().
		ExistsUser(gomock.Any(), &userUUID).
		Return(&answer, nil)

	mockRepo.EXPECT().
		FindUser(gomock.Any(), gomock.Any()).
		Return(receiver, nil)

	err := service.TransferCoins(context.Background(), &userUUID, sendingCoinsInfo)

	assert.ErrorIs(t, err, errorsx.ErrForbiddenTransaction)
}

func TestTransferCoins_TransferError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	mockTokenGenerator := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(mockRepo, slog.Default(), mockTokenGenerator)

	userUUID := "a831f52d-9de2-4af1-8677-4f3d1226fed2"
	receiverUUID := "a831f52d-9de2-4af1-8677-4f3d1226fyg5"
	sendingCoinsInfo := &entity.SendingCoins{User: "receiver", Amount: 100}
	receiver := &entity.User{UUID: receiverUUID, Username: "receiver"}

	answer := true
	mockRepo.EXPECT().
		ExistsUser(gomock.Any(), &userUUID).
		Return(&answer, nil)

	mockRepo.EXPECT().
		FindUser(gomock.Any(), gomock.Any()).
		Return(receiver, nil)

	expectedErr := errors.New("transfer failed")
	mockRepo.EXPECT().
		TransferCoins(gomock.Any(), &userUUID, &receiverUUID, &sendingCoinsInfo.Amount).
		Return(expectedErr)

	err := service.TransferCoins(context.Background(), &userUUID, sendingCoinsInfo)

	assert.ErrorIs(t, err, expectedErr)
}

func TestTransferCoins_WrongUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	mockTokenGenerator := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(mockRepo, slog.Default(), mockTokenGenerator)

	userUUID := "wrongUUID"
	sendingCoinsInfo := &entity.SendingCoins{User: "receiver", Amount: 100}

	err := service.TransferCoins(context.Background(), &userUUID, sendingCoinsInfo)

	assert.ErrorIs(t, err, errorsx.ErrWrongUUID)
}
