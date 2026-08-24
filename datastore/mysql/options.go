package mysql

import (
	"time"

	"github.com/spf13/pflag"
)

const (
	AddressConfigKey       = "mysql-address"
	TLSCAFileConfigKey     = "mysql-tls-ca"
	TLSServerNameConfigKey = "mysql-tls-server-name"
)

type Options struct {
	MigrationPath         string
	MaxOpenConnections    int
	MaxIdleConnections    int
	MaxConnectionLifetime time.Duration
	// DisableTLS skips automatic TLS enablement. Use only for local setups
	// that do not support TLS. Prefer setting tls=false in the DSN instead.
	DisableTLS bool
	// TLSCAFile is an optional path to a PEM-encoded CA certificate used to
	// verify the MySQL server certificate.
	TLSCAFile string
	// TLSServerName overrides hostname verification when the server certificate
	// does not match the DSN host (e.g. connect via IP but verify a DNS SAN).
	TLSServerName string
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

// DisableTLS disables automatic TLS for the MySQL connection.
// Prefer an explicit tls=... DSN parameter when possible.
func DisableTLS() Option {
	return func(options *Options) {
		options.DisableTLS = true
	}
}

// TLSCAFile sets a custom CA certificate file used to verify the MySQL server.
func TLSCAFile(path string) Option {
	return func(options *Options) {
		options.TLSCAFile = path
	}
}

// TLSServerName sets the expected server name for TLS hostname verification.
func TLSServerName(serverName string) Option {
	return func(options *Options) {
		options.TLSServerName = serverName
	}
}

func Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("mysql", pflag.ContinueOnError)
	// tls=false keeps the local default usable; production DSNs without tls=
	// get TLS enabled automatically by New().
	fs.String(AddressConfigKey, "root:root@tcp(localhost:3306)/database?parseTime=true&tls=false", "mysql connection string")
	fs.String(TLSCAFileConfigKey, "", "optional PEM CA certificate file for MySQL TLS verification")
	fs.String(TLSServerNameConfigKey, "", "optional TLS server name override for hostname verification")
	return fs
}

func newOptions(opts ...Option) *Options {
	opt := &Options{
		MigrationPath:         "migrations",
		MaxOpenConnections:    25,
		MaxIdleConnections:    25,
		MaxConnectionLifetime: 5 * time.Minute,
	}
	for _, o := range opts {
		o(opt)
	}

	return opt
}
