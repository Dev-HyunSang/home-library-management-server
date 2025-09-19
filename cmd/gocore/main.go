package main

import (
	"flag"
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
	"github.com/dev-hyunsang/home-library/logger"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
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
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: encryptcookie.GenerateKey(),
	}))

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger.Init(),
	}))

	app.Use(cors.New(cors.Config{
		// TODO: production 에서 수정
		AllowOrigins: "http://localhost:3000, http://localhost:5173/, http://192.168.0.13:5173/",
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

	dbConn, err := db.NewDBConnection(cfg)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	logger.Init().Sugar().Info("성공적으로 데이터베이스(MySQL)에 연결되었습니다.")

	store := NewSessionStore(cfg)

	logger.Init().Sugar().Info("성공적으로 Redis에 세션 저장소가 초기화되었습니다.")

	// csrfConfig := csrf.Config{
	// 	Session:        store,
	// 	KeyLookup:      "json:csrf",
	// 	CookieName:     "__Host-csrf",
	// 	CookieSameSite: "Lax",
	// 	CookieSecure:   true,
	// 	CookieHTTPOnly: true,
	// 	ContextKey:     "csrf",
	// 	ErrorHandler: func(c *fiber.Ctx, err error) error {
	// 		return c.Status(fiber.StatusForbidden).JSON(handler.ErrorHandler(domain.ErrInvalidCSRFToken))
	// 	},
	// 	Expiration: time.Minute * 30,
	// }

	// _ := csrf.New(csrfConfig)

	// 사용자 관련 의존성 주입
	authRepo := redisRepository.NewAuthRepository(store)
	userRepo := repository.NewUserRepository(dbConn, store)
	userUseCase := usecase.NewUserUseCase(userRepo, authRepo)
	userHandler := handler.NewUserHandler(userUseCase, userUseCase)

	// 책 관련 의존성 주입
	bookRepo := repository.NewBookRepository(dbConn)
	bookUseCase := usecase.NewBookUseCase(bookRepo)
	bookHandler := handler.NewBookHandler(bookUseCase, authRepo)

	api := app.Group("/api")
	user := api.Group("/users")
	user.Post("/register", userHandler.UserRegisterHandler)
	user.Post("/login", userHandler.UserLoginHandler)
	user.Post("/me", userHandler.UserVerifyHandler)
	user.Get("/:id", userHandler.UserGetByIdHandler)
	user.Put("/:id", userHandler.UserEditHandler)
	user.Delete("/:id", userHandler.UserDeleteHandler)

	books := api.Group("/books")
	books.Post("/", bookHandler.SaveBookHandler)
	books.Get("/", bookHandler.GetBooksHandler)
	books.Delete("/:id", bookHandler.BookDeleteHandler)
	books.Get("/:name", bookHandler.GetBooksByUserNameHandler)
	books.Post("/search", bookHandler.SearchBookIsbnHandler)

	if err := app.Listen(":3000"); err != nil {
		logger.Init().Sugar().Fatalf("서버를 시작하는 도중 오류가 발생했습니다: %v", err)
	}
}
