package config

import (
	"strings"

	"github.com/fatih/structs"
	"github.com/spf13/viper"
)

const EnvPrefix = "MO"
const MonitorablePrefix = "MONITORABLE"

type (
	// Config contain backend Configuration
	Config struct {
		// --- General Configuration ---
		Port int
		Env  string

		// --- Cache Configuration ---
		// UpstreamCacheExpiration is used to respond before executing the request. Avoid overloading services.
		UpstreamCacheExpiration int
		// DownstreamCacheExpiration is used to respond after executing the request in case of timeout error.
		DownstreamCacheExpiration int

		// InitialMaxDelay is used to add delay on first methode to avoid bursting x requets in same time on start
		InitialMaxDelay int // in Millisecond
	}
)

var defaultConfig = &Config{
	Port:                      8080,
	Env:                       "production",
	UpstreamCacheExpiration:   10000,
	DownstreamCacheExpiration: 120000,
	InitialMaxDelay:           1700,
}

// InitConfig from configuration file / env / default value
func InitConfig() *Config {
	var config Config

	// Setup Env
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix(EnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Setup default values
	for _, field := range structs.Fields(defaultConfig) {
		v.SetDefault(field.Name(), field.Value())
	}

	_ = v.Unmarshal(&config)

	return &config
}
