package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"stats/consumer"
	"stats/handlers"
	"stats/repository"
	"stats/repository/pg"
	"stats/service"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()

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
	statsHandler := handlers.NewStatsHandler(statsService)

	stats := router.Group("/stats")
	{
		stats.GET("/users/:user_id", statsHandler.GetUserStats)
		stats.GET("/quizzes/:quiz_id", statsHandler.GetQuizGlobalStats)
	}

	PORT := os.Getenv("PORT")

	srv := &http.Server{
		Addr:    PORT,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("listen and serve failed", slog.Any("err", err))
			os.Exit(1)
		}
	}()
	slog.Info("server started succesfully", slog.String("port", PORT))

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
