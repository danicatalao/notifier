package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	l "github.com/danicatalao/notifier/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Url            string
	ExchangeName   string
	ExchangeType   string
	ReconnectDelay time.Duration
	MaxRetries     int
}

type Message struct {
	RoutingKey string
	Body       interface{}
}

type DeliveryMessage struct {
	Body []byte
	Ack  func(multiple bool) error
	Nack func(multiple bool, requeue bool) error
}

type Service interface {
	Publish(ctx context.Context, msg Message) error
	Consume(ctx context.Context, queueName string) (<-chan DeliveryMessage, error)
	Close() error
}

type service struct {
	config Config
	conn   *amqp.Connection
	ch     *amqp.Channel
	mu     sync.RWMutex
	closed bool
	log    l.Logger
}

func NewService(config Config, l l.Logger) (*service, error) {
	if config.ReconnectDelay == 0 {
		config.ReconnectDelay = 5 * time.Second
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 5
	}

	s := &service{
		config: config,
		log:    l,
	}

	if err := s.connect(); err != nil {
		return nil, err
	}

	go s.monitorConnection()

	return s, nil
}

func (s *service) connect() error {
	conn, err := amqp.Dial(s.config.Url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %v", err)
	}

	err = ch.ExchangeDeclare(
		s.config.ExchangeName,
		s.config.ExchangeType,
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	s.mu.Lock()
	s.conn = conn
	s.ch = ch
	s.mu.Unlock()

	return nil
}

func (s *service) monitorConnection() {
	for {
		s.mu.RLock()
		if s.closed {
			s.mu.RUnlock()
			return
		}
		s.mu.RUnlock()

		connErrChan := make(chan *amqp.Error)
		s.conn.NotifyClose(connErrChan)

		err := <-connErrChan

		if err != nil {
			for i := 0; i < s.config.MaxRetries; i++ {
				if err := s.connect(); err == nil {
					break
				}
				time.Sleep(s.config.ReconnectDelay)
			}
		}
	}
}

func (s *service) Publish(ctx context.Context, msg Message) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return fmt.Errorf("service is closed")
	}

	body, err := json.Marshal(msg.Body)
	if err != nil {
		return fmt.Errorf("failed to marshal message body: %v", err)
	}

	err = s.ch.PublishWithContext(ctx,
		s.config.ExchangeName,
		msg.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}

func (s *service) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true

	var errors []error
	if s.ch != nil {
		if err := s.ch.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close channel: %v", err))
		}
	}

	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close connection: %v", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors closing service: %v", errors)
	}

	return nil
}

func (s *service) Consume(ctx context.Context, queueName string) (<-chan DeliveryMessage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, fmt.Errorf("service is closed")
	}

	// Declare queue if it doesn't exist
	q, err := s.ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	// Bind queue to exchange
	err = s.ch.QueueBind(
		q.Name,
		q.Name, // Using queue name as routing key
		s.config.ExchangeName,
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %v", err)
	}

	// Start consuming
	deliveries, err := s.ch.Consume(
		q.Name,
		"",    // consumer tag
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume from queue: %v", err)
	}

	// Create a channel to pass converted messages
	deliveryChan := make(chan DeliveryMessage)

	// Start a goroutine to handle context cancellation and message conversion
	go func() {
		defer close(deliveryChan)

		for {
			select {
			case <-ctx.Done():
				s.log.InfoContext(ctx, "context cancelled, stopping consumption")
				return
			case d, ok := <-deliveries:
				if !ok {
					s.log.InfoContext(ctx, "deliveries channel closed, stopping consumption")
					return
				}

				// Convert amqp.Delivery to DeliveryMessage
				deliveryChan <- DeliveryMessage{
					Body: d.Body,
					Ack: func(multiple bool) error {
						return d.Ack(multiple)
					},
					Nack: func(multiple bool, requeue bool) error {
						return d.Nack(multiple, requeue)
					},
				}
			}
		}
	}()

	return deliveryChan, nil
}
