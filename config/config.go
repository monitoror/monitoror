package config

import (
	"strings"

	"github.com/fatih/structs"
	"github.com/spf13/viper"
)

const EnvPrefix = "MO"
const MonitorablePrefix = "MONITORABLE"

const DefaultInitialMaxDelay = 1700

type (
	// Config contain backend Configuration
	Config struct {
		// --- General Configuration ---
		Port int    // Default: 8080
		Env  string // Default: production

		// --- Cache Configuration ---
		// UpstreamCacheExpiration is used to respond before executing the request. Avoid overloading services.
		UpstreamCacheExpiration int
		// DownstreamCacheExpiration is used to respond after executing the request in case of timeout error.
		DownstreamCacheExpiration int
	}
)

var defaultConfig = &Config{
	Port:                      8080,
	Env:                       "production",
	UpstreamCacheExpiration:   10000,
	DownstreamCacheExpiration: 120000,
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
