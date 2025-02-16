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
	BuyMerch(ctx context.Context, userUUID *string, item *entity.Merch) error
	TransferCoins(ctx context.Context, userUUID *string, sendingCoinsInfo *entity.SendingCoins) error
	GetUserInfo(ctx context.Context, userUUID *string) (*entity.UserInfo, error)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": errorsx.ErrInvalidInput.Error()})
		return
	}

	token, err := h.service.Login(ctx, &userDto)
	if err != nil {
		switch {
		case errors.Is(err, errorsx.ErrInvalidPassword):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case errors.Is(err, errorsx.ErrDB), errors.Is(err, errorsx.ErrService):
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		default:
			h.logger.Error("unknown error service.Login",
				slog.String("method", "handler.Login"),
				slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorsx.ErrUnknown.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, token)
}

func (h *Handlers) BuyMerch(c *gin.Context) {
	ctx := c.Request.Context()
	userUUID, exists := c.Get("userUUID")
	if !exists {
		h.logger.Error("userUUID not found in context",
			slog.String("method", "handler.BuyMerch"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorsx.ErrService.Error()})
		return
	}
	userUUIDStr, ok := userUUID.(string)
	if !ok {
		h.logger.Error("userUUID cannot parse to string",
			slog.String("method", "handler.BuyMerch"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorsx.ErrService.Error()})
		return
	}

	item := c.Param("item")
	if item == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorsx.ErrInvalidInput.Error()})
		return
	}
	var merch entity.Merch
	merch.Name = item

	err := h.service.BuyMerch(ctx, &userUUIDStr, &merch)
	if err != nil {
		switch {
		case errors.Is(err, errorsx.ErrNotEnoughMoney):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, errorsx.ErrUnknownUser), errors.Is(err, errorsx.ErrItemNotFound),
			errors.Is(err, errorsx.ErrWrongUUID):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, errorsx.ErrDB), errors.Is(err, errorsx.ErrService):
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		default:
			h.logger.Error("unknown error service.BuyMerch",
				slog.String("method", "handler.BuyMerch"),
				slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorsx.ErrUnknown.Error()})
		}
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handlers) SendCoin(c *gin.Context) {
	ctx := c.Request.Context()
	userUUID, exists := c.Get("userUUID")

	if !exists {
		h.logger.Error("userUuid not found in context",
			slog.String("method", "handler.BuyMerch"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorsx.ErrService.Error()})
		return
	}
	userUUIDdStr, ok := userUUID.(string)
	if !ok {
		h.logger.Error("userUuid cannot parse to string",
			slog.String("method", "handler.BuyMerch"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorsx.ErrService.Error()})
		return
	}

	sendingCoins := entity.SendingCoins{}

	err := c.ShouldBindJSON(&sendingCoins)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorsx.ErrInvalidInput.Error()})
		return
	}

	err = h.service.TransferCoins(ctx, &userUUIDdStr, &sendingCoins)
	if err != nil {
		switch {
		case errors.Is(err, errorsx.ErrReceiverNotFound), errors.Is(err, errorsx.ErrForbiddenTransaction),
			errors.Is(err, errorsx.ErrUnknownUser), errors.Is(err, errorsx.ErrWrongUUID):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, errorsx.ErrNotEnoughMoney):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, errorsx.ErrDB), errors.Is(err, errorsx.ErrService):
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		default:
			h.logger.Error("unknown error service.TransferCoins",
				slog.String("method", "handler.TransferCoins"),
				slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorsx.ErrUnknown.Error()})
		}
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handlers) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	userUUID, exists := c.Get("userUUID")

	if !exists {
		h.logger.Error("userUUID not found in context",
			slog.String("method", "handler.GetUserInfo"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorsx.ErrService.Error()})
		return
	}
	userUUIDStr, ok := userUUID.(string)
	if !ok {
		h.logger.Error("userUUID cannot parse to string",
			slog.String("method", "handler.GetUserInfo"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorsx.ErrService.Error()})
		return
	}

	info, err := h.service.GetUserInfo(ctx, &userUUIDStr)
	if err != nil {
		switch {
		case errors.Is(err, errorsx.ErrUnknownUser), errors.Is(err, errorsx.ErrWrongUUID):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, errorsx.ErrDB), errors.Is(err, errorsx.ErrService):
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		default:
			h.logger.Error("unknown error",
				slog.String("method", "handler.GetUserInfo"),
				slog.String("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": errorsx.ErrUnknown.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, info)
}
