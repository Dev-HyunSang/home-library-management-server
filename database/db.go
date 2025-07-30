package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB는 데이터베이스 연결을 위한 전역 변수입니다.
var DB *gorm.DB

// InitDatabase는 데이터베이스 연결을 초기화합니다.
func InitDatabase() {
	// 로깅 설정
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // 느린 쿼리 기준 시간
			LogLevel:                  logger.Info, // 로그 레벨
			IgnoreRecordNotFoundError: true,        // RecordNotFound 오류 무시
			Colorful:                  true,        // 컬러 출력 활성화
		},
	)

	// SQLite 데이터베이스 연결
	database, err := gorm.Open(sqlite.Open("home-library.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("데이터베이스 연결에 실패했습니다: %v", err)
	}

	// SQLite에서 외래 키 제약 조건 활성화
	_ = database.Exec("PRAGMA foreign_keys = ON")

	// 전역 변수에 데이터베이스 인스턴스 할당
	DB = database

	// 데이터베이스 모델 마이그레이션
	MigrateModels(DB)

	log.Println("데이터베이스 연결 및 마이그레이션에 성공했습니다")
}

// GetDB는 데이터베이스 인스턴스를 반환합니다.
func GetDB() *gorm.DB {
	return DB
}
