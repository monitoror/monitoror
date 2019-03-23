package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type (
	Config struct {
		// General Configuration
		Port  int   `json:"port"` // Default: 8080
		Cache Cache `json:"cache"`

		// Gitlab Configuration
		Gitlab GitlabConfig `json:"gitlab"`
	}

	Cache struct {
		Duration        int `json:"duration"`         // In Seconde. Default: 10
		CleanupInterval int `json:"cleanup-Interval"` // In Seconde. Default: 1
	}

	GitlabConfig struct {
		Token string `json:"token,omitempty"`
	}
)

// Load confiuration from configuration file / env / default value
func Load() (*Config, error) {
	var config Config

	// Setup config filename / path
	viper.SetConfigName("config")

	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/monitowall/")

	// Setup Env
	viper.AutomaticEnv()
	viper.SetEnvPrefix("MW")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Setup default values
	viper.SetDefault("Port", 8080)
	viper.SetDefault("Cache.Duration", 10)
	viper.SetDefault("Cache.CleanupInterval", 1)

	// Read Configuration
	_ = viper.ReadInConfig()

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode configuration into struct, %v", err)
	}

	return &config, nil
}
