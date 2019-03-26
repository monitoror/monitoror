package errors

import (
	"fmt"

	. "github.com/jsdidierlaurent/monitowall/renderings"
)

type TimeoutError struct {
	Type    TileType
	Message string //Source of timeout
}

func NewTimeoutError(t TileType, message string) *TimeoutError {
	return &TimeoutError{Type: t, Message: message}
}

func (err *TimeoutError) Error() string {
	return fmt.Sprintf("%s", err.Message)
}
