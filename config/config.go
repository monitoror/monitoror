package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/spf13/viper"

	"github.com/monitoror/monitoror/internal/pkg/env"
)

const (
	EnvPrefix         = "MO"
	MonitorablePrefix = "MONITORABLE"

	DefaultConfigName ConfigName = "default"
)

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

		// InitialMaxDelay is used to add delay on first method to avoid bursting x requests in same time on start
		InitialMaxDelay int // in Millisecond

		// Config can contains ui config (path or url)
		// Can contains default or named config file
		// Like:
		// 		MO_CONFIG=./default.json
		//		MO_CONFIG_SCREEN1=./screen1.json
		//		MO_CONFIG_SCREEN2=https://config.example.com/screen2.json
		//
		// Note: it's the only way to load config file outside of monitoror directory
		Config map[ConfigName]string
	}

	//nolint:golint
	ConfigName string
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

	// Setup default named config
	for _, field := range structs.Fields(&config) {
		if reflect.DeepEqual(field.Value(), config.Config) {
			envPrefix := strings.ToUpper(fmt.Sprintf("%s_%s", EnvPrefix, field.Name()))
			env.InitEnvDefaultLabel(envPrefix, "", string(DefaultConfigName))
			break
		}
	}

	// Setup default values
	for _, field := range structs.Fields(defaultConfig) {
		v.SetDefault(field.Name(), field.Value())
	}

	_ = v.Unmarshal(&config)

	return &config
}
