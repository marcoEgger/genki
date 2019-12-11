package metadata

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
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

func NewOutgoingContext(ctx context.Context) context.Context {
	md := metadata.MD{}

	 ctxMeta, ok := FromContext(ctx)
	 if ok {
		 for key, value := range ctxMeta {
			md.Set(key, value)
		 }
	 }
	 outCtx := metadata.NewOutgoingContext(ctx, md)
	 return outCtx
}

func NewRequestID() string {
	return uuid.New().String()
}

func GetFromContext(ctx context.Context, key string) string {
	md, ok := FromContext(ctx)
	if !ok {
		return ""
	}
	if val, ok := md[key]; ok {
		return val
	}
	return ""
}