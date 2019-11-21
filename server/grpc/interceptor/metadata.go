package interceptor

import (
	"context"
	"encoding/base64"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	md "github.com/lukasjarosch/genki/metadata"
)

const RequestIdMetadataKey = "requestId"
const AccountIdMetadataKey = "accountId"
const UserIdMetadataKey = "userId"
const EmailMetadataKey = "email"
const FirstNameMetadataKey = "firstName"
const LastNameMetadataKey = "lastName"
const TypeMetadataKey = "type"
const SubTypeMetadataKey = "subType"
const RolesMetadataKey = "roles"

func UnaryServerMetadata() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		meta := md.Metadata{}
		ensureRequestId(ctx, &meta)
		findAccountId(ctx, &meta)
		findUserId(ctx, &meta)
		findEmail(ctx, &meta)
		findFirstName(ctx, &meta)
		findLastName(ctx, &meta)
		findType(ctx, &meta)
		findSubType(ctx, &meta)
		findRoles(ctx, &meta)
		ctx = md.NewContext(ctx, meta)

		return handler(ctx, req)
	}
}

func UnaryClientMetadata() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		meta := md.Metadata{}
		ensureRequestId(ctx, &meta)
		findAccountId(ctx, &meta)
		findUserId(ctx, &meta)
		findEmail(ctx, &meta)
		findFirstName(ctx, &meta)
		findLastName(ctx, &meta)
		findType(ctx, &meta)
		findSubType(ctx, &meta)
		findRoles(ctx, &meta)
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
		reqId := md.NewRequestID()
		(*meta)[md.RequestIDKey] = reqId
		ctx = metadata.AppendToOutgoingContext(ctx, md.RequestIDKey, reqId)
	}
}

func findAccountId(ctx context.Context, meta *md.Metadata) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		accountID := header.Get(AccountIdMetadataKey)
		if len(accountID) > 0 {
			(*meta)[md.AccountIDKey] = accountID[0]
			ctx = metadata.AppendToOutgoingContext(ctx, md.RequestIDKey, accountID[0])
		}
	}
}

func findUserId(ctx context.Context, meta *md.Metadata) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		userID := header.Get(UserIdMetadataKey)
		if len(userID) > 0 {
			(*meta)[md.UserIDKey] = userID[0]
			ctx = metadata.AppendToOutgoingContext(ctx, md.RequestIDKey, userID[0])
		}
	}
}

func findEmail(ctx context.Context, meta *md.Metadata) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		email := header.Get(EmailMetadataKey)
		if len(email) > 0 {
			emailEncoded := base64.StdEncoding.EncodeToString([]byte(email[0]))
			(*meta)[md.EmailKey] = emailEncoded
			ctx = metadata.AppendToOutgoingContext(ctx, md.RequestIDKey, emailEncoded)
		}
	}
}

func findFirstName(ctx context.Context, meta *md.Metadata) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		firstName := header.Get(FirstNameMetadataKey)
		if len(firstName) > 0 {
			firstNameEncoded := base64.StdEncoding.EncodeToString([]byte(firstName[0]))
			(*meta)[md.FirstNameKey] = firstNameEncoded
			ctx = metadata.AppendToOutgoingContext(ctx, md.RequestIDKey, firstNameEncoded)
		}
	}
}

func findLastName(ctx context.Context, meta *md.Metadata) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		lastName := header.Get(LastNameMetadataKey)
		if len(lastName) > 0 {
			lastNameEncoded := base64.StdEncoding.EncodeToString([]byte(lastName[0]))
			(*meta)[md.LastNameKey] = lastNameEncoded
			ctx = metadata.AppendToOutgoingContext(ctx, md.RequestIDKey, lastNameEncoded)
		}
	}
}

func findType(ctx context.Context, meta *md.Metadata) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		typ := header.Get(TypeMetadataKey)
		if len(typ) > 0 {
			(*meta)[md.TypeKey] = typ[0]
			ctx = metadata.AppendToOutgoingContext(ctx, md.RequestIDKey, typ[0])
		}
	}
}

func findSubType(ctx context.Context, meta *md.Metadata) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		subType := header.Get(SubTypeMetadataKey)
		if len(subType) > 0 {
			(*meta)[md.SubTypeKey] = subType[0]
			ctx = metadata.AppendToOutgoingContext(ctx, md.RequestIDKey, subType[0])
		}
	}
}

func findRoles(ctx context.Context, meta *md.Metadata) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		roles := header.Get(RolesMetadataKey)
		if len(roles) > 0 {
			(*meta)[md.SubTypeKey] = roles[0]
			ctx = metadata.AppendToOutgoingContext(ctx, md.RequestIDKey, roles[0])
		}
	}
}
