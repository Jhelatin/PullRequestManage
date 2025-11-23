package app

import (
	"PullRequestManage/internal/config"
	"PullRequestManage/internal/handler"
	"PullRequestManage/internal/repository/postgres"
	"PullRequestManage/internal/service"
	"PullRequestManage/pkg/database"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func Run() {
	logger := log.New(os.Stdout, "reviewer-service: ", log.LstdFlags)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("FATAL: failed to load config: %v", err)
	}

	// Подключение к БД
	pool, err := database.NewPool(cfg)
	if err != nil {
		logger.Fatalf("FATAL: failed to connect to database: %v", err)
	}
	defer pool.Close()
	logger.Println("INFO: Successfully connected to PostgreSQL")

	logger.Println("INFO: Initializing repositories...")
	userRepo := postgres.NewUserRepository(pool)
	teamRepo := postgres.NewTeamRepository(pool)
	prRepo := postgres.NewPullRequestRepository(pool)

	logger.Println("INFO: Initializing services...")
	teamService := service.NewTeamService(teamRepo)
	userService := service.NewUserService(userRepo, prRepo)
	prService := service.NewPRService(userRepo, teamRepo, prRepo)

	logger.Println("INFO: Initializing HTTP handler...")
	handler := handler.NewHandler(teamService, userService, prService, logger)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	handler.RegisterRoutes(router)

	server := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		logger.Printf("INFO: Starting server on http://localhost:%s", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("FATAL: Could not start server: %v", err)
		}
	}()

	// завершение работы от docker
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Println("INFO: Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("FATAL: Server forced to shutdown: %v", err)
	}

	logger.Println("INFO: Server exited properly")
}
