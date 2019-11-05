package logger

import (
	"context"

	"github.com/spf13/pflag"

	"github.com/lukasjarosch/genki/config"
)

// Global log instance to be able to directly access the log functions
var log Logger

const DefaultLevel = InfoLevel

type Logger interface {
	Debug(fields ...interface{})
	Debugf(format string, args ...interface{})
	Info(fields ...interface{})
	Infof(format string, args ...interface{})
	Warn(fields ...interface{})
	Warnf(format string, args ...interface{})
	Error(fields ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(fields ...interface{})
	Fatalf(format string, args ...interface{})
	WithFields(keyValues Fields) Logger
	WithMetadata(ctx context.Context) Logger
}

type Fields map[string]interface{}

const (
	DebugLevel        = "debug"
	InfoLevel         = "info"
	WarnLevel         = "warn"
	ErrorLevel        = "error"
	FatalLevel        = "fatal"
	DefaultCallerSkip = 2
	LogLevelConfigKey = "log-level"
)

func NewLogger(level string) error {
	logger, err := newZapLogger(level, DefaultCallerSkip)
	if err != nil {
		return err
	}
	log = logger
	return nil
}

func Debug(fields ...interface{}) {
	log.Debug(fields...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Info(fields ...interface{}) {
	log.Info(fields...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warn(fields ...interface{}) {
	log.Warn(fields...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Error(fields ...interface{}) {
	log.Error(fields...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatal(fields ...interface{}) {
	log.Fatal(fields...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func WithFields(keyValues Fields) Logger {
	return log.WithFields(keyValues)
}

func WithMetadata(ctx context.Context) Logger {
	return log.WithMetadata(ctx)
}

// Flags is a convenience function to quickly add the log options as CLI flags
// Implements the cli.FlagProvider type
func Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("logger", pflag.ContinueOnError)

	fs.String(
		LogLevelConfigKey,
		DefaultLevel,
		"log level defines the lowest level of logs printed (debug, info, warn, error, fatal)",
	)

	return fs
}

// EnsureLoggerFromConfig is a convenience function to quickly create a new logger from
// the configuration.
// This requires that the configuration has already be bound from the flags.
// The function will FATAL if the logger could not be created.
func EnsureLoggerFromConfig() {
	if err := NewLogger(config.GetString(LogLevelConfigKey)); err != nil {
		log.Fatal(err.Error())
	}
}
