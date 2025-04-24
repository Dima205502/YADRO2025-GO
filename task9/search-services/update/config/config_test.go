package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var configPath = "../config.yaml"

func TestMustLoad(t *testing.T) {
	defer func() {
		r := recover()
		assert.Nil(t, r)
	}()

	cfg := MustLoad(configPath)

	assert.NotEmpty(t, cfg.Address)
	assert.NotEmpty(t, cfg.DBAddress)
	assert.NotEmpty(t, cfg.LogLevel)
	assert.NotEmpty(t, cfg.WordsAddress)
	assert.NotEmpty(t, cfg.URL)

	assert.Greater(t, cfg.Concurrency, 0)
	assert.Greater(t, cfg.CheckPeriod, int64(0))
	assert.Greater(t, cfg.Timeout, int64(0))
}
