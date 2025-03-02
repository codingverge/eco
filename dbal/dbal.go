package dbal

import (
	"context"
	"strings"
	"sync"

	"github.com/codingverge/axon"
	"github.com/pkg/errors"
)

var (
	drivers = make([]func() axon.Dbal, 0)
	dmtx    sync.Mutex
	// ErrNoResponsibleDriverFound is returned when no driver was found for the provided DSN.
	ErrNoResponsibleDriverFound = errors.New("dsn value requested an unknown driver")
	ErrSQLiteSupportMissing     = errors.New(`the DSN connection string looks like a SQLite connection, but SQLite support was not built into the binary. Please check if you have downloaded the correct binary or are using the correct Docker Image. Binary archives and Docker Images indicate SQLite support by appending the -sqlite suffix`)
)

// RegisterDriver registers a driver
func RegisterDriver(d func() axon.Dbal) {
	dmtx.Lock()
	drivers = append(drivers, d)
	dmtx.Unlock()
}

// GetDriverFor returns a driver for the given DSN or ErrNoResponsibleDriverFound if no driver was found.
func GetDriverFor(dsn string) (axon.Dbal, error) {
	for _, f := range drivers {
		driver := f()
		if driver.CanHandle(dsn) {
			return driver, nil
		}
	}

	if IsSQLite(dsn) {
		return nil, ErrSQLiteSupportMissing
	}

	return nil, ErrNoResponsibleDriverFound
}

// IsSQLite returns true if the connection is a SQLite string.
func IsSQLite(dsn string) bool {
	scheme := strings.Split(dsn, "://")[0]
	return scheme == "sqlite" || scheme == "sqlite3"
}

func NewDriverFromDbal(ctx context.Context, c axon.DriverConfigure, l axon.Logger) (axon.Driver, error) {
	driver, err := GetDriverFor(c.DSN(ctx))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	d, ok := driver.(axon.Driver)
	if !ok {
		return nil, errors.Errorf("driver of type %T does not implement interface Registry", driver)
	}
	return d.WithLogger(l).WithConfig(c), nil
}
