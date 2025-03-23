package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"codebase-go/internal/config"
	"codebase-go/internal/database/mongo"
	"codebase-go/internal/database/redis"
	"codebase-go/internal/middleware/cors"
	"codebase-go/pkg/logger"
	"codebase-go/pkg/wrapper"

	"github.com/gin-gonic/gin"
)

func monitorConnection(ctx context.Context, mongodb *mongo.MongoDB, redisClient *redis.Redis) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// check mongodb conn
			if err := mongodb.CheckConnection(ctx); err != nil {
				logger.Error("mongodb connection lost", err)
			}

			// check redis conn
			if err := redisClient.CheckConnection(ctx); err != nil {
				logger.Error("redis connection lost", err)
			}
		}
	}
}

func main() {
	// load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger.Init(cfg.Env)

	if cfg.Env != "development" && cfg.Env != "local" {
		gin.SetMode(gin.ReleaseMode)
	}

	// initiate router
	router := gin.Default()
	router.Use(cors.CORS(cfg))

	router.GET("/health-check", func(c *gin.Context) {
		wrapper.Success(c, gin.H{
			"env":  cfg.Env,
			"time": time.Now(),
		})
	})

	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// connect to database (mongoDB)
	mongodb, err := mongo.NewMongoDB(cfg)
	if err != nil {
		logger.Error("failed to connect MongoDB", err)
		os.Exit(1)
	}

	logger.Info("success connect to MongoDB", "database", cfg.Database.Name)

	defer func() {
		if err := mongodb.Close(context.Background()); err != nil {
			fmt.Printf("failed to close MongoDB conn: %v\n", err)
		}
	}()

	// connect to redis
	redisClient, err := redis.NewRedis(cfg)
	if err != nil {
		logger.Error("failed to connect redis", err)
		os.Exit(1)
	}

	logger.Info("success connect to redis", "host", cfg.Redis.Host, "port", cfg.Redis.Port)
	defer redisClient.Close()

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	go monitorConnection(serverCtx, mongodb, redisClient)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutDownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutDownCtx.Done()
			if shutDownCtx.Err() == context.DeadlineExceeded {
				fmt.Println("graceful shutdown timed out..., forcing exit")
			}
		}()

		err := server.Shutdown(shutDownCtx)
		if err != nil {
			fmt.Printf("error shutting down server: %s\n", err)
		}
		serverStopCtx()
	}()

	logger.Info("server is starting on", "address", serverAddr, "environment", cfg.Env)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Error("failed to start server", err)
	}

	<-serverCtx.Done()
}
