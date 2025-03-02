package driver

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/codingverge/axon"
	"github.com/codingverge/axon/config"
	"github.com/codingverge/axon/dbal"
	"github.com/codingverge/axon/logrus"
	"github.com/ory/graceful"
	"github.com/urfave/negroni"
	"golang.org/x/sync/errgroup"
)

// registering DefaultDriver to dbal
func init() {
	dbal.RegisterDriver(func() axon.Dbal {
		return NewDefaultDriver()
	})
}

func NewDefaultDriver() *DefaultDriver {
	return &DefaultDriver{}
}

type DefaultDriver struct {
	l axon.Logger
	c axon.DriverConfigure
}

func (r *DefaultDriver) WithLogger(l axon.Logger) axon.Driver {
	r.l = l
	return r
}

func (r *DefaultDriver) WithConfig(c axon.DriverConfigure) axon.Driver {
	r.c = c
	return r
}

var _ axon.Driver = &DefaultDriver{}

func New(ctx context.Context, stdOutOrErr io.Writer, dOpts []axon.DriverOption, cOpts []config.Option) (axon.Driver, error) {
	opts := axon.NewOptions(dOpts)
	l := opts.Logger()
	if l == nil {
		l = logrus.New().WithOutStream(stdOutOrErr)
	}
	c := opts.Config()
	if c == nil {
		var err error
		c, err = NewDriverConfig(ctx, cOpts...)
		if err != nil {
			l.WithError(err).Error("Unable to instantiate configuration.")
			return nil, err
		}
	}
	r, err := dbal.NewDriverFromDbal(ctx, c, l)
	if err != nil {
		l.WithError(err).Error("Unable to instantiate service registry.")
		return nil, err
	}
	return r, nil
}

func (r *DefaultDriver) CanHandle(dsn string) bool {
	return true
}

func (r *DefaultDriver) Ping() error {
	return nil
}

func (r *DefaultDriver) PingContext(ctx context.Context) error {
	return nil
}

func (r *DefaultDriver) Logger() axon.Logger {
	if r.l == nil {
		r.l = logrus.New()
	}
	return r.l
}

func (r *DefaultDriver) RunE(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	r.server(ctx, g)
	return g.Wait()
}

func (r *DefaultDriver) server(ctx context.Context, eg *errgroup.Group) {
	n := negroni.New()
	var handler http.Handler = n
	server := graceful.WithDefaults(&http.Server{
		Handler:           handler,
		TLSConfig:         &tls.Config{MinVersion: tls.VersionTLS12},
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      120 * time.Second,
		IdleTimeout:       600 * time.Second,
	})
	addr := "0.0.0.0:8080"
	eg.Go(func() error {
		r.l.Printf("Starting the httpd on: %s", addr)
		if err := graceful.GracefulContext(ctx, func() error {
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				return err
			}
			return server.Serve(listener)
		}, server.Shutdown); err != nil {
			if !errors.Is(err, context.Canceled) {
				r.l.Errorf("Failed to gracefully shutdown httpd: %s", err)
				return err
			}
		}
		r.l.Printf("httpd was shutdown gracefully")
		return nil
	})
}
