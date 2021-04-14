package interceptor

import (
	b64 "encoding/base64"

	"github.com/marcoEgger/genki/broker"
	"github.com/marcoEgger/genki/metadata"
)

func SubscriberMetadataInterceptor(next broker.Handler) broker.Handler {
	return func(event broker.Event) {
		meta := metadata.Metadata{}
		ensureRequestId(&meta, event)
		findAccountID(&meta, event)
		findUserID(&meta, event)
		findEmail(&meta, event)
		findFirstName(&meta, event)
		findLastName(&meta, event)
		findType(&meta, event)
		findSubType(&meta, event)
		findRoles(&meta, event)
		findInternal(&meta, event)

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

func findUserID(meta *metadata.Metadata, event broker.Event) {
	userID := metadata.GetFromContext(event.Message().Context, metadata.UserIDKey)
	if userID != "" {
		(*meta)[metadata.UserIDKey] = userID
	}
}

func findEmail(meta *metadata.Metadata, event broker.Event) {
	email := metadata.GetFromContext(event.Message().Context, metadata.EmailKey)
	if email != "" {
		decoded, err := b64.StdEncoding.DecodeString(email)
		if err != nil {
			return
		}
		(*meta)[metadata.EmailKey] = string(decoded)
	}
}

func findFirstName(meta *metadata.Metadata, event broker.Event) {
	firstName := metadata.GetFromContext(event.Message().Context, metadata.FirstNameKey)
	if firstName != "" {
		decoded, err := b64.StdEncoding.DecodeString(firstName)
		if err != nil {
			return
		}
		(*meta)[metadata.FirstNameKey] = string(decoded)
	}
}

func findLastName(meta *metadata.Metadata, event broker.Event) {
	lastName := metadata.GetFromContext(event.Message().Context, metadata.LastNameKey)
	if lastName != "" {
		decoded, err := b64.StdEncoding.DecodeString(lastName)
		if err != nil {
			return
		}
		(*meta)[metadata.LastNameKey] = string(decoded)
	}
}

func findType(meta *metadata.Metadata, event broker.Event) {
	typeValue := metadata.GetFromContext(event.Message().Context, metadata.TypeKey)
	if typeValue != "" {
		(*meta)[metadata.TypeKey] = typeValue
	}
}

func findSubType(meta *metadata.Metadata, event broker.Event) {
	subType := metadata.GetFromContext(event.Message().Context, metadata.SubTypeKey)
	if subType != "" {
		(*meta)[metadata.SubTypeKey] = subType
	}
}

func findRoles(meta *metadata.Metadata, event broker.Event) {
	roles := metadata.GetFromContext(event.Message().Context, metadata.RolesKey)
	if roles != "" {
		(*meta)[metadata.RolesKey] = roles
	}
}

func findInternal(meta *metadata.Metadata, event broker.Event) {
	internal := metadata.GetFromContext(event.Message().Context, metadata.InternalKey)
	if internal != "" {
		(*meta)[metadata.InternalKey] = internal
	}
}
