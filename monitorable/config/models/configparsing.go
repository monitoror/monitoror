package models

import (
	"bytes"
	"encoding/json"
)

// The goal here is to raise an error if a key is sent that is not supported.
// This should stop many problems, like misspelling a parameter.
type TempConfig Config // Use Temp config to avoid infinite UnmarshalJSON loop

// UnmarshalJSON should error if there is something unexpected
func (c *Config) UnmarshalJSON(data []byte) error {
	var tc TempConfig
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields() // Force

	if err := dec.Decode(&tc); err != nil {
		return err
	}
	*c = Config(tc)
	return nil
}
