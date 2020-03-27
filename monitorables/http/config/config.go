package config

import coreConfig "github.com/monitoror/monitoror/config"

type (
	HTTP struct {
		Timeout         int // In Millisecond
		SSLVerify       bool
		InitialMaxDelay int // In Millisecond
	}
)

var Default = &HTTP{
	Timeout:         2000,
	SSLVerify:       true,
	InitialMaxDelay: coreConfig.DefaultInitialMaxDelay,
}
