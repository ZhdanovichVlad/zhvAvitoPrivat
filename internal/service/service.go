package service

import (
	"context"
	"github.com/google/uuid"
	"log/slog"

	"github.com/ZhdanovichVlad/go_final_project/internal/entity"
	"github.com/ZhdanovichVlad/go_final_project/internal/pkg/errorsx"
	"golang.org/x/crypto/bcrypt"
)

type repository interface {
	SaveUser(ctx context.Context, user *entity.User) (*string, error)
	FindUser(ctx context.Context, userInfo *entity.User) (*entity.User, error)
	ExistsUser(ctx context.Context, userUUID *string) (*bool, error)

	BuyMerch(ctx context.Context, userUUID *string, merchInfo *entity.Merch) error
	TransferCoins(ctx context.Context, userUUID *string, receiverUUID *string, amount *int) error
	GetUserInfo(ctx context.Context, userUUID *string) (*entity.UserInfo, error)
}

type tokenGenerator interface {
	GenerateToken(userUUID string) (string, error)
}

type Service struct {
	repo           repository
	logger         *slog.Logger
	tokenGenerator tokenGenerator
}

func NewService(repo repository, logger *slog.Logger, generator tokenGenerator) *Service {
	return &Service{repo: repo, logger: logger, tokenGenerator: generator}
}

func (s *Service) Login(ctx context.Context, userDto *entity.UserDTO) (*entity.JwtToken, error) {
	user := &entity.User{}
	user.Username = userDto.Username

	user, err := s.repo.FindUser(ctx, user)
	if err != nil {
		return nil, err
	}

	var accessToken string
	if user == nil {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), 8)
		if err != nil {
			s.logger.Error("failed to generate password hash",
				slog.String("method", "service.SaveUser"),
				slog.String("error", err.Error()))
			return nil, errorsx.ErrService
		}
		user = &entity.User{
			Username:     userDto.Username,
			PasswordHash: string(passwordHash), // Преобразуем в строку
		}
		userUUID, err := s.repo.SaveUser(ctx, user)
		if err != nil {
			return nil, err
		}

		accessToken, err = s.tokenGenerator.GenerateToken(*userUUID)
		if err != nil {
			s.logger.Error("failed to generate access token",
				slog.String("method", "service.SaveUser"),
				slog.String("error", err.Error()))
			return nil, err
		}

	} else {
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(userDto.Password))
		if err != nil {
			return nil, errorsx.ErrInvalidPassword
		}
		accessToken, err = s.tokenGenerator.GenerateToken(user.UUID)
		if err != nil {
			s.logger.Error("failed to generate access token",
				slog.String("method", "service.SaveUser"),
				slog.String("error", err.Error()))
			return nil, err
		}

	}
	return &entity.JwtToken{Token: accessToken}, nil
}

func (s *Service) BuyMerch(ctx context.Context, userUUID *string, merchInfo *entity.Merch) error {

	err := uuid.Validate(*userUUID)
	if err != nil {
		return errorsx.ErrWrongUUID
	}

	exist, err := s.repo.ExistsUser(ctx, userUUID)
	if err != nil {
		return err
	}
	if !*exist {
		return errorsx.ErrUnknownUser
	}

	err = s.repo.BuyMerch(ctx, userUUID, merchInfo)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) TransferCoins(ctx context.Context, userUUID *string, sendingCoinsInfo *entity.SendingCoins) error {

	err := uuid.Validate(*userUUID)
	if err != nil {
		return errorsx.ErrWrongUUID
	}

	exist, err := s.repo.ExistsUser(ctx, userUUID)
	if err != nil {
		return err
	}
	if !*exist {
		return errorsx.ErrUnknownUser
	}

	receiver := &entity.User{}
	receiver.Username = sendingCoinsInfo.User
	receiver, err = s.repo.FindUser(ctx, receiver)
	if err != nil {
		return err
	}
	if receiver == nil {
		return errorsx.ErrReceiverNotFound
	}
	if *userUUID == receiver.UUID {
		return errorsx.ErrForbiddenTransaction
	}

	err = s.repo.TransferCoins(ctx, userUUID, &receiver.UUID, &sendingCoinsInfo.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetUserInfo(ctx context.Context, userUUID *string) (*entity.UserInfo, error) {

	err := uuid.Validate(*userUUID)
	if err != nil {
		return nil, errorsx.ErrWrongUUID
	}

	exist, err := s.repo.ExistsUser(ctx, userUUID)
	if err != nil {
		return nil, err
	}
	if !*exist {
		return nil, errorsx.ErrUnknownUser
	}

	info, err := s.repo.GetUserInfo(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	return info, nil
}
