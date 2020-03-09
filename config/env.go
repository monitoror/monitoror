package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/labstack/gommon/log"
)

type (
	VariantsByMonitorable map[string]Variants
	Variants              map[string]bool
)

func initEnvAndVariant() VariantsByMonitorable {
	// We need to Identify every Variant
	labels := make(VariantsByMonitorable)

	// We need empty Monitorable for identify Monitorable with variant and Env to edit
	monitorableStruct := &Monitorable{}
	monitorableFields := structs.Fields(monitorableStruct)

	// Every env currently set
	envs := os.Environ()

	for _, field := range monitorableFields {
		// Find struct identified as map (can has variant)
		if field.Kind() == reflect.Map {
			labels[field.Name()] = make(Variants)
			labels[field.Name()][DefaultVariant] = true // By default, every config has default value

			// Define EnvPrefix for each env of struct ex: MO_MONITORABLE_JENKINS_
			envPrefix := strings.ToUpper(fmt.Sprintf("%s_%s_%s_", EnvPrefix, structs.Name(monitorableStruct), field.Name()))

			// List all env that matches prefix
			for _, env := range envs {
				if strings.HasPrefix(env, envPrefix) {

					splitedEnv := strings.Split(env, "=")

					envKey := splitedEnv[0]
					envValue := splitedEnv[1]

					envKeyWithoutPrefix := strings.Replace(envKey, envPrefix, "", 1)
					splittedEnvKeyWithoutPrefix := strings.Split(envKeyWithoutPrefix, "_")
					monitorable := reflect.TypeOf(field.Value()).
						Elem(). // Get type of value inside map[string]*XXX (*XXX)
						Elem()  // Get Type pointed by pointer *XXX (XXX)

					// Look if env without prefix start with struct field or label
					needDefault := false
					hasVariant := false
					for i := 0; i < monitorable.NumField(); i++ {
						if strings.ToUpper(monitorable.Field(i).Name) == splittedEnvKeyWithoutPrefix[0] {
							needDefault = true
							break
						} else if len(splittedEnvKeyWithoutPrefix) > 1 && strings.ToUpper(monitorable.Field(i).Name) == splittedEnvKeyWithoutPrefix[1] {
							hasVariant = true
						}
					}

					if needDefault {
						addDefaultVariant(envKey, fmt.Sprintf("%s%s_%s", envPrefix, strings.ToUpper(DefaultVariant), envKeyWithoutPrefix), envValue)
					}

					if hasVariant {
						labels[field.Name()][strings.ToLower(splittedEnvKeyWithoutPrefix[0])] = true
					}
				}
			}
		}
	}

	return labels
}

func addDefaultVariant(oldEnv, newEnv, value string) {
	if _, exist := os.LookupEnv(newEnv); exist {
		log.Warnf("Env %s can't be used as default, %s already exist", oldEnv, newEnv)
		return
	}

	_ = os.Setenv(newEnv, value)
	_ = os.Unsetenv(oldEnv)
}
