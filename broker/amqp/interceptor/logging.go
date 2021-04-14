package interceptor

import (
	"time"

	"github.com/marcoEgger/genki/broker"
	"github.com/marcoEgger/genki/logger"
)

func SubscriberLoggerInterceptor(next broker.Handler) broker.Handler {
	return func(event broker.Event) {
		log := logger.WithMetadata(event.Message().Context)
		log.Infof("incoming event with routing key '%s'", event.RoutingKey())

		defer func(started time.Time) {
			log = log.WithFields(logger.Fields{
				"took": time.Since(started),
			})
			log.Infof("event handling finished")
		}(time.Now())

		next(event)
	}
}
