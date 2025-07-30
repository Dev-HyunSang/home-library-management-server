package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"github.com/dev-hyunsang/home-library/internal/config"
	"github.com/dev-hyunsang/home-library/internal/db"
	"github.com/dev-hyunsang/home-library/internal/handler"
	repository "github.com/dev-hyunsang/home-library/internal/repository/mysql"
	redisRepository "github.com/dev-hyunsang/home-library/internal/repository/redis"
	"github.com/dev-hyunsang/home-library/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis/v3"
)

// USING REDIS IN SESSION STORAGE
func NewSessionStore(cfg *config.Config) *session.Store {
	storage := redis.New(redis.Config{
		Host:      cfg.DB.Redis.Host,
		Port:      cfg.DB.Redis.Port,
		Username:  cfg.DB.Redis.Username,
		Password:  cfg.DB.Redis.Password,
		Database:  cfg.DB.Redis.DB,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	})

	store := session.New(session.Config{
		Storage:        storage,
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: "None",
	})

	return store
}

func main() {
	app := fiber.New()

	env := flag.String("env", "dev", "Environment (dev, qa, stg, prod)")
	flag.Parse()

	validEnvs := map[string]bool{"dev": true, "qa": true, "stg": true, "prod": true}
	if !validEnvs[*env] {
		log.Fatalf("Invalid environment: %s. Valid environments are: dev, qa, stg, prod", *env)
	}

	cfg, err := config.LoadConfig(*env)
	if err != nil {
		log.Fatalf("Config load error: %v", err)
	}

	fmt.Printf("config: %+v\n", cfg)

	dbConn, err := db.NewDBConnection(cfg)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	store := NewSessionStore(cfg)
	// 사용자 관련 의존성 주입
	authRepo := redisRepository.NewAuthRepository(store)
	userRepo := repository.NewUserRepository(dbConn, store)
	userUseCase := usecase.NewUserUseCase(userRepo, authRepo)
	userHandler := handler.NewUserHandler(userUseCase)

	user := app.Group("/user")
	user.Post("/register", userHandler.Register)
	user.Post("/login", userHandler.Login)
	user.Get("/:id", userHandler.GetByID)
	user.Put("/:id", userHandler.Edit)
	user.Delete("/:id", userHandler.Delete)

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}
