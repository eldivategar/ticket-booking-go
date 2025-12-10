package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher interface {
	GetConnection() *amqp.Connection
	CreateQueue(name string) error
	Publish(ctx context.Context, queueName string, payload interface{}) error
	Close()
}

type rabbitMQPublisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQPublisher(url string) (*rabbitMQPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &rabbitMQPublisher{conn: conn, ch: ch}, nil
}

func (r *rabbitMQPublisher) GetConnection() *amqp.Connection {
	return r.conn
}

func (r *rabbitMQPublisher) CreateQueue(name string) error {
	_, err := r.ch.QueueDeclare(
		name,  // name of the queue
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

func (r *rabbitMQPublisher) Publish(ctx context.Context, queueName string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	return r.ch.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
}

func (r *rabbitMQPublisher) Close() {
	r.ch.Close()
	r.conn.Close()
}
