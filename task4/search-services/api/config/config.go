package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel     string `yaml:"log_level" env:"LOG_LEVEL"`
	WordsAddress string `yaml:"words_address" env:"WORDS_ADDRESS"`
	ServerConfig `yaml:"http_server" `
}

type ServerConfig struct {
	Address string        `yaml:"address" env:"HTTP_SERVER_ADDRESS"`
	Timeout time.Duration `yaml:"timeout" env:"HTTP_SERVER_TIMEOUT"`
}

func LoadConfig(configPath string) Config {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config %q: %s", configPath, err)
	}
	return cfg
}
