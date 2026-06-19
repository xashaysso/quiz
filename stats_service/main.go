package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"stats/consumer"
	"stats/repository"
	"stats/repository/pg"
	"stats/service"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("No .env file found")
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if os.Getenv("APP_ENV") == "prod" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}
	slog.SetDefault(logger)

	slog.Info("stats service starting...")

	globalPool := repository.Serve()
	defer globalPool.Close()

	ctx, cancel := context.WithCancel(context.Background())

	statsRepo := pg.NewStatsRepo(globalPool)
	statsService := service.NewStatsService(statsRepo)

	quizConsumer := consumer.NewQuizPassedConsumer([]string{"127.0.0.1:9092"}, "quiz-results", "stats-processors", statsService)

	defer func() {
		if err := quizConsumer.Close(); err != nil {
			slog.Error("failed to close consumer", slog.Any("err", err))
		}
	}()

	go quizConsumer.Start(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	slog.Info("shutting down stats service...")

	cancel()
	timeToShutdown := 200 * time.Millisecond

	time.Sleep(timeToShutdown)
	slog.Info("stats service exited cleanly")
}
