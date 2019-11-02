package logger

import (
	"context"

	"github.com/lukasjarosch/genki/server/grpc/metadata"
	"github.com/opentracing/opentracing-go"
	zipkintracer "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	sugared *zap.SugaredLogger
}

func newZapLogger(level string) (Logger, error) {
	logLevel := parseStringLevel(level)

	zapEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "severity",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	zapConfig := zap.Config{
		Level:       logLevel,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zapEncoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}

	logger, err := zapConfig.Build()
	if err != nil {
	    return nil, err
	}

	// skip the two newest calls in the call-stack, they are from this logging package and of no use.
	logger = logger.WithOptions(zap.AddCallerSkip(2))


	return &zapLogger{
		sugared: logger.Sugar(),
	}, nil
}

func (l *zapLogger) Debug(fields ...interface{})  {
	l.sugared.Debug(fields...)
}

func (l *zapLogger) Debugf(format string, fields ...interface{})  {
	l.sugared.Debugf(format, fields...)
}

func (l *zapLogger) Info(fields ...interface{})  {
	l.sugared.Info(fields...)
}

func (l *zapLogger) Infof(format string, fields ...interface{})  {
	l.sugared.Infof(format, fields...)
}

func (l *zapLogger) Warn(fields ...interface{})  {
	l.sugared.Info(fields...)
}

func (l *zapLogger) Warnf(format string, fields ...interface{})  {
	l.sugared.Warnf(format, fields...)
}

func (l *zapLogger) Error(fields ...interface{})  {
	l.sugared.Info(fields...)
}

func (l *zapLogger) Errorf(format string, fields ...interface{})  {
	l.sugared.Errorf(format, fields...)
}

func (l *zapLogger) Fatal(fields ...interface{})  {
	l.sugared.Info(fields...)
}

func (l *zapLogger) Fatalf(format string, fields ...interface{})  {
	l.sugared.Fatalf(format, fields...)
}

func (l *zapLogger) WithFields(keyValues Fields) Logger {
	var f = make([]interface{}, 0)
	for k, v := range keyValues {
		f = append(f, k)
		f = append(f, v)
	}

	newLogger := l.sugared.With(f...)
	return &zapLogger{newLogger}
}

func (l *zapLogger) WithContext(ctx context.Context) Logger {
	var logger Logger
	span, ok := opentracing.SpanFromContext(ctx).Context().(zipkintracer.SpanContext)
	if ok {
		logger = log.WithFields(Fields{
			metadata.TraceID: span.TraceID.String(),
		})
	}

	return logger.WithFields(
		Fields{
			metadata.RequestID: metadata.GetRequestID(ctx),
		},
	)
}

func parseStringLevel(logLevel string) zap.AtomicLevel{
	var level zap.AtomicLevel

	switch logLevel {
	case DebugLevel:
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case InfoLevel:
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case WarnLevel:
		level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case ErrorLevel:
		level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case FatalLevel:
		level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	}

	return level
}