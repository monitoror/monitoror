package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/labstack/gommon/log"
)

type (
	Labels map[string]bool
)

// InitEnvDefaultLabel extract Label of monitoror Env variable
// and inject Default label if missing
// Like:
//	 MO_CONFIG=xxx => MO_CONFIG_DEFAULT
//   MO_MONITORABLE_JENKINS_URL=xxx => MO_MONITORABLE_JENKINS_DEFAULT_URL=xxx =>
func InitEnvDefaultLabel(envPrefix string, envSuffix string, defaultLabel string) Labels {
	// We need to Identify every Variant
	variantNames := make(Labels)
	variantNames[defaultLabel] = true

	// Clean params
	envPrefix = strings.TrimSuffix(envPrefix, "_")
	envSuffix = strings.TrimPrefix(envSuffix, "_")
	defaultLabel = strings.ToUpper(defaultLabel)

	// Load env variables
	envs := os.Environ()

	// List all env that matches prefix
	for _, env := range envs {
		if strings.HasPrefix(env, envPrefix) {
			splitedEnv := strings.SplitN(env, "=", 2)

			envKey := splitedEnv[0]
			envValue := splitedEnv[1]

			if strings.HasSuffix(envKey, envSuffix) {
				// extract variant name
				variant := strings.TrimPrefix(envKey, envPrefix)
				variant = strings.TrimSuffix(variant, envSuffix)
				variant = strings.Trim(variant, "_")
				variant = strings.ToLower(variant)

				if variant == "" {
					newEnv := strings.TrimSuffix(fmt.Sprintf("%s_%s_%s", envPrefix, defaultLabel, envSuffix), "_")
					addDefaultLabelToEnv(envKey, newEnv, envValue)
				} else {
					variantNames[variant] = true
				}
			}
		}
	}

	return variantNames
}

func addDefaultLabelToEnv(oldEnv, newEnv, value string) {
	if _, exist := os.LookupEnv(newEnv); exist {
		log.Warnf("Env %s can't be used as default, %s already exist", oldEnv, newEnv)
		return
	}

	_ = os.Setenv(newEnv, value)
	_ = os.Unsetenv(oldEnv)
}
