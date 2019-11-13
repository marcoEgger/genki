package amqp

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/streadway/amqp"
)

var (
	InboundGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "amqp_inbound",
		Help:        "Increased on incoming deliveries, decreased on ack/nack",
		ConstLabels: nil,
	}, []string{"routing_key"})
	TransportErrorCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        "amqp_transport_error",
		Help:        "Increased when a message could not be decoded",
		ConstLabels: nil,
	}, []string{"routing_key"})
	NackCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        "amqp_not_acknowledged",
		Help:        "Incremented on every NACK",
	}, []string{"routing_key", "requeue"})
)

func init() {
	_ = prometheus.Register(InboundGauge)
	_ = prometheus.Register(TransportErrorCounter)
	_ = prometheus.Register(NackCounter)
}

func Nack(delivery amqp.Delivery, multiple, requeue bool) error {
	InboundGauge.With(prometheus.Labels{
		"routing_key": delivery.RoutingKey,
	}).Dec()

	req := "0"
	if requeue {
		req = "1"
	}

	NackCounter.With(prometheus.Labels{
		"routing_key": delivery.RoutingKey,
		"requeue": req,
	}).Inc()
	return delivery.Nack(multiple, requeue)
}

func Ack(delivery amqp.Delivery, multiple bool) error {
	InboundGauge.With(prometheus.Labels{
		"routing_key": delivery.RoutingKey,
	}).Dec()
	return delivery.Ack(multiple)
}
