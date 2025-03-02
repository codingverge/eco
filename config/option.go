package config

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/codingverge/axon"
	"github.com/codingverge/axon/logrus"
	"github.com/knadh/koanf/v2"
	"github.com/ory/jsonschema/v3"
	"github.com/ory/x/watcherx"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

type Option func(c *Config)

func WithContext(ctx context.Context) Option {
	return func(p *Config) {
		for _, o := range ConfigOptionsFromContext(ctx) {
			o(p)
		}
	}
}

func WithConfigFiles(files ...string) Option {
	return func(p *Config) {
		p.files = append(p.files, files...)
	}
}

func WithImmutables(immutables ...string) Option {
	return func(p *Config) {
		p.immutables = append(p.immutables, immutables...)
	}
}

func WithExceptImmutables(exceptImmutables ...string) Option {
	return func(p *Config) {
		p.exceptImmutables = append(p.exceptImmutables, exceptImmutables...)
	}
}

func WithFlags(flags *pflag.FlagSet) Option {
	return func(p *Config) {
		p.flags = flags
	}
}

func WithLogger(l axon.Logger) Option {
	return func(p *Config) {
		p.logger = l
	}
}

func SkipValidation() Option {
	return func(p *Config) {
		p.skipValidation = true
	}
}

func DisableEnvLoading() Option {
	return func(p *Config) {
		p.disableEnvLoading = true
	}
}

func WithValue(key string, value interface{}) Option {
	return func(p *Config) {
		p.forcedValues = append(p.forcedValues, tuple{Key: key, Value: value})
	}
}

func WithValues(values map[string]interface{}) Option {
	return func(p *Config) {
		for key, value := range values {
			p.forcedValues = append(p.forcedValues, tuple{Key: key, Value: value})
		}
	}
}

func WithBaseValues(values map[string]interface{}) Option {
	return func(p *Config) {
		for key, value := range values {
			p.baseValues = append(p.baseValues, tuple{Key: key, Value: value})
		}
	}
}

func WithUserProviders(providers ...koanf.Provider) Option {
	return func(p *Config) {
		p.userProviders = providers
	}
}

func AttachWatcher(watcher func(event watcherx.Event, err error)) Option {
	return func(p *Config) {
		p.onChanges = append(p.onChanges, watcher)
	}
}

func WithLogrusWatcher(l *logrus.Logger) Option {
	return AttachWatcher(LogrusWatcher(l))
}

func LogrusWatcher(l *logrus.Logger) func(e watcherx.Event, err error) {
	return func(e watcherx.Event, err error) {
		l.WithField("file", e.Source()).
			WithField("event_type", fmt.Sprintf("%T", e)).
			Info("A change to a configuration file was detected.")

		if et := new(jsonschema.ValidationError); errors.As(err, &et) {
			l.WithField("event", fmt.Sprintf("%#v", et)).
				Errorf("The changed configuration is invalid and could not be loaded. Rolling back to the last working configuration revision. Please address the validation errors before restarting the process.")
		} else if et := new(ImmutableError); errors.As(err, &et) {
			l.WithError(err).
				WithField("key", et.Key).
				WithField("old_value", fmt.Sprintf("%v", et.From)).
				WithField("new_value", fmt.Sprintf("%v", et.To)).
				Errorf("A configuration value marked as immutable has changed. Rolling back to the last working configuration revision. To reload the values please restart the process.")
		} else if err != nil {
			l.WithError(err).Errorf("An error occurred while watching config file %s", e.Source())
		} else {
			l.WithField("file", e.Source()).
				WithField("event_type", fmt.Sprintf("%T", e)).
				Info("Configuration change processed successfully.")
		}
	}
}

func WithStderrValidationReporter() Option {
	return func(p *Config) {
		p.onValidationError = func(k *koanf.Koanf, err error) {
			p.printHumanReadableValidationErrors(k, os.Stderr, err)
		}
	}
}

func WithStandardValidationReporter(w io.Writer) Option {
	return func(p *Config) {
		p.onValidationError = func(k *koanf.Koanf, err error) {
			p.printHumanReadableValidationErrors(k, w, err)
		}
	}
}
