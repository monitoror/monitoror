package models

import (
	"errors"
	"net"
	"os"
	"strings"
)

type (
	MonitororError struct {
		// Err is the error that occurred during the operation.
		Err error

		// Message used to override Err message
		Message string

		// Tile is used in error handler to return errored tile to request
		Tile *Tile

		// ErrorStatus is used for override current tile Status when error happen
		// Default : ErrorStatus
		ErrorStatus TileStatus
	}
)

var (
	QueryParamsError = &MonitororError{Message: "invalid configuration, unable to parse/check queryParams"}
)

func (e *MonitororError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}
func (e *MonitororError) Unwrap() error { return e.Err }
func (e *MonitororError) Timeout() bool {
	// timeout, host unreachable, deadline exceeded are considered "timeout"
	// it mean we will found previous status in cache to answer

	// Timeout
	if os.IsTimeout(e.Err) {
		return true
	}

	// Host unreachable
	err := e.Err
	for {
		if _, ok := err.(*net.DNSError); ok {
			return true
		}

		if err = errors.Unwrap(err); err == nil {
			break
		}
	}

	// Deadline Exceeded aka context cancellation
	if strings.Contains(e.Err.Error(), "net/http: request canceled while waiting for connection") {
		return true
	}

	return false
}
