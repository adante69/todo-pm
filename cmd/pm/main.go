package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"todo-project/internal/app"
	"todo-project/internal/config"
	"todo-project/internal/redis"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting server")

	rdb := redis.NewRedisClient(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
	err := rdb.Set(context.Background(), "key", "value", 0).Err()
	if err != nil {
		log.Info("Failed to set key in Redis", slog.String("error", err.Error()))
	}
	val, err := rdb.Get(context.Background(), "key").Result()
	if err != nil {
		log.Info("Failed to get key from Redis", slog.String("error", err.Error()))
	}
	fmt.Println("key:", val)

	// Initialize the application
	application := app.New(log, cfg.Server.GRPC.Port, cfg.Database.DSN)
	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	sign := <-stop
	application.GRPCServer.Stop()
	log.Info("shutting down server", slog.String("reason", sign.String()))

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
