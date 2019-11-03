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
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_ms",
			Help:    "Duration of gRPC requests in ms",
			Buckets: []float64{50, 100, 250, 1000},
		},
		[]string{"method", "status_code"})
	requestsCurrent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "grpc_current_request_count",
			Help: "The current amount of requests which are being handled",
	}, []string{"method"})
)

func init() {
	_ = prometheus.Register(requestDuration)
	_ = prometheus.Register(requestsCurrent)
}

// Prometheus adds basic RED metrics on all endpoints. The transport layer (gRPC) should also have metrics attached and
// will then take care of monitoring grpc endpoints including their status.
func Prometheus() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		meth := strings.Split(info.FullMethod, "/")
		rpcName := meth[len(meth)-1]

		requestsCurrent.With(prometheus.Labels{
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
			requestDuration.With(prometheus.Labels{
				"method":      rpcName,
				"status_code": stat,
			}).Observe(float64(duration.Milliseconds()))

			requestsCurrent.With(prometheus.Labels{
				"method": rpcName,
			}).Dec()
		}(time.Now())

		return handler(ctx, req)
	}
}
