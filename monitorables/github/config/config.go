package config

import coreConfig "github.com/monitoror/monitoror/config"

type (
	Github struct {
		Timeout              int // In Millisecond
		Token                string
		CountCacheExpiration int // In Millisecond
		InitialMaxDelay      int // In Millisecond
	}
)

var Default = &Github{
	Timeout:              5000,
	Token:                "",
	CountCacheExpiration: 30000,
	InitialMaxDelay:      coreConfig.DefaultInitialMaxDelay,
}
