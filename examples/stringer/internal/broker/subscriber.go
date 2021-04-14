package broker

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/streadway/amqp"

	"github.com/marcoEgger/genki/broker"
	pb "github.com/marcoEgger/genki/examples/stringer/internal/proto"
	"github.com/marcoEgger/genki/examples/stringer/internal/stringer"
	example "github.com/marcoEgger/genki/examples/stringer/proto"
)

func GreetingHappenedSubscriber(service stringer.Service) broker.Subscriber {
	return func(delivery interface{}) {
		event := delivery.(amqp.Delivery)
		var greeting example.Greeting
		_ = proto.Unmarshal(event.Body, &greeting)

		greet := pb.GreetingFromProto(&greeting)
		_ = service.Render(context.Background(), *greet)

		_ = event.Ack(false)
	}
}
