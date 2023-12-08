package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"auth-service/configs"
	"auth-service/handlers"
	"auth-service/repositories"

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

	configs.LoadConfig()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		configs.Cfg.Database.Host,
		configs.Cfg.Database.Username,
		configs.Cfg.Database.Password,
		configs.Cfg.Database.Name,
		configs.Cfg.Database.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Sugar().Errorf("connect database error: %s", err.Error())
		return
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			configs.Cfg.Redis.Host,
			configs.Cfg.Redis.Port,
		),
	})

	// gracefull shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	router := gin.Default()

	// TODO: add ping and health
	userHandlersDeps := &handlers.UserHandlersDeps{
		DB:          db,
		RedisClient: redisClient,
		UserRepo:    repositories.NewUserRepository(),
	}
	userHandlers := handlers.NewUserHandlers(userHandlersDeps)
	userHandlers.RouteGroup(router)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", configs.Cfg.AuthService.Port),
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
