package servicetest_test

import (
	"context"
	"github.com/ZhdanovichVlad/go_final_project/internal/entity"
	"github.com/ZhdanovichVlad/go_final_project/internal/pkg/errorsx"
	"github.com/ZhdanovichVlad/go_final_project/internal/service"
	"github.com/ZhdanovichVlad/go_final_project/tests/unit/servicetest/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestServiceBuyMerch_UserNotFound(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	mockTokenGenerator := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(mockRepo, slog.Default(), mockTokenGenerator)

	ctx := context.Background()
	userUUID := "a831f52d-9de2-4af1-8677-4f3d1226fe22"

	merchInfo := &entity.Merch{Name: "pen"}

	answFalse := false
	mockRepo.EXPECT().
		ExistsUser(ctx, &userUUID).
		Return(&answFalse, nil)

	err := service.BuyMerch(ctx, &userUUID, merchInfo)
	assert.ErrorIs(t, err, errorsx.ErrUnknownUser)
}

func TestServiceBuyMerch_MerchNotFound(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	mockTokenGenerator := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(mockRepo, slog.Default(), mockTokenGenerator)

	ctx := context.Background()
	userUUID := "a831f52d-9de2-4af1-8677-4f3d1226fed2"

	merchInfo := &entity.Merch{Name: "pen"}

	answTrue := true
	mockRepo.EXPECT().
		ExistsUser(ctx, &userUUID).
		Return(&answTrue, nil)

	mockRepo.EXPECT().
		BuyMerch(ctx, gomock.Any(), gomock.Any()).
		Return(errorsx.ErrItemNotFound)

	err := service.BuyMerch(ctx, &userUUID, merchInfo)
	assert.ErrorIs(t, err, errorsx.ErrItemNotFound)
}

func TestServiceBuyMerch_Successful(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	mockTokenGenerator := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(mockRepo, slog.Default(), mockTokenGenerator)

	ctx := context.Background()
	userName := "a831f52d-9de2-4af1-8677-4f3d1226fed2"

	merchInfo := &entity.Merch{Name: "pen"}

	answTrue := true
	mockRepo.EXPECT().
		ExistsUser(ctx, &userName).
		Return(&answTrue, nil)

	mockRepo.EXPECT().
		BuyMerch(ctx, &userName, merchInfo).
		Return(nil)

	err := service.BuyMerch(ctx, &userName, merchInfo)

	assert.NoError(t, err)
}

func TestServiceBuyMerch_WrongUUID(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	mockTokenGenerator := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(mockRepo, slog.Default(), mockTokenGenerator)

	ctx := context.Background()
	userUUID := "wrongUUID"

	merchInfo := &entity.Merch{Name: "pen"}

	err := service.BuyMerch(ctx, &userUUID, merchInfo)
	assert.ErrorIs(t, err, errorsx.ErrWrongUUID)
}
