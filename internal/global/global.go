package global

import (
	"chain/internal/config"
)

var (
	AppConfig        config.AppConfig
	DatabaseConfig   config.DatabaseConfig
	BlockChainConfig config.BlockChainConfig
	RedisConfig      config.RedisConfig
	JWT              config.JwtConfig
)
