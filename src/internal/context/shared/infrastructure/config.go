package infrastructure

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App         AppConfig
	Database    DatabaseConfig
	GEOFilePath string
	Server      ServerConfig
	Redis       RedisConfig
	Loki        LokiConfig
	OAuth       OAuthConfig
	JWT         JWTConfig
}

type AppConfig struct {
	Name    string
	Version string
	Env     string
}

type LokiConfig struct {
	Url string
}

type ServerConfig struct {
	Port         uint16
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

type DatabaseConfig struct {
	Url string
}

type RedisConfig struct {
	Url      string
	Password string
}

type OAuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	GitHubClientID     string
	GitHubClientSecret string
	GitHubRedirectURL  string
	StateSecret        string
}

type JWTConfig struct {
	Secret              string
	ExpirationMinutes   int
	RefreshTokenTTLDays int
}

func Load() *Config {
	err := godotenv.Load()
	_, envExist := os.LookupEnv("ENV")
	if err != nil && !envExist {
		panic("Error loading .env file")
	}

	serverPort, err := strconv.ParseUint(getEnv("SERVER_PORT", "8081"), 10, 16)
	if err != nil {
		panic("Invalid SERVER_PORT value in .env file")
	}

	jwtExpirationMinutes, err := strconv.Atoi(getEnv("JWT_EXPIRATION_MINUTES", "60"))
	if err != nil {
		panic("Invalid JWT_EXPIRATION_MINUTES value in .env file")
	}

	refreshTokenTTLDays, err := strconv.Atoi(getEnv("REFRESH_TOKEN_TTL_DAYS", "30"))
	if err != nil {
		panic("Invalid REFRESH_TOKEN_TTL_DAYS value in .env file")
	}

	cfg := &Config{
		GEOFilePath: getEnv("GEO_FILE_PATH", "../db/geo/GeoLite2-Country.mmdb"),
		App: AppConfig{
			Name:    getEnv("APP_NAME", "snipme"),
			Version: getEnv("APP_VERSION", "1.0.0"),
			Env:     getEnv("ENV", "dev"),
		},
		Server: ServerConfig{
			Port: uint16(serverPort),
		},
		Database: DatabaseConfig{
			Url: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/mydb"),
		},
		Redis: RedisConfig{
			Url:      getEnv("REDIS_URL", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		Loki: LokiConfig{
			Url: getEnv("LOKI_URL", ""),
		},
		OAuth: OAuthConfig{
			GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8081/auth/google/callback"),
			GitHubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
			GitHubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
			GitHubRedirectURL:  getEnv("GITHUB_REDIRECT_URL", "http://localhost:8081/auth/github/callback"),
			StateSecret:        getEnv("OAUTH_STATE_SECRET", ""),
		},
		JWT: JWTConfig{
			Secret:              getEnv("JWT_SECRET", ""),
			ExpirationMinutes:   jwtExpirationMinutes,
			RefreshTokenTTLDays: refreshTokenTTLDays,
		},
	}

	if cfg.JWT.Secret == "" {
		panic("JWT_SECRET is required in .env file")
	}

	if cfg.OAuth.StateSecret == "" {
		panic("OAUTH_STATE_SECRET is required in .env file")
	}

	return cfg
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
