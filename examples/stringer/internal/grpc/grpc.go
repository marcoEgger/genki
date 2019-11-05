package grpc

import (
	"context"

	"github.com/lukasjarosch/genki/examples/stringer/internal/greeting"
	example "github.com/lukasjarosch/genki/examples/stringer/proto"
)

type ExampleService struct {
	greeting greeting.Service
}

func NewExampleService(greetingService greeting.Service) *ExampleService {
	return &ExampleService{
		greeting: greetingService,
	}
}

func (svc *ExampleService) Hello(ctx context.Context, request *example.HelloRequest) (*example.HelloResponse, error) {
	greeting, err := svc.greeting.Hello(ctx, request.Name)

	return &example.HelloResponse{
		Greeting: GreetingToProto(greeting),
	}, err
}
