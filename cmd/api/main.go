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

	"github.com/gin-gonic/gin"
)

func main() {
	// load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	if cfg.Env != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	// initiate router
	router := gin.Default()

	router.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"code":   200,
			"env":    cfg.Env,
			"time":   time.Now(),
		})
	})

	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	mongodb, err := mongo.NewMongoDB(cfg)
	if err != nil {
		fmt.Printf("failed to connect to MongoDB: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("success connect to MongoDB: %s\n", cfg.Database.Name)

	defer func() {
		if err := mongodb.Close(context.Background()); err != nil {
			fmt.Printf("failed to close MongoDB conn: %v\n", err)
		}
	}()

	redisClient, err := redis.NewRedis(cfg)
	if err != nil {
		fmt.Printf("failed to connect Redis: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("success connect to Redis on %s:%s\n", cfg.Redis.Host, cfg.Redis.Port)
	defer redisClient.Close()

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

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

	fmt.Printf("server is running on %s in %s mode\n", serverAddr, cfg.Env)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Printf("error starting server: %s\n", err)
	}

	<-serverCtx.Done()
}
