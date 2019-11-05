package metadata

import (
	"context"

	"github.com/google/uuid"
)

// Metadata is the internal metadata abstraction. It is used to provide a single way of handling metadata
// throughout different transport layers (gRPC, HTTP, AMQP, ...).
type Metadata map[string]string

type key struct {}

func FromContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(key{}).(Metadata)
	return md, ok
}

func NewContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, key{}, md)
}

func NewRequestID() string {
	return uuid.New().String()
}
