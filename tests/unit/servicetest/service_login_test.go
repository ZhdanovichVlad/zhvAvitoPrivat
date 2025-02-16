package servicetest_test

import (
	"context"
	service "github.com/ZhdanovichVlad/go_final_project/internal/service"
	"github.com/ZhdanovichVlad/go_final_project/tests/unit/servicetest/mocks"
	"log/slog"
	"testing"

	"github.com/ZhdanovichVlad/go_final_project/internal/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestServiceLogin_UserNotFound(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	mockTokenGenerator := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(mockRepo, slog.Default(), mockTokenGenerator)

	ctx := context.Background()
	userDto := &entity.UserDTO{
		Username: "testuser",
		Password: "testpassword",
	}

	mockRepo.EXPECT().
		FindUser(ctx, gomock.Any()).
		Return(nil, nil)

	mockRepo.EXPECT().
		SaveUser(ctx, gomock.Any()).
		Return(func() *string { s := "user-uuid"; return &s }(), nil)

	mockTokenGenerator.EXPECT().
		GenerateToken("user-uuid").
		Return("access-token", nil)

	result, err := service.Login(ctx, userDto)

	assert.NoError(t, err)
	assert.Equal(t, "access-token", result.Token)
}

func TestService_Login_UserFound(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepository(ctrl)
	mockTokenGenerator := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(mockRepo, slog.Default(), mockTokenGenerator)

	ctx := context.Background()
	userDto := &entity.UserDTO{
		Username: "testuser",
		Password: "testpassword",
	}

	existingUser := &entity.User{
		UUID:         "user-uuid",
		Username:     "testuser",
		PasswordHash: "$2a$08$nzrvW0ExCfom8pMhvhTOI.purmNWeIZK0w.N3omQifOwJYrkahK4q",
	}
	mockRepo.EXPECT().
		FindUser(ctx, gomock.Any()).
		Return(existingUser, nil)

	mockTokenGenerator.EXPECT().
		GenerateToken("user-uuid").
		Return("access-token", nil)

	result, err := service.Login(ctx, userDto)

	assert.NoError(t, err)
	assert.Equal(t, "access-token", result.Token)
}
