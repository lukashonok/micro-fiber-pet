package mq

import (
	"fmt"
	"log"

	"github.com/lukashonok/micro-fiber-pet/pkg/book"
	"github.com/lukashonok/micro-fiber-pet/pkg/firebase"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

// RabbitMQ handler will pass services to every handler which takes only package for request
type MqHandler func(DefaultServices) func(amqp.Delivery)

type DefaultServices struct {
	BookPublisher *Publisher
	RabbitMQ      RabbitMQ
	BookService   book.Service
	Db            *mongo.Database
	FirebaseAuth  *firebase.Auth
}

type Consumer struct {
	rabbitmq RabbitMQ
	queue    string
	msgs     <-chan amqp.Delivery
	handlers map[string]MqHandler
	services DefaultServices
}

// NewConsumer creates a new RabbitMQ consumer bound to a queue
func NewConsumer(rabbitmq RabbitMQ, exchange string, queue string, services DefaultServices, keysHandler map[string]MqHandler) (*Consumer, error) {
	// Declare queue
	q, err := rabbitmq.Channel.QueueDeclare(
		queue,
		false, // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)

	if err != nil {
		rabbitmq.Connection.Close()
		rabbitmq.Channel.Close()
		return nil, fmt.Errorf("declare queue: %w", err)
	}

	// creating only bindings to queue without touching handlers
	for key := range keysHandler {
		if err := rabbitmq.Channel.QueueBind(q.Name, key, exchange, false, nil); err != nil {
			rabbitmq.Connection.Close()
			rabbitmq.Channel.Close()
			return nil, fmt.Errorf("queue bind: %w", err)
		}
	}

	msgs, err := rabbitmq.Channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		rabbitmq.Connection.Close()
		rabbitmq.Channel.Close()
		return nil, fmt.Errorf("consume: %w", err)
	}

	return &Consumer{rabbitmq: rabbitmq, msgs: msgs, queue: q.Name, services: services, handlers: keysHandler}, nil
}

func (c *Consumer) Start() {
	go func() {
		for msg := range c.msgs {
			handler, ok := c.handlers[msg.RoutingKey]
			if !ok {
				log.Printf("No handler for this action: %s", msg.RoutingKey)
				continue
			}
			handler(c.services)(msg)
		}
	}()
}

// Close closes the consumer Connection
func (c *Consumer) Close() {
	if c.rabbitmq.Channel != nil {
		c.rabbitmq.Channel.Close()
	}
	if c.rabbitmq.Connection != nil {
		c.rabbitmq.Connection.Close()
	}
}
