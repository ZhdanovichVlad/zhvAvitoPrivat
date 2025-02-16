package app

import (
	"net/http"

	"github.com/ZhdanovichVlad/go_final_project/internal/controllers"
	"github.com/ZhdanovichVlad/go_final_project/internal/middleware"
	"github.com/gin-gonic/gin"
)

const (
	basePath     = "/api"
	infoPath     = "/info"
	sendCoinPath = "/sendCoin"
	buyItemPath  = "/buy/:item"
	authPath     = "/auth"
)

type Router struct {
	router *gin.Engine
}

func NewRouter() *Router {
	return &Router{router: gin.Default()}
}

func (a *Router) RegisterHandlers(h *controllers.Handlers) {

	api := a.router.Group(basePath)
	api.POST(authPath, h.Login)

	api.Use(middleware.AuthMiddleware())

	api.GET(buyItemPath, h.BuyMerch)
	api.POST(sendCoinPath, h.SendCoin)
	api.GET(infoPath, h.GetUserInfo)

	a.router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
	})
}

func (a *Router) Run(host string) error {
	err := a.router.Run(host)
	if err != nil {
		return err
	}
	return err
}
