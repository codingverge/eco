package driver

import (
	"context"
	"github.com/codingverge/axon"
	"github.com/codingverge/axon/config"
)

const (
	KeyDSN                 = "dsn"
	DefaultSQLiteMemoryDSN = "sqlite://file::memory:?_fk=true&cache=shared"
	KeyServerPort          = "server.port"
	KeyServerHost          = "server.host"
)

func NewDriverConfig(ctx context.Context, l axon.Logger, opts ...config.Option) (axon.DriverConfigure, error) {
	c, err := config.New(ctx, []byte(ConfigSchema), opts...)
	if err != nil {
		return nil, err
	}
	l.UseConfig(c)
	return &DriverConfig{
		Config: c,
		l:      l,
	}, nil
}

type (
	DriverConfig struct {
		*config.Config
		l axon.Logger
	}
)

func (c *DriverConfig) ServerSHost(ctx context.Context) string {
	return c.String(KeyServerHost)
}

func (c *DriverConfig) DSN(ctx context.Context) string {
	dsn := c.String(KeyDSN)

	if dsn == "memory" {
		return DefaultSQLiteMemoryDSN
	}

	if len(dsn) > 0 {
		return dsn
	}

	c.l.Errorf("dsn must be set")
	return ""
}
