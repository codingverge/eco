package axon

import (
	"context"
)

type Axon interface {
	DbalDriver
	LoggingProvider
	ConfigProvider

	WithConfig(c Configurator)
	WithLogger(l Logger)
	RegisterRoutes(ctx context.Context, router *Router)
}
