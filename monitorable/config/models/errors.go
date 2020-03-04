package models

import (
	"fmt"
)

// ConfigFileNotFoundError
type ConfigFileNotFoundError struct {
	PathOrURL string
	Err       error
}

func (e *ConfigFileNotFoundError) Error() string {
	return fmt.Sprintf(`Config not found at: %s, %v`, e.PathOrURL, e.Err.Error())
}
func (e *ConfigFileNotFoundError) Unwrap() error { return e.Err }

// ConfigVersionFormatError
type ConfigVersionFormatError struct {
	WrongVersion string
}

func (e *ConfigVersionFormatError) Error() string {
	return fmt.Sprintf(`json: cannot unmarshal %s into Go struct field Config.ConfigVersion of type string and X.y format`, e.WrongVersion)
}

//ConfigUnmarshalError
type ConfigUnmarshalError struct {
	Err       error
	RawConfig string
}

func (e *ConfigUnmarshalError) Error() string {
	return e.Err.Error()
}
func (e *ConfigUnmarshalError) Unwrap() error { return e.Err }
