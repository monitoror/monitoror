package debug

import (
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

// Enable sets the viper debug var to true
// and makes the logger to log at debug level.
func Enable() {
	viper.Set("debug", true)
	log.SetLevel(log.DEBUG)
}

// Disable sets the viper debug var to false
// and makes the logger to log at info level.
func Disable() {
	viper.Set("debug", false)
	log.SetLevel(log.INFO)
}

// IsEnabled checks whether the debug flag is set or not.
func IsEnabled() bool {
	return viper.GetBool("debug")
}
