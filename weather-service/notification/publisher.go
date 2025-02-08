package notification

import (
	"encoding/json"
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

var ErrFailedToConnect = errors.New("failed to connect to broker")
var ErrChannel = errors.New("failed to open channel")
var ErrQueue = errors.New("failed to declare queue")
var ErrPublish = errors.New("failed to publish message")

type Notification struct {
	UserID  string `json:"userId"`
	Content string `json:"content"`
}

type Publisher struct {
	amqpConnection *amqp.Connection
	amqpChannel    *amqp.Channel
	queue          amqp.Queue
}

func NewPublisher(brokerURL, queueName string) (*Publisher, error) {
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

	return &Publisher{
		amqpConnection: conn,
		amqpChannel:    ch,
		queue:          q,
	}, nil
}

func (p *Publisher) Notify(userId, content string) error {
	notification := Notification{
		UserID:  userId,
		Content: content,
	}

	body, err := json.Marshal(notification)

	if err != nil {
		return fmt.Errorf("%w: failed to encode notification", ErrPublish)
	}

	err = p.amqpChannel.Publish(
		"",           // exchange
		p.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		},
	)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrPublish, err)
	}

	return nil
}

func (p *Publisher) Close() {
	p.amqpChannel.Close()
	p.amqpConnection.Close()
}
