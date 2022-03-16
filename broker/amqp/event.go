package amqp

import (
	"context"
	"fmt"
	"github.com/marcoEgger/genki/metadata"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/streadway/amqp"

	"github.com/marcoEgger/genki/broker"
	"github.com/marcoEgger/genki/broker/amqp/interceptor"
)

const (
	RequestIDHeader = "requestId"
	AccountIDHeader = "accountId"
	UserIDHeader    = "userId"
)

type Event struct {
	delivery   amqp.Delivery
	queue      string
	routingKey string
	context    context.Context
}

func NewEvent(queue, routingKey string, delivery amqp.Delivery) *Event {
	// extract amqp headers and push them into metadata
	md := metadata.Metadata{}
	for k, v := range delivery.Headers {
		md[k] = fmt.Sprint(v)
	}
	ctx := metadata.NewContext(context.Background(), md)

	return &Event{
		delivery:   delivery,
		queue:      queue,
		routingKey: routingKey,
		context:    ctx,
	}
}

func (evt *Event) Message() *broker.Message {
	return &broker.Message{
		Context: evt.context,
		Body:    evt.delivery.Body,
	}
}

func (evt *Event) SetContext(ctx context.Context) {
	evt.context = ctx
}

func (evt *Event) Ack() {
	_ = evt.delivery.Ack(false)
}

func (evt *Event) Nack(requeue bool) {
	req := "0"
	if requeue {
		req = "1"
	}

	interceptor.NackCounter.With(prometheus.Labels{
		"routing_key": evt.routingKey,
		"requeue":     req,
	}).Inc()

	_ = evt.delivery.Nack(false, requeue)
}

func (evt *Event) Reject(requeue bool) {
	req := "0"
	if requeue {
		req = "1"
	}

	interceptor.NackCounter.With(prometheus.Labels{
		"routing_key": evt.routingKey,
		"requeue":     req,
	}).Inc()

	_ = evt.delivery.Reject(requeue)
}

func (evt *Event) QueueName() string {
	return evt.queue
}

func (evt *Event) RoutingKey() string {
	return evt.routingKey
}
