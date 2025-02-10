package schedule

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fgouvea/weather/weather-service/queue"
	"github.com/fgouvea/weather/weather-service/user"
	"github.com/fgouvea/weather/weather-service/weather"
	"github.com/google/uuid"
	"go.uber.org/zap"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ScheduleProcessor interface {
	Process(schedule Schedule) error
}

type Consumer struct {
	Processor ScheduleProcessor
	Consumers int
	Logger    *zap.Logger

	amqpConnection  *amqp.Connection
	amqpReadChannel *amqp.Channel
	queue           amqp.Queue
}

func NewConsumer(brokerURL, queueName string, consumers int, processor ScheduleProcessor, logger *zap.Logger) (*Consumer, error) {
	conn, err := amqp.Dial(brokerURL)

	if err != nil {
		return nil, fmt.Errorf("%w on %s: %w", queue.ErrFailedToConnect, brokerURL, err)
	}

	ch, err := conn.Channel()

	if err != nil {
		return nil, fmt.Errorf("%w: %w", queue.ErrChannel, err)
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", queue.ErrQueue, err)
	}

	return &Consumer{
		Processor: processor,
		Consumers: consumers,
		Logger:    logger,

		amqpConnection:  conn,
		amqpReadChannel: ch,
		queue:           q,
	}, nil
}

func (c *Consumer) consume(consumerName string, delivery amqp.Delivery) {
	var schedule Schedule

	err := json.Unmarshal(delivery.Body, &schedule)

	if err != nil {
		c.Logger.Error("error reading message body", zap.String("consumer", consumerName), zap.String("body", string(delivery.Body)))
		delivery.Nack(false, false)
		return
	}

	err = c.Processor.Process(schedule)

	if errors.Is(user.ErrUserNotFound, err) || errors.Is(weather.ErrCityNotFound, err) {
		c.Logger.Error("non retryable error processing schedule", zap.String("consumer", consumerName), zap.String("userID", string(schedule.UserID)), zap.Error(err))
		delivery.Nack(false, false)
		return
	}

	if err != nil {
		c.Logger.Error("error processing schedule", zap.String("consumer", consumerName), zap.String("userID", string(schedule.UserID)), zap.Error(err))
		delivery.Nack(false, true)
		return
	}

	delivery.Ack(false)
}

func (c *Consumer) Start() error {
	for i := 0; i < c.Consumers; i++ {
		consumerName := fmt.Sprintf("%s-consumer-%s", c.queue.Name, uuid.New())

		msgs, err := c.amqpReadChannel.Consume(
			c.queue.Name,
			consumerName,
			false,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			return fmt.Errorf("%w: %w", queue.ErrStartConsumer, err)
		}

		go func(consumer string) {
			c.Logger.Info("starting consumer", zap.String("consumer", consumerName))

			for delivery := range msgs {
				c.consume(consumer, delivery)
			}
		}(consumerName)
	}

	return nil
}
