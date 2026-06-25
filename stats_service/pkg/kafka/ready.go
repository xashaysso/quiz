package kafka

import (
	"log/slog"
	"net"
	"time"
)

func WaitForKafka(brokers []string, timeout time.Duration) bool {
	slog.Info("waiting for kafka brokers to become available...", slog.Any("brokers", brokers))
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		for _, broker := range brokers {
			conn, err := net.DialTimeout("tcp", broker, 1*time.Second)
			if err == nil {
				conn.Close()
				slog.Info("successfully connected to kafka broker", slog.String("broker", broker))
				return true
			}
		}
		time.Sleep(2 * time.Second)
	}
	return false
}
