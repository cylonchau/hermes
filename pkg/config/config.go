package config

import (
	"fmt"
	"strings"

	"github.com/cylonchau/hermes/pkg/store"
	"github.com/spf13/viper"
)

// Config 全局配置结构
type Config struct {
	AppName        string               `mapstructure:"app_name"`
	DatabaseDriver string               `mapstructure:"database_driver"`
	MySQL          store.DatabaseConfig `mapstructure:"mysql"`
	SQLite         store.DatabaseConfig `mapstructure:"sqlite"`
	Database       store.DatabaseConfig `mapstructure:"database"` // 保留旧的兼容性
	Server         ServerConfig         `mapstructure:"server"`
}

// ServerConfig 服务器基础配置
type ServerConfig struct {
	LogLevel string `mapstructure:"log_level"`
	Port     int    `mapstructure:"port"`
}

var (
	globalConfig *Config
	CONFIG       *Config
)

// InitConfiguration 初始化配置 (兼容新架构)
func InitConfiguration(cfgPath string) error {
	_, err := Load(cfgPath)
	return err
}

// Load 从指定路径加载配置
func Load(cfgPath string) (*Config, error) {
	v := viper.New()
	if cfgPath != "" {
		v.SetConfigFile(cfgPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
	}

	v.SetEnvPrefix("HERMES")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if cfg.AppName == "" {
		cfg.AppName = "hermes"
	}

	globalConfig = cfg
	CONFIG = cfg
	return cfg, nil
}

// Get 获取全局配置单例
func Get() *Config {
	return globalConfig
}
