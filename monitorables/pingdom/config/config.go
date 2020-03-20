package config

import coreConfig "github.com/monitoror/monitoror/config"

type (
	Pingdom struct {
		URL             string
		Token           string
		Timeout         int // In Millisecond
		CacheExpiration int // In Millisecond
		InitialMaxDelay int // In Millisecond
	}
)

var Default = &Pingdom{
	URL:             "https://api.pingdom.com/api/3.1",
	Token:           "",
	Timeout:         2000,
	CacheExpiration: 30000,
	InitialMaxDelay: coreConfig.DefaultInitialMaxDelay,
}
