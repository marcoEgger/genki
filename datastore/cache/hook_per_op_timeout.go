package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisOperationTimeout is the maximum time to wait on a single Redis command
// (GET/SET/DEL/…). Applied per operation via hook, independent of gRPC latency.
const RedisOperationTimeout = 200 * time.Millisecond

type cancelCtxKey struct{}

// perOpTimeoutHook attaches a fresh deadline to each Redis command so slow cache
// I/O does not share one wall-clock budget with later steps (e.g. cache SET
// after a long upstream call in go-redis/cache Once).
type perOpTimeoutHook struct {
	timeout time.Duration
}

func newPerOpTimeoutHook(d time.Duration) *perOpTimeoutHook {
	return &perOpTimeoutHook{timeout: d}
}

func (h *perOpTimeoutHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	opCtx, cancel := context.WithTimeout(ctx, h.timeout)
	return context.WithValue(opCtx, cancelCtxKey{}, cancel), nil
}

func (h *perOpTimeoutHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if cancel, ok := ctx.Value(cancelCtxKey{}).(context.CancelFunc); ok {
		cancel()
	}
	return nil
}

func (h *perOpTimeoutHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	opCtx, cancel := context.WithTimeout(ctx, h.timeout)
	return context.WithValue(opCtx, cancelCtxKey{}, cancel), nil
}

func (h *perOpTimeoutHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if cancel, ok := ctx.Value(cancelCtxKey{}).(context.CancelFunc); ok {
		cancel()
	}
	return nil
}
