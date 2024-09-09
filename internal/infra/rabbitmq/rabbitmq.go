package rabbitmq

import (
	"bytes"
	"context"
	"encoding/gob"
	"manga/internal/domain/dtos"
	"manga/pkg"
	"time"

	"github.com/streadway/amqp"
)

// Task represents the repository used for publishing Task records.
type Task struct {
	ch *amqp.Channel
}

// NewTask instantiates the Task repository.
func NewTask(channel *amqp.Channel) (*Task, error) {
	return &Task{
		ch: channel,
	}, nil
}

// Created publishes a message indicating a task was created.
func (t *Task) CreatedManga(ctx context.Context, task dtos.IndexedManga) error {
	return t.publish(ctx, "Task.Created", "tasks.event.created", task)
}

// Deleted publishes a message indicating a task was deleted.
func (t *Task) Deleted(ctx context.Context, id string) error {
	return t.publish(ctx, "Task.Deleted", "tasks.event.deleted", id)
}


func (t *Task) publish(ctx context.Context, spanName, routingKey string, event interface{}) error {
	// _, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, spanName)
	// defer span.End()

	// span.SetAttributes(
	// 	attribute.KeyValue{
	// 		Key:   semconv.MessagingSystemKey,
	// 		Value: attribute.StringValue("rabbitmq"),
	// 	},
	// 	attribute.KeyValue{
	// 		Key:   semconv.MessagingRabbitMQRoutingKeyKey,
	// 		Value: attribute.StringValue(routingKey),
	// 	},
	// )

	//-

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(event); err != nil {
		return pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "gob.Encode")
	}

	err := t.ch.Publish(
		"tasks",    // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			AppId:       "tasks-rest-server",
			ContentType: "application/x-encoding-gob", // XXX: We will revisit this in future episodes
			Body:        b.Bytes(),
			Timestamp:   time.Now(),
		})
	if err != nil {
		return pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "ch.Publish")
	}

	return nil
}
