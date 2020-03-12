package models

import (
	"bytes"
	"encoding/json"
)

// The goal here is to raise an error if a key is sent that is not supported.
// This should stop many problems, like misspelling a parameter.
type ConfigWrapper Config // Use wrapper to avoid infinite UnmarshalJSON loop

// UnmarshalJSON should error if there is something unexpected
func (c *Config) UnmarshalJSON(data []byte) error {
	var cw ConfigWrapper
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // Force

	if err := dec.Decode(&cw); err != nil {
		return err
	}
	*c = Config(cw)
	return nil
}
