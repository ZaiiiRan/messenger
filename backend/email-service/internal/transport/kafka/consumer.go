package kafka

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ZaiiiRan/messenger/backend/email-service/internal/config/settings"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

type Message struct {
	Key  string
	Body string
}

type batchMsg struct {
	raw *kafka.Message
	val Message
}

type Consumer struct {
	kConsumer    *kafka.Consumer
	topic        string
	msgCh        chan batchMsg
	done         chan struct{}
	closeOnce    sync.Once
	batchSize    int
	batchTimeout time.Duration
	log          *zap.SugaredLogger
}

func NewConsumer(cfg settings.KafkaConsumerSettings, kafkaClient *KafkaClient, log *zap.SugaredLogger) (*Consumer, error) {
	cm := kafkaClient.ConfigMap()
	cm["group.id"] = cfg.GroupID
	cm["auto.offset.reset"] = "earliest"
	cm["enable.auto.commit"] = false

	kConsumer, err := kafka.NewConsumer(&cm)
	if err != nil {
		return nil, fmt.Errorf("kafka consumer: create: %w", err)
	}

	if err := kConsumer.Subscribe(cfg.Topic, nil); err != nil {
		_ = kConsumer.Close()
		return nil, fmt.Errorf("kafka consumer: subscribe: %w", err)
	}

	return &Consumer{
		kConsumer:    kConsumer,
		topic:        cfg.Topic,
		msgCh:        make(chan batchMsg, int(cfg.BatchSize)*2),
		done:         make(chan struct{}),
		batchSize:    int(cfg.BatchSize),
		batchTimeout: time.Duration(cfg.BatchTimeoutMs) * time.Millisecond,
		log:          log,
	}, nil
}

func (c *Consumer) Close() {
	c.closeOnce.Do(func() {
		close(c.done)
	})
}

func (c *Consumer) Run(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.readLoop(ctx)
	}()

	c.batchLoop()
	wg.Wait()
	_ = c.kConsumer.Close()
}

func (c *Consumer) readLoop(ctx context.Context) {
	defer c.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.done:
			return
		default:
		}

		ev, err := c.kConsumer.ReadMessage(100 * time.Millisecond)
		if err != nil {
			if kerr, ok := err.(kafka.Error); ok && kerr.Code() == kafka.ErrTimedOut {
				continue
			}
			c.log.Errorw("kafka_consumer.read_error", "err", err, "topic", c.topic)
			continue
		}

		msg := batchMsg{
			raw: ev,
			val: Message{Key: string(ev.Key), Body: string(ev.Value)},
		}

		select {
		case c.msgCh <- msg:
		case <-ctx.Done():
			return
		case <-c.done:
			return
		}
	}
}

func (c *Consumer) batchLoop() {
	batch := make([]batchMsg, 0, c.batchSize)
	ticker := time.NewTicker(c.batchTimeout)
	defer ticker.Stop()

	for {
		select {
		case msg := <-c.msgCh:
			batch = append(batch, msg)
			if len(batch) >= c.batchSize {
				c.process(batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				c.process(batch)
				batch = batch[:0]
			}
		case <-c.done:
		drain:
			for {
				select {
				case msg := <-c.msgCh:
					batch = append(batch, msg)
				default:
					break drain
				}
			}
			if len(batch) > 0 {
				c.process(batch)
			}
			return
		}
	}
}

func (c *Consumer) process(batch []batchMsg) {
	for _, msg := range batch {
		c.log.Infow("kafka_consumer.received_message", "key", msg.val.Key, "body", msg.val.Body)
	}
	last := batch[len(batch)-1]
	if _, err := c.kConsumer.CommitMessage(last.raw); err != nil {
		c.log.Errorw("kafka_consumer.commit_failed", "err", err, "topic", c.topic)
	}
}
