package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nielwyn/inventory-system/config"
	"github.com/nielwyn/inventory-system/internal/database"
	"github.com/nielwyn/inventory-system/internal/handlers"
	"github.com/nielwyn/inventory-system/internal/middleware"
	"github.com/nielwyn/inventory-system/internal/repository"
	"github.com/nielwyn/inventory-system/internal/service"
	"github.com/nielwyn/inventory-system/pkg/logger"
	"github.com/nielwyn/inventory-system/pkg/validator"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	if err := logger.Init(cfg.Log.Level, cfg.Log.Encoding); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Licenser API")

	gin.SetMode(cfg.Server.Mode)

	db, err := database.New(cfg.Database.GetDSN())
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	if err := db.AutoMigrate(); err != nil {
		logger.Fatal("Failed to run database migrations", zap.Error(err))
	}

	validator.RegisterCustomValidations()

	userRepo := repository.NewUserRepository(db.DB)
	subscriptionRepo := repository.NewSubscriptionRepository(db.DB)

	authService := service.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.ExpiryHours)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)

	healthHandler := handlers.NewHealthHandler(db)
	authHandler := handlers.NewAuthHandler(authService)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)

	// 100 requests/second per IP, burst of 20
	rateLimiter := middleware.NewRateLimiter(rate.Limit(100), 20)

	router := setupRouter(healthHandler, authHandler, subscriptionHandler, authService, rateLimiter)

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	go func() {
		logger.Info("Server starting", zap.String("address", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server stopped")
}

func setupRouter(
	healthHandler *handlers.HealthHandler,
	authHandler *handlers.AuthHandler,
	subscriptionHandler *handlers.SubscriptionHandler,
	authService service.AuthService,
	rateLimiter *middleware.RateLimiter,
) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(rateLimiter.Middleware())

	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		subscriptions := v1.Group("/subscriptions")
		subscriptions.Use(middleware.Auth(authService))
		{
			subscriptions.POST("", subscriptionHandler.CreateSubscription)
			subscriptions.GET("", subscriptionHandler.GetAllSubscriptions)
			subscriptions.GET("/:id", subscriptionHandler.GetSubscriptionByID)
			subscriptions.PUT("/:id", subscriptionHandler.UpdateSubscription)
			subscriptions.DELETE("/:id", subscriptionHandler.DeleteSubscription)
		}
	}

	return router
}
