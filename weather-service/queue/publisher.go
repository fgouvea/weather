package queue

import (
	"encoding/json"
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ErrFailedToConnect = errors.New("failed to connect to broker")
	ErrChannel         = errors.New("failed to open channel")
	ErrQueue           = errors.New("failed to declare queue")
	ErrStartConsumer   = errors.New("failed to start consumer")
	ErrPublish         = errors.New("failed to publish message")
)

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

func (p *Publisher) Publish(message string) error {
	err := p.amqpChannel.Publish(
		"",           // exchange
		p.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(message),
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

func Publish[T any](p *Publisher, message T) error {
	body, err := json.Marshal(message)

	if err != nil {
		return fmt.Errorf("%w: failed to encode message", ErrPublish)
	}

	return p.Publish(string(body))
}
