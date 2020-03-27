package config

import coreConfig "github.com/monitoror/monitoror/config"

type (
	AzureDevOps struct {
		URL             string
		Timeout         int // In Millisecond
		Token           string
		InitialMaxDelay int // In Millisecond
	}
)

var Default = &AzureDevOps{
	URL:             "",
	Timeout:         4000,
	Token:           "",
	InitialMaxDelay: coreConfig.DefaultInitialMaxDelay,
}
