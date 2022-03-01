package interceptor

import (
	"context"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	serverRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_ms",
			Help:    "Duration of gRPC requests in ms",
			Buckets: []float64{50, 100, 250, 1000},
		},
		[]string{"method", "status_code"})
	serverRequestsCurrent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "grpc_current_request_count",
			Help: "The current amount of requests which are being handled",
	}, []string{"method"})
	clientRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_client_request_duration_ms",
			Help:    "Duration of gRPC client requests in ms",
			Buckets: []float64{50, 100, 250, 1000},
		},
		[]string{"method", "status_code"})
	clientRequestErrorCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        "grpc_client_request_errors",
		Help:        "Amount of errors on gRPC requests",
	}, []string{"endpoint", "method", "error_code"})
)

func init() {
	_ = prometheus.Register(serverRequestDuration)
	_ = prometheus.Register(serverRequestsCurrent)
	_ = prometheus.Register(clientRequestErrorCounter)
	_ = prometheus.Register(clientRequestDuration)
}

// UnaryServerPrometheus adds basic RED metrics on all endpoints. The transport layer (gRPC) should also have metrics attached and
// will then take care of monitoring grpc endpoints including their status.
func UnaryServerPrometheus() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		meth := strings.Split(info.FullMethod, "/")
		rpcName := meth[len(meth)-1]

		serverRequestsCurrent.With(prometheus.Labels{
			"method": rpcName,
		}).Inc()

		defer func(begin time.Time) {
			stat := codes.OK.String()
			if err != nil {
				if s, ok := status.FromError(err); ok {
					stat = s.Code().String()
				}
			}

			duration := time.Since(begin)
			serverRequestDuration.With(prometheus.Labels{
				"method":      rpcName,
				"status_code": stat,
			}).Observe(float64(duration.Milliseconds()))

			serverRequestsCurrent.With(prometheus.Labels{
				"method": rpcName,
			}).Dec()
		}(time.Now())

		return handler(ctx, req)
	}
}

func UnaryClientPrometheus() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		defer func(started time.Time) {
			duration := time.Since(started)
			stat := codes.OK.String()
			if err != nil {
				if s, ok := status.FromError(err); ok {
					stat = s.Code().String()
				} else {
					stat = codes.Unknown.String()
				}
			}
			clientRequestDuration.With(prometheus.Labels{
				"method": method,
				"status_code": stat,
			}).Observe(float64(duration.Milliseconds()))

			if err != nil {
				clientRequestErrorCounter.With(prometheus.Labels{
					"method": method,
					"error_code": stat,
					"endpoint": cc.Target(),
				}).Inc()
			}
		}(time.Now())

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
