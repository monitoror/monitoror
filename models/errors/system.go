package errors

import (
	"fmt"
)

type SystemError struct {
	message string
	err     error
}

func NewSystemError(message string, err error) *SystemError {
	return &SystemError{message, err}
}

func (se *SystemError) Error() string {
	if se.err != nil {
		return fmt.Sprintf("%s, %v", se.message, se.err)
	} else {
		return fmt.Sprintf("%s", se.message)
	}
}
