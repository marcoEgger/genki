package authz

import "context"

type Authorizer interface {
	Authorize(ctx context.Context, resource, action interface{}, externalData interface{}) error
}