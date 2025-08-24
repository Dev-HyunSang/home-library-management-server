package config

import (
	"fmt"
	"os"

	"github.com/dev-hyunsang/home-library/logger"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"
)

type Config struct {
	App  AppConfig  `yaml:"app"`
	DB   DBConfig   `yaml:"db"`
	Auth AuthConfig `yaml:"auth"`
}

type AppConfig struct {
	Env   string `yaml:"env"`
	Port  int    `yaml:"port"`
	Debug bool   `yaml:"debug"`
}

type DBConfig struct {
	MySQL MySQLConfig `yaml:"mysql"`
	Redis RedisConfig `yaml:"redis"`
}

type MySQLConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type AuthConfig struct {
	CookieKey string `yaml:"cookie_key"`
}

// TODO: 추후에 YAML에서 ENV 파일로 통일 예정.
// 구조체 참고 // https://helicopter55.tistory.com/96
func LoadConfig(env string) (*Config, error) {
	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.AddConfigPath("./config")     // 실행 파일 기준 경로
	viper.AddConfigPath("../../config") // 테스트 환경에서 상대 경로로 접근 시 대비

	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func LoadEnv(key string) string {
	if err := godotenv.Load("config/.env.dev"); err != nil {
		logger.Init().Sugar().Errorf("failed to load env file: %w", err)
	}

	value := os.Getenv(key)
	if value == "" || len(value) == 0 {
		logger.Init().Sugar().Errorf("env variable %s is not set", key)
	}

	return value
}
