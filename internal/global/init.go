package global

import (
	"chain/internal/config"
	"log"
)

func Init() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	cfg.SetUpConfig("app", &AppConfig)
	AppConfig.Env = cfg.GetEnv()
	cfg.SetUpConfig("Database", &DatabaseConfig)
	cfg.SetUpConfig("BlockChain", &BlockChainConfig)
	cfg.SetUpConfig("Redis", &RedisConfig)
	cfg.SetUpConfig("Jwt", &JWT)
}
