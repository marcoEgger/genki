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
	Initialize() error
	Consume(group *sync.WaitGroup)
	Disconnect() error
	Subscribe(exchange, routingKey string, handler Handler) error
}

// Publisher can publish messages with a routingKey to the message-broker
type Publisher interface {
	Publish(exchange, routingKey string, message *Message) error
}

type Handler func(event Event)

// Event abstraction
type Event interface {
	Message() *Message
	Ack()
	Nack(retry bool)
	QueueName() string
	RoutingKey() string
	SetContext(ctx context.Context)
}

/*
type SubscriptionCreator interface {
	AddSubscription(exchangeName, queueName, routingKey string, handler Subscriber) error
}

type PublishCreator interface {
	AddPublisher(exchangeName, routingKey string) error
}

type Broker interface {
	PublishProvider
	SubscriptionCreator
	PublishCreator
	Declare() error
	Consume(wg *sync.WaitGroup)
	Shutdown()
}

type Subscriber func(ctx context.Context, delivery amqp.Delivery)

type Publisher func(ctx context.Context, routingKey string, event interface{}) error

type PublishProvider interface {
	Publish(ctx context.Context, routingKey string, event interface{}) error
}

type PublishInterceptor func(next Publisher) Publisher
type SubscriptionInterceptor func(next Subscriber) Subscriber

*/
