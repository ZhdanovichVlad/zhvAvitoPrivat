package app

import (
	"context"
	"fmt"
	"github.com/ZhdanovichVlad/go_final_project/internal/config"
	"github.com/ZhdanovichVlad/go_final_project/internal/controllers"
	repository "github.com/ZhdanovichVlad/go_final_project/internal/repository"
	"github.com/ZhdanovichVlad/go_final_project/internal/service"
	"github.com/jackc/pgx/v4"
	"log/slog"
	"os"
	"time"
)

const (
	exitCodeError = 1
)

func Run() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	config, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config", slog.String("err", err.Error()))
		os.Exit(exitCodeError)
	}

	ctxDbOpen, cancelDbOpen := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelDbOpen()

	db, err := pgx.Connect(ctxDbOpen, config.PgDSN)
	if err != nil {
		logger.Error("error opening database", slog.String("err", err.Error()))
		os.Exit(exitCodeError)
	}

	ctxDbClose, cancelDbClose := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelDbClose()
	defer db.Close(ctxDbClose)

	repository := repository.NewStorage(db, logger)
	service := service.NewService(repository, logger)
	handlers := controllers.NewHandlers(service, logger)

	appRouter := NewRouter()
	appRouter.RegisterHandlers(handlers)
	host := fmt.Sprintf("0.0.0.0:%s", config.ServerPort)
	appRouter.Run(host)

}
