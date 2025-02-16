package app

import (
	"context"
	"fmt"
	"github.com/ZhdanovichVlad/go_final_project/internal/pkg/jwttoken"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
	"time"

	"github.com/ZhdanovichVlad/go_final_project/internal/config"
	"github.com/ZhdanovichVlad/go_final_project/internal/controllers"
	repository "github.com/ZhdanovichVlad/go_final_project/internal/repository"
	"github.com/ZhdanovichVlad/go_final_project/internal/service"
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

	db, err := pgxpool.New(ctxDbOpen, config.PgDSN)

	if err != nil {
		logger.Error("error opening database", slog.String("err", err.Error()))
		os.Exit(exitCodeError)
	}

	defer db.Close()

	repository := repository.NewStorage(db, logger)
	generator := jwttoken.NewJwtTokenGenerator()
	service := service.NewService(repository, logger, generator)
	handlers := controllers.NewHandlers(service, logger)

	appRouter := NewRouter()
	appRouter.RegisterHandlers(handlers)
	host := fmt.Sprintf("0.0.0.0:%s", config.ServerPort)

	err = appRouter.Run(host)
	if err != nil {
		logger.Error("error running server", slog.String("err", err.Error()))
		os.Exit(exitCodeError)
	}

}
