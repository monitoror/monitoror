package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/monitoror/monitoror/models"
	pkgConfig "github.com/monitoror/monitoror/pkg/monitoror/config"

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

// -------- Config Utility function ---------
func LoadMonitorableConfig(conf interface{}, defaultConf interface{}) {
	pkgConfig.LoadConfigWithVariant(fmt.Sprintf("%s_%s", EnvPrefix, MonitorablePrefix), models.DefaultVariant, conf, defaultConf)
}

func GetVariantsFromConfig(conf interface{}) []models.Variant {
	var variants []models.Variant
	if reflect.TypeOf(conf).Kind() == reflect.Map {
		keys := reflect.ValueOf(conf).MapKeys()
		for _, k := range keys {
			variants = append(variants, models.Variant(k.String()))
		}
	} else {
		variants = append(variants, models.DefaultVariant)
	}

	return variants
}

func GetEnvFromMonitorableVariable(conf interface{}, variant models.Variant, variableName string) string {
	// Verify Params
	if reflect.ValueOf(conf).Kind() != reflect.Ptr {
		panic(fmt.Sprintf("wrong GetConfigVariableEnv parameters: conf need to be a pointer of struct not a %s", reflect.ValueOf(conf).Kind()))
	}

	var env string
	confName := reflect.TypeOf(conf).Elem().Name()
	if variant == models.DefaultVariant {
		env = strings.ToUpper(fmt.Sprintf("%s_%s_%s_%s", EnvPrefix, MonitorablePrefix, confName, variableName))
	} else {
		env = strings.ToUpper(fmt.Sprintf("%s_%s_%s_%s_%s", EnvPrefix, MonitorablePrefix, confName, variant, variableName))
	}

	return env
}
