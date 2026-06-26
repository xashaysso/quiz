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
	"stats/pkg/kafka"
	"stats/repository"
	"stats/repository/pg"
	"stats/service"
	"strings"
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
		users := stats.Group("/users")
		{
			users.GET("/:user_id", statsHandler.GetUserStats)
			users.GET("/:user_id/analytics", statsHandler.GetUserAnalytics)
		}
		quizzes := stats.Group("/quizzes")
		{
			quizzes.GET("/:quiz_id", statsHandler.GetQuizStats)
			quizzes.GET("/:quiz_id/analytics", statsHandler.GetQuizAnalytics)
		}
		stats.GET("/leaderboard", statsHandler.GetUserLeaderboard)
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

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "127.0.0.1:9092"
	}
	brokersSlice := strings.Split(kafkaBrokers, ",")

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		kafkaTopic = "quiz-results"
	}

	kafkaGroupID := os.Getenv("KAFKA_GROUP_ID")
	if kafkaGroupID == "" {
		kafkaGroupID = "stats-processors"
	}

	quizConsumer := consumer.NewQuizPassedConsumer(brokersSlice, kafkaTopic, kafkaGroupID, statsService)

	defer func() {
		if err := quizConsumer.Close(); err != nil {
			slog.Error("failed to close consumer", slog.Any("err", err))
		}
	}()

	if !kafka.WaitForKafka(brokersSlice, 30*time.Second) {
		slog.Error("kafka is not available after 30 seconds, exiting")
		os.Exit(1)
	}

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
