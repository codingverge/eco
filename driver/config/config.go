package config

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

func New(ctx context.Context) (axon.DriverConfigure, error) {
	c, err := config.New(ctx, []byte(ConfigSchema))
	if err != nil {
		return nil, err
	}
	return &DriverConfig{
		Config: c,
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
