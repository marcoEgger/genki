package grpc

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"

	"github.com/lukasjarosch/genki/config"
	"github.com/lukasjarosch/genki/logger"
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

	c.conn, err = grpc.Dial(config.GetString(fmt.Sprintf("%s-grpc-client-address", c.name)), grpc.WithInsecure())
	if err != nil {
	    return errors.Wrap(err, fmt.Sprintf("grpc client connection '%s' failed", c.name))
	}

	return nil
}

func (c *Client) Disconnect() {
	if err := c.conn.Close(); err != nil {
		logger.Warnf("unable to close %s-client connection: %s", c.name, err)
	}
}

func (c *Client) Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet(fmt.Sprintf("%s-client", c.name), pflag.ContinueOnError)
	fs.String(fmt.Sprintf("%s-grpc-client-address", c.name), "localhost:50052", fmt.Sprintf("the upstream addess to which the %s-client will connect", c.name))
	return fs
}
