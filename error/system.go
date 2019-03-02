package error

import (
	"fmt"
)

// TODO Seb : Changing this ^^

type SystemError struct {
	Message    string
	LogMessage string
}

func (err *SystemError) Error() string {
	return fmt.Sprintf("%s", err.Message)
}
