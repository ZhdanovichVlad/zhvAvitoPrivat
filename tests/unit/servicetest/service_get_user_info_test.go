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

func TestGetUserInfo_UserNotFound(t *testing.T) {

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

func TestGetUserInfo_DBError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	mockTokenGenerator := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(mockRepo, slog.Default(), mockTokenGenerator)

	ctx := context.Background()
	userUUID := "testuser"

	answTrue := true
	mockRepo.EXPECT().
		ExistsUser(ctx, gomock.Any()).
		Return(&answTrue, nil)

	mockRepo.EXPECT().GetUserInfo(ctx, &userUUID).
		Return(nil, errorsx.ErrDB)

	_, err := service.GetUserInfo(ctx, &userUUID)

	assert.ErrorIs(t, err, errorsx.ErrDB)
}

func TestServiceInfo_SuccessfulRequest(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	mockTokenGenerator := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(mockRepo, slog.Default(), mockTokenGenerator)

	ctx := context.Background()
	userUUID := "testuser"

	info := &entity.UserInfo{Coins: 50}

	answTrue := true
	mockRepo.EXPECT().
		ExistsUser(ctx, gomock.Any()).
		Return(&answTrue, nil)

	mockRepo.EXPECT().GetUserInfo(ctx, &userUUID).
		Return(info, nil)

	infoAnsw, err := service.GetUserInfo(ctx, &userUUID)

	assert.NoError(t, err)
	assert.Equal(t, info, infoAnsw)
}
