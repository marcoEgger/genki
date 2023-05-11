package metadata

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

// Metadata is the internal metadata abstraction. It is used to provide a single way of handling metadata
// throughout different transport layers (gRPC, HTTP, AMQP, ...).
type Metadata map[string]string

type key struct{}

func FromContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(key{}).(Metadata)
	return md, ok
}

func NewContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, key{}, md)
}

func NewInternalContext(ctx context.Context) context.Context {
	meta, _ := FromContext(ctx)
	newMeta := make(Metadata)
	replaced := false
	for k, v := range meta {
		if k == InternalKey {
			newMeta[k] = "true"
			replaced = true
		} else {
			newMeta[k] = v
		}
	}
	if !replaced {
		newMeta[InternalKey] = "true"
	}
	return NewContext(ctx, newMeta)
}

func NewOutgoingContext(ctx context.Context) context.Context {
	md := metadata.MD{}

	ctxMeta, ok := FromContext(ctx)
	if ok {
		for k, v := range ctxMeta {
			if k == EmailKey || k == FirstNameKey || k == LastNameKey {
				md.Set(k, base64.StdEncoding.EncodeToString([]byte(v)))
			} else {
				md.Set(k, v)
			}
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

// HasRole returns true if the given role string occurs in metadata.roles
func HasRole(ctx context.Context, role string) bool {
	roleString := GetFromContext(ctx, RolesKey)
	if strings.Contains(roleString, role) {
		return true
	}
	return false
}
