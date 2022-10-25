package grpc

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"sync"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"

	"github.com/marcoEgger/genki/logger"
	"github.com/marcoEgger/genki/server/grpc/interceptor"
)

type server struct {
	opts Options
	grpc *grpc.Server

	// only set if the gRPC health server is enabled
	healthz *health.Server
}

type BasicHealthChecker struct {
	opts Options
}

type MySQLHealthChecker struct {
	db   *sqlx.DB
	opts Options
}

type MongodbHealthChecker struct {
	db   *mongo.Client
	opts Options
}

func (s *BasicHealthChecker) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	if req.GetService() != "" && req.GetService() != s.opts.Name {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_SERVICE_UNKNOWN,
		}, nil
	}

	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *BasicHealthChecker) Watch(req *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	return server.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}

func (s *MySQLHealthChecker) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	if req.GetService() != "" && req.GetService() != s.opts.Name {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_SERVICE_UNKNOWN,
		}, nil
	}

	if err := s.db.Ping(); err == nil {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_SERVING,
		}, nil
	} else {
		logger.Errorf("not serving: %s", err)
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}
}

func (s *MySQLHealthChecker) Watch(req *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	return server.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}

func (s *MongodbHealthChecker) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	if req.GetService() != "" && req.GetService() != s.opts.Name {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_SERVICE_UNKNOWN,
		}, nil
	}

	if err := s.db.Ping(ctx, readpref.Primary()); err == nil {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_SERVING,
		}, nil
	} else {
		logger.Errorf("not serving: %s", err)
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}
}

func (s *MongodbHealthChecker) Watch(req *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	return server.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}

func NewBasicHealthChecker(opts ...Option) *BasicHealthChecker {
	return &BasicHealthChecker{
		opts: newOptions(opts...),
	}
}

func NewMySQLHealthChecker(db *sqlx.DB, opts ...Option) *MySQLHealthChecker {
	return &MySQLHealthChecker{
		db:   db,
		opts: newOptions(opts...),
	}
}

func NewMongoDBHealthChecker(db *mongo.Client, opts ...Option) *MongodbHealthChecker {
	return &MongodbHealthChecker{
		db:   db,
		opts: newOptions(opts...),
	}
}

func NewServer(opts ...Option) Server {
	options := newOptions(opts...)
	srv := &server{opts: options}

	var unaryInterceptors []grpc.UnaryServerInterceptor

	unaryInterceptors = append(unaryInterceptors, otelgrpc.UnaryServerInterceptor())

	unaryInterceptors = append(unaryInterceptors, interceptor.UnaryServerMetadata())
	logger.Debugf("gRPC server '%s': metadata interceptor enabled", srv.opts.Name)

	if srv.opts.enabledUnaryInterceptor.logging {
		logger.Debugf("gRPC server '%s': logging interceptor enabled", srv.opts.Name)
		unaryInterceptors = append(unaryInterceptors, interceptor.UnaryServerLogging())
	}
	if srv.opts.enabledUnaryInterceptor.prometheus {
		logger.Debugf("gRPC server '%s': prometheus interceptor enabled", srv.opts.Name)
		unaryInterceptors = append(unaryInterceptors, interceptor.UnaryServerPrometheus())
	}

	srv.grpc = grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		unaryInterceptors...,
	)))

	return srv
}

// ListenAndServe ties everything together and runs the gRPC server in a separate goroutine.
// The method then blocks until the passed context is cancelled, so this method should also be started
// as goroutine if more work is needed after starting the gRPC server.
func (srv *server) ListenAndServe(ctx context.Context, wg *sync.WaitGroup, healthServer grpc_health_v1.HealthServer) {
	defer wg.Done()

	if srv.opts.HealthServerEnabled {
		srv.registerHealthServer(healthServer)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", srv.opts.Port))
	if err != nil {
		srv.setServingStatus(HealthNotServing)
		logger.Fatalf("gRPC server '%s' could not be started: %s", srv.opts.Name, err.Error())
	}

	go func() {
		logger.Infof("gRPC server '%s' running on port %s", srv.opts.Name, srv.opts.Port)
		srv.setServingStatus(HealthServing)
		if err := srv.grpc.Serve(listener); err != nil {
			srv.setServingStatus(HealthNotServing)
			logger.Errorf("gRPC server '%s' crashed: %s", srv.opts.Name, err.Error())
			return
		}
	}()

	<-ctx.Done()
	srv.setServingStatus(HealthNotServing)

	srv.shutdown()
}

// Server returns the raw grpc Server instance. It is required to register services.
func (srv *server) Server() *grpc.Server {
	return srv.grpc
}

// shutdown is responsible of gracefully shutting down the gRPC server
// First, GracefulStop() is executed. If the call doesn't return
// until the ShutdownGracePeriod is exceeded, the server is terminated directly.
func (srv *server) shutdown() {
	stopped := make(chan struct{})
	go func() {
		srv.grpc.GracefulStop()
		close(stopped)
	}()
	t := time.NewTicker(srv.opts.ShutdownGracePeriod)
	select {
	case <-t.C:
		logger.Warnf("gRPC server '%s' graceful shutdown timed-out", srv.opts.Name)
	case <-stopped:
		logger.Infof("gRPC server '%s' shut-down gracefully", srv.opts.Name)
		t.Stop()
	}
}

type Option func(*Options)
