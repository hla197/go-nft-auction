package redis

import (
	"context"
	"fmt"
	"time"

	"chain/internal/config"

	"github.com/go-redis/redis/v8"
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

func Init(cfg config.RedisConfig) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), // 从配置文件读取
		Password: cfg.Password,                             // 密码
		DB:       cfg.DB,

		// --- 生产级核心配置 ---
		PoolSize:     cfg.PoolSize, // 连接池大小，根据并发量调整
		MinIdleConns: 10,           // 最小空闲连接，减少建立连接的开销
		MaxRetries:   3,            // 网络错误重试次数

		// 超时控制 (防止 goroutine 泄露)
		DialTimeout:  5 * time.Second, // 连接超时
		ReadTimeout:  3 * time.Second, // 读超时
		WriteTimeout: 3 * time.Second, // 写超时
		PoolTimeout:  4 * time.Second, // 从连接池获取连接的超时时间
	})

	// 测试连接
	_, err := Client.Ping(Ctx).Result()
	return err
}

// Close 优雅关闭 (在 main.go 的 defer 中调用)
func Close() {
	_ = Client.Close()
}
