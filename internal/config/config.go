package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App AppConfig `yaml:"app"`
	DB  DBConfig  `yaml:"db"`
	RedisConfig
	Jwt JwtConfig `yaml:"jwt"`
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

type JwtConfig struct {
	AccessToken  string `yaml:"access_token"`
	RefreshToken string `yaml:"refresh_token"`
}

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
