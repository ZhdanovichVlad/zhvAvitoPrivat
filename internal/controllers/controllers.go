package controllers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ZhdanovichVlad/go_final_project/internal/entity"
	"github.com/ZhdanovichVlad/go_final_project/internal/pkg/errorsx"

	"github.com/gin-gonic/gin"
)

type api interface {
	Login(ctx context.Context, userDto *entity.UserDTO) (*entity.JwtToken, error)
	BuyMerch(ctx context.Context, userUuid *string, item *entity.Merch) error
	TransferCoins(ctx context.Context, userUuid *string, sendingCoinsInfo *entity.SendingCoins) error
	GetUserInfo(ctx context.Context, userUuid *string) (*entity.UserInfo, error)
}

type Handlers struct {
	service api
	logger  *slog.Logger
}

func NewHandlers(api api, logger *slog.Logger) *Handlers {
	return &Handlers{api, logger}
}

func (h *Handlers) Login(c *gin.Context) {
	ctx := c.Request.Context()

	userDto := entity.UserDTO{}
	err := c.ShouldBindJSON(&userDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorsx.InvalidInput.Error()})
		return
	}

	token, err := h.service.Login(ctx, &userDto)
	if err != nil {
		switch {
		case errors.Is(err, errorsx.InvalidPassword):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case errors.Is(err, errorsx.DBError), errors.Is(err, errorsx.ServiceError):
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		default:
			h.logger.Error("unknown error service.Login",
				slog.String("method", "handler.Login"),
				slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorsx.UnknownError.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, token)
}

func (h *Handlers) BuyMerch(c *gin.Context) {
	ctx := c.Request.Context()
	userUuid, exists := c.Get("userUuid")
	if !exists {
		h.logger.Error("userUuid not found in context",
			slog.String("method", "handler.BuyMerch"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorsx.ServiceError.Error()})
		return
	}
	userUuidStr, ok := userUuid.(string)
	if !ok {
		h.logger.Error("userUuid cannot parse to string",
			slog.String("method", "handler.BuyMerch"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorsx.ServiceError.Error()})
		return
	}

	item := c.Param("item")
	if item == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorsx.InvalidInput.Error()})
		return
	}
	var merch entity.Merch
	merch.Name = item

	err := h.service.BuyMerch(ctx, &userUuidStr, &merch)
	if err != nil {
		switch {
		case errors.Is(err, errorsx.NotEnoughMoney):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, errorsx.UnknownUser), errors.Is(err, errorsx.ItemNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, errorsx.DBError), errors.Is(err, errorsx.ServiceError):
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		default:
			h.logger.Error("unknown error service.BuyMerch",
				slog.String("method", "handler.BuyMerch"),
				slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorsx.UnknownError.Error()})
		}
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handlers) SendCoin(c *gin.Context) {
	ctx := c.Request.Context()
	userUuid, exists := c.Get("userUuid")

	if !exists {
		h.logger.Error("userUuid not found in context",
			slog.String("method", "handler.BuyMerch"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorsx.ServiceError.Error()})
		return
	}
	userUuidStr, ok := userUuid.(string)
	if !ok {
		h.logger.Error("userUuid cannot parse to string",
			slog.String("method", "handler.BuyMerch"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorsx.ServiceError.Error()})
		return
	}

	sendingCoins := entity.SendingCoins{}

	err := c.ShouldBindJSON(&sendingCoins)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorsx.InvalidInput.Error()})
		return
	}

	err = h.service.TransferCoins(ctx, &userUuidStr, &sendingCoins)
	if err != nil {
		switch {
		case errors.Is(err, errorsx.NotEnoughMoney):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, errorsx.ReceiverNotFound), errors.Is(err, errorsx.ForbiddenTransaction):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, errorsx.DBError), errors.Is(err, errorsx.ServiceError):
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		default:
			h.logger.Error("unknown error service.TransferCoins",
				slog.String("method", "handler.TransferCoins"),
				slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorsx.UnknownError.Error()})
		}
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handlers) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	userUuid, exists := c.Get("userUuid")

	if !exists {
		h.logger.Error("userUuid not found in context",
			slog.String("method", "handler.BuyMerch"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorsx.ServiceError.Error()})
		return
	}
	userUuidStr, ok := userUuid.(string)
	if !ok {
		h.logger.Error("userUuid cannot parse to string",
			slog.String("method", "handler.BuyMerch"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorsx.ServiceError.Error()})
		return
	}

	info, err := h.service.GetUserInfo(ctx, &userUuidStr)
	if err != nil {
		switch {
		case errors.Is(err, errorsx.DBError), errors.Is(err, errorsx.ServiceError):
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		default:
			h.logger.Error("unknown error service.Login",
				slog.String("method", "handler.BuyMerch"),
				slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorsx.UnknownError.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, info)
}
