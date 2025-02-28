package log

import (
	"io"
	"net/http"
)

type LoggerProvider interface {
	Logger() Logger
}

type Logger interface {
	WithField(key string, value interface{}) Logger
	WithRequest(r *http.Request) Logger
	WithError(err error) Logger
	WithOutStream(w io.Writer)

	Info(v ...interface{})
	Infof(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Error(args ...interface{})
	Printf(format string, args ...interface{})
}
