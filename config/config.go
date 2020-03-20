package config

import (
	"fmt"
	"reflect"
	"strings"

	pkgConfig "github.com/monitoror/monitoror/pkg/monitoror/config"

	"github.com/fatih/structs"
	"github.com/spf13/viper"
)

const EnvPrefix = "MO"
const MonitorablePrefix = "MONITORABLE"

const DefaultVariant = "default"
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
	viper.AutomaticEnv()
	viper.SetEnvPrefix(EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Setup default values
	for _, field := range structs.Fields(defaultConfig) {
		viper.SetDefault(field.Name(), field.Value())
	}

	_ = viper.Unmarshal(&config)

	return &config
}

// -------- Config Utility function ---------
func LoadMonitorableConfig(conf interface{}, defaultConf interface{}) {
	pkgConfig.LoadConfigWithVariant(fmt.Sprintf("%s_%s", EnvPrefix, MonitorablePrefix), DefaultVariant, conf, defaultConf)
}

func GetVariantsFromConfig(conf interface{}) []string {
	var variants []string
	if reflect.TypeOf(conf).Kind() == reflect.Map {
		keys := reflect.ValueOf(conf).MapKeys()
		for _, k := range keys {
			variants = append(variants, k.String())
		}
	} else {
		variants = append(variants, DefaultVariant)
	}

	return variants
}
