package humanize

import (
	"fmt"

	"github.com/dustin/go-humanize"
)

//Interface transform interface to string
func Interface(value interface{}) string {
	// TODO handle slice properly if needed

	switch value := value.(type) {
	case float64:
		return humanize.Ftoa(value)
	default:
		return fmt.Sprintf("%v", value)
	}
}
