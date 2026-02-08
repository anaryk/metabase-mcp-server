package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Flags(t *testing.T) {
	cfg, err := Load([]string{
		"--metabase-url", "http://localhost:3000",
		"--api-key", "mb_key123",
		"--log-level", "debug",
	})
	require.NoError(t, err)
	assert.Equal(t, "http://localhost:3000", cfg.MetabaseURL)
	assert.Equal(t, "mb_key123", cfg.APIKey)
	assert.Equal(t, "debug", cfg.LogLevel)
}

func TestLoad_EnvVars(t *testing.T) {
	t.Setenv("METABASE_URL", "http://metabase:3000")
	t.Setenv("METABASE_API_KEY", "mb_envkey")
	t.Setenv("LOG_LEVEL", "warn")

	cfg, err := Load(nil)
	require.NoError(t, err)
	assert.Equal(t, "http://metabase:3000", cfg.MetabaseURL)
	assert.Equal(t, "mb_envkey", cfg.APIKey)
	assert.Equal(t, "warn", cfg.LogLevel)
}

func TestLoad_FlagsPrecedence(t *testing.T) {
	t.Setenv("METABASE_URL", "http://env:3000")
	t.Setenv("METABASE_API_KEY", "env_key")

	cfg, err := Load([]string{
		"--metabase-url", "http://flag:3000",
		"--api-key", "flag_key",
	})
	require.NoError(t, err)
	assert.Equal(t, "http://flag:3000", cfg.MetabaseURL)
	assert.Equal(t, "flag_key", cfg.APIKey)
}

func TestLoad_SessionAuth(t *testing.T) {
	cfg, err := Load([]string{
		"--metabase-url", "http://localhost:3000",
		"--username", "admin@test.com",
		"--password", "secret",
	})
	require.NoError(t, err)
	assert.Equal(t, "admin@test.com", cfg.Username)
	assert.Equal(t, "secret", cfg.Password)
}

func TestLoad_TrailingSlashRemoved(t *testing.T) {
	cfg, err := Load([]string{
		"--metabase-url", "http://localhost:3000/",
		"--api-key", "key",
	})
	require.NoError(t, err)
	assert.Equal(t, "http://localhost:3000", cfg.MetabaseURL)
}

func TestLoad_MissingURL(t *testing.T) {
	_, err := Load([]string{"--api-key", "key"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "metabase URL is required")
}

func TestLoad_MissingAuth(t *testing.T) {
	_, err := Load([]string{"--metabase-url", "http://localhost:3000"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "either API key")
}

func TestLoad_PartialSessionAuth(t *testing.T) {
	_, err := Load([]string{
		"--metabase-url", "http://localhost:3000",
		"--username", "admin@test.com",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "either API key")
}

func TestLoad_DefaultLogLevel(t *testing.T) {
	cfg, err := Load([]string{
		"--metabase-url", "http://localhost:3000",
		"--api-key", "key",
	})
	require.NoError(t, err)
	assert.Equal(t, "info", cfg.LogLevel)
}

func TestLoad_DefaultTransportAndPort(t *testing.T) {
	cfg, err := Load([]string{
		"--metabase-url", "http://localhost:3000",
		"--api-key", "key",
	})
	require.NoError(t, err)
	assert.Equal(t, "stdio", cfg.Transport)
	assert.Equal(t, 8808, cfg.Port)
}

func TestLoad_SSETransport(t *testing.T) {
	cfg, err := Load([]string{
		"--metabase-url", "http://localhost:3000",
		"--api-key", "key",
		"--transport", "sse",
		"--port", "9090",
	})
	require.NoError(t, err)
	assert.Equal(t, "sse", cfg.Transport)
	assert.Equal(t, 9090, cfg.Port)
}

func TestLoad_TransportEnvVars(t *testing.T) {
	t.Setenv("METABASE_URL", "http://metabase:3000")
	t.Setenv("METABASE_API_KEY", "key")
	t.Setenv("TRANSPORT", "sse")
	t.Setenv("PORT", "7777")

	cfg, err := Load(nil)
	require.NoError(t, err)
	assert.Equal(t, "sse", cfg.Transport)
	assert.Equal(t, 7777, cfg.Port)
}

func TestLoad_InvalidTransport(t *testing.T) {
	_, err := Load([]string{
		"--metabase-url", "http://localhost:3000",
		"--api-key", "key",
		"--transport", "grpc",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "transport must be")
}
