package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dev-hyunsang/home-library/logger"
	"github.com/joho/godotenv"
)

type Config struct {
	App   AppConfig   `json:"app"`
	DB    DBConfig    `json:"db"`
	Auth  AuthConfig  `json:"auth"`
	JWT   JWTConfig   `json:"jwt"`
	Kafka KafkaConfig `json:"kafka"`
	FCM   FCMConfig   `json:"fcm"`
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

type JWTConfig struct {
	Secret string `json:"secret"`
	TTL    string `json:"ttl"`
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
		JWT: JWTConfig{
			Secret: getEnvOrDefault("JWT_SECRET", ""),
			TTL:    getEnvOrDefault("JWT_TTL", "24h"),
		},
		Kafka: KafkaConfig{
			Brokers: strings.Split(getEnvOrDefault("KAFKA_BROKERS", "localhost:29092"), ","),
			Topic:   getEnvOrDefault("KAFKA_TOPIC", "notifications"),
			GroupID: getEnvOrDefault("KAFKA_GROUP_ID", "home-library-notification-group"),
		},
		FCM: FCMConfig{
			ServiceAccountPath: getEnvOrDefault("FCM_SERVICE_ACCOUNT_PATH", "serviceAccountKey.json"),
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
	if config.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}

type KafkaConfig struct {
	Brokers []string `json:"brokers"`
	Topic   string   `json:"topic"`
	GroupID string   `json:"group_id"`
}

type FCMConfig struct {
	ServiceAccountPath string `json:"service_account_path"`
}

func GetEnv(key string) string {
	err := godotenv.Load(".env.dev")
	if err != nil {
		logger.Init().Sugar().Warnf("Warning: failed to load .env file: %v", err)
	}

	return os.Getenv(key)
}
