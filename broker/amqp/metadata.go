package amqp

import (
	"context"

	"github.com/streadway/amqp"

	"github.com/lukasjarosch/genki/metadata"
)

const (
	RequestIdHeaderKey = "requestId"
)

func MetadataFromDelivery(delivery amqp.Delivery) context.Context {
	header := delivery.Headers

	md := make(metadata.Metadata)

	if requestId, ok := header[RequestIdHeaderKey]; ok {
		md[metadata.RequestIDKey] = requestId.(string)
	}

	return metadata.NewContext(context.Background(), md)
}
