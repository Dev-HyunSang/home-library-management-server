package logger

import (
	"log"

	"go.uber.org/zap"
)

func InitLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln(err)
	}

	return logger
}

func ErrLog(msg string) {
	logger := InitLogger()

	defer logger.Sync()

	logger.Sugar().Errorf(msg)
}

func UserInfoLog(userID, msg string) {
	logger := InitLogger()

	defer logger.Sync()

	logger.Sugar().Infow(msg,
		"user_id", userID,
	)
}
