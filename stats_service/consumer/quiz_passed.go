package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"stats/entities"
	"stats/service"
	"time"

	"github.com/segmentio/kafka-go"
)

type QuizPassedConsumer struct {
	reader       *kafka.Reader
	statsService service.StatsServiceInterface
}

func NewQuizPassedConsumer(brokers []string, topic, groupID string, statService service.StatsServiceInterface) *QuizPassedConsumer {
	return &QuizPassedConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 1,
			MaxBytes: 10e6,
			MaxWait:  1 * time.Second,
		}),
		statsService: statService,
	}
}

func (c *QuizPassedConsumer) Start(ctx context.Context) {
	slog.Info("quiz_passed consumer started, listening for messages...")
	for {
		select {
		case <-ctx.Done():
			slog.Info("stopping consumer loop due to context cancellation")
			return
		default:
			readCtx, readCancel := context.WithTimeout(ctx, 5*time.Second)
			msg, err := c.reader.ReadMessage(readCtx)
			readCancel()
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					continue
				}
				if errors.Is(err, context.Canceled) {
					return
				}
				slog.Error("failed to read message from kafka", slog.Any("err", err))
				continue
			}

			processCtx, processCancel := context.WithTimeout(context.Background(), 3*time.Second)
			var event entities.QuizPassedEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				slog.Error("unable to unmarshal kafka message JSON", slog.Any("err", err), slog.Int64("offset", msg.Offset))
				processCancel()
				continue
			}

			err = c.statsService.ProcessQuizPassed(processCtx, event)
			processCancel()
			if err != nil {
				slog.Error("failed to process event in service", slog.Any("err", err))
				continue
			}
		}
	}
}

func (c *QuizPassedConsumer) Close() error {
	return c.reader.Close()
}
