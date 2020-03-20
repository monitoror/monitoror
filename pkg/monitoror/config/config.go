package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"

	"github.com/spf13/viper"
)

func LoadConfigWithVariant(envPrefix, defaultVariant string, conf interface{}, defaultConf interface{}) {
	// Verify Params
	if reflect.ValueOf(conf).Kind() != reflect.Ptr {
		panic(fmt.Sprintf("wrong LoadConfigWithVariant parameters: conf need to be a pointer of map[string]struct not a %s", reflect.ValueOf(conf).Kind()))
	}
	if reflect.ValueOf(conf).Elem().Kind() != reflect.Map || reflect.ValueOf(conf).Elem().Type().Key().Kind() != reflect.String {
		panic(fmt.Sprintf("wrong LoadConfigWithVariant parameters: conf need to be a pointer of map[string]struct not a %s", reflect.ValueOf(conf).Elem().Type()))
	}
	if reflect.ValueOf(defaultConf).Kind() != reflect.Ptr {
		panic(fmt.Sprintf("wrong LoadConfigWithVariant parameters: defaultConf need to be a pointer of struct not a %s", reflect.ValueOf(conf).Kind()))
	}

	// Unbox defaultConf
	unboxedDefaultConfig := reflect.ValueOf(defaultConf).Elem()
	unboxedDefaultConfigType := unboxedDefaultConfig.Type()

	// Add Config struct name to prefix
	envPrefix = strings.ToUpper(fmt.Sprintf("%s_%s", envPrefix, unboxedDefaultConfigType.Name()))

	// Setup Env
	viper.AutomaticEnv()
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Transform Env and define Label for setting default value
	variants := initEnvAndVariant(envPrefix, defaultVariant, unboxedDefaultConfigType)

	// Setup default value
	for variant := range variants {
		for _, field := range structs.Fields(defaultConf) {
			viper.SetDefault(fmt.Sprintf("%s.%s", variant, field.Name()), field.Value())
		}
	}

	_ = viper.Unmarshal(&conf)
}
