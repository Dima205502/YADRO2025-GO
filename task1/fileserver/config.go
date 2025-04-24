package main

import (
	"flag"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port        int    `yaml:"port" env:"FILESERVER_PORT"`
	StoragePath string `yaml:"storage_path" env-default:"./FileStorage"`
}

func InitConfig() (*Config, error) {
	configPath := flag.String("config", "", "Path to config.yaml file")
	flag.Parse()

	var cfg Config
	var err error

	if *configPath != "" {
		err = cleanenv.ReadConfig(*configPath, &cfg)
	} else {
		err = cleanenv.ReadEnv(&cfg)
	}

	return &cfg, err
}
