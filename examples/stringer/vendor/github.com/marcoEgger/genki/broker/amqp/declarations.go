package amqp

import (
	"github.com/streadway/amqp"
)

type Declaration func(Declarator) error

// Declarator is implemented by amqp.Channel
type Declarator interface {
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
}

type Binding struct {
	exchange   Exchange
	queue      Queue
	routingKey string
	args       amqp.Table
}

type Exchange struct {
	name       string
	kind       string
	durable    bool
	autoDelete bool
	exclusive  bool
	noWait     bool
	args       amqp.Table
}

type Queue struct {
	name       string
	durable    bool
	autoDelete bool
	exclusive  bool
	noWait     bool
	args       amqp.Table
}

func AutoQueue(name string) Declaration {
	return DeclareQueue(&Queue{
		name:       name,
		durable:    true,
		autoDelete: false,
		exclusive:  false,
		noWait:     false,
		args:       nil,
	})
}

func DeclareQueue(q *Queue) Declaration {
	return func(d Declarator) error {
		_, err := d.QueueDeclare(
			q.name,
			q.durable,
			q.autoDelete,
			q.exclusive,
			q.noWait,
			q.args,
		)
		return err
	}
}

func AutoExchange(name string) Declaration {
	return DeclareExchange(&Exchange{
		name:       name,
		kind:       "topic",
		durable:    true,
		autoDelete: false,
		exclusive:  false,
		noWait:     false,
		args:       nil,
	})
}

func DeclareExchange(e *Exchange) Declaration {
	return func(d Declarator) error {
		return d.ExchangeDeclare(
			e.name,
			e.kind,
			e.durable,
			e.exclusive,
			e.autoDelete,
			e.noWait,
			e.args,
		)
	}
}

func AutoBinding(routingKey, queue, exchange string) Declaration {
	return DeclareBinding(&Binding{
		exchange:   Exchange{name: exchange},
		queue:      Queue{name: queue},
		routingKey: routingKey,
		args:       nil,
	})
}

func DeclareBinding(b *Binding) Declaration {
	return func(d Declarator) error {
		return d.QueueBind(
			b.queue.name,
			b.routingKey,
			b.exchange.name,
			false,
			b.args,
		)
	}
}

type Delivery struct {
	amqp.Delivery
}

func (d *Delivery) Ack(multiple bool) {
	_ = d.Delivery.Ack(multiple)
}

func (d *Delivery) Nack(multiple, requeue bool) {
	_ = d.Delivery.Nack(multiple, requeue)
}
