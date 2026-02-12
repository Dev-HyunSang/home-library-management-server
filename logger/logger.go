package logger

import (
	"log"
	"sync"

	"go.uber.org/zap"
)

var (
	instance *zap.Logger
	sugar    *zap.SugaredLogger
	once     sync.Once
)

// Init returns a singleton logger instance
func Init() *zap.Logger {
	once.Do(func() {
		var err error
		instance, err = zap.NewProduction()
		if err != nil {
			log.Fatalln(err)
		}
		sugar = instance.Sugar()
	})
	return instance
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
