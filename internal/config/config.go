package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	CORS     CORSConfig
	Env      string
}

type ServerConfig struct {
	Host string
	Port string
}

type DatabaseConfig struct {
	URI      string
	Name     string
	PoolSize uint64
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
}

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			URI:      getEnv("MONGO_DATABASE_URI", "mongodb://localhost:27017"),
			Name:     getEnv("MONGO_DATABASE_NAME", "sample"),
			PoolSize: uint64(getEnvAsInt("MONGO_DATABASE_POOL_SIZE", 100)),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
			PoolSize: getEnvAsInt("REDIS_POOL_SIZE", 100),
		},
		CORS: CORSConfig{
			AllowOrigins:     strings.Split(getEnv("CORS_ALLOW_ORIGINS", "*"), ","),
			AllowMethods:     strings.Split(getEnv("CORS_ALLOW_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS"), ","),
			AllowHeaders:     strings.Split(getEnv("CORS_ALLOW_HEADERS", "Content-Type,Authorization,X-CSRF_Token"), ","),
			ExposeHeaders:    strings.Split(getEnv("CORS_EXPOSE_HEADERS", ""), ","),
			AllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", true),
			MaxAge:           getEnvAsInt("CORS_MAX_AGE", 300),
		},
		Env: getEnv("ENV", "development"),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intVal
}

func getEnvAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	boolVar, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolVar
}
