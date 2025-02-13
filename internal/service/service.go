package service

import (
	"context"
	"github.com/ZhdanovichVlad/go_final_project/internal/entity"
	"github.com/ZhdanovichVlad/go_final_project/internal/pkg/errorsx"
	"github.com/ZhdanovichVlad/go_final_project/internal/pkg/jwt_token"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

type repository interface {
	SaveUser(ctx context.Context, user *entity.User) (*string, error)
	FindUser(ctx context.Context, userInfo *entity.User) (*entity.User, error)
	ExistsUser(ctx context.Context, userUuid *string) (*bool, error)

	BuyMerch(ctx context.Context, userUuid *string, merchInfo *entity.Merch) error
	TransferCoins(ctx context.Context, userUuid *string, receiverUuid *string, amount *int) error
	GetUserInfo(ctx context.Context, userUUID *string) (*entity.UserInfo, error)
}

type service struct {
	repo   repository
	logger *slog.Logger
}

func NewService(repo repository, logger *slog.Logger) *service {
	return &service{repo: repo, logger: logger}
}

func (s *service) Login(ctx context.Context, userDto *entity.UserDTO) (*entity.JwtToken, error) {
	user := &entity.User{}
	user.Username = userDto.Username
	user, err := s.repo.FindUser(ctx, user)
	if err != nil {
		return nil, err
	}
	if user == nil {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), 10)
		if err != nil {
			s.logger.Error("failed to generate password hash",
				slog.String("method", "service.SaveUser"),
				slog.String("error", err.Error()))
			return nil, errorsx.ServiceError
		}
		user = &entity.User{
			Username:     userDto.Username,
			PasswordHash: string(passwordHash), // Преобразуем в строку
		}
		userUuid, err := s.repo.SaveUser(ctx, user)
		if err != nil {
			return nil, err
		}

		accessToken, err := jwt_token.GenerateToken(*userUuid)
		if err != nil {
			s.logger.Error("failed to generate access token",
				slog.String("method", "service.SaveUser"),
				slog.String("error", err.Error()))
			return nil, err
		}

		return &entity.JwtToken{accessToken}, nil

	} else {
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(userDto.Password))
		if err != nil {
			return nil, errorsx.InvalidPassword
		}
		accessToken, err := jwt_token.GenerateToken(user.Uuid)
		if err != nil {
			s.logger.Error("failed to generate access token",
				slog.String("method", "service.SaveUser"),
				slog.String("error", err.Error()))
			return nil, err
		}

		return &entity.JwtToken{accessToken}, nil
	}
}

func (s *service) BuyMerch(ctx context.Context, userUuid *string, merchInfo *entity.Merch) error {
	exist, err := s.repo.ExistsUser(ctx, userUuid)
	if err != nil {
		return err
	}
	if !*exist {
		return errorsx.UnknownUser
	}

	err = s.repo.BuyMerch(ctx, userUuid, merchInfo)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) TransferCoins(ctx context.Context, userUuid *string, sendingCoinsInfo *entity.SendingCoins) error {

	exist, err := s.repo.ExistsUser(ctx, userUuid)
	if err != nil {
		return err
	}
	if !*exist {
		return errorsx.UnknownUser
	}

	receiver := &entity.User{}
	receiver.Username = sendingCoinsInfo.User
	receiver, err = s.repo.FindUser(ctx, receiver)
	if err != nil {
		return err
	}
	if receiver == nil {
		return errorsx.ReceiverNotFound
	}
	if *userUuid == receiver.Uuid {
		return errorsx.ForbiddenTransaction
	}

	err = s.repo.TransferCoins(ctx, userUuid, &receiver.Uuid, &sendingCoinsInfo.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetUserInfo(ctx context.Context, userUuid *string) (*entity.UserInfo, error) {
	exist, err := s.repo.ExistsUser(ctx, userUuid)
	if err != nil {
		return nil, err
	}
	if !*exist {
		return nil, errorsx.UnknownUser
	}

	info, err := s.repo.GetUserInfo(ctx, userUuid)
	if err != nil {
		return nil, err
	}

	return info, nil
}
