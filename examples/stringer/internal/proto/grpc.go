package proto

import (
	"context"

	"github.com/marcoEgger/genki/examples/stringer/internal/stringer"
	example "github.com/marcoEgger/genki/examples/stringer/proto"
)

type ExampleService struct {
	greeting stringer.Service
}

func NewExampleService(greetingService stringer.Service) *ExampleService {
	return &ExampleService{
		greeting: greetingService,
	}
}

func (svc *ExampleService) Hello(ctx context.Context, request *example.HelloRequest) (*example.HelloResponse, error) {
	greeting, err := svc.greeting.Hello(ctx, request.Name)
	if err != nil {
		return nil, ErrorToProto(err)
	}

	return &example.HelloResponse{
		Greeting: GreetingToProto(greeting),
	}, err
}
