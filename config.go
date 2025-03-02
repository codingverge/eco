package axon

import (
	"context"
	"io"
	"time"

	"github.com/inhies/go-bytesize"
)

type (
	Configure interface {
		DirtyPatch(key string, value any) error
		Set(key string, value interface{}) error
		Bool(path string) bool
		BoolF(key string, fallback bool) bool
		String(key string) string
		StringF(key string, fallback string) string
		StringsF(key string, fallback []string) (val []string)
		IntF(key string, fallback int) (val int)
		Float64F(key string, fallback float64) (val float64)
		DurationF(key string, fallback time.Duration) (val time.Duration)
		ByteSizeF(key string, fallback bytesize.ByteSize) bytesize.ByteSize
		GetF(key string, fallback interface{}) (val interface{})
		PrintHumanReadableValidationErrors(w io.Writer, err error)
	}

	ConfigProvider interface {
		Config() DriverConfigure
	}

	DriverConfigure interface {
		Configure
		DSN(ctx context.Context) string
		ServerSHost(ctx context.Context) string
	}
)
