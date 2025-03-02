package axon

import "context"

type Driver interface {
	Dbal
	LoggingProvider

	WithLogger(l Logger) Driver
	WithConfig(c DriverConfigure) Driver

	RunE(ctx context.Context) error
}

type Options struct {
	logger Logger
	config DriverConfigure
}

type DriverOption func(*Options)

func WithLogger(l Logger) DriverOption {
	return func(o *Options) {
		o.logger = l
	}
}

func WithConfig(c DriverConfigure) DriverOption {
	return func(o *Options) {
		o.config = c
	}
}

func (o *Options) Logger() Logger {
	return o.logger
}

func (o *Options) Config() DriverConfigure {
	return o.config
}

func NewOptions(os []DriverOption) *Options {
	o := new(Options)
	for _, f := range os {
		f(o)
	}
	return o
}
