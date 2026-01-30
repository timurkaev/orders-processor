package consumer

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/timurkaev/order-processor/internal/storage"
	"log/slog"
	"time"
)

type Consumer struct {
	consumer *kafka.Consumer
	storage  storage.Storage
	logger   *slog.Logger
	topic    string
}

type Config struct {
	BootstrapServers string
	Topic            string
	GroupID          string
	Storage          storage.Storage
	Logger           *slog.Logger
}

func New(cfg Config) (*Consumer, error) {
	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers":      cfg.BootstrapServers,
		"group.id":               cfg.GroupID,
		"auto.offset.reset":      "earliest",
		"enable.auto.commit":     false,
		"session.timeout.ma":     30000,
		"heartbeat.intervals.ms": 3000,
		"max.poll.interval.ms":   300000,
		"fetch.min.bytes":        1,
		"fetch.wait.max.ms":      500,
	}

	consumer, err := kafka.NewConsumer(kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("faied to create kafka consumer: %w", err)
	}

	err = consumer.Subscribe(cfg.Topic, nil)
	if err != nil {
		consumer.Close()
		return nil, fmt.Errorf("failed to subscribe to topic %s: %w", cfg.Topic, err)
	}

	return &Consumer{
		consumer: consumer,
		storage:  cfg.Storage,
		logger:   cfg.Logger,
		topic:    cfg.Topic,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	c.logger.Info("starting kafka consumer", slog.String("topic", c.topic))

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("context cancelled, stopping consumer")
			return ctx.Err()
		default:

		}

		msg, err := c.consumer.ReadMessage(100 * time.Millisecond)

		if err != nil {
			if kafkaErr, ok := err.(kafka.Error); ok {
				if kafkaErr.Code() == kafka.ErrTimedOut {
					continue
				}
			}
			c.logger.Error("error reading message", slog.String("error", err.Error()))
			continue
		}

	}
}
