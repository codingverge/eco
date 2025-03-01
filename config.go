package axon

import (
	"io"
	"time"

	"github.com/inhies/go-bytesize"
)

type ConfigProvider interface {
	Configurator() Configurator
}

type Configurator interface {
	DirtyPatch(key string, value any) error
	Set(key string, value interface{}) error
	BoolF(key string, fallback bool) bool
	StringF(key string, fallback string) string
	StringsF(key string, fallback []string) (val []string)
	IntF(key string, fallback int) (val int)
	Float64F(key string, fallback float64) (val float64)
	DurationF(key string, fallback time.Duration) (val time.Duration)
	ByteSizeF(key string, fallback bytesize.ByteSize) bytesize.ByteSize
	GetF(key string, fallback interface{}) (val interface{})
	PrintHumanReadableValidationErrors(w io.Writer, err error)
}
