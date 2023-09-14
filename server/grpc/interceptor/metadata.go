package interceptor

import (
	"context"
	"encoding/base64"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	md "github.com/marcoEgger/genki/metadata"
)

const RequestIdMetadataKey = "requestId"
const M2MMetadataKey = "m2m"
const AccountIdsMetadataKey = "accountIds"
const AccountIdMetadataKey = "accountId"
const UserIdMetadataKey = "userId"
const EmailMetadataKey = "email"
const FirstNameMetadataKey = "firstName"
const LastNameMetadataKey = "lastName"
const TypeMetadataKey = "type"
const SubTypeMetadataKey = "subType"
const RolesMetadataKey = "roles"
const InternalMetadataKey = "internal"

func UnaryServerMetadata() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		meta := md.Metadata{}
		ensureRequestId(ctx, &meta)
		findM2M(ctx, &meta)
		findAccountIds(ctx, &meta)
		findAccountId(ctx, &meta)
		findUserId(ctx, &meta)
		findEmail(ctx, &meta, false)
		findFirstName(ctx, &meta, false)
		findLastName(ctx, &meta, false)
		findType(ctx, &meta)
		findSubType(ctx, &meta)
		findRoles(ctx, &meta)
		findInternal(ctx, &meta)
		ctx = md.NewContext(ctx, meta)

		return handler(ctx, req)
	}
}

func UnaryClientMetadata() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		meta := md.Metadata{}
		ensureRequestId(ctx, &meta)
		findM2M(ctx, &meta)
		findAccountIds(ctx, &meta)
		findAccountId(ctx, &meta)
		findUserId(ctx, &meta)
		findEmail(ctx, &meta, true)
		findFirstName(ctx, &meta, true)
		findLastName(ctx, &meta, true)
		findType(ctx, &meta)
		findSubType(ctx, &meta)
		findRoles(ctx, &meta)
		findInternal(ctx, &meta)
		ctx = md.NewOutgoingContext(ctx)
		ctx = md.NewContext(ctx, meta)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func ensureRequestId(ctx context.Context, meta *md.Metadata) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		requestID := header.Get(RequestIdMetadataKey)
		if len(requestID) > 0 {
			(*meta)[md.RequestIDKey] = requestID[0]
			ctx = metadata.AppendToOutgoingContext(ctx, md.RequestIDKey, requestID[0])
			return
		}
	}

	// set existing genki metadata
	if err := setAppContextValue(ctx, meta, md.RequestIDKey); err == nil {
		return
	}

	reqId := md.NewRequestID()
	(*meta)[md.RequestIDKey] = reqId
	ctx = metadata.AppendToOutgoingContext(ctx, md.RequestIDKey, reqId)
}

func findM2M(ctx context.Context, meta *md.Metadata) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		m2m := header.Get(M2MMetadataKey)
		if len(m2m) > 0 {
			(*meta)[md.M2MKey] = m2m[0]
			ctx = metadata.AppendToOutgoingContext(ctx, md.M2MKey, m2m[0])
		}
	}
	// eventually the app context is filled with metadata
	_ = setAppContextValue(ctx, meta, md.M2MKey)
}

func findAccountIds(ctx context.Context, meta *md.Metadata) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		accountIds := header.Get(AccountIdsMetadataKey)
		if len(accountIds) > 0 {
			(*meta)[md.AccountIDsKey] = accountIds[0]
			ctx = metadata.AppendToOutgoingContext(ctx, md.AccountIDsKey, accountIds[0])
		}
	}
	// eventually the app context is filled with metadata
	_ = setAppContextValue(ctx, meta, md.AccountIDsKey)
}

func findAccountId(ctx context.Context, meta *md.Metadata) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		accountID := header.Get(AccountIdMetadataKey)
		if len(accountID) > 0 {
			(*meta)[md.AccountIDKey] = accountID[0]
			ctx = metadata.AppendToOutgoingContext(ctx, md.AccountIDKey, accountID[0])
		}
	}

	// eventually the app context is filled with metadata
	_ = setAppContextValue(ctx, meta, md.AccountIDKey)
}

func findUserId(ctx context.Context, meta *md.Metadata) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		userID := header.Get(UserIdMetadataKey)
		if len(userID) > 0 {
			(*meta)[md.UserIDKey] = userID[0]
			ctx = metadata.AppendToOutgoingContext(ctx, md.UserIDKey, userID[0])
		}
	}

	// eventually the app context is filled with metadata
	_ = setAppContextValue(ctx, meta, md.UserIDKey)
}

func findEmail(ctx context.Context, meta *md.Metadata, encode bool) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		email := header.Get(EmailMetadataKey)
		if len(email) > 0 {
			if encode {
				emailEncoded := base64.StdEncoding.EncodeToString([]byte(email[0]))
				(*meta)[md.EmailKey] = emailEncoded
				ctx = metadata.AppendToOutgoingContext(ctx, md.EmailKey, emailEncoded)
			} else {
				(*meta)[md.EmailKey] = email[0]
				ctx = metadata.AppendToOutgoingContext(ctx, md.EmailKey, email[0])
			}
		}
	}

	// eventually the app context is filled with metadata
	_ = setAppContextValue(ctx, meta, md.EmailKey)
}

func findFirstName(ctx context.Context, meta *md.Metadata, encode bool) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		firstName := header.Get(FirstNameMetadataKey)
		if len(firstName) > 0 {
			if encode {
				firstNameEncoded := base64.StdEncoding.EncodeToString([]byte(firstName[0]))
				(*meta)[md.FirstNameKey] = firstNameEncoded
				ctx = metadata.AppendToOutgoingContext(ctx, md.FirstNameKey, firstNameEncoded)
			} else {
				(*meta)[md.FirstNameKey] = firstName[0]
				ctx = metadata.AppendToOutgoingContext(ctx, md.FirstNameKey, firstName[0])
			}
		}
	}
	// eventually the app context is filled with metadata
	_ = setAppContextValue(ctx, meta, md.FirstNameKey)
}

func findLastName(ctx context.Context, meta *md.Metadata, encode bool) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		lastName := header.Get(LastNameMetadataKey)
		if len(lastName) > 0 {
			if encode {
				lastNameEncoded := base64.StdEncoding.EncodeToString([]byte(lastName[0]))
				(*meta)[md.LastNameKey] = lastNameEncoded
				ctx = metadata.AppendToOutgoingContext(ctx, md.LastNameKey, lastNameEncoded)
			} else {
				(*meta)[md.LastNameKey] = lastName[0]
				ctx = metadata.AppendToOutgoingContext(ctx, md.LastNameKey, lastName[0])
			}
		}
	}
	// eventually the app context is filled with metadata
	_ = setAppContextValue(ctx, meta, md.LastNameKey)
}

func findType(ctx context.Context, meta *md.Metadata) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		typ := header.Get(TypeMetadataKey)
		if len(typ) > 0 {
			(*meta)[md.TypeKey] = typ[0]
			ctx = metadata.AppendToOutgoingContext(ctx, md.TypeKey, typ[0])
		}
	}
	// eventually the app context is filled with metadata
	_ = setAppContextValue(ctx, meta, md.TypeKey)
}

func findSubType(ctx context.Context, meta *md.Metadata) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		subType := header.Get(SubTypeMetadataKey)
		if len(subType) > 0 {
			(*meta)[md.SubTypeKey] = subType[0]
			ctx = metadata.AppendToOutgoingContext(ctx, md.SubTypeKey, subType[0])
		}
	}
	// eventually the app context is filled with metadata
	_ = setAppContextValue(ctx, meta, md.SubTypeKey)
}

func findRoles(ctx context.Context, meta *md.Metadata) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		roles := header.Get(RolesMetadataKey)
		if len(roles) > 0 {
			(*meta)[md.RolesKey] = roles[0]
			ctx = metadata.AppendToOutgoingContext(ctx, md.RolesKey, roles[0])
		}
	}
	// eventually the app context is filled with metadata
	_ = setAppContextValue(ctx, meta, md.RolesKey)
}

func findInternal(ctx context.Context, meta *md.Metadata) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		internal := header.Get(InternalMetadataKey)
		if len(internal) > 0 {
			(*meta)[md.InternalKey] = internal[0]
			ctx = metadata.AppendToOutgoingContext(ctx, md.InternalKey, internal[0])
		}
	}
	// eventually the app context is filled with metadata
	_ = setAppContextValue(ctx, meta, md.InternalKey)
}

// setAppContextValue will lookup the given metadata key in the GenkiMetadata struct
// if the value is not empty, it's added to the Metadata as well as the outgoing metadata struct of gRPC
func setAppContextValue(ctx context.Context, meta *md.Metadata, key string) error {
	val := md.GetFromContext(ctx, key)
	if val != "" {
		(*meta)[key] = val
		ctx = metadata.AppendToOutgoingContext(ctx, key, val)
		return nil
	}
	return fmt.Errorf("nothing found")
}
