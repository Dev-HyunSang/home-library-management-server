package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/dev-hyunsang/home-library/internal/config"
	"github.com/dev-hyunsang/home-library/internal/db"
	"github.com/dev-hyunsang/home-library/internal/handler"
	repository "github.com/dev-hyunsang/home-library/internal/repository/mysql"
	redisRepository "github.com/dev-hyunsang/home-library/internal/repository/redis"
	"github.com/dev-hyunsang/home-library/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis/v3"
)

// USING REDIS IN SESSION STORAGE
func NewSessionStore(cfg *config.Config) *session.Store {
	storage := redis.New(redis.Config{
		Host:      cfg.DB.Redis.Host,
		Port:      cfg.DB.Redis.Port,
		Password:  cfg.DB.Redis.Password,
		Database:  cfg.DB.Redis.DB,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	})

	store := session.New(session.Config{
		Storage:           storage,
		Expiration:        24 * time.Hour,
		KeyLookup:         "cookie:session",
		CookieSessionOnly: true,
		CookieSecure:      true,
		CookieHTTPOnly:    true,
		CookieSameSite:    "None",
	})

	return store
}

func main() {
	app := fiber.New()
	// 안전한 쿠키 사용을 위해 쿠키를 암호화 함.
	// Key는 32자의 문자열이며, 무작위 값으로 생성됨.
	// app.Use(encryptcookie.New(encryptcookie.Config{
	// 	Key: encryptcookie.GenerateKey(),
	// }))

	app.Use(cors.New(cors.Config{
		// TODO: production 에서 수정
		AllowOrigins: "https://localhost:3000",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowCredentials: true,
	}))

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
	userHandler := handler.NewUserHandler(userUseCase, userUseCase)

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
