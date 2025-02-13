package app

import (
	"github.com/ZhdanovichVlad/go_final_project/internal/controllers"
	"github.com/ZhdanovichVlad/go_final_project/internal/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	basePath     = "/api"
	infoPath     = "/info"
	sendCoinPath = "/sendCoin"
	buyItemPath  = "/buy/:item"
	authPath     = "/auth"
)

type router struct {
	router *gin.Engine
}

func NewRouter() *router {
	return &router{router: gin.Default()}
}

func (a *router) RegisterHandlers(h *controllers.Handlers) {

	api := a.router.Group(basePath)
	api.POST(authPath, h.Login)

	api.Use(middleware.AuthMiddleware())

	api.POST(buyItemPath, h.BuyMerch)
	api.POST(sendCoinPath, h.SendCoin)
	api.GET(infoPath, h.GetUserInfo)

	a.router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
	})
}

func (a *router) Run(host string) {
	log.Println("server started on", host)
	a.router.Run(host)
	log.Println("stopping server")
}
