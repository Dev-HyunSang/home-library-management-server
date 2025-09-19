package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/dev-hyunsang/home-library/logger"
	"github.com/joho/godotenv"
)

type Config struct {
	App  AppConfig  `json:"app"`
	DB   DBConfig   `json:"db"`
	Auth AuthConfig `json:"auth"`
}

type AppConfig struct {
	Env   string `json:"env"`
	Port  int    `json:"port"`
	Debug bool   `json:"debug"`
}

type DBConfig struct {
	MySQL MySQLConfig `json:"mysql"`
	Redis RedisConfig `json:"redis"`
}

type MySQLConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type AuthConfig struct {
	CookieKey string `json:"cookie_key"`
}

func LoadConfig(env string) (*Config, error) {
	// ENV 파일 로드
	envFile := fmt.Sprintf(".env.%s", env)
	if err := godotenv.Load(envFile); err != nil {
		logger.Init().Sugar().Warnf("Warning: failed to load env file %s: %v", envFile, err)
	}

	// APP 설정
	appPort, err := strconv.Atoi(getEnvOrDefault("APP_PORT", "3000"))
	if err != nil {
		return nil, fmt.Errorf("invalid APP_PORT: %w", err)
	}

	appDebug, err := strconv.ParseBool(getEnvOrDefault("APP_DEBUG", "false"))
	if err != nil {
		return nil, fmt.Errorf("invalid APP_DEBUG: %w", err)
	}

	// MySQL 설정
	mysqlPort, err := strconv.Atoi(getEnvOrDefault("MYSQL_PORT", "3306"))
	if err != nil {
		return nil, fmt.Errorf("invalid MYSQL_PORT: %w", err)
	}

	// Redis 설정
	redisPort, err := strconv.Atoi(getEnvOrDefault("REDIS_PORT", "6379"))
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_PORT: %w", err)
	}

	redisDB, err := strconv.Atoi(getEnvOrDefault("REDIS_DB", "0"))
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB: %w", err)
	}

	config := &Config{
		App: AppConfig{
			Env:   getEnvOrDefault("APP_ENV", "development"),
			Port:  appPort,
			Debug: appDebug,
		},
		DB: DBConfig{
			MySQL: MySQLConfig{
				Driver:   getEnvOrDefault("MYSQL_DRIVER", "mysql"),
				Host:     getEnvOrDefault("MYSQL_HOST", "localhost"),
				Port:     mysqlPort,
				User:     getEnvOrDefault("MYSQL_USER", "root"),
				Password: getEnvOrDefault("MYSQL_PASSWORD", ""),
				DBName:   getEnvOrDefault("MYSQL_DBNAME", ""),
				SSLMode:  getEnvOrDefault("MYSQL_SSLMODE", "disable"),
			},
			Redis: RedisConfig{
				Host:     getEnvOrDefault("REDIS_HOST", "localhost"),
				Port:     redisPort,
				Username: getEnvOrDefault("REDIS_USERNAME", ""),
				Password: getEnvOrDefault("REDIS_PASSWORD", ""),
				DB:       redisDB,
			},
		},
		Auth: AuthConfig{
			CookieKey: getEnvOrDefault("AUTH_COOKIE_KEY", ""),
		},
	}

	// 필수 값 검증
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func validateConfig(config *Config) error {
	if config.DB.MySQL.DBName == "" {
		return fmt.Errorf("MYSQL_DBNAME is required")
	}
	if config.Auth.CookieKey == "" {
		return fmt.Errorf("AUTH_COOKIE_KEY is required")
	}
	return nil
}
