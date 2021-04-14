package interceptor

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/marcoEgger/genki/broker"
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
		Name: "amqp_not_acknowledged",
		Help: "Incremented on every NACK",
	}, []string{"routing_key", "requeue"})
)

func PrometheusInterceptor(next broker.Handler) broker.Handler {
	return func(event broker.Event) {
		InboundGauge.With(prometheus.Labels{
			"routing_key": event.RoutingKey(),
		}).Inc()

		defer func() {
			InboundGauge.With(prometheus.Labels{
				"routing_key": event.RoutingKey(),
			}).Dec()
		}()

		next(event)
	}
}

func init() {
	_ = prometheus.Register(InboundGauge)
	_ = prometheus.Register(TransportErrorCounter)
	_ = prometheus.Register(NackCounter)
}
