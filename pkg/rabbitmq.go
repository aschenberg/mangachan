package pkg

import (
	"fmt"
	"log"
	"manga/config"

	"github.com/streadway/amqp"
)

// RabbitMQ ...
type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// NewRabbitMQ instantiates the RabbitMQ instances using configuration defined in environment variables.
func NewRabbitMQ(cfg *config.Config, queueNames ...string) (*RabbitMQ, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s",
		cfg.RabbitMq.User,
		cfg.RabbitMq.Password,
		cfg.RabbitMq.Host,
		cfg.RabbitMq.Port)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, WrapErrorf(err, ErrorCodeUnknown, "amqp.Dial")
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		"task",  // exchange name
		"topic", // exchange type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return nil, WrapErrorf(err, ErrorCodeUnknown, "ch.ExchangeDeclare")
	}

	if err := ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		return nil, WrapErrorf(err, ErrorCodeUnknown, "ch.Qos")
	}

	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}, nil
}

// Publish sends a message to a specified RabbitMQ queue
func (r *RabbitMQ) Publish(queueName, body string) error {

	err := r.Channel.Publish(
		"",      // exchange
		"manga", // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	return err
}

// Consume listens for messages from a specified RabbitMQ queue and processes them with a handler
func (r *RabbitMQ) Consume(queueName string, handler func(amqp.Delivery)) error {

	msgs, err := r.Channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			handler(d)
		}
	}()

	return nil
}

// Close ...
func (r *RabbitMQ) Close() {
	r.Conn.Close()
	r.Channel.Close()
}
