package monitorable

import (
	"fmt"
	"reflect"
	"strings"

	coreConfig "github.com/monitoror/monitoror/config"
	pkgConfig "github.com/monitoror/monitoror/internal/pkg/monitorable/config"
	"github.com/monitoror/monitoror/models"
	coreModels "github.com/monitoror/monitoror/models"
)

//LoadConfig load config wrapper for monitorable
func LoadConfig(conf interface{}, defaultConf interface{}) {
	pkgConfig.LoadConfigWithVariant(fmt.Sprintf("%s_%s", coreConfig.EnvPrefix, coreConfig.MonitorablePrefix), coreModels.DefaultVariant, conf, defaultConf)
}

//GetVariants extract variants from monitorable config
func GetVariants(conf interface{}) []models.Variant {
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

//GetEnv rebuild Env variable from config variable
//a little dirty, but I don't know how to do better
func GetEnvName(conf interface{}, variant models.Variant, variableName string) string {
	// Verify Params
	if reflect.ValueOf(conf).Kind() != reflect.Ptr {
		panic(fmt.Sprintf("wrong GetConfigVariableEnv parameters: conf need to be a pointer of struct not a %s", reflect.ValueOf(conf).Kind()))
	}

	var env string
	confName := reflect.TypeOf(conf).Elem().Name()
	if variant == models.DefaultVariant {
		env = strings.ToUpper(fmt.Sprintf("%s_%s_%s_%s", coreConfig.EnvPrefix, coreConfig.MonitorablePrefix, confName, variableName))
	} else {
		env = strings.ToUpper(fmt.Sprintf("%s_%s_%s_%s_%s", coreConfig.EnvPrefix, coreConfig.MonitorablePrefix, confName, variant, variableName))
	}

	return env
}
