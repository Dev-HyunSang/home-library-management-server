package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"time"

	"github.com/dev-hyunsang/home-library/internal/cache"
	"github.com/dev-hyunsang/home-library/internal/config"
	"github.com/dev-hyunsang/home-library/internal/db"
	"github.com/dev-hyunsang/home-library/internal/handler"
	"github.com/dev-hyunsang/home-library/internal/infrastructure/fcm"
	"github.com/dev-hyunsang/home-library/internal/infrastructure/kafka"
	"github.com/dev-hyunsang/home-library/internal/middleware"
	repository "github.com/dev-hyunsang/home-library/internal/repository/mysql"
	redisRepository "github.com/dev-hyunsang/home-library/internal/repository/redis"
	"github.com/dev-hyunsang/home-library/internal/usecase"
	"github.com/dev-hyunsang/home-library/logger"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
)

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
		AllowOrigins: "http://localhost:3000, http://localhost:5173/, http://192.168.219.100:5173",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodDelete,
			fiber.MethodPut,
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

	// Redis 클라이언트 초기화
	redisClient := cache.NewRedisClient(
		cfg.DB.Redis.Host,
		cfg.DB.Redis.Port,
		cfg.DB.Redis.Password,
		cfg.DB.Redis.DB,
	)

	logger.Init().Sugar().Info("성공적으로 Redis 클라이언트가 초기화되었습니다.")

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

	// Kafka & FCM 초기화
	// 1. FCM Service
	fcmService, err := fcm.NewFCMService(cfg.FCM.ServiceAccountPath)
	if err != nil {
		logger.Init().Sugar().Warnf("FCM 서비스를 초기화하는데 실패했습니다(알림 기능 비활성화): %v", err)
	} else {
		logger.Init().Sugar().Info("FCM 서비스가 성공적으로 초기화되었습니다.")
	}

	// 2. Kafka Producer
	kafkaProducer := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	defer kafkaProducer.Close()
	logger.Init().Sugar().Info("Kafka Producer가 초기화되었습니다.")

	// 3. Kafka Consumer (백그라운드 실행)
	kafkaConsumer := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.GroupID, fcmService)
	go func() {
		logger.Init().Sugar().Info("Kafka Consumer가 시작되었습니다.")
		kafkaConsumer.Start(context.Background())
	}()
	defer kafkaConsumer.Close()

	// 사용자 관련 의존성 주입
	authRepo := redisRepository.NewAuthRepository(cfg.JWT.Secret, 1*time.Hour, 24*time.Hour, redisClient)
	userRepo := repository.NewUserRepository(dbConn, nil)
	authUseCase := usecase.NewAuthUseCase(authRepo)
	userUseCase := usecase.NewUserUseCase(userRepo, authUseCase)
	userHandler := handler.NewUserHandler(userUseCase, authUseCase)
	authHandler := handler.NewAuthHandler(authUseCase)

	// 책 관련 의존성 주입
	bookRepo := repository.NewBookRepository(dbConn)
	bookUseCase := usecase.NewBookUseCase(bookRepo, kafkaProducer)
	bookHandler := handler.NewBookHandler(bookUseCase, authUseCase)

	api := app.Group("/api")
	user := api.Group("/users")
	user.Post("/signup", userHandler.UserSignUpHandler)
	user.Get("/check/nickname", userHandler.UserVerifyByNicknameHandler)
	user.Post("/signin", userHandler.UserSignInHandler)
	user.Post("/signout", middleware.JWTAuthMiddleware(authUseCase), userHandler.UserSignOutHandler)
	user.Post("/forgot-password", userHandler.UserRestPasswordHandler)
	user.Post("/me", middleware.JWTAuthMiddleware(authUseCase), userHandler.UserVerifyHandler)
	user.Get("/:id", middleware.JWTAuthMiddleware(authUseCase), userHandler.UserGetByIdHandler)
	user.Put("/update/:id", middleware.JWTAuthMiddleware(authUseCase), userHandler.UserEditHandler)
	user.Delete("/:id", middleware.JWTAuthMiddleware(authUseCase), userHandler.UserDeleteHandler)

	books := api.Group("/books")
	books.Post("/add", middleware.JWTAuthMiddleware(authUseCase), bookHandler.SaveBookHandler)
	books.Get("/get", middleware.JWTAuthMiddleware(authUseCase), bookHandler.GetBooksHandler)
	books.Delete("/delete/:id", middleware.JWTAuthMiddleware(authUseCase), bookHandler.BookDeleteHandler)
	books.Get("/:name", middleware.JWTAuthMiddleware(authUseCase), bookHandler.GetBooksByUserNameHandler)
	books.Post("/search", middleware.JWTAuthMiddleware(authUseCase), bookHandler.SearchBookIsbnHandler)

	reviews := books.Group("/reviews")
	reviews.Post("/", middleware.JWTAuthMiddleware(authUseCase), bookHandler.SaveBookReviewHandler)
	reviews.Get("/get", middleware.JWTAuthMiddleware(authUseCase), bookHandler.GetBookReviewByUserIDHandler)

	bookmarks := books.Group("/bookmarks")
	bookmarks.Post("/add/:id", middleware.JWTAuthMiddleware(authUseCase), bookHandler.AddBookmarkHandler)
	bookmarks.Get("/get", middleware.JWTAuthMiddleware(authUseCase), bookHandler.GetBookmarksByUserIDHandler)
	bookmarks.Delete("/delete/:id", middleware.JWTAuthMiddleware(authUseCase), bookHandler.DeleteBookmarkHandler)

	auth := api.Group("/auth")
	auth.Post("/refresh", authHandler.RefreshTokenHandler)
	auth.Post("/revoke-all", middleware.JWTAuthMiddleware(authUseCase), authHandler.RevokeAllTokensHandler)
	auth.Get("/rate-limit", middleware.JWTAuthMiddleware(authUseCase), authHandler.CheckRateLimitHandler)

	if err := app.Listen(":3000"); err != nil {
		logger.Init().Sugar().Fatalf("서버를 시작하는 도중 오류가 발생했습니다: %v", err)
	}
}
