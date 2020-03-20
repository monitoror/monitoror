package config

import coreConfig "github.com/monitoror/monitoror/config"

type (
	Ping struct {
		Count           int
		Timeout         int // In Millisecond
		Interval        int // In Millisecond
		InitialMaxDelay int // In Millisecond
	}
)

var Default = &Ping{
	Count:           2,
	Timeout:         1000,
	Interval:        100,
	InitialMaxDelay: coreConfig.DefaultInitialMaxDelay,
}
