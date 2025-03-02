package driver

import (
	"context"
	"fmt"
	"io"

	"github.com/codingverge/axon"
	"github.com/codingverge/axon/dbal"
	"github.com/codingverge/axon/driver/config"
	"github.com/codingverge/axon/logrus"
)

// registering DefaultDriver to dbal
//func init() {
//	dbal.RegisterDriver(func() axon.Dbal {
//		return NewDefaultDriver()
//	})
//}

func NewDefaultDriver() *DefaultDriver {
	return &DefaultDriver{}
}

type DefaultDriver struct {
	l axon.Logger
	c axon.DriverConfigure
}

func (d *DefaultDriver) WithLogger(l axon.Logger) axon.Driver {
	d.l = l
	return d
}

func (d *DefaultDriver) WithConfig(c axon.DriverConfigure) axon.Driver {
	d.c = c
	return d
}

var _ axon.Driver = &DefaultDriver{}

func New(ctx context.Context, stdOutOrErr io.Writer, dOpts ...axon.DriverOption) (axon.Driver, error) {
	opts := axon.NewOptions(dOpts)
	l := opts.Logger()
	if l == nil {
		l = logrus.New().WithOutStream(stdOutOrErr)
	}
	c := opts.Config()
	if c == nil {
		var err error
		c, err = config.New(ctx)
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

func (d *DefaultDriver) CanHandle(dsn string) bool {
	return true
}

func (d *DefaultDriver) Ping() error {
	return nil
}

func (d *DefaultDriver) PingContext(ctx context.Context) error {
	return nil
}

func (d *DefaultDriver) Logger() axon.Logger {
	if d.l == nil {
		d.l = logrus.New()
	}
	return d.l
}

func (d *DefaultDriver) RunE() error {
	fmt.Println("starting server")
	return nil
}
