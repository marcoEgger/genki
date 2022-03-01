package config

import (
	"time"

	"github.com/spf13/pflag"
)

var cfg Config

func init() {
	cfg = NewConfig()
}

func NewConfig() Config {
	return newViperConfig()
}

type Config interface {
	BindFlags(set *pflag.FlagSet)
	GetString(key string) string
	GetBool(key string) bool
	GetDuration(key string) time.Duration
}

func BindFlagSet(set *pflag.FlagSet) {
	cfg.BindFlags(set)
}

func GetString(key string) string {
	return cfg.GetString(key)
}

func GetBool(key string) bool {
	return cfg.GetBool(key)
}

func GetDuration(key string) time.Duration {
	return cfg.GetDuration(key)
}
