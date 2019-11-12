package authz

import "context"

type Authorizer interface {
	Authorize(ctx context.Context, resourceId, action string, externalData interface{}) error
}