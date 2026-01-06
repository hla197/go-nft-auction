package db

import (
	"chain/internal/config"
	"chain/internal/logger"
	"chain/internal/models"
	"fmt"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(config config.DatabaseConfig) {
	var err error
	// 从环境变量获取MySQL连接配置
	dbHost := config.Host
	dbPort := config.Port
	dbUser := config.Username
	dbPassword := config.Password
	dbName := config.Name

	// 构建MySQL连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, strconv.Itoa(dbPort), dbName)

	// 连接MySQL数据库
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.Log.Errorf("Failed to connect to MySQL database:", err)
	}
}

func Migrate() {
	DB.AutoMigrate(&models.Events{})
	DB.AutoMigrate(&models.Auction{})
	DB.AutoMigrate(&models.AuctionLog{})
	DB.AutoMigrate(&models.Nft{})
}

// GetDB 获取数据库连接实例
func GetDB() *gorm.DB {
	return DB
}
