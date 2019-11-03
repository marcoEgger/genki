package http

import (
	"context"
	"sync"
)

type Server interface {
	ListenAndServe(ctx context.Context, wg *sync.WaitGroup)
}
