package broker

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/streadway/amqp"

	"github.com/lukasjarosch/genki/broker"
	pb "github.com/lukasjarosch/genki/examples/stringer/internal/proto"
	"github.com/lukasjarosch/genki/examples/stringer/internal/service"
	example "github.com/lukasjarosch/genki/examples/stringer/proto"
)

func GreetingHappenedSubscriber(service service.Service) broker.Subscriber {
	return func(delivery interface{}) {
		event := delivery.(amqp.Delivery)
		var greeting example.Greeting
		_ = proto.Unmarshal(event.Body, &greeting)

		greet := pb.GreetingFromProto(&greeting)
		_ = service.Render(context.Background(), *greet)

		_ = event.Ack(false)
	}
}
