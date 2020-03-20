package config

import coreConfig "github.com/monitoror/monitoror/config"

type (
	TravisCI struct {
		URL             string
		Timeout         int // In Millisecond
		Token           string
		GithubToken     string
		InitialMaxDelay int // In Millisecond
	}
)

var Default = &TravisCI{
	URL:             "https://api.travis-ci.com/",
	Timeout:         2000,
	Token:           "",
	GithubToken:     "",
	InitialMaxDelay: coreConfig.DefaultInitialMaxDelay,
}
