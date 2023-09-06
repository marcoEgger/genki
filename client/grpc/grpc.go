package grpc

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/marcoEgger/genki/config"
	"github.com/marcoEgger/genki/logger"
	"github.com/marcoEgger/genki/server/grpc/interceptor"
)

type Client struct {
	conn *grpc.ClientConn
	name string
	addr string
}

//goland:noinspection GoUnusedExportedFunction
func NewClient(name string) *Client {
	return &Client{
		conn: nil,
		name: name,
	}
}

//goland:noinspection GoUnusedExportedFunction
func NewClientWithAddress(name, address string) *Client {
	return &Client{
		conn: nil,
		name: name,
		addr: address,
	}
}

func (c *Client) Connect() (err error) {
	if c.name == "" {
		return errors.New("missing client name")
	}

	// do nothing if a ready connection is already available
	if c.conn != nil {
		if c.conn.GetState() == connectivity.Ready {
			return nil
		}
	}

	if c.addr == "" {
		c.addr = config.GetString(fmt.Sprintf("%s-grpc-client-address", c.name))
	}

	if c.addr == "" {
		return fmt.Errorf("missing address for client '%s'", c.name)
	}
	c.conn, err = grpc.Dial(
		c.addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			interceptor.UnaryClientPrometheus(),
			interceptor.UnaryClientMetadata(),
			interceptor.UnaryClientLogging(),
			otelgrpc.UnaryClientInterceptor(),
		),
	)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("gRPC client connection '%s' (%s) failed", c.name, c.addr))
	}

	return nil
}

func (c *Client) Connection() *grpc.ClientConn {
	return c.conn
}

func (c *Client) Disconnect() {
	if err := c.conn.Close(); err != nil {
		logger.Warnf("unable to close %s-client connection: %s", c.name, err)
	}
}

func (c *Client) Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet(fmt.Sprintf("%s-client", c.name), pflag.ContinueOnError)
	fs.String(fmt.Sprintf("%s-grpc-client-address", c.name), "localhost:50051", fmt.Sprintf("the upstream addess to which the %s-client will connect", c.name))
	return fs
}
