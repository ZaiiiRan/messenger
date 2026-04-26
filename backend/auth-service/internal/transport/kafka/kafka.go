package kafka

import (
	"fmt"
	"strings"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/config/settings"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaClient struct {
	cfg settings.KafkaSettings
}

func New(cfg settings.KafkaSettings) (*KafkaClient, error) {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": strings.Join(cfg.Brokers, ","),
		"security.protocol": "SASL_PLAINTEXT",
		"sasl.mechanism":    "PLAIN",
		"sasl.username":     cfg.User,
		"sasl.password":     cfg.Password,
		"socket.timeout.ms": int(cfg.DialTimeout) * 1000,
	})
	if err != nil {
		return nil, fmt.Errorf("kafka: create admin client: %w", err)
	}
	defer adminClient.Close()

	if _, err = adminClient.GetMetadata(nil, false, int(cfg.DialTimeout)*1000); err != nil {
		return nil, fmt.Errorf("kafka: connect: %w", err)
	}

	return &KafkaClient{cfg: cfg}, nil
}

func (k *KafkaClient) ConfigMap() kafka.ConfigMap {
	return kafka.ConfigMap{
		"bootstrap.servers": strings.Join(k.cfg.Brokers, ","),
		"security.protocol": "SASL_PLAINTEXT",
		"sasl.mechanism":    "PLAIN",
		"sasl.username":     k.cfg.User,
		"sasl.password":     k.cfg.Password,
		"socket.timeout.ms": int(k.cfg.DialTimeout) * 1000,
	}
}

func (k *KafkaClient) Close() {}
