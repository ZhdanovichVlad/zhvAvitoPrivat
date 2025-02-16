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

	userUUID := "user-123"
	receiverUUID := "user-456"
	sendingCoinsInfo := &entity.SendingCoins{User: "receiver", Amount: 100}
	receiver := &entity.User{UUID: receiverUUID, Username: "receiver"}

	// Ожидаемые вызовы
	mockRepo.EXPECT().
		ExistsUser(gomock.Any(), &userUUID).
		Return(ptrBool(true), nil) // Пользователь существует

	mockRepo.EXPECT().
		FindUser(gomock.Any(), gomock.Any()).
		Return(receiver, nil) // Получатель найден

	mockRepo.EXPECT().
		TransferCoins(gomock.Any(), &userUUID, &receiverUUID, &sendingCoinsInfo.Amount).
		Return(nil) // Перевод успешен

	// Вызов тестируемой функции
	err := service.TransferCoins(context.Background(), &userUUID, sendingCoinsInfo)

	// Проверка результата
	assert.NoError(t, err)
}

func TestTransferCoins_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	mockTokenGenerator := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(mockRepo, slog.Default(), mockTokenGenerator)

	userUUID := "user-123"
	sendingCoinsInfo := &entity.SendingCoins{User: "receiver", Amount: 100}

	mockRepo.EXPECT().
		ExistsUser(gomock.Any(), &userUUID).
		Return(ptrBool(false), nil)

	err := service.TransferCoins(context.Background(), &userUUID, sendingCoinsInfo)

	assert.ErrorIs(t, err, errorsx.ErrUnknownUser)
}

func TestTransferCoins_ReceiverNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	mockTokenGenerator := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(mockRepo, slog.Default(), mockTokenGenerator)

	userUUID := "user-123"
	sendingCoinsInfo := &entity.SendingCoins{User: "receiver", Amount: 100}

	mockRepo.EXPECT().
		ExistsUser(gomock.Any(), &userUUID).
		Return(ptrBool(true), nil) // Пользователь существует

	mockRepo.EXPECT().
		FindUser(gomock.Any(), gomock.Any()).
		Return(nil, nil) // Получатель не найден

	// Вызов тестируемой функции
	err := service.TransferCoins(context.Background(), &userUUID, sendingCoinsInfo)

	// Проверка результата
	assert.ErrorIs(t, err, errorsx.ErrReceiverNotFound)
}

func TestTransferCoins_ForbiddenTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	mockTokenGenerator := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(mockRepo, slog.Default(), mockTokenGenerator)

	userUUID := "user-123"
	sendingCoinsInfo := &entity.SendingCoins{User: "receiver", Amount: 100}
	receiver := &entity.User{UUID: userUUID, Username: "receiver"} // Получатель == отправитель

	// Ожидаемые вызовы
	mockRepo.EXPECT().
		ExistsUser(gomock.Any(), &userUUID).
		Return(ptrBool(true), nil) // Пользователь существует

	mockRepo.EXPECT().
		FindUser(gomock.Any(), gomock.Any()).
		Return(receiver, nil) // Получатель найден

	// Вызов тестируемой функции
	err := service.TransferCoins(context.Background(), &userUUID, sendingCoinsInfo)

	// Проверка результата
	assert.ErrorIs(t, err, errorsx.ErrForbiddenTransaction)
}

func TestTransferCoins_TransferError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	mockTokenGenerator := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(mockRepo, slog.Default(), mockTokenGenerator)

	userUUID := "user-123"
	receiverUUID := "user-456"
	sendingCoinsInfo := &entity.SendingCoins{User: "receiver", Amount: 100}
	receiver := &entity.User{UUID: receiverUUID, Username: "receiver"}

	// Ожидаемые вызовы
	mockRepo.EXPECT().
		ExistsUser(gomock.Any(), &userUUID).
		Return(ptrBool(true), nil) // Пользователь существует

	mockRepo.EXPECT().
		FindUser(gomock.Any(), gomock.Any()).
		Return(receiver, nil) // Получатель найден

	expectedErr := errors.New("transfer failed")
	mockRepo.EXPECT().
		TransferCoins(gomock.Any(), &userUUID, &receiverUUID, &sendingCoinsInfo.Amount).
		Return(expectedErr) // Ошибка при переводе

	// Вызов тестируемой функции
	err := service.TransferCoins(context.Background(), &userUUID, sendingCoinsInfo)

	// Проверка результата
	assert.ErrorIs(t, err, expectedErr)
}

// Вспомогательная функция для создания указателя на bool
func ptrBool(b bool) *bool {
	return &b
}
