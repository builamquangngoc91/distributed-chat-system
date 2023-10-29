package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tylerbui430/user-service/handlers"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("create logger error: %s", err.Error()))
	}

	dsn := "host=localhost user=postgres password=postgres dbname=socialnetwork port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Sugar().Errorf("connect database error: %s", err.Error())
		return
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// gracefull shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "anc",
		})
	})

	userHandlersDeps := &handlers.UserHandlersDeps{
		DB:          db,
		RedisClient: redisClient,
	}
	userHandlers := handlers.NewUserHandlers(userHandlersDeps)
	userHandlers.RouteGroup(router)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Sugar().Errorf("shutdown http.Server error: %s", err.Error())
			return
		}
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Sugar().Errorf("ListenAndServe error: %s", err.Error())
		return
	}
}
