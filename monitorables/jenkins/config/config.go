package config

import coreConfig "github.com/monitoror/monitoror/config"

type (
	Jenkins struct {
		URL             string
		Timeout         int // In Millisecond
		SSLVerify       bool
		Login           string
		Token           string
		InitialMaxDelay int // In Millisecond
	}
)

var Default = &Jenkins{
	URL:             "",
	Timeout:         2000,
	SSLVerify:       true,
	Login:           "",
	Token:           "",
	InitialMaxDelay: coreConfig.DefaultInitialMaxDelay,
}
