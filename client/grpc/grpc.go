package grpc

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	"github.com/lukasjarosch/genki/config"
	"github.com/lukasjarosch/genki/logger"
	"github.com/lukasjarosch/genki/server/grpc/interceptor"
)

type Client struct {
	conn *grpc.ClientConn
	name string
}

func NewClient(name string) Client {
	return Client{
		conn: nil,
		name: name,
	}
}

func (c *Client) Connect() (err error) {
	if c.name == "" {
		return errors.New("missing client name")
	}

	// do nothing if a ready connection is already available
	if c.conn != nil && c.conn.GetState() == connectivity.Ready {
		return nil
	}

	addr := config.GetString(fmt.Sprintf("%s-grpc-client-address", c.name))
	c.conn, err = grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithChainUnaryInterceptor(
			interceptor.UnaryClientMetadata(),
			interceptor.UnaryClientLogging(),
		),
	)
	if err != nil {
	    return errors.Wrap(err, fmt.Sprintf("gRPC client connection '%s' (%s) failed", c.name, addr))
	}
	logger.Infof("gRPC client connection to '%s' (%s) established", c.name, addr)

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
