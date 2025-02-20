package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server `yaml:"server"`
	DB     `yaml:"db"`
	Logger `yaml:"logger"`
	Secret string `yaml:"secret"`
}

type Server struct {
	Address         string        `yaml:"address" env-default:":8080"`
	Timeout         time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env-default:"60s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"10s"`
}

type DB struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	Username string `yaml:"username" env-default:"postgres"`
	Name     string `yaml:"name" env-default:"house_service"`
	Password string `yaml:"password" env-default:"postgres"`
}

type Logger struct {
	Level string `yaml:"level"`
}

func ParseConfig(path string) (*Config, error) {
	var cfg *Config

	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config error: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("ummarshal to config struct is failed: %w", err)
	}

	return cfg, nil
}
