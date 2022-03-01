package service

import (
	"os"
	"os/signal"
	"syscall"
)

var onlyOneSignalHandler = make(chan struct{})
var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

// NewSignalHandler registers a SIGTERM and SIGINT handler.
// A stop channel is returned, which is closed when one of these signals are caught.
// If a second signal is caught, the application is terminated immediately with exit code 1.
func NewSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler)

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)

	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal: terminate immediately
	}()

	return stop
}
