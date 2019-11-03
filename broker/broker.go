package broker

import "sync"

type SubscriptionCreator interface {
	AddSubscription(exchangeName, queueName, routingKey string, handler Subscriber) error
}

type PublishCreator interface {
	AddPublisher(exchangeName, routingKey string) error
}

type Broker interface {
	SubscriptionCreator
	PublishCreator
	Publish(routingKey string, event interface{}) error
	Declare() error
	Consume(wg *sync.WaitGroup)
	Shutdown()
}

type Subscriber func(delivery interface{})
