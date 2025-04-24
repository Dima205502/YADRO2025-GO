package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var configPath = "../config.yaml"

func TestMustLoad(t *testing.T) {
	defer func() {
		r := recover()
		require.Nil(t, r)
	}()

	cfg := MustLoad(configPath)

	/*
		LogLevel          string        `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
		SearchConcurrency int           `yaml:"search_concurrency" env:"SEARCH_CONCURRENCY" env-default:"1"`
		SearchRate        int           `yaml:"search_rate" env:"SEARCH_RATE" env-default:"1"`
		WordsAddress      string        `yaml:"words_address" env:"WORDS_ADDRESS" env-default:"words:8081"`
		UpdateAddress     string        `yaml:"update_address" env:"UPDATE_ADDRESS" env-default:"update:8082"`
		SearchAddress     string        `yaml:"search_address" env:"SEARCH_ADDRESS" env-default:"search:8083"`
		TokenTTL          time.Duration `yaml:"token_ttl" env:"TOKEN_TTL" env-default:"24h"`

		Address string        `yaml:"address" env:"API_ADDRESS" env-default:"localh
		Timeout time.Duration `yaml:"timeout" env:"API_TIMEOUT" env-default:"5s"`
	*/

	require.NotEmpty(t, cfg.LogLevel)
	require.NotEmpty(t, cfg.WordsAddress)
	require.NotEmpty(t, cfg.UpdateAddress)
	require.NotEmpty(t, cfg.SearchAddress)
	require.NotEmpty(t, cfg.Address)

	require.Greater(t, cfg.TokenTTL, int64(0))
	require.Greater(t, cfg.Timeout, int64(0))

	require.Greater(t, cfg.SearchConcurrency, 0)
	require.Greater(t, cfg.SearchRate, 0)
}
