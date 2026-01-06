package main

import (
	"chain/internal/global"
	"chain/internal/infra/db"
	"chain/internal/infra/eth"
	"chain/internal/infra/redis"
	"chain/internal/logger"
	"chain/internal/middleware"
	"chain/internal/routers"
	"chain/internal/scanner"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化日志
	logger.InitLogger()
	// 初始化配置文件
	global.Init()
	// 初始化数据库、Redis、区块链客户端等
	initMain()
	// 初始化 Gin
	initGin()
}

func initGin() {
	write := logger.Log.GetIoWriter()
	// 将 Gin 的日志输出指向 Zap
	// 重定向必须在 gin.New() 前
	gin.DefaultWriter = write
	gin.DefaultErrorWriter = write

	logger.Log.Infoln("starting handlers")

	router := gin.New()
	router.Use(middleware.GinLogMiddleware())
	router.Use(middleware.GinRecoveryWithLogger())

	routers.InitApi(router)

	logger.Log.Infoln("starting routers", global.AppConfig.Port)
	// 启动服务
	port := ":" + global.AppConfig.Port
	router.Run(port)
}

func initMain() {
	db.InitDB(global.DatabaseConfig)
	redis.Init(global.RedisConfig)
	// 数据迁移
	if global.AppConfig.Env == "dev" {
		db.Migrate()
	}
	eth.NewEthClientManager(global.BlockChainConfig.RpcUrl, global.BlockChainConfig.WssUrl)

	scanner.Run()
}
