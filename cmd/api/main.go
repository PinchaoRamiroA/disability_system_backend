package main

import (
	"context"
	"fmt"
	
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"disability_system_backend/internal/shared/config"
	"disability_system_backend/internal/shared/database"
	"disability_system_backend/internal/shared/logger"
	"disability_system_backend/internal/shared/middleware"
	"disability_system_backend/internal/shared/response"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.App.Env)

	log.Info("Starting server", 
		"name", cfg.App.Name,
		"env", cfg.App.Env,
		"version", cfg.App.Version,
	)

	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	if err := database.HealthCheck(db); err != nil {
		log.Error("database health check failed", "error", err)
		os.Exit(1)
	}

	log.Info("database connected successfully")

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := setupRouter(cfg, log)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Info("server listening", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", "error", err)
	}

	log.Info("server exited")
}

func setupRouter(cfg *config.Config, log *logger.Logger) *gin.Engine {
	router := gin.New()

	router.Use(middleware.Recovery(nil))
	router.Use(middleware.RequestID())
	router.Use(middleware.CORS([]string{"*"}))
	router.Use(middleware.Logger(func(msg string, args ...any) {
		log.Info(msg, args...)
	}))

	router.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status":  "ok",
			"version": cfg.App.Version,
		}, "healthy")
	})

	router.GET("/", func(c *gin.Context) {
		response.Success(c, gin.H{
			"name":    cfg.App.Name,
			"env":     cfg.App.Env,
			"version": cfg.App.Version,
		}, "welcome")
	})

	return router
}