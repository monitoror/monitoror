package monitorable

import (
	"fmt"
	"reflect"
	"strings"

	coreConfig "github.com/monitoror/monitoror/config"
	pkgConfig "github.com/monitoror/monitoror/internal/pkg/monitorable/config"
	"github.com/monitoror/monitoror/internal/pkg/validator/validate"
	"github.com/monitoror/monitoror/models"
	coreModels "github.com/monitoror/monitoror/models"
)

// LoadConfig load config wrapper for monitorable
func LoadConfig(conf interface{}, defaultConf interface{}) {
	pkgConfig.LoadConfigWithVariant(fmt.Sprintf("%s_%s", coreConfig.EnvPrefix, coreConfig.MonitorablePrefix), coreModels.DefaultVariant, conf, defaultConf)
}

func ValidateConfig(conf interface{}, variantName coreModels.VariantName) []error {
	var result []error

	if errors := validate.Struct(conf); len(errors) > 0 {
		for _, err := range errors {
			// Replace fieldName by env variable
			err.SetFieldName(buildMonitorableEnvKey(conf, variantName, strings.ToUpper(err.GetFieldName())))

			result = append(result, err)
		}
	}

	return result
}

// buildMonitorableEnvKey rebuild Env variable from config variable
// a little dirty, but I don't know how to do better
func buildMonitorableEnvKey(conf interface{}, variantName models.VariantName, variableName string) string {
	// Verify Params
	if reflect.ValueOf(conf).Kind() != reflect.Ptr {
		panic(fmt.Sprintf("wrong GetConfigVariableEnv parameters: conf need to be a pointer of struct not a %s", reflect.ValueOf(conf).Kind()))
	}

	var env string
	confName := reflect.TypeOf(conf).Elem().Name()
	if variantName == models.DefaultVariant {
		env = strings.ToUpper(fmt.Sprintf("%s_%s_%s_%s", coreConfig.EnvPrefix, coreConfig.MonitorablePrefix, confName, variableName))
	} else {
		env = strings.ToUpper(fmt.Sprintf("%s_%s_%s_%s_%s", coreConfig.EnvPrefix, coreConfig.MonitorablePrefix, confName, variantName, variableName))
	}

	return env
}

// GetVariantsNames extract variants from monitorable config
func GetVariantsNames(conf interface{}) []models.VariantName {
	// Verify Params
	if reflect.ValueOf(conf).Kind() != reflect.Map {
		panic(fmt.Sprintf("wrong GetVariantsNames parameters: conf need to be a map[coreModels.VariantName] not a %s", reflect.ValueOf(conf).Kind()))
	}

	var variants []models.VariantName
	keys := reflect.ValueOf(conf).MapKeys()
	for _, k := range keys {
		variants = append(variants, models.VariantName(k.String()))
	}

	return variants
}
