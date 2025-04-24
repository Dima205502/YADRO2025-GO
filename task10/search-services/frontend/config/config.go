package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPConfig struct {
	Address string        `yaml:"address" env:"FRONTEND_ADDRESS" env-default:"localhost:8080"`
	Timeout time.Duration `yaml:"timeout" env:"FRONTEND_TIMEOUT" env-default:"5s"`
}

type Config struct {
	LogLevel    string `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	Api_address string `yaml:"api_address" env:"API_ADDRESS" env-default:"localhost:8080"`
	HTTPConfig  `yaml:"frontend_server"`
}

func MustLoad(configPath string) Config {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(err)
	}
	return cfg
}
