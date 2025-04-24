package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel     string        `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	Address      string        `yaml:"search_address" env:"SEARCH_ADDRESS" env-default:"localhost:8080"`
	WordsAddress string        `yaml:"words_address" env:"WORDS_ADDRESS" env-default:"localhost:8081"`
	IndexTTL     time.Duration `yaml:"index_ttl" env:"INDEX_TTL" env-default:"20s"`
	DBAddress    string        `yaml:"db_address" env:"DB_ADDRESS" env-default:"localhost:5431"`
}

func MustLoad(configPath string) Config {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(err)
	}
	return cfg
}
