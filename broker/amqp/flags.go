package amqp

import "github.com/spf13/pflag"

const UrlConfigKey = "amqp-url"

func Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("amqp", pflag.ContinueOnError)
	fs.String(UrlConfigKey, "amqp://guest:guest@localhost:5672/", "the amqp connection string: amqp://<user>:<pass>@host:port/<vhost>")
	return fs
}