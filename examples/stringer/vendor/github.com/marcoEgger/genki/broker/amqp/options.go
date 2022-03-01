package amqp

import (
	"fmt"
	"os"

	"github.com/rs/xid"
	"github.com/spf13/pflag"
)

const AddressConfigKey = "amqp-url"

type Options struct {
	Address         string
	PrefetchCount   int
	SubscriberQueue string
	ConsumerName    string
}

func Address(address string) Option {
	return func(opts *Options) {
		opts.Address = address
	}
}

func PrefetchCount(prefetch int) Option {
	return func(opts *Options) {
		opts.PrefetchCount = prefetch
	}
}

func ConsumerQueue(queue string) Option {
	return func(opts *Options) {
		opts.SubscriberQueue = queue
	}
}

func ConsumerName(name string) Option {
	return func(opts *Options) {
		opts.ConsumerName = name
	}
}

func newOptions(opts ...Option) *Options {
	defaults := &Options{
		Address:         "",
		SubscriberQueue: "default-queue",
		ConsumerName:    defaultConsumerName(),
		PrefetchCount:   10,
	}

	for _, o := range opts {
		o(defaults)
	}

	return defaults
}

func Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("amqp", pflag.ContinueOnError)
	fs.String(AddressConfigKey, "amqp://guest:guest@localhost:5672/", "the amqp connection string: amqp://<user>:<pass>@host:port/<vhost>")
	return fs
}

func defaultConsumerName() string {
	hostname, _ := os.Hostname()
	binaryName := os.Args[0]
	random := xid.New().String()

	return fmt.Sprintf("%s_%s-%s", binaryName, hostname, random)
}

type Option func(*Options)
