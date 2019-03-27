package errors

import (
	"fmt"
)

type ConfigError struct {
	err error
}

func NewConfigError(err error) *ConfigError {
	return &ConfigError{err}
}

func (ce *ConfigError) Error() string {
	return fmt.Sprintf("unable to init config, %v", ce.err)
}
