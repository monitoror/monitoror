package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/structs"
	"github.com/spf13/viper"

	"github.com/monitoror/monitoror/internal/pkg/env"
)

const (
	EnvPrefix         = "MO"
	MonitorablePrefix = "MONITORABLE"
	ConfigPrefix      = "CONFIG"

	DefaultConfigName ConfigName = "default"
)

type (
	// NamedConfig contain backend Configuration
	CoreConfig struct {
		// --- General Configuration ---
		Port      int
		Address   string
		DisableUI bool

		// --- Cache Configuration ---
		// UpstreamCacheExpiration is used to respond before executing the request. Avoid overloading services.
		UpstreamCacheExpiration int
		// DownstreamCacheExpiration is used to respond after executing the request in case of timeout error.
		DownstreamCacheExpiration int

		// InitialMaxDelay is used to add delay on first method to avoid bursting x requests in same time on start
		InitialMaxDelay int // in Millisecond

		// NamedConfig can contains ui config (path or url)
		// Can contains default or named config file
		// Like:
		// 		MO_CONFIG=./default.json
		//		MO_CONFIG_SCREEN1=./screen1.json
		//		MO_CONFIG_SCREEN2=https://config.example.com/screen2.json
		//
		// Note: it's the only way to load config file outside of monitoror directory
		NamedConfigs map[ConfigName]string
	}

	//nolint:golint
	ConfigName string
)

var defaultConfig = &CoreConfig{
	Port:                      8080,
	Address:                   "0.0.0.0",
	DisableUI:                 false,
	UpstreamCacheExpiration:   10000,
	DownstreamCacheExpiration: 120000,
	InitialMaxDelay:           1700,
}

// InitConfig from configuration file / env / default value
func InitConfig() *CoreConfig {
	coreConfig := &CoreConfig{}

	// Setup Env
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix(EnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Setup default values
	for _, field := range structs.Fields(defaultConfig) {
		v.SetDefault(field.Name(), field.Value())
	}

	_ = v.Unmarshal(coreConfig)

	// Setup NamedConfig without viper
	loadNamedConfig(coreConfig)

	return coreConfig
}

// loadUiConfig load NamedConfig
// Note: it's to "hacky" and complicated with viper so i do it manually
func loadNamedConfig(config *CoreConfig) {
	config.NamedConfigs = make(map[ConfigName]string)

	// Setup default named config
	envPrefix := strings.ToUpper(fmt.Sprintf("%s_%s", EnvPrefix, ConfigPrefix))
	env.InitEnvDefaultLabel(envPrefix, "", string(DefaultConfigName))

	for _, env := range os.Environ() {
		if strings.HasPrefix(env, envPrefix) {
			splittedEnv := strings.SplitN(env, "=", 2)

			configName := strings.TrimPrefix(splittedEnv[0], envPrefix)
			configName = strings.Trim(configName, "_")
			configName = strings.ToLower(configName)
			value := splittedEnv[1]

			config.NamedConfigs[ConfigName(configName)] = value
		}
	}
}
