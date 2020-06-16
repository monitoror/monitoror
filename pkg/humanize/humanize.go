package humanize

import (
	"fmt"
	"reflect"

	"github.com/dustin/go-humanize"
)

//Interface transform interface to string
func Interface(value interface{}) string {
	// TODO handle slice properly if needed
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		rValue := reflect.ValueOf(value)
		if rValue.IsNil() {
			return ""
		}

		value = rValue.Elem()
	}

	switch value := value.(type) {
	case float64:
		return humanize.Ftoa(value)
	default:
		return fmt.Sprintf("%v", value)
	}
}
