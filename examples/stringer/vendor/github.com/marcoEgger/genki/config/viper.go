package config

import (
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type config struct {
}

func newViperConfig() *config {
	return &config{}
}

func (c config) BindFlags(set *pflag.FlagSet) {
	err := viper.BindPFlags(set)
	if err != nil {
	    return
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func (c config) GetString(key string) string {
	return viper.GetString(key)
}

func (c config) GetBool(key string) bool {
	return viper.GetBool(key)
}

func (c config) GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}
