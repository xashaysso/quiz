package kafka_producers

import (
	"context"
	"encoding/json"
	"fmt"
	"quiz/entities/dto"
	pkgKafka "quiz/pkg/kafka"
	"strconv"
	"time"
)

type QuizKafkaProducer struct {
	Producer *pkgKafka.Producer
}

func NewQuizKafkaProducer(producer *pkgKafka.Producer) *QuizKafkaProducer {
	return &QuizKafkaProducer{
		Producer: producer,
	}
}

func (p *QuizKafkaProducer) PublishQuizPassed(ctx context.Context, userID int64, quizID int64, score int) error {
	event := &dto.QuizPassedEvent{
		UserID:   userID,
		QuizID:   quizID,
		Score:    score,
		PassedAt: time.Now(),
	}

	valueByte, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal quiz passed event: %w", err)
	}

	keyStr := strconv.FormatInt(userID, 10)
	keyBytes := []byte(keyStr)

	err = p.Producer.SendMessage(ctx, keyBytes, valueByte)
	if err != nil {
		return fmt.Errorf("failed to send quiz event via raw producer: %w", err)
	}

	return nil
}
