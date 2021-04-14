package amqp

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/marcoEgger/genki/logger"
)

// Connection is a wrapper for amqp.Connection but adding reconnection functionality.
type Connection struct {
	addr                  string
	conn                  *amqp.Connection
	connMutex             sync.Mutex
	ctx                   context.Context
	cancel                context.CancelFunc
	connected             bool
	notifyCloseConnection chan *amqp.Error
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
	}

	return c
}

// Connect will dial to the specified AMQP server addr.
func (c *Connection) Connect() (err error) {
	c.conn, err = c.dial()
	if err != nil {
		return errors.Wrap(err, "unable to connect to amqp server")
	}

	go c.monitorConnection()

	return nil
}

// Shutdown the reconnector and terminate any existing connections
func (c *Connection) Shutdown() {
	c.setConnected(false)
	c.cancel()

	if c.IsConnected() {
		err := c.conn.Close()
		if err != nil {
			logger.Warnf("error while closing amqp connection: %s", err)
			return
		}
	}
}

// dial and return the connection and any occurred error
func (c *Connection) dial() (*amqp.Connection, error) {
	c.setConnected(false)

	conn, err := amqp.Dial(c.addr)
	if err != nil {
		return nil, err
	}
	c.changeConnection(conn)
	c.setConnected(true)
	return conn, nil
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
				logger.Warnf("amqp connection error (code %s): %s", amqpErr.Code, amqpErr.Reason)
				c.reconnect()
			}
		}
	}
}

// reconnect will, once started, try to connect to amqp forever
// the method only returns if a connection is established or the ctxReconnect context is cancelled by Shutdown()
func (c *Connection) reconnect() {
	var err error

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}
		c.conn, err = c.dial()
		if err != nil {
			logger.Warnf("unable to connect to amqp server: %s", err)
			time.Sleep(ReconnectDelay)
			continue
		}
		logger.Info("reconnected to amqp server")
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
	c.connMutex.Lock()
	defer c.connMutex.Unlock()
	return c.connected
}

func (c *Connection) setConnected(status bool) {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()
	c.connected = status
}

func (c *Connection) Channel() (*amqp.Channel, error) {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()
	return c.conn.Channel()
}
