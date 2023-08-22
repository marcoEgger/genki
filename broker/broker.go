package broker

import (
	"context"
	"sync"
)

// Message is the core structure of the broker. It abstract's away the concrete broker message
// implementations (e.g. amqp.Delivery).
type Message struct {
	Context context.Context
	Body    []byte
}

// Broker is the interface which defines the client connection(s) to a concrete broker like amqp
type Broker interface {
	Publisher
	Subscriber
	Initialize() error
	Disconnect() error
	Consume(group *sync.WaitGroup)
	HasConsumer() bool
}

// Subscriber can subscribe to exchanges with a routingKey
type Subscriber interface {
	Subscribe(exchange, routingKey string, handler Handler) error
}

// Publisher can publish messages with a routingKey to the message-broker
type Publisher interface {
	Publish(exchange, routingKey string, message *Message) error
	EnsureExchange(exchange string)
}

type Handler func(event Event)

// Event abstraction
type Event interface {
	Message() *Message
	Ack()
	Nack(requeue bool)
	QueueName() string
	RoutingKey() string
	SetContext(ctx context.Context)
}
