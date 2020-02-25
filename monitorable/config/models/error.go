package models

import "fmt"

type (
	ConfigNotFoundError struct {
		Err  error
		URL  string
		Path string
	}
)

func (e *ConfigNotFoundError) Error() string {
	if e.URL != "" {
		return fmt.Sprintf("Config not found at URL: %s", e.URL)
	}
	if e.Path != "" {
		return fmt.Sprintf("Config not found at: %s", e.Path)
	}

	return "Config not found"
}

func (e *ConfigNotFoundError) Unwrap() error { return e.Err }
func (e *ConfigNotFoundError) Timeout() bool { return false }
