package infrastructure

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Load(t *testing.T) {
	t.Run("loads config with valid environment variables", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-jwt-secret")
		os.Setenv("OAUTH_STATE_SECRET", "test-oauth-secret")
		os.Setenv("SERVER_PORT", "9090")
		os.Setenv("JWT_EXPIRATION_MINUTES", "30")
		os.Setenv("REFRESH_TOKEN_TTL_DAYS", "15")
		os.Setenv("DATABASE_URL", "postgres://testuser:testpass@localhost:5432/testdb")
		os.Setenv("REDIS_URL", "localhost:6380")
		os.Setenv("APP_NAME", "test-app")
		os.Setenv("ENV", "test")

		defer func() {
			os.Unsetenv("JWT_SECRET")
			os.Unsetenv("OAUTH_STATE_SECRET")
			os.Unsetenv("SERVER_PORT")
			os.Unsetenv("JWT_EXPIRATION_MINUTES")
			os.Unsetenv("REFRESH_TOKEN_TTL_DAYS")
			os.Unsetenv("DATABASE_URL")
			os.Unsetenv("REDIS_URL")
			os.Unsetenv("APP_NAME")
			os.Unsetenv("ENV")
		}()

		cfg := Load()

		assert.Equal(t, "test-jwt-secret", cfg.JWT.Secret)
		assert.Equal(t, "test-oauth-secret", cfg.OAuth.StateSecret)
		assert.Equal(t, uint16(9090), cfg.Server.Port)
		assert.Equal(t, 30, cfg.JWT.ExpirationMinutes)
		assert.Equal(t, 15, cfg.JWT.RefreshTokenTTLDays)
		assert.Equal(t, "postgres://testuser:testpass@localhost:5432/testdb", cfg.Database.Url)
		assert.Equal(t, "localhost:6380", cfg.Redis.Url)
		assert.Equal(t, "test-app", cfg.App.Name)
		assert.Equal(t, "test", cfg.App.Env)
	})

	t.Run("loads config with default values", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-jwt-secret")
		os.Setenv("OAUTH_STATE_SECRET", "test-oauth-secret")
		os.Setenv("ENV", "test")

		defer func() {
			os.Unsetenv("JWT_SECRET")
			os.Unsetenv("OAUTH_STATE_SECRET")
			os.Unsetenv("ENV")
		}()

		cfg := Load()

		assert.Equal(t, uint16(8081), cfg.Server.Port)
		assert.Equal(t, 60, cfg.JWT.ExpirationMinutes)
		assert.Equal(t, 30, cfg.JWT.RefreshTokenTTLDays)
		assert.Equal(t, "snipme", cfg.App.Name)
		assert.Equal(t, "1.0.0", cfg.App.Version)
	})

	t.Run("panics when JWT_SECRET is missing", func(t *testing.T) {
		os.Setenv("OAUTH_STATE_SECRET", "test-oauth-secret")
		os.Setenv("ENV", "test")

		defer func() {
			os.Unsetenv("OAUTH_STATE_SECRET")
			os.Unsetenv("ENV")
		}()

		assert.Panics(t, func() {
			Load()
		})
	})

	t.Run("panics when OAUTH_STATE_SECRET is missing", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-jwt-secret")
		os.Setenv("ENV", "test")

		defer func() {
			os.Unsetenv("JWT_SECRET")
			os.Unsetenv("ENV")
		}()

		assert.Panics(t, func() {
			Load()
		})
	})

	t.Run("panics when SERVER_PORT is invalid", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-jwt-secret")
		os.Setenv("OAUTH_STATE_SECRET", "test-oauth-secret")
		os.Setenv("SERVER_PORT", "invalid")
		os.Setenv("ENV", "test")

		defer func() {
			os.Unsetenv("JWT_SECRET")
			os.Unsetenv("OAUTH_STATE_SECRET")
			os.Unsetenv("SERVER_PORT")
			os.Unsetenv("ENV")
		}()

		assert.Panics(t, func() {
			Load()
		})
	})

	t.Run("panics when JWT_EXPIRATION_MINUTES is invalid", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-jwt-secret")
		os.Setenv("OAUTH_STATE_SECRET", "test-oauth-secret")
		os.Setenv("JWT_EXPIRATION_MINUTES", "not-a-number")
		os.Setenv("ENV", "test")

		defer func() {
			os.Unsetenv("JWT_SECRET")
			os.Unsetenv("OAUTH_STATE_SECRET")
			os.Unsetenv("JWT_EXPIRATION_MINUTES")
			os.Unsetenv("ENV")
		}()

		assert.Panics(t, func() {
			Load()
		})
	})

	t.Run("panics when REFRESH_TOKEN_TTL_DAYS is invalid", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-jwt-secret")
		os.Setenv("OAUTH_STATE_SECRET", "test-oauth-secret")
		os.Setenv("REFRESH_TOKEN_TTL_DAYS", "not-a-number")
		os.Setenv("ENV", "test")

		defer func() {
			os.Unsetenv("JWT_SECRET")
			os.Unsetenv("OAUTH_STATE_SECRET")
			os.Unsetenv("REFRESH_TOKEN_TTL_DAYS")
			os.Unsetenv("ENV")
		}()

		assert.Panics(t, func() {
			Load()
		})
	})
}

func TestGetEnv(t *testing.T) {
	t.Run("returns environment variable value when set", func(t *testing.T) {
		t.Parallel()
		os.Setenv("TEST_VAR", "test-value")
		defer os.Unsetenv("TEST_VAR")

		value := getEnv("TEST_VAR", "default-value")
		assert.Equal(t, "test-value", value)
	})

	t.Run("returns default value when environment variable is not set", func(t *testing.T) {
		t.Parallel()
		value := getEnv("NON_EXISTENT_VAR", "default-value")
		assert.Equal(t, "default-value", value)
	})

	t.Run("returns empty string when env var is set to empty", func(t *testing.T) {
		t.Parallel()
		os.Setenv("EMPTY_VAR", "")
		defer os.Unsetenv("EMPTY_VAR")

		value := getEnv("EMPTY_VAR", "default-value")
		assert.Equal(t, "", value)
	})
}
