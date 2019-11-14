package mysql


import (
	"time"

	"github.com/spf13/pflag"
)

const AddressConfigKey = "mysql-addres"

type Options struct {
	MigrationPath         string
	MaxOpenConnections    int
	MaxIdleConnections    int
	MaxConnectionLifetime time.Duration
}

type Option func(*Options)

func MigrationPath(path string) Option {
	return func(options *Options) {
		options.MigrationPath = path
	}
}

func MaxOpenConnections(connLimit int) Option {
	return func(options *Options) {
		options.MaxOpenConnections = connLimit
	}
}

func MaxIdleConnections(connLimit int) Option {
	return func(options *Options) {
		options.MaxIdleConnections = connLimit
	}
}

func MaxConnectionLifetime(maxLifetime time.Duration) Option {
	return func(options *Options) {
		options.MaxConnectionLifetime = maxLifetime
	}
}

func Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("mysql", pflag.ContinueOnError)
	fs.String(AddressConfigKey, "root:root@tcp(localhost:3306)/database?parseTime=true", "mysql connection string")
	return fs
}

func newOptions(opts ...Option) *Options {
	opt := &Options{
		MigrationPath:         "migrations",
		MaxOpenConnections:    10,
		MaxIdleConnections:    0,
		MaxConnectionLifetime: 600 * time.Second,
	}
	for _, o := range opts {
		o(opt)
	}

	return opt
}

