package axon

import (
	"context"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type LoggingProvider interface {
	Logger() Logger
}

type Logger interface {
	WithField(key string, value interface{}) Logger
	WithRequest(r *http.Request) Logger
	WithError(err error) Logger
	WithOutStream(w io.Writer) Logger
	UseConfig(configure Configure)
	WithFields(f logrus.Fields) Logger
	WithSensitiveField(key string, value interface{}) Logger
	WithSpanFromContext(ctx context.Context) Logger
	WithContext(ctx context.Context) Logger
	LeakSensitiveData() bool
	Info(v ...interface{})
	Infof(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warn(args ...interface{})
	Warningf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})
	Tracef(format string, args ...interface{})
}
