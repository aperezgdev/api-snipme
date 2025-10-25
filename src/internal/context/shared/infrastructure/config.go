package infrastructure

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Server   ServerConfig
	Redis    RedisConfig
}
type AppConfig struct {
	Name    string
	Version string
	Env     string
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
	Url string
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

	cfg := &Config{
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
	}

	return cfg
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
