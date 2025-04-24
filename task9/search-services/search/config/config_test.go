package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var configPath = "../config.yaml"

func TestMustLoad(t *testing.T) {
	defer func() {
		r := recover()
		assert.Nil(t, r)
	}()

	cfg := MustLoad(configPath)

	require.NotEmpty(t, cfg.LogLevel)
	require.NotEmpty(t, cfg.Address)
	require.NotEmpty(t, cfg.WordsAddress)
	require.NotEmpty(t, cfg.DBAddress)

	require.Greater(t, cfg.IndexTTL, int64(0))
}
