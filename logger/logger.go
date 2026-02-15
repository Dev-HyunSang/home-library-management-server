package logger

import (
	"log"
	"sync"

	"go.uber.org/zap"
)

// LogConfig holds configuration for the logger
type LogConfig struct {
	Service     string
	Environment string
}

var (
	instance *zap.Logger
	sugar    *zap.SugaredLogger
	once     sync.Once
)

// InitWithConfig returns a singleton logger with service and environment fields
func InitWithConfig(cfg LogConfig) *zap.Logger {
	once.Do(func() {
		var err error
		base, err := zap.NewProduction()
		if err != nil {
			log.Fatalln(err)
		}
		instance = base.With(
			zap.String("service", cfg.Service),
			zap.String("env", cfg.Environment),
		)
		sugar = instance.Sugar()
	})
	return instance
}

// Init returns a singleton logger instance with default config
func Init() *zap.Logger {
	return InitWithConfig(LogConfig{
		Service:     "home-library",
		Environment: "dev",
	})
}

// Sugar returns the sugared logger singleton
func Sugar() *zap.SugaredLogger {
	if sugar == nil {
		Init()
	}
	return sugar
}

// Sync flushes any buffered log entries
func Sync() error {
	if instance != nil {
		return instance.Sync()
	}
	return nil
}

func ErrLog(msg string) {
	Sugar().Error(msg)
}

func UserInfoLog(userID, msg string) {
	Sugar().Infow(msg,
		"user_id", userID,
	)
}
