package mongodb

import (
	"github.com/spf13/pflag"
)

const AddressConfigKey = "mongodb-address"

func Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("mongodb", pflag.ContinueOnError)
	fs.String(AddressConfigKey, "mongodb://root:root@localhost:27017/database?authSource=admin", "mongodb connection string")
	return fs
}
