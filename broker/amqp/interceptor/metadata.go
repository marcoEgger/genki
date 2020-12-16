package interceptor

import (
	"github.com/lukasjarosch/genki/broker"
	"github.com/lukasjarosch/genki/metadata"
)

func SubscriberMetadataInterceptor(next broker.Handler) broker.Handler {
	return func(event broker.Event) {
		meta := metadata.Metadata{}
		ensureRequestId(&meta, event)
		findAccountID(&meta, event)

		ctx := metadata.NewContext(event.Message().Context, meta)
		event.SetContext(ctx)

		next(event)
	}
}

func ensureRequestId(meta *metadata.Metadata, event broker.Event) {
	requestID := metadata.GetFromContext(event.Message().Context, metadata.RequestIDKey)
	if requestID == "" {
		(*meta)[metadata.RequestIDKey] = metadata.NewRequestID()
		return
	}
	(*meta)[metadata.RequestIDKey] = requestID
}


func findAccountID(meta *metadata.Metadata, event broker.Event) {
	accID := metadata.GetFromContext(event.Message().Context, metadata.AccountIDKey)
	if accID != "" {
		(*meta)[metadata.AccountIDKey] = accID
	}
}