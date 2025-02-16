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

func TestBuyMerch_Successful(t *testing.T) {
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
	assert.NoError(t, err)

	merchInfo := entity.Merch{Name: "pen"}

	err = db.QueryRow(ctx, findMerchPrice, merchInfo.Name).Scan(&merchInfo.UUID, &merchInfo.Price)
	assert.NoError(t, err)

	err = service.BuyMerch(ctx, &userBeforeBuyMerch.UUID, &merchInfo)
	assert.NoError(t, err)

	items := []entity.Item{}

	rows, err := db.Query(ctx, findMerchInInventory, &userBeforeBuyMerch.UUID)
	assert.NoError(t, err)

	defer rows.Close()

	for rows.Next() {
		var item entity.Item
		err = rows.Scan(&item.Type, &item.Quantity)
		assert.NoError(t, err)
		items = append(items, item)
	}

	assert.Equal(t, len(items), 1)
	assert.Equal(t, merchInfo.UUID, items[0].Type)

	var userAfterBuyMerch entity.User

	err = db.QueryRow(ctx, findUserPath, testerNameInDB).Scan(&userAfterBuyMerch.UUID,
		&userAfterBuyMerch.Username,
		&userAfterBuyMerch.Balance)
	assert.NoError(t, err)

	balanceDifference := userBeforeBuyMerch.Balance - userAfterBuyMerch.Balance

	assert.Equal(t, balanceDifference, merchInfo.Price)
}

func TestBuyMerch_MerchNotFound(t *testing.T) {
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
	assert.NoError(t, err)

	merchInfo := entity.Merch{Name: "parker"}

	err = service.BuyMerch(ctx, &userBeforeBuyMerch.UUID, &merchInfo)
	assert.ErrorIs(t, err, errorsx.ErrItemNotFound)

	items := []entity.Item{}

	rows, err := db.Query(ctx, findMerchInInventory, &userBeforeBuyMerch.UUID)
	assert.NoError(t, err)

	defer rows.Close()

	for rows.Next() {
		var item entity.Item
		err = rows.Scan(&item.Type, &item.Quantity)
		assert.NoError(t, err)
		items = append(items, item)
	}

	assert.Equal(t, len(items), 0)

	var userAfterBuyMerch entity.User

	err = db.QueryRow(ctx, findUserPath, testerNameInDB).Scan(&userAfterBuyMerch.UUID,
		&userAfterBuyMerch.Username,
		&userAfterBuyMerch.Balance)
	assert.NoError(t, err)

	balanceDifference := userBeforeBuyMerch.Balance - userAfterBuyMerch.Balance

	assert.Equal(t, balanceDifference, 0)
}

func TestBuyMerch_WrongUUID(t *testing.T) {
	ctrl := gomock.NewController(t)

	db := test_repository.SetupTestDB(t)
	logger := slog.Default()
	repo := repository.NewStorage(db, logger)

	tokenGen := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(repo, logger, tokenGen)

	ctx := context.Background()

	testerUUID := "badTesterUUID"

	merchInfo := entity.Merch{Name: "pen"}

	err := service.BuyMerch(ctx, &testerUUID, &merchInfo)
	if err != nil {
		t.Log("err", err.Error())
	}
	assert.ErrorIs(t, err, errorsx.ErrWrongUUID)

}

func TestBuyMerch_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	db := test_repository.SetupTestDB(t)
	logger := slog.Default()
	repo := repository.NewStorage(db, logger)

	tokenGen := mocks.NewMocktokenGenerator(ctrl)

	service := service.NewService(repo, logger, tokenGen)

	ctx := context.Background()

	testerUUID := "badTesterUUID"

	merchInfo := entity.Merch{Name: "pen"}

	err := service.BuyMerch(ctx, &testerUUID, &merchInfo)
	if err != nil {
		t.Log("err", err.Error())
	}
	assert.ErrorIs(t, err, errorsx.ErrWrongUUID)
}
