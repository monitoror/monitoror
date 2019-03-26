package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type (
	Config struct {
		// --- General Configuration ---
		Port int `json:"port"` // Default: 8080

		// --- Cache Configuration ---
		//UpstreamCache is used to respond before executing the request. Avoid overloading services.
		UpstreamCache Cache `json:"upstream-cache"` // Default: Duration=10 CleanupInterval=1
		//DownstreamCache is used to respond after executing the request in case of timeout error.
		DownstreamCache Cache `json:"downstream-cache"` // Default: Duration=60 CleanupInterval=10

		// Gitlab Configuration
		Gitlab GitlabConfig `json:"gitlab"`
	}

	Cache struct {
		Expire          int `json:"expire"`
		CleanupInterval int `json:"cleanup-Interval"`
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
	viper.SetDefault("UpstreamCache.Duration", 10)
	viper.SetDefault("UpstreamCache.CleanupInterval", 1)
	viper.SetDefault("DownstreamCache.Duration", 120)
	viper.SetDefault("DownstreamCache.CleanupInterval", 10)

	// Read Configuration
	_ = viper.ReadInConfig()

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode configuration into struct, %v", err)
	}

	return &config, nil
}
