package proto

import (
	"context"

	"github.com/lukasjarosch/genki/examples/stringer/internal/service"
	example "github.com/lukasjarosch/genki/examples/stringer/proto"
)

type ExampleService struct {
	greeting service.Service
}

func NewExampleService(greetingService service.Service) *ExampleService {
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
