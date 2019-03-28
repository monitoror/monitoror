package errors

import (
	"fmt"

	. "github.com/jsdidierlaurent/monitoror/models/tiles"
)

type TimeoutError struct {
	Tile   *Tile
	Reason string
}

func NewTimeoutError(tile *Tile, reason string) *TimeoutError {
	return &TimeoutError{tile, reason}
}

func (te *TimeoutError) Error() string {
	return fmt.Sprintf("timeout on %s request. %s", te.Tile.SubType, te.Reason)
}
