package kafka

import (
	"context"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	err := createTopicWithRetry(brokers[0], topic)
	if err != nil {
		slog.Error("failed to create kafka topic, attempting to continue...", slog.Any("err", err))
	}
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			WriteTimeout: 10 * time.Second,
			Async:        false,
		},
	}
}

func (p *Producer) SendMessage(ctx context.Context, key, value []byte) error {
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
	if err != nil {
		slog.Error("failed to send message to kafka", slog.Any("err", err))
		return err
	}

	slog.Info("message sent to kafka successfully", slog.String("key", string(key)))
	return nil
}

func (p *Producer) Close() error {
	if err := p.writer.Close(); err != nil {
		slog.Error("failed to close kafka writer", slog.Any("err", err))
		return err
	}
	slog.Info("kafka writer closed clearly")
	return nil
}

func createTopicWithRetry(brokerAddr, topic string) error {
	conn, err := kafka.Dial("tcp", brokerAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	err = conn.CreateTopics(topicConfig)
	if err != nil {
		return err
	}

	slog.Info("kafka topic verified or created successfully", slog.String("topic", topic))
	return nil
}
