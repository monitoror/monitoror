package config

import (
	"strings"

	"github.com/jsdidierlaurent/monitowall/models/errors"

	"github.com/spf13/viper"
)

type (
	Config struct {
		// --- General Configuration ---
		Port int `json:"port"` // Default: 8080

		// --- Cache Configuration ---
		//UpstreamCache is used to respond before executing the request. Avoid overloading services.
		UpstreamCache Cache `json:"upstream-cache"`
		//DownstreamCache is used to respond after executing the request in case of timeout error.
		DownstreamCache Cache `json:"downstream-cache"`

		// --- Ping Configuration ---
		Ping PingConfig `json:"ping"`

		// --- Gitlab Configuration ---
		Gitlab GitlabConfig `json:"gitlab"`
	}

	Cache struct {
		Expire          int `json:"expire"`           // In Millisecond
		CleanupInterval int `json:"cleanup-interval"` // In Millisecond
	}

	PingConfig struct {
		Count    int `json:"count"`
		Timeout  int `json:"timeout"`  // In Millisecond
		Interval int `json:"interval"` // In Millisecond
	}

	GitlabConfig struct {
		Token string `json:"token,omitempty"`
	}
)

// Load confiuration from configuration file / env / default value
func InitConfig() (*Config, error) {
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
	// --- General Configuration ---
	viper.SetDefault("Port", 8080)

	// --- Cache Configuration ---
	viper.SetDefault("UpstreamCache.Expire", 10000)
	viper.SetDefault("UpstreamCache.CleanupInterval", 1000)
	viper.SetDefault("DownstreamCache.Expire", 120000)
	viper.SetDefault("DownstreamCache.CleanupInterval", 10000)

	// --- Ping Configuration ---
	viper.SetDefault("Ping.Count", 2)
	viper.SetDefault("Ping.Timeout", 1000)
	viper.SetDefault("Ping.Interval", 100)

	// Read Configuration
	_ = viper.ReadInConfig()

	if err := viper.Unmarshal(&config); err != nil {
		return nil, errors.NewConfigError(err)
	}

	return &config, nil
}
