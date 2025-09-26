package mq

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	rabbitmq RabbitMQ
	exchange string
}

func NewPublisher(rabbitmq RabbitMQ, exchange string) (*Publisher, error) {
	if err := rabbitmq.Channel.ExchangeDeclare(
		exchange,
		"topic",
		false, // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	); err != nil {
		rabbitmq.Channel.Close()
		rabbitmq.Connection.Close()
		return nil, err
	}

	return &Publisher{rabbitmq: rabbitmq, exchange: exchange}, nil
}

func (p *Publisher) Publish(routingKey string, payload interface{}) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return p.rabbitmq.Channel.Publish(
		p.exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		},
	)
}

// Close closes the consumer Connection
func (c *Publisher) Close() {
	if c.rabbitmq.Channel != nil {
		c.rabbitmq.Channel.Close()
	}
	if c.rabbitmq.Connection != nil {
		c.rabbitmq.Connection.Close()
	}
}
