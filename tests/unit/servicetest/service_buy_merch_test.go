package servicetest

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
	userName := "testuser"

	merchInfo := &entity.Merch{Name: "pen"}

	answFalse := false
	mockRepo.EXPECT().
		ExistsUser(ctx, gomock.Any()).
		Return(&answFalse, nil)

	err := service.BuyMerch(ctx, &userName, merchInfo)
	assert.ErrorIs(t, err, errorsx.ErrUnknownUser)
}

func TestServiceBuyMerch_ItemNotFound(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	mockTokenGenerator := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(mockRepo, slog.Default(), mockTokenGenerator)

	ctx := context.Background()
	userName := "testuser"

	merchInfo := &entity.Merch{Name: "pen"}

	answTrue := true
	mockRepo.EXPECT().
		ExistsUser(ctx, gomock.Any()).
		Return(&answTrue, nil)

	mockRepo.EXPECT().
		BuyMerch(ctx, gomock.Any(), gomock.Any()).
		Return(errorsx.ErrItemNotFound)

	err := service.BuyMerch(ctx, &userName, merchInfo)
	assert.ErrorIs(t, err, errorsx.ErrItemNotFound)
}

func TestServiceBuyMerch_UserFound(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	mockTokenGenerator := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(mockRepo, slog.Default(), mockTokenGenerator)

	ctx := context.Background()
	userName := "testuser"

	merchInfo := &entity.Merch{Name: "pen"}

	answTrue := true
	mockRepo.EXPECT().
		ExistsUser(ctx, gomock.Any()).
		Return(&answTrue, nil)

	mockRepo.EXPECT().
		BuyMerch(ctx, &userName, merchInfo).
		Return(nil)

	err := service.BuyMerch(ctx, &userName, merchInfo)

	assert.NoError(t, err)
}

func TestServiceBuyMerch_ItemFound(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	mockTokenGenerator := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(mockRepo, slog.Default(), mockTokenGenerator)

	ctx := context.Background()
	userName := "testuser"

	merchInfo := &entity.Merch{Name: "pen"}

	answTrue := true
	mockRepo.EXPECT().
		ExistsUser(ctx, gomock.Any()).
		Return(&answTrue, nil)

	mockRepo.EXPECT().
		BuyMerch(ctx, &userName, merchInfo).
		Return(nil)

	err := service.BuyMerch(ctx, &userName, merchInfo)

	assert.NoError(t, err)
}
