package amqp

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/lukasjarosch/genki/logger"
)

// Connection is a wrapper for amqp.Connection but adding reconnection functionality.
type Connection struct {
	addr                  string
	conn                  *amqp.Connection
	connMutex             sync.Mutex
	channel               *amqp.Channel
	ctx                   context.Context
	cancel                context.CancelFunc
	connected             bool
	notifyCloseConnection chan *amqp.Error
	name                  string
}

const ReconnectDelay = 5 * time.Second

func NewConnection(addr string) *Connection {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Connection{
		ctx:                   ctx,
		cancel:                cancel,
		addr:                  addr,
		connMutex:             sync.Mutex{},
		notifyCloseConnection: make(chan *amqp.Error),
		connected:             false,
		name:                  "default",
	}

	return c
}

func (c *Connection) SetName(name string) {
	c.name = name
}

// Consume will dial to the specified AMQP server addr.
func (c *Connection) Connect() (err error) {
	err = c.dial()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("%s: unable to connect to amqp server", c.name))
	}

	go c.monitorConnection()

	return nil
}

func (c *Connection) WaitForConnection() {
	t := time.Tick(200 * time.Millisecond)
	for {
		select {
		case <-t:
			if c.IsConnected() {
				c.setConnected(true)
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// Shutdown the reconnector and terminate any existing connections
func (c *Connection) Shutdown() {
	c.setConnected(false)
	c.cancel()
	if c.channel != nil {
		_ = c.channel.Close()
	}

	if c.IsConnected() {
		err := c.conn.Close()
		if err != nil {
			logger.Warnf("%s: error while closing amqp connection: %s", c.name, err)
			return
		}
	}
}

// dial and return the connection and any occurred error
func (c *Connection) dial() error {
	c.setConnected(false)
	conn, err := amqp.Dial(c.addr)
	if err != nil {
		return err
	}
	c.changeConnection(conn)
	c.setConnected(true)
	return nil
}

// monitorConnection ensures that the amqp connection is recovered on failures.
// if an error can be read from the amqp connectionClosed channel, then reconnect() is called
func (c *Connection) monitorConnection() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case amqpErr, ok := <-c.notifyCloseConnection:
			if ok {
				logger.Warnf("%s: amqp connection closed, attempting reconnect (code %d): %s", c.name, amqpErr.Code, amqpErr.Reason)
				c.reconnect()
			}
		}
	}
}

// reconnect will, once started, try to connect to amqp forever
// the method only returns if a connection is established or the ctxReconnect context is cancelled by Shutdown()
func (c *Connection) reconnect() {
	c.setConnected(false)

	// close the current channel and set to nil. the channel is only initialized when needed by Channel()
	if c.channel != nil {
		_ = c.channel.Close()
		c.channel = nil
	}

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}
		err := c.dial()
		if err != nil {
			logger.Warnf("%s: unable to connect to amqp server: %s", c.name, err)
			c.setConnected(false)
			time.Sleep(ReconnectDelay)
			continue
		}
		logger.Infof("%s: reconnected to amqp server", c.name)
		c.setConnected(true)
		return
	}
}

// changeConnection sets a new amqp.Connection and renews the notification channel
func (c *Connection) changeConnection(connection *amqp.Connection) {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()

	c.conn = connection
	c.notifyCloseConnection = make(chan *amqp.Error)
	c.conn.NotifyClose(c.notifyCloseConnection)
}

func (c *Connection) IsConnected() bool {
	return c.connected
}

func (c *Connection) setConnected(status bool) {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()

	c.connected = status
}

func (c *Connection) Channel() (channel *amqp.Channel, err error) {
	if !c.IsConnected() {
		return nil, fmt.Errorf("amqp connection not alive")
	}

	if c.channel == nil {
		c.channel, err = c.conn.Channel()
		if err != nil {
			return nil, err
		}
	}
	return c.channel, nil
}
