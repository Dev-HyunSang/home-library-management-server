package logger

import (
	"log"

	"go.uber.org/zap"
)

func Init() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln(err)
	}

	return logger
}

func ErrLog(msg string) {
	logger := Init()

	defer logger.Sync()

	logger.Sugar().Errorf(msg)
}

func UserInfoLog(userID, msg string) {
	logger := Init()

	defer logger.Sync()

	logger.Sugar().Infow(msg,
		"user_id", userID,
	)
}
