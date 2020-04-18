package versions

import "fmt"

// ConfigVersionFormatError
type ConfigVersionFormatError struct {
	WrongVersion string
}

func (e *ConfigVersionFormatError) Error() string {
	return fmt.Sprintf(`json: cannot unmarshal %s into Go struct field Config.Version of type string and X.y format`, e.WrongVersion)
}
