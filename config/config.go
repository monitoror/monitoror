package config

import (
	"strings"

	"github.com/jsdidierlaurent/monitoror/models/errors"

	"github.com/spf13/viper"
)

const ServerConfigFileName = "server-config"
const EnvPrefix = "MO"

type (
	Config struct {
		// --- General Configuration ---
		Port int    `json:"port"` // Default: 8080
		Mode string `json:"mode"` // Default: production

		// --- Cache Configuration ---
		//UpstreamCache is used to respond before executing the request. Avoid overloading services.
		UpstreamCache Cache `json:"upstream-cache"`
		//DownstreamCache is used to respond after executing the request in case of timeout error.
		DownstreamCache Cache `json:"downstream-cache"`

		// --- Ping Configuration ---
		PingConfig PingConfig `json:"ping-config"`

		// --- Port Configuration ---
		PortConfig PortConfig `json:"port-config"`

		// --- Gitlab Configuration ---
		GitlabConfig GitlabConfig `json:"gitlab-config"`
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

	PortConfig struct {
		Timeout int `json:"timeout"` // In Millisecond
	}

	GitlabConfig struct {
		Token string `json:"token,omitempty"`
	}
)

// Load confiuration from configuration file / env / default value
func InitConfig() (*Config, error) {
	var config Config

	// Setup config filename / path
	viper.SetConfigName(ServerConfigFileName)

	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/monitoror/")

	// Setup Env
	viper.AutomaticEnv()
	viper.SetEnvPrefix(EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Setup default values
	// --- General Configuration ---
	viper.SetDefault("Port", 8080)
	viper.SetDefault("Mode", "production")

	// --- Cache Configuration ---
	viper.SetDefault("UpstreamCache.Expire", 10000)
	viper.SetDefault("UpstreamCache.CleanupInterval", 1000)
	viper.SetDefault("DownstreamCache.Expire", 120000)
	viper.SetDefault("DownstreamCache.CleanupInterval", 10000)

	// --- Ping Configuration ---
	viper.SetDefault("PingConfig.Count", 2)
	viper.SetDefault("PingConfig.Timeout", 1000)
	viper.SetDefault("PingConfig.Interval", 100)

	// --- Port Configuration ---
	viper.SetDefault("PortConfig.Timeout", 1000)

	// Read Configuration
	err := viper.ReadInConfig()
	if _, ok := err.(viper.ConfigParseError); ok {
		return nil, errors.NewConfigError(err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, errors.NewConfigError(err)
	}

	return &config, nil
}
