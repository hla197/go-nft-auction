package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	vp *viper.Viper
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	fmt.Printf("🚀 启动应用，当前环境: %s\n", env)

	vp := viper.New()
	vp.AutomaticEnv()
	vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	vp.SetConfigName(fmt.Sprintf("config.%s", env))
	vp.SetConfigType("yaml")
	vp.AddConfigPath(".")

	if err := vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取通用配置失败: %w", err)
	}
	if err := vp.MergeInConfig(); err != nil {
		// 如果没有特定环境配置文件，忽略错误
		fmt.Printf("合并配置文件失败 %s\n", err)
	}

	vp.SetDefault("env", env)

	return &Config{vp}, nil
}

func (c *Config) GetEnv() string {
	return c.vp.GetString("env")
}

func (c *Config) SetUpConfig(k string, config interface{}) error {
	err := c.vp.UnmarshalKey(k, config)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) Get(k string) interface{} {
	return c.vp.Get(k)
}
