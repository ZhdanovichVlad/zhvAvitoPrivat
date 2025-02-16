package integration

import (
	"context"
	"github.com/ZhdanovichVlad/go_final_project/internal/entity"
	"github.com/ZhdanovichVlad/go_final_project/internal/repository"
	"github.com/ZhdanovichVlad/go_final_project/internal/service"
	"github.com/ZhdanovichVlad/go_final_project/tests/integration/test_repository"
	"github.com/ZhdanovichVlad/go_final_project/tests/unit/servicetest/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

const (
	findUserPath         = "SELECT uuid, username, balance FROM users WHERE username = $1"
	findMerchInInventory = "SELECT merchandise_uuid, quantity FROM owned_inventory WHERE  user_uuid = $1"
	findMerchPrice       = "SELECT uuid, price FROM merchandise WHERE name = $1"
)

func TestBuyMerchIntegration_(t *testing.T) {
	ctrl := gomock.NewController(t)

	db := test_repository.SetupTestDB(t)
	logger := slog.Default()
	repo := repository.NewStorage(db, logger)

	tokenGen := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(repo, logger, tokenGen)

	ctx := context.Background()

	testerNameInDB := "tester"
	var userBeforeBuyMerch entity.User

	err := db.QueryRow(ctx, findUserPath, testerNameInDB).Scan(&userBeforeBuyMerch.UUID,
		&userBeforeBuyMerch.Username,
		&userBeforeBuyMerch.Balance)
	if err != nil {
		t.Error("Cannot find user before buy merch err:", err.Error())
	}

	merchInfo := entity.Merch{Name: "pen"}

	err = db.QueryRow(ctx, findMerchPrice, merchInfo.Name).Scan(&merchInfo.UUID, &merchInfo.Price)
	if err != nil {
		t.Error("Cannot find merch price err:", err.Error())
	}

	err = service.BuyMerch(ctx, &userBeforeBuyMerch.UUID, &merchInfo)
	if err != nil {
		t.Error("Cannot ByMerch")
	}

	items := []entity.Item{}

	rows, err := db.Query(ctx, findMerchInInventory, &userBeforeBuyMerch.UUID)
	if err != nil {
		t.Error("field to get rows :", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.Item
		err = rows.Scan(&item.Type, &item.Quantity)
		if err != nil {
			t.Error("field to scan row :", err.Error())
		}
		items = append(items, item)
	}

	assert.Equal(t, len(items), 1)
	assert.Equal(t, merchInfo.UUID, items[0].Type)

	var userAfterBuyMerch entity.User

	err = db.QueryRow(ctx, findUserPath, testerNameInDB).Scan(&userAfterBuyMerch.UUID,
		&userAfterBuyMerch.Username,
		&userAfterBuyMerch.Balance)
	if err != nil {
		t.Error("Cannot find user before buy merch err:", err.Error())
	}

	balanceDifference := userBeforeBuyMerch.Balance - userAfterBuyMerch.Balance

	assert.Equal(t, balanceDifference, merchInfo.Price)
	assert.True(t, true)
}

//func TestTransferCoinsIntegration(t *testing.T) {
//	ctrl := gomock.NewController(t)
//
//	db := test_repository.SetupTestDB(t)
//	logger := slog.Default()
//	repo := repository.NewStorage(db, logger)
//
//	tokenGen := mocks.NewMocktokenGenerator(ctrl)
//
//	service := service.NewService(repo, logger, tokenGen)
//
//	// Создаём пользователей, вызываем service_test.TransferCoins и проверяем БД
//	// ...
//	assert.True(t, true)
//}
