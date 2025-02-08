package queue

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fgouvea/weather/notification-service/notification"
	"github.com/fgouvea/weather/notification-service/user"
	"github.com/google/uuid"
	"go.uber.org/zap"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ErrFailedToConnect = errors.New("failed to connect to broker")
	ErrChannel         = errors.New("failed to open channel")
	ErrQueue           = errors.New("failed to declare queue")
	ErrStartConsumer   = errors.New("failed to start consumer")
)

type NotificationProcessor interface {
	Process(n notification.Notification) error
}

type NotificationConsumer struct {
	Processor NotificationProcessor
	Consumers int
	Logger    *zap.Logger

	amqpConnection *amqp.Connection
	amqpChannel    *amqp.Channel
	queue          amqp.Queue
}

func NewConsumer(brokerURL, queueName string, consumers int, processor NotificationProcessor, logger *zap.Logger) (*NotificationConsumer, error) {
	conn, err := amqp.Dial(brokerURL)

	if err != nil {
		return nil, fmt.Errorf("%w on %s: %w", ErrFailedToConnect, brokerURL, err)
	}

	ch, err := conn.Channel()

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrChannel, err)
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
		return nil, fmt.Errorf("%w: %w", ErrQueue, err)
	}

	return &NotificationConsumer{
		Processor: processor,
		Consumers: consumers,
		Logger:    logger,

		amqpConnection: conn,
		amqpChannel:    ch,
		queue:          q,
	}, nil
}

func (c *NotificationConsumer) consume(consumerName string, delivery amqp.Delivery) {
	var userNotification notification.Notification

	err := json.Unmarshal(delivery.Body, &userNotification)

	if err != nil {
		c.Logger.Error("error reading message body", zap.String("consumer", consumerName), zap.String("body", string(delivery.Body)))
		delivery.Nack(false, false)
		return
	}

	err = c.Processor.Process(userNotification)

	if errors.Is(user.ErrUserNotFound, err) {
		c.Logger.Error("user does not exist", zap.String("consumer", consumerName), zap.String("userID", string(userNotification.UserID)))
		delivery.Nack(false, false)
		return
	}

	if err != nil {
		c.Logger.Error("error processing notification", zap.String("consumer", consumerName), zap.String("userID", string(userNotification.UserID)), zap.Error(err))
		delivery.Nack(false, true)
		return
	}

	delivery.Ack(false)
}

func (c *NotificationConsumer) Start() error {
	for i := 0; i < c.Consumers; i++ {
		consumerName := fmt.Sprintf("%s-consumer-%s", c.queue.Name, uuid.New())

		msgs, err := c.amqpChannel.Consume(
			c.queue.Name,
			consumerName,
			false,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			return fmt.Errorf("%w: %w", ErrStartConsumer, err)
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
