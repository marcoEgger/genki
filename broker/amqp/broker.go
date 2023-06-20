package amqp

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/marcoEgger/genki/broker"
	"github.com/marcoEgger/genki/broker/amqp/interceptor"
	"github.com/marcoEgger/genki/logger"
	"github.com/marcoEgger/genki/metadata"
)

type Broker struct {
	opts           *Options
	consumeConn    *Connection
	produceConn    *Connection
	consumerDecls  []Declaration
	producerDecls  []Declaration
	consumeContext context.Context
	stopConsuming  context.CancelFunc
	waitGroup      *sync.WaitGroup
	subscriptions  map[string]broker.Handler
}

func NewBroker(options ...Option) *Broker {
	opts := newOptions(options...)

	ctx, cancel := context.WithCancel(context.Background())

	return &Broker{
		opts:           opts,
		stopConsuming:  cancel,
		consumeContext: ctx,
		subscriptions:  make(map[string]broker.Handler),
	}
}

// Initialize will setup the connections and declare all required amqp bindings for producers and consumers
func (b *Broker) Initialize() error {
	if err := b.ensureConnections(); err != nil {
		return err
	}

	// declare all the subscriber things!
	if b.HasConsumer() {
		ch, _ := b.consumeConn.Channel()
		for _, declare := range b.consumerDecls {
			if err := declare(ch); err != nil {
				return fmt.Errorf("failed to declare for consumer: %s", err.Error())
			}
		}
	}

	// declare all the consumer things!
	if len(b.producerDecls) > 0 {
		ch, _ := b.produceConn.Channel()
		for _, declare := range b.producerDecls {
			if err := declare(ch); err != nil {
				return fmt.Errorf("failed to declare for producer: %s", err.Error())
			}
		}
	}

	return nil
}

func (b *Broker) HasConsumer() bool {
	return len(b.consumerDecls) > 0
}

func (b *Broker) Consume(wg *sync.WaitGroup) {
	b.waitGroup = wg
	for {
		select {
		case <-b.consumeContext.Done():
			logger.Debug("amqp broker stopped consuming events")
			return
		default:
		}

		if !b.consumeConn.IsConnected() {
			logger.Infof("amqp consumer connection offline, waiting for reconnect")
			b.consumeConn.WaitForConnection()
			logger.Infof("amqp consumer connection back online, consuming events")
		}

		channel, err := b.consumeConn.Channel()
		if err != nil {
			logger.Warnf("unable to fetch AMQP channel for consumer: %s", err.Error())
			time.Sleep(500 * time.Millisecond)
			continue
		}

		err = channel.Qos(b.opts.PrefetchCount, 0, false)
		if err != nil {
			logger.Warnf("unable to set Qos on channel for consuming (prefetchCount=%d): %s", b.opts.PrefetchCount, err)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		deliveries, err := channel.Consume(b.opts.SubscriberQueue, b.opts.ConsumerName, false, false, false, false, nil)
		if err != nil {
			logger.Error("amqp consumer error: %s", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		for delivery := range deliveries {
			routingKey := delivery.RoutingKey

			// TODO: metrics

			if handler, ok := b.subscriptions[routingKey]; ok {
				event := NewEvent(b.opts.SubscriberQueue, routingKey, delivery)

				handler = interceptor.SubscriberLoggerInterceptor(handler)
				handler = interceptor.SubscriberMetadataInterceptor(handler)
				handler(event)
			} else {
				logger.Errorf("handler not defined for %s", routingKey)
			}
		}
	}
}

func (b *Broker) Disconnect() error {
	defer b.stopConsuming()

	if b.HasConsumer() && b.consumeConn != nil {
		defer b.waitGroup.Done()
		b.consumeConn.Shutdown()
		logger.Debug("amqp consumer connection closed")
	}
	if b.produceConn != nil {
		b.produceConn.Shutdown()
		logger.Debug("amqp producer connection closed")
	}
	logger.Info("amqp session terminated")
	return nil
}

func (b *Broker) Publish(exchange, routingKey string, message *broker.Message) error {
	pub := amqp.Publishing{
		Headers: amqp.Table{
			RequestIDHeader: metadata.GetFromContext(message.Context, metadata.RequestIDKey),
			AccountIDHeader: metadata.GetFromContext(message.Context, metadata.AccountIDKey),
			UserIDHeader:    metadata.GetFromContext(message.Context, metadata.UserIDKey),
		},
		ContentType:  "application/octet-stream",
		DeliveryMode: 0,
		Priority:     0,
		Body:         message.Body,
	}

	if !b.produceConn.IsConnected() {
		logger.Infof("amqp publish connection offline, waiting for reconnect")
		b.produceConn.WaitForConnection()
		logger.Infof("amqp publish connection back online, consuming events")
	}

	channel, err := b.produceConn.Channel()
	if err != nil {
		return err
	}

	if err := channel.Publish(exchange, routingKey, false, false, pub); err != nil {
		return errors.Wrap(err, "unable to publish event")
	}
	return nil
}

func (b *Broker) Subscribe(exchange, routingKey string, handler broker.Handler) error {
	b.consumerDecls = append(b.consumerDecls, AutoExchange(exchange))
	b.consumerDecls = append(b.consumerDecls, AutoQueue(b.opts.SubscriberQueue))
	b.consumerDecls = append(b.consumerDecls, AutoBinding(routingKey, b.opts.SubscriberQueue, exchange))
	b.subscriptions[routingKey] = handler
	logger.Infof("subscribed to events with routing key '%s' from exchange '%s'", routingKey, exchange)
	return nil
}

func (b *Broker) EnsureExchange(exchange string) {
	b.producerDecls = append(b.producerDecls, AutoExchange(exchange))
	logger.Infof("ensured amqp exchange %s", exchange)
}

// ensureConnections is responsible for creating and establishing the required connections.
// We use separate AMQP connections for publish and subscribe. The publish connection is active by default,
// whereas the consumer connection is only started if a subscriber is added.
func (b *Broker) ensureConnections() error {
	if b.HasConsumer() && b.consumeConn == nil {
		b.consumeConn = NewConnection(b.opts.Address)
		b.consumeConn.SetName("consume")
		if err := b.consumeConn.Connect(); err != nil {
			return fmt.Errorf("failed to create amqp connection: %s", err)
		}
		logger.Debug("AMQP consumer connection established")
	}
	if b.produceConn == nil {
		b.produceConn = NewConnection(b.opts.Address)
		b.produceConn.SetName("produce")
		if err := b.produceConn.Connect(); err != nil {
			return fmt.Errorf("failed to create amqp connection: %s", err)
		}
		logger.Debug("AMQP producer connection established")
	}
	logger.Infof("AMQP session alive")
	return nil
}
