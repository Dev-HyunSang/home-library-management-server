package main

import (
	"flag"
	"log"
	"time"

	"github.com/dev-hyunsang/my-own-library-backend/internal/cache"
	"github.com/dev-hyunsang/my-own-library-backend/internal/config"
	"github.com/dev-hyunsang/my-own-library-backend/internal/db"
	"github.com/dev-hyunsang/my-own-library-backend/internal/handler"
	"github.com/dev-hyunsang/my-own-library-backend/internal/infrastructure/fcm"
	"github.com/dev-hyunsang/my-own-library-backend/internal/infrastructure/scheduler"
	"github.com/dev-hyunsang/my-own-library-backend/internal/middleware"
	repository "github.com/dev-hyunsang/my-own-library-backend/internal/repository/mysql"
	redisRepository "github.com/dev-hyunsang/my-own-library-backend/internal/repository/redis"
	"github.com/dev-hyunsang/my-own-library-backend/internal/usecase"
	"github.com/dev-hyunsang/my-own-library-backend/logger"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
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

	logger.Sugar().Info("성공적으로 데이터베이스(MySQL)에 연결되었습니다.")

	// Redis 클라이언트 초기화
	redisClient := cache.NewRedisClient(
		cfg.DB.Redis.Host,
		cfg.DB.Redis.Port,
		cfg.DB.Redis.Password,
		cfg.DB.Redis.DB,
	)

	logger.Sugar().Info("성공적으로 Redis 클라이언트가 초기화되었습니다.")

	// FCM 초기화
	fcmService, err := fcm.NewFCMService(cfg.FCM.ServiceAccountPath)
	if err != nil {
		logger.Sugar().Warnf("FCM 서비스를 초기화하는데 실패했습니다(알림 기능 비활성화): %v", err)
	} else {
		logger.Sugar().Info("FCM 서비스가 성공적으로 초기화되었습니다.")
	}

	// 사용자 관련 의존성 주입
	authRepo := redisRepository.NewAuthRepository(cfg.JWT.Secret, 1*time.Hour, 24*time.Hour, redisClient, cfg.JWT.Issuer, cfg.JWT.Audience)
	userRepo := repository.NewUserRepository(dbConn, nil)
	emailVerificationRepo := redisRepository.NewEmailVerificationRepository(redisClient)
	authUseCase := usecase.NewAuthUseCase(authRepo)
	userUseCase := usecase.NewUserUseCase(userRepo, authUseCase)
	userHandler := handler.NewUserHandler(userUseCase, authUseCase, emailVerificationRepo)
	authHandler := handler.NewAuthHandler(authUseCase)

	// 책 관련 의존성 주입
	bookRepo := repository.NewBookRepository(dbConn)
	bookUseCase := usecase.NewBookUseCase(bookRepo)
	bookHandler := handler.NewBookHandler(bookUseCase, authUseCase)

	// 리뷰 관련 의존성 주입
	reviewRepo := repository.NewReviewRepository(dbConn)
	reviewUseCase := usecase.NewReviewUseCase(reviewRepo)
	reviewHandler := handler.NewReviewHandler(reviewUseCase, authUseCase, bookUseCase)

	// 읽기 리마인더 관련 의존성 주입
	reminderRepo := repository.NewReadingReminderRepository(dbConn)
	reminderUseCase := usecase.NewReadingReminderUseCase(reminderRepo)
	reminderHandler := handler.NewReadingReminderHandler(reminderUseCase, authUseCase)

	// 관리자 API Key 관련 의존성 주입
	apiKeyRepo := repository.NewAdminAPIKeyRepository(dbConn)
	apiKeyUseCase := usecase.NewAdminAPIKeyUseCase(apiKeyRepo)

	// 관리자 핸들러 초기화
	adminHandler := handler.NewAdminHandler(userRepo, apiKeyUseCase)

	// 리마인더 스케줄러 시작
	reminderScheduler, err := scheduler.NewReminderScheduler(reminderRepo, userRepo, fcmService)
	if err != nil {
		logger.Sugar().Warnf("리마인더 스케줄러 초기화 실패: %v", err)
	} else {
		if err := reminderScheduler.Start(); err != nil {
			logger.Sugar().Warnf("리마인더 스케줄러 시작 실패: %v", err)
		} else {
			defer reminderScheduler.Stop()
		}
	}

	api := app.Group("/api")
	user := api.Group("/users")
	user.Post("/signup", userHandler.UserSignUpHandler)
	user.Get("/check/nickname", userHandler.UserVerifyByNicknameHandler)
	user.Get("/verify/email/:email", userHandler.UserVerifyByEmailHandler)
	user.Post("/verify/code", userHandler.UserVerifyCodeHandler)
	user.Post("/signin", userHandler.UserSignInHandler)
	user.Post("/signout", middleware.JWTAuthMiddleware(authUseCase), userHandler.UserSignOutHandler)
	user.Post("/forgot-password", userHandler.UserRestPasswordHandler)
	user.Put("/change-password", middleware.JWTAuthMiddleware(authUseCase), userHandler.UserChangePasswordHandler)
	user.Post("/me", middleware.JWTAuthMiddleware(authUseCase), userHandler.UserVerifyHandler)
	user.Get("/:id", middleware.JWTAuthMiddleware(authUseCase), userHandler.UserGetByIdHandler)
	user.Put("/update/:id", middleware.JWTAuthMiddleware(authUseCase), userHandler.UserEditHandler)
	user.Delete("/:id", middleware.JWTAuthMiddleware(authUseCase), userHandler.UserDeleteHandler)
	user.Put("/fcm-token", middleware.JWTAuthMiddleware(authUseCase), userHandler.UpdateFCMTokenHandler)
	user.Put("/timezone", middleware.JWTAuthMiddleware(authUseCase), userHandler.UpdateTimezoneHandler)

	books := api.Group("/books")
	books.Post("/add", middleware.JWTAuthMiddleware(authUseCase), bookHandler.SaveBookHandler)
	books.Get("/get", middleware.JWTAuthMiddleware(authUseCase), bookHandler.GetBooksHandler)
	books.Get("/get/:user_id/:book_id", middleware.JWTAuthMiddleware(authUseCase), bookHandler.GetBookHandler)
	books.Put("/update/:id", middleware.JWTAuthMiddleware(authUseCase), bookHandler.UpdateBookHandler)
	books.Delete("/delete/:id", middleware.JWTAuthMiddleware(authUseCase), bookHandler.BookDeleteHandler)
	books.Get("/:name", middleware.JWTAuthMiddleware(authUseCase), bookHandler.GetBooksByUserNameHandler)
	books.Post("/search", middleware.JWTAuthMiddleware(authUseCase), bookHandler.SearchBookIsbnHandler)

	// ISBN 기반 리뷰 API
	reviewsAPI := api.Group("/reviews")
	reviewsAPI.Get("/me", middleware.JWTAuthMiddleware(authUseCase), reviewHandler.GetMyReviewsHandler)
	reviewsAPI.Post("/:isbn", middleware.JWTAuthMiddleware(authUseCase), reviewHandler.CreateReviewHandler)
	reviewsAPI.Get("/:isbn", reviewHandler.GetReviewsByISBNHandler)
	reviewsAPI.Get("/:isbn/:id", reviewHandler.GetReviewByIDHandler)
	reviewsAPI.Put("/:isbn/:id", middleware.JWTAuthMiddleware(authUseCase), reviewHandler.UpdateReviewHandler)
	reviewsAPI.Delete("/:isbn/:id", middleware.JWTAuthMiddleware(authUseCase), reviewHandler.DeleteReviewHandler)

	bookmarks := books.Group("/bookmarks")
	bookmarks.Post("/add/:id", middleware.JWTAuthMiddleware(authUseCase), bookHandler.AddBookmarkHandler)
	bookmarks.Get("/get", middleware.JWTAuthMiddleware(authUseCase), bookHandler.GetBookmarksByUserIDHandler)
	bookmarks.Delete("/delete/:id", middleware.JWTAuthMiddleware(authUseCase), bookHandler.DeleteBookmarkHandler)

	reminders := api.Group("/reminders")
	reminders.Post("/", middleware.JWTAuthMiddleware(authUseCase), reminderHandler.CreateReminderHandler)
	reminders.Get("/", middleware.JWTAuthMiddleware(authUseCase), reminderHandler.GetRemindersHandler)
	reminders.Put("/:id", middleware.JWTAuthMiddleware(authUseCase), reminderHandler.UpdateReminderHandler)
	reminders.Patch("/:id/toggle", middleware.JWTAuthMiddleware(authUseCase), reminderHandler.ToggleReminderHandler)
	reminders.Delete("/:id", middleware.JWTAuthMiddleware(authUseCase), reminderHandler.DeleteReminderHandler)

	auth := api.Group("/auth")
	auth.Post("/refresh", authHandler.RefreshTokenHandler)
	auth.Post("/revoke-all", middleware.JWTAuthMiddleware(authUseCase), authHandler.RevokeAllTokensHandler)
	auth.Get("/rate-limit", middleware.JWTAuthMiddleware(authUseCase), authHandler.CheckRateLimitHandler)

	// 관리자 API Key 부트스트랩 (최초 API Key 생성용)
	adminBootstrap := api.Group("/admin/bootstrap")
	adminBootstrap.Use(middleware.AdminBootstrapMiddleware(cfg.Admin.BootstrapKey))
	adminBootstrap.Post("/api-keys", adminHandler.CreateAPIKeyHandler)

	// 관리자 API (API Key 인증)
	admin := api.Group("/admin")
	admin.Use(middleware.AdminAPIKeyMiddleware(apiKeyUseCase))
	admin.Post("/notifications/broadcast", adminHandler.BroadcastNotificationHandler)
	admin.Get("/api-keys", adminHandler.GetAPIKeysHandler)
	admin.Post("/api-keys", adminHandler.CreateAPIKeyHandler)
	admin.Patch("/api-keys/:id/deactivate", adminHandler.DeactivateAPIKeyHandler)
	admin.Delete("/api-keys/:id", adminHandler.DeleteAPIKeyHandler)

	if err := app.Listen(":3000"); err != nil {
		logger.Sugar().Fatalf("서버를 시작하는 도중 오류가 발생했습니다: %v", err)
	}
}
