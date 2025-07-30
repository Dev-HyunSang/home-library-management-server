package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/dev-hyunsang/home-library/internal/config"
	"github.com/dev-hyunsang/home-library/lib/ent"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

// DB 연동
func NewDBConnection(config *config.Config) (*ent.Client, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.DB.MySQL.User,
		config.DB.MySQL.Password,
		config.DB.MySQL.Host,
		config.DB.MySQL.Port,
		config.DB.MySQL.DBName,
	)
	db, err := sql.Open(config.DB.MySQL.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// DB 설정 때문에 처음 호출 시엔 sql.open을 사용함.
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// sql.open을 ent.Client으로 변환함.
	drv := entsql.OpenDB(config.DB.MySQL.Driver, db)

	if err := NewUserTable(ent.NewClient(ent.Driver(drv))); err != nil {
		return nil, fmt.Errorf("failed to create user table: %w", err)
	}

	return ent.NewClient(ent.Driver(drv)), nil
}

// DB 생성 시 User 테이블을 생성하는 함수 / User 테이블이 존재하지 않은 경우에만 작동
func NewUserTable(client *ent.Client) error {
	if err := client.Schema.Create(context.Background()); err != nil {
		return fmt.Errorf("failed to create user table: %w", err)
	}

	return nil
}

func NewRedisConnection(config *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.DB.Redis.Host, config.DB.Redis.Port),
		Password: config.DB.Redis.Password,
		DB:       config.DB.Redis.DB,
	})

	if err := client.Ping().Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}
