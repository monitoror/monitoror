package config

import coreConfig "github.com/monitoror/monitoror/config"

type (
	Port struct {
		Timeout         int // In Millisecond
		InitialMaxDelay int // In Millisecond
	}
)

var Default = &Port{
	Timeout:         2000,
	InitialMaxDelay: coreConfig.DefaultInitialMaxDelay,
}
