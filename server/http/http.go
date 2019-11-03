package http

import (
	"context"
	"net/http"
	"sync"
)

type Server interface {
	ListenAndServe(ctx context.Context, wg *sync.WaitGroup)
	Handle(endpoint string, handler http.Handler)
}

const (
	PortConfigKey        = "http-port"
	GracePeriodConfigKey = "http-grace-period"
)