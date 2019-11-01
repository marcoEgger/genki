package logger

// Global log instance to be able to directly access the log functions
var log Logger

// ensures that the global logger is never null
func init() {
	_ = NewLogger(InfoLevel)
}

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
}

type Fields map[string]interface{}

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
)

func NewLogger(level string) error {
	logger, err := newZapLogger(level)
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
