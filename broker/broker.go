package broker

import (
	"context"
	"sync"
)

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

type Subscriber func(delivery interface{})

type Publisher func(ctx context.Context,  string, event interface{}) error

type PublishProvider interface {
	Publish(ctx context.Context, routingKey string, event interface{}) error
}
