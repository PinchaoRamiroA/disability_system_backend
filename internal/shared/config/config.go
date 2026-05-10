package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	DB       DBConfig
	JWT      JWTConfig
	Server   ServerConfig
}

type AppConfig struct {
	Name    string
	Env     string
	Version string
}

type DBConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
	URL             string
}

type JWTConfig struct {
	Secret          string
	Expiration      time.Duration
	RefreshExpiry   time.Duration
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using system environment variables")
	}

	return &Config{
		App:    loadAppConfig(),
		DB:     loadDBConfig(),
		JWT:    loadJWTConfig(),
		Server: loadServerConfig(),
	}
}

func loadAppConfig() AppConfig {
	return AppConfig{
		Name:    getEnv("APP_NAME", "disability_system"),
		Env:     getEnv("APP_ENV", "development"),
		Version: getEnv("APP_VERSION", "1.0.0"),
	}
}

func loadDBConfig() DBConfig {
	port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	maxOpen, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	maxIdle, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "10"))
	maxLifetime, _ := strconv.Atoi(getEnv("DB_CONN_MAX_LIFETIME", "5"))

	return DBConfig{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            port,
		User:            getEnv("DB_USER", "postgres"),
		Password:        getEnv("DB_PASSWORD", "postgres"),
		Name:            getEnv("DB_NAME", "app"),
		SSLMode:         getEnv("DB_SSLMODE", "disable"),
		MaxOpenConns:    maxOpen,
		MaxIdleConns:    maxIdle,
		ConnMaxLifetime: maxLifetime,
		URL:             getEnv("DATABASE_URL", ""),
	}
}

func loadJWTConfig() JWTConfig {
	expiration, _ := time.ParseDuration(getEnv("JWT_EXPIRATION", "24h"))
	refreshExpiry, _ := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h"))

	return JWTConfig{
		Secret:        getEnv("JWT_SECRET", "supersecret"),
		Expiration:    expiration,
		RefreshExpiry: refreshExpiry,
	}
}

func loadServerConfig() ServerConfig {
	readTimeout, _ := time.ParseDuration(getEnv("SERVER_READ_TIMEOUT", "30s"))
	writeTimeout, _ := time.ParseDuration(getEnv("SERVER_WRITE_TIMEOUT", "30s"))

	return ServerConfig{
		Port:         getEnv("APP_PORT", "8080"),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func (c *Config) DSN() string {
	if c.DB.URL != "" {
		return c.DB.URL
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DB.Host,
		c.DB.Port,
		c.DB.User,
		c.DB.Password,
		c.DB.Name,
		c.DB.SSLMode,
	)
}