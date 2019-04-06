package config

import (
	"strings"

	"github.com/spf13/viper"
)

const EnvPrefix = "MO"

type (
	Config struct {
		// --- General Configuration ---
		Port int    // Default: 8080
		Env  string // Default: production

		// --- Cache Configuration ---
		//UpstreamCache is used to respond before executing the request. Avoid overloading services.
		UpstreamCache Cache
		//DownstreamCache is used to respond after executing the request in case of timeout error.
		DownstreamCache Cache

		// --- Ping Configuration ---
		PingConfig PingConfig

		// --- Port Configuration ---
		PortConfig PortConfig

		// --- Gitlab Configuration ---
		GitlabConfig GitlabConfig
	}

	Cache struct {
		Expire          int // In Millisecond
		CleanupInterval int // In Millisecond
	}

	PingConfig struct {
		Count    int
		Timeout  int // In Millisecond
		Interval int // In Millisecond
	}

	PortConfig struct {
		Timeout int // In Millisecond
	}

	GitlabConfig struct {
		Token string
	}
)

// Load confiuration from configuration file / env / default value
func InitConfig() (*Config, error) {
	var config Config

	// Setup Env
	viper.AutomaticEnv()
	viper.SetEnvPrefix(EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Setup default values
	// --- General Configuration ---
	viper.SetDefault("Port", 8080)
	viper.SetDefault("Env", "production")

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
	_ = viper.Unmarshal(&config)

	return &config, nil
}
