package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/labstack/gommon/log"
)

type (
	Variants map[string]bool
)

func initEnvAndVariant(envPrefix, defaultVariant string, configType reflect.Type) Variants {
	// We need to Identify every Variant
	variants := make(Variants)
	variants[defaultVariant] = true

	// Every env currently set
	envs := os.Environ()

	// List all env that matches prefix
	for _, env := range envs {
		if strings.HasPrefix(env, envPrefix) {
			splitedEnv := strings.Split(env, "=")

			envKey := splitedEnv[0]
			envValue := splitedEnv[1]

			envKeyWithoutPrefix := strings.Replace(envKey, fmt.Sprintf("%s_", envPrefix), "", 1)
			splittedEnvKeyWithoutPrefix := strings.Split(envKeyWithoutPrefix, "_")

			// Look if env without prefix start with struct field or variant
			for i := 0; i < configType.NumField(); i++ {
				if len(splittedEnvKeyWithoutPrefix) > 1 && strings.ToUpper(configType.Field(i).Name) == splittedEnvKeyWithoutPrefix[1] {
					// Env has a variant add it to map
					variants[strings.ToLower(splittedEnvKeyWithoutPrefix[0])] = true
					break
				} else if strings.ToUpper(configType.Field(i).Name) == splittedEnvKeyWithoutPrefix[0] {
					// Env don't have variant, add default
					addDefaultVariant(envKey, fmt.Sprintf("%s_%s_%s", envPrefix, strings.ToUpper(defaultVariant), envKeyWithoutPrefix), envValue)
					break
				}
			}
		}
	}

	return variants
}

func addDefaultVariant(oldEnv, newEnv, value string) {
	if _, exist := os.LookupEnv(newEnv); exist {
		log.Warnf("Env %s can't be used as default, %s already exist", oldEnv, newEnv)
		return
	}

	_ = os.Setenv(newEnv, value)
	_ = os.Unsetenv(oldEnv)
}
