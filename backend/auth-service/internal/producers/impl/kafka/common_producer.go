package implkafkaproducer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/config/settings"
	kafkatransport "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

type Message struct {
	Key   string
	Value string
}

type Producer struct {
	kProducer     *kafka.Producer
	topic         string
	msgCh         chan Message
	done          chan struct{}
	closeOnce     sync.Once
	batchSize     int
	flushInterval time.Duration
	writeTimeout  int
	log           *zap.SugaredLogger
}

func New(cfg settings.KafkaProducerSettings, kafkaClient *kafkatransport.KafkaClient, log *zap.SugaredLogger) (*Producer, error) {
	cm := kafkaClient.ConfigMap()
	cm["client.id"] = cfg.ClientID
	cm["message.timeout.ms"] = int(cfg.WriteTimeout) * 1000

	kProducer, err := kafka.NewProducer(&cm)
	if err != nil {
		return nil, fmt.Errorf("kafka producer: create: %w", err)
	}

	p := &Producer{
		kProducer:     kProducer,
		topic:         cfg.Topic,
		msgCh:         make(chan Message, int(cfg.BatchSize)*2),
		done:          make(chan struct{}),
		batchSize:     int(cfg.BatchSize),
		flushInterval: time.Duration(cfg.FlushFrequency) * time.Millisecond,
		writeTimeout:  int(cfg.WriteTimeout) * 1000,
		log:           log,
	}

	go p.run()
	go p.handleDeliveryEvents()

	return p, nil
}

func (p *Producer) Produce(ctx context.Context, msg Message) error {
	select {
	case p.msgCh <- msg:
		return nil
	case <-p.done:
		return fmt.Errorf("kafka producer: closed")
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *Producer) Close() {
	p.closeOnce.Do(func() {
		close(p.done)
	})
}

func (p *Producer) run() {
	batch := make([]Message, 0, p.batchSize)
	ticker := time.NewTicker(p.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case msg := <-p.msgCh:
			batch = append(batch, msg)
			if len(batch) >= p.batchSize {
				p.send(batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				p.send(batch)
				batch = batch[:0]
			}
		case <-p.done:
		drain:
			for {
				select {
				case msg := <-p.msgCh:
					batch = append(batch, msg)
				default:
					break drain
				}
			}
			if len(batch) > 0 {
				p.send(batch)
			}
			p.kProducer.Flush(p.writeTimeout)
			p.kProducer.Close()
			return
		}
	}
}

func (p *Producer) send(batch []Message) {
	for _, msg := range batch {
		err := p.kProducer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &p.topic,
				Partition: kafka.PartitionAny,
			},
			Key:   []byte(msg.Key),
			Value: []byte(msg.Value),
		}, nil)
		if err != nil {
			p.log.Errorw("kafka producer: enqueue failed", "err", err, "topic", p.topic)
		}
	}
}

func (p *Producer) handleDeliveryEvents() {
	for e := range p.kProducer.Events() {
		m, ok := e.(*kafka.Message)
		if !ok {
			continue
		}
		if m.TopicPartition.Error != nil {
			p.log.Errorw("kafka producer: delivery failed",
				"err", m.TopicPartition.Error,
				"topic", *m.TopicPartition.Topic,
			)
		}
	}
}
