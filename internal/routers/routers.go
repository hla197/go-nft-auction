package routers

import (
	"chain/internal/handlers"
	"chain/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitApi(router *gin.Engine) {
	// 静态资源（JS / CSS / 图片）
	router.Static("/assets", "./web/assets")

	// SPA 入口（React 路由兜底）
	router.NoRoute(func(c *gin.Context) {
		c.File("./web/index.html")
	})
	// 跨域
	router.Use(middleware.CORSMiddleware())

	authHandler := &handlers.AuthHandler{}
	auctionHandler := &handlers.AuctionHandler{}

	// 公共接口（不需要 token）
	public := router.Group("/auth")
	{
		public.POST("/login", authHandler.Login)
		public.POST("/register", authHandler.Register)
	}

	auth := router.Group("")
	auth.Use(middleware.JWTAuthMiddleware())
	{
	}

	auction := router.Group("auction")
	{
		auction.POST("/getNftHistory", auctionHandler.GetNftHistory)
	}
}
