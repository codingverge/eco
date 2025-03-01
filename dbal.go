package axon

import "context"

// Dbal represents a driver
type Dbal interface {
	// CanHandle returns true if the driver is capable of handling the given DSN or false otherwise.
	CanHandle(dsn string) bool

	// Ping returns nil if the driver has connectivity and is healthy or an error otherwise.
	Ping() error
	PingContext(context.Context) error
}
