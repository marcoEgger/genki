package amqp

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/streadway/amqp"

	"github.com/lukasjarosch/genki/broker"
	"github.com/lukasjarosch/genki/broker/amqp/interceptor"
)

const (
	RequestIDHeader = "requestId"
)

type Event struct {
	delivery   amqp.Delivery
	queue      string
	routingKey string
	context    context.Context
}

func NewEvent(queue, routingKey string, delivery amqp.Delivery) *Event {
	return &Event{
		delivery:   delivery,
		queue:      queue,
		routingKey: routingKey,
		context:    context.Background(),
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

func (evt *Event) Nack(retry bool) {
	req := "0"
	if retry {
		req = "1"
	}

	interceptor.NackCounter.With(prometheus.Labels{
		"routing_key": evt.routingKey,
		"requeue": req,
	}).Inc()

	_ = evt.delivery.Nack(false, retry)
}

func (evt *Event) QueueName() string {
	return evt.queue
}

func (evt *Event) RoutingKey() string {
	return evt.routingKey
}
