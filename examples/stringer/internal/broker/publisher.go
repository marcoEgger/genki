package broker

import (
	"github.com/marcoEgger/genki/broker"
	"github.com/marcoEgger/genki/examples/stringer/internal/models"
	"github.com/marcoEgger/genki/examples/stringer/internal/proto"
)

const GreetingRenderedTopic = "greeting.rendered"

func NewExamplePublisher(publisher broker.PublishProvider) *examplePublisher {
	return &examplePublisher{broker:publisher}
}

type examplePublisher struct {
	broker broker.PublishProvider
}

func (pub *examplePublisher) GreetingRendered(greeting *models.Greeting) error {
	pbGreeting := proto.GreetingToProto(greeting)
	return pub.broker.Publish(GreetingRenderedTopic, pbGreeting)
}