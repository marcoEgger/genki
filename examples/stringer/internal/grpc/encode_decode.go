package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lukasjarosch/genki/examples/stringer/internal/models"
	example "github.com/lukasjarosch/genki/examples/stringer/proto"
)

func ErrorToProto(err error) error {
	switch err {
	case models.ErrGreetingEmptyName:
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Unknown, err.Error())
	}
}

func GreetingToProto(greeting *models.Greeting) (pbGreeting *example.Greeting) {
	pbGreeting = &example.Greeting{
		Template: greeting.Template,
		Name:     greeting.Name,
		Rendered: greeting.Rendered,
	}

	return pbGreeting
}

func GreetingFromProto(pbGreeting *example.Greeting) (greeting *models.Greeting) {
	greeting = &models.Greeting{
		Template: pbGreeting.Template,
		Name:     pbGreeting.Name,
		Rendered: pbGreeting.Rendered,
	}

	return greeting
}
