package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"manga/config"
	"manga/internal/domain/dtos"
	"manga/internal/infra/meili"
	"manga/pkg"
	"time"

	"manga/pkg/logging"

	"os"
	"os/signal"
	"syscall"
)


const rabbitMQConsumerName = "meilisearch-indexer"

func main() {
	cfg := config.NewConfig()
	



	errC, err := run(cfg)
	if err != nil {
		log.Fatalf("Couldn't run: %s", err)
	}

	if err := <-errC; err != nil {
		log.Fatalf("Error while running: %s", err)
	}

}

func run(cfg *config.Config) (<-chan error, error)  {
	log := logging.NewLogger(cfg)

	rmq,err:= pkg.NewRabbitMQ(cfg)
	if err != nil {
		return nil, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "internal.newRabbitMQ")
	}
	srv := &Server{
		log: log,
		rmq:rmq,
	
		done:   make(chan struct{}),
	}

	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

		go func() {
			<-ctx.Done()
	
			log.Info(logging.General,logging.Shutdown,"Shutdown signal received",nil)
	
			ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	
			defer func() {
				// _ = log.Fatal()
	
				rmq.Close()
				stop()
				cancel()
				close(errC)
			}()
	
			if err := srv.Shutdown(ctxTimeout); err != nil { //nolint: contextcheck
				errC <- err
			}
	
			log.Info(logging.General,logging.Shutdown,"Shutdown completed",nil)
		}()
	
		go func() {
			log.Info(logging.General,logging.Startup,"Listening and serving",nil)
	
			if err := srv.ListenAndServe(); err != nil {
				errC <- err
			}
		}()

		return errC, nil
}



type Server struct {
	log    logging.Logger
	rmq    *pkg.RabbitMQ
	meili  *meili.Manga
	done   chan struct{}
}

func (s *Server) ListenAndServe() error {
	queue, err := s.rmq.Channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "channel.QueueDeclare")
	}

	err = s.rmq.Channel.QueueBind(
		queue.Name,      // queue name
		"tasks.event.*", // routing key
		"tasks",         // exchange
		false,
		nil,
	)
	if err != nil {
		return pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "channel.QueueBind")
	}

	msgs, err := s.rmq.Channel.Consume(
		queue.Name,           // queue
		rabbitMQConsumerName, // consumer
		false,                // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	if err != nil {
		return pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "channel.Consume")
	}

	go func() {
		for msg := range msgs {
			s.log.Info(logging.Rabbit,logging.Received,fmt.Sprintf("Received message: %s", msg.RoutingKey),nil)

			var nack bool

			// XXX: Instrumentation to be added in a future episode

			// XXX: We will revisit defining these topics in a better way in future episodes
			switch msg.RoutingKey {
			case "tasks.event.updated", "tasks.event.created":
				task, err := decodeData[dtos.IndexedManga](msg.Body)
				if err != nil {
					return
				}

				if err := s.meili.Index(context.Background(), task); err != nil {
					nack = true
				}
			case "tasks.event.deleted":
				id, err := decodeID(msg.Body)
				if err != nil {
					return
				}

				if err := s.meili.Delete(context.Background(), id); err != nil {
					nack = true
				}
			default:
				nack = true
			}

			if nack {
				s.log.Info(logging.Rabbit,logging.Received,"NAcking :(",nil)

				_ = msg.Nack(false, nack)
			} else {
				s.log.Info(logging.Rabbit,logging.Received,"Acking :)",nil)

				_ = msg.Ack(false)
			}
		}

		s.log.Info(logging.Rabbit,logging.Received,"No more messages to consume. Exiting.",nil)

		s.done <- struct{}{}
	}()

	return nil
}


// Shutdown ...
func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info(logging.General,logging.Shutdown,"Shutting down server",nil)

	_ = s.rmq.Channel.Cancel(rabbitMQConsumerName, false)

	for {
		select {
		case <-ctx.Done():
			return pkg.WrapErrorf(ctx.Err(), pkg.ErrorCodeUnknown, "context.Done")
		case <-s.done:
			return nil
		}
	}
}



func decodeData[T any](b []byte) (T, error) {
	var res T

	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&res); err != nil {
		var zero T
		return zero, pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "gob.Decode")
	}

	return res, nil
}
func decodeID(b []byte) (string, error) {
	var res string

	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&res); err != nil {
		return "", pkg.WrapErrorf(err, pkg.ErrorCodeUnknown, "gob.Decode")
	}

	return res, nil
}
