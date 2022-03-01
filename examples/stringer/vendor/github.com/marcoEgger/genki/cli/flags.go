package cli

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

// FlagProvider defines a function which returns a FlagSet.
// The FlagProvider can be used to add new FlagSets in bulk.
type FlagProvider func() *pflag.FlagSet

type FlagSet struct {
	fs *pflag.FlagSet
}

func NewFlagSet(name string) *FlagSet {
	return &FlagSet{
		fs: pflag.NewFlagSet(name, pflag.ContinueOnError),
	}
}

func (fs *FlagSet) Add(providers ...FlagProvider) {
	for _, set := range providers {
		fs.fs.AddFlagSet(set())
	}
}

// Parse the flags passed to the binary.
// If an error occurs, the error is logged and the help displayed.
func (fs *FlagSet) Parse() {
	err := fs.fs.Parse(os.Args[1:])
	switch {
	case err == pflag.ErrHelp:
		os.Exit(0)
	case err != nil:
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n\n", err.Error())
		fs.fs.PrintDefaults()
		os.Exit(2)
	}
}

// Set returns the FlagSet
func (fs *FlagSet) Set() *pflag.FlagSet {
	return fs.fs
}
