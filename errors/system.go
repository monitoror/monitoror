package errors

import (
	"fmt"
)

type SystemError struct {
	Message string
	Err     error
}

func NewSystemError(message string) *SystemError {
	return &SystemError{Message: message}
}

func NewSystemErrorWithError(message string, err error) *SystemError {
	return &SystemError{Message: message, Err: err}
}

func (err *SystemError) Error() string {
	if err.Err != nil {
		return fmt.Sprintf("%s, %v", err.Message, err.Err)
	} else {
		return fmt.Sprintf("%s", err.Message)
	}
}
