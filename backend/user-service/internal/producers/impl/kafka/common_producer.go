package implkafkaproducer

import (
	"context"
	"fmt"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/config/settings"
	kafkatransport "github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Message struct {
	Key   string
	Value string
}

type Producer struct {
	kProducer    *kafka.Producer
	topic        string
	writeTimeout int
	name         string
}

func New(cfg settings.KafkaProducerSettings, kafkaClient *kafkatransport.KafkaClient) (*Producer, error) {
	cm := kafkaClient.ConfigMap()
	cm["client.id"] = cfg.ClientID
	cm["message.timeout.ms"] = int(cfg.WriteTimeout) * 1000

	kProducer, err := kafka.NewProducer(&cm)
	if err != nil {
		return nil, fmt.Errorf("kafka producer: create: %w", err)
	}

	return &Producer{
		kProducer:    kProducer,
		topic:        cfg.Topic,
		writeTimeout: int(cfg.WriteTimeout) * 1000,
		name:         cfg.Name,
	}, nil
}

func (p *Producer) Produce(ctx context.Context, msg Message) error {
	deliveryChan := make(chan kafka.Event, 1)

	if err := p.kProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(msg.Key),
		Value: []byte(msg.Value),
	}, deliveryChan); err != nil {
		return fmt.Errorf("kafka producer %s: enqueue: %w", p.name, err)
	}

	select {
	case e := <-deliveryChan:
		m, ok := e.(*kafka.Message)
		if !ok {
			return fmt.Errorf("kafka producer %s: unexpected event type", p.name)
		}
		if m.TopicPartition.Error != nil {
			return fmt.Errorf("kafka producer %s: delivery failed: %w", p.name, m.TopicPartition.Error)
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *Producer) Close() {
	p.kProducer.Flush(p.writeTimeout)
	p.kProducer.Close()
}
