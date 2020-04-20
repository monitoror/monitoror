package models

import (
	"fmt"
	"reflect"
	"strings"
)

// ConfigFileNotFoundError
type ConfigFileNotFoundError struct {
	PathOrURL string
	Err       error
}

func (e *ConfigFileNotFoundError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf(`Config not found at: %s, %v`, e.PathOrURL, e.Err.Error())
	}
	return fmt.Sprintf(`Config not found at: %s`, e.PathOrURL)
}
func (e *ConfigFileNotFoundError) Unwrap() error { return e.Err }

//ConfigUnmarshalError
type ConfigUnmarshalError struct {
	Err       error
	RawConfig string
}

func (e *ConfigUnmarshalError) Error() string {
	// Hack to hide ConfigWrapper wrapper
	strError := strings.ReplaceAll(e.Err.Error(), reflect.TypeOf(TempConfig{}).Name(), reflect.TypeOf(Config{}).Name())
	return strError
}
func (e *ConfigUnmarshalError) Unwrap() error { return e.Err }
