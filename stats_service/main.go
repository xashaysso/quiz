package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("stats service starting...")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "quiz-results",
		GroupID:  "stats-processors",
		MinBytes: 1,
		MaxBytes: 10e6,
		MaxWait:  1 * time.Second,
	})

	defer func() {
		if err := reader.Close(); err != nil {
			slog.Error("failed to close kafka reader", slog.Any("err", err))
		} else {
			slog.Info("kafka reader closed cleanly")
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		slog.Info("kafka consumer started, listening for messages...")
		for {
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				slog.Error("failed to read message from kafka", slog.Any("err", err))
				continue
			}

			slog.Info("RECEIVED MESSAGE FROM KAFKA", slog.String("key", string(msg.Key)), slog.String("value", string(msg.Value)), slog.Int64("offset", msg.Offset))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	slog.Info("shutting down stats service...")

	cancel()

	timeToShutdown := 200 * time.Millisecond

	time.Sleep(timeToShutdown)
	slog.Info("stats service exited cleanly")
}
